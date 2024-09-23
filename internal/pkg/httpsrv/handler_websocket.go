package httpsrv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"goapp/internal/pkg/watcher"

	"github.com/gorilla/websocket"
)

func (s *Server) handlerWebSocket(w http.ResponseWriter, r *http.Request) {
	// Check the CSRF token only if the request is non-CLI
	isCliReq := r.URL.Query().Get("cli")
	if isCliReq != "true" {
		// Get session based on the received request
		session, err := s.cookieStore.Get(r, "ws-session")
		if err != nil {
			http.Error(w, "Unable to retrieve session", http.StatusInternalServerError)
			return
		}

		// Get the CSRF token from the session
		csrfTokenSession, ok := session.Values["csfr_token"].(string)
		if !ok || csrfTokenSession == "" {
			http.Error(w, "Missing CSRF token", http.StatusForbidden)
			return
		}

		// Get the CSRF token from the query parameter
		csrfTokenQueryParam := r.URL.Query().Get("csrf_token")

		// Compare the tokens
		if csrfTokenSession != csrfTokenQueryParam {
			http.Error(w, "Invalid CSRF token", http.StatusForbidden)
			return
		}
	}

	// Create and start a watcher.
	var watch = s.watcherPool.Get().(*watcher.Watcher)

	s.addWatcher(watch)
	defer s.removeWatcher(watch)

	// Start WS.
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.error(w, http.StatusInternalServerError, fmt.Errorf("failed to upgrade connection: %w", err))
		return
	}
	defer func() { _ = c.Close() }()

	log.Printf("websocket started for watcher %s\n", watch.GetWatcherId())
	defer func() {
		log.Printf("websocket stopped for watcher %s\n", watch.GetWatcherId())
	}()

	// Read done.
	readDoneCh := make(chan struct{})

	// All done.
	doneCh := make(chan struct{})
	defer close(doneCh)

	go func() {
		defer close(readDoneCh)
		for {
			select {
			default:
				_, p, err := c.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
						log.Printf("failed to read message: %v\n", err)
					}
					return
				}
				var m watcher.CounterReset
				if err := json.Unmarshal(p, &m); err != nil {
					log.Printf("failed to unmarshal message: %v\n", err)
					continue
				}
				watch.ResetCounter()
			case <-doneCh:
				return
			case <-s.quitChannel:
				return
			}
		}
	}()

	for {
		select {
		case cv := <-watch.Recv():
			data, _ := json.Marshal(cv)
			err = c.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("failed to write message: %v\n", err)
				}
				return
			}
		case <-readDoneCh:
			return
		case <-s.quitChannel:
			return
		}
	}
}
