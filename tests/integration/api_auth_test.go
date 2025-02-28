package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Request and response types
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// Test setup
const (
	apiURL        = "http://localhost:8082/chat/completions"
	validAPIKey   = "test-key-2"
	invalidAPIKey = "invalid-key"
)

var apiCmd *exec.Cmd
var binaryPath string
var cleanupBinary bool

func TestAPIAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Authentication Test Suite")
}

var _ = BeforeSuite(func() {
	// Find project root
	projectRoot, err := filepath.Abs("../..")
	Expect(err).NotTo(HaveOccurred(), "Failed to get project root")

	// Set up binary path
	binaryPath = filepath.Join(projectRoot, "api-server")

	// Check if server binary already exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		// Build the API server if it doesn't exist
		fmt.Println("Building API server...")
		buildCmd := exec.Command("go", "build", "-o", binaryPath)
		buildCmd.Dir = projectRoot
		output, err := buildCmd.CombinedOutput()
		Expect(err).NotTo(HaveOccurred(), "Failed to build API server: %s", string(output))

		// Mark binary for cleanup if we built it during this test run
		cleanupBinary = true
	}

	// Start the API server
	fmt.Println("Starting API server...")
	apiCmd = exec.Command(binaryPath)
	apiCmd.Dir = projectRoot
	err = apiCmd.Start()
	Expect(err).NotTo(HaveOccurred(), "Failed to start API server")

	// Wait for API server to start
	fmt.Println("Waiting for API server to start...")
	waitForServer()
	fmt.Println("API server started successfully")
})

var _ = AfterSuite(func() {
	// Kill the API server
	if apiCmd != nil && apiCmd.Process != nil {
		fmt.Println("Stopping API server...")
		apiCmd.Process.Kill()

		// Wait for the process to fully exit
		apiCmd.Wait()
	}

	// Clean up the binary if we built it as part of this test run
	if cleanupBinary && binaryPath != "" {
		fmt.Println("Cleaning up API server binary...")
		err := os.Remove(binaryPath)
		if err != nil {
			fmt.Printf("Warning: Failed to clean up binary at %s: %v\n", binaryPath, err)
		}
	}
})

func waitForServer() {
	// Try to connect to the server for up to 5 seconds
	timeout := time.After(5 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	for {
		select {
		case <-timeout:
			Fail("Timed out waiting for API server to start")
		case <-tick:
			_, err := http.Get("http://localhost:8082/")
			if err == nil || !isConnectionRefused(err) {
				// Allow a little more time for the server to initialize fully
				time.Sleep(500 * time.Millisecond)
				return
			}
		}
	}
}

func isConnectionRefused(err error) bool {
	return err != nil && err.Error() != "" &&
		(err.Error() == "connect: connection refused" ||
			err.Error() == "Get \"http://localhost:8082/\": connect: connection refused")
}

var _ = Describe("API Authentication", func() {
	var chatReq ChatRequest

	BeforeEach(func() {
		// Set up a standard request for all tests
		chatReq = ChatRequest{
			Model: "deepseek-r1-distill-llama-70b",
			Messages: []Message{
				{
					Role:    "user",
					Content: "Hello, how are you?",
				},
			},
			Temperature: 0.7,
			MaxTokens:   50,
		}
	})

	Context("when using a valid API key", func() {
		It("should return a successful response", func() {
			// Convert request to JSON
			reqBody, err := json.Marshal(chatReq)
			Expect(err).NotTo(HaveOccurred())

			// Create request
			req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
			Expect(err).NotTo(HaveOccurred())

			// Set headers
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validAPIKey))

			// Send request
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			// Check response
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			// Read response body
			body, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			// Check for some key fields in response
			Expect(string(body)).To(ContainSubstring("id"))
			Expect(string(body)).To(ContainSubstring("choices"))
		})
	})

	Context("when using an invalid API key", func() {
		It("should return an unauthorized error", func() {
			// Convert request to JSON
			reqBody, err := json.Marshal(chatReq)
			Expect(err).NotTo(HaveOccurred())

			// Create request
			req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
			Expect(err).NotTo(HaveOccurred())

			// Set headers
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", invalidAPIKey))

			// Send request
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			// Check response
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))

			// Read response body
			body, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			// Parse error response
			var errorResp ErrorResponse
			err = json.Unmarshal(body, &errorResp)
			Expect(err).NotTo(HaveOccurred())

			// Check error message
			Expect(errorResp.Error.Message).To(Equal("Invalid API key"))
			Expect(errorResp.Error.Type).To(Equal("unauthorized"))
		})
	})

	Context("when no API key is provided", func() {
		It("should return an unauthorized error", func() {
			// Convert request to JSON
			reqBody, err := json.Marshal(chatReq)
			Expect(err).NotTo(HaveOccurred())

			// Create request
			req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
			Expect(err).NotTo(HaveOccurred())

			// Set headers (no Authorization header)
			req.Header.Set("Content-Type", "application/json")

			// Send request
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			// Check response
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))

			// Read response body
			body, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			// Parse error response
			var errorResp ErrorResponse
			err = json.Unmarshal(body, &errorResp)
			Expect(err).NotTo(HaveOccurred())

			// Check error message
			Expect(errorResp.Error.Message).To(Equal("Missing API key"))
			Expect(errorResp.Error.Type).To(Equal("unauthorized"))
		})
	})
})
