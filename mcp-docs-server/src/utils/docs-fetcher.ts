import { promises as fs } from 'fs';
import path from 'path';
import yaml from 'yaml';
import axios from 'axios';
import os from 'os';
import { VectorSearch } from './vector-search.js';

export interface DocPage {
  title: string;
  path: string;
  icon?: string;
  content?: string;
  url: string;
  category: string;
  section: string;
  subsection?: string;
  level: number;
}

export interface DocStructure {
  pages: DocPage[];
  sections: Map<string, DocPage[]>;
  categories: Map<string, DocPage[]>;
  lastUpdated: Date;
  version: string;
}

interface CacheEntry {
  data: DocStructure;
  timestamp: number;
  etag?: string;
}

export class DocsFetcher {
  private docsCache: DocStructure | null = null;
  private readonly GITHUB_RAW_BASE = 'https://raw.githubusercontent.com/VapiAI/docs/main';
  private readonly DOCS_YML_URL = 'https://raw.githubusercontent.com/VapiAI/docs/main/fern/docs.yml';
  private readonly CACHE_TTL = 1000 * 60 * 60; // 1 hour
  private readonly CACHE_FILE_PATH = path.join(os.tmpdir(), 'vapi-docs-cache.json');
  private backgroundRefreshPromise: Promise<void> | null = null;
  private vectorSearch: VectorSearch = new VectorSearch();
  private vectorIndexingPromise: Promise<void> | null = null;

  async getDocumentationStructure(): Promise<DocStructure> {
    // Check memory cache first
    if (this.docsCache && Date.now() - this.docsCache.lastUpdated.getTime() < this.CACHE_TTL) {
      return this.docsCache;
    }

    // Try to load from persistent cache
    const cachedData = await this.loadFromDiskCache();
    if (cachedData && Date.now() - cachedData.timestamp < this.CACHE_TTL) {
      this.docsCache = cachedData.data;
      console.log('📦 Loaded documentation from disk cache');
      
      // Start background vector indexing if not already done
      this.ensureVectorIndexing();
      
      // Start background refresh if cache is getting old (> 30 minutes)
      const cacheAge = Date.now() - cachedData.timestamp;
      if (cacheAge > 30 * 60 * 1000 && !this.backgroundRefreshPromise) {
        this.backgroundRefreshPromise = this.backgroundRefresh();
      }
      
      return this.docsCache;
    }

    // Fetch fresh data
    console.log('🔄 Fetching fresh documentation from Vapi...');
    return await this.fetchFreshDocumentation();
  }

  private async fetchFreshDocumentation(): Promise<DocStructure> {
    try {
      // Fetch docs.yml structure from the official source
      const response = await axios.get(this.DOCS_YML_URL, {
        headers: {
          'User-Agent': 'Vapi-MCP-Docs-Server/1.0'
        }
      });
      
      const docsConfig = yaml.parse(response.data);
      
      // Parse the navigation structure
      const pages = await this.parseNavigationStructure(docsConfig.navigation);
      
      // Organize into sections and categories
      const sections = new Map<string, DocPage[]>();
      const categories = new Map<string, DocPage[]>();
      
      pages.forEach((page: DocPage) => {
        // Group by section
        if (!sections.has(page.section)) {
          sections.set(page.section, []);
        }
        sections.get(page.section)!.push(page);
        
        // Group by category
        if (!categories.has(page.category)) {
          categories.set(page.category, []);
        }
        categories.get(page.category)!.push(page);
      });

      const docStructure: DocStructure = {
        pages,
        sections,
        categories,
        lastUpdated: new Date(),
        version: 'live'
      };

      // Save to memory and disk cache
      this.docsCache = docStructure;
      await this.saveToDiskCache({
        data: docStructure,
        timestamp: Date.now(),
        etag: response.headers.etag
      });

      // Start vector indexing in background
      this.ensureVectorIndexing();

      console.log(`✅ Fetched ${pages.length} documentation pages`);
      return docStructure;
      
    } catch (error) {
      console.error('❌ Failed to fetch docs structure:', error);
      
      // Try to return stale cache as fallback
      const staleCache = await this.loadFromDiskCache();
      if (staleCache) {
        console.log('⚠️  Using stale cache as fallback');
        this.docsCache = staleCache.data;
        return staleCache.data;
      }
      
      throw new Error(`Failed to fetch docs structure: ${error}`);
    }
  }

