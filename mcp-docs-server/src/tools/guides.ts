import { getGuides as getGuidesDocs } from "../utils/documentation-data";

/**
 * Get step-by-step guides for implementing Vapi features
 */
export async function getGuides(
  topic: string,
  level: string = "all"
): Promise<string> {
  try {
    // Get guides from Vapi documentation
    const guidesLevel = level === "all" ? undefined : level;
    const guides = await getGuidesDocs(guidesLevel);
    
    // Search for topic in guides
    const searchTerm = topic.toLowerCase();
    const relevantGuides = guides.filter(guide =>
      guide.title.toLowerCase().includes(searchTerm) ||
      guide.description.toLowerCase().includes(searchTerm) ||
      guide.url.toLowerCase().includes(searchTerm)
    );

    if (relevantGuides.length === 0) {
      return `# 📖 No Guides Found

No guides found for "${topic}" with level "${level}".

## Available Guides:
${guides.slice(0, 10).map(guide => `- **${guide.title}** - ${guide.description}`).join('\n')}

## Popular Topics:
- **introduction** - Getting started with Vapi
- **phone** - Making phone calls 
- **workflows** - Building conversation flows
- **assistants** - Creating voice assistants
- **tools** - Adding custom functions
- **campaigns** - Outbound call campaigns
- **webhooks** - Handling real-time events

Try searching for one of these topics!

📚 **All Guides:** https://docs.vapi.ai/guides`;
    }

    let response = `# 📖 Vapi Implementation Guides\n\n`;
    response += `Found ${relevantGuides.length} guide(s) for "${topic}"\n`;
    if (level !== "all") {
      response += `**Level:** ${level}\n`;
    }
    response += `\n`;

    relevantGuides.forEach((guide, index) => {
      response += `## ${index + 1}. ${guide.title}\n\n`;
      response += `${guide.description}\n\n`;
      response += `**Category:** ${guide.category}\n`;
      response += `**📖 View Guide:** ${guide.url}\n\n`;
      response += "---\n\n";
    });

    response += `## 🎯 Next Steps\n\n`;
    response += `After reviewing these guides:\n`;
    response += `- Use \`get_examples\` to see code implementations\n`;
    response += `- Use \`get_api_reference\` for detailed API documentation\n`;
    response += `- Visit [Vapi Dashboard](https://dashboard.vapi.ai) to test your implementation\n`;
    response += `- Check out [Quickstart Guide](https://docs.vapi.ai/quickstart/introduction)\n\n`;
    
    response += `📋 **Popular Getting Started Guides:**\n`;
    response += `- **Introduction:** https://docs.vapi.ai/quickstart/introduction\n`;
    response += `- **Phone Calls:** https://docs.vapi.ai/quickstart/phone\n`;
    response += `- **Web Calls:** https://docs.vapi.ai/quickstart/web\n`;
    response += `- **Workflows:** https://docs.vapi.ai/workflows/quickstart`;

    return response;
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `❌ Error fetching guides: ${errorMessage}\n\nPlease visit https://docs.vapi.ai/guides for guides`;
  }
} 