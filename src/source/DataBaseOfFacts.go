package source

import "math"

const IsAsked = 1

const ProbablyModifier = 0.5
const MinimumProbability = 0.05

type DataBaseOfFacts struct {
	answers              []int
	answeredIndexes      []int
	recordProbability    []float64
	recordedAnswerNumber int
}

func DataBaseOfFactsConstructor(questionCount, recordsCount int) *DataBaseOfFacts {
	obj := new(DataBaseOfFacts)
	obj.answers = make([]int, questionCount)
	obj.answeredIndexes = make([]int, questionCount)

	obj.recordProbability = make([]float64, recordsCount)
	for index := range obj.recordProbability {
		obj.recordProbability[index] = 0
	}

	return obj
}

func (obj *DataBaseOfFacts) processResponse(questionIndex int, records []Record, response Response) {
	switch response {
	case Response(Yes):
		obj.calculateAllProbabilitiesOfAnswer(questionIndex, records, 0, 1)
	case Response(No):
		obj.calculateAllProbabilitiesOfAnswer(questionIndex, records, 1, 1)
	case Response(DontKnowOrIrrelevant):
		obj.calculateAllProbabilitiesOfAnswer(questionIndex, records, 2, 1)
	case Response(ProbablyYes):
		obj.calculateAllProbabilitiesOfAnswer(questionIndex, records, 0, ProbablyModifier)
	case Response(ProbablyNo):
		obj.calculateAllProbabilitiesOfAnswer(questionIndex, records, 1, ProbablyModifier)
	}
}

func (obj *DataBaseOfFacts) calculateAllProbabilitiesOfAnswer(
	questionIndex int,
	records []Record,
	answer int, // this is technically the array index of the answer in statistics' data...
	modifier float64,
) {
	obj.record(answer, questionIndex)
	for i := range records {
		valueForMultiplication := records[i].statistics[questionIndex].getProbability(answer)
		if valueForMultiplication < MinimumProbability {
			valueForMultiplication = MinimumProbability
		}
		obj.recordProbability[i] += math.Log2(valueForMultiplication) * modifier
	}
}

func (obj *DataBaseOfFacts) record(value, index int) {
	if index >= len(obj.answers) {
		return
	}
	obj.answers[index] = value
	obj.answeredIndexes[index] = IsAsked
	obj.recordedAnswerNumber++
}

func (obj *DataBaseOfFacts) isAsked(index int) bool {
	if index >= len(obj.answers) || obj.recordedAnswerNumber >= len(obj.answers) {
		return true
	}
	return obj.answeredIndexes[index] == IsAsked
}

func (obj *DataBaseOfFacts) hasBeenAskedEveryQuestion() bool {
	return obj.recordedAnswerNumber >= len(obj.answeredIndexes)
}

//func (obj *DataBaseOfFacts) saveInputToKnowledgeBase(record *Record) {
//
//}
