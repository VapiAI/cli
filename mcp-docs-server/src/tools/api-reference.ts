import { DocsFetcher } from "../utils/docs-fetcher.js";

const docsFetcher = new DocsFetcher();

/**
 * Get detailed API reference information for Vapi endpoints
 */
export async function getApiReference(
  endpoint: string,
  method: string = "all",
  includeExamples: boolean = true
): Promise<string> {
  try {
    // Get all API reference pages
    const allApiPages = await docsFetcher.getApiReference();
    
    // Search for endpoint in API pages
    const searchTerm = endpoint.toLowerCase();
    let relevantApiPages = allApiPages.filter(page =>
      page.title.toLowerCase().includes(searchTerm) ||
      page.section.toLowerCase().includes(searchTerm) ||
      page.url.toLowerCase().includes(searchTerm) ||
      page.url.toLowerCase().includes(endpoint.toLowerCase())
    );

    // If no direct matches, try broader search
    if (relevantApiPages.length === 0) {
      const broadSearchResults = await docsFetcher.searchDocumentation(endpoint + " api");
      relevantApiPages = broadSearchResults.slice(0, 3);
    }

    if (relevantApiPages.length === 0) {
      return `# 🔧 No API Reference Found

No API reference found for "${endpoint}".

## 📚 Available API Endpoints:

${allApiPages.slice(0, 8).map(page => `- **${page.title}** - ${page.section}`).join('\n')}

## 🎯 Popular API Endpoints:

- **Assistants** - Create and manage voice assistants
- **Calls** - Make and manage phone calls
- **Phone Numbers** - Manage phone numbers
- **Tools** - Define custom functions
- **Webhooks** - Configure event notifications
- **Analytics** - Call analytics and insights
- **Files** - File upload and management
- **Squads** - Team management

## 💡 Tips:
- Try searching for broader terms (e.g., "assistant" instead of "assistants")
- Use the \`search_documentation\` tool for more general searches
- Check the full API reference at https://docs.vapi.ai/api-reference

Try searching for one of the popular endpoints above!`;
    }

    let response = `# 🔧 API Reference for "${endpoint}"\n\n`;
    response += `Found ${relevantApiPages.length} API reference(s) for "${endpoint}"\n`;
    if (method !== "all") {
      response += `**Method:** ${method.toUpperCase()}\n`;
    }
    response += `**Include Examples:** ${includeExamples ? 'Yes' : 'No'}\n\n`;

    // Fetch and return actual content for each API reference
    for (let i = 0; i < Math.min(relevantApiPages.length, 3); i++) {
      const apiPage = relevantApiPages[i];
      if (!apiPage) continue;
      
      try {
        const content = await docsFetcher.fetchPageContent(apiPage);
        
        response += `## 📄 ${i + 1}. ${apiPage.title}\n\n`;
        response += `**Section:** ${apiPage.section}\n`;
        response += `**Category:** ${apiPage.category}\n`;
        response += `**URL:** ${apiPage.url}\n\n`;
        
        // Add the actual content
        response += `### Content:\n\n${content}\n\n`;
        response += `---\n\n`;
        
      } catch (error) {
        response += `## 📄 ${i + 1}. ${apiPage.title}\n\n`;
        response += `**Section:** ${apiPage.section}\n`;
        response += `**URL:** ${apiPage.url}\n\n`;
        response += `⚠️ Content temporarily unavailable. Please visit the URL above.\n\n`;
        response += `---\n\n`;
      }
    }

    response += `## 🎯 Next Steps\n\n`;
    response += `After reviewing this API reference:\n`;
    response += `- Use \`get_examples\` to see code implementations\n`;
    response += `- Use \`get_guides\` for step-by-step tutorials\n`;
    response += `- Visit the URLs above for interactive API testing\n`;
    response += `- Check the **Quickstart** guides for basic setup\n\n`;
    
    response += `## 🔗 Additional Resources\n\n`;
    response += `- **Full API Reference:** https://docs.vapi.ai/api-reference\n`;
    response += `- **Interactive API:** https://api.vapi.ai/api\n`;
    response += `- **OpenAPI Spec:** https://api.vapi.ai/api-json\n`;
    response += `- **Dashboard:** https://dashboard.vapi.ai\n`;
    response += `- **Discord Community:** https://discord.gg/vapi`;

    return response;
    
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `# ❌ API Reference Error

Failed to fetch API reference: ${errorMessage}

## 🛠️ Troubleshooting:
- The documentation server might be temporarily unavailable
- Try again in a few moments
- Check your internet connection

## 📋 Manual Resources:
- **Full API Reference:** https://docs.vapi.ai/api-reference
- **Interactive API:** https://api.vapi.ai/api
- **OpenAPI Spec:** https://api.vapi.ai/api-json
- **Postman Collection:** Available in the dashboard

## 🎯 Popular API Endpoints:
- **Assistants:** https://docs.vapi.ai/api-reference/assistants
- **Calls:** https://docs.vapi.ai/api-reference/calls
- **Phone Numbers:** https://docs.vapi.ai/api-reference/phone-numbers
- **Tools:** https://docs.vapi.ai/api-reference/tools`;
  }
} 