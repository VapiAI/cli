import { VapiDocumentation, CodeExample } from "../utils/documentation-data.js";

/**
 * Get code examples for specific Vapi features or use cases
 */
export async function getExamples(
  feature: string,
  language: string = "typescript",
  framework: string = "all"
): Promise<string> {
  try {
    let examples = VapiDocumentation.getAllExamples();
    
    // Filter by language
    if (language !== "all") {
      examples = VapiDocumentation.getExamplesByLanguage(language);
    }
    
    // Filter by framework
    if (framework !== "all") {
      examples = VapiDocumentation.getExamplesByFramework(framework);
    }
    
    // Search for feature in title, description, or tags
    const searchTerm = feature.toLowerCase();
    const filteredExamples = examples.filter((example: CodeExample) => 
      example.title.toLowerCase().includes(searchTerm) ||
      example.description.toLowerCase().includes(searchTerm) ||
      example.category.toLowerCase().includes(searchTerm) ||
      example.tags.some(tag => tag.toLowerCase().includes(searchTerm))
    );

    if (filteredExamples.length === 0) {
      return `# 📝 No Examples Found

No code examples found for "${feature}" with:
- Language: ${language}
- Framework: ${framework}

## Available Examples:
${examples.map((ex: CodeExample) => `- **${ex.title}** (${ex.language})`).join('\n')}

## Popular Features:
- **assistants** - Create and manage voice assistants
- **calls** - Make outbound phone calls
- **functions** - Add custom function calling
- **webhooks** - Handle real-time events
- **voice** - Configure voice settings

Try searching for one of these features!`;
    }

    let response = `# 💻 Vapi Code Examples\n\n`;
    response += `Found ${filteredExamples.length} example(s) for "${feature}"\n`;
    response += `**Language:** ${language} | **Framework:** ${framework}\n\n`;

    filteredExamples.forEach((example: CodeExample, index: number) => {
      response += `## ${index + 1}. ${example.title}\n\n`;
      response += `${example.description}\n\n`;
      response += `**Language:** ${example.language}`;
      if (example.framework) {
        response += ` | **Framework:** ${example.framework}`;
      }
      response += `\n\n`;
      
      response += "```" + example.language + "\n";
      response += example.code + "\n";
      response += "```\n\n";
      
      if (example.tags.length > 0) {
        response += `**Tags:** ${example.tags.join(', ')}\n\n`;
      }
      
      response += "---\n\n";
    });

    response += `💡 **Need more examples?** Check out:\n`;
    response += `- [Vapi Examples Repository](https://github.com/VapiAI/examples)\n`;
    response += `- [Documentation](https://docs.vapi.ai)\n`;
    response += `- Use \`get_guides\` for step-by-step tutorials`;

    return response;
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `❌ Error fetching examples: ${errorMessage}\n\nPlease try again or visit https://docs.vapi.ai/examples`;
  }
} 