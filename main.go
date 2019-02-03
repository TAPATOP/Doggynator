package main

import (
	"bufio"
	"fmt"
	"os"
	"source"
)

func main() {
	game, err := source.DoggynatorConstructor(
		"DogQuestions.txt",
		"DogRecords.txt",
		bufio.NewWriter(os.Stdout),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	game.Play()
}

/*
Accidental and unneeded object copy checks
Capital letters of constructors?
fix these pieces of crap obj.dbf.recordProbability[candidateIndex]
return pointers where possible
reduce the probability of a wrong answer
slices: store object addresses, not objects?
*/
