import { getExamples as getExampleDocs } from "../utils/documentation-data";

/**
 * Get code examples for specific Vapi features or use cases
 */
export async function getExamples(
  feature: string,
  language: string = "typescript",
  framework: string = "all"
): Promise<string> {
  try {
    // Get example documentation pages from Vapi docs
    const exampleDocs = await getExampleDocs(framework !== "all" ? framework : undefined);
    
    // Filter by feature
    const searchTerm = feature.toLowerCase();
    const filteredExamples = exampleDocs.filter(doc => 
      doc.title.toLowerCase().includes(searchTerm) ||
      doc.description.toLowerCase().includes(searchTerm) ||
      doc.url.toLowerCase().includes(searchTerm)
    );

    if (filteredExamples.length === 0) {
      return `# 📝 No Examples Found

No examples found for "${feature}" with framework "${framework}".

## Available Example Categories:
${exampleDocs.slice(0, 10).map(ex => `- **${ex.title}** - ${ex.description}`).join('\n')}

## Popular Features:
- **assistants** - Create and manage voice assistants  
- **calls** - Make outbound phone calls
- **workflows** - Build conversation flows
- **tools** - Add custom function calling
- **webhooks** - Handle real-time events
- **campaigns** - Outbound call campaigns

Try searching for one of these features!

📚 **Full Examples:** https://docs.vapi.ai/guides`;
    }

    let response = `# 💻 Vapi Code Examples\n\n`;
    response += `Found ${filteredExamples.length} example(s) for "${feature}"\n`;
    if (framework !== "all") {
      response += `**Framework:** ${framework}\n`;
    }
    response += `**Language:** ${language}\n\n`;

    filteredExamples.forEach((doc, index) => {
      response += `## ${index + 1}. ${doc.title}\n\n`;
      response += `${doc.description}\n\n`;
      response += `**Category:** ${doc.category}\n`;
      response += `**📖 View Example:** ${doc.url}\n\n`;
      response += "---\n\n";
    });

    response += `💡 **Getting Started:**\n`;
    response += `- **Quickstart Guide:** https://docs.vapi.ai/quickstart/introduction\n`;
    response += `- **Phone Calls:** https://docs.vapi.ai/quickstart/phone\n`;
    response += `- **Web Calls:** https://docs.vapi.ai/quickstart/web\n`;
    response += `- **All Guides:** https://docs.vapi.ai/guides\n\n`;
    
    response += `🔧 **Need specific ${language} examples?** Visit the documentation links above for:\n`;
    response += `- Complete, working code samples\n`;
    response += `- Step-by-step implementation guides\n`;
    response += `- Best practices and tips\n`;
    response += `- Framework-specific integrations`;

    return response;
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `❌ Error fetching examples: ${errorMessage}\n\nPlease visit https://docs.vapi.ai/guides for examples`;
  }
} 