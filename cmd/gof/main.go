package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sergio-Saraiva/go-frontend-framework/compiler"
	"github.com/fsnotify/fsnotify"
)

const (
	Port         = "8080"
	BuildCmd     = "go build -o main.wasm ."
	CompilerPath = "cmd/compiler/main.go"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "build":
		runFullBuild()
	case "serve":
		startDevServer()
	case "generate":
		fmt.Println("TODO: Component generation")
	case "new":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing project name.")
			fmt.Println("Usage: gof new <project-name>")
			return
		}
		scaffoldProject(os.Args[2])
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Println("Go Frontend CLI (gof)")
	fmt.Println("  new <name> Create a new project")
	fmt.Println("  serve    Start dev server with live reload")
	fmt.Println("  build    Compile project to WASM")
	fmt.Println("  generate Generate components")
}

func runFullBuild() error {
	fmt.Print("Compiling Templates... ")

	if err := compiler.Run(); err != nil {
		fmt.Printf("\nCompiler Error: %v\n", err)
		return err
	}

	fmt.Print("Building WASM... ")

	buildCmd := exec.Command("sh", "-c", "GOOS=js GOARCH=wasm go build -o main.wasm .")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		fmt.Printf("\nBuild Error: %v\n", err)
		return err
	}
	fmt.Println("✅")
	return nil
}

func startDevServer() {
	if err := runFullBuild(); err != nil {
		fmt.Println("Build failed, but starting server anyway...")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err := filepath.Walk("./src", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	go func() {
		debounce := time.NewTimer(100 * time.Millisecond)
		debounce.Stop()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fileName := event.Name
				if strings.HasSuffix(fileName, "_gen.go") ||
					strings.HasSuffix(fileName, ".wasm") ||
					strings.HasPrefix(filepath.Base(fileName), ".") {
					continue
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					debounce.Reset(100 * time.Millisecond)
				}
			case <-debounce.C:
				fmt.Println("File changed detected. Rebuilding...")
				runFullBuild()
				fmt.Println("waiting for changes...")
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// 3. Start HTTP Server
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	fmt.Printf("Server running at http://localhost:%s\n", Port)
	log.Fatal(http.ListenAndServe(":"+Port, nil))
}

func scaffoldProject(name string) {
	fmt.Printf("Scaffolding new project: %s...\n", name)

	dirs := []string{
		name,
		filepath.Join(name, "src", "app", "home"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	files := map[string]string{
		filepath.Join(name, "go.mod"): fmt.Sprintf(`module %s

go 1.25

require github.com/Sergio-Saraiva/go-frontend-framework v0.1.6
`, name),

		filepath.Join(name, "index.html"): `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Frontend App</title>
    <script src="wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
            go.run(result.instance);
        });
    </script>
    <style>
        body { font-family: sans-serif; margin: 0; padding: 20px; }
    </style>
</head>
<body>
    </body>
</html>`,

		filepath.Join(name, "main.go"): `//go:build js && wasm
package main

import (
	"syscall/js"
	"github.com/Sergio-Saraiva/go-frontend-framework/router"
	
	home_pkg "` + name + `/src/app/home"
)

func main() {
	// Keep WASM alive
	c := make(chan struct{}, 0)

	// Register Routes
	router.Register("/", home_pkg.New)

	// Mount the router to the DOM
	root := router.Outlet()
	js.Global().Get("document").Get("body").Call("appendChild", root)

	// Start listening for URL changes
	router.Start()

	<-c
}`,

		filepath.Join(name, "src", "app", "home", "home.go"): `//go:build js && wasm
package home

import (
	"github.com/Sergio-Saraiva/go-frontend-framework/component" 
	"github.com/Sergio-Saraiva/go-frontend-framework/signal"
)

type Component struct {
	component.BaseComponent
	Title *signal.Signal[string]
}

func New() component.Interface {
	c := &Component{
		Title: signal.New("Welcome to Go Frontend! 🐹"),
	}
	c.BaseComponent.Init(c)
	return c
}
`,

		filepath.Join(name, "src", "app", "home", "home.html"): `<div class="hero">
    <h1>{{ Title.Get() }}</h1>
    <p>Your Go WebAssembly application is running.</p>
    <input type="text" bind-value="Title" placeholder="Edit title..." />
</div>

<style>
    .hero {
        text-align: center;
        padding: 50px;
        background: #f4f4f9;
        border-radius: 8px;
    }
    input {
        padding: 10px;
        font-size: 16px;
        margin-top: 20px;
    }
</style>`,
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			log.Fatalf("Failed to write file %s: %v", path, err)
		}
	}

	fmt.Print("Fetching wasm_exec.js... ")
	copyWasmExec(name)

	fmt.Printf("\nProject '%s' created successfully!\n\n", name)
	fmt.Printf("Next steps:\n")
	fmt.Printf("  cd %s\n", name)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  gof serve\n")
}

func copyWasmExec(projectDir string) {
	out, err := exec.Command("go", "env", "GOROOT").Output()
	if err == nil {
		goroot := strings.TrimSpace(string(out))
		localPath := filepath.Join(goroot, "misc", "wasm", "wasm_exec.js")

		if input, err := os.ReadFile(localPath); err == nil {
			destPath := filepath.Join(projectDir, "wasm_exec.js")
			os.WriteFile(destPath, input, 0644)
			fmt.Println("(Copied from GOROOT)")
			return
		}
	}

	fmt.Print("(Downloading from GitHub...) ")
	verOut, err := exec.Command("go", "env", "GOVERSION").Output()
	if err != nil {
		log.Fatalf("\nFailed to determine Go version: %v", err)
	}
	goVer := strings.TrimSpace(string(verOut))

	urlsToTry := []string{
		fmt.Sprintf("https://raw.githubusercontent.com/golang/go/%s/misc/wasm/wasm_exec.js", goVer),
		fmt.Sprintf("https://raw.githubusercontent.com/golang/go/%s/misc/wasm/wasm_exec.js", strings.TrimSuffix(goVer, ".0")),
		"https://raw.githubusercontent.com/golang/go/master/misc/wasm/wasm_exec.js",
	}

	var resp *http.Response
	for _, u := range urlsToTry {
		r, err := http.Get(u)
		if err == nil && r.StatusCode == http.StatusOK {
			resp = r
			break
		}
		if r != nil {
			r.Body.Close()
		}
	}

	if resp == nil {
		log.Fatalf("\nFailed to download wasm_exec.js. All fallback URLs returned 404.")
	}
	defer resp.Body.Close()

	destPath := filepath.Join(projectDir, "wasm_exec.js")
	outf, err := os.Create(destPath)
	if err != nil {
		log.Fatalf("\nFailed to create wasm_exec.js: %v", err)
	}
	defer outf.Close()

	if _, err = io.Copy(outf, resp.Body); err != nil {
		log.Fatalf("\nFailed to write wasm_exec.js: %v", err)
	}
}
