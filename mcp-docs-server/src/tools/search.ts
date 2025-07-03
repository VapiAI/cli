import { searchDocumentation as searchDocs, DocSection } from "../utils/documentation-data";

/**
 * Search Vapi documentation for specific topics, features, or concepts
 */
export async function searchDocumentation(
  query: string,
  category: string = "all",
  limit: number = 5
): Promise<string> {
  try {
    // Use the Vapi documentation search
    const searchCategory = category === "all" ? undefined : category;
    const searchResults = await searchDocs(query, searchCategory);
    
    // Limit results
    const limitedResults = searchResults.slice(0, limit);

    if (limitedResults.length === 0) {
      return `No documentation found for "${query}" in category "${category}". Try:\n\n` +
        "• Using different keywords\n" +
        "• Searching in 'all' categories\n" +
        "• Check our complete documentation at https://docs.vapi.ai\n\n" +
        "Popular topics: assistants, phone calls, tools, webhooks, voice settings, workflows, campaigns";
    }

    let response = `# 📚 Vapi Documentation Search Results\n\n`;
    response += `Found ${limitedResults.length} result(s) for "${query}"\n\n`;

    limitedResults.forEach((doc: DocSection, index: number) => {
      response += `## ${index + 1}. ${doc.title}\n`;
      response += `**Category:** ${doc.category}\n\n`;
      response += `${doc.description}\n\n`;
      
      if (doc.url) {
        response += `**📖 Read more:** ${doc.url}\n\n`;
      }
      
      response += "---\n\n";
    });

    response += `💡 **Tip:** For more detailed information, use \`get_guides\`, \`get_examples\`, or \`get_api_reference\` tools.\n\n`;
    response += `📄 **Full Documentation:** https://docs.vapi.ai`;

    return response;
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `❌ Error searching documentation: ${errorMessage}\n\nPlease try again or visit https://docs.vapi.ai`;
  }
} 