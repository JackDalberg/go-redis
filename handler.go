package main

import (
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

var SETs = map[string]string{} //key-value pairs
var SETsMu = sync.RWMutex{}

var HSETs = map[string]map[string]string{} //key-table pairs
var HSETsMu = sync.RWMutex{}

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
