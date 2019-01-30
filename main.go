package main

import (
	"source"
)

func main() {
	game := source.DoggynatorConstructor("questions.txt", "records.txt")
	game.Play()
}

/*
Capital letters of constructors?
*/
