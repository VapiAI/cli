import { getApiEndpoints } from "../utils/documentation-data";

/**
 * Get detailed API reference information for Vapi endpoints
 */
export async function getApiReference(
  endpoint: string,
  method: string = "all",
  includeExamples: boolean = true
): Promise<string> {
  try {
    // Get API documentation from real Vapi docs
    const apiDocs = await getApiEndpoints();
    
    // Search for endpoint
    const searchTerm = endpoint.toLowerCase();
    
    // Search in API documentation
    const relevantDocs = apiDocs.filter(doc =>
      doc.title.toLowerCase().includes(searchTerm) ||
      doc.description.toLowerCase().includes(searchTerm) ||
      doc.url.toLowerCase().includes(searchTerm)
    );

    if (relevantDocs.length === 0) {
      return `# 🔍 No API Reference Found

No API reference found for "${endpoint}".

## Available API Documentation:
${apiDocs.slice(0, 15).map(doc => `- **${doc.title}** - ${doc.description}`).join('\n')}

## Popular API Categories:
- **assistants** - Create and manage voice assistants
- **calls** - Make outbound phone calls  
- **tools** - Manage custom functions
- **phone-numbers** - Manage phone numbers
- **campaigns** - Outbound call campaigns
- **workflows** - Conversation flow management
- **chats** - Text-based conversations

📚 **Complete API Reference:** https://docs.vapi.ai/api-reference

Try searching for one of these categories!`;
    }

    let response = `# 🔧 Vapi API Reference\n\n`;
    response += `Found ${relevantDocs.length} API reference(s) for "${endpoint}"\n`;
    if (method !== "all") {
      response += `**Method Filter:** ${method}\n`;
    }
    response += `\n`;

    // Show API documentation
    relevantDocs.forEach((doc, index) => {
      response += `## ${index + 1}. ${doc.title}\n\n`;
      response += `${doc.description}\n\n`;
      response += `**Category:** ${doc.category}\n`;
      response += `**📖 View API Reference:** ${doc.url}\n\n`;
      
      // Extract method from title if possible
      const methodMatch = doc.title.match(/^(GET|POST|PUT|DELETE|PATCH)/i);
      if (methodMatch && methodMatch[1]) {
        response += `**HTTP Method:** ${methodMatch[1].toUpperCase()}\n\n`;
      }
      
      response += "---\n\n";
    });

    response += `## 🛠️ API Resources\n\n`;
    response += `- **[Complete API Reference](https://docs.vapi.ai/api-reference)** - Full API documentation\n`;
    response += `- **[API Authentication](https://docs.vapi.ai/api-reference)** - How to authenticate API requests\n`;
    response += `- **[Rate Limits](https://docs.vapi.ai/api-reference)** - API usage limits\n`;
    response += `- **[SDKs](https://docs.vapi.ai)** - Official client libraries\n\n`;
    
    if (includeExamples) {
      response += `💡 **Need examples?** Visit the documentation links above for:\n`;
      response += `- Complete request/response examples\n`;
      response += `- Authentication samples\n`;
      response += `- Error handling patterns\n`;
      response += `- SDK usage examples\n\n`;
    }
    
    response += `🔧 **For code examples**, use the \`get_examples\` tool with your specific use case.`;

    return response;
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `❌ Error fetching API reference: ${errorMessage}\n\nPlease visit https://docs.vapi.ai/api-reference for complete API documentation`;
  }
} 