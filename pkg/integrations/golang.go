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

// GenerateGoIntegration creates Go files for Vapi integration
func GenerateGoIntegration(projectPath string, info *ProjectInfo) error {
	// Create examples directory
	examplesDir := filepath.Join(projectPath, "examples", "vapi")
	if err := os.MkdirAll(examplesDir, 0o750); err != nil {
		return fmt.Errorf("failed to create examples directory: %w", err)
	}

	// Generate basic example
	if err := generateGoBasicExample(examplesDir); err != nil {
		return err
	}

	// Generate HTTP server example
	if err := generateGoHTTPExample(examplesDir); err != nil {
		return err
	}

	// Generate Gin framework example
	if err := generateGoGinExample(examplesDir); err != nil {
		return err
	}

	// Generate environment template
	if err := generateGoEnvTemplate(projectPath); err != nil {
		return err
	}

	return nil
}

func generateGoBasicExample(dir string) error {
	content := `package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/VapiAI/vapi-go"
	"github.com/VapiAI/vapi-go/option"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize Vapi client
	client := vapi.NewClient(
		option.WithAPIKey(os.Getenv("VAPI_API_KEY")),
	)

	fmt.Println("🚀 Vapi Go SDK Example")
	fmt.Println(strings.Repeat("-", 50))

	// List assistants
	fmt.Println("\n📋 Listing assistants...")
	listAssistants(client)

	// Create a new assistant
	fmt.Println("\n✨ Creating a new assistant...")
	assistant := createAssistant(client)

	if assistant != nil {
		// Example of making a call (commented out for safety)
		// fmt.Println("\n📞 Making a phone call...")
		// makePhoneCall(client, assistant.ID, "+1234567890")
		fmt.Println("\n💡 To make a phone call, uncomment the code above and provide a valid phone number")
	}

	fmt.Println("\n✅ Example completed!")
}

func listAssistants(client *vapi.Client) {
	ctx := context.Background()
	
	assistants, err := client.Assistants.List(ctx, &vapi.AssistantListParams{
		Limit: vapi.Int(10),
	})
	if err != nil {
		log.Printf("Error listing assistants: %v", err)
		return
	}

	fmt.Printf("Found %d assistants:\n", len(assistants.Items))
	for _, assistant := range assistants.Items {
		fmt.Printf("  - %s (ID: %s)\n", assistant.Name, assistant.ID)
	}
}

func createAssistant(client *vapi.Client) *vapi.Assistant {
	ctx := context.Background()

	assistant, err := client.Assistants.Create(ctx, &vapi.AssistantCreateParams{
		Name: "Go SDK Assistant",
		Model: &vapi.Model{
			Provider: "openai",
			Model:    "gpt-4",
			Messages: []vapi.Message{
				{
					Role:    "system",
					Content: "You are a helpful assistant created via the Go SDK.",
				},
			},
		},
		Voice: &vapi.Voice{
			Provider: "elevenlabs",
			VoiceID:  "21m00Tcm4TlvDq8ikWAM", // Rachel voice
		},
	})
	if err != nil {
		log.Printf("Error creating assistant: %v", err)
		return nil
	}

	fmt.Printf("Created assistant: %s (ID: %s)\n", assistant.Name, assistant.ID)
	return assistant
}

func makePhoneCall(client *vapi.Client, assistantID, phoneNumber string) {
	ctx := context.Background()

	call, err := client.Calls.Create(ctx, &vapi.CallCreateParams{
		AssistantID: assistantID,
		Customer: &vapi.Customer{
			Number: phoneNumber,
		},
	})
	if err != nil {
		log.Printf("Error making call: %v", err)
		return
	}

	fmt.Printf("Call initiated! Call ID: %s\n", call.ID)
	fmt.Printf("Status: %s\n", call.Status)
}
`

	return os.WriteFile(filepath.Join(dir, "basic_example.go"), []byte(content), 0o600)
}

