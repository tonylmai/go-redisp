package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"github.com/braintree/manners"
	"github.com/go-redis/redis"
)

// Application configuration
type Config struct {
	BackingRedisUrl string
	Capacity        int64
	Expiry          int64
	Port            string
}

// Managed cache
var cache *managedCache

// Backing Redis
var backingRedis *redis.Client

// Set up logger and register for OS shutdown signal
func init() {
	// Log to stdout
	log.SetOutput(os.Stdout)

	// Initialize a channel to listen for signal from the OS
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)

	// Kick off listening for signal from OS asynchronously
	go listenForShutdown(ch)


}

// Start REST server and redis client
func Start(config Config) {
	log.Printf("Starting RedisP with\n\tUrl: %s\n\tCapacity: %d\n\tExpiry: %d\n\tPort:%s", config.BackingRedisUrl, config.Capacity, config.Expiry, config.Port)

	cache = NewManagedCache(config.Capacity, config.Expiry)

	log.Printf("Connecting to backing Redis at %s\n", config.BackingRedisUrl)
	backingRedis = redis.NewClient(&redis.Options{
		Addr:     config.BackingRedisUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Register all handlers for all path now
	http.HandleFunc("/get", get)
	// Fall through. If request path does not match anything above this, it will come to this last one
	http.HandleFunc("/", notSupportURL)

	// Ready
	log.Printf("RedisP is now ready to accept REST request on port %s\n", config.Port)

	// launch web services now
	err := manners.ListenAndServe(string(config.Port), nil)
	if err != nil {
		panic(err)
	}
}

// Service the /get endpoint
func get(res http.ResponseWriter, req *http.Request) {
	log.Printf("Got a request \n")

	// single client from the start
	cache.Lock()
	defer cache.Unlock()

	// get the key TODO add defensive code for invalid path
	path := req.URL.Path
	parts := strings.Split(path, "/")
	key := parts[2]
	if key == "" {
		fmt.Print(res, "key is empty", 400)
		return
	}

	// now Get from the managedCache
	var value = cache.Get(key)

	if value == nil {
		// Get from backing Redis
		value, err := backingRedis.Get(key).Result()
		if err == redis.Nil {
			fmt.Print(res, "Not Found", 404)
		} else if err != nil {
			panic(err)
		} else {
			cache.Add(key, value)
			// return as a string
			fmt.Print(res, value, 200)
		}
	} else {
		fmt.Print(res, value, 200)
	}
	return
}

// Not supported endpoints
func notSupportURL(res http.ResponseWriter, req *http.Request) {
	log.Println("Not a valid path")
	fmt.Print(res, "Invalid path")
}

// Listens to OS signal for a shutdown message
func listenForShutdown(ch <-chan os.Signal) {
	log.Println("Registering for OS SHUTDOWN signal...")
	// Listen for signal from OS
	<- ch

	// Got a signal. Time to go
	log.Println("Got a ShutDown signal from OS. Going down NOW...")

	// Send signal to no longer accept connection request and shutdown after all current requests are completed.
	manners.Close()
}