package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/token"
)

const defaultArgsNum = 3

func main() {
	// You can use print statements as follows for debugging, they'll be visible
	// when running tests.
	// fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")
	if len(os.Args) < defaultArgsNum {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	code := 0

	switch command := os.Args[1]; command {
	case "tokenize":
		filename := os.Args[2]
		for t := range token.Tokenize(filename) {
			log.Println("got token", t)
			fmt.Println(t)

			if t.IsError() {
				code = 65
			}
		}
	default:
		log.Fatalf("Unknown command: %s\n", command)
	}

	log.Println("exiting")

	os.Exit(code)
}