func generateGoHTTPExample(dir string) error {
	content := `package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/VapiAI/vapi-go"
	"github.com/VapiAI/vapi-go/option"
	"github.com/joho/godotenv"
)

var client *vapi.Client

// HTML template for the web interface
const htmlTemplate = ` + "`" + `
<!DOCTYPE html>
<html>
<head>
    <title>Vapi Go Example</title>
    <script src="https://cdn.jsdelivr.net/npm/@vapi-ai/web@latest/dist/vapi.js"></script>
</head>
<body>
    <h1>Vapi Voice Assistant (Go)</h1>
    <button id="startCall">Start Call</button>
    <button id="endCall" disabled>End Call</button>
    <div id="status"></div>
    
    <script>
        const vapi = new Vapi("{{.PublicKey}}");
        const assistantId = "{{.AssistantID}}";
        
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
` + "`" + `

type PageData struct {
	PublicKey   string
	AssistantID string
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize Vapi client
	client = vapi.NewClient(
		option.WithAPIKey(os.Getenv("VAPI_API_KEY")),
	)

	// Set up routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/assistants", handleListAssistants)
	http.HandleFunc("/api/calls", handleCreateCall)
	http.HandleFunc("/webhook/vapi", handleVapiWebhook)

	fmt.Println("🚀 Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("home").Parse(htmlTemplate))
	
	data := PageData{
		PublicKey:   os.Getenv("VAPI_PUBLIC_KEY"),
		AssistantID: os.Getenv("VAPI_ASSISTANT_ID"),
	}
	
	tmpl.Execute(w, data)
}

func handleListAssistants(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	assistants, err := client.Assistants.List(ctx, &vapi.AssistantListParams{
		Limit: vapi.Int(20),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assistants.Items)
}

func handleCreateCall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AssistantID string ` + "`json:\"assistantId\"`" + `
		PhoneNumber string ` + "`json:\"phoneNumber\"`" + `
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	call, err := client.Calls.Create(ctx, &vapi.CallCreateParams{
		AssistantID: req.AssistantID,
		Customer: &vapi.Customer{
			Number: req.PhoneNumber,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(call)
}

func handleVapiWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	eventType, _ := event["type"].(string)
	log.Printf("Received Vapi webhook: %s", eventType)

	switch eventType {
	case "call.started":
		if call, ok := event["call"].(map[string]interface{}); ok {
			callID, _ := call["id"].(string)
			log.Printf("Call started: %s", callID)
		}
	case "call.ended":
		if call, ok := event["call"].(map[string]interface{}); ok {
			callID, _ := call["id"].(string)
			log.Printf("Call ended: %s", callID)
		}
	case "transcript":
		if transcript, ok := event["transcript"].(string); ok {
			log.Printf("Transcript: %s", transcript)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
`

	return os.WriteFile(filepath.Join(dir, "http_example.go"), []byte(content), 0o600)
}

func generateGoGinExample(dir string) error {
	content := `package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/VapiAI/vapi-go"
	"github.com/VapiAI/vapi-go/option"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var client *vapi.Client

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize Vapi client
	client = vapi.NewClient(
		option.WithAPIKey(os.Getenv("VAPI_API_KEY")),
	)

	// Set up Gin router
	r := gin.Default()

	// Routes
	r.GET("/", handleHome)
	r.GET("/api/assistants", handleListAssistants)
	r.POST("/api/calls", handleCreateCall)
	r.POST("/webhook/vapi", handleVapiWebhook)
	r.GET("/health", handleHealth)

	log.Println("🚀 Server starting on http://localhost:8080")
	r.Run(":8080")
}

func handleHome(c *gin.Context) {
	html := ` + "`" + `
<!DOCTYPE html>
<html>
<head>
    <title>Vapi Gin Example</title>
    <script src="https://cdn.jsdelivr.net/npm/@vapi-ai/web@latest/dist/vapi.js"></script>
</head>
<body>
    <h1>Vapi Voice Assistant (Go + Gin)</h1>
    <button id="startCall">Start Call</button>
    <button id="endCall" disabled>End Call</button>
    <div id="status"></div>
    
    <script>
        const vapi = new Vapi("` + "`" + ` + os.Getenv("VAPI_PUBLIC_KEY") + ` + "`" + `");
        const assistantId = "` + "`" + ` + os.Getenv("VAPI_ASSISTANT_ID") + ` + "`" + `";
        
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
` + "`" + `
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func handleListAssistants(c *gin.Context) {
	ctx := context.Background()
	
	assistants, err := client.Assistants.List(ctx, &vapi.AssistantListParams{
		Limit: vapi.Int(20),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assistants.Items)
}

func handleCreateCall(c *gin.Context) {
	var req struct {
		AssistantID string ` + "`json:\"assistantId\"`" + `
		PhoneNumber string ` + "`json:\"phoneNumber\"`" + `
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	call, err := client.Calls.Create(ctx, &vapi.CallCreateParams{
		AssistantID: req.AssistantID,
		Customer: &vapi.Customer{
			Number: req.PhoneNumber,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, call)
}

func handleVapiWebhook(c *gin.Context) {
	var event map[string]interface{}
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventType, _ := event["type"].(string)
	log.Printf("Received Vapi webhook: %s", eventType)

	switch eventType {
	case "call.started":
		if call, ok := event["call"].(map[string]interface{}); ok {
			callID, _ := call["id"].(string)
			log.Printf("Call started: %s", callID)
		}
	case "call.ended":
		if call, ok := event["call"].(map[string]interface{}); ok {
			callID, _ := call["id"].(string)
			log.Printf("Call ended: %s", callID)
		}
	case "transcript":
		if transcript, ok := event["transcript"].(string); ok {
			log.Printf("Transcript: %s", transcript)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "vapi-gin-example",
	})
}
`

	return os.WriteFile(filepath.Join(dir, "gin_example.go"), []byte(content), 0o600)
}

func generateGoEnvTemplate(projectPath string) error {
	content := `# Vapi Configuration
VAPI_API_KEY=your_api_key_here
VAPI_PUBLIC_KEY=your_public_key_here
VAPI_ASSISTANT_ID=your_assistant_id_here

# Server Configuration
PORT=8080
`

	return os.WriteFile(filepath.Join(projectPath, ".env.example"), []byte(content), 0o600)
}
