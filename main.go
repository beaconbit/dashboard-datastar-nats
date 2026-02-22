package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/starfederation/datastar-go/datastar"
)

type DashboardData struct {
	Title string
}

type QuarterAStore struct {
	sync.RWMutex
	Data map[string]interface{}
}

type UpdateMessage struct {
	Path  string `json:"path"`
	Value int    `json:"value"`
}

var (
	store = &QuarterAStore{
		Data: map[string]interface{}{
			"screen": "tv",
			"counters": map[string]int{
				"q1": 0,
				"q2": 0,
				"q3": 0,
				"q4": 0,
			},
			"quarterA": map[string]interface{}{
				"col1": map[string]interface{}{
					"title": "CBW 1",
					"num1":  45, "num2": 23, "num3": 78, "num4": 12,
					"num5": 89, "num6": 34, "num7": 67, "num8": 91,
				},
				"col2": map[string]interface{}{
					"title": "CBW 2",
					"num1":  56, "num2": 33, "num3": 81, "num4": 19,
					"num5": 72, "num6": 44, "num7": 95, "num8": 28,
				},
				"col3": map[string]interface{}{
					"title": "CBW 3",
					"num1":  37, "num2": 64, "num3": 22, "num4": 88,
					"num5": 15, "num6": 53, "num7": 76, "num8": 41,
				},
			},
			"quarterB": map[string]interface{}{
				"col1": map[string]interface{}{
					"title": "RCW 1",
					"num1":  82, "num2": 17, "num3": 93, "num4": 41,
					"num5": 26, "num6": 58, "num7": 39, "num8": 74,
				},
				"col2": map[string]interface{}{
					"title": "RCW 2",
					"num1":  65, "num2": 29, "num3": 84, "num4": 36,
					"num5": 71, "num6": 13, "num7": 92, "num8": 47,
				},
				"col3": map[string]interface{}{
					"title": "RCW 3",
					"num1":  18, "num2": 55, "num3": 27, "num4": 63,
					"num5": 89, "num6": 32, "num7": 76, "num8": 44,
				},
			},
			"quarterC": map[string]interface{}{
				"cbw1": map[string]interface{}{
					"waitTime": 5,
				},
				"cbw2": map[string]interface{}{
					"waitTime": 3,
				},
				"cbw3": map[string]interface{}{
					"waitTime": 7,
				},
				"wastedMinutes": map[string]interface{}{
					"hour4":   120,
					"hour3":   95,
					"hour2":   80,
					"hour1":   65,
					"current": 150,
				},
			},
			"quarterD": map[string]interface{}{
				"col1": map[string]interface{}{
					"title":   "Quota₁",
					"current": 34,
					"target":  100,
				},
				"col2": map[string]interface{}{
					"title":   "Quota₂",
					"current": 89,
					"target":  100,
				},
				"col3": map[string]interface{}{
					"title":   "Quota₃",
					"current": 50,
					"target":  100,
				},
				"col4": map[string]interface{}{
					"title":   "Quota₄",
					"current": 45,
					"target":  100,
				},
				"col5": map[string]interface{}{
					"title":   "Quota₅",
					"current": 67,
					"target":  100,
				},
				"col6": map[string]interface{}{
					"title":   "Quota₆",
					"current": 23,
					"target":  100,
				},
			},
		},
	}
	// SSE broadcast
	sseClients   = make(map[*datastar.ServerSentEventGenerator]bool)
	sseClientsMu sync.RWMutex
)

var nc *nats.Conn

