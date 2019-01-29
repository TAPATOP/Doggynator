package main

import (
	"fmt"
	"source"
)

func main() {
	game := source.DoggynatorConstructor("questions.txt")
	fmt.Println(game.Records[0].ToString())
}
