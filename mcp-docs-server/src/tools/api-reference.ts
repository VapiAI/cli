import { Tool } from '@modelcontextprotocol/sdk/types.js';
import axios from 'axios';

interface OpenAPISpec {
  openapi: string;
  paths: Record<string, Record<string, any>>;
  components?: {
    schemas?: Record<string, any>;
  };
}

interface APIEndpoint {
  path: string;
  method: string;
  operationId?: string;
  summary?: string;
  description?: string;
  parameters?: any[];
  requestBody?: any;
  responses?: Record<string, any>;
  tags?: string[];
}

export const apiReferenceTool: Tool = {
  name: 'get_api_reference',
  description: 'Get detailed API reference information for Vapi endpoints using the actual OpenAPI specification',
  inputSchema: {
    type: 'object',
    properties: {
      endpoint: {
        type: 'string',
        description: 'API endpoint or resource to get reference for (e.g., "assistants", "calls", "phone-numbers")'
      },
      method: {
        type: 'string',
        enum: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'all'],
        default: 'all',
        description: 'HTTP method (optional)'
      },
      includeExamples: {
        type: 'boolean',
        default: true,
        description: 'Include request/response examples'
      }
    },
    required: ['endpoint']
  }
};

export async function handleApiReference(args: any): Promise<string> {
  const { endpoint, method = 'all', includeExamples = true } = args;

  try {
    // Fetch the OpenAPI spec
    const response = await axios.get('https://api.vapi.ai/api-json');
    const spec: OpenAPISpec = response.data;
    
    // Find matching endpoints
    const matchingEndpoints = findMatchingEndpoints(spec, endpoint, method);
    
    if (matchingEndpoints.length === 0) {
      return generateNoResultsResponse(endpoint, method);
    }

    // Generate comprehensive API reference
    return generateApiReference(matchingEndpoints, spec, includeExamples);
    
  } catch (error) {
    console.error('Failed to fetch OpenAPI spec:', error);
    return generateErrorResponse(endpoint);
  }
}

function findMatchingEndpoints(spec: OpenAPISpec, endpoint: string, method: string): APIEndpoint[] {
  const endpoints: APIEndpoint[] = [];
  const searchTerm = endpoint.toLowerCase();
  
  for (const [path, pathMethods] of Object.entries(spec.paths)) {
    for (const [httpMethod, operation] of Object.entries(pathMethods)) {
      // Skip if method filter doesn't match
      if (method !== 'all' && httpMethod.toUpperCase() !== method.toUpperCase()) {
        continue;
      }
      
      // Check if endpoint matches path, operationId, summary, or tags
      const pathMatch = path.toLowerCase().includes(searchTerm);
      const operationMatch = operation.operationId?.toLowerCase().includes(searchTerm);
      const summaryMatch = operation.summary?.toLowerCase().includes(searchTerm);
      const tagMatch = operation.tags?.some((tag: string) => 
        tag.toLowerCase().includes(searchTerm)
      );
      
      if (pathMatch || operationMatch || summaryMatch || tagMatch) {
        endpoints.push({
          path,
          method: httpMethod.toUpperCase(),
          operationId: operation.operationId,
          summary: operation.summary,
          description: operation.description,
          parameters: operation.parameters,
          requestBody: operation.requestBody,
          responses: operation.responses,
          tags: operation.tags
        });
      }
    }
  }
  
  return endpoints;
}

