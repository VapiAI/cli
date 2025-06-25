/*
Copyright © 2025 Vapi, Inc.

Licensed under the MIT License (the "License");
you may not use this file except in compliance with the License.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

Authors:
	Dan Goosewin <dan@vapi.ai>
*/

package integrations

import (
	"fmt"
	"os"
	"path/filepath"
)

// GenerateNodeIntegration generates SDK examples and configuration for Node.js/TypeScript projects
func GenerateNodeIntegration(projectPath string, info *ProjectInfo) error {
	examplesDir := filepath.Join(projectPath, "vapi-examples")
	if err := os.MkdirAll(examplesDir, 0o755); err != nil {
		return fmt.Errorf("failed to create examples directory: %w", err)
	}

	if info.IsTypeScript {
		if err := generateNodeTSBasicExample(examplesDir); err != nil {
			return err
		}
		if err := generateNodeTSExpressExample(examplesDir); err != nil {
			return err
		}
		if err := generateNodeTSFastifyExample(examplesDir); err != nil {
			return err
		}
	} else {
		if err := generateNodeJSBasicExample(examplesDir); err != nil {
			return err
		}
		if err := generateNodeJSExpressExample(examplesDir); err != nil {
			return err
		}
	}

	// Generate environment template
	if err := generateNodeEnvTemplate(projectPath); err != nil {
		return err
	}

	return nil
}

func generateNodeJSBasicExample(dir string) error {
	content := `/**
 * Basic Vapi Node.js SDK Example
 * 
 * This example demonstrates how to use the Vapi Node.js SDK to manage assistants
 * and make phone calls programmatically.
 */

require('dotenv').config();
const Vapi = require('@vapi-ai/server-sdk').Vapi;

// Initialize Vapi client
const client = new Vapi({
  apiKey: process.env.VAPI_API_KEY,
});

async function listAssistants() {
  try {
    const assistants = await client.assistants.list();
    console.log(` + "`Found ${assistants.length} assistants:`" + `);
    assistants.forEach(assistant => {
      console.log(` + "`  - ${assistant.name} (ID: ${assistant.id})`" + `);
    });
    return assistants;
  } catch (error) {
    console.error('Error listing assistants:', error);
    return [];
  }
}

async function createAssistant() {
  try {
    const assistant = await client.assistants.create({
      name: 'Node.js SDK Assistant',
      model: {
        provider: 'openai',
        model: 'gpt-4',
        messages: [
          {
            role: 'system',
            content: 'You are a helpful assistant created via the Node.js SDK.'
          }
        ]
      },
      voice: {
        provider: 'elevenlabs',
        voiceId: '21m00Tcm4TlvDq8ikWAM' // Rachel voice
      }
    });
    console.log(` + "`Created assistant: ${assistant.name} (ID: ${assistant.id})`" + `);
    return assistant;
  } catch (error) {
    console.error('Error creating assistant:', error);
    return null;
  }
}

async function makePhoneCall(assistantId, phoneNumber) {
  try {
    const call = await client.calls.create({
      assistantId,
      customer: {
        number: phoneNumber
      }
    });
    console.log(` + "`Call initiated! Call ID: ${call.id}`" + `);
    console.log(` + "`Status: ${call.status}`" + `);
    return call;
  } catch (error) {
    console.error('Error making call:', error);
    return null;
  }
}

async function main() {
  console.log('🚀 Vapi Node.js SDK Example');
  console.log('-'.repeat(50));
  
  // List existing assistants
  console.log('\n📋 Listing assistants...');
  const assistants = await listAssistants();
  
  // Create a new assistant
  console.log('\n✨ Creating a new assistant...');
  const newAssistant = await createAssistant();
  
  if (newAssistant) {
    // Example of making a call (commented out for safety)
    // console.log('\n📞 Making a phone call...');
    // await makePhoneCall(newAssistant.id, '+1234567890');
    console.log('\n💡 To make a phone call, uncomment the code above and provide a valid phone number');
  }
  
  console.log('\n✅ Example completed!');
}

// Run the example
main().catch(console.error);
`

	return os.WriteFile(filepath.Join(dir, "basic-example.js"), []byte(content), 0o644)
}

