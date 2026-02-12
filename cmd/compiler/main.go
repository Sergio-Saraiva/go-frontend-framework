package main

import (
	"fmt"
	"go-frontend-framework/compiler"
	"os"
	"path/filepath"
	"strings"
	// Adjust this import to match your module name (e.g., go.mod)
)

func main() {
	// 1. Open the Template
	inputPath := "src/app/about/about.html"
	outputPath := "src/app/about/about_gen.go"
	cssPath := strings.Replace(inputPath, ".html", ".css", 1)

	// 1. Parse
	file, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	root, err := compiler.Parse(file)
	if err != nil {
		panic(err)
	}

	cssContent := ""
	cssBytes, err := os.ReadFile(cssPath)
	if err == nil {
		cssContent = string(cssBytes)
		fmt.Printf("🎨 Found CSS: %s\n", cssPath)
	} else {
		fmt.Println("⚠️  No CSS found (skipping)")
	}

	// 2. Determine Package Name (directory name)
	// e.g., src/app/counter -> "counter"
	absPath, _ := filepath.Abs(inputPath)
	dirName := filepath.Base(filepath.Dir(absPath))

	// 3. Generate Content
	code := compiler.GenerateFullFile(*root, dirName, cssContent)

	// 4. Write File
	err = os.WriteFile(outputPath, []byte(code), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("✅ Compiled %s -> %s\n", inputPath, outputPath)
}
