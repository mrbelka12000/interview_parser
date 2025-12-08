# Wails + Go Integration Guide

This guide covers how to effectively integrate Wails with Go for building desktop applications. Based on your existing `interview_parser_app` project.

## Overview

Wails is a framework that allows you to build desktop applications using Go for the backend and web technologies (HTML/CSS/JS) for the frontend. It provides a seamless bridge between Go and JavaScript.

## Project Structure

Your current Wails project follows this structure:

```
interview_parser_app/
├── app.go              # Main application logic and methods
├── main.go             # Application entry point and Wails configuration
├── wails.json          # Wails project configuration
├── go.mod              # Go module dependencies
├── frontend/           # Vue.js frontend application
│   ├── src/
│   │   ├── components/
│   │   ├── App.vue
│   │   └── main.js
│   ├── dist/           # Built frontend assets
│   └── package.json    # Frontend dependencies
└── wailsjs/            # Auto-generated bindings
    ├── go/
    └── runtime/
```

## Core Components

### 1. Main Application Structure (app.go)

```go
package main

import (
    "context"
    "fmt"
)

// App struct - main application logic
type App struct {
    ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
    return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
}

// Example method exposed to frontend
func (a *App) Greet(name string) string {
    return fmt.Sprintf("Hello %s!", name)
}

// Additional service struct
type AppCheck struct {
    ctx context.Context
}

func NewAppCheck() *AppCheck {
    return &AppCheck{}
}

func (a *AppCheck) startup(ctx context.Context) {
    a.ctx = ctx
}

func (c *AppCheck) Check(name string) string {
    return "Checked: " + name
}
```

### 2. Application Entry Point (main.go)

```go
package main

import (
    "embed"
    
    "github.com/wailsapp/wails/v2"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    // Create instances of your app structures
    app := NewApp()
    appCheck := NewAppCheck()
    
    // Create application with options
    err := wails.Run(&options.App{
        Title:  "interview_parser_app",
        Width:  1024,
        Height: 768,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
        OnStartup:        app.startup,
        Bind: []interface{}{
            app,      // Bind App struct methods
            appCheck, // Bind AppCheck struct methods
        },
    })
    
    if err != nil {
        println("Error:", err.Error())
    }
}
```

### 3. Frontend Integration (Vue.js)

```vue
<script setup>
import {reactive} from 'vue'
import {App, AppCheck} from '../../wailsjs/go/main/App'

const data = reactive({
  name: "",
  resultText: "Please enter your name below 👇",
})

function greet() {
  App.Greet(data.name).then(result => {
    data.resultText = result
  })
}

function check() {
  AppCheck.Check(data.name).then(result => {
    data.resultText = result
  })
}
</script>

<template>
  <main>
    <div id="result" class="result">{{ data.resultText }}</div>
    <div id="input" class="input-box">
      <input id="name" v-model="data.name" autocomplete="off" class="input" type="text"/>
      <button class="btn" @click="greet">Greet</button>
      <button class="btn" @click="check">Check</button>
    </div>
  </main>
</template>
```

## Key Integration Concepts

### 1. Method Binding

Any public method in a bound Go struct becomes available to the frontend:

```go
// In Go
func (a *App) ProcessData(data string) (string, error) {
    // Process data
    return result, nil
}
```

```javascript
// In JavaScript
import {App} from '../../wailsjs/go/main/App'

App.ProcessData("some data").then(result => {
    console.log(result);
}).catch(err => {
    console.error(err);
});
```

### 2. Context Usage

The context provided in startup allows you to use Wails runtime functions:

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

func (a *App) ShowDialog() {
    runtime.Dialog(a.ctx, runtime.DialogOptions{
        Title: "Dialog",
        Message: "Hello from Go!",
    })
}
```

### 3. Data Types

Wails automatically converts between Go and JavaScript types:

| Go Type | JavaScript Type |
|---------|-----------------|
| string | string |
| int, int64 | number |
| float64 | number |
| bool | boolean |
| []T | Array |
| map[string]T | Object |
| struct | Object |
| error | throws Error |

### 4. Asynchronous Operations

Go methods that return results are called asynchronously from JavaScript:

```go
// Go
func (a *App) LongOperation() string {
    time.Sleep(2 * time.Second)
    return "Done"
}
```

```javascript
// JavaScript
async function doLongOperation() {
    try {
        const result = await App.LongOperation();
        console.log(result);
    } catch (err) {
        console.error(err);
    }
}
```

## Advanced Integration Patterns

### 1. File Operations

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

func (a *App) OpenFile() (string, error) {
    selectedFile, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
        Title: "Select File",
        Filters: []runtime.FileFilter{
            {
                DisplayName: "Text Files",
                Pattern:     "*.txt",
            },
        },
    })
    
    if err != nil {
        return "", err
    }
    
    return selectedFile, nil
}
```

