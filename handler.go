package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"HSET": hset,
	"HGET": hget,
}

var SETs = map[string]string{}
var SETsMU = sync.RWMutex{}

var HSETs = map[string]map[string]string{}
var HSETsMU = sync.RWMutex{}

// PING Command
func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number for SET args"}
	}
	key := args[0].bulk
	value := args[1].bulk
	SETsMU.Lock()
	SETs[key] = value
	SETsMU.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number for SET args"}
	}

	key := args[0].bulk
	SETsMU.RLock()
	value, ok := SETs[key]
	SETsMU.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number for HSET args"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMU.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMU.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {

	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number for HGET args"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMU.RLock()
	value, ok := HSETs[hash][key]
	HSETsMU.RUnlock()
	
	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}
