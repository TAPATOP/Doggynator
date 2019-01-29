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

func DoggynatorConstructor(questionsURL string) *Doggynator {
	newObj := new(Doggynator)
	newObj.loadQuestions(questionsURL)

	stat1 := Statistic{positive: 1, negative: 22, irrelevant: 0}
	stat2 := Statistic{positive: 1, negative: 20, irrelevant: 0}
	rec0 := RecordConstructor("rec0", []Statistic{
		stat1,
		stat2,
	})
	rec1 := RecordConstructor("rec1", []Statistic{
		{positive: 1, negative: 22, irrelevant: 0},
		{positive: 1, negative: 22, irrelevant: 0},
	})
	rec2 := RecordConstructor("rec2", []Statistic{
		{positive: 1, negative: 220, irrelevant: 0},
		{positive: 1, negative: 220, irrelevant: 0},
	})
	newObj.records = []Record{*rec0, *rec1, *rec2}
	newObj.addQuestion("Is it brown")
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
}

func (obj *Doggynator) saveRecords(recordsURL string) {

}

func (obj *Doggynator) QuestionsToString() (stringified string) {
	for _, elem := range obj.questions {
		stringified += elem + "\n"
	}
	return stringified
}