func main() {
	var err error
	// Connect to NATS
	nc, err = nats.Connect("nats://nats:4222")
	if err != nil {
		log.Printf("Failed to connect to NATS: %v", err)
	} else {
		defer nc.Close()
		log.Println("Connected to NATS")

			// Subscribe to all Quarter A number topics
		subscribeToQuarterATopics(nc)

			// Subscribe to Quarter D topics

		//subscribeToQuarterDTopics(nc)

			// Subscribe to Quarter C topics

		//subscribeToQuarterCTopics(nc)
	}

	// Start periodic updates for Quarter D (fallback if no NATS messages)
	// Disabled as we now use NATS topics for Quarter D updates
	// startQuarterDUpdates()

	// HTTP routes
	// Serve static files from /static path
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/stream", streamHandler)

	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/sse", sseHandler)
	http.HandleFunc("/api/store", storeHandler)
	http.HandleFunc("/sse-test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("templates", "sse-test.html"))
	})
	http.HandleFunc("/api/debug", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Debug: %s %s", r.Method, r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	log.Println("Server starting on :3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func subscribeToQuarterATopics(nc *nats.Conn) {
	// Subscribe to topics for each Quarter A number (24 topics)
	for col := 1; col <= 3; col++ {
		for num := 1; num <= 8; num++ {
			topic := fmt.Sprintf("quarterA.col%d.num%d", col, num)
			_, err := nc.Subscribe(topic, func(msg *nats.Msg) {
				var value int
				if err := json.Unmarshal(msg.Data, &value); err != nil {
					log.Printf("Error unmarshaling value from %s: %v", msg.Subject, err)
					return
				}

				// Update store
				store.Lock()
				quarterA := store.Data["quarterA"].(map[string]interface{})
				colKey := fmt.Sprintf("col%d", col)
				colData := quarterA[colKey].(map[string]interface{})
				numKey := fmt.Sprintf("num%d", num)
				colData[numKey] = value
				store.Unlock()

				// Broadcast update to SSE clients
				path := fmt.Sprintf("quarterA.%s.%s", colKey, numKey)
				broadcastUpdate(UpdateMessage{Path: path, Value: value})
				//log.Printf("Broadcast Updated %s to %d", msg.Subject, value)
			})

			if err != nil {
				log.Printf("Failed to subscribe to %s: %v", topic, err)
			} else {
				log.Printf("Subscribed to %s", topic)
			}
		}
	}
}

func subscribeToQuarterDTopics(nc *nats.Conn) {
	// Subscribe to topics for each Quarter D column (6 topics)
	for col := 1; col <= 6; col++ {
		topic := fmt.Sprintf("quarterD.col%d", col)
		_, err := nc.Subscribe(topic, func(msg *nats.Msg) {
			var value int
			if err := json.Unmarshal(msg.Data, &value); err != nil {
				log.Printf("Error unmarshaling value from %s: %v", msg.Subject, err)
				return
			}

			// Ensure value is between 1 and 99
			if value < 1 {
				value = 1
			}
			if value > 99 {
				value = 99
			}

			// Update store
			store.Lock()
			quarterD := store.Data["quarterD"].(map[string]interface{})
			colKey := fmt.Sprintf("col%d", col)
			colData := quarterD[colKey].(map[string]interface{})
			colData["current"] = value
			store.Unlock()

			// Broadcast update to SSE clients
			path := fmt.Sprintf("quarterD.%s.current", colKey)
			broadcastUpdate(UpdateMessage{Path: path, Value: value})
			log.Printf("QuarterD updated %s to %d", msg.Subject, value)
		})

		if err != nil {
			log.Printf("Failed to subscribe to %s: %v", topic, err)
		} else {
			log.Printf("Subscribed to %s", topic)
		}
	}
}

func subscribeToQuarterCTopics(nc *nats.Conn) {
	log.Printf("Setting up Quarter C subscriptions")
	// Subscribe to topics for Quarter C wait times (3 topics)
	for cbw := 1; cbw <= 3; cbw++ {
		topic := fmt.Sprintf("quarterC.cbw%d.waitTime", cbw)
		_, err := nc.Subscribe(topic, func(msg *nats.Msg) {
			var value int
			if err := json.Unmarshal(msg.Data, &value); err != nil {
				log.Printf("Error unmarshaling value from %s: %v", msg.Subject, err)
				return
			}
			// Update store
			store.Lock()
			quarterC := store.Data["quarterC"].(map[string]interface{})
			cbwKey := fmt.Sprintf("cbw%d", cbw)
			cbwData := quarterC[cbwKey].(map[string]interface{})
			cbwData["waitTime"] = value
			store.Unlock()
			// Broadcast update
			path := fmt.Sprintf("quarterC.%s.waitTime", cbwKey)
			broadcastUpdate(UpdateMessage{Path: path, Value: value})
			log.Printf("Updated %s to %d", msg.Subject, value)
		})
		if err != nil {
			log.Printf("Failed to subscribe to %s: %v", topic, err)
		} else {
			log.Printf("Subscribed to %s", topic)
		}
	}
	// Subscribe to topics for wasted minutes (5 topics)
	wastedKeys := []string{"hour4", "hour3", "hour2", "hour1", "current"}
	for _, key := range wastedKeys {
		topic := fmt.Sprintf("quarterC.wastedMinutes.%s", key)
		_, err := nc.Subscribe(topic, func(msg *nats.Msg) {
			var value int
			if err := json.Unmarshal(msg.Data, &value); err != nil {
				log.Printf("Error unmarshaling value from %s: %v", msg.Subject, err)
				return
			}
			// Update store
			store.Lock()
			quarterC := store.Data["quarterC"].(map[string]interface{})
			wastedMinutes := quarterC["wastedMinutes"].(map[string]interface{})
			wastedMinutes[key] = value
			store.Unlock()
			// Broadcast update
			path := fmt.Sprintf("quarterC.wastedMinutes.%s", key)
			broadcastUpdate(UpdateMessage{Path: path, Value: value})
			log.Printf("Updated %s to %d", msg.Subject, value)
		})
		if err != nil {
			log.Printf("Failed to subscribe to %s: %v", topic, err)
		} else {
			log.Printf("Subscribed to %s", topic)
		}
	}
}

func startQuarterDUpdates() {
	go func() {
		log.Println("Starting Quarter D periodic updates (every 30 seconds)")
		// Seed random
		rand.Seed(time.Now().UnixNano())
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// Pick random column 1-6
			col := rand.Intn(6) + 1
			colKey := fmt.Sprintf("col%d", col)

			store.Lock()
			quarterD := store.Data["quarterD"].(map[string]interface{})
			colData := quarterD[colKey].(map[string]interface{})
			// target could be int or float64
			var target int
			switch v := colData["target"].(type) {
			case int:
				target = v
			case float64:
				target = int(v)
			default:
				target = 150
			}
			// Generate new current between 0 and target
			newCurrent := rand.Intn(target + 1)
			colData["current"] = newCurrent
			store.Unlock()

			// Broadcast update
			path := fmt.Sprintf("quarterD.%s.current", colKey)
			broadcastUpdate(UpdateMessage{Path: path, Value: newCurrent})
			log.Printf("QuarterD updated %s to %d", path, newCurrent)
		}
	}()
}

