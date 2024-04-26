package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/graphql-go/graphql"
)

func formatErrorMessage(message string) []byte {
	return []byte(fmt.Sprintf(`{"error": "%s"}`, message))
}

func executeQuery(query string, variableValues map[string]interface{}, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
		VariableValues: variableValues,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

type QueryBody struct {
	Query          string                 `json:"query"`
	VariableValues map[string]interface{} `json:"variables"`
}

type Request struct {
	Headers map[string]interface{} `json:"headers"`
	Body    string                 `json:"body"`
	Date    string                 `json:"date"`
}

var (
	requests    []Request
	requestsMtx sync.RWMutex
)

func recordRequest(r *http.Request, body []byte) {
	// Log the request for debugging purposes
	requestsMtx.Lock()
	defer requestsMtx.Unlock()

	headers := make(map[string]interface{})
	for name, value := range r.Header {
		headers[name] = strings.Join(value, ",")
	}
	request := Request{
		Headers: headers,
		Body:    string(body),
		Date:    time.Now().Format(time.RFC850),
	}
	requests = append(requests, request)
}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write(formatErrorMessage(fmt.Sprintf("Failed to read GraphQL query: %s", err)))
		if err != nil {
			log.Fatalf("Failed to write error message to response: %s", err)
		}
	}
	defer r.Body.Close()

	queryBody := QueryBody{}
	err = json.Unmarshal(body, &queryBody)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		message := fmt.Sprintf("Failed to read query body: %s", err)
		_, _ = w.Write(formatErrorMessage(message))
		return
	}

	var query string
	var variableValues map[string]interface{}
	if queryBody.Query != "" {
		query = queryBody.Query
		variableValues = queryBody.VariableValues
	} else {
		query = string(body)
	}

	result := executeQuery(query, variableValues, schema)

	recordRequest(r, body)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Fatalf("Failed to write query response: %s", err)
	}
}

func debugHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(formatErrorMessage("invalid method"))
		return
	}
	requestsMtx.RLock()
	defer requestsMtx.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(requests)
}

func init() {
	// Initialize the Requests slice
	requests = []Request{}
}

func main() {
	help := flag.Bool("help", false, "usage")
	address := flag.String("address", "localhost", "address to listen")
	port := flag.Int("port", 8080, "port to bind")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	http.HandleFunc("/graphql", graphqlHandler)
	http.HandleFunc("/graphql/", graphqlHandler)
	http.HandleFunc("/debug/requests", debugHandler)
	http.HandleFunc("/debug/requests/", debugHandler)

	addr := net.JoinHostPort(*address, strconv.Itoa(*port))
	fmt.Printf("Server is running on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