func generateNodeJSExpressExample(dir string) error {
	content := `/**
 * Express + Vapi Integration Example
 * 
 * This example shows how to integrate Vapi with an Express.js application
 * to handle webhooks and create web-based voice interfaces.
 */

require('dotenv').config();
const express = require('express');
const { Vapi } = require('@vapi-ai/server-sdk');

const app = express();
app.use(express.json());

// Initialize Vapi client
const client = new Vapi({
  apiKey: process.env.VAPI_API_KEY,
});

// Serve the main page
app.get('/', (req, res) => {
  const html = ` + "`" + `
<!DOCTYPE html>
<html>
<head>
    <title>Vapi Express Example</title>
    <script src="https://cdn.jsdelivr.net/npm/@vapi-ai/web@latest/dist/vapi.js"></script>
</head>
<body>
    <h1>Vapi Voice Assistant (Express)</h1>
    <button id="startCall">Start Call</button>
    <button id="endCall" disabled>End Call</button>
    <div id="status"></div>
    
    <script>
        const vapi = new Vapi("${process.env.VAPI_PUBLIC_KEY}");
        const assistantId = "${process.env.VAPI_ASSISTANT_ID}";
        
        document.getElementById('startCall').onclick = async () => {
            try {
                await vapi.start(assistantId);
                document.getElementById('startCall').disabled = true;
                document.getElementById('endCall').disabled = false;
                document.getElementById('status').innerText = 'Call started';
            } catch (error) {
                console.error('Error starting call:', error);
            }
        };
        
        document.getElementById('endCall').onclick = async () => {
            try {
                await vapi.stop();
                document.getElementById('startCall').disabled = false;
                document.getElementById('endCall').disabled = true;
                document.getElementById('status').innerText = 'Call ended';
            } catch (error) {
                console.error('Error ending call:', error);
            }
        };
        
        vapi.on('message', (message) => {
            console.log('Vapi message:', message);
        });
    </script>
</body>
</html>
` + "`" + `;
  res.send(html);
});

// API endpoint to list assistants
app.get('/api/assistants', async (req, res) => {
  try {
    const assistants = await client.assistants.list();
    res.json(assistants.map(a => ({
      id: a.id,
      name: a.name
    })));
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// API endpoint to create a call
app.post('/api/calls', async (req, res) => {
  try {
    const { assistantId, phoneNumber } = req.body;
    const call = await client.calls.create({
      assistantId,
      customer: {
        number: phoneNumber
      }
    });
    res.json({
      id: call.id,
      status: call.status
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// Webhook endpoint for Vapi events
app.post('/webhook/vapi', (req, res) => {
  const { type } = req.body;
  console.log(` + "`Received Vapi webhook: ${type}`" + `);
  
  switch (type) {
    case 'call.started':
      const callStarted = req.body.call;
      console.log(` + "`Call started: ${callStarted.id}`" + `);
      break;
      
    case 'call.ended':
      const callEnded = req.body.call;
      console.log(` + "`Call ended: ${callEnded.id}`" + `);
      break;
      
    case 'transcript':
      console.log(` + "`Transcript: ${req.body.transcript}`" + `);
      break;
  }
  
  res.json({ status: 'ok' });
});

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({ status: 'healthy', service: 'vapi-express-example' });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(` + "`🚀 Server running on http://localhost:${PORT}`" + `);
});
`

	return os.WriteFile(filepath.Join(dir, "express-example.js"), []byte(content), 0o644)
}

