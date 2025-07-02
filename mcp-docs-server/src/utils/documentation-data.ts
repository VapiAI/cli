export interface DocSection {
  title: string;
  url: string;
  description: string;
  category: string;
  content?: string;
}

export interface DocIndex {
  sections: DocSection[];
  lastFetched: Date;
  version: string;
}

export interface CodeExample {
  id: string;
  title: string;
  description: string;
  code: string;
  language: string;
  framework?: string;
  category: string;
  tags: string[];
}

export interface ApiEndpoint {
  id: string;
  path: string;
  method: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
  description: string;
  parameters?: Record<string, any>;
  requestBody?: Record<string, any>;
  responses?: Record<string, any>;
  examples?: {
    request?: string;
    response?: string;
  };
}

// Cache for documentation
let docCache: DocIndex | null = null;
const CACHE_TTL = 1000 * 60 * 60; // 1 hour

export async function fetchDocumentationIndex(): Promise<DocIndex> {
  // Return cached version if still fresh
  if (docCache && Date.now() - docCache.lastFetched.getTime() < CACHE_TTL) {
    return docCache;
  }

  console.log('Fetching fresh documentation from Vapi...');
  
  try {
    // Fetch the structured documentation index
    const response = await fetch('https://docs.vapi.ai/llms.txt');
    if (!response.ok) {
      throw new Error(`Failed to fetch docs: ${response.status}`);
    }
    
    const docsText = await response.text();
    const sections = parseDocumentationIndex(docsText);
    
    docCache = {
      sections,
      lastFetched: new Date(),
      version: 'live'
    };
    
    console.log(`Fetched ${sections.length} documentation sections`);
    return docCache;
    
  } catch (error) {
    console.error('Failed to fetch live documentation:', error);
    
    // Return minimal fallback if fetch fails
    return {
      sections: [{
        title: "Vapi Documentation",
        url: "https://docs.vapi.ai",
        description: "Visit docs.vapi.ai for the latest documentation",
        category: "fallback"
      }],
      lastFetched: new Date(),
      version: 'fallback'
    };
  }
}

function parseDocumentationIndex(docsText: string): DocSection[] {
  const sections: DocSection[] = [];
  const lines = docsText.split('\n');
  
  let currentCategory = 'general';
  
  for (const line of lines) {
    const trimmedLine = line.trim();
    
    // Skip empty lines and main headers
    if (!trimmedLine || trimmedLine === '# Vapi' || trimmedLine === '____') {
      continue;
    }
    
    // Detect category headers
    if (trimmedLine.startsWith('## ')) {
      currentCategory = trimmedLine.replace('## ', '').toLowerCase().replace(/\s+/g, '-');
      continue;
    }
    
    // Parse documentation links
    // Format: - [Title](url): Description
    const linkMatch = trimmedLine.match(/^- \[(.*?)\]\((.*?)\):\s*(.*)/);
    if (linkMatch && linkMatch[1] && linkMatch[2] && linkMatch[3]) {
      const title = linkMatch[1];
      const url = linkMatch[2];
      const description = linkMatch[3];
      
      sections.push({
        title: title.trim(),
        url: url.startsWith('http') ? url : `https://docs.vapi.ai${url}`,
        description: description.trim(),
        category: currentCategory
      });
    }
  }
  
  return sections;
}

export async function fetchDocumentationContent(url: string): Promise<string> {
  try {
    // For now, return a helpful summary with the real URL
    const index = await fetchDocumentationIndex();
    const section = index.sections.find(s => s.url === url);
    
    if (section) {
      return `# ${section.title}

${section.description}

**Category:** ${section.category}
**Full documentation:** ${section.url}

Visit the link above for the complete, interactive documentation with examples, code samples, and detailed explanations.

This is the official Vapi documentation, always up-to-date with the latest features and API changes.`;
    }
    
    return `Documentation not found for URL: ${url}`;
    
  } catch (error) {
    console.error('Failed to fetch documentation content:', error);
    return `Error fetching documentation: ${error}`;
  }
}

// Utility function to search documentation
export async function searchDocumentation(query: string, category?: string): Promise<DocSection[]> {
  const index = await fetchDocumentationIndex();
  const searchTerm = query.toLowerCase();
  
  return index.sections.filter(section => {
    const matchesQuery = 
      section.title.toLowerCase().includes(searchTerm) ||
      section.description.toLowerCase().includes(searchTerm) ||
      section.url.toLowerCase().includes(searchTerm);
    
    const matchesCategory = !category || section.category === category;
    
    return matchesQuery && matchesCategory;
  }).slice(0, 20); // Limit results
}

// Get sections by category
export async function getDocumentationByCategory(category: string): Promise<DocSection[]> {
  const index = await fetchDocumentationIndex();
  return index.sections.filter(section => section.category === category);
}

