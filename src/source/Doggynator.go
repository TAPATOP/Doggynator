package source

import (
	"bufio"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Doggynator struct {
	questions []string
	records   []Record
	dbf       DataBaseOfFacts
	ie        InferenceEngine
	lm        LearningMechanism
	em        ExplainingMechanism

	output       *bufio.Writer
	input        *bufio.Scanner
	questionsURL string
	recordsURL   string
}

func DoggynatorConstructor(questionsURL, recordsURL string, input *bufio.Reader, output *bufio.Writer) (*Doggynator, error) {
	newObj := new(Doggynator)
	newObj.output = output
	newObj.input = bufio.NewScanner(input)
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

func (obj *Doggynator) addRecord(scanner *bufio.Scanner) {
	obj.writeln("Looks like I don't know what you're talking about. Please, tell me what that is")
	recordName := strings.ToLower(receiveInput(scanner))
	recordWithSameName := obj.contains(recordName)
	if recordWithSameName == nil {
		recordWithSameName = EmptyRecordConstructor(recordName, len(obj.questions))
		obj.records = append(obj.records, *recordWithSameName)
	}
	obj.lm.learn(recordWithSameName)
	obj.writeln("Thank you")
}

func (obj *Doggynator) contains(str string) *Record {
	for i := range obj.records {
		if obj.records[i].name == str {
			return &obj.records[i]
		}
	}
	return nil
}

// Playing Section //

func (obj *Doggynator) Play() {
	obj.initializeGame()

	for questionIndex := obj.askQuestion(); true; {
		if questionIndex == -1 {
			bestGuess := obj.ie.getBestGuess()
			hasGuessed := obj.makeGuess(bestGuess, obj.input)
			if hasGuessed {
				obj.processCorrectGuess(bestGuess, obj.input)
			} else {
				obj.addRecord(obj.input)
			}
			break
		}
		obj.writeln(obj.questions[questionIndex])
		rawResponse := receiveInput(obj.input)
		response := stringToResponse(rawResponse)
		if response == Response(IncorrectResponse) {
			obj.writeln("Bad answer!")
			continue
		}
		obj.processResponse(questionIndex, response)
		answer, indexOfAnswer := obj.ie.concludeAnAnswer()
		if answer != nil {
			hasGuessed := obj.makeGuess(answer, obj.input)
			if hasGuessed {
				obj.processCorrectGuess(answer, obj.input)
				break
			} else {
				obj.ie.reduceProbability(indexOfAnswer)
				obj.writeln("Do you want to keep playing?")
				wantsToContinue := obj.askForYesOrNo(obj.input)
				if wantsToContinue == Response(No) {
					obj.addRecord(obj.input)
					break
				}
			}
		}
		questionIndex = obj.askQuestion()
	}
	obj.finalizeGame()
}

func (obj *Doggynator) initializeGame() {
	obj.dbf = *DataBaseOfFactsConstructor(len(obj.questions), len(obj.records))
	obj.ie = *InferenceEngineConstructor(obj.records, &obj.dbf)
	obj.lm = *LearningMechanismConstructor(&obj.dbf)
	obj.em = *ExplainingMechanismConstructor(obj.questions, &obj.dbf)
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
	obj.writeln("Please say \"yes\" if I'm correct and \"no\" if I'm not")
	return obj.askIfGuessIsCorrect(scanner)
}

func (obj *Doggynator) askIfGuessIsCorrect(scanner *bufio.Scanner) bool {
	response := obj.askForYesOrNo(scanner)
	switch response {
	case Response(Yes):
		return true
	case Response(No):
		return false
	default:
		return false
	}
}

func (obj *Doggynator) askForYesOrNo(scanner *bufio.Scanner) Response {
	for true {
		rawAnswer := receiveInput(scanner)
		response := stringToResponse(rawAnswer)
		if response == Response(Yes) || response == Response(No) {
			return response
		} else {
			obj.writeln("Please answer with a \"yes\" or \"no\"")
		}
	}
	return Response(IncorrectResponse)
}

func (obj *Doggynator) boast(scanner *bufio.Scanner) {
	obj.writeln("Heh, I'm so smart")
}

func (obj *Doggynator) processCorrectGuess(guess *Record, scanner *bufio.Scanner) {
	obj.writeln("I made this guess because you gave the following answers:")
	obj.writeln(*obj.em.explain(guess))
	obj.lm.learn(guess)
	//obj.boast(scanner)
}

// Menu //

func (obj *Doggynator) Start() {
	obj.writeln("Hello to Doggynator!")
	for true {
		obj.writeln("Please press the following numbers to execute the commands")
		obj.writeln("1: Play")
		obj.writeln("2: Add Question")
		obj.writeln("3: Exit")
		input, err := strconv.Atoi(receiveInput(obj.input))
		if err != nil {
			obj.writeln("Oops, looks like you messed up! Try the commands again")
			continue
		}
		switch input {
		case 1:
			obj.Play()
		case 2:
			obj.writeln("Input your question")
			question := receiveInput(obj.input)
			obj.AddQuestion(question)
			obj.finalizeGame()
		case 3:
			obj.writeln("Goodgye")
			return
		}
	}
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
