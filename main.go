package main

import (
	"fmt"
	"source"
)

func main() {
	game := source.DoggynatorConstructor("questions.txt", "records.txt")
	fmt.Println(game.QuestionsToString())
}
