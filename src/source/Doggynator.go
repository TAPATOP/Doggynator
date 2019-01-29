package source

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Doggynator struct {
	Questions []string
	Records   []Record
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
	newObj.Records = []Record{*rec0, *rec1, *rec2}
	//newObj.saveRecords("records.txt")
	return newObj
}

func (obj *Doggynator) loadQuestions(questionsURL string) {
	data, err := ioutil.ReadFile(questionsURL)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	obj.Questions = strings.Split(string(data), "\n")
}

func (obj *Doggynator) loadRecords(recordsURL string) {
}

func (obj *Doggynator) saveRecords(recordsURL string) {
}

func (obj *Doggynator) ToString() []string {
	return obj.Questions
}
