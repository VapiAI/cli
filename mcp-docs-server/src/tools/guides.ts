import { VapiDocumentation, DocItem } from "../utils/documentation-data.js";

/**
 * Get step-by-step guides for implementing Vapi features
 */
export async function getGuides(
  topic: string,
  level: string = "all"
): Promise<string> {
  try {
    // Get guides category docs
    const guides = VapiDocumentation.getDocsByCategory("guides");
    
    // Search for topic in guides
    const searchTerm = topic.toLowerCase();
    const relevantGuides = guides.filter((guide: DocItem) =>
      guide.title.toLowerCase().includes(searchTerm) ||
      guide.description.toLowerCase().includes(searchTerm) ||
      guide.content.toLowerCase().includes(searchTerm) ||
      guide.tags.some(tag => tag.toLowerCase().includes(searchTerm))
    );

    if (relevantGuides.length === 0) {
      return `# 📖 No Guides Found

No guides found for "${topic}".

## Available Guides:
${guides.map((guide: DocItem) => `- **${guide.title}** - ${guide.description}`).join('\n')}

## Popular Topics:
- **getting started** - Create your first voice assistant
- **phone calls** - Make outbound calls
- **tools** - Add custom functions
- **voice settings** - Configure voice providers
- **webhooks** - Handle real-time events
- **assistants** - Create and manage assistants

Try searching for one of these topics!`;
    }

    let response = `# 📖 Vapi Implementation Guides\n\n`;
    response += `Found ${relevantGuides.length} guide(s) for "${topic}"\n`;
    if (level !== "all") {
      response += `**Level:** ${level}\n`;
    }
    response += `\n`;

    relevantGuides.forEach((guide: DocItem, index: number) => {
      response += `## ${index + 1}. ${guide.title}\n\n`;
      response += `${guide.description}\n\n`;
      
      // Add the full content
      response += guide.content + "\n\n";
      
      response += `**📅 Last Updated:** ${guide.lastUpdated}\n`;
      response += `**🔗 View Online:** ${guide.url}\n`;
      
      if (guide.tags.length > 0) {
        response += `**🏷️ Tags:** ${guide.tags.join(', ')}\n`;
      }
      
      response += "\n---\n\n";
    });

    response += `## 🎯 Next Steps\n\n`;
    response += `After following this guide:\n`;
    response += `- Use \`get_examples\` to see code implementations\n`;
    response += `- Use \`get_api_reference\` for detailed API docs\n`;
    response += `- Visit [Vapi Dashboard](https://dashboard.vapi.ai) to test your implementation\n`;
    response += `- Join our [Discord Community](https://discord.gg/vapi) for support`;

    return response;
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `❌ Error fetching guides: ${errorMessage}\n\nPlease try again or visit https://docs.vapi.ai`;
  }
} 