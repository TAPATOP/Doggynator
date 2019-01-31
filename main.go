package main

import (
	"bufio"
	"os"
	"source"
)

func main() {
	game := source.DoggynatorConstructor(
		"questions.txt",
		"records.txt",
		bufio.NewWriter(os.Stdout),
	)
	game.Play()
}

/*
Capital letters of constructors?
*/
