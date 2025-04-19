#!/bin/bash
set -e # Exit immediately if a command exits with a non-zero status.

# Create frontend and backend directories if they don't exist
mkdir -p frontend backend

# --- Frontend Initialization (React + Vite + TypeScript + Tailwind CSS) ---
if [ ! -f "frontend/package.json" ]; then
  echo "Initializing frontend with Tailwind CSS..."
  # Use npm create vite, force overwrite if exists (though we check above)
  npm create vite@latest frontend -- --template react-ts
  cd frontend
  # Install Tailwind CSS and dependencies
  npm install --save-dev tailwindcss postcss autoprefixer
  # Initialize Tailwind CSS config files
  npx tailwindcss init -p
  # Configure Tailwind template paths in tailwind.config.js
  echo "/** @type {import('tailwindcss').Config} */
export default {
  content: [
    \"./index.html\",
    \"./src/**/*.{js,ts,jsx,tsx}\",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}" > tailwind.config.js
  # Add Tailwind directives to index.css
  echo "@tailwind base;
@tailwind components;
@tailwind utilities;" > src/index.css

  # Install other dependencies
  npm install
  # Add concurrently for running frontend and backend together easily later
  npm install -D concurrently
  cd ..
else
  echo "Frontend already initialized. Checking Tailwind setup..."
  cd frontend
  if ! npm list tailwindcss > /dev/null 2>&1; then
    echo "Installing Tailwind CSS..."
    npm install --save-dev tailwindcss postcss autoprefixer
    npx tailwindcss init -p
    echo "/** @type {import('tailwindcss').Config} */
export default {
  content: [
    \"./index.html\",
    \"./src/**/*.{js,ts,jsx,tsx}\",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}" > tailwind.config.js
    echo "@tailwind base;
@tailwind components;
@tailwind utilities;" > src/index.css
  fi
  npm install # Ensure dependencies are up-to-date
  cd ..
fi

# --- Backend Initialization (Go + Gin) ---
if [ ! -f "backend/go.mod" ]; then
  echo "Initializing Go backend..."
  cd backend
  # Initialize Go module
  go mod init backend # You might want to replace 'backend' with your actual module path later (e.g., github.com/youruser/yourrepo/backend)
  # Get Gin framework
  go get -u github.com/gin-gonic/gin
  # Create basic folder structure and main.go
  mkdir -p cmd/server
  echo 'package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Simple route
	r.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from Go backend!",
		})
	})

	// Allow all origins (for development)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	port := "8080"
	println("Backend server listening at http://localhost:" + port)
	r.Run(":" + port) // listen and serve on 0.0.0.0:8080
}' > cmd/server/main.go
  cd ..
else
  echo "Backend already initialized. Skipping."
  cd backend
  go mod tidy # Ensure dependencies are up-to-date
  cd ..
fi

echo "Post-create script finished."
echo "-----------------------------------------------------"
echo "To start development:"
echo "1. Open two terminals."
echo "2. In the first terminal, run: cd frontend && npm run dev"
echo "3. In the second terminal, run: cd backend && air cmd/server/main.go"
echo "   (Requires 'air' installed: go install github.com/cosmtrek/air@latest)"
echo "   Alternatively, run: cd backend && go run cmd/server/main.go"
echo "-----------------------------------------------------"

