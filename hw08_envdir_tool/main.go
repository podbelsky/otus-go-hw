package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatal("too few arguments")
	}

	dir, err := os.Stat(args[1])
	if err != nil {
		log.Fatal(err)
	}

	if !dir.IsDir() {
		log.Fatal("First argument is not a dir")
	}

	env, err := ReadDir(args[1])
	if err != nil {
		log.Fatal(err)
	}

	rc := RunCmd(args[2:], env)

	os.Exit(rc)
}
