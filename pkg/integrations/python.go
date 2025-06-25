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

// GeneratePythonIntegration generates SDK examples and configuration for Python projects
func GeneratePythonIntegration(projectPath string, info *ProjectInfo) error {
	examplesDir := filepath.Join(projectPath, "vapi_examples")
	if err := os.MkdirAll(examplesDir, 0o750); err != nil {
		return fmt.Errorf("failed to create examples directory: %w", err)
	}

	// Create all example files
	if err := generatePythonBasicExample(examplesDir); err != nil {
		return err
	}

	if err := generatePythonFlaskExample(examplesDir); err != nil {
		return err
	}

	if err := generatePythonFastAPIExample(examplesDir); err != nil {
		return err
	}

	if err := generatePythonEnvTemplate(projectPath); err != nil {
		return err
	}

	if err := generatePythonRequirements(projectPath); err != nil {
		return err
	}

	return nil
}

func generatePythonBasicExample(dir string) error {
	content := `"""
Basic Vapi Python SDK Example

This example demonstrates how to use the Vapi Python SDK to manage assistants
and make phone calls programmatically.
"""

import os
from vapi import Vapi
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Initialize Vapi client
client = Vapi(api_key=os.getenv("VAPI_API_KEY"))

def list_assistants():
    """List all assistants"""
    try:
        assistants = client.assistants.list()
        print(f"Found {len(assistants)} assistants:")
        for assistant in assistants:
            print(f"  - {assistant.name} (ID: {assistant.id})")
        return assistants
    except Exception as e:
        print(f"Error listing assistants: {e}")
        return []

def create_assistant():
    """Create a new assistant"""
    try:
        assistant = client.assistants.create(
            name="Python SDK Assistant",
            model={
                "provider": "openai",
                "model": "gpt-4",
                "messages": [
                    {
                        "role": "system",
                        "content": "You are a helpful assistant created via the Python SDK."
                    }
                ]
            },
            voice={
                "provider": "elevenlabs",
                "voiceId": "21m00Tcm4TlvDq8ikWAM"  # Rachel voice
            }
        )
        print(f"Created assistant: {assistant.name} (ID: {assistant.id})")
        return assistant
    except Exception as e:
        print(f"Error creating assistant: {e}")
        return None

def make_phone_call(assistant_id, phone_number):
    """Make an outbound phone call"""
    try:
        call = client.calls.create(
            assistantId=assistant_id,
            customer={
                "number": phone_number
            }
        )
        print(f"Call initiated! Call ID: {call.id}")
        print(f"Status: {call.status}")
        return call
    except Exception as e:
        print(f"Error making call: {e}")
        return None

def main():
    print("🚀 Vapi Python SDK Example")
    print("-" * 50)
    
    # List existing assistants
    print("\n📋 Listing assistants...")
    assistants = list_assistants()
    
    # Create a new assistant
    print("\n✨ Creating a new assistant...")
    new_assistant = create_assistant()
    
    if new_assistant:
        # Example of making a call (commented out for safety)
        # print("\n📞 Making a phone call...")
        # call = make_phone_call(new_assistant.id, "+1234567890")
        print("\n💡 To make a phone call, uncomment the code above and provide a valid phone number")
    
    print("\n✅ Example completed!")

if __name__ == "__main__":
    main()
`

	return os.WriteFile(filepath.Join(dir, "basic_example.py"), []byte(content), 0o600)
}

func generatePythonFlaskExample(dir string) error {
	content := `"""
Flask + Vapi Integration Example

This example shows how to integrate Vapi with a Flask web application
to handle webhooks and create web-based voice interfaces.
"""

import os
from flask import Flask, request, jsonify, render_template_string
from vapi import Vapi
from dotenv import load_dotenv

load_dotenv()

app = Flask(__name__)
client = Vapi(api_key=os.getenv("VAPI_API_KEY"))

# Simple HTML template for the web interface
HTML_TEMPLATE = """
<!DOCTYPE html>
<html>
<head>
    <title>Vapi Flask Example</title>
    <script src="https://cdn.jsdelivr.net/npm/@vapi-ai/web@latest/dist/vapi.js"></script>
</head>
<body>
    <h1>Vapi Voice Assistant</h1>
    <button id="startCall">Start Call</button>
    <button id="endCall" disabled>End Call</button>
    <div id="status"></div>
    
    <script>
        const vapi = new Vapi("{{ vapi_public_key }}");
        const assistantId = "{{ assistant_id }}";
        
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
"""

@app.route('/')
def index():
    """Render the main page with Vapi web interface"""
    return render_template_string(
        HTML_TEMPLATE,
        vapi_public_key=os.getenv("VAPI_PUBLIC_KEY"),
        assistant_id=os.getenv("VAPI_ASSISTANT_ID")
    )

@app.route('/webhook/vapi', methods=['POST'])
def vapi_webhook():
    """Handle Vapi webhooks for call events"""
    data = request.json
    event_type = data.get('type')
    
    print(f"Received Vapi webhook: {event_type}")
    
    if event_type == 'call.started':
        call_id = data.get('call', {}).get('id')
        print(f"Call started: {call_id}")
        # Handle call start logic
        
    elif event_type == 'call.ended':
        call_id = data.get('call', {}).get('id')
        print(f"Call ended: {call_id}")
        # Handle call end logic
        
    elif event_type == 'transcript':
        transcript = data.get('transcript')
        print(f"Transcript: {transcript}")
        # Handle transcript
    
    return jsonify({"status": "ok"})

@app.route('/api/assistants', methods=['GET'])
def list_assistants():
    """API endpoint to list all assistants"""
    try:
        assistants = client.assistants.list()
        return jsonify([{
            "id": a.id,
            "name": a.name
        } for a in assistants])
    except Exception as e:
        return jsonify({"error": str(e)}), 500

@app.route('/api/calls', methods=['POST'])
def create_call():
    """API endpoint to create a new call"""
    try:
        data = request.json
        call = client.calls.create(
            assistantId=data.get('assistantId'),
            customer={
                "number": data.get('phoneNumber')
            }
        )
        return jsonify({
            "id": call.id,
            "status": call.status
        })
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(debug=True, port=5000)
`

	return os.WriteFile(filepath.Join(dir, "flask_example.py"), []byte(content), 0o600)
}

