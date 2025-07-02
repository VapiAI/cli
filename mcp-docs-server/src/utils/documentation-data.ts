export interface DocItem {
  id: string;
  title: string;
  description: string;
  content: string;
  category: "api" | "guides" | "examples" | "changelog";
  url: string;
  tags: string[];
  lastUpdated: string;
}

export interface CodeExample {
  id: string;
  title: string;
  description: string;
  code: string;
  language: string;
  framework?: string;
  category: string;
  tags: string[];
}

export interface ApiEndpoint {
  id: string;
  path: string;
  method: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
  description: string;
  parameters?: Record<string, any>;
  requestBody?: Record<string, any>;
  responses?: Record<string, any>;
  examples?: {
    request?: string;
    response?: string;
  };
}

export class VapiDocumentation {
  private static docs: DocItem[] = [
    {
      id: "getting-started",
      title: "Getting Started with Vapi",
      description: "Learn how to create your first voice AI assistant with Vapi in minutes",
      content: `# Getting Started with Vapi

Vapi is the voice AI platform that lets you build voice assistants that can make and receive phone calls. Here's how to get started:

## Quick Start

1. Sign up at https://dashboard.vapi.ai
2. Get your API key from the dashboard
3. Create your first assistant
4. Make your first call

## Installation

### JavaScript/TypeScript
\`\`\`bash
npm install @vapi-ai/web
# or
npm install @vapi-ai/server
\`\`\`

### Python
\`\`\`bash
pip install vapi-python
\`\`\`

## Your First Assistant

\`\`\`typescript
import Vapi from '@vapi-ai/web';

const vapi = new Vapi('your-api-key');

const assistant = await vapi.assistants.create({
  name: "My First Assistant",
  voice: {
    provider: "openai",
    voiceId: "alloy"
  },
  model: {
    provider: "openai",
    model: "gpt-3.5-turbo",
    messages: [{
      role: "system",
      content: "You are a helpful assistant."
    }]
  }
});
\`\`\``,
      category: "guides",
      url: "https://docs.vapi.ai/quickstart",
      tags: ["quickstart", "setup", "assistant", "beginner"],
      lastUpdated: "2024-01-15"
    },
    {
      id: "assistants-api",
      title: "Assistants API",
      description: "Create and manage AI voice assistants programmatically",
      content: `# Assistants API

The Assistants API allows you to create, update, and manage your voice AI assistants.

## Creating an Assistant

\`\`\`typescript
const assistant = await vapi.assistants.create({
  name: "Customer Support Bot",
  voice: {
    provider: "elevenlabs",
    voiceId: "21m00Tcm4TlvDq8ikWAM"
  },
  model: {
    provider: "openai",
    model: "gpt-4",
    messages: [{
      role: "system",
      content: "You are a customer support representative."
    }]
  },
  tools: [
    {
      type: "function",
      function: {
        name: "get_order_status",
        description: "Get the status of a customer order"
      }
    }
  ]
});
\`\`\`

## Voice Configuration

Choose from multiple voice providers:
- OpenAI (alloy, echo, fable, onyx, nova, shimmer)
- ElevenLabs (premium voices)
- Azure (neural voices)
- Deepgram (Aura voices)`,
      category: "api",
      url: "https://docs.vapi.ai/api-reference/assistants",
      tags: ["assistants", "api", "voice", "models"],
      lastUpdated: "2024-01-20"
    },
    {
      id: "phone-calls",
      title: "Making Phone Calls",
      description: "Learn how to make outbound phone calls with your voice assistants",
      content: `# Making Phone Calls

Vapi allows you to make outbound phone calls using your voice assistants.

## Outbound Calls

\`\`\`typescript
const call = await vapi.calls.create({
  phoneNumber: "+1234567890",
  assistantId: "assistant-id",
  // Optional: Override assistant settings
  assistantOverrides: {
    voice: {
      provider: "openai",
      voiceId: "alloy"
    }
  }
});
\`\`\`

## Call Status

Monitor your calls in real-time:

\`\`\`typescript
// Get call details
const call = await vapi.calls.get('call-id');

// List all calls
const calls = await vapi.calls.list({
  limit: 100,
  assistantId: 'assistant-id'
});
\`\`\`

## Webhooks

Set up webhooks to receive real-time call events:

\`\`\`typescript
// Configure webhook URL in dashboard
// Receive events: call.started, call.ended, function.called
\`\`\``,
      category: "guides",
      url: "https://docs.vapi.ai/phone-calls",
      tags: ["calls", "outbound", "webhooks", "monitoring"],
      lastUpdated: "2024-01-18"
    },
    {
      id: "tools-functions",
      title: "Tools and Functions",
      description: "Add custom functions and tools to your voice assistants",
      content: `# Tools and Functions

Extend your assistants with custom functions that can interact with external APIs and services.

## Function Calling

\`\`\`typescript
const assistant = await vapi.assistants.create({
  name: "Weather Assistant",
  tools: [
    {
      type: "function",
      function: {
        name: "get_weather",
        description: "Get current weather for a location",
        parameters: {
          type: "object",
          properties: {
            location: {
              type: "string",
              description: "The city and state"
            }
          },
          required: ["location"]
        }
      }
    }
  ]
});
\`\`\`

## Webhook Implementation

Handle function calls via webhooks:

\`\`\`typescript
app.post('/webhook', (req, res) => {
  const { type, functionCall } = req.body;
  
  if (type === 'function-call') {
    const { name, parameters } = functionCall;
    
    if (name === 'get_weather') {
      const weather = getWeatherData(parameters.location);
      res.json({
        result: \`The weather in \${parameters.location} is \${weather}\`
      });
    }
  }
});
\`\`\``,
      category: "guides",
      url: "https://docs.vapi.ai/tools",
      tags: ["tools", "functions", "webhooks", "integration"],
      lastUpdated: "2024-01-22"
    },
    {
      id: "voice-settings",
      title: "Voice Settings and Providers",
      description: "Configure voice settings, providers, and speech parameters",
      content: `# Voice Settings and Providers

Customize the voice experience for your assistants.

## Voice Providers

### OpenAI Voices
\`\`\`typescript
voice: {
  provider: "openai",
  voiceId: "alloy", // alloy, echo, fable, onyx, nova, shimmer
  speed: 1.0,
  emotion: "neutral"
}
\`\`\`

### ElevenLabs Voices
\`\`\`typescript
voice: {
  provider: "elevenlabs",
  voiceId: "21m00Tcm4TlvDq8ikWAM",
  stability: 0.5,
  similarityBoost: 0.75,
  style: 0.0,
  useSpeakerBoost: true
}
\`\`\`

### Azure Voices
\`\`\`typescript
voice: {
  provider: "azure",
  voiceId: "jenny",
  style: "cheerful",
  rate: "medium",
  pitch: "medium"
}
\`\`\`

## Speech Settings

\`\`\`typescript
{
  voice: { /* voice config */ },
  // Interruption handling
  interruptionsEnabled: true,
  responseDelaySeconds: 0.5,
  
  // Background sound
  backgroundSound: "office",
  backgroundDenoisingEnabled: true,
  
  // Call recording
  recordingEnabled: true
}
\`\`\``,
      category: "api",
      url: "https://docs.vapi.ai/voice-settings",
      tags: ["voice", "speech", "providers", "settings"],
      lastUpdated: "2024-01-25"
    }
  ];

