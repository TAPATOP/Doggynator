package source

import "math"

const IsAsked = 1

//const ProbablyModifier = 0.5
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
		obj.calculateAllProbabilitiesOfAnswer(questionIndex, records, 0)
	case Response(No):
		obj.calculateAllProbabilitiesOfAnswer(questionIndex, records, 1)
	case Response(DontKnowOrIrrelevant):
		obj.calculateAllProbabilitiesOfAnswer(questionIndex, records, 2)
	case Response(ProbablyYes): //obj.processResponse(questionIndex, records, Response(Yes))
	case Response(ProbablyNo):
	}
}

func (obj *DataBaseOfFacts) calculateAllProbabilitiesOfAnswer(questionIndex int, records []Record, answer int) {
	obj.record(answer, questionIndex)
	for i := range records {
		valueForMultiplication := records[i].statistics[questionIndex].getProbability(answer)
		if valueForMultiplication < MinimumProbability {
			valueForMultiplication = MinimumProbability
		}
		obj.recordProbability[i] += math.Log10(valueForMultiplication)
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