func generateNodeTSBasicExample(dir string) error {
	content := `/**
 * Basic Vapi TypeScript SDK Example
 * 
 * This example demonstrates how to use the Vapi TypeScript SDK to manage assistants
 * and make phone calls programmatically.
 */

import 'dotenv/config';
import { Vapi } from '@vapi-ai/server-sdk';

// Initialize Vapi client
const client = new Vapi({
  apiKey: process.env.VAPI_API_KEY!,
});

async function listAssistants() {
  try {
    const assistants = await client.assistants.list();
    console.log(` + "`Found ${assistants.length} assistants:`" + `);
    assistants.forEach(assistant => {
      console.log(` + "`  - ${assistant.name} (ID: ${assistant.id})`" + `);
    });
    return assistants;
  } catch (error) {
    console.error('Error listing assistants:', error);
    return [];
  }
}

async function createAssistant() {
  try {
    const assistant = await client.assistants.create({
      name: 'TypeScript SDK Assistant',
      model: {
        provider: 'openai',
        model: 'gpt-4',
        messages: [
          {
            role: 'system',
            content: 'You are a helpful assistant created via the TypeScript SDK.'
          }
        ]
      },
      voice: {
        provider: 'elevenlabs',
        voiceId: '21m00Tcm4TlvDq8ikWAM' // Rachel voice
      }
    });
    console.log(` + "`Created assistant: ${assistant.name} (ID: ${assistant.id})`" + `);
    return assistant;
  } catch (error) {
    console.error('Error creating assistant:', error);
    return null;
  }
}

async function makePhoneCall(assistantId: string, phoneNumber: string) {
  try {
    const call = await client.calls.create({
      assistantId,
      customer: {
        number: phoneNumber
      }
    });
    console.log(` + "`Call initiated! Call ID: ${call.id}`" + `);
    console.log(` + "`Status: ${call.status}`" + `);
    return call;
  } catch (error) {
    console.error('Error making call:', error);
    return null;
  }
}

async function main() {
  console.log('🚀 Vapi TypeScript SDK Example');
  console.log('-'.repeat(50));
  
  // List existing assistants
  console.log('\n📋 Listing assistants...');
  const assistants = await listAssistants();
  
  // Create a new assistant
  console.log('\n✨ Creating a new assistant...');
  const newAssistant = await createAssistant();
  
  if (newAssistant) {
    // Example of making a call (commented out for safety)
    // console.log('\n📞 Making a phone call...');
    // await makePhoneCall(newAssistant.id, '+1234567890');
    console.log('\n💡 To make a phone call, uncomment the code above and provide a valid phone number');
  }
  
  console.log('\n✅ Example completed!');
}

// Run the example
main().catch(console.error);
`

	return os.WriteFile(filepath.Join(dir, "basic-example.ts"), []byte(content), 0o644)
}

