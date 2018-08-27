go-redisp is a Redis proxy service. This proxy implements a subset of the Redis protocol, with the ability to add
additional features on top of Redis (e.g. caching and sharding).

**Runtime**

Build go-redisp by execute the following command

`$ go build`

Launch the backing Redis by issuing this command

```$ redis-server --daemonize yes```

Launch Redis Proxy by issuing this command

```$ go-redisp ```

This will start the web server listening on port 9997 (configurable using the conf.yaml file).

Once the server is up and running, send a `curl http://localhost:9997/get?key=abc` command to fetch for a stored value. If a value is not found, a 404 will return to the client.

*Launch Redis with Docker*

go-redisp is configured with a Dockerfile and docker-compose.yml where you may build a Docker image and launch the service in a swamp.

To build the Docker image, 

`$ go build`

`$ docker build -t go-redisp .`

To run with Docker,

`docker run -p 9997:9997 -ti -v /tmp:tmp go-redisp /bin/bash`

**High-level architecture**

go-redisp will be implemented as a REST API service. At the start, there will be a single endpoint '/get?key=xyz'.

As for Redis client, I chose `github.com/go-redis/redis` for its simplicity. 

In order to enforce a single client only, I will use a Mutex lock from the entry. Only one client may obtain the lock at a time.

For LRU eviction, instead of using a sorted map (by int), I opted for using a regular map and will perform a sort when I need to drop an entry. I am assuming the turn over is small hence I trade it for not sorting the map upon every get.

**Requirements**

*Single backing instance*

Each instance of the proxy service is associated with a single Redis service instance  called the “backing Redis” . 

*Cached GET*

A GET command, directed at the proxy, returns the value of the specified key from the proxy’s local cache if the local cache contains a value for that key. If the local cache does not contain a value for the specified key, it fetches the value from the backing Redis instance and stores it in the local cache, associated with the specified key.

*Global expiry*

Entries added to the proxy cache are expired after being in the cache for a time duration that is globally configured  per instance . After an entry is expired, a GET command will act as if the value associated with the key was never stored in the cache.

*LRU eviction*

One the cache fills to capacity, the least recently used key is evicted each time a new key needs to be added to the cache.

*Fixed key size*

The cache capacity is configured in terms of number of keys it retains.

*Single client*

The proxy is able to process at least one concurrent client request  i.e. when more than one client makes a request to the proxy at the same time, the second client’s request only starts processing once the first one has completed .

*Configuration*
