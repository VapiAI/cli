import { DocsFetcher } from "../utils/docs-fetcher.js";

const docsFetcher = new DocsFetcher();

/**
 * Get code examples for specific Vapi features or use cases
 */
export async function getExamples(
  feature: string,
  language: string = "typescript",
  framework: string = "all"
): Promise<string> {
  try {
    // Get all examples first
    const allExamples = await docsFetcher.getExamples();
    
    // Filter by feature
    const searchTerm = feature.toLowerCase();
    const filteredExamples = allExamples.filter(page => 
      page.title.toLowerCase().includes(searchTerm) ||
      page.section.toLowerCase().includes(searchTerm) ||
      page.url.toLowerCase().includes(searchTerm)
    );

    if (filteredExamples.length === 0) {
      // If no direct matches, try a broader search
      const broadSearchResults = await docsFetcher.searchDocumentation(feature + " example");
      const exampleResults = broadSearchResults.results.slice(0, 3);
      
      if (exampleResults.length === 0) {
        return `# 📝 No Examples Found

No examples found for "${feature}".

## 📚 Available Example Categories:

${allExamples.slice(0, 8).map(ex => `- **${ex.title}** - ${ex.section}`).join('\n')}

## 🎯 Popular Example Topics:

- **MCP server** - Claude Desktop configuration with @vapi-ai/mcp-server
- **Voice assistants** - Complete assistant creation and configuration  
- **Phone calls** - Outbound calls, inbound handling, phone number setup
- **Function tools** - Custom API integrations and webhooks
- **Workflows** - Conversation flows and automation
- **Voice providers** - ElevenLabs, OpenAI TTS, and voice configuration
- **Real-time features** - Streaming, WebSocket transport
- **Integration examples** - Google Calendar, Sheets, CRM systems

## 💡 Tips:
- Try specific terms like "MCP server", "assistant create", "phone call"
- Use the \`search_documentation\` tool for broader searches
- Check the **API Reference** for endpoint-specific examples

Try searching for one of the popular topics above!`;
      }
      
      // Use broader search results
      return await formatExamplesResponse(exampleResults, feature, language, framework, true);
    }

    return await formatExamplesResponse(filteredExamples, feature, language, framework, false);
    
      } catch (error) {
      const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `# ❌ Examples Error

Failed to fetch examples: ${errorMessage}

## 🛠️ Troubleshooting:
- The documentation server might be temporarily unavailable
- Try again in a few moments
- Check your internet connection

## 📋 Manual Resources:
- **Examples Gallery:** https://docs.vapi.ai/guides
- **Quickstart Guide:** https://docs.vapi.ai/quickstart/introduction
- **GitHub Examples:** https://github.com/VapiAI/docs

## 🎯 Popular Examples:
- **Phone Calls:** https://docs.vapi.ai/quickstart/phone
- **Web Integration:** https://docs.vapi.ai/quickstart/web
- **Workflows:** https://docs.vapi.ai/workflows/quickstart`;
  }
}

async function formatExamplesResponse(
  examples: any[],
  feature: string,
  language: string,
  framework: string,
  isBroadSearch: boolean
): Promise<string> {
  const responseTitle = isBroadSearch ? 
    `# 🔍 Related Examples for "${feature}"` : 
    `# 💻 Examples for "${feature}"`;
  
  let response = `${responseTitle}\n\n`;
  
  if (isBroadSearch) {
    response += `No direct examples found for "${feature}", but here are related examples:\n\n`;
  } else {
    response += `Found ${examples.length} example(s) for "${feature}"\n\n`;
  }
  
  if (framework !== "all") {
    response += `**Framework:** ${framework}\n`;
  }
  response += `**Language:** ${language}\n\n`;

  // Fetch and return actual content for each example
  for (let i = 0; i < Math.min(examples.length, 3); i++) {
    const example = examples[i];
    if (!example) continue;
    
    try {
      response += `## 💻 ${i + 1}. ${example.title}\n\n`;
      response += `**Section:** ${example.section}\n`;
      response += `**URL:** ${example.url}\n\n`;
      
      // Use the already-extracted content
      if (example.content && example.content.length > 50) {
        // Truncate very long content for readability
        let contentToShow = example.content;
        if (contentToShow.length > 1500) {
          contentToShow = contentToShow.substring(0, 1500) + '...\n\n*[Content truncated - visit URL for complete example]*';
        }
        
        response += `### Content:\n\n${contentToShow}\n\n`;
      } else {
        response += `*Content extraction in progress - visit URL for complete example*\n\n`;
      }
      
      response += `---\n\n`;
      
    } catch (error) {
      response += `## 📄 ${i + 1}. ${example.title}\n\n`;
      response += `**Section:** ${example.section}\n`;
      response += `**URL:** ${example.url}\n\n`;
      response += `⚠️ Content temporarily unavailable. Please visit the URL above.\n\n`;
      response += `---\n\n`;
    }
  }

  response += `## 🎯 Next Steps\n\n`;
  response += `- Use \`get_guides [topic]\` for step-by-step implementation guides\n`;
  response += `- Use \`get_api_reference [endpoint]\` for detailed API documentation\n`;
  response += `- Visit the URLs above for complete interactive code examples\n`;
  response += `- Try \`search_documentation [feature]\` for related content\n\n`;
  
  response += `## 🔗 Additional Resources\n\n`;
  response += `- **All Examples:** https://docs.vapi.ai/guides\n`;
  response += `- **Quickstart:** https://docs.vapi.ai/quickstart/introduction\n`;
  response += `- **GitHub:** https://github.com/VapiAI/docs\n`;
  response += `- **Discord Community:** https://discord.gg/vapi`;

  return response;
} 