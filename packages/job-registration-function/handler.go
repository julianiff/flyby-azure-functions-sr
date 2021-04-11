package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	message := "Need post to add new item to the queue..\n"
	name := r.URL.Query().Get("name")
	if name != "" {
		message = fmt.Sprintf("Need post to add new item to the queue.\n", name)
	}
	fmt.Fprint(w, message)
}

// connect to Queue service and put new item on the queue

func main() {
	listenAddr := "localhost:8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	http.HandleFunc("/api/job-registration", helloHandler)
	log.Printf("About to listen on %s. Go to %s", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
