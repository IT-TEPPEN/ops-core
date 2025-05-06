package main

import (
	"flag"
	"fmt"
	"os"
)

// Command-line arguments for specific endpoint testing
func parseEndpointTestArgs() (string, string, string, string, bool) {
	// Define command line flags
	cmdEndpoint := flag.NewFlagSet("endpoint", flag.ExitOnError)
	
	// Common arguments
	baseURL := cmdEndpoint.String("baseURL", "http://localhost:8080", "Base URL of the API")
	repoID := cmdEndpoint.String("repoId", "", "Repository ID for operations that need it")
	verbose := cmdEndpoint.Bool("verbose", false, "Enable verbose logging")
	
	// Arguments for specific operations
	repoURL := cmdEndpoint.String("repoURL", "https://github.com/docker/compose", "URL for repository registration")
	accessToken := cmdEndpoint.String("token", "", "Access token for private repositories")
	
	// Parse arguments
	cmdEndpoint.Parse(os.Args[2:])
	
	return *baseURL, *repoID, *repoURL, *accessToken, *verbose
}

// RunEndpointTest executes a specific API endpoint test
func RunEndpointTest(endpoint string) {
	baseURL, repoID, repoURL, accessToken, verbose := parseEndpointTestArgs()
	
	// Create API client with the configuration
	config := Config{
		BaseURL: baseURL,
		Verbose: verbose,
	}
	client := NewAPIClient(config)
	
	var err error
	
	// Execute the appropriate test based on the specified endpoint
	switch endpoint {
	case "register":
		_, err = client.TestRegisterRepository(repoURL, accessToken)
	case "list":
		_, err = client.TestListRepositories()
	case "get":
		if repoID == "" {
			fmt.Println("Error: Repository ID is required for this operation")
			os.Exit(1)
		}
		_, err = client.TestGetRepository(repoID)
	case "update-token":
		if repoID == "" || accessToken == "" {
			fmt.Println("Error: Repository ID and access token are required for this operation")
			os.Exit(1)
		}
		err = client.TestUpdateAccessToken(repoID, accessToken)
	case "files":
		if repoID == "" {
			fmt.Println("Error: Repository ID is required for this operation")
			os.Exit(1)
		}
		_, err = client.TestListRepositoryFiles(repoID)
	case "select-files":
		if repoID == "" {
			fmt.Println("Error: Repository ID is required for this operation")
			os.Exit(1)
		}
		// First get files to select from
		files, filesErr := client.TestListRepositoryFiles(repoID)
		if filesErr != nil {
			fmt.Printf("Failed to list files: %v\n", filesErr)
			os.Exit(1)
		}
		
		// Select up to 3 markdown files
		var markdownFiles []string
		for _, file := range files {
			if file.Type == "file" && len(markdownFiles) < 3 {
				if strings.HasSuffix(file.Path, ".md") {
					markdownFiles = append(markdownFiles, file.Path)
				}
			}
		}
		
		if len(markdownFiles) == 0 {
			fmt.Println("No markdown files found in the repository")
			os.Exit(0)
		}
		
		err = client.TestSelectRepositoryFiles(repoID, markdownFiles)
	case "markdown":
		if repoID == "" {
			fmt.Println("Error: Repository ID is required for this operation")
			os.Exit(1)
		}
		_, err = client.TestGetSelectedMarkdown(repoID)
	default:
		fmt.Printf("Unknown endpoint: %s\n", endpoint)
		os.Exit(1)
	}
	
	if err != nil {
		fmt.Printf("Test failed: %v\n", err)
		os.Exit(1)
	}
}