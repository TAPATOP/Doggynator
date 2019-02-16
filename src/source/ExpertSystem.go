package source

import (
	"bufio"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// ExpertSystem is the main Expert System file
// It's called like this becaue of it's original application sphere and
// because it would take too long to fix it now
// It's main components are KnowledgeBase, DatabaseOfFacts,
// InferenceEngine, LearningMechanism and ExplainingMechanism
// KnowledgeBase, represented by "questions" and "records and should
// be separate object, but again, time is against me
// The other components will be described in their separate definitions
// Other than these building blocks, the System also has an "output" and
// "input" objects in order for the informtion flow to be easily redirectable
// (for example this would come very handy in the case of implementing a GUI)
// It also memorises where its questions and records are stored in order
// to have better control over the actions over them.
//
// If you want to start the system, you would call ExpertSystem.Start(), which
// shows a menu in the console, at least for now.
type ExpertSystem struct {
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

// DoggynatorConstructor is the standard constructor
// and other than copying the passed parameters, it also
// loads the questions and records data from their files
func DoggynatorConstructor(questionsURL, recordsURL string, input *bufio.Reader, output *bufio.Writer) (*ExpertSystem, error) {
	newObj := new(ExpertSystem)
	newObj.output = output
	newObj.input = bufio.NewScanner(input)
	newObj.questionsURL = questionsURL
	newObj.recordsURL = recordsURL

	if err := newObj.loadQuestions(); err != nil {
		newObj.writeln("Error loading questions!")
		newObj.writeErr(err)
		return nil, err
	}

	if err := newObj.loadRecords(); err != nil {
		newObj.writeln("Error loading records!")
		newObj.writeErr(err)
		return nil, err
	}
	return newObj, nil
}

// Question section //

// loadQuestions uses the saved inside the System
// variable questionsURL to fetch the information
// from the corresponding file
func (obj *ExpertSystem) loadQuestions() error {
	data, err := ioutil.ReadFile(obj.questionsURL)
	if err != nil {
		return err
	}
	obj.questions = filter(strings.Split(string(data), "\n"))
	return nil
}

// saveQuestions saves the questions slice into the
// file, corresponding to questionURL
func (obj *ExpertSystem) saveQuestions() (err error) {
	err = ioutil.WriteFile(obj.questionsURL, []byte(*(obj.questionsToString())), 0644)
	if err != nil {
		return
	}
	return
}

// AddQuestion is used by the player to add a new
// question. It increases the size of all records data by one
// and fills it with preset data in order to guarantee
// high starting entropy, thus ensuring we quickly gain information
// about the new question
func (obj *ExpertSystem) AddQuestion(question string) {
	if question != "" && question != "\n" && question != "\t" && question != " " {
		obj.questions = append(obj.questions, question)
		for i := 0; i < len(obj.records); i++ {
			obj.records[i].AddField()
		}
	}
}

// Records Section //

// loadRecords loads all records from the saved
// recordsURL in the ExpertSystem's members
func (obj *ExpertSystem) loadRecords() error {
	data, err := ioutil.ReadFile(obj.recordsURL)
	if err != nil {
		return err
	}
	rawRecords := filter(strings.Split(string(data), "\n"))
	obj.records, err = processRawRecords(rawRecords, len(obj.questions))
	if err != nil {
		obj.writeln("Error processing raw records")
		obj.writeErr(err)
		return err
	}
	return nil
}

// processRawRecords receives a string that has been
// directly read from a file and constructs/ returns a slice
// of Records
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

// saveRecords saves the records in their current state
// into the file that is saved within the ExpertSystem
func (obj *ExpertSystem) saveRecords() (err error) {
	var stringified string
	for _, elem := range obj.records {
		stringified += elem.ToString()
	}
	err = ioutil.WriteFile(obj.recordsURL, []byte(stringified), 0644)
	if err != nil {
		return
	}
	return nil
}

// addRecord creates a new Record with
// nullified data into the slice
func (obj *ExpertSystem) addRecord(scanner *bufio.Scanner) {
	obj.writeln("Looks like I don't know what you're talking about. Please, tell me what that is")
	recordName := strings.ToLower(receiveInput(scanner))
	recordWithSameName := obj.contains(recordName)
	if recordWithSameName == nil {
		recordWithSameName = EmptyRecordConstructor(recordName, len(obj.questions))
		obj.records = append(obj.records, *recordWithSameName)
		obj.writeln("Hm, I didn't know about this...")
	} else {
		obj.writeln("Hey, I already know about this! I will update my records.")
		obj.printExplanation(recordWithSameName)
	}
	obj.lm.learn(recordWithSameName)
	obj.writeln("Thank you\n")
}

// contains is used when checking whether a Record that
// is about to be added already exists. Returns the found Record
// if so
func (obj *ExpertSystem) contains(str string) *Record {
	for i := range obj.records {
		if obj.records[i].name == str {
			return &obj.records[i]
		}
	}
	return nil
}

// Playing Section //

// Play is where the System tries to guess
// what the player is thinking about. It begins by
// initializing the required resources. Then it starts asking
// questions according to some criteria, defined by askQuestions.
// The game makes a guess when it's out of questions to ask or
// when a certain record is a lot more probable than another one.
//
// It gets its input from the input object, declared in the System's methods.
// In case of bad input, asks the same question again
//
// In case of a correct guess, the System records the new values for the guessed Record
// In case of an incorrect guess, it asks if the player wants to keep playing.
// If he doesn't want to, asks for the name of the Record and either updates an
// old Record or creates a new one
// Same thing happens if there are no more questions to be asked
func (obj *ExpertSystem) Play() {
	obj.initializeGame()

	for questionIndex := obj.ie.askQuestion(); true; {
		if questionIndex == -1 {
			obj.processIfGameIsOver()
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
		questionIndex = obj.ie.askQuestion()
	}
	obj.finalizeGame()
}

// processIfGameIsOver handles the logic related to finishing
// the game after a final guess has been made
func (obj *ExpertSystem) processIfGameIsOver() {
	bestGuess := obj.ie.getBestGuess()
	hasGuessed := obj.makeGuess(bestGuess, obj.input)
	if hasGuessed {
		obj.processCorrectGuess(bestGuess, obj.input)
	} else {
		obj.addRecord(obj.input)
	}
}

// initializeGame creates the needed data structures
// in order for the system to work properly.
// It's used during game resets to ensure the game doesn't
// have junk left over from the last play
func (obj *ExpertSystem) initializeGame() {
	obj.dbf = *DataBaseOfFactsConstructor(len(obj.questions))
	obj.ie = *InferenceEngineConstructor(obj.records, obj.questions, &obj.dbf, &DefaultRandomGenerator{})

	obj.lm = *LearningMechanismConstructor(&obj.dbf)
	obj.em = *ExplainingMechanismConstructor(obj.questions, &obj.dbf)

	rand.Seed(time.Now().UTC().UnixNano())
}

func (obj *ExpertSystem) processResponse(questionIndex int, response Response) {
	obj.ie.processResponse(questionIndex, response)
}

// finalizeGame is the final step of finishing a game,
// it saves the currently known data into the files
func (obj *ExpertSystem) finalizeGame() {
	obj.saveQuestions()
	obj.saveRecords()
}

// makeGuess is just a formatting function for asking the player if
// the system is correct in its guess
func (obj *ExpertSystem) makeGuess(answer *Record, scanner *bufio.Scanner) bool {
	obj.writeln("I believe you are thinking about: " + answer.name)
	obj.writeln("Please say \"yes\" if I'm correct and \"no\" if I'm not")
	return obj.askIfGuessIsCorrect(scanner)
}

// askIfGuessIsCorrects asks and receives input from the player
// whether the guess was correct
func (obj *ExpertSystem) askIfGuessIsCorrect(scanner *bufio.Scanner) bool {
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

// askForYesOrNo acquires only a Yes or No answer from the player
// and in enters a loop until it receives one
func (obj *ExpertSystem) askForYesOrNo(scanner *bufio.Scanner) Response {
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

func (obj *ExpertSystem) boast(scanner *bufio.Scanner) {
	obj.writeln("Heh, I'm so smart")
}

// processCorrectGuess modifies the data of a correctly guessed
// Record with the answers it has received
func (obj *ExpertSystem) processCorrectGuess(guess *Record, scanner *bufio.Scanner) {
	obj.printExplanation(guess)
	obj.lm.learn(guess)
	//obj.boast(scanner)
}

// printExplanation prints whatever the ExplanationMechanism
// has concluded
func (obj *ExpertSystem) printExplanation(guess *Record) {
	explanation, surprised := obj.em.explain(guess)
	if len(*explanation) > 0 {
		obj.writeln("I made this guess because you gave the following answers:")
		obj.writeln(*explanation)
		obj.writeln("")
	}
	if len(*surprised) > 0 {
		obj.writeln("The following answers you gave surprised me:")
		obj.writeln(*surprised)
		obj.writeln("")
	}
}

// Menu //

// Start shows the menu up on the console and
// waits for input from the player
func (obj *ExpertSystem) Start() {
	obj.writeln("Hello to the garage- made Akinator- like thing!")
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
			obj.writeln("Goodbye")
			return
		}
	}
}

// Helper Methods //

// TODO: Could be faster?

// filter filters the empty lines
func filter(input []string) (output []string) {
	for i := range input {
		if input[i] != "" {
			if input[i][len(input[i])-1] == '\r' {
				input[i] = input[i][:len(input[i])-1]
			}
			output = append(output, input[i])
		}
	}
	return output
}

// questionToString puts all available questions in a single String
func (obj *ExpertSystem) questionsToString() *string {
	stringified := ""
	for _, elem := range obj.questions {
		stringified += elem + "\n"
	}
	return &stringified
}

// write writes to the dedicated output object
func (obj *ExpertSystem) write(message string) {
	obj.output.WriteString(message)
	obj.output.Flush()
}

// write + '\n'
func (obj *ExpertSystem) writeln(message string) {
	obj.write(message + "\n")
}

// writeErr writes an error into the dedicated
// output object
func (obj *ExpertSystem) writeErr(err error) {
	obj.writeln(err.Error())
}

// receiveInput receives input from the dedicated input object
func receiveInput(scanner *bufio.Scanner) string {
	scanner.Scan()
	return scanner.Text()
}

// RandomGenerator is used for mocking the RNG in the Unit tests
type RandomGenerator interface {
	Intn(limit int) int
}

type DefaultRandomGenerator struct{}

func (obj *DefaultRandomGenerator) Intn(limit int) int {
	return rand.Intn(limit)
}

// : Properties : //
// reduce records
// minimum answers
// minimum distance between guesses
// logarithms
// reduce probability of wrong answer
