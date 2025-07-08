import { DocsFetcher } from "../utils/docs-fetcher.js";

const docsFetcher = new DocsFetcher();

/**
 * Get step-by-step guides for implementing Vapi features
 */
export async function getGuides(
  topic: string,
  level: string = "all"
): Promise<string> {
  try {
    // Get all guides first
    const allGuides = await docsFetcher.getGuides();
    
    // Search for topic in guides
    const searchTerm = topic.toLowerCase();
    let relevantGuides = allGuides.filter(guide =>
      guide.title.toLowerCase().includes(searchTerm) ||
      guide.section.toLowerCase().includes(searchTerm) ||
      guide.url.toLowerCase().includes(searchTerm)
    );

    // If no direct matches, try broader search
    if (relevantGuides.length === 0) {
      const broadSearchResults = await docsFetcher.searchDocumentation(topic + " guide");
      relevantGuides = broadSearchResults.results.slice(0, 3);
    }

    if (relevantGuides.length === 0) {
      return `# 📖 No Guides Found

No guides found for "${topic}".

## 📚 Available Guide Categories:

${allGuides.slice(0, 8).map(guide => `- **${guide.title}** - ${guide.section}`).join('\n')}

## 🎯 Popular Guide Topics:

- **Introduction** - Getting started with Vapi
- **Phone calls** - Making and receiving calls
- **Web calls** - Web integration
- **Workflows** - Building conversation flows
- **Assistants** - Creating voice assistants
- **Tools** - Adding custom functions
- **Webhooks** - Handling real-time events
- **Best practices** - Prompting and debugging

## 💡 Tips:
- Try searching for broader terms (e.g., "phone" instead of "telephone")
- Use the \`search_documentation\` tool for more general searches
- Check the **Examples** section for code samples

Try searching for one of the popular topics above!`;
    }

    let response = `# 📖 Implementation Guides for "${topic}"\n\n`;
    response += `Found ${relevantGuides.length} guide(s) for "${topic}"\n`;
    if (level !== "all") {
      response += `**Level:** ${level}\n`;
    }
    response += `\n`;

    // Fetch and return actual content for each guide
    for (let i = 0; i < Math.min(relevantGuides.length, 3); i++) {
      const guide = relevantGuides[i];
      if (!guide) continue;
      
              try {
          response += `## 📚 ${i + 1}. ${guide.title}\n\n`;
        response += `**Section:** ${guide.section}\n`;
        response += `**URL:** ${guide.url}\n\n`;
        
        // Use the already-extracted content
        if (guide.content && guide.content.length > 50) {
          // Truncate very long content for readability
          let contentToShow = guide.content;
          if (contentToShow.length > 1500) {
            contentToShow = contentToShow.substring(0, 1500) + '...\n\n*[Content truncated - visit URL for complete guide]*';
          }
          
          response += `### Content:\n\n${contentToShow}\n\n`;
        } else {
          response += `*Content extraction in progress - visit URL for complete guide*\n\n`;
        }
        
        response += `---\n\n`;
        
      } catch (error) {
        response += `## 📄 ${i + 1}. ${guide.title}\n\n`;
        response += `**Section:** ${guide.section}\n`;
        response += `**URL:** ${guide.url}\n\n`;
        response += `⚠️ Content temporarily unavailable. Please visit the URL above.\n\n`;
        response += `---\n\n`;
      }
    }

    response += `## 🎯 Next Steps\n\n`;
    response += `After reviewing these guides:\n`;
    response += `- Use \`get_examples\` to see code implementations\n`;
    response += `- Use \`get_api_reference\` for detailed API documentation\n`;
    response += `- Visit the URLs above for interactive content\n`;
    response += `- Check the **Quickstart** guides for basic setup\n\n`;
    
    response += `## 🔗 Additional Resources\n\n`;
    response += `- **All Guides:** https://docs.vapi.ai/guides\n`;
    response += `- **Quickstart:** https://docs.vapi.ai/quickstart/introduction\n`;
    response += `- **API Reference:** https://docs.vapi.ai/api-reference\n`;
    response += `- **Dashboard:** https://dashboard.vapi.ai\n`;
    response += `- **Discord Community:** https://discord.gg/vapi`;

    return response;
    
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `# ❌ Guides Error

Failed to fetch guides: ${errorMessage}

## 🛠️ Troubleshooting:
- The documentation server might be temporarily unavailable
- Try again in a few moments
- Check your internet connection

## 📋 Manual Resources:
- **All Guides:** https://docs.vapi.ai/guides
- **Quickstart:** https://docs.vapi.ai/quickstart/introduction
- **Phone Calls:** https://docs.vapi.ai/quickstart/phone
- **Web Integration:** https://docs.vapi.ai/quickstart/web

## 🎯 Popular Guides:
- **Getting Started:** https://docs.vapi.ai/quickstart/introduction
- **Phone Calls:** https://docs.vapi.ai/quickstart/phone
- **Web Calls:** https://docs.vapi.ai/quickstart/web
- **Workflows:** https://docs.vapi.ai/workflows/quickstart`;
  }
} 