function generateApiReference(endpoints: APIEndpoint[], spec: OpenAPISpec, includeExamples: boolean): string {
  let result = `# 🔧 API Reference for "${endpoints[0]?.path.split('/')[1] || 'endpoint'}"\n\n`;
  
  result += `Found ${endpoints.length} endpoint(s)\n\n`;
  
  endpoints.forEach((endpoint, index) => {
    result += `## 📄 ${index + 1}. ${endpoint.method} ${endpoint.path}\n\n`;
    
    if (endpoint.summary) {
      result += `**Summary:** ${endpoint.summary}\n\n`;
    }
    
    if (endpoint.description) {
      result += `**Description:** ${endpoint.description}\n\n`;
    }
    
    if (endpoint.operationId) {
      result += `**Operation ID:** \`${endpoint.operationId}\`\n\n`;
    }
    
    if (endpoint.tags && endpoint.tags.length > 0) {
      result += `**Tags:** ${endpoint.tags.join(', ')}\n\n`;
    }
    
    // Parameters
    if (endpoint.parameters && endpoint.parameters.length > 0) {
      result += `### Parameters\n\n`;
      endpoint.parameters.forEach((param: any) => {
        result += `- **${param.name}** (${param.in})`;
        if (param.required) result += ` *required*`;
        result += `\n`;
        if (param.description) result += `  - ${param.description}\n`;
        if (param.schema?.type) result += `  - Type: \`${param.schema.type}\`\n`;
        if (param.schema?.enum) result += `  - Allowed values: \`${param.schema.enum.join('`, `')}\`\n`;
        result += `\n`;
      });
    }
    
    // Request Body
    if (endpoint.requestBody) {
      result += `### Request Body\n\n`;
      if (endpoint.requestBody.description) {
        result += `${endpoint.requestBody.description}\n\n`;
      }
      
      if (endpoint.requestBody.content) {
        const contentTypes = Object.keys(endpoint.requestBody.content);
        result += `**Content Types:** ${contentTypes.join(', ')}\n\n`;
        
        // Show schema for application/json if available
        const jsonContent = endpoint.requestBody.content['application/json'];
        if (jsonContent?.schema && includeExamples) {
          result += `**Schema:**\n\`\`\`json\n${JSON.stringify(jsonContent.schema, null, 2)}\`\`\`\n\n`;
        }
      }
    }
    
    // Responses
    if (endpoint.responses) {
      result += `### Responses\n\n`;
      Object.entries(endpoint.responses).forEach(([statusCode, response]: [string, any]) => {
        result += `**${statusCode}** - ${response.description || 'Success'}\n`;
        
        if (response.content && includeExamples) {
          const contentTypes = Object.keys(response.content);
          if (contentTypes.length > 0) {
            result += `  - Content Types: ${contentTypes.join(', ')}\n`;
          }
        }
        result += `\n`;
      });
    }
    
    // Add link to interactive API
    result += `### 🔗 Try it out\n\n`;
    result += `**Interactive API:** https://api.vapi.ai/api#${endpoint.operationId || endpoint.method.toLowerCase() + endpoint.path.replace(/[{}]/g, '')}\n\n`;
    
    if (index < endpoints.length - 1) {
      result += `---\n\n`;
    }
  });
  
  // Footer with additional resources
  result += `## 🎯 Additional Resources\n\n`;
  result += `- **Full API Reference:** https://docs.vapi.ai/api-reference\n`;
  result += `- **Interactive API Explorer:** https://api.vapi.ai/api\n`;
  result += `- **OpenAPI Spec:** https://api.vapi.ai/api-json\n`;
  result += `- **Dashboard:** https://dashboard.vapi.ai\n`;
  result += `- **Discord Community:** https://discord.gg/vapi\n\n`;
  
  result += `💡 **Pro Tip:** Use the interactive API explorer to test endpoints with your actual API key and see real request/response examples.\n`;
  
  return result;
}

function generateNoResultsResponse(endpoint: string, method: string): string {
  return `# 🔍 No API Reference Found\n\n` +
    `No API endpoints found for "${endpoint}" with method "${method}".\n\n` +
    `## 💡 Suggestions:\n` +
    `- Try broader search terms (e.g., "call" instead of "voice-call")\n` +
    `- Use "all" for method to see all HTTP methods\n` +
    `- Check for typos in your search query\n` +
    `- Try these common endpoints: assistants, calls, phone-numbers, tools, webhooks\n\n` +
    `## 📚 Available Resources:\n` +
    `- **Full API Reference:** https://docs.vapi.ai/api-reference\n` +
    `- **Interactive API:** https://api.vapi.ai/api\n` +
    `- **OpenAPI Spec:** https://api.vapi.ai/api-json\n`;
}

function generateErrorResponse(endpoint: string): string {
  return `# ❌ API Reference Error\n\n` +
    `Sorry, there was an error fetching the API reference for "${endpoint}".\n\n` +
    `## 🔗 Alternative Resources:\n` +
    `- **Full API Reference:** https://docs.vapi.ai/api-reference\n` +
    `- **Interactive API:** https://api.vapi.ai/api\n` +
    `- **OpenAPI Spec:** https://api.vapi.ai/api-json\n` +
    `- **Dashboard:** https://dashboard.vapi.ai\n\n` +
    `Please try again later or visit the links above for complete API documentation.`;
} 