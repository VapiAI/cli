import { promises as fs } from 'fs';
import path from 'path';
import yaml from 'yaml';
import axios from 'axios';

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
}

export class DocsFetcher {
  private docsCache: DocStructure | null = null;
  private readonly GITHUB_RAW_BASE = 'https://raw.githubusercontent.com/VapiAI/docs/main';
  private readonly DOCS_YML_URL = 'https://raw.githubusercontent.com/VapiAI/docs/main/fern/docs.yml';
  private readonly CACHE_TTL = 1000 * 60 * 60; // 1 hour

  async getDocumentationStructure(): Promise<DocStructure> {
    if (this.docsCache && Date.now() - this.docsCache.lastUpdated.getTime() < this.CACHE_TTL) {
      return this.docsCache;
    }

    try {
      // Fetch docs.yml structure
      const response = await axios.get(this.DOCS_YML_URL);
      const docsConfig = yaml.parse(response.data);
      
      // Parse the navigation structure
      const pages = await this.parseNavigationStructure(docsConfig.navigation);
      
      // Organize into sections and categories
      const sections = new Map<string, DocPage[]>();
      const categories = new Map<string, DocPage[]>();
      
      pages.forEach(page => {
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

      this.docsCache = {
        pages,
        sections,
        categories,
        lastUpdated: new Date()
      };

      return this.docsCache;
      
    } catch (error) {
      throw new Error(`Failed to fetch docs structure: ${error}`);
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
      const contentUrl = `${this.GITHUB_RAW_BASE}/fern/docs/${docPage.path}`;
      const response = await axios.get(contentUrl);
      
      // Parse MDX content and extract meaningful text
      const content = this.parseMDXContent(response.data);
      
      // Cache the content
      docPage.content = content;
      
      return content;
      
    } catch (error) {
      return `# ${docPage.title}\n\nContent temporarily unavailable. Please visit: ${docPage.url}`;
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

  async searchDocumentation(query: string, category?: string): Promise<DocPage[]> {
    const docs = await this.getDocumentationStructure();
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
} 