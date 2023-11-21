package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

// Object represents the structure of your data.
type Object struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func main() {
	// Initialize a Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Check if the connection to Redis is successful
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	fmt.Println("connection has been established successfully")

	// Create a new router
	router := mux.NewRouter()

	// Define REST API routes
	router.HandleFunc("/publish", func(w http.ResponseWriter, _ *http.Request) {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go startPublisher(wg, rdb)
		wg.Wait()
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	router.HandleFunc("/subscribe", func(w http.ResponseWriter, _ *http.Request) {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go startSubscriber(wg, rdb)
		wg.Wait()
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	// Start the REST API server in a goroutine
	go func() {
		log.Println("Starting server on :8080...")
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	// Use a wait group to wait for goroutines to finish
	var wg sync.WaitGroup

	// Start the publisher goroutine
	wg.Add(1)
	go startPublisher(&wg, rdb)
	fmt.Println("Objects have been published to redis")

	// Start the subscriber goroutine
	wg.Add(1)
	go startSubscriber(&wg, rdb)
	fmt.Println("Objects have been consumed by the subscriber from redis")

	// Wait for goroutines to finish
	wg.Wait()

	fmt.Println("Publisher and subscriber finished.")
}

func startPublisher(wg *sync.WaitGroup, rdb *redis.Client) {
	defer wg.Done()

	// Generate and store 100000 objects in the Redis queue
	for i := 1; i <= 100000; i++ {
		object := generateObject(i)
		err := enqueueObject(rdb, object)
		if err != nil {
			log.Printf("Error enqueuing object: %v", err)
		}
	}

	// Signal that no more objects will be enqueued
	closeObjectQueue(rdb)
}

func startSubscriber(wg *sync.WaitGroup, rdb *redis.Client) {
	defer wg.Done()

	// Subscribe to the "object_queue" channel
	pubsub := rdb.Subscribe("object_queue")
	defer pubsub.Close()

	for {
		// Wait for a message to be received
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			if err == redis.Nil {
				// No more messages in the channel
				break
			}
			log.Printf("Error receiving message: %v", err)
			return
		}

		// If the received message is "CLOSED", exit the goroutine
		if msg.Payload == "CLOSED" {
			fmt.Println("Subscriber received 'CLOSED', exiting.")
			return
		}

		// Unmarshal the JSON message into an Object
		var obj Object
		err = json.Unmarshal([]byte(msg.Payload), &obj)
		if err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			return
		}

		// Process the received object
		processObject(obj)

		// Log the received message
		log.Printf("Subscriber received message: %+v", obj)
	}
}

func processObject(obj Object) {
	// Simulate processing the object (replace this with your actual processing logic)
	fmt.Printf("Processing object: %+v\n", obj)
	time.Sleep(time.Millisecond) // Simulate processing time
}

func closeObjectQueue(rdb *redis.Client) {
	// Signal that no more objects will be enqueued
	err := rdb.Publish("object_queue", "CLOSED").Err()
	if err != nil {
		log.Printf("Error publishing 'CLOSED' message: %v", err)
	}
}

func generateObject(id int) Object {
	return Object{
		ID:    id,
		Name:  fmt.Sprintf("Object%d", id),
		Value: rand.Intn(100),
	}
}

func enqueueObject(rdb *redis.Client, obj Object) error {
	// Convert the object to JSON
	objJSON, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("error marshalling object to JSON: %v", err)
	}

	// Enqueue the JSON object in the Redis queue
	err = rdb.LPush("object_queue", string(objJSON)).Err()
	if err != nil {
		return fmt.Errorf("error enqueuing object in Redis: %v", err)
	}

	return nil
}
