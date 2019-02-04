package source

import (
	"bufio"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

type Doggynator struct {
	questions []string
	records   []Record
	dbf       DataBaseOfFacts
	ie        InferenceEngine
	lm        LearningMechanism

	output       *bufio.Writer
	input        *bufio.Reader
	questionsURL string
	recordsURL   string
}

func DoggynatorConstructor(questionsURL, recordsURL string, input *bufio.Reader, output *bufio.Writer) (*Doggynator, error) {
	newObj := new(Doggynator)
	newObj.output = output
	newObj.input = input
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

func (obj *Doggynator) AddQuestion(question string) {
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
	scanner := bufio.NewScanner(obj.input)

	for questionIndex := obj.askQuestion(); true; {
		if questionIndex == -1 {
			obj.makeGuess(obj.ie.getBestGuess(), scanner)
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
			obj.makeGuess(answer, scanner)
			//obj.lm.learn(answer)
			//obj.finalizeGame()
			//return
		}
		questionIndex = obj.askQuestion()
	}
}

func (obj *Doggynator) initializeGame() {
	obj.dbf = *DataBaseOfFactsConstructor(len(obj.questions), len(obj.records))
	obj.ie = *InferenceEngineConstructor(obj.records, &obj.dbf)
	obj.lm = *LearningMechanismConstructor(&obj.dbf)
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

func (obj *Doggynator) makeGuess(answer *Record, scanner *bufio.Scanner) bool {
	obj.writeln("I believe you are thinking about: " + answer.name)
	return obj.askIfGuessIsCorrect(answer, scanner)
}

func (obj *Doggynator) askIfGuessIsCorrect(record *Record, scanner *bufio.Scanner) bool {
	for true {
		rawAnswer := receiveInput(scanner)
		response := toResponse(rawAnswer)
		switch response {
		case Response(Yes):
			return true
		case Response(No):
			return false
		}
	}
	return false
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

func (obj *Doggynator) QuestionsToString() *string {
	stringified := ""
	for _, elem := range obj.questions {
		stringified += elem + "\n"
	}
	return &stringified
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
