package main

import (
	"strconv"
	"sync"
)

// implementation from https://redis.io/docs/latest/commands
var Handlers = map[string]func([]Value) Value{
	"PING":   ping,
	"SET":    set,
	"GET":    get,
	"EXISTS": exists,
	"APPEND": _append, // bc "append" in std
	"DEL":    del,
	"COPY":   copy,
	"INCR":   incr,
	// Hashmap commands
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
	"HDEL":    hdel,
	// List commands (L)eft-CMD and (R)ight-CMD
	"RPUSH":  rpush,
	"LPUSH":  lpush,
	"LLEN":   llen,
	"LSET":   lset,
	"LRANGE": lrange,
}

var ModifiesDB = []string{
	"SET",
	"DEL",
	"APPEND",
	"COPY",
	"INCR",
	// Hashmap commands
	"HSET",
	"HDEL",
	// List commands
	"RPUSH",
	"LPUSH",
	"LSET",
}

var SETs = map[string]string{} //key-value pairs
var SETsMu = sync.RWMutex{}

var HSETs = map[string]map[string]string{} //key-hash pairs
var HSETsMu = sync.RWMutex{}

var LSETs = map[string][]string{} //key-list pairs
var LSETsMu = sync.RWMutex{}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for SET command"}
	}
	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for GET command"}
	}
	key := args[0].bulk

	SETsMu.Lock()
	value, ok := SETs[key]
	SETsMu.Unlock()
	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for HSET command"}
	}
	table := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[table]; !ok { //create table
		HSETs[table] = map[string]string{}
	}
	HSETs[table][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for HGET command"}
	}
	table := args[0].bulk
	key := args[1].bulk

	HSETsMu.Lock()
	value, ok := HSETs[table][key]
	HSETsMu.Unlock()
	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for HGETALL command"}
	}
	key := args[0].bulk
	HSETsMu.Lock()
	table, ok := HSETs[key]
	HSETsMu.Unlock()
	if !ok {
		return Value{typ: "null"}
	}
	v := Value{typ: "array"}
	v.array = make([]Value, 2*len(table))
	i := 0
	for key, val := range table {
		v.array[i] = Value{typ: "bulk", bulk: key}
		i++
		v.array[i] = Value{typ: "bulk", bulk: val}
		i++
	}
	return v
}

func exists(args []Value) Value {
	length := len(args)
	if length < 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for EXISTS"}
	}
	count := 0
	for i := range length {
		key := args[i].bulk
		SETsMu.Lock()
		_, ok := SETs[key]
		SETsMu.Unlock()
		if ok {
			count++
		}
	}
	return Value{typ: "integer", num: count}
}

func _append(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for APPEND"}
	}
	key := args[0].bulk
	addedValue := args[1].bulk
	SETsMu.Lock()
	SETs[key] = SETs[key] + addedValue
	newVal := SETs[key]
	SETsMu.Unlock()
	return Value{typ: "integer", num: len(newVal)}
}
func del(args []Value) Value {
	length := len(args)
	if length < 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for DEL"}
	}
	count := 0
	for i := range length {
		key := args[i].bulk
		SETsMu.Lock()
		_, ok := SETs[key]
		if ok {
			delete(SETs, key)
			count++
		}
		SETsMu.Unlock()
	}
	return Value{typ: "integer", num: count}
}

func copy(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for COPY"}
	}
	sourceKey := args[0].bulk
	destKey := args[1].bulk
	SETsMu.Lock()
	val, ok := SETs[sourceKey]
	if !ok {
		SETsMu.Unlock()
		return Value{typ: "integer", num: 0}
	}
	SETs[destKey] = val
	SETsMu.Unlock()
	return Value{typ: "integer", num: 1}
}

func hdel(args []Value) Value {
	length := len(args)
	if length < 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for HDEL"}
	}
	hash := args[0].bulk
	count := 0
	for i := range length - 1 {
		key := args[i+1].bulk
		HSETsMu.Lock()
		_, ok := HSETs[hash][key]
		if ok { // delete key
			delete(HSETs[hash], key)
			count++
		}
		kv, ok := HSETs[hash]
		if ok && len(kv) == 0 { // delete hash if empty
			delete(HSETs, hash)
		}
		HSETsMu.Unlock()
	}
	return Value{typ: "integer", num: count}
}