func generatePythonFastAPIExample(dir string) error {
	content := `"""
FastAPI + Vapi Integration Example

This example shows how to integrate Vapi with FastAPI for high-performance
async web applications with automatic API documentation.
"""

import os
from fastapi import FastAPI, HTTPException, Request
from fastapi.responses import HTMLResponse
from pydantic import BaseModel
from vapi import Vapi
from dotenv import load_dotenv
import uvicorn

load_dotenv()

app = FastAPI(title="Vapi FastAPI Example")
client = Vapi(api_key=os.getenv("VAPI_API_KEY"))

# Request/Response models
class CallRequest(BaseModel):
    assistant_id: str
    phone_number: str

class WebhookEvent(BaseModel):
    type: str
    call: dict = None
    transcript: str = None

# HTML template for the web interface
HTML_TEMPLATE = """
<!DOCTYPE html>
<html>
<head>
    <title>Vapi FastAPI Example</title>
    <script src="https://cdn.jsdelivr.net/npm/@vapi-ai/web@latest/dist/vapi.js"></script>
</head>
<body>
    <h1>Vapi Voice Assistant (FastAPI)</h1>
    <button id="startCall">Start Call</button>
    <button id="endCall" disabled>End Call</button>
    <div id="status"></div>
    
    <script>
        const vapi = new Vapi("{vapi_public_key}");
        const assistantId = "{assistant_id}";
        
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
"""

@app.get("/", response_class=HTMLResponse)
async def root():
    """Serve the main page with Vapi web interface"""
    html = HTML_TEMPLATE.format(
        vapi_public_key=os.getenv("VAPI_PUBLIC_KEY", ""),
        assistant_id=os.getenv("VAPI_ASSISTANT_ID", "")
    )
    return HTMLResponse(content=html)

@app.get("/api/assistants")
async def list_assistants():
    """List all assistants"""
    try:
        assistants = client.assistants.list()
        return [{
            "id": a.id,
            "name": a.name,
            "created_at": a.created_at
        } for a in assistants]
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/calls")
async def create_call(call_request: CallRequest):
    """Create a new outbound call"""
    try:
        call = client.calls.create(
            assistantId=call_request.assistant_id,
            customer={
                "number": call_request.phone_number
            }
        )
        return {
            "id": call.id,
            "status": call.status,
            "created_at": call.created_at
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/webhook/vapi")
async def vapi_webhook(event: WebhookEvent):
    """Handle Vapi webhooks"""
    print(f"Received Vapi webhook: {event.type}")
    
    if event.type == "call.started":
        # Handle call started
        call_id = event.call.get("id") if event.call else None
        print(f"Call started: {call_id}")
        
    elif event.type == "call.ended":
        # Handle call ended
        call_id = event.call.get("id") if event.call else None
        print(f"Call ended: {call_id}")
        
    elif event.type == "transcript":
        # Handle transcript
        print(f"Transcript: {event.transcript}")
    
    return {"status": "ok"}

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "service": "vapi-fastapi-example"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
`

	return os.WriteFile(filepath.Join(dir, "fastapi_example.py"), []byte(content), 0o600)
}

func generatePythonEnvTemplate(projectPath string) error {
	content := `# Vapi Configuration
VAPI_API_KEY=your_api_key_here
VAPI_PUBLIC_KEY=your_public_key_here
VAPI_ASSISTANT_ID=your_assistant_id_here

# Optional: Webhook URL for receiving call events
VAPI_WEBHOOK_URL=https://your-domain.com/webhook/vapi
`

	return os.WriteFile(filepath.Join(projectPath, ".env.example"), []byte(content), 0o600)
}

func generatePythonRequirements(projectPath string) error {
	content := `# Vapi Python SDK
vapi-python>=1.0.0

# Web frameworks (optional - uncomment as needed)
# flask>=2.0.0
# fastapi>=0.100.0
# uvicorn>=0.23.0

# Utilities
python-dotenv>=1.0.0
`

	// Check if requirements.txt already exists
	reqPath := filepath.Join(projectPath, "requirements-vapi.txt")
	return os.WriteFile(reqPath, []byte(content), 0o600)
}
