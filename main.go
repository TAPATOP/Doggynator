package main

import (
	"bufio"
	"fmt"
	"os"
	"source"
)

func main() {
	game, err := source.DoggynatorConstructor(
		"questions.txt",
		"records.txt",
		bufio.NewWriter(os.Stdout),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	game.Play()
}

/*
Capital letters of constructors?
*/
