package source

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

type Doggynator struct {
	questions []string
	records   []Record
	kb        KnowledgeBase
}

func DoggynatorConstructor(questionsURL, recordsURL string) *Doggynator {
	newObj := new(Doggynator)
	err := newObj.loadQuestions(questionsURL)
	if err != nil {
		fmt.Println("Error loading questions!")
		return nil
	}

	newObj.loadRecords(recordsURL)
	if err != nil {
		fmt.Println("Error loading records!")
		return nil
	}

	newObj.saveQuestions("questions.txt")
	newObj.saveRecords("records.txt")
	return newObj
}

// Question section //

func (obj *Doggynator) loadQuestions(questionsURL string) error {
	data, err := ioutil.ReadFile(questionsURL)
	if err != nil {
		fmt.Println("Questions reading error", err)
		return err
	}
	obj.questions = filter(strings.Split(string(data), "\n"))
	return nil
}

func (obj *Doggynator) saveQuestions(questionsURL string) {
	err := ioutil.WriteFile(questionsURL, []byte(obj.QuestionsToString()), 0644)
	if err != nil {
		fmt.Println("Questions saving error", err)
		return
	}
}

func (obj *Doggynator) addQuestion(question string) {
	if question != "" && question != "\n" && question != "\t" && question != " " {
		obj.questions = append(obj.questions, question)
		for i := 0; i < len(obj.records); i++ {
			obj.records[i].AddField()
		}
	}
}

// Records Section //

func (obj *Doggynator) loadRecords(recordsURL string) error {
	data, err := ioutil.ReadFile(recordsURL)
	if err != nil {
		fmt.Println("Records reading error", err)
		return err
	}
	rawRecords := filter(strings.Split(string(data), "\n"))
	obj.records, err = processRawRecords(rawRecords, len(obj.questions))
	return nil
}

func processRawRecords(rawRecords []string, numberOfQuestions int) (records []Record, err error) {
	currRecordName := rawRecords[0]
	currRecordData := []Statistic{}
	counter := 1
	numberOfQuestions++

	for _, elem := range rawRecords[1:] {
		if counter%numberOfQuestions == 0 {
			records = append(records, *RecordConstructor(currRecordName, currRecordData))
			currRecordName = elem
			currRecordData = []Statistic{}
		} else {
			newStat, err := RawStatisticConstructor(elem[:(len(elem) - 1)])
			if err != nil {
				fmt.Println("Failed loading raw record of Statistic", err)
				return nil, err
			}
			currRecordData = append(currRecordData, *newStat)
		}
		counter++
	}
	return append(records, *RecordConstructor(currRecordName, currRecordData)), nil
}

func (obj *Doggynator) saveRecords(recordsURL string) {
	var stringified string
	for _, elem := range obj.records {
		stringified += elem.ToString()
	}
	err := ioutil.WriteFile(recordsURL, []byte(stringified), 0644)
	if err != nil {
		fmt.Println("There was an issue with saving the records to a file")
	}
}

// Playing Section //

func (obj *Doggynator) Play() {
	obj.initializeGame()
	//obj.askQuestion()
	//obj.askQuestion()
	//obj.askQuestion()
	//obj.askQuestion()
}

func (obj *Doggynator) initializeGame() {
	obj.kb = *KnowledgeBaseConstructor(len(obj.questions))
	rand.Seed(time.Now().UTC().UnixNano())
}

func (obj *Doggynator) askQuestion() {
	if obj.kb.hasBeenAskedEveryQuestion() {
		fmt.Println("DONT ASK ME ANYMORE, I'VE ALREADY SAID EVERYTHING I KNOW!!!")
		return
	}
	questionIndex := obj.chooseQuestionIndex()
	fmt.Println(obj.questions[questionIndex])
	obj.kb.record(0, questionIndex)
}

func (obj *Doggynator) chooseQuestionIndex() int {
	num := rand.Intn(len(obj.questions))
	for obj.kb.isAsked(num) {
		num = rand.Intn(len(obj.questions))
	}
	return num
}

// Helper Methods //

// TODO: Could be faster?
func filter(input []string) (output []string) {
	for _, elem := range input {
		if elem != "" {
			output = append(output, elem)
		}
	}
	return output
}

func (obj *Doggynator) QuestionsToString() (stringified string) {
	for _, elem := range obj.questions {
		stringified += elem + "\n"
	}
	return
}