func broadcastUpdate(update UpdateMessage) {
	sseClientsMu.RLock()
	defer sseClientsMu.RUnlock()
	clientCount := len(sseClients)
	if false {
	    log.Printf("Broadcasting update %s=%d to %d clients", update.Path, update.Value, clientCount)
	}

	// Convert to Datastar patch signals format
	signals := map[string]any{update.Path: update.Value}

	for sse := range sseClients {
		// Patch signals using Datastar SDK
		err := sse.MarshalAndPatchSignals(signals)
		if err != nil {
			log.Printf("Error patching signals: %v", err)
			// Connection may be closed, we'll clean up in sseHandler defer
		}
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "dashboard.html")))
	data := DashboardData{
		Title: "Dashboard",
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func streamHandler(w http.ResponseWriter, r *http.Request) {

    log.Printf("streamHandler activated")
    sse := datastar.NewSSE(w, r)

    // Buffer of 64 so slow clients don't block the NATS subscription
    msgCh := make(chan *nats.Msg, 64)

    col := 1
    num := 1
    topic := fmt.Sprintf("quarterA.col%d.num%d", col, num)
    sub, err := nc.ChanSubscribe(topic, msgCh)
    if err != nil {
        http.Error(w, "Failed to subscribe to NATS", http.StatusInternalServerError)
        return
    }
    defer sub.Unsubscribe()

    // Keep the connection open until the client disconnects
    ctx := r.Context()
    for {
        select {
        case <-ctx.Done():
            return
        case msg, ok := <-msgCh:
            if !ok {
		log.Printf("Not good msg received, breaking")
                return
            }

            // Expect the NATS message payload to be an int
            var value int
            if err := json.Unmarshal(msg.Data, &value); err != nil {
                log.Printf("Failed to unmarshal NATS message on subject %s: %v", msg.Subject, err)
                continue
            }
	    log.Printf("Stream NATS message was: %s", msg.Data)
	    log.Printf("Stream NATS value was: %s", value)

            if err := sse.MarshalAndPatchSignals(map[string]any{
		"quarterA": map[string]any{
		    "col1": map[string]any{
			"num1": value,
		    },
		},
	    }); err != nil {
                log.Printf("Failed to send SSE: %v", err)
                return
            }
	    if f, ok := w.(http.Flusher); ok {
                log.Printf("Flush")
		f.Flush()
	    }
        }
    }
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("SSE connection established from %s", r.RemoteAddr)

	// Create Datastar SSE generator
	sse := datastar.NewSSE(w, r)
	sseClientsMu.Lock()
	sseClients[sse] = true
	sseClientsMu.Unlock()

	// Clean up when connection closes
	defer func() {
		log.Printf("SSE connection closed from %s", r.RemoteAddr)
		sseClientsMu.Lock()
		delete(sseClients, sse)
		sseClientsMu.Unlock()
	}()

	// Send initial store as init event (nested format for backward compatibility)
	store.RLock()
	initialData, _ := json.Marshal(store.Data)
	store.RUnlock()

	err := sse.Send(datastar.EventType("init"), []string{string(initialData)})
	if err != nil {
		log.Printf("Error sending initial data: %v", err)
		return
	}

	// Listen for heartbeats and connection closure
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send heartbeat comment to keep connection alive
			// Use the underlying writer directly
			if _, err := w.Write([]byte(": heartbeat\n\n")); err != nil {
				log.Printf("Heartbeat write error: %v", err)
				return
			}
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		case <-r.Context().Done():
			return
		}
	}
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	store.RLock()
	defer store.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.Data)
}
