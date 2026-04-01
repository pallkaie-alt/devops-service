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
// statusWriter mhib tavalise ResponseWriter-i, et me saaksime staatust logida

type statusWriter struct {
    http.ResponseWriter
    status int
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
// Creating wrapper methods and MiddleWare 
// Me kirjutame üle WriteHeader meetodi, et salvestada kood enne selle teelesaatmist

func (sw *statusWriter) WriteHeader(statusCode int) {
    sw.status = statusCode
    sw.ResponseWriter.WriteHeader(statusCode)
}
// Handle the rare case where Write is called without WriteHeader (extra safety)

func (sw *statusWriter) Write(b []byte) (int, error) {
    if sw.status == 0 {
        sw.status = http.StatusOK
    }
    return sw.ResponseWriter.Write(b)
}
//* Middleware *//

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

		// Nüüd on päring tehtud ja sw.status sisaldab õiget koodi
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
// CORSMiddleware lisab vajalikud päised ja vastab OPTIONS päringutele

func CORSMiddleware(next http.Handler, origin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Seame lubatud päritolu vastavalt konfiguratsioonile 
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Kui on OPTIONS päring (preflight), vastame kohe 204-ga
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Muul juhul liigume järgmise mähise juurde
		next.ServeHTTP(w, r)
	})
}

func main() {
	app := &App{
		Port:		 		getEnv("PORT", "8000"),
		ResponseMessage: 	getEnv("RESPONSE_MESSAGE", "Service request succeeded!"),
		AllowOrigin:		getEnv("ALLOW_ORIGIN", "*"), 
}	

// Loo ruuter ServeMux
	mux := http.NewServeMux()
	
	// Registreeri handler meetodid, kasutame app.helloWorld

    mux.HandleFunc("/", app.helloWorld)
	mux.HandleFunc("/health", app.healthHandler)
    mux.HandleFunc("/ready", app.readyHandler)

    // Wrap mux with middleware. Single logging pass.
// chained functions, where recovery midlw. protects also logging.
// 1. Kõigepealt mähime mux-i logimisse
loggedMux := LoggingMiddleware(mux)

// 2. Nüüd mähime logitud mux-i CORS-i sisse. 
// SIIN ongi kaks argumenti: 
// 1) loggedMux (http.Handler) 
// 2) app.AllowOrigin (string)
corsMux := CORSMiddleware(loggedMux, app.AllowOrigin)

// 3. Ja kõige lõpuks lisame RecoveryMiddleware, et see kaitseks kogu ahelat
finalHandler := RecoveryMiddleware(corsMux)

	//  Loo eraldi Serveri objekt, et saaksime seda hiljem sujuvalt sulgeda.
    // Kasuta porti oma App struktuurist, lisades ette kooloni!

	server := &http.Server{
		Addr:	":" + app.Port,
		Handler: finalHandler,		
	}

// Käivita server eraldi gorutiinis
	go func() {
		log.Printf("Server running %s...", app.Port)
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("ListenAndServe error: %v", err)
    }
}()
	
// Ootame shutdown signaali (Ctrl+C või kill)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    <-quit
    log.Println("Closing server in 9 seconds...")

	// creating new 'context' for smooth shutdown in 9 seconds.
    ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
    defer cancel() // hea praktika on vabastada

	// Graceful shutdown ooteaeg aktiivsete päringute lõpetamiseks
    if err := server.Shutdown(ctx); err != nil {
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
 