func generateNodeTSExpressExample(dir string) error {
	content := `/**
 * Express + Vapi TypeScript Integration Example
 * 
 * This example shows how to integrate Vapi with an Express.js application
 * using TypeScript for type safety.
 */

import 'dotenv/config';
import express, { Request, Response } from 'express';
import { Vapi } from '@vapi-ai/server-sdk';

const app = express();
app.use(express.json());

// Initialize Vapi client
const client = new Vapi({
  apiKey: process.env.VAPI_API_KEY!,
});

// Serve the main page
app.get('/', (req: Request, res: Response) => {
  const html = ` + "`" + `
<!DOCTYPE html>
<html>
<head>
    <title>Vapi Express TypeScript Example</title>
    <script src="https://cdn.jsdelivr.net/npm/@vapi-ai/web@latest/dist/vapi.js"></script>
</head>
<body>
    <h1>Vapi Voice Assistant (Express + TypeScript)</h1>
    <button id="startCall">Start Call</button>
    <button id="endCall" disabled>End Call</button>
    <div id="status"></div>
    
    <script>
        const vapi = new Vapi("${process.env.VAPI_PUBLIC_KEY}");
        const assistantId = "${process.env.VAPI_ASSISTANT_ID}";
        
        document.getElementById('startCall').onclick = async () => {
            try {
                await vapi.start(assistantId);
                document.getElementById('startCall').disabled = true;
                document.getElementById('endCall').disabled = false;
                document.getElementById('status').innerText = 'Call started';
            } catch (error) {
                console.error('Error starting call:', error);
            }
        };
        
        document.getElementById('endCall').onclick = async () => {
            try {
                await vapi.stop();
                document.getElementById('startCall').disabled = false;
                document.getElementById('endCall').disabled = true;
                document.getElementById('status').innerText = 'Call ended';
            } catch (error) {
                console.error('Error ending call:', error);
            }
        };
    </script>
</body>
</html>
` + "`" + `;
  res.send(html);
});

// API endpoint to list assistants
app.get('/api/assistants', async (req: Request, res: Response) => {
  try {
    const assistants = await client.assistants.list();
    res.json(assistants.map(a => ({
      id: a.id,
      name: a.name
    })));
  } catch (error) {
    res.status(500).json({ error: (error as Error).message });
  }
});

// API endpoint to create a call
interface CreateCallRequest {
  assistantId: string;
  phoneNumber: string;
}

app.post('/api/calls', async (req: Request<{}, {}, CreateCallRequest>, res: Response) => {
  try {
    const { assistantId, phoneNumber } = req.body;
    const call = await client.calls.create({
      assistantId,
      customer: {
        number: phoneNumber
      }
    });
    res.json({
      id: call.id,
      status: call.status
    });
  } catch (error) {
    res.status(500).json({ error: (error as Error).message });
  }
});

// Webhook endpoint for Vapi events
interface VapiWebhookEvent {
  type: string;
  call?: {
    id: string;
    [key: string]: any;
  };
  transcript?: string;
}

app.post('/webhook/vapi', (req: Request<{}, {}, VapiWebhookEvent>, res: Response) => {
  const { type } = req.body;
  console.log(` + "`Received Vapi webhook: ${type}`" + `);
  
  switch (type) {
    case 'call.started':
      if (req.body.call) {
        console.log(` + "`Call started: ${req.body.call.id}`" + `);
      }
      break;
      
    case 'call.ended':
      if (req.body.call) {
        console.log(` + "`Call ended: ${req.body.call.id}`" + `);
      }
      break;
      
    case 'transcript':
      if (req.body.transcript) {
        console.log(` + "`Transcript: ${req.body.transcript}`" + `);
      }
      break;
  }
  
  res.json({ status: 'ok' });
});

// Health check endpoint
app.get('/health', (req: Request, res: Response) => {
  res.json({ status: 'healthy', service: 'vapi-express-ts-example' });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(` + "`🚀 Server running on http://localhost:${PORT}`" + `);
});
`

	return os.WriteFile(filepath.Join(dir, "express-example.ts"), []byte(content), 0o644)
}

