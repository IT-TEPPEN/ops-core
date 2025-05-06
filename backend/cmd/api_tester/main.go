package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Configuration for the API tester
type Config struct {
	BaseURL string
	Verbose bool
}

// Repository represents a repository from the API
type Repository struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FileNode represents a file or directory within a repository
type FileNode struct {
	Path string `json:"path"`
	Type string `json:"type"` // "file" or "dir"
}

// Error response from the API
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Test client for making API requests
type APIClient struct {
	config     Config
	httpClient *http.Client
}

// NewAPIClient creates a new API client with the given configuration
func NewAPIClient(config Config) *APIClient {
	return &APIClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// logVerbose logs a message if verbose mode is enabled
func (c *APIClient) logVerbose(format string, args ...interface{}) {
	if c.config.Verbose {
		fmt.Printf(format+"\n", args...)
	}
}

// makeRequest performs an HTTP request and returns the response
func (c *APIClient) makeRequest(method, path string, body interface{}) (*http.Response, []byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
		c.logVerbose("Request Body: %s", string(jsonData))
	}

	url := c.config.BaseURL + path
	c.logVerbose("%s %s", method, url)

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return resp, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	c.logVerbose("Response Status: %s", resp.Status)
	c.logVerbose("Response Body: %s", string(respBody))

	return resp, respBody, nil
}

// TestRegisterRepository tests repository registration
func (c *APIClient) TestRegisterRepository(url, accessToken string) (*Repository, error) {
	fmt.Println("\n=== Testing Repository Registration ===")

	// Prepare request body
	reqBody := map[string]string{
		"url": url,
	}
	if accessToken != "" {
		reqBody["accessToken"] = accessToken
	}

	// Make request
	resp, body, err := c.makeRequest("POST", "/api/v1/repositories", reqBody)
	if err != nil {
		return nil, err
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("failed to register repository: %s - %s", errResp.Code, errResp.Message)
		}
		return nil, fmt.Errorf("failed to register repository: HTTP %d", resp.StatusCode)
	}

	// Parse response
	var repo Repository
	if err := json.Unmarshal(body, &repo); err != nil {
		return nil, fmt.Errorf("failed to parse repository response: %w", err)
	}

	fmt.Printf("Successfully registered repository: ID=%s, Name=%s\n", repo.ID, repo.Name)
	return &repo, nil
}

// TestListRepositories tests listing all repositories
func (c *APIClient) TestListRepositories() ([]Repository, error) {
	fmt.Println("\n=== Testing List Repositories ===")

	// Make request
	resp, body, err := c.makeRequest("GET", "/api/v1/repositories", nil)
	if err != nil {
		return nil, err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("failed to list repositories: %s - %s", errResp.Code, errResp.Message)
		}
		return nil, fmt.Errorf("failed to list repositories: HTTP %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Repositories []Repository `json:"repositories"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse repositories response: %w", err)
	}

	fmt.Printf("Successfully listed %d repositories\n", len(response.Repositories))
	for i, repo := range response.Repositories {
		fmt.Printf("  %d. ID=%s, Name=%s, URL=%s\n", i+1, repo.ID, repo.Name, repo.URL)
	}

	return response.Repositories, nil
}

// TestGetRepository tests fetching a specific repository
func (c *APIClient) TestGetRepository(repoID string) (*Repository, error) {
	fmt.Printf("\n=== Testing Get Repository (ID: %s) ===\n", repoID)

	// Make request
	resp, body, err := c.makeRequest("GET", fmt.Sprintf("/api/v1/repositories/%s", repoID), nil)
	if err != nil {
		return nil, err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("failed to get repository: %s - %s", errResp.Code, errResp.Message)
		}
		return nil, fmt.Errorf("failed to get repository: HTTP %d", resp.StatusCode)
	}

	// Parse response
	var repo Repository
	if err := json.Unmarshal(body, &repo); err != nil {
		return nil, fmt.Errorf("failed to parse repository response: %w", err)
	}

	fmt.Printf("Successfully retrieved repository: ID=%s, Name=%s, URL=%s\n", repo.ID, repo.Name, repo.URL)
	return &repo, nil
}

