# Go-Redis

A toy redis clone in go. Currently only supports a small palette of commands with limited functionality:
- PING (word)
- SET (key) (value)
- GET  (key)
- DEL  (key1) (key2) ...
- COPY (sourceKey) (destinationKey)
- APPEND (key) (addedValue)
- EXISTS (key1) (key2) ...
- HSET (hash) (key) (value)
- HGET (hash) (key)
- HGETALL (hash)
- HDEL (key) (field1) (field2) ...

The current goal is to explore implementing other commands.
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
