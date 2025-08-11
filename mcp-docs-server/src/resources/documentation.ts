export class DocumentationSource {
  /**
   * Get resource content by URI
   */
  async getResource(uri: string): Promise<string> {
    try {
      switch (uri) {
        case "vapi://docs/overview":
          return this.getDocumentationOverview();

        case "vapi://docs/quickstart":
          return this.getQuickStartGuide();

        case "vapi://examples/collection":
          return this.getExamplesCollection();

        case "vapi://api/reference":
          return this.getApiReference();

        case "vapi://changelog/latest":
          return this.getLatestChanges();

        default:
          throw new Error(`Unknown resource URI: ${uri}`);
      }
    } catch (error) {
      throw new Error(`Failed to load resource ${uri}: ${error instanceof Error ? error.message : "Unknown error"}`);
    }
  }

  private getDocumentationOverview(): string {
    return `# 📚 Vapi Documentation Overview

Welcome to the comprehensive Vapi documentation! This overview will help you navigate through all available resources.

## 🚀 Getting Started
- **Quick Start Guide** - Get up and running in 5 minutes
- **Installation** - SDK installation for all platforms
- **Authentication** - API key setup and management
- **Your First Call** - Make your first voice AI call

## 🤖 Core Concepts

### Voice Assistants
Create AI assistants that can:
- Make and receive phone calls
- Understand natural language
- Execute custom functions
- Integrate with your business logic

### Key Components
- **Models** - LLM providers (OpenAI, Anthropic, etc.)
- **Voices** - TTS providers (OpenAI, ElevenLabs, Azure)
- **Tools** - Custom functions and API integrations
- **Webhooks** - Real-time event handling

## 📖 Documentation Sections

### 🔧 API Reference
Complete reference for all Vapi APIs:
- **Assistants API** - Create and manage voice assistants
- **Calls API** - Make outbound calls and manage call history
- **Tools API** - Define custom functions
- **Phone Numbers API** - Manage phone numbers
- **Webhooks API** - Configure event notifications

### 💻 Code Examples
Ready-to-use examples:
- **Basic Assistant** - Simple voice assistant setup
- **Function Calling** - Add custom tools and functions
- **Webhook Handling** - Process real-time events
- **Advanced Voice Settings** - Configure voice parameters

### 📝 Guides & Tutorials
Step-by-step implementation guides:
- **Making Phone Calls** - Outbound call setup
- **Voice Configuration** - Voice provider setup
- **Function Integration** - Connect external APIs
- **Error Handling** - Best practices for error management

### 🔄 Changelog
Stay updated with:
- **New Features** - Latest Vapi capabilities
- **Bug Fixes** - Recent issue resolutions
- **Breaking Changes** - Important API updates
- **Migration Guides** - Upgrade instructions

## 🛠️ SDKs & Tools

### Official SDKs
- **JavaScript/TypeScript** - Web and Node.js
- **Python** - Server-side integration
- **Go** - High-performance applications
- **REST API** - Direct HTTP integration

### Development Tools
- **CLI Tool** - Command-line interface
- **Postman Collection** - API testing
- **OpenAPI Specification** - API documentation
- **Webhook Testing** - Local development tools

## 🎯 Common Use Cases

### Customer Support
- Automated phone support
- Call routing and escalation
- FAQ handling
- Appointment scheduling

### Sales & Marketing
- Lead qualification calls
- Follow-up automation
- Product demos
- Survey collection

### Internal Operations
- Meeting scheduling
- Internal notifications
- Data collection
- Process automation

## 🔗 Additional Resources

- **[Dashboard](https://dashboard.vapi.ai)** - Web interface
- **[Community Discord](https://discord.gg/vapi)** - Get help and connect
- **[GitHub](https://github.com/VapiAI)** - Open source examples
- **[Blog](https://vapi.ai/blog)** - Latest updates and tutorials

## 💡 Need Help?

Use the MCP tools available in your IDE:
- \`search_documentation\` - Find specific information
- \`get_examples\` - Get code examples
- \`get_guides\` - Access step-by-step tutorials
- \`get_api_reference\` - API documentation
- \`get_changelog\` - Latest updates

Happy building with Vapi! 🎉`;
  }

  private getQuickStartGuide(): string {
    return `# ⚡ Vapi Quick Start Guide

Get your first voice AI assistant up and running in 5 minutes!

## Step 1: Get Your API Key

1. Sign up at [dashboard.vapi.ai](https://dashboard.vapi.ai)
2. Navigate to the API Keys section
3. Create a new API key and copy it

## Step 2: Install the SDK

Choose your platform:

### JavaScript/TypeScript (Web)
\`\`\`bash
npm install @vapi-ai/web
\`\`\`

### JavaScript/TypeScript (Server)
\`\`\`bash
npm install @vapi-ai/server
\`\`\`

### Python
\`\`\`bash
pip install vapi-python
\`\`\`

## Step 3: Create Your First Assistant

### JavaScript/TypeScript
\`\`\`typescript
import Vapi from '@vapi-ai/web';

const vapi = new Vapi('your-api-key-here');

const assistant = await vapi.assistants.create({
  name: "My First Assistant",
  voice: {
    provider: "openai",
    voiceId: "alloy"
  },
  model: {
    provider: "openai", 
    model: "gpt-4o",
    messages: [{
      role: "system",
      content: "You are a helpful assistant. Be concise and friendly."
    }]
  }
});

console.log('Assistant created:', assistant.id);
\`\`\`

### Python
\`\`\`python
from vapi_python import Vapi

vapi = Vapi(api_key="your-api-key-here")

assistant = vapi.assistants.create(
    name="My First Assistant",
    voice={
        "provider": "openai",
        "voice_id": "alloy"
    },
    model={
        "provider": "openai",
        "model": "gpt-4o", 
        "messages": [{
            "role": "system",
            "content": "You are a helpful assistant. Be concise and friendly."
        }]
    }
)

print(f"Assistant created: {assistant.id}")
\`\`\`

## Step 4: Make Your First Call

### JavaScript/TypeScript
\`\`\`typescript
const call = await vapi.calls.create({
  phoneNumber: "+1234567890",  // Replace with actual number
  assistantId: assistant.id
});

console.log('Call started:', call.id);
\`\`\`

### Python
\`\`\`python
call = vapi.calls.create(
    phone_number="+1234567890",  # Replace with actual number
    assistant_id=assistant.id
)

print(f"Call started: {call.id}")
\`\`\`

## Step 5: Monitor Your Call

Check the call status in your dashboard or via API:

\`\`\`typescript
const callStatus = await vapi.calls.get(call.id);
console.log('Call status:', callStatus.status);
\`\`\`

## 🎉 Congratulations!

You've successfully:
- ✅ Created your first voice assistant
- ✅ Made an outbound phone call
- ✅ Integrated with Vapi's API

## Next Steps

### Add Custom Functions
Make your assistant more powerful with custom tools:

\`\`\`typescript
const assistant = await vapi.assistants.create({
  name: "Enhanced Assistant",
  // ... other config
  tools: [{
    type: "function",
    function: {
      name: "get_weather",
      description: "Get current weather for a location",
      parameters: {
        type: "object",
        properties: {
          location: { type: "string", description: "City name" }
        },
        required: ["location"]
      }
    }
  }]
});
\`\`\`

### Set Up Webhooks
Handle real-time events:

\`\`\`typescript
// Configure webhook URL in dashboard
// Receive events like call.started, call.ended, function.called
\`\`\`

### Explore Voice Options
Try different voice providers:

\`\`\`typescript
voice: {
  provider: "elevenlabs",
  voiceId: "21m00Tcm4TlvDq8ikWAM"  // Premium voice
}
\`\`\`

## 🔗 Learn More

- **[Complete Documentation](https://docs.vapi.ai)** - Full guides
- **[API Reference](https://docs.vapi.ai/api-reference)** - All endpoints
- **[Examples Repository](https://github.com/VapiAI/examples)** - Code samples
- **[Dashboard](https://dashboard.vapi.ai)** - Web interface
- **[Discord Community](https://discord.gg/vapi)** - Get help

Ready to build something amazing? Let's go! 🚀`;
  }

  private getExamplesCollection(): string {
    return `# 💻 Vapi Code Examples Collection

A comprehensive collection of ready-to-use Vapi code examples for all platforms and use cases.

## 🚀 Basic Examples

### Simple Voice Assistant
\`\`\`typescript
import Vapi from '@vapi-ai/web';

const vapi = new Vapi('your-api-key');

// Create basic assistant
const assistant = await vapi.assistants.create({
  name: "Basic Assistant",
  voice: { provider: "openai", voiceId: "alloy" },
  model: {
    provider: "openai",
    model: "gpt-4o",
    messages: [{ role: "system", content: "You are a helpful assistant." }]
  }
});

// Make call
const call = await vapi.calls.create({
  phoneNumber: "+1234567890",
  assistantId: assistant.id
});
\`\`\`

### Customer Support Bot
\`\`\`typescript
const supportBot = await vapi.assistants.create({
  name: "Support Bot",
  voice: { provider: "openai", voiceId: "echo" },
  model: {
    provider: "openai",
    model: "gpt-4o",
    messages: [{
      role: "system",
      content: \`You are a customer support representative for Acme Corp.
      Be helpful, professional, and empathetic. 
      If you need to transfer the call, ask for permission first.\`
    }]
  },
  tools: [{
    type: "function",
    function: {
      name: "lookup_order",
      description: "Look up customer order status",
      parameters: {
        type: "object",
        properties: {
          order_id: { type: "string", description: "Order ID" }
        },
        required: ["order_id"]
      }
    }
  }]
});
\`\`\`

## 🛠️ Advanced Examples

### Function Calling with Webhooks
\`\`\`typescript
import express from 'express';
import Vapi from '@vapi-ai/server';

const app = express();
app.use(express.json());

const vapi = new Vapi('your-api-key');

// Create assistant with functions
const assistant = await vapi.assistants.create({
  name: "Smart Assistant",
  voice: { provider: "elevenlabs", voiceId: "21m00Tcm4TlvDq8ikWAM" },
  model: {
    provider: "openai",
    model: "gpt-4o",
    messages: [{
      role: "system",
      content: "You are a smart assistant that can help with weather and scheduling."
    }]
  },
  tools: [
    {
      type: "function",
      function: {
        name: "get_weather",
        description: "Get current weather",
        parameters: {
          type: "object",
          properties: {
            location: { type: "string" }
          },
          required: ["location"]
        }
      }
    },
    {
      type: "function", 
      function: {
        name: "schedule_meeting",
        description: "Schedule a meeting",
        parameters: {
          type: "object",
          properties: {
            title: { type: "string" },
            date: { type: "string" },
            duration: { type: "number" }
          },
          required: ["title", "date"]
        }
      }
    }
  ]
});

// Webhook handler
app.post('/webhook', async (req, res) => {
  const { type, functionCall, call } = req.body;
  
  if (type === 'function-call') {
    const { name, parameters } = functionCall;
    
    try {
      let result;
      
      switch (name) {
        case 'get_weather':
          result = await getWeather(parameters.location);
          break;
          
        case 'schedule_meeting':
          result = await scheduleMeeting(parameters);
          break;
          
        default:
          result = "Unknown function";
      }
      
      res.json({ result });
    } catch (error) {
      res.json({ error: error.message });
    }
  }
});

async function getWeather(location) {
  // Weather API integration
  return \`The weather in \${location} is sunny and 72°F\`;
}

async function scheduleMeeting({ title, date, duration = 30 }) {
  // Calendar API integration  
  return \`Meeting "\${title}" scheduled for \${date} (\${duration} minutes)\`;
}

app.listen(3000);
\`\`\`

### Real-time Call Monitoring
\`\`\`typescript
// Set up webhook for call events
app.post('/webhook', (req, res) => {
  const { type, call, message } = req.body;
  
  switch (type) {
    case 'call-started':
      console.log(\`Call \${call.id} started with \${call.phoneNumber}\`);
      // Send to monitoring dashboard
      break;
      
    case 'call-ended':
      console.log(\`Call \${call.id} ended. Duration: \${call.duration}s\`);
      // Log call data
      break;
      
    case 'speech-update':
      console.log(\`Transcript: \${message.transcript}\`);
      // Real-time transcription processing
      break;
      
    case 'function-call':
      console.log(\`Function called: \${message.functionCall.name}\`);
      // Function execution logging
      break;
  }
  
  res.sendStatus(200);
});
\`\`\`

## 🐍 Python Examples

### Flask Integration
\`\`\`python
from flask import Flask, request, jsonify
from vapi_python import Vapi

app = Flask(__name__)
vapi = Vapi(api_key="your-api-key")

@app.route('/create-assistant', methods=['POST'])
def create_assistant():
    data = request.json
    
    assistant = vapi.assistants.create(
        name=data['name'],
        voice={
            "provider": "openai",
            "voice_id": "alloy"
        },
        model={
            "provider": "openai",
            "model": "gpt-4o",
            "messages": [{
                "role": "system", 
                "content": data['system_prompt']
            }]
        }
    )
    
    return jsonify({"assistant_id": assistant.id})

@app.route('/make-call', methods=['POST'])
def make_call():
    data = request.json
    
    call = vapi.calls.create(
        phone_number=data['phone_number'],
        assistant_id=data['assistant_id']
    )
    
    return jsonify({"call_id": call.id})

@app.route('/webhook', methods=['POST'])
def webhook():
    data = request.json
    
    if data['type'] == 'function-call':
        function_name = data['functionCall']['name']
        parameters = data['functionCall']['parameters']
        
        # Handle function calls
        if function_name == 'get_user_data':
            result = get_user_data(parameters['user_id'])
            return jsonify({"result": result})
    
    return '', 200

def get_user_data(user_id):
    # Database lookup
    return f"User {user_id} data retrieved"

if __name__ == '__main__':
    app.run(debug=True)
\`\`\`

### Django Integration
\`\`\`python
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_http_methods
import json
from vapi_python import Vapi

vapi = Vapi(api_key="your-api-key")

@csrf_exempt
@require_http_methods(["POST"])
def webhook(request):
    data = json.loads(request.body)
    
    if data['type'] == 'call-started':
        # Log call start
        call_id = data['call']['id']
        phone_number = data['call']['phoneNumber']
        
        # Save to database
        Call.objects.create(
            vapi_call_id=call_id,
            phone_number=phone_number,
            status='started'
        )
    
    elif data['type'] == 'call-ended':
        # Update call record
        call_id = data['call']['id']
        duration = data['call']['duration']
        
        Call.objects.filter(vapi_call_id=call_id).update(
            status='ended',
            duration=duration
        )
    
    return JsonResponse({'status': 'ok'})
\`\`\`

## 🌐 Frontend Examples

### React Integration
\`\`\`tsx
import React, { useState } from 'react';
import Vapi from '@vapi-ai/web';

const VapiDemo = () => {
  const [vapi] = useState(() => new Vapi('your-public-key'));
  const [isCallActive, setIsCallActive] = useState(false);
  
  const startCall = async () => {
    try {
      await vapi.start('assistant-id');
      setIsCallActive(true);
    } catch (error) {
      console.error('Failed to start call:', error);
    }
  };
  
  const endCall = () => {
    vapi.stop();
    setIsCallActive(false);
  };
  
  return (
    <div>
      <h1>Vapi Voice Demo</h1>
      
      {!isCallActive ? (
        <button onClick={startCall}>
          🎤 Start Voice Call
        </button>
      ) : (
        <button onClick={endCall}>
          ⏹️ End Call
        </button>
      )}
      
      {isCallActive && (
        <div>
          <p>Call is active...</p>
          <div className="audio-visualizer">
            {/* Add audio visualizer */}
          </div>
        </div>
      )}
    </div>
  );
};

export default VapiDemo;
\`\`\`

## 📱 React Native Example
\`\`\`typescript
import React, { useState, useEffect } from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import Vapi from '@vapi-ai/web';

const VapiMobileDemo = () => {
  const [vapi] = useState(() => new Vapi('your-public-key'));
  const [isCallActive, setIsCallActive] = useState(false);
  const [transcript, setTranscript] = useState('');
  
  useEffect(() => {
    vapi.on('speech-start', () => console.log('Speech started'));
    vapi.on('speech-end', () => console.log('Speech ended'));
    vapi.on('message', (message) => {
      if (message.type === 'transcript') {
        setTranscript(message.transcript);
      }
    });
    
    return () => vapi.removeAllListeners();
  }, []);
  
  const startCall = async () => {
    try {
      await vapi.start('assistant-id');
      setIsCallActive(true);
    } catch (error) {
      console.error('Call failed:', error);
    }
  };
  
  const endCall = () => {
    vapi.stop();
    setIsCallActive(false);
    setTranscript('');
  };
  
  return (
    <View style={styles.container}>
      <Text style={styles.title}>Vapi Mobile Demo</Text>
      
      <TouchableOpacity 
        style={[styles.button, isCallActive && styles.activeButton]}
        onPress={isCallActive ? endCall : startCall}
      >
        <Text style={styles.buttonText}>
          {isCallActive ? '⏹️ End Call' : '🎤 Start Call'}
        </Text>
      </TouchableOpacity>
      
      {transcript && (
        <View style={styles.transcriptContainer}>
          <Text style={styles.transcriptLabel}>Live Transcript:</Text>
          <Text style={styles.transcript}>{transcript}</Text>
        </View>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: { flex: 1, padding: 20, justifyContent: 'center' },
  title: { fontSize: 24, fontWeight: 'bold', textAlign: 'center', marginBottom: 30 },
  button: { backgroundColor: '#007AFF', padding: 15, borderRadius: 10, marginBottom: 20 },
  activeButton: { backgroundColor: '#FF3B30' },
  buttonText: { color: 'white', textAlign: 'center', fontSize: 18, fontWeight: 'bold' },
  transcriptContainer: { marginTop: 20, padding: 15, backgroundColor: '#f5f5f5', borderRadius: 10 },
  transcriptLabel: { fontWeight: 'bold', marginBottom: 5 },
  transcript: { fontSize: 16, lineHeight: 22 }
});

export default VapiMobileDemo;
\`\`\`

## 🔗 Additional Resources

- **[GitHub Examples](https://github.com/VapiAI/examples)** - Complete example projects
- **[API Reference](https://docs.vapi.ai/api-reference)** - Detailed API docs  
- **[SDKs](https://docs.vapi.ai/sdks)** - All platform SDKs
- **[Tutorials](https://docs.vapi.ai/tutorials)** - Step-by-step guides

Happy coding! 🚀`;
  }

  private getApiReference(): string {
    return `# 🔧 Vapi API Reference

Complete reference for all Vapi REST API endpoints.

## Base URL
\`https://api.vapi.ai\`

## Authentication
Include your API key in the Authorization header:
\`Authorization: Bearer your-api-key\`

## 🤖 Assistants API

### Create Assistant
\`POST /assistants\`

Create a new voice assistant.

**Request Body:**
\`\`\`json
{
  "name": "string",
  "voice": {
    "provider": "openai|elevenlabs|azure|deepgram",
    "voiceId": "string"
  },
  "model": {
    "provider": "openai|anthropic|groq",
    "model": "string",
    "messages": [
      {
        "role": "system|user|assistant",
        "content": "string"
      }
    ]
  }
}
\`\`\`

**Response:**
\`\`\`json
{
  "id": "assistant_123",
  "name": "string",
  "voice": { ... },
  "model": { ... },
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
\`\`\`

### Get Assistant
\`GET /assistants/{id}\`

Retrieve assistant details.

### List Assistants
\`GET /assistants\`

**Query Parameters:**
- \`limit\` (optional): Number of results (default: 10, max: 100)
- \`offset\` (optional): Pagination offset

### Update Assistant
\`PUT /assistants/{id}\`

Update assistant configuration.

### Delete Assistant
\`DELETE /assistants/{id}\`

Delete an assistant.

## 📞 Calls API

### Create Call
\`POST /calls\`

Make an outbound phone call.

**Request Body:**
\`\`\`json
{
  "phoneNumber": "+1234567890",
  "assistantId": "assistant_123",
  "assistantOverrides": {
    "voice": {
      "provider": "openai",
      "voiceId": "alloy"
    }
  },
  "metadata": {
    "customField": "value"
  }
}
\`\`\`

**Response:**
\`\`\`json
{
  "id": "call_456",
  "phoneNumber": "+1234567890",
  "assistantId": "assistant_123",
  "status": "queued|ringing|in-progress|ended|failed",
  "startedAt": "2024-01-15T10:35:00Z",
  "endedAt": null,
  "duration": null,
  "cost": null,
  "metadata": { ... }
}
\`\`\`

### Get Call
\`GET /calls/{id}\`

Get call details including transcript and recordings.

### List Calls
\`GET /calls\`

**Query Parameters:**
- \`limit\` (optional): Number of results
- \`assistantId\` (optional): Filter by assistant
- \`status\` (optional): Filter by status
- \`startDate\` (optional): Filter by date range
- \`endDate\` (optional): Filter by date range

### End Call
\`POST /calls/{id}/end\`

End an active call.

## 🛠️ Tools API

### Create Tool
\`POST /tools\`

Create a custom function tool.

**Request Body:**
\`\`\`json
{
  "name": "get_weather",
  "description": "Get current weather for a location",
  "parameters": {
    "type": "object",
    "properties": {
      "location": {
        "type": "string",
        "description": "City and state"
      },
      "units": {
        "type": "string",
        "enum": ["celsius", "fahrenheit"],
        "default": "fahrenheit"
      }
    },
    "required": ["location"]
  },
  "webhookUrl": "https://yourapi.com/webhook"
}
\`\`\`

### Get Tool
\`GET /tools/{id}\`

### List Tools
\`GET /tools\`

### Update Tool
\`PUT /tools/{id}\`

### Delete Tool
\`DELETE /tools/{id}\`

## 📱 Phone Numbers API

### List Phone Numbers
\`GET /phone-numbers\`

List all your phone numbers.

### Get Phone Number
\`GET /phone-numbers/{id}\`

### Purchase Phone Number
\`POST /phone-numbers\`

**Request Body:**
\`\`\`json
{
  "areaCode": "415",
  "country": "US"
}
\`\`\`

### Release Phone Number
\`DELETE /phone-numbers/{id}\`

## 🔗 Webhooks API

### Create Webhook
\`POST /webhooks\`

**Request Body:**
\`\`\`json
{
  "url": "https://yourapi.com/webhook",
  "events": ["call.started", "call.ended", "function.called"],
  "secret": "optional-secret-for-verification"
}
\`\`\`

### List Webhooks
\`GET /webhooks\`

### Update Webhook
\`PUT /webhooks/{id}\`

### Delete Webhook
\`DELETE /webhooks/{id}\`

### Test Webhook
\`POST /webhooks/{id}/test\`

## 📊 Analytics API

### Get Call Analytics
\`GET /analytics/calls\`

**Query Parameters:**
- \`startDate\`: Start date (ISO 8601)
- \`endDate\`: End date (ISO 8601)
- \`groupBy\`: Group results (day|week|month)
- \`assistantId\`: Filter by assistant

**Response:**
\`\`\`json
{
  "totalCalls": 150,
  "totalDuration": 7200,
  "averageDuration": 48,
  "successRate": 0.95,
  "totalCost": 25.50,
  "data": [
    {
      "date": "2024-01-15",
      "calls": 25,
      "duration": 1200,
      "cost": 4.25
    }
  ]
}
\`\`\`

## 🔄 Webhook Events

### Call Events
- \`call.started\` - Call initiated
- \`call.ringing\` - Phone is ringing
- \`call.answered\` - Call answered
- \`call.ended\` - Call completed
- \`call.failed\` - Call failed

### Speech Events
- \`speech.started\` - User started speaking
- \`speech.ended\` - User stopped speaking
- \`transcript.partial\` - Partial transcript
- \`transcript.final\` - Final transcript

### Function Events
- \`function.called\` - Function was invoked
- \`function.result\` - Function completed

### Example Webhook Payload
\`\`\`json
{
  "type": "call.started",
  "timestamp": "2024-01-15T10:35:00Z",
  "call": {
    "id": "call_456",
    "phoneNumber": "+1234567890",
    "assistantId": "assistant_123",
    "status": "in-progress"
  }
}
\`\`\`

## ❌ Error Handling

All errors return appropriate HTTP status codes with error details:

\`\`\`json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Phone number is required",
    "details": {
      "field": "phoneNumber",
      "reason": "missing_required_field"
    }
  }
}
\`\`\`

### Common Error Codes
- \`400\` - Bad Request
- \`401\` - Unauthorized (invalid API key)
- \`403\` - Forbidden (insufficient permissions)
- \`404\` - Not Found
- \`429\` - Rate Limited
- \`500\` - Internal Server Error

## 📈 Rate Limits

- **API Calls**: 1000 requests per minute
- **Outbound Calls**: 100 concurrent calls
- **Webhook Deliveries**: 10,000 per hour

Rate limit headers included in responses:
- \`X-RateLimit-Limit\`
- \`X-RateLimit-Remaining\`
- \`X-RateLimit-Reset\`

## 🔗 Additional Resources

- **[Postman Collection](https://postman.vapi.ai)** - Test APIs
- **[OpenAPI Spec](https://api.vapi.ai/openapi.json)** - Machine-readable
- **[SDKs](https://docs.vapi.ai/sdks)** - Official client libraries
- **[Examples](https://github.com/VapiAI/examples)** - Code samples`;
  }

  private getLatestChanges(): string {
    return `# 📋 Latest Vapi Changes

## ✨ v1.8.0 - Enhanced Voice Settings & Deepgram Aura Support
**Released: January 25, 2024**

### New Features
- **Deepgram Aura Support**: Added support for Deepgram's new Aura voice models
- **Voice Interruption Controls**: New settings for managing conversation flow
- **Background Noise Cancellation**: Enhanced noise detection and cancellation
- **Custom Voice Cloning**: Support for ElevenLabs voice cloning

### Improvements  
- Better voice quality on mobile networks
- Faster response times for voice generation
- Enhanced audio processing pipeline

---

## 🐛 v1.7.5 - Bug Fixes & Performance Improvements  
**Released: January 22, 2024**

### Bug Fixes
- Fixed webhook delivery reliability issues
- Resolved call connection stability problems
- Fixed audio quality issues on mobile networks
- Resolved function calling timeout errors

### Performance
- 30% faster API response times
- Improved memory usage for long calls
- Better error recovery mechanisms

---

## ✨ v1.7.0 - Advanced Function Calling & Tool Integration
**Released: January 18, 2024**

### New Features
- **Parallel Function Calling**: Execute multiple functions simultaneously
- **Enhanced Parameter Validation**: Better type checking and validation
- **Webhook Retry Mechanism**: Exponential backoff for failed deliveries
- **Streaming Responses**: Support for streaming function responses

### API Changes
- New \`parallel\` parameter for function tools
- Enhanced webhook payload format
- Additional validation options

---

## 🐛 v1.6.8 - Call Management Improvements
**Released: January 15, 2024**

### Bug Fixes
- Fixed call status tracking inconsistencies
- Improved error messages for failed calls
- Enhanced call recording quality
- Fixed timezone issues in call logs

### Dashboard Updates
- New real-time call monitoring
- Improved call analytics
- Better error reporting

---

## ✨ v1.6.5 - Multi-Language Support & Localization
**Released: January 12, 2024**

### New Features
- **15+ Languages Supported**: Expanded language support
- **Automatic Language Detection**: Detect caller language automatically
- **Enhanced Accent Handling**: Better voice recognition for accents
- **Localized Messages**: Error messages and prompts in multiple languages

### Supported Languages
- English, Spanish, French, German, Italian
- Portuguese, Dutch, Russian, Japanese, Korean
- Chinese (Simplified & Traditional), Arabic, Hindi

---

## ⚠️ v1.6.0 - API V2 Release (Breaking Changes)
**Released: January 8, 2024**

### Breaking Changes
- **New API Structure**: Updated RESTful API design
- **Authentication Changes**: New refresh token system
- **Error Format Updates**: Standardized error responses

### Migration Required
- Update authentication implementation
- Modify error handling logic
- Review API endpoint changes

### Migration Guide
See [Migration Guide](https://docs.vapi.ai/migrations/v2) for detailed instructions.

---

## 🔗 Stay Updated

- **[Complete Changelog](https://docs.vapi.ai/changelog)** - Full version history
- **[Migration Guides](https://docs.vapi.ai/migrations)** - Upgrade instructions
- **[Breaking Changes](https://docs.vapi.ai/breaking-changes)** - Important updates
- **[Discord](https://discord.gg/vapi)** - Community updates
- **[GitHub Releases](https://github.com/VapiAI/releases)** - Release notes

## 📬 Release Notifications

Get notified about new releases:
- Watch our [GitHub repository](https://github.com/VapiAI/vapi)
- Join our [Discord community](https://discord.gg/vapi)
- Follow [@VapiAI](https://twitter.com/VapiAI) on Twitter
- Subscribe to release notifications in your dashboard`;
  }
} 