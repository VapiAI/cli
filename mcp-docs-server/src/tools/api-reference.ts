import { VapiDocumentation, ApiEndpoint, DocItem } from "../utils/documentation-data.js";

/**
 * Get detailed API reference information for Vapi endpoints
 */
export async function getApiReference(
  endpoint: string,
  method: string = "all",
  includeExamples: boolean = true
): Promise<string> {
  try {
    // Get API documentation
    const apiDocs = VapiDocumentation.getDocsByCategory("api");
    const apiEndpoints = VapiDocumentation.getAllApiEndpoints();
    
    // Search for endpoint
    const searchTerm = endpoint.toLowerCase();
    
    // Search in API endpoints
    const relevantEndpoints = apiEndpoints.filter((ep: ApiEndpoint) =>
      ep.path.toLowerCase().includes(searchTerm) ||
      ep.description.toLowerCase().includes(searchTerm) ||
      ep.id.toLowerCase().includes(searchTerm)
    );
    
    // Filter by method if specified
    let filteredEndpoints = relevantEndpoints;
    if (method !== "all") {
      filteredEndpoints = VapiDocumentation.getApiEndpointsByMethod(method);
    }
    
    // Search in API docs
    const relevantDocs = apiDocs.filter((doc: DocItem) =>
      doc.title.toLowerCase().includes(searchTerm) ||
      doc.description.toLowerCase().includes(searchTerm) ||
      doc.content.toLowerCase().includes(searchTerm) ||
      doc.tags.some(tag => tag.toLowerCase().includes(searchTerm))
    );

    if (filteredEndpoints.length === 0 && relevantDocs.length === 0) {
      return `# 🔍 No API Reference Found

No API reference found for "${endpoint}".

## Available API Endpoints:
${apiEndpoints.map((ep: ApiEndpoint) => `- **${ep.method} ${ep.path}** - ${ep.description}`).join('\n')}

## Available API Documentation:
${apiDocs.map((doc: DocItem) => `- **${doc.title}** - ${doc.description}`).join('\n')}

## Popular Endpoints:
- **assistants** - Create and manage voice assistants
- **calls** - Make outbound phone calls  
- **tools** - Manage custom functions
- **phone-numbers** - Manage phone numbers
- **webhooks** - Configure webhook endpoints

Try searching for one of these!`;
    }

    let response = `# 🔧 Vapi API Reference\n\n`;
    response += `API information for "${endpoint}"\n`;
    if (method !== "all") {
      response += `**Method:** ${method}\n`;
    }
    response += `\n`;

    // Show API endpoints first
    if (filteredEndpoints.length > 0) {
      response += `## 🚀 API Endpoints\n\n`;
      
      filteredEndpoints.forEach((ep: ApiEndpoint, index: number) => {
        response += `### ${index + 1}. ${ep.method} ${ep.path}\n\n`;
        response += `${ep.description}\n\n`;
        
        // Request body
        if (ep.requestBody) {
          response += `**Request Body:**\n`;
          response += "```json\n";
          response += JSON.stringify(ep.requestBody, null, 2) + "\n";
          response += "```\n\n";
        }
        
        // Parameters
        if (ep.parameters) {
          response += `**Parameters:**\n`;
          Object.entries(ep.parameters).forEach(([key, value]) => {
            response += `- **${key}**: ${value}\n`;
          });
          response += `\n`;
        }
        
        // Examples
        if (includeExamples && ep.examples) {
          if (ep.examples.request) {
            response += `**Example Request:**\n`;
            response += "```json\n";
            response += ep.examples.request + "\n";
            response += "```\n\n";
          }
          
          if (ep.examples.response) {
            response += `**Example Response:**\n`;
            response += "```json\n";
            response += ep.examples.response + "\n";
            response += "```\n\n";
          }
        }
        
        response += "---\n\n";
      });
    }

    // Show relevant documentation
    if (relevantDocs.length > 0) {
      response += `## 📚 Related Documentation\n\n`;
      
      relevantDocs.forEach((doc: DocItem, index: number) => {
        response += `### ${index + 1}. ${doc.title}\n\n`;
        response += `${doc.description}\n\n`;
        
        // Add content preview (first 500 chars)
        if (doc.content && doc.content.length > 500) {
          response += `**Preview:**\n${doc.content.substring(0, 500)}...\n\n`;
        } else if (doc.content) {
          response += doc.content + "\n\n";
        }
        
        response += `**📖 Full Documentation:** ${doc.url}\n`;
        response += `**📅 Last Updated:** ${doc.lastUpdated}\n\n`;
        
        response += "---\n\n";
      });
    }

    response += `## 🛠️ Additional Resources\n\n`;
    response += `- **[Complete API Reference](https://docs.vapi.ai/api-reference)** - Full API documentation\n`;
    response += `- **[Postman Collection](https://postman.vapi.ai)** - Test APIs directly\n`;
    response += `- **[OpenAPI Spec](https://api.vapi.ai/openapi.json)** - Machine-readable API spec\n`;
    response += `- Use \`get_examples\` for code implementations\n`;
    response += `- Use \`search_documentation\` for general information`;

    return response;
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `❌ Error fetching API reference: ${errorMessage}\n\nPlease try again or visit https://docs.vapi.ai/api-reference`;
  }
} 