  private ensureVectorIndexing(): void {
    if (this.vectorIndexingPromise || !this.docsCache) {
      return;
    }

    this.vectorIndexingPromise = this.performVectorIndexing();
  }

  private async performVectorIndexing(): Promise<void> {
    try {
      if (!this.docsCache) return;

      console.log('🧠 Starting vector indexing...');
      await this.vectorSearch.initialize();
      
      // Only reindex if we don't have embeddings or docs have changed
      if (!this.vectorSearch.isReady()) {
        await this.vectorSearch.indexDocuments(this.docsCache.pages);
      }
      
      console.log('✅ Vector indexing completed');
    } catch (error) {
      console.error('❌ Vector indexing failed:', error);
    } finally {
      this.vectorIndexingPromise = null;
    }
  }

  private async backgroundRefresh(): Promise<void> {
    try {
      console.log('🔄 Background refresh started...');
      await this.fetchFreshDocumentation();
      console.log('✅ Background refresh completed');
    } catch (error) {
      console.error('❌ Background refresh failed:', error);
    } finally {
      this.backgroundRefreshPromise = null;
    }
  }

  private async loadFromDiskCache(): Promise<CacheEntry | null> {
    try {
      const cacheData = await fs.readFile(this.CACHE_FILE_PATH, 'utf-8');
      const parsed = JSON.parse(cacheData) as CacheEntry;
      
      // Reconstruct Maps from plain objects
      parsed.data.sections = new Map(Object.entries(parsed.data.sections || {}));
      parsed.data.categories = new Map(Object.entries(parsed.data.categories || {}));
      parsed.data.lastUpdated = new Date(parsed.data.lastUpdated);
      
      return parsed;
    } catch (error) {
      // Cache file doesn't exist or is corrupted
      return null;
    }
  }

  private async saveToDiskCache(cacheEntry: CacheEntry): Promise<void> {
    try {
      // Convert Maps to plain objects for JSON serialization
      const serializable = {
        ...cacheEntry,
        data: {
          ...cacheEntry.data,
          sections: Object.fromEntries(cacheEntry.data.sections),
          categories: Object.fromEntries(cacheEntry.data.categories)
        }
      };
      
      await fs.writeFile(this.CACHE_FILE_PATH, JSON.stringify(serializable, null, 2));
    } catch (error) {
      console.warn('⚠️  Failed to save cache to disk:', error);
    }
  }

  async invalidateCache(): Promise<void> {
    this.docsCache = null;
    await this.vectorSearch.invalidateIndex();
    try {
      await fs.unlink(this.CACHE_FILE_PATH);
      console.log('🗑️  Cache invalidated');
    } catch (error) {
      // Cache file doesn't exist, that's fine
    }
  }

  private async parseNavigationStructure(navigation: any[]): Promise<DocPage[]> {
    const pages: DocPage[] = [];
    
    for (const tab of navigation) {
      if (tab.tab === 'documentation' && tab.layout) {
        await this.parseLayoutSection(tab.layout, pages, tab.tab, '', 0);
      }
    }
    
    return pages;
  }

  private async parseLayoutSection(
    layout: any[], 
    pages: DocPage[], 
    category: string, 
    parentSection: string, 
    level: number
  ): Promise<void> {
    for (const item of layout) {
      if (item.section) {
        // This is a section
        const sectionName = item.section;
        const currentSection = parentSection ? `${parentSection} > ${sectionName}` : sectionName;
        
        if (item.contents) {
          await this.parseLayoutSection(item.contents, pages, category, currentSection, level + 1);
        }
      } else if (item.page && item.path) {
        // This is a page
        const docPage: DocPage = {
          title: item.page,
          path: item.path,
          icon: item.icon,
          url: `https://docs.vapi.ai/${item.path.replace('.mdx', '')}`,
          category,
          section: parentSection || 'General',
          level
        };
        
        pages.push(docPage);
      } else if (item.link) {
        // This is an external link
        const docPage: DocPage = {
          title: item.link,
          path: '',
          url: item.href,
          category,
          section: parentSection || 'External',
          level
        };
        
        pages.push(docPage);
      }
    }
  }

