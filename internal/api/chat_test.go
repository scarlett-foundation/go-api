package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go-api/internal/api"
	"go-api/internal/types"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Chat API", func() {
	var (
		e       *echo.Echo
		handler *api.Handler
	)

	BeforeEach(func() {
		// Load environment variables
		err := godotenv.Load("../../.env")
		Expect(err).NotTo(HaveOccurred())

		apiKey := os.Getenv("GROQ_API_KEY")
		Expect(apiKey).NotTo(BeEmpty(), "GROQ_API_KEY must be set")

		// Setup Echo and handler
		e = echo.New()
		handler = api.NewHandler(apiKey)
	})

	Context("POST /chat/completions", func() {
		It("should successfully process a chat completion request", func() {
			// Prepare request body
			reqBody := types.ChatRequest{
				Messages: []types.Message{
					{
						Role:    "system",
						Content: "You are a helpful AI assistant who provides concise responses.",
					},
					{
						Role:    "user",
						Content: "What is the capital of France?",
					},
				},
				Model:       "deepseek-r1-distill-llama-70b",
				Temperature: 0.7,
				MaxTokens:   150,
				Stream:      false,
			}

			jsonBody, err := json.Marshal(reqBody)
			Expect(err).NotTo(HaveOccurred())

			// Log request
			GinkgoWriter.Printf("\n🌟 Test Request:\n%s\n", prettyJSON(jsonBody))

			// Create test request
			req := httptest.NewRequest(http.MethodPost, "/chat/completions", bytes.NewBuffer(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Make request
			err = handler.HandleChatCompletions(c)
			Expect(err).NotTo(HaveOccurred())

			// Check response
			Expect(rec.Code).To(Equal(http.StatusOK))

			// Log response status
			GinkgoWriter.Printf("\n📊 Response Status: %d\n", rec.Code)

			// Parse and log response
			var response types.ChatResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())

			// Log formatted response
			GinkgoWriter.Printf("\n✅ Success Response:\n%s\n", prettyJSON(rec.Body.Bytes()))

			// Log specific fields of interest
			GinkgoWriter.Printf("\n📝 Response Details:\n")
			GinkgoWriter.Printf("- ID: %s\n", response.ID)
			GinkgoWriter.Printf("- Model: %s\n", response.Model)
			GinkgoWriter.Printf("- Total Tokens: %d\n", response.Usage.TotalTokens)
			GinkgoWriter.Printf("- Assistant's Response: %s\n", response.Choices[0].Message.Content)

			// Validate response structure
			Expect(response.ID).NotTo(BeEmpty())
			Expect(response.Model).To(Equal("deepseek-r1-distill-llama-70b"))
			Expect(response.Choices).To(HaveLen(1))
			Expect(response.Choices[0].Message.Role).To(Equal("assistant"))
			Expect(response.Choices[0].Message.Content).NotTo(BeEmpty())
			Expect(response.Usage.TotalTokens).To(BeNumerically(">", 0))
		})

		It("should handle invalid requests", func() {
			invalidJSON := "invalid json"
			// Log invalid request
			GinkgoWriter.Printf("\n❌ Testing Invalid Request:\n%s\n", invalidJSON)

			// Send invalid JSON
			req := httptest.NewRequest(http.MethodPost, "/chat/completions", bytes.NewBufferString(invalidJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Make request
			err := handler.HandleChatCompletions(c)
			Expect(err).NotTo(HaveOccurred())

			// Log response status
			GinkgoWriter.Printf("\n📊 Error Response Status: %d\n", rec.Code)

			// Check response
			Expect(rec.Code).To(Equal(http.StatusBadRequest))

			// Parse and log error response
			var errorResponse types.ErrorResponse
			err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
			Expect(err).NotTo(HaveOccurred())

			// Log formatted error response
			GinkgoWriter.Printf("\n❌ Error Response:\n%s\n", prettyJSON(rec.Body.Bytes()))

			Expect(errorResponse.Error.Type).To(Equal("invalid_request_error"))
		})
	})
})

// prettyJSON formats JSON with indentation for better readability in logs
func prettyJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, data, "", "  ")
	if err != nil {
		return string(data) // Return original if formatting fails
	}
	return prettyJSON.String()
}

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}
