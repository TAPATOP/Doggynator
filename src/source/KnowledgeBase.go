package source

const IsAsked = 1

type DataBaseOfFacts struct {
	answers              []int
	answeredIndexes      []int
	recordedAnswerNumber int
}

func DataBaseOfFactsConstructor(size int) *DataBaseOfFacts {
	obj := new(DataBaseOfFacts)
	obj.answers = make([]int, size)
	obj.answeredIndexes = make([]int, size)
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