  private static examples: CodeExample[] = [
    {
      id: "basic-assistant",
      title: "Basic Voice Assistant",
      description: "Create a simple voice assistant that can answer questions",
      code: `import Vapi from '@vapi-ai/web';

const vapi = new Vapi('your-api-key');

// Create assistant
const assistant = await vapi.assistants.create({
  name: "Basic Assistant",
  voice: {
    provider: "openai",
    voiceId: "alloy"
  },
  model: {
    provider: "openai",
    model: "gpt-3.5-turbo",
    messages: [{
      role: "system",
      content: "You are a helpful assistant that answers questions clearly and concisely."
    }]
  }
});

// Make a call
const call = await vapi.calls.create({
  phoneNumber: "+1234567890",
  assistantId: assistant.id
});

console.log('Call started:', call.id);`,
      language: "typescript",
      framework: "node",
      category: "getting-started",
      tags: ["assistant", "basic", "calls"]
    },
    {
      id: "weather-function",
      title: "Weather Assistant with Function Calling",
      description: "Voice assistant that can get weather information using function calls",
      code: `import Vapi from '@vapi-ai/server';
import express from 'express';

const app = express();
app.use(express.json());

const vapi = new Vapi('your-api-key');

// Create weather assistant
const assistant = await vapi.assistants.create({
  name: "Weather Assistant",
  voice: {
    provider: "openai",
    voiceId: "alloy"
  },
  model: {
    provider: "openai",
    model: "gpt-4",
    messages: [{
      role: "system",
      content: "You are a weather assistant. Use the get_weather function to provide accurate weather information."
    }]
  },
  tools: [{
    type: "function",
    function: {
      name: "get_weather",
      description: "Get current weather for a location",
      parameters: {
        type: "object",
        properties: {
          location: {
            type: "string",
            description: "The city and state or country"
          }
        },
        required: ["location"]
      }
    }
  }]
});

// Webhook to handle function calls
app.post('/webhook', async (req, res) => {
  const { type, functionCall } = req.body;
  
  if (type === 'function-call' && functionCall.name === 'get_weather') {
    const { location } = functionCall.parameters;
    
    // In a real app, call a weather API
    const weather = await getWeatherData(location);
    
    res.json({
      result: \`The current weather in \${location} is \${weather.temperature}°F with \${weather.conditions}.\`
    });
  }
});

async function getWeatherData(location: string) {
  // Simulate weather API call
  return {
    temperature: 72,
    conditions: "sunny"
  };
}

app.listen(3000, () => {
  console.log('Webhook server running on port 3000');
});`,
      language: "typescript",
      framework: "express",
      category: "functions",
      tags: ["functions", "weather", "webhook", "express"]
    }
  ];

