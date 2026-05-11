# Main.go
> **Change_01**

Before
```go
handlers.InitializeTemplates("templates")
```
After
```go
logger := log.New(os.Stdout, "APP: ", log.LstdFlags|log.Lshortfile) // Set up a logger
	templates, err := handlers.InitializeTemplates("templates", logger) // Initialize templates
	if err != nil {
		logger.Fatalf("Failed to initialize templates: %v", err)
	}
	handlers.SetTemplates(templates)     // Pass the templates to your handlers or use globally        
	logger.Println("Server started successfully") // Log for starting your server successfully
```

  * <code>logger := log.New(os.Stdout, "APP: ", log.LstdFlags|log.Lshortfile)</code> \
  The log.New function in Go creates a new logger instance that accepts three values:
    * Where to print 
    * Prefix before the output
    * Format of printing
  
  ---
