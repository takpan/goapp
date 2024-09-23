package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"goapp/internal/pkg/watcher"
	"log"
	"net"
	"net/url"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

func main() {
	// Flags
	n := flag.Int("n", 0, "Specify the number of parallel connections")
	host := flag.String("host", "localhost", "(optional) Specify the host to connect to")
	port := flag.String("port", "8080", "(optional) Specify the port to connect to")
	path := flag.String("path", "/goapp/ws", "(optional) Websocket path")

	flag.Parse()

	// Check wether a valid value provided for n
	if *n <= 0 {
		log.Println("Error: The -n flag has been ommited or an invalid value provided")
		flag.Usage() // Prints usage information
		os.Exit(1)   // Exit with non-zero status to indicate an error
	}

	// Prepare URL
	hostPort := net.JoinHostPort(*host, *port)
	wsUrl := url.URL{
		Scheme: "ws",
		Host:   hostPort,
		Path:   *path,
	}

	// Add a CLI query parameter to identify the origin of the request in the websocket handler
	query := wsUrl.Query()
	query.Set("cli", "true")
	wsUrl.RawQuery = query.Encode()

	var wg sync.WaitGroup

	// Start n concurrent websocket connections
	for i := 0; i < *n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			conn, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), nil)
			if err != nil {
				log.Printf("Failed to connect to thewebsocket server, connection id: %d, error: %v\n", id, err)
				return
			}
			defer conn.Close()

			// Read and handle messages received by the server on the websocket connection
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Printf("Error reading message: %v\n", err)
					break
				}

				var receivedMsg watcher.Counter
				if err := json.Unmarshal(msg, &receivedMsg); err != nil {
					log.Printf("Error unmarshalling message: %v\n", err)
					continue
				}

				// Print message in a specific format (e.g.: [conn #0] iteration: 1, value: 66D53ED788)
				fmt.Printf("[conn #%d] iteration: %d, value: %s\n", id, receivedMsg.Iteration, receivedMsg.Value)
			}

		}(i)
	}

	// Using a wait group to indefinitely prevent the function from terminating until a stop signal is triggered by the user
	wg.Wait()
}