// Get all available categories
export async function getDocumentationCategories(): Promise<string[]> {
  const index = await fetchDocumentationIndex();
  const categories = [...new Set(index.sections.map(s => s.category))];
  return categories.sort();
}

// Extract API references
export async function getApiEndpoints(): Promise<DocSection[]> {
  const index = await fetchDocumentationIndex();
  return index.sections.filter(section => 
    section.category === 'api-docs' || 
    section.title.toLowerCase().includes('api reference') ||
    section.url.includes('/api-reference/')
  );
}

// Get examples from documentation
export async function getExamples(framework?: string): Promise<DocSection[]> {
  const index = await fetchDocumentationIndex();
  let examples = index.sections.filter(section => 
    section.category === 'docs' && (
      section.title.toLowerCase().includes('example') ||
      section.title.toLowerCase().includes('quickstart') ||
      section.title.toLowerCase().includes('guide') ||
      section.description.toLowerCase().includes('example')
    )
  );

  if (framework) {
    examples = examples.filter(section => 
      section.title.toLowerCase().includes(framework.toLowerCase()) ||
      section.description.toLowerCase().includes(framework.toLowerCase())
    );
  }

  return examples;
}

// Get guides from documentation  
export async function getGuides(level?: string): Promise<DocSection[]> {
  const index = await fetchDocumentationIndex();
  let guides = index.sections.filter(section => 
    section.title.toLowerCase().includes('guide') ||
    section.title.toLowerCase().includes('tutorial') ||
    section.title.toLowerCase().includes('quickstart') ||
    section.title.toLowerCase().includes('workflow') ||
    section.description.toLowerCase().includes('learn')
  );

  if (level) {
    // Filter by difficulty level (beginner, intermediate, advanced)
    const isBeginnerLevel = (s: DocSection) => 
      s.title.toLowerCase().includes('quickstart') ||
      s.title.toLowerCase().includes('introduction') ||
      s.title.toLowerCase().includes('getting started');
    
    const isAdvancedLevel = (s: DocSection) =>
      s.title.toLowerCase().includes('advanced') ||
      s.title.toLowerCase().includes('enterprise') ||
      s.title.toLowerCase().includes('custom');

    switch (level) {
      case 'beginner':
        guides = guides.filter(isBeginnerLevel);
        break;
      case 'advanced':
        guides = guides.filter(isAdvancedLevel);
        break;
      case 'intermediate':
        guides = guides.filter(s => !isBeginnerLevel(s) && !isAdvancedLevel(s));
        break;
    }
  }

  return guides;
}

// Get changelog entries
export async function getChangelog(limit = 10): Promise<DocSection[]> {
  const index = await fetchDocumentationIndex();
  const changelog = index.sections.filter(section => 
    section.category === 'docs' && section.url.includes('/changelog/')
  );

  // Sort by date (assuming URLs contain dates)
  changelog.sort((a, b) => {
    const dateA = extractDateFromUrl(a.url);
    const dateB = extractDateFromUrl(b.url);
    return dateB.getTime() - dateA.getTime();
  });

  return changelog.slice(0, limit);
}

function extractDateFromUrl(url: string): Date {
  // Extract date from URL like /changelog/2025/1/15.mdx
  const match = url.match(/\/changelog\/(\d{4})\/(\d{1,2})\/(\d{1,2})/);
  if (match && match[1] && match[2] && match[3]) {
    const year = match[1];
    const month = match[2];
    const day = match[3];
    return new Date(parseInt(year), parseInt(month) - 1, parseInt(day));
  }
  return new Date(0); // fallback to epoch
}

// Main class for backward compatibility and easier access
export class VapiDocumentation {
  static async getAllDocs(): Promise<DocSection[]> {
    const index = await fetchDocumentationIndex();
    return index.sections;
  }

  static async getDocsByCategory(category: string): Promise<DocSection[]> {
    return getDocumentationByCategory(category);
  }

  static async searchDocs(query: string): Promise<DocSection[]> {
    return searchDocumentation(query);
  }

  static async getAllExamples(): Promise<DocSection[]> {
    return getExamples();
  }

  static async getExamplesByFramework(framework: string): Promise<DocSection[]> {
    return getExamples(framework);
  }

  static async getAllApiEndpoints(): Promise<DocSection[]> {
    return getApiEndpoints();
  }

  static async getGuides(level?: string): Promise<DocSection[]> {
    return getGuides(level);
  }

  static async getChangelog(limit = 10): Promise<DocSection[]> {
    return getChangelog(limit);
  }

  static async refreshCache(): Promise<void> {
    docCache = null;
    await fetchDocumentationIndex();
  }
} 