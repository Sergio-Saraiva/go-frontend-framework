# Go Frontend Framework (gof) 🐹✨

A blazing-fast, WebAssembly-powered frontend framework for Go. Build reactive Single Page Applications (SPAs) using strongly typed Go, HTML components, and a powerful signal-based state management system—no JavaScript required.

## Features
* **🧩 Component-Based Architecture:** Write clean, modular views separating your `.go` logic from your `.html` templates.
* **🎨 Separate CSS Support:** Keep your styles clean by optionally pairing a `.css` file directly with your component.
* **🗺️ SPA Routing:** Built-in router to handle URL changes and seamlessly swap components without reloading the page.
* **⚡ Reactive Signals:** State management built-in. Update a signal, and the DOM updates automatically.
* **💉 Dependency Injection:** Share state and logic across components with auto-wired services (Angular/NestJS style).
* **🛠️ Powerful CLI:** Scaffold, build, and serve your applications with live-reloading out of the box.

---

## 🚀 Quick Start

### 1. Install the CLI globally
Make sure your Go `bin` directory is in your system's PATH.

`bash
go install github.com/Sergio-Saraiva/go-frontend-framework/cmd/gof@latest
`

### 2. Scaffold a new project
Use the CLI to instantly generate a fully configured starter application.

`bash
gof new my-app
cd my-app
`

### 3. Start the Development Server
Fire up the dev server with automatic template compilation and live-reloading.

`bash
gof serve
`
Open your browser to `http://localhost:8080`.

---

## 💻 Writing Components

A component in this framework consists of up to three files grouped together in a folder: a Go file for logic, an HTML file for the view, and an optional CSS file for styling.

**`src/app/home/home.go`**
```go
package home

import (
	"github.com/Sergio-Saraiva/go-frontend-framework/component"
	"github.com/Sergio-Saraiva/go-frontend-framework/signal"
)

type Component struct {
	component.BaseComponent
	
	// Reactive State
	Title *signal.Signal[string]
	Count *signal.Signal[int]
}

func New() component.Interface {
	c := &Component{
		Title: signal.New("Hello, Go WebAssembly!"),
		Count: signal.New(0),
	}
	c.Base.Init(c)
	return c
}

// Event Handler
func (c *Component) Increment() {
	c.Count.Set(c.Count.Get() + 1)
}
```

**`src/app/home/home.html`**
```html
<div class="container">
    <h1>{{ Title.Get() }}</h1>
    
    <input type="text" bind-value="Title" />

    <p>Counter: {{ Count.Get() }}</p>
    <button (click)="Increment()">Add +1</button>
</div>
```

**`src/app/home/home.css` (Optional)**
```css
.container { 
    padding: 20px; 
    font-family: sans-serif; 
}

button {
    background-color: #007bff;
    color: white;
    border: none;
    padding: 10px 15px;
    border-radius: 4px;
    cursor: pointer;
}

button:hover {
    background-color: #0056b3;
}
```

---

## 🗺️ Routing (SPA)



The framework includes a built-in router that mounts your components based on the browser's URL. When the URL changes, the router destroys the old component and seamlessly mounts the new one—no page reload required!

Register your routes inside your `main.go` file:

**`main.go`**
```go
//go:build js && wasm
package main

import (
	"syscall/js"
	"github.com/Sergio-Saraiva/go-frontend-framework/router"
	
	home_pkg "my-app/src/app/home"
	about_pkg "my-app/src/app/about"
)

func main() {
	c := make(chan struct{}, 0)

	// 1. Register your routes
	router.Register("/", home_pkg.New)
	router.Register("/about", about_pkg.New)

	// 2. Mount the router outlet to the DOM
	root := router.Outlet()
	js.Global().Get("document").Get("body").Call("appendChild", root)

	// 3. Start listening for navigation events
	router.Start()

	<-c // Keep WASM alive
}
```

---

## 🛠️ CLI Commands

| Command | Description |
| :--- | :--- |
| `gof new <name>` | Scaffolds a new project with all necessary files and `wasm_exec.js`. |
| `gof serve` | Starts a development server with automatic recompilation and live reload. |
| `gof build` | Compiles the project into a production-ready `main.wasm` binary. |

## 📦 Deployment
To deploy your application, run `gof build`. This will generate your `main.wasm` file. You can then host the project folder statically on any web server (Nginx, Vercel, Netlify, GitHub Pages, etc.) by serving the `index.html`, `main.wasm`, and `wasm_exec.js` files.