  async fetchPageContent(docPage: DocPage): Promise<string> {
    if (!docPage.path) {
      return `# ${docPage.title}\n\nExternal link: ${docPage.url}`;
    }

    try {
      // Try to fetch from GitHub raw content
      const contentUrl = `${this.GITHUB_RAW_BASE}/fern/docs/${docPage.path}`;
      const response = await axios.get(contentUrl);
      
      // Parse MDX content and extract meaningful text
      let content = this.parseMDXContent(response.data);
      
      // If content is too short, provide helpful summary instead
      if (content.length < 100) {
        content = this.generatePageSummary(docPage);
      }
      
      // Cache the content
      docPage.content = content;
      
      return content;
      
    } catch (error) {
      // Return helpful summary instead of error
      return this.generatePageSummary(docPage);
    }
  }

  private parseMDXContent(mdxContent: string): string {
    // Remove frontmatter
    const content = mdxContent.replace(/^---[\s\S]*?---\n/, '');
    
    // Remove import statements
    const cleanContent = content.replace(/^import.*$/gm, '');
    
    // Remove JSX components but keep their content
    const withoutJSX = cleanContent
      .replace(/<[^>]*>/g, '')
      .replace(/\{[^}]*\}/g, '')
      .replace(/^\s*$/gm, '')
      .trim();
    
