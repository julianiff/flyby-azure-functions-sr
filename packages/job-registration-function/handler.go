package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-queue-go/azqueue"
)

const (
	// QueueMaxMessagesDequeue indicates the maximum number of messages
	// you can retrieve with each call to Dequeue (32).
	QueueMaxMessagesDequeue = 32

	// QueueMaxMessagesPeek indicates the maximum number of messages
	// you can retrieve with each call to Peek (32).
	QueueMaxMessagesPeek = 32

	// QueueMessageMaxBytes indicates the maximum number of bytes allowed for a message's UTF-8 text.
	QueueMessageMaxBytes = 64 * 1024 // 64KB
)

const SASTimeFormat = "2006-01-02T15:04:05Z" //"2017-07-27T00:00:00Z" // ISO 8601
const SASVersion = ServiceVersion
const (
	// ServiceVersion specifies the version of the operations used in this package.
	ServiceVersion = "2018-03-28"
)

// Please set the ACCOUNT_NAME and ACCOUNT_KEY environment variables to your storage account's
// name and account key, before running the examples.
func accountInfo() (string, string) {
	return os.Getenv("devstoreaccount1"), os.Getenv("Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	message := "Need post to add new  item to the queue..\n"
	name := r.URL.Query().Get("name")
	if name != "" {
		message = fmt.Sprintf("Need post to add new item to the quesfdue  astest.\n", name)
	}
	fmt.Fprint(w, message)

	// Here comes pipeline queue stuff
	doQueueStuff()

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

func doQueueStuff() {
	// From the Azure portal, get your Storage account's name and account key.

	storageAccountName := "devstoreaccount1"
	storageAccountKey := "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="

	credential, err := azqueue.NewSharedKeyCredential(storageAccountName, storageAccountKey)
	if err != nil {
		log.Fatal("Error creating credentials: ", err)
	}

	p := azqueue.NewPipeline(credential, azqueue.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:10001/devstoreaccount1"))
	serviceURL := azqueue.NewServiceURL(*u, p)
	ctx := context.TODO()
	queueURL := serviceURL.NewQueueURL("funqueue")

	props, err := queueURL.GetProperties(ctx)

	if err != nil {
		// https://godoc.org/github.com/Azure/azure-storage-queue-go/azqueue#StorageErrorCodeType
		errorType := err.(azqueue.StorageError).ServiceCode()

		if errorType == azqueue.ServiceCodeQueueNotFound {

			log.Print("Queue does not exist, creating")

			_, err = queueURL.Create(ctx, azqueue.Metadata{})
			if err != nil {
				log.Fatal("Error creating queue: ", err)
			}

			props, err = queueURL.GetProperties(ctx)
			if err != nil {
				log.Fatal("Error parsing url: ", err)
			}

		} else {
			log.Fatal("Error getting queue properties: ", err)
		}
	}

	messageCount := props.ApproximateMessagesCount()
	log.Printf("Appx number of messages: %d", messageCount)
}
