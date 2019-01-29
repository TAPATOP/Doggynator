package source

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Doggynator struct {
	questions []string
	records   []Record
}

func DoggynatorConstructor(questionsURL, recordsURL string) *Doggynator {
	newObj := new(Doggynator)
	newObj.loadQuestions(questionsURL)
	newObj.loadRecords(recordsURL)

	//stat1 := Statistic{[...]int{1, 22, 0}}
	//stat2 := Statistic{[...]int{1, 20, 0}}
	//rec0 := RecordConstructor("rec0", []Statistic{
	//	stat1,
	//	stat2,
	//})
	//rec1 := RecordConstructor("rec1", []Statistic{
	//	{[...]int{1, 22, 0}},
	//	{[...]int{1, 22, 0}},
	//})
	//rec2 := RecordConstructor("rec2", []Statistic{
	//	{[...]int{1, 220, 0}},
	//	{[...]int{1, 220, 0}},
	//})
	//newObj.records = []Record{*rec0, *rec1, *rec2}
	newObj.saveQuestions("questions.txt")

	//newObj.saveRecords("records.txt")
	return newObj
}

func (obj *Doggynator) loadQuestions(questionsURL string) {
	data, err := ioutil.ReadFile(questionsURL)
	if err != nil {
		fmt.Println("Questions reading error", err)
		return
	}
	obj.questions = filter(strings.Split(string(data), "\n"))
}

// TODO: Could be faster?
func filter(input []string) (output []string) {
	for _, elem := range input {
		if elem != "" {
			output = append(output, elem)
		}
	}
	return output
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
	}
}

func (obj *Doggynator) loadRecords(recordsURL string) {
	data, err := ioutil.ReadFile(recordsURL)
	if err != nil {
		fmt.Println("Records reading error", err)
		return
	}
	rawRecords := filter(strings.Split(string(data), "\n"))
	obj.records = processRawRecords(rawRecords, len(obj.questions))
	fmt.Println("Didn't die")
}

func processRawRecords(rawRecords []string, numberOfQuestions int) (records []Record) {
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
				return nil
			}
			currRecordData = append(currRecordData, *newStat)
		}
		counter++
	}
	return append(records, *RecordConstructor(currRecordName, currRecordData))
}

func (obj *Doggynator) saveRecords(recordsURL string) {

}

func (obj *Doggynator) QuestionsToString() (stringified string) {
	for _, elem := range obj.questions {
		stringified += elem + "\n"
	}
	return
}
