package hotreload

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

// InjectScriptMiddleware wraps an HTTP handler and injects the SSE client script into HTML responses
func InjectScriptMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip non-GET requests
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		// Skip non-HTML requests
		if !strings.HasSuffix(r.URL.Path, ".html") && r.URL.Path != "/" {
			next.ServeHTTP(w, r)
			return
		}

		// Capture the response
		recorder := &responseRecorder{ResponseWriter: w}
		next.ServeHTTP(recorder, r)

		// Check if response is HTML
		contentType := recorder.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/html") && recorder.statusCode != 0 {
			// If content type is not HTML, write original response
			// Write headers first
			if recorder.statusCode != 0 {
				w.WriteHeader(recorder.statusCode)
			}
			// Then write body
			w.Write(recorder.body.Bytes())
			return
		}

		// Inject script before </body> or </head> tag
		body := recorder.body.String()
		script := GetClientScript()

		var injected string
		if strings.Contains(body, "</body>") {
			injected = strings.Replace(body, "</body>", script+"</body>", 1)
		} else if strings.Contains(body, "</html>") {
			injected = strings.Replace(body, "</html>", script+"</html>", 1)
		} else {
			injected = body + script
		}

		// Copy headers from recorder to original response
		for k, v := range recorder.Header() {
			for _, val := range v {
				w.Header().Add(k, val)
			}
		}

		// Update Content-Length to reflect new body size
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(injected)))

		// Write headers (use 200 if no status code was set)
		statusCode := recorder.statusCode
		if statusCode == 0 {
			statusCode = 200
		}
		w.WriteHeader(statusCode)

		// Write injected body
		w.Write([]byte(injected))
	})
}

// responseRecorder captures the response status code and body
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

// WriteHeader captures the status code without forwarding it
func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

// Write captures the response body
func (r *responseRecorder) Write(b []byte) (int, error) {
	return r.body.Write(b)
}
