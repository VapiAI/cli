import { DocsFetcher } from "../utils/docs-fetcher.js";

const docsFetcher = new DocsFetcher();

/**
 * Search Vapi documentation for specific topics, features, or concepts
 */
export async function searchDocumentation(
  query: string,
  category: string = "all",
  limit: number = 5
): Promise<string> {
  try {
    // Search the documentation
    const searchResult = await docsFetcher.searchDocumentation(query, category);
    const searchResults = searchResult.results;
    const usedVectorSearch = searchResult.usedVectorSearch;
    
    // Limit results
    const limitedResults = searchResults.slice(0, limit);

    if (limitedResults.length === 0) {
      return `# 🔍 No Results Found

No documentation found for "${query}" in category "${category}".

## 💡 Suggestions:
- Try different keywords (e.g., "phone calls" instead of "calling")
- Search in all categories instead of specific ones
- Check for typos in your search query
- Try broader terms first, then narrow down

## 📚 Popular Topics:
- **Phone calls** - Making and receiving calls
- **Assistants** - Creating voice assistants
- **Workflows** - Building conversation flows
- **Tools** - Custom function calling
- **Webhooks** - Real-time event handling
- **Best practices** - Prompting and debugging guides

Try searching for one of these topics!`;
    }

    // Display search method used
    const vectorIndexSize = docsFetcher.getVectorIndexSize();
    
    let response = `# 🔍 Search Results for "${query}"\n\n`;
    
    if (usedVectorSearch) {
      response += `🧠 **Enhanced AI Search** - Found ${limitedResults.length} semantically relevant page(s) from ${vectorIndexSize} indexed documents:\n\n`;
    } else if (vectorIndexSize > 0) {
      response += `📝 **Text Search** - Found ${limitedResults.length} relevant page(s) (vector search available but no semantic matches above threshold):\n\n`;
    } else {
      response += `📝 **Text Search** - Found ${limitedResults.length} relevant page(s) (vector search initializing...):\n\n`;
    }

    // Fetch and return actual content for each result
    for (let i = 0; i < limitedResults.length; i++) {
      const page = limitedResults[i];
      if (!page) continue;
      
      try {
        const content = await docsFetcher.fetchPageContent(page);
        
        response += `## 📄 ${i + 1}. ${page.title}\n\n`;
        response += `**Section:** ${page.section}\n`;
        response += `**Category:** ${page.category}\n`;
        response += `**URL:** ${page.url}\n\n`;
        
        // Add the actual content
        response += `### Content:\n\n${content}\n\n`;
        response += `---\n\n`;
        
      } catch (error) {
        response += `## 📄 ${i + 1}. ${page.title}\n\n`;
        response += `**Section:** ${page.section}\n`;
        response += `**URL:** ${page.url}\n\n`;
        response += `⚠️ Content temporarily unavailable. Please visit the URL above.\n\n`;
        response += `---\n\n`;
      }
    }

    response += `## 🎯 Next Steps\n\n`;
    response += `- Use \`get_examples\` to see code examples\n`;
    response += `- Use \`get_guides\` for step-by-step tutorials\n`;
    response += `- Use \`get_api_reference\` for API documentation\n`;
    response += `- Visit the URLs above for interactive content\n\n`;

    if (usedVectorSearch) {
      response += `## 🧠 AI-Powered Search\n\n`;
      response += `This search used semantic similarity to find the most relevant results, even if they don't contain your exact keywords. `;
      response += `Vector search indexed ${vectorIndexSize} documents for intelligent matching.\n\n`;
    }

    response += `## 🔗 Resources\n\n`;
    response += `- **Main Documentation:** https://docs.vapi.ai\n`;
    response += `- **API Reference:** https://docs.vapi.ai/api-reference\n`;
    response += `- **Community Discord:** https://discord.gg/vapi`;

    return response;
    
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `# ❌ Search Error

Failed to search documentation: ${errorMessage}

## 🛠️ Troubleshooting:
- The documentation server might be temporarily unavailable
- Vector search model might be initializing (first run takes longer)
- Try again in a few moments
- Check your internet connection
- Contact support if the issue persists

## 📋 Manual Resources:
- **Main Documentation:** https://docs.vapi.ai
- **API Reference:** https://docs.vapi.ai/api-reference
- **Guides:** https://docs.vapi.ai/guides
- **Examples:** https://docs.vapi.ai/guides`;
  }
} 