// TestUpdateAccessToken tests updating a repository's access token
func (c *APIClient) TestUpdateAccessToken(repoID, accessToken string) error {
	fmt.Printf("\n=== Testing Update Access Token (ID: %s) ===\n", repoID)

	// Prepare request body
	reqBody := map[string]string{
		"accessToken": accessToken,
	}

	// Make request
	resp, body, err := c.makeRequest("PUT", fmt.Sprintf("/api/v1/repositories/%s/token", repoID), reqBody)
	if err != nil {
		return err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return fmt.Errorf("failed to update access token: %s - %s", errResp.Code, errResp.Message)
		}
		return fmt.Errorf("failed to update access token: HTTP %d", resp.StatusCode)
	}

	fmt.Printf("Successfully updated access token for repository ID=%s\n", repoID)
	return nil
}

// TestListRepositoryFiles tests listing files in a repository
func (c *APIClient) TestListRepositoryFiles(repoID string) ([]FileNode, error) {
	fmt.Printf("\n=== Testing List Repository Files (ID: %s) ===\n", repoID)

	// Make request
	resp, body, err := c.makeRequest("GET", fmt.Sprintf("/api/v1/repositories/%s/files", repoID), nil)
	if err != nil {
		return nil, err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("failed to list repository files: %s - %s", errResp.Code, errResp.Message)
		}
		return nil, fmt.Errorf("failed to list repository files: HTTP %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Files []FileNode `json:"files"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse files response: %w", err)
	}

	fmt.Printf("Successfully listed %d files\n", len(response.Files))
	for i, file := range response.Files {
		if i < 10 { // Only show first 10 to avoid flooding the output
			fmt.Printf("  %d. %s (%s)\n", i+1, file.Path, file.Type)
		}
	}
	if len(response.Files) > 10 {
		fmt.Printf("  ... and %d more files\n", len(response.Files)-10)
	}

	return response.Files, nil
}

// TestSelectRepositoryFiles tests selecting files in a repository
func (c *APIClient) TestSelectRepositoryFiles(repoID string, filePaths []string) error {
	fmt.Printf("\n=== Testing Select Repository Files (ID: %s) ===\n", repoID)

	// Prepare request body
	reqBody := map[string][]string{
		"filePaths": filePaths,
	}

	// Make request
	resp, body, err := c.makeRequest("POST", fmt.Sprintf("/api/v1/repositories/%s/files/select", repoID), reqBody)
	if err != nil {
		return err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return fmt.Errorf("failed to select files: %s - %s", errResp.Code, errResp.Message)
		}
		return fmt.Errorf("failed to select files: HTTP %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Message       string `json:"message"`
		RepoID        string `json:"repoId"`
		SelectedFiles int    `json:"selectedFiles"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse selection response: %w", err)
	}

	fmt.Printf("Successfully selected %d files: %s\n", response.SelectedFiles, response.Message)
	return nil
}

// TestGetSelectedMarkdown tests getting selected markdown content
func (c *APIClient) TestGetSelectedMarkdown(repoID string) (string, error) {
	fmt.Printf("\n=== Testing Get Selected Markdown (ID: %s) ===\n", repoID)

	// Make request
	resp, body, err := c.makeRequest("GET", fmt.Sprintf("/api/v1/repositories/%s/markdown", repoID), nil)
	if err != nil {
		return "", err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return "", fmt.Errorf("failed to get markdown content: %s - %s", errResp.Code, errResp.Message)
		}
		return "", fmt.Errorf("failed to get markdown content: HTTP %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		RepoID  string `json:"repoId"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse markdown response: %w", err)
	}

	contentPreview := response.Content
	if len(contentPreview) > 100 {
		contentPreview = contentPreview[:100] + "..."
	}
	
	fmt.Printf("Successfully retrieved markdown content (%d bytes)\n", len(response.Content))
	fmt.Printf("Content preview: %s\n", contentPreview)
	
	return response.Content, nil
}