func incr(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for INCR"}
	}
	key := args[0].bulk
	SETsMu.Lock()
	val, ok := SETs[key]
	SETsMu.Unlock()

	if !ok { //defaults to val 1 if not exist
		SETsMu.Lock()
		SETs[key] = "1"
		SETsMu.Unlock()
		return Value{typ: "integer", num: 1}
	}

	ival, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return Value{typ: "error", str: "ERR wrong type of argument for INCR"}
	}
	SETsMu.Lock()
	SETs[key] = strconv.Itoa(int(ival + 1))
	SETsMu.Unlock()
	return Value{typ: "integer", num: int(ival + 1)}
}

func rpush(args []Value) Value {
	length := len(args)
	if length < 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for RPUSH"}
	}
	key := args[0].bulk
	LSETsMu.Lock()
	if _, ok := LSETs[key]; !ok {
		LSETs[key] = []string{}
	}
	LSETsMu.Unlock()
	count := 0
	for i := range length - 1 {
		val := args[i+1].bulk
		LSETsMu.Lock()
		LSETs[key] = append(LSETs[key], val)
		if i == (length - 2) {
			count = len(LSETs[key])
		}
		LSETsMu.Unlock()
	}
	return Value{typ: "integer", num: count}
}

func lpush(args []Value) Value {
	length := len(args)
	if length < 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for RPUSH"}
	}
	key := args[0].bulk
	LSETsMu.Lock()
	if _, ok := LSETs[key]; !ok {
		LSETs[key] = []string{}
	}
	LSETsMu.Unlock()
	count := 0
	for i := range length - 1 {
		val := args[i+1].bulk
		LSETsMu.Lock()
		LSETs[key] = append([]string{val}, LSETs[key]...)
		if i == (length - 2) {
			count = len(LSETs[key])
		}
		LSETsMu.Unlock()
	}
	return Value{typ: "integer", num: count}
}

func llen(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for LLEN"}
	}
	key := args[0].bulk
	LSETsMu.Lock()
	list, ok := LSETs[key]
	LSETsMu.Unlock()
	if !ok {
		return Value{typ: "integer", num: 0}
	}
	return Value{typ: "integer", num: len(list)}
}

func lset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for LSET"}
	}
	key := args[0].bulk
	idx, err := strconv.Atoi(args[1].bulk)
	if err != nil {
		return Value{typ: "error", str: "ERR index must be a number"}
	}
	val := args[2].bulk
	LSETsMu.Lock()
	list, ok := LSETs[key]
	LSETsMu.Unlock()
	if !ok {
		return Value{typ: "error", str: "ERR no such key"}
	}
	length := len(list)
	if idx > (length-1) || idx < -length {
		return Value{typ: "error", str: "ERR index out of range"}
	}
	LSETsMu.Lock()
	LSETs[key][(length+idx)%length] = val
	LSETsMu.Unlock()
	return Value{typ: "string", str: "OK"}
}

func lrange(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for LRANGE"}
	}
	key := args[0].bulk
	start, err := strconv.Atoi(args[1].bulk)
	if err != nil {
		return Value{typ: "error", str: "ERR wrong type of argument for LRANGE"}
	}
	end, err := strconv.Atoi(args[2].bulk)
	if err != nil {
		return Value{typ: "error", str: "ERR wrong type of argument for LRANGE"}
	}
	LSETsMu.Lock()
	list, ok := LSETs[key]
	LSETsMu.Unlock()
	if !ok {
		return Value{typ: "array", array: []Value{}}
	}
	length := len(list)
	// slice bounds checking
	if start >= length || end < -length {
		return Value{typ: "array", array: []Value{}}
	}
	start += length
	if start <= 0 {
		start = 0
	} else {
		start %= length
	}
	end -= length
	if end > 0 {
		end = length - 1
	} else {
		end = (end + 2*length) % length
	}
	if end < start {
		return Value{typ: "array", array: []Value{}}
	}
	v := Value{typ: "array"}
	v.array = make([]Value, end+1-start)
	for i, val := range list[start : end+1] {
		v.array[i] = Value{typ: "bulk", bulk: val}
	}
	return v
}
