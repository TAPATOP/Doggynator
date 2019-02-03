package source

import (
	"bufio"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Doggynator struct {
	questions []string
	records   []Record
	dbf       DataBaseOfFacts
	ie        InferenceEngine

	output       *bufio.Writer
	questionsURL string
	recordsURL   string
}

type Response int

const (
	Yes Response = iota
	No
	DontKnowOrIrrelevant
	ProbablyYes
	ProbablyNo
	IncorrectResponse
)

func (resp Response) Integer() int {
	return [...]int{1, 2, 3, 4, 5}[resp]
}

func DoggynatorConstructor(questionsURL, recordsURL string, output *bufio.Writer) (*Doggynator, error) {
	newObj := new(Doggynator)
	newObj.output = output
	newObj.questionsURL = questionsURL
	newObj.recordsURL = recordsURL

	err := newObj.loadQuestions(questionsURL)
	if err != nil {
		newObj.writeln("Error loading questions!")
		return nil, err
	}

	err = newObj.loadRecords(recordsURL)
	if err != nil {
		newObj.writeln("Error loading records!")
		return nil, err
	}
	return newObj, nil
}

// Question section //

func (obj *Doggynator) loadQuestions(questionsURL string) error {
	data, err := ioutil.ReadFile(questionsURL)
	if err != nil {
		return err
	}
	obj.questions = filter(strings.Split(string(data), "\n"))
	return nil
}

func (obj *Doggynator) saveQuestions(questionsURL string) (err error) {
	err = ioutil.WriteFile(questionsURL, []byte(*(obj.QuestionsToString())), 0644)
	if err != nil {
		return
	}
	return
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
		return err
	}
	rawRecords := filter(strings.Split(string(data), "\n"))
	obj.records, err = processRawRecords(rawRecords, len(obj.questions))
	if err != nil {
		obj.writeln("Error processing raw records")
		return err
	}
	return nil
}

func processRawRecords(rawRecords []string, numberOfQuestions int) (records []Record, err error) {
	currRecordName := rawRecords[0]
	var currRecordData []Statistic
	counter := 1
	numberOfQuestions++

	for _, elem := range rawRecords[1:] {
		if counter%numberOfQuestions == 0 {
			records = append(records, *RecordConstructor(currRecordName, currRecordData))
			currRecordName = elem
			currRecordData = []Statistic{}
		} else {
			newStat, err := RawStatisticConstructor(elem)
			if err != nil {
				return nil, err
			}
			currRecordData = append(currRecordData, *newStat)
		}
		counter++
	}
	return append(records, *RecordConstructor(currRecordName, currRecordData)), nil
}

func (obj *Doggynator) saveRecords(recordsURL string) (err error) {
	var stringified string
	for _, elem := range obj.records {
		stringified += elem.ToString()
	}
	err = ioutil.WriteFile(recordsURL, []byte(stringified), 0644)
	if err != nil {
		return
	}
	return nil
}

// Playing Section //

func (obj *Doggynator) Play() {
	obj.initializeGame()
	scanner := bufio.NewScanner(os.Stdin)

	for questionIndex := obj.askQuestion(); true; {
		if questionIndex == -1 {
			obj.writeln(
				"I'm out of questions, so this is my final conclusion: " +
					obj.ie.getBestGuess().name,
			)
			break
		}
		obj.writeln(obj.questions[questionIndex])
		rawResponse := receiveInput(scanner)
		response := toResponse(rawResponse)
		if response == Response(IncorrectResponse) {
			obj.writeln("Bad answer!")
			continue
		}
		obj.processResponse(questionIndex, response)
		answer := obj.ie.concludeAnAnswer()
		if answer != nil {
			obj.writeln("I think you are thinking about: " + answer.name)
		}
		questionIndex = obj.askQuestion()
	}
}

func (obj *Doggynator) initializeGame() {
	obj.dbf = *DataBaseOfFactsConstructor(len(obj.questions), len(obj.records))
	obj.ie = *InferenceEngineConstructor(obj.records, &obj.dbf)
	rand.Seed(time.Now().UTC().UnixNano())
}

func (obj *Doggynator) askQuestion() (index int) {
	if obj.dbf.hasBeenAskedEveryQuestion() {
		return -1
	}
	questionIndex := obj.chooseQuestionIndex()
	return questionIndex
}

func (obj *Doggynator) chooseQuestionIndex() int {
	num := rand.Intn(len(obj.questions))
	for obj.dbf.isAsked(num) {
		num = rand.Intn(len(obj.questions))
	}
	return num
}

func (obj *Doggynator) processResponse(questionIndex int, response Response) {
	obj.dbf.processResponse(questionIndex, obj.records, response)
}

func (obj *Doggynator) finalizeGame() {
	obj.saveQuestions(obj.questionsURL)
	obj.saveRecords(obj.recordsURL)
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

func (obj *Doggynator) QuestionsToString() (stringified *string) {
	for _, elem := range obj.questions {
		*stringified += elem + "\n"
	}
	return
}

func (obj *Doggynator) write(message string) {
	obj.output.WriteString(message)
	obj.output.Flush()
}

func (obj *Doggynator) writeln(message string) {
	obj.write(message + "\n")
}

func receiveInput(scanner *bufio.Scanner) string {
	scanner.Scan()
	return scanner.Text()
}

func toResponse(forConverting string) (value Response) {
	switch forConverting {
	case "yes", "y":
		return Response(Yes)
	case "no", "n":
		return Response(No)
	case "irrelevant", "don't know", "no idea", "irr", "dk":
		return Response(DontKnowOrIrrelevant)
	case "probably", "p", "prob":
		return Response(ProbablyYes)
	case "probably not", "pn", "prob no":
		return Response(ProbablyNo)
	default:
		return Response(IncorrectResponse)
	}
}
