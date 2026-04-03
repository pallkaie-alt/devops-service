

# Go HTTP Service – Cybernetica Trial Task 2026

This is a minimalist yet production-ready HTTP service created according to the trial task instructions. The service is optimized for security, observability, and smooth management.

##  Functionality

* **Dynamic Configuration**: Configurable via environment variables (`PORT`, `RESPONSE_MESSAGE`, `ALLOW_ORIGIN`).
* **Structured Endpoints**:
    * `GET /`: Returns the configured greeting message.
    * `GET /health`: Application health check (Liveness probe).
    * `GET /ready`: Application readiness check (Readiness probe).
* **Observability**: Method, path, HTTP status code, and processing duration are logged for every request.
* **Resilience**: Automatic Panic Recovery to prevent server crashes.
* **Graceful Shutdown**: The service waits up to 9 seconds for in-flight requests to complete before terminating.

## 🛠 Technical Architecture

The application uses the **"Middleware Onion"** pattern and a **Dependency Injection** style structure (`App struct`) to keep the code testable and clean.

### Middleware Chain (from inside out):
1.  **Mux**: Standard Go router.
2.  **CORS**: Handles `OPTIONS` requests and adds allowed origin headers.
3.  **Security**: Adds critical security headers (`X-Frame-Options`, `CSP`, `X-Content-Type-Options`).
4.  **Logging**: Measures and logs the entire request lifecycle.
5.  **Recovery**: The outermost layer that catches all unexpected errors.

##  Containerization (Docker)

The container is built following the **Multi-stage build** principle to ensure minimal image size and security.

* **Builder**: `golang:1.21-alpine`.
* **Final Image**: `alpine:latest` (contains only the static binary).
* **Security**: The application runs as a non-root user with ID 10001.

### Instructions:
```bash
# 1. Build the container
docker build -t cyber-service .

# 2. Run the container (with default values)
docker run -p 8000:8000 cyber-service

# 3. Run the container with custom settings
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e RESPONSE_MESSAGE="Hello Cybernetica!" \
  cyber-service
```

## Technical Decisions

* **Pure Go Standard Library**: Avoided external dependencies to keep the code lightweight and auditable.
* **Incremental Chain**: Used a step-by-step middleware chain in the `main` function to improve readability and avoid a "bracket maze."
* **CGO_ENABLED=0**: Static compilation ensures the application does not depend on host system libraries.

##  Collaboration Note
This solution was developed with the assistance of artificial intelligence (Gemini) as a consultative partner for structuring code, refining idiomatic Go syntax, and implementing best security practices. All final decisions and code verification were performed by me.


Author - Kaie Päll, Junior Developer course student at /kood/Jõhvi 
04.2026