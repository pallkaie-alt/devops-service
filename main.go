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

	 // Handlers for endpoints. 

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
	// Creating wrapper, 
	// statusWriter mhib tavalise ResponseWriter-i, et me saaksime staatust logida

type statusWriter struct {
    http.ResponseWriter
    status int
}
// Me kirjutame üle WriteHeader meetodi, et salvestada kood enne selle teelesaatmist

func (sw *statusWriter) WriteHeader(statusCode int) {
    sw.status = statusCode
    sw.ResponseWriter.WriteHeader(statusCode)
}
// Handle the rare case where Write is called without WriteHeader (extra safety):
func (sw *statusWriter) Write(b []byte) (int, error) {
    if sw.status == 0 {
        sw.status = http.StatusOK
    }
    return sw.ResponseWriter.Write(b)
}

	func LoggingMiddleware(next http.Handler)
 http.Handler {
	return http.HandlerFunc(func(w
		 http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 1. Loo instants ja mri vaike-staatuskoodiks 200
        sw := &statusWriter {
            ResponseWriter: w,
            status:         http.StatusOK, // ehk 200
        }

		// 2. Anna päring edasi, kasutades meie uut "sw" wrapperit
		next.ServeHTTP(sw, r)

		// 3. Nüüd on päring tehtud ja sw.status sisaldab õiget koodi
		duration := time.Since(start)

		logger.Printf("METHOD: %s | PATH: %s | STATUS: %d | DURATION: %v",
		r.Method, r.URL.Path, sw.status, duration)
			})
	}




func main() {

app := &App{
	Port:		 		getEnv("Port", "8000"),
	ResponseMessage: 	getEnv("RESPONSE_MESSAGE", "Service request succeeded!"),
	AllowOrigin:		getEnv("ALLOW_ORIGIN", "*"), 
}	


// Loo ruuter ServeMux
	mux := http.NewServeMux()
	
	// Registreeri handlerid ja mhi need Middleware'i sisse [cite: 43-49, 53]
    // Kuna need on meetodid, kasutame app.helloWorld

    mux.Handle("/", LoggingMiddleware(http.HandlerFunc(app.helloWorld)))
	mux.Handle("/health", LoggingMiddleware(http.HandlerFunc(app.healthHandler)))
    mux.Handle("/ready", LoggingMiddleware(http.HandlerFunc(app.readyHandler)))

    // Wrap mux with middleware
    wrapped := LoggingMiddleware(mux)

	log.Println("Server running on :8080")
    http.ListenAndServe(":8080", wrapped)

	// Graceful shutdown func!
// Käivita server eraldi gorutiinis
	go func() {
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("ListenAndServe error: %v", err)
    }
}()
	
// Ootame graceful shutdown signaali (Ctrl+C või kill)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Closing server in 10 seconds...")

	// creating new 'context' with short waiting time.
    ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
    defer cancel() // hea praktika on vabastada

	// Graceful shutdown – annab serverile 10 sekundit aktiivsete päringute lõpetamiseks
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("Server shutdown error: %v", err)
    } else {
        log.Println("Server was shut down safely.")
    }
}
// Abifunktsioon vaikeväärtuste haldamiseks
func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