### 2. Database Integration

```go
import "database/sql"
import _ "github.com/mattn/go-sqlite3"

type DatabaseService struct {
    ctx context.Context
    db  *sql.DB
}

func (ds *DatabaseService) InitDB() error {
    var err error
    ds.db, err = sql.Open("sqlite3", "./app.db")
    return err
}

func (ds *DatabaseService) QueryData() ([]map[string]interface{}, error) {
    rows, err := ds.db.Query("SELECT * FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var results []map[string]interface{}
    for rows.Next() {
        var id int
        var name string
        
        err := rows.Scan(&id, &name)
        if err != nil {
            return nil, err
        }
        
        results = append(results, map[string]interface{}{
            "id":   id,
            "name": name,
        })
    }
    
    return results, nil
}
```

### 3. HTTP Client Integration

```go
import "net/http"
import "encoding/json"

func (a *App) FetchData(url string) (map[string]interface{}, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        return nil, err
    }
    
    return result, nil
}
```

### 4. Event System

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

func (a *App) StartBackgroundTask() {
    go func() {
        for i := 0; i < 5; i++ {
            time.Sleep(1 * time.Second)
            runtime.EventsEmit(a.ctx, "progress", i+1)
        }
        runtime.EventsEmit(a.ctx, "completed", true)
    }()
}
```

```javascript
// JavaScript
import {EventsOn} from '../../wailsjs/runtime/runtime'

EventsOn("progress", (progress) => {
    console.log(`Progress: ${progress}`);
});

EventsOn("completed", () => {
    console.log("Task completed!");
});
```

## Development Workflow

### 1. Development Mode

```bash
cd interview_parser_app
wails dev
```

This starts the application with hot reload for both Go and frontend changes.

### 2. Building for Production

```bash
wails build
```

This creates a production binary with embedded frontend assets.

### 3. Frontend Development

```bash
cd frontend
npm install
npm run dev  # For frontend-only development
npm run build  # Build frontend assets
```

## Best Practices

### 1. Error Handling

Always return errors from Go methods:

```go
func (a *App) ProcessData(data string) (string, error) {
    if data == "" {
        return "", errors.New("data cannot be empty")
    }
    
    // Process data
    return result, nil
}
```

### 2. Input Validation

Validate inputs in Go methods:

```go
func (a *App) ValidateEmail(email string) (bool, error) {
    if email == "" {
        return false, errors.New("email is required")
    }
    
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email), nil
}
```

### 3. Resource Management

Use defer for cleanup:

```go
func (a *App) ProcessFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    // Process file
    return nil
}
```

### 4. Structured Logging

```go
import "log"

func (a *App) LogOperation(operation string, data interface{}) {
    log.Printf("[%s] %v", operation, data)
}
```

## Common Pitfalls and Solutions

### 1. Method Not Found
- Ensure methods are exported (start with capital letter)
- Verify structs are properly bound in main.go
- Run `wails generate module` if bindings are outdated

### 2. Context Issues
- Always store the context in startup method
- Use the stored context for runtime calls

### 3. Frontend Build Issues
- Ensure frontend assets are built before running wails build
- Check package.json dependencies

### 4. Cross-Platform Considerations
- Test on target platforms
- Handle platform-specific paths and behaviors
- Use runtime.Environment().Platform for conditional logic

## Integration with Your Interview Parser

Based on your existing codebase, here's how to integrate your interview parsing logic:

```go
// Add to app.go
func (a *App) ProcessInterview(audioPath string) (map[string]interface{}, error) {
    // Use your existing client functionality
    transcribe, err := client.TranscribeAudio(audioPath)
    if err != nil {
        return nil, err
    }
    
    analysis, err := client.AnalyzeTranscript(transcribe)
    if err != nil {
        return nil, err
    }
    
    return map[string]interface{}{
        "transcript": transcribe,
        "analysis":   analysis,
    }, nil
}

func (a *App) SaveResults(data map[string]interface{}, filename string) error {
    // Use your existing saver functionality
    return saver.SaveToFile(data, filename)
}
```

This comprehensive guide should help you effectively integrate Wails with Go for your interview parser application. The framework provides excellent tooling for building cross-platform desktop applications with Go's powerful backend capabilities.
