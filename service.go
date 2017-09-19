package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	err := http.ListenAndServe(string(config.Port), nil)
	if err != nil {
		panic(err)
	}
}

// Service the /get endpoint
func get(res http.ResponseWriter, req *http.Request) {
	// single client from the start
	cache.Lock()
	defer cache.Unlock()

	// get the key TODO add defensive code for invalid path
	key := req.URL.Query().Get("key")
	if key == "" {
		log.Println("key is empty\n")
		fmt.Print(res, "key is empty", 400)
		return
	}

	// now Get from the managedCache
	log.Printf("Fetching for key=%s\n", key)
	var value = cache.Get(key)

	if value == nil {
		log.Printf("cache does not contain key=%s\n", key)

		if backingRedis == nil {
			log.Println("No backing Redis")
			fmt.Print(res, "Backing Redis not connected", 500)
			return
		}

		// Get from backing Redis
		value, err := backingRedis.Get(key).Result()
		if err == redis.Nil {
			log.Printf("Backing Redis has no key=%s\n", key)
			fmt.Print(res, "Not Found", 404)
			return
		} else if err != nil {
			panic(err)
		} else {
			log.Printf("Found value from backing Redis: key=%s returns %s\n", key, value)
			cache.Add(key, value)
			// return as a string
			fmt.Print(res, value, 200)
			return
		}
	} else {
		log.Printf("Found value from cache: key=%s returns %s\n", key, *value)
		fmt.Print(res, *value, 200)
		return
	}
}

// Not supported endpoints
func notSupportURL(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	log.Printf("(400) %s - Not Supported.\n", path)
	fmt.Fprint(res, "Not Supported", 400)
}

// Listens to OS signal for a shutdown message
func listenForShutdown(ch <-chan os.Signal) {
	log.Println("Registering for OS SHUTDOWN signal...")
	// Listen for signal from OS
	<- ch

	// Got a signal. Time to go
	log.Println("Got a ShutDown signal from OS. Going down NOW...")

	// Send signal to no longer accept connection request and shutdown after all current requests are completed.
	//manners.Close()
	os.Exit(1)
}