func generateNodeTSFastifyExample(dir string) error {
	content := `/**
 * Fastify + Vapi TypeScript Integration Example
 * 
 * This example shows how to integrate Vapi with Fastify for
 * high-performance Node.js applications.
 */

import 'dotenv/config';
import Fastify from 'fastify';
import { Vapi } from '@vapi-ai/server-sdk';

// Initialize Fastify
const fastify = Fastify({
  logger: true
});

// Initialize Vapi client
const client = new Vapi({
  apiKey: process.env.VAPI_API_KEY!,
});

// HTML response for main page
const htmlTemplate = ` + "`" + `
<!DOCTYPE html>
<html>
<head>
    <title>Vapi Fastify TypeScript Example</title>
    <script src="https://cdn.jsdelivr.net/npm/@vapi-ai/web@latest/dist/vapi.js"></script>
</head>
<body>
    <h1>Vapi Voice Assistant (Fastify + TypeScript)</h1>
    <button id="startCall">Start Call</button>
    <button id="endCall" disabled>End Call</button>
    <div id="status"></div>
    
    <script>
        const vapi = new Vapi("${process.env.VAPI_PUBLIC_KEY}");
        const assistantId = "${process.env.VAPI_ASSISTANT_ID}";
        
        document.getElementById('startCall').onclick = async () => {
            try {
                await vapi.start(assistantId);
                document.getElementById('startCall').disabled = true;
                document.getElementById('endCall').disabled = false;
                document.getElementById('status').innerText = 'Call started';
            } catch (error) {
                console.error('Error starting call:', error);
            }
        };
        
        document.getElementById('endCall').onclick = async () => {
            try {
                await vapi.stop();
                document.getElementById('startCall').disabled = false;
                document.getElementById('endCall').disabled = true;
                document.getElementById('status').innerText = 'Call ended';
            } catch (error) {
                console.error('Error ending call:', error);
            }
        };
    </script>
</body>
</html>
` + "`" + `;

// Routes
fastify.get('/', async (request, reply) => {
  reply.type('text/html').send(htmlTemplate);
});

// List assistants
fastify.get('/api/assistants', async (request, reply) => {
  try {
    const assistants = await client.assistants.list();
    return assistants.map(a => ({
      id: a.id,
      name: a.name
    }));
  } catch (error) {
    reply.status(500).send({ error: (error as Error).message });
  }
});

// Create call
interface CreateCallBody {
  assistantId: string;
  phoneNumber: string;
}

fastify.post<{ Body: CreateCallBody }>('/api/calls', async (request, reply) => {
  try {
    const { assistantId, phoneNumber } = request.body;
    const call = await client.calls.create({
      assistantId,
      customer: {
        number: phoneNumber
      }
    });
    return {
      id: call.id,
      status: call.status
    };
  } catch (error) {
    reply.status(500).send({ error: (error as Error).message });
  }
});

// Webhook handler
interface VapiWebhookBody {
  type: string;
  call?: {
    id: string;
  };
  transcript?: string;
}

fastify.post<{ Body: VapiWebhookBody }>('/webhook/vapi', async (request, reply) => {
  const { type } = request.body;
  fastify.log.info(` + "`Received Vapi webhook: ${type}`" + `);
  
  switch (type) {
    case 'call.started':
      if (request.body.call) {
        fastify.log.info(` + "`Call started: ${request.body.call.id}`" + `);
      }
      break;
      
    case 'call.ended':
      if (request.body.call) {
        fastify.log.info(` + "`Call ended: ${request.body.call.id}`" + `);
      }
      break;
      
    case 'transcript':
      if (request.body.transcript) {
        fastify.log.info(` + "`Transcript: ${request.body.transcript}`" + `);
      }
      break;
  }
  
  return { status: 'ok' };
});

// Health check
fastify.get('/health', async (request, reply) => {
  return { status: 'healthy', service: 'vapi-fastify-ts-example' };
});

// Start server
const start = async () => {
  try {
    const PORT = process.env.PORT || 3000;
    await fastify.listen({ port: Number(PORT), host: '0.0.0.0' });
    console.log(` + "`🚀 Server running on http://localhost:${PORT}`" + `);
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
};

start();
`

	return os.WriteFile(filepath.Join(dir, "fastify-example.ts"), []byte(content), 0o644)
}

func generateNodeEnvTemplate(projectPath string) error {
	content := `# Vapi Configuration
VAPI_API_KEY=your_api_key_here
VAPI_PUBLIC_KEY=your_public_key_here
VAPI_ASSISTANT_ID=your_assistant_id_here

# Server Configuration
PORT=3000

# Optional: Webhook URL for receiving call events
VAPI_WEBHOOK_URL=https://your-domain.com/webhook/vapi
`

	return os.WriteFile(filepath.Join(projectPath, ".env.example"), []byte(content), 0o644)
}
