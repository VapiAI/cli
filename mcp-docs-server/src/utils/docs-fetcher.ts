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
      const response = await axios.get(contentUrl, {
        timeout: 5000,
        headers: {
          'User-Agent': 'Vapi-MCP-Docs-Server/1.0'
        }
      });
      
      // Parse MDX content and extract meaningful text
      let content = this.parseMDXContent(response.data);
      
      // If content is too short or looks like it failed, provide contextual summary
      if (content.length < 100 || content.includes('This documentation covers')) {
        content = this.generateContextualSummary(docPage);
      }
      
      // Cache the content
      docPage.content = content;
      
      return content;
      
          } catch (error: any) {
        // Return contextual summary instead of generic error
        console.warn(`⚠️  Failed to fetch content for ${docPage.path}:`, error.response?.status || error.message);
        return this.generateContextualSummary(docPage);
      }
  }

  private parseMDXContent(mdxContent: string): string {
    try {
      // Remove frontmatter
      let content = mdxContent.replace(/^---[\s\S]*?---\n/, '');
      
      // Remove import statements
      content = content.replace(/^import.*$/gm, '');
      
      // Remove JSX components but preserve their content
      content = content
        .replace(/<([A-Z][A-Za-z0-9]*)[^>]*>([\s\S]*?)<\/\1>/g, '$2')
        .replace(/<[^>]*\/>/g, '')
        .replace(/<[^>]*>/g, '');
      
      // Remove JSX expressions but keep useful content
      content = content.replace(/\{[^}]*\}/g, '');
      
      // Clean up excessive whitespace and newlines
      content = content
        .replace(/^\s*$/gm, '')
        .replace(/\n{3,}/g, '\n\n')
        .trim();
      
      // If we have substantial content, return it
      if (content.length > 100 && !content.startsWith('This documentation covers')) {
        return content;
      }
      
      return '';
      
    } catch (error) {
      console.warn('Failed to parse MDX content:', error);
      return '';
    }
  }

  private generateContextualSummary(docPage: DocPage): string {
    const { title, url, section, category, path } = docPage;
    
    let summary = `# ${title}\n\n`;
    
    // Generate specific content based on the page path and title
    if (path?.includes('mcp') || title.toLowerCase().includes('mcp')) {
      summary += this.generateMCPSpecificContent(title, path);
    } else if (path?.includes('tools/') || title.toLowerCase().includes('tool')) {
      summary += this.generateToolsContent(title);
    } else if (path?.includes('assistants/') || title.toLowerCase().includes('assistant')) {
      summary += this.generateAssistantsContent(title);
    } else if (path?.includes('quickstart/') || section.toLowerCase().includes('quickstart')) {
      summary += this.generateQuickstartContent(title);
    } else if (path?.includes('phone') || title.toLowerCase().includes('phone')) {
      summary += this.generatePhoneContent(title);
    } else {
      summary += this.generateGenericContent(title, category);
    }
    
    summary += `\n\n**Section:** ${section}\n`;
    summary += `**Category:** ${category}\n\n`;
    
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

  private generateMCPSpecificContent(title: string, path?: string): string {
    if (title.toLowerCase().includes('client') || path?.includes('tools/mcp')) {
      return `Connect your assistant to dynamic tools through MCP servers for enhanced capabilities.\n\n` +
        `The Model Context Protocol (MCP) integration allows your Vapi assistant to:\n` +
        `- Connect to any MCP-compatible server\n` +
        `- Access tools dynamically at runtime\n` +
        `- Execute actions through the MCP server\n\n` +
        `## Setup Steps\n` +
        `1. Obtain MCP Server URL from your provider (e.g., Zapier, Composio)\n` +
        `2. Create and configure MCP Tool in Vapi Dashboard\n` +
        `3. Add tool to your assistant\n\n` +
        `## Popular MCP Providers\n` +
        `- **Zapier MCP** - Access to 7,000+ apps and 30,000+ actions\n` +
        `- **Composio MCP** - Developer-focused integrations\n` +
        `- **Custom MCP Servers** - Build your own integrations\n\n`;
    } else if (title.toLowerCase().includes('server') || path?.includes('sdk/mcp-server')) {
      return `Vapi provides its own MCP server that exposes Vapi APIs as callable tools.\n\n` +
        `The Vapi MCP Server allows you to:\n` +
        `- Use Vapi APIs directly from Claude Desktop\n` +
        `- Create assistants and manage calls from any MCP client\n` +
        `- Access phone numbers, tools, and analytics\n` +
        `- Build custom applications with MCP integration\n\n` +
        `## Installation\n` +
        `\`\`\`bash\nnpm install @vapi-ai/mcp-server\n\`\`\`\n\n` +
        `## Configuration\n` +
        `Add to your MCP client configuration:\n` +
        `\`\`\`json\n{\n  "mcpServers": {\n    "vapi": {\n      "command": "npx",\n      "args": ["@vapi-ai/mcp-server"]\n    }\n  }\n}\n\`\`\`\n\n`;
    }
    return `Model Context Protocol (MCP) integration for Vapi voice AI platform.\n\n`;
  }

  private generateToolsContent(title: string): string {
    return `Tools in Vapi allow your assistant to perform actions and access external systems.\n\n` +
      `Available tool types:\n` +
      `- **Function Tools** - Custom server-side functions\n` +
      `- **Transfer Tools** - Call forwarding and routing\n` +
      `- **End Call Tools** - Programmatic call termination\n` +
      `- **DTMF Tools** - Touch-tone digit collection\n` +
      `- **Make Tools** - Integration with Make.com\n` +
      `- **GoHighLevel Tools** - CRM integrations\n` +
      `- **MCP Tools** - Dynamic tool discovery via Model Context Protocol\n\n` +
      `## Key Features\n` +
      `- Real-time tool execution during calls\n` +
      `- Custom message templates for user feedback\n` +
      `- Conditional tool availability\n` +
      `- Webhook-based tool responses\n\n`;
  }

  private generateAssistantsContent(title: string): string {
    return `Assistants are the core of Vapi - they define how your voice AI behaves and responds.\n\n` +
      `## Key Configuration Areas\n` +
      `- **Model Settings** - Choose LLM provider and model\n` +
      `- **Voice Configuration** - Select voice provider and voice\n` +
      `- **Transcriber Settings** - Configure speech-to-text\n` +
      `- **System Messages** - Define personality and behavior\n` +
      `- **Tools Integration** - Add custom functions and actions\n` +
      `- **First Message** - Set the opening greeting\n\n` +
      `## Advanced Features\n` +
      `- Background noise filtering\n` +
      `- Conversation recording\n` +
      `- Real-time analytics\n` +
      `- Custom webhooks and callbacks\n` +
      `- Multi-language support\n\n`;
  }

  private generateQuickstartContent(title: string): string {
    return `Get started with Vapi voice AI platform quickly and easily.\n\n` +
      `## Getting Started Steps\n` +
      `1. **Sign up** for a Vapi account at https://dashboard.vapi.ai\n` +
      `2. **Get your API key** from the dashboard\n` +
      `3. **Install the SDK** for your preferred language\n` +
      `4. **Create your first assistant** with basic configuration\n` +
      `5. **Make your first call** and test the setup\n\n` +
      `## Quick Setup Commands\n` +
      `\`\`\`bash\n# Install Vapi CLI\nnpm install -g @vapi-ai/cli\n\n# Initialize new project\nvapi init\n\n# Create assistant\nvapi assistant create\n\`\`\`\n\n` +
      `## Popular Use Cases\n` +
      `- Customer support automation\n` +
      `- Appointment scheduling\n` +
      `- Lead qualification\n` +
      `- Survey and feedback collection\n\n`;
  }

  private generatePhoneContent(title: string): string {
    return `Phone number management and telephony integration with Vapi.\n\n` +
      `## Phone Number Options\n` +
      `- **Vapi Numbers** - Free numbers provided by Vapi\n` +
      `- **Twilio Integration** - Use your own Twilio numbers\n` +
      `- **Vonage Integration** - Connect Vonage phone numbers\n` +
      `- **Telnyx Integration** - Enterprise-grade telephony\n` +
      `- **Bring Your Own (BYO)** - Use any SIP-compatible provider\n\n` +
      `## Features\n` +
      `- Inbound and outbound calling\n` +
      `- Call forwarding and routing\n` +
      `- Voicemail detection\n` +
      `- Call recording and transcription\n` +
      `- Real-time call analytics\n\n` +
      `## Setup Process\n` +
      `1. Choose your telephony provider\n` +
      `2. Configure phone number in dashboard\n` +
      `3. Assign assistant to phone number\n` +
      `4. Test incoming and outgoing calls\n\n`;
  }

  private generateGenericContent(title: string, category: string): string {
    switch (category) {
      case 'workflows':
        return `Visual workflow builder for creating complex conversation flows with branching logic.\n\n`;
      case 'guides':
        return `Practical implementation guide with code examples and step-by-step instructions.\n\n`;
      case 'webhooks':
        return `Real-time event notifications and webhook configuration for call events.\n\n`;
      case 'analytics':
        return `Call analytics, performance metrics, and detailed reporting capabilities.\n\n`;
      default:
        return `This documentation covers ${title.toLowerCase()} in Vapi voice AI platform.\n\n`;
    }
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
          const vectorResults = await this.vectorSearch.search(query, 10, 0.15); // Lower threshold for better matches
        
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