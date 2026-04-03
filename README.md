

# Go HTTP Service 

This is a minimal, production-ready HTTP service.    
The service is optimized for security, observability, and smooth management.

##  Functionality

* **Dynamic Configuration**: Configurable via environment variables (`PORT`, `RESPONSE_MESSAGE`, `ALLOW_ORIGIN`).
* **Structured Endpoints**:
    * `GET /`:        Returns the configured greeting message.
    * `GET /health`:  Application liveness probe returning "ok".
    * `GET /ready`:   Application readiness probe returning "ready".
* **Observability**: Method, path, HTTP status code, and processing duration are logged for every request.
* **Resilience**: Automatic Panic Recovery to prevent server crashes.
* **Graceful Shutdown**: The service waits up to 9 seconds for in-flight requests to complete before terminating.

## 🛠 Technical Architecture

The application uses the **"Middleware Onion"** pattern and a **Dependency Injection** style structure (`App struct`) to keep the code testable and clean.

### Middleware Chain (from inside out):
1.  **Mux**: Standard Go router.
2.  **CORS**: Handles `OPTIONS` preflight requests and adds allowed origin headers.
**CORS Security:** The service supports dynamic origin configuration via the ALLOW_ORIGIN environment variable. While the default is set to `*` for ease of testing, it is recommended to restrict this to specific trusted domains in a production environment to prevent unauthorized cross-origin requests.
3.  **Security**: Adds critical security headers (`X-Frame-Options`, `CSP`, `X-Content-Type-Options`).
4.  **Logging**: Measures and logs the entire request lifecycle.
5.  **Recovery**: The outermost layer that catches all unexpected errors.

##  Containerization (Docker)

The container is built following the **Multi-stage build** principle to ensure minimal image size and security.

* **Builder**: `golang:1.21-alpine`.
* **Final Image**: `alpine:latest` (contains only the static binary).
* **Security**: The application runs as a non-root user (nonroot, UID 10001) for enhanced   container security.

### Build and Run Instructions
```bash
# 1. Build the image:
docker build -t devops-service .

# 2. Run the container (with default values)
docker run -p 8000:8000 devops-service

# 3. Run the container with custom settings
docker run \
  -p 8080:8080 \
  -e PORT=8080 \
  -e RESPONSE_MESSAGE='Hello World!' \
  devops-service

**For stopping this process later**
docker ps
docker stop <container_id>

### Running in Detached Mode
To run the container in the background (detached mode), use the `-d` flag. This allows the service to run without occupying your terminal session.

```bash
# Start the container in the background
docker run -d -p 8000:8000 --name devops-app devops-service
```
## 🔍 Testing the Endpoints

Once the container is running, you can verify the service using `curl`.

* **Root Endpoint (Welcome Message)**
```bash
curl -v http://localhost:8000/    
```
* **Health Check Endpoint**
```
curl -v http://localhost:8000/health
```
* **Readiness Check Endpoint**
```
curl -v http://localhost:8000/ready
```
## Technical Decisions

* **Pure Go Standard Library**: Avoided external dependencies to keep the code lightweight and auditable.
* **Incremental Chain**: Used a step-by-step middleware chain in the `main` function to improve readability and avoid a "bracket maze."
* **CGO_ENABLED=0**: Static compilation ensures the application does not depend on host system libraries.

##  Collaboration Note
This Go program was developed with the assistance of artificial intelligence (Gemini) as a consultative partner for structuring code, refining idiomatic Go syntax, and implementing best security practices. All final decisions and code verification were performed by me.
 **My first prompt was:** 1. Be a professional instructor for a junior Java developer coding in GoLang. 2. Provide action plan points to start the project described in the file. 3. Do not write code in advance. 4. Minimal template files: 1) go.mod and 2) main.go are provided. 4. Recommend materials and videos from the web that contain the necessary instructions. Building step-by-step I made many trials and errors during coding first Go HTTP Service.  
  Getting Docker working on MacOS 12 Monterey after the coding part was a great challenge, which I managed with Copilot's guidance in almost 3 hours.

Author - Kaie Päll, Junior Developer course student at /kood/Jõhvi 
04.2026