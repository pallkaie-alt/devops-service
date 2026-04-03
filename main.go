package main

import (
	"log"
	"net/http"
    "os"
	"os/signal"
    "time"
	"context"
	"syscall"
)

type App struct {
    Port            string
    ResponseMessage string
    AllowOrigin     string
}

/** statusWriter to reliably capture HTTP status codes for logging,
handles WriteHeader and default Write cases. */

type statusWriter struct {
    http.ResponseWriter
    status int
}

func (a *App) helloWorld(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", a.AllowOrigin)
    w.Write([]byte(a.ResponseMessage))
}
func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK) // Tagastab HTTP 200 
    w.Write([]byte("ok"))
}
func (a *App) readyHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK) 
    w.Write([]byte("ready"))
}

// Overriding the WriteHeader method to save the code before sending it.

func (sw *statusWriter) WriteHeader(statusCode int) {
    sw.status = statusCode
    sw.ResponseWriter.WriteHeader(statusCode)
}

// Handle case where Write is called without WriteHeader.

func (sw *statusWriter) Write(b []byte) (int, error) {
    if sw.status == 0 {
        sw.status = http.StatusOK
    }
    return sw.ResponseWriter.Write(b)
}


func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		//  Create instance and mri by defult always 200 
        sw := &statusWriter {
            ResponseWriter: w,
            status:         http.StatusOK, 
        }
		// Call the next handler in chain.
		next.ServeHTTP(sw, r)

		duration := time.Since(start)

		log.Printf("METHOD: %s | PATH: %s | STATUS: %d | DURATION: %v",
		r.Method, r.URL.Path, sw.status, duration)
			})
	}

// Panic recovery (prevents whole-server crash)
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("PANIC recovered: %v\n", rec)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func SecurityMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY") 
        
        next.ServeHTTP(w, r)
    })
}

// CORSMiddleware adds the necessary headers and responds to OPTIONS requests

func CORSMiddleware(next http.Handler, origin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Seame lubatud päritolu vastavalt konfiguratsioonile 
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	app := &App{
		Port:		 		getEnv("PORT", "8000"),
		ResponseMessage: 	getEnv("RESPONSE_MESSAGE", "Service request succeeded!"),
		AllowOrigin:		getEnv("ALLOW_ORIGIN", "*"), 
}	


	mux := http.NewServeMux()
	


    mux.HandleFunc("/", app.helloWorld)
	mux.HandleFunc("/health", app.healthHandler)
    mux.HandleFunc("/ready", app.readyHandler)

/** Chain of func:
Srequest-> Recovery -> Logging-> Security -> CORS -> mux routing;
 and vice-versa to give response out */

loggedMux := LoggingMiddleware(mux)

corsMux := CORSMiddleware(loggedMux, app.AllowOrigin)

secureMux := SecurityMiddleware(corsMux)

finalHandler := RecoveryMiddleware(secureMux)

	server := &http.Server{
		Addr:	":" + app.Port,
		Handler: finalHandler,		
	}

	go func() {
		log.Printf("Server running at localhost: %s ", app.Port)
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("ListenAndServe error: %v", err)
    }
}()
	
// Graceful shutdown wait time to complete active requests

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    <-quit
    log.Println("Closing server in 9 seconds...")

	
    ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
    defer cancel() 

	
    if err := server.Shutdown(ctx); err != nil {
        log.Printf("Server shutdown error: %v", err)
    } else {
        log.Println("Server was shut down safely.")
    }
}

// Helper function for managing default values
func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
