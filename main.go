package main

import (
	"bufio"
	"fmt"
	"os"
	"source"
)

func main() {
	game, err := source.DoggynatorConstructor(
		"AnimalQuestions.txt",
		"AnimalRecords.txt",
		bufio.NewReader(os.Stdin),
		bufio.NewWriter(os.Stdout),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	game.Start()
}

/*
Accidental and unneeded object copy checks
Capital letters of constructors?
return pointers where possible
remove "very intelligent"
fucking simplify the play function
IE: valueForMultiplication == 0
hide the objects call in a function
*/
