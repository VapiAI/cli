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
- Try different keywords (e.g., "MCP server" instead of "model context protocol")  
- Use "get_examples" tool to find code examples specifically
- Try broader terms first, then narrow down
- Check for typos in your search query

## 📚 Popular Topics Available:
- **MCP integration** - Claude Desktop configuration and setup
- **Voice assistants** - Creating and configuring assistants
- **Phone calls** - Outbound calls, inbound handling, phone numbers
- **Tools** - Function tools, MCP tools, custom integrations
- **Webhooks** - Real-time events and server configuration  
- **API reference** - Assistants, calls, tools, phone numbers
- **Voice providers** - ElevenLabs, OpenAI, Cartesia, PlayHT
- **Workflows** - Building conversation flows and automation

Try searching for one of these topics or use 'get_examples [topic]' for code samples!`;
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

    // Return results with already-extracted content
    for (let i = 0; i < limitedResults.length; i++) {
      const page = limitedResults[i];
      if (!page) continue;
      
      response += `## 📄 ${i + 1}. ${page.title}\n\n`;
      response += `**Section:** ${page.section}\n`;
      response += `**Category:** ${page.category}\n`;
      response += `**URL:** ${page.url}\n\n`;
      
      // Use the already-extracted content
      if (page.content && page.content.length > 50) {
        // Truncate very long content for readability
        let contentToShow = page.content;
        if (contentToShow.length > 2000) {
          contentToShow = contentToShow.substring(0, 2000) + '...\n\n*[Content truncated - visit URL for complete documentation]*';
        }
        
        response += `### Content:\n\n${contentToShow}\n\n`;
      } else {
        response += `*Content extraction in progress - visit URL for complete documentation*\n\n`;
      }
      
      response += `---\n\n`;
    }

    response += `## 🎯 Next Steps\n\n`;
    response += `- Use \`get_examples [feature]\` to see code examples for specific features\n`;
    response += `- Use \`get_guides [topic]\` for step-by-step tutorials\n`;
    response += `- Use \`get_api_reference [endpoint]\` for API documentation\n`;
    response += `- Visit the URLs above for complete interactive documentation\n\n`;

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