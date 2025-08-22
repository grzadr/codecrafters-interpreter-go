package main

import (
	"log"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/scanner"
)

const defaultArgsNum = 3

func main() {
	if len(os.Args) < defaultArgsNum {
		log.Fatalln("Usage: ./your_program.sh tokenize|parse <filename>")
	}

	switch command := os.Args[1]; command {
	case "tokenize":
		filename := os.Args[2]

		code, err := scanner.CmdTokenize(filename)
		if err != nil {
			log.Fatalf("error tokenizing %q: %s", filename, err)
		}

		os.Exit(code)

	default:
		log.Fatalf("Unknown command: %s\n", command)
	}
}
