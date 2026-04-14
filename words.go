package main

import "math/rand"

var easyWords = []string{
	"go", "run", "get", "set", "map", "key", "val", "put", "top", "low",
	"net", "web", "api", "url", "log", "err", "nil", "var", "int", "str",
	"add", "sub", "mul", "div", "mod", "max", "min", "sum", "cat", "dog",
	"fox", "bit", "hex", "bin", "bus", "cmd", "env", "zip", "tar", "raw",
}

var mediumWords = []string{
	"serve", "build", "parse", "query", "route", "redis", "token", "cache",
	"retry", "fetch", "write", "close", "flush", "reset", "spawn", "check",
	"stdin", "proxy", "debug", "merge", "async", "error", "panic", "defer",
	"mutex", "scope", "stack", "queue", "graph", "slice", "chunk", "batch",
}

var hardWords = []string{
	"goroutine", "interface", "serialize", "middleware", "benchmark",
	"websocket", "marshaler", "interrupt", "signature", "recursive",
	"algorithm", "container", "subscriber", "scheduler", "heartbeat",
	"bootstrap", "implement", "concurrent", "dispatcher", "propagate",
}

func wordList(score int) []string {
	switch {
	case score >= 30:
		return hardWords
	case score >= 10:
		return mediumWords
	default:
		return easyWords
	}
}

func spawnWord(score, width int) Word {
	words := wordList(score)
	text := words[rand.Intn(len(words))]

	maxX := width - len(text)
	if maxX < 0 {
		maxX = 0
	}
	x := 0
	if maxX > 0 {
		x = rand.Intn(maxX)
	}

	return Word{Text: text, X: x, Y: 0}
}
