package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

// HeadersResponse represents the JSON response for headers
type HeadersResponse struct {
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	RemoteAddr string              `json:"remote_addr"`
	Headers    map[string][]string `json:"headers"`
}

// ShowHeaders handles requests and displays all HTTP headers
// Supports ?format=json for JSON output and ?format=plain for plain text
func ShowHeaders(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")

	switch format {
	case "json":
		showHeadersJSON(w, r)
	case "plain":
		showHeadersPlain(w, r)
	default:
		showHeadersHTML(w, r)
	}
}

// showHeadersJSON returns headers as JSON
func showHeadersJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := HeadersResponse{
		Method:     r.Method,
		URL:        r.URL.String(),
		RemoteAddr: r.RemoteAddr,
		Headers:    r.Header,
	}

	json.NewEncoder(w).Encode(response)
}

// showHeadersPlain returns headers as plain text
func showHeadersPlain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "URL: %s\n", r.URL.String())
	fmt.Fprintf(w, "Remote Address: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "\n--- Headers ---\n")

	// Sort headers alphabetically for consistent display
	var headerNames []string
	for name := range r.Header {
		headerNames = append(headerNames, name)
	}
	sort.Strings(headerNames)

	for _, name := range headerNames {
		values := r.Header[name]
		for _, value := range values {
			fmt.Fprintf(w, "%s: %s\n", name, value)
		}
	}
}

// showHeadersHTML returns headers as styled HTML
func showHeadersHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, "<html><head><title>Show Headers</title>")
	fmt.Fprintf(w, "<style>")
	fmt.Fprintf(w, "body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }")
	fmt.Fprintf(w, "h1 { color: #333; }")
	fmt.Fprintf(w, "table { border-collapse: collapse; width: 100%%; max-width: 800px; background-color: white; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }")
	fmt.Fprintf(w, "th, td { padding: 12px 15px; text-align: left; border-bottom: 1px solid #ddd; }")
	fmt.Fprintf(w, "th { background-color: #4CAF50; color: white; }")
	fmt.Fprintf(w, "tr:hover { background-color: #f5f5f5; }")
	fmt.Fprintf(w, ".header-name { font-weight: bold; color: #333; }")
	fmt.Fprintf(w, ".header-value { color: #666; word-break: break-all; }")
	fmt.Fprintf(w, "</style></head><body>")

	fmt.Fprintf(w, "<h1>HTTP Request Headers</h1>")
	fmt.Fprintf(w, "<p>Request Method: <strong>%s</strong></p>", r.Method)
	fmt.Fprintf(w, "<p>Request URL: <strong>%s</strong></p>", r.URL.String())
	fmt.Fprintf(w, "<p>Remote Address: <strong>%s</strong></p>", r.RemoteAddr)

	fmt.Fprintf(w, "<table>")
	fmt.Fprintf(w, "<tr><th>Header Name</th><th>Header Value</th></tr>")

	// Sort headers alphabetically for consistent display
	var headerNames []string
	for name := range r.Header {
		headerNames = append(headerNames, name)
	}
	sort.Strings(headerNames)

	for _, name := range headerNames {
		values := r.Header[name]
		for _, value := range values {
			fmt.Fprintf(w, "<tr><td class=\"header-name\">%s</td><td class=\"header-value\">%s</td></tr>", name, value)
		}
	}

	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "</body></html>")
}