    return withoutJSX;
  }

  private generatePageSummary(docPage: DocPage): string {
    const { title, url, section, category } = docPage;
    
    let summary = `# ${title}\n\n`;
    
    // Add contextual information based on category
    switch (category) {
      case 'quickstart':
        summary += `This is a getting started guide for ${title.toLowerCase()}. `;
        summary += `It covers the basic setup and implementation steps for building voice AI applications with Vapi.\n\n`;
        break;
      case 'assistants':
        summary += `This guide covers ${title.toLowerCase()} for voice assistants. `;
        summary += `Learn how to configure and customize your Vapi voice agents.\n\n`;
        break;
      case 'workflows':
        summary += `This section explains ${title.toLowerCase()} in Vapi's visual workflow system. `;
        summary += `Build complex conversation flows with branching logic.\n\n`;
        break;
      case 'guides':
        summary += `This is a practical guide showing ${title.toLowerCase()}. `;
        summary += `Includes code examples and step-by-step instructions.\n\n`;
        break;
      default:
        summary += `This documentation covers ${title.toLowerCase()} in Vapi.\n\n`;
    }
    
    summary += `**Section:** ${section}\n`;
    summary += `**Category:** ${category}\n\n`;
    
    // Add relevant getting started tips
    if (category === 'quickstart') {
      summary += `## Quick Start Tips:\n\n`;
      summary += `1. Sign up for a Vapi account at https://dashboard.vapi.ai\n`;
      summary += `2. Get your API key from the dashboard\n`;
      summary += `3. Install the Vapi SDK for your preferred language\n`;
      summary += `4. Follow the examples in the complete documentation\n\n`;
    }
    
    summary += `## 📖 Complete Documentation\n\n`;
    summary += `For detailed instructions, code examples, and interactive content, visit:\n`;
    summary += `**${url}**\n\n`;
    
    summary += `This official documentation includes:\n`;
    summary += `- Step-by-step tutorials\n`;
    summary += `- Complete code examples\n`;
    summary += `- API reference details\n`;
    summary += `- Troubleshooting guides\n`;
    summary += `- Interactive demos\n\n`;
    
    summary += `💡 **Pro Tip:** Use the CLI to generate code templates:\n`;
    summary += `\`\`\`bash\nvapi init  # Set up a new project\nvapi assistant create  # Create an assistant\n\`\`\``;
    
    return summary;
  }

  async searchDocumentation(query: string, category?: string): Promise<{results: DocPage[], usedVectorSearch: boolean}> {
    const docs = await this.getDocumentationStructure();
    
    // Wait for vector indexing to complete if it's in progress
    if (this.vectorIndexingPromise) {
      console.log('⏳ Waiting for vector indexing to complete...');
      await this.vectorIndexingPromise;
    }
    
    // Try vector search first if available
    if (this.vectorSearch.isReady()) {
      try {
        const vectorResults = await this.vectorSearch.search(query, 10, 0.2); // Lower threshold for better matches
        
        // Filter by category if specified
        let filteredResults = vectorResults;
        if (category && category !== 'all') {
          filteredResults = vectorResults.filter(page => page.category === category);
        }
        
        if (filteredResults.length > 0) {
          console.log(`🧠 Vector search returned ${filteredResults.length} results`);
          return {results: filteredResults, usedVectorSearch: true};
        } else {
          console.log(`🔍 Vector search found no matches above threshold, falling back to text search`);
        }
      } catch (error) {
        console.warn('⚠️  Vector search failed, falling back to text search:', error);
      }
    } else {
      console.log(`📝 Vector search not ready (${this.vectorSearch.getIndexSize()} embeddings), using text search`);
    }

    // Fallback to text search
    const textResults = this.textSearchDocumentation(query, category, docs);
    return {results: textResults, usedVectorSearch: false};
  }

  // Keep the old method for backward compatibility, but mark it as deprecated
  async searchDocumentationLegacy(query: string, category?: string): Promise<DocPage[]> {
    const result = await this.searchDocumentation(query, category);
    return result.results;
  }

  private textSearchDocumentation(query: string, category: string | undefined, docs: DocStructure): DocPage[] {
    const searchTerm = query.toLowerCase();
    
    let searchPages = docs.pages;
    
    // Filter by category if specified
    if (category && category !== 'all') {
      searchPages = docs.categories.get(category) || [];
    }
    
    // Search in title, section, and content
    const results = searchPages.filter(page => {
      const titleMatch = page.title.toLowerCase().includes(searchTerm);
      const sectionMatch = page.section.toLowerCase().includes(searchTerm);
      const urlMatch = page.url.toLowerCase().includes(searchTerm);
      
      return titleMatch || sectionMatch || urlMatch;
    });
    
    // Sort by relevance (title matches first)
    return results.sort((a, b) => {
      const aTitle = a.title.toLowerCase().includes(searchTerm) ? 1 : 0;
      const bTitle = b.title.toLowerCase().includes(searchTerm) ? 1 : 0;
      return bTitle - aTitle;
    });
  }

  async getExamples(): Promise<DocPage[]> {
    const docs = await this.getDocumentationStructure();
    return docs.pages.filter(page => 
      page.section.toLowerCase().includes('example') ||
      page.category.toLowerCase().includes('example') ||
      page.title.toLowerCase().includes('example')
    );
  }

  async getGuides(): Promise<DocPage[]> {
    const docs = await this.getDocumentationStructure();
    return docs.pages.filter(page => 
      page.section.toLowerCase().includes('guide') ||
      page.section.toLowerCase().includes('quickstart') ||
      page.section.toLowerCase().includes('get started') ||
      page.section.toLowerCase().includes('tutorial')
    );
  }

  async getApiReference(): Promise<DocPage[]> {
    const docs = await this.getDocumentationStructure();
    return docs.pages.filter(page => 
      page.category.toLowerCase().includes('api') ||
      page.section.toLowerCase().includes('api') ||
      page.title.toLowerCase().includes('api')
    );
  }

  async getBestPractices(): Promise<DocPage[]> {
    const docs = await this.getDocumentationStructure();
    return docs.pages.filter(page => 
      page.section.toLowerCase().includes('best practices') ||
      page.section.toLowerCase().includes('prompting') ||
      page.section.toLowerCase().includes('debugging') ||
      page.section.toLowerCase().includes('testing')
    );
  }

  isVectorSearchReady(): boolean {
    return this.vectorSearch.isReady();
  }

  getVectorIndexSize(): number {
    return this.vectorSearch.getIndexSize();
  }
} 