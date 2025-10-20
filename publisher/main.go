package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Msg string `json:"msg"`
}

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")

	if rabbitURL == "" {
		rabbitURL = "amqp://user:pass@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v - rabbitURL (%v)", err, rabbitURL)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("tasks", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var m Message
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if m.Msg == "" {
			http.Error(w, "missing msg in body", http.StatusBadRequest)
			return
		}

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(m.Msg),
			},
		)
		if err != nil {
			http.Error(w, "Failed to publish", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"sent":"%s"}`, m.Msg)
	})

	log.Println("Publisher running on :8080")
	http.ListenAndServe(":8080", nil)
}
