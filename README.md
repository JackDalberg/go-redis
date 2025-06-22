# Go-Redis

A toy redis clone in go. Currently only supports a small palette of commands with limited functionality:
- PING (word)
- SET (key) (value)
- GET  (key)
- DEL  (key1) (key2) ...
- COPY (sourceKey) (destinationKey)
- APPEND (key) (addedValue)
- EXISTS (key1) (key2) ...
- INCR (key)
- HSET (hash) (key) (value)
- HGET (hash) (key)
- HGETALL (hash)
- HDEL (hash) (key1) (key2) ...
- RPUSH (list) (value1) (value2)
- LPUSH (list) (value1) (value2)
- LLEN (list)
- LSET (list) (index) (value)
- LRANGE (list) (start) (end)

The current goal is to explore implementing other commands. Also worth noting that this is not a "true" Redis server. You can assign multiple keys different values, but only if the underlying type is different. For example, you can do `SET example val` and `HSET example key val` with no problem as each are stored in separated maps `SETs` and `HSETs`. This will be something to remedy later on with better type validation to make truly unique key value pairs, regardless of the type.
## How To Run

In order to be able to connect, you need a Redis client.  The simplest is the Redis CLI from [official Redis](https://redis.io/docs/latest/operate/oss_and_stack/install/install-stack/).  

On WSL2, I accomplished this by
```
$ sudo apt install redis
```
Then disabling the Redis server so we can connect to go-redis by
```
$ sudo systemctl stop redis
```
Finally, we can run the go-redis server by
```
$ go run . 
```
And connect with
```
$ redis-cli
```