  private static apiEndpoints: ApiEndpoint[] = [
    {
      id: "create-assistant",
      path: "/assistants",
      method: "POST",
      description: "Create a new voice assistant",
      requestBody: {
        name: "string",
        voice: "object",
        model: "object",
        tools: "array"
      },
      examples: {
        request: `{
  "name": "Customer Support Bot",
  "voice": {
    "provider": "openai",
    "voiceId": "alloy"
  },
  "model": {
    "provider": "openai",
    "model": "gpt-3.5-turbo",
    "messages": [{
      "role": "system",
      "content": "You are a helpful customer support representative."
    }]
  }
}`,
        response: `{
  "id": "assistant_123",
  "name": "Customer Support Bot",
  "createdAt": "2024-01-15T10:30:00Z",
  "voice": { ... },
  "model": { ... }
}`
      }
    },
    {
      id: "make-call",
      path: "/calls",
      method: "POST",
      description: "Make an outbound phone call",
      requestBody: {
        phoneNumber: "string",
        assistantId: "string",
        assistantOverrides: "object"
      },
      examples: {
        request: `{
  "phoneNumber": "+1234567890",
  "assistantId": "assistant_123"
}`,
        response: `{
  "id": "call_456",
  "phoneNumber": "+1234567890",
  "assistantId": "assistant_123",
  "status": "queued",
  "createdAt": "2024-01-15T10:35:00Z"
}`
      }
    }
  ];

  static getAllDocs(): DocItem[] {
    return this.docs;
  }

  static getDocsByCategory(category: string): DocItem[] {
    return this.docs.filter(doc => doc.category === category);
  }

  static getDocById(id: string): DocItem | undefined {
    return this.docs.find(doc => doc.id === id);
  }

  static getAllExamples(): CodeExample[] {
    return this.examples;
  }

  static getExamplesByLanguage(language: string): CodeExample[] {
    return this.examples.filter(ex => ex.language === language || language === "all");
  }

  static getExamplesByFramework(framework: string): CodeExample[] {
    return this.examples.filter(ex => ex.framework === framework || framework === "all");
  }

  static getAllApiEndpoints(): ApiEndpoint[] {
    return this.apiEndpoints;
  }

  static getApiEndpointsByMethod(method: string): ApiEndpoint[] {
    return this.apiEndpoints.filter(ep => ep.method === method || method === "all");
  }

  static searchApiEndpoints(query: string): ApiEndpoint[] {
    const lowerQuery = query.toLowerCase();
    return this.apiEndpoints.filter(ep => 
      ep.path.toLowerCase().includes(lowerQuery) ||
      ep.description.toLowerCase().includes(lowerQuery)
    );
  }
} 