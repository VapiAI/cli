import Fuse from "fuse.js";
import { VapiDocumentation } from "../utils/documentation-data.js";

/**
 * Search Vapi documentation for specific topics, features, or concepts
 */
export async function searchDocumentation(
  query: string,
  category: string = "all",
  limit: number = 5
): Promise<string> {
  try {
    const docs = VapiDocumentation.getAllDocs();
    
    // Filter by category if specified
    let filteredDocs = docs;
    if (category !== "all") {
      filteredDocs = docs.filter(doc => doc.category === category);
    }

    // Set up Fuse.js for fuzzy searching
    const fuse = new Fuse(filteredDocs, {
      keys: [
        { name: "title", weight: 0.4 },
        { name: "description", weight: 0.3 },
        { name: "content", weight: 0.2 },
        { name: "tags", weight: 0.1 },
      ],
      threshold: 0.4,
      includeScore: true,
      includeMatches: true,
    });

    const results = fuse.search(query, { limit });

    if (results.length === 0) {
      return `No documentation found for "${query}" in category "${category}". Try:\n\n` +
        "• Using different keywords\n" +
        "• Searching in 'all' categories\n" +
        "• Check our complete documentation at https://docs.vapi.ai\n\n" +
        "Popular topics: assistants, phone calls, tools, webhooks, voice settings";
    }

    let response = `# 📚 Vapi Documentation Search Results\n\n`;
    response += `Found ${results.length} result(s) for "${query}"\n\n`;

    results.forEach((result, index) => {
      const doc = result.item;
      const score = Math.round((1 - (result.score || 0)) * 100);
      
      response += `## ${index + 1}. ${doc.title}\n`;
      response += `**Relevance:** ${score}% | **Category:** ${doc.category}\n\n`;
      response += `${doc.description}\n\n`;
      
      // Add content preview (first 200 chars)
      if (doc.content && doc.content.length > 200) {
        response += `**Preview:** ${doc.content.substring(0, 200)}...\n\n`;
      } else if (doc.content) {
        response += `**Content:** ${doc.content}\n\n`;
      }
      
      if (doc.url) {
        response += `**📖 Read more:** ${doc.url}\n\n`;
      }
      
      response += "---\n\n";
    });

    response += `💡 **Tip:** For more detailed information, use \`get_guides\`, \`get_examples\`, or \`get_api_reference\` tools.`;

    return response;
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `❌ Error searching documentation: ${errorMessage}\n\nPlease try again or visit https://docs.vapi.ai`;
  }
} 