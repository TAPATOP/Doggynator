package source

const IsAsked = 1

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
		obj.recordProbability[index] = 1
	}

	return obj
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