// RunFullTest runs a complete test of all API endpoints
func (c *APIClient) RunFullTest(repoURL, accessToken string) error {
	fmt.Println("\n========== STARTING FULL API TEST ==========")

	// Test repository registration
	repo, err := c.TestRegisterRepository(repoURL, accessToken)
	if err != nil {
		// Check if it failed because the repository might already exist
		fmt.Printf("Registration failed: %v\n", err)
		fmt.Println("Trying to list repositories to find an existing one...")
		
		// Try to list repositories instead
		repos, err := c.TestListRepositories()
		if err != nil {
			return fmt.Errorf("failed to register repository and failed to list repositories: %w", err)
		}
		
		if len(repos) == 0 {
			return fmt.Errorf("no repositories found after registration failed")
		}
		
		// Use the first repository for further tests
		repo = &repos[0]
		fmt.Printf("Using existing repository ID=%s for further tests\n", repo.ID)
	}

	// Test getting repository details
	_, err = c.TestGetRepository(repo.ID)
	if err != nil {
		return fmt.Errorf("failed to get repository details: %w", err)
	}

	// Test updating access token (if provided)
	if accessToken != "" {
		if err := c.TestUpdateAccessToken(repo.ID, accessToken); err != nil {
			return fmt.Errorf("failed to update access token: %w", err)
		}
	}

	// Test listing repository files
	files, err := c.TestListRepositoryFiles(repo.ID)
	if err != nil {
		return fmt.Errorf("failed to list repository files: %w", err)
	}

	// Select a subset of the files (up to 3 markdown files)
	var markdownFiles []string
	for _, file := range files {
		if file.Type == "file" && len(markdownFiles) < 3 {
			// Simple check for markdown files by extension
			if len(file.Path) > 3 && file.Path[len(file.Path)-3:] == ".md" {
				markdownFiles = append(markdownFiles, file.Path)
			}
		}
	}

	if len(markdownFiles) == 0 {
		fmt.Println("No markdown files found in the repository, skipping file selection and content tests")
	} else {
		// Test selecting files
		if err := c.TestSelectRepositoryFiles(repo.ID, markdownFiles); err != nil {
			return fmt.Errorf("failed to select repository files: %w", err)
		}

		// Test getting markdown content
		_, err = c.TestGetSelectedMarkdown(repo.ID)
		if err != nil {
			return fmt.Errorf("failed to get markdown content: %w", err)
		}
	}

	fmt.Println("\n========== API TEST COMPLETED SUCCESSFULLY ==========")
	return nil
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  api_tester full [flags]    - Run the full API test suite")
	fmt.Println("  api_tester endpoint [cmd]  - Test a specific API endpoint")
	fmt.Println("\nAvailable endpoints:")
	fmt.Println("  register      - Register a new repository")
	fmt.Println("  list          - List all repositories")
	fmt.Println("  get           - Get details of a specific repository")
	fmt.Println("  update-token  - Update repository access token")
	fmt.Println("  files         - List files in a repository")
	fmt.Println("  select-files  - Select markdown files in a repository")
	fmt.Println("  markdown      - Get selected markdown content")
	fmt.Println("\nRun 'api_tester [command] -h' for details on available flags")
	os.Exit(1)
}

func main() {
	// Check for command-line arguments
	if len(os.Args) < 2 {
		printUsage()
	}

	// Determine which mode to run in
	mode := os.Args[1]

	switch mode {
	case "full":
		// Run the full test suite
		runFullTest()
	case "endpoint":
		// Check if an endpoint is specified
		if len(os.Args) < 3 {
			fmt.Println("Error: No endpoint specified")
			printUsage()
		}
		// Run the specific endpoint test
		endpoint := os.Args[2]
		runEndpointTest(endpoint)
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", mode)
		printUsage()
	}
}

func runFullTest() {
	// Parse command-line flags for full test mode
	cmdFull := flag.NewFlagSet("full", flag.ExitOnError)
	baseURL := cmdFull.String("baseURL", "http://localhost:8080", "Base URL of the API")
	repoURL := cmdFull.String("repoURL", "https://github.com/docker/compose", "URL of the repository to test with")
	accessToken := cmdFull.String("token", "", "Access token for private repositories")
	verbose := cmdFull.Bool("verbose", false, "Enable verbose logging")
	cmdFull.Parse(os.Args[2:])

	// Create config
	config := Config{
		BaseURL: *baseURL,
		Verbose: *verbose,
	}

	// Create API client
	client := NewAPIClient(config)

	// Run the full test
	err := client.RunFullTest(*repoURL, *accessToken)
	if err != nil {
		fmt.Printf("Test failed: %v\n", err)
		os.Exit(1)
	}
}

func runEndpointTest(endpoint string) {
	// Call the endpoint testing function from endpoint_tests.go
	RunEndpointTest(endpoint)
}