package source

const IsAnswered = 1
const ProbablyModifier = 0.5
const MinimumProbability = 0.05

type DataBaseOfFacts struct {
	answers              []int
	answeredIndexes      []int
	recordedAnswerNumber int
}

func DataBaseOfFactsConstructor(questionCount int) *DataBaseOfFacts {
	obj := new(DataBaseOfFacts)
	obj.answers = make([]int, questionCount)
	obj.answeredIndexes = make([]int, questionCount)

	return obj
}

func CustomDataBaseOfFactsConstructor(answers, answeredIndexes []int) *DataBaseOfFacts {
	numberOfAnswers := 0
	for i := range answeredIndexes {
		if answeredIndexes[i] == IsAnswered {
			numberOfAnswers++
		}
	}

	obj := new(DataBaseOfFacts)
	obj.answers = answers
	obj.answeredIndexes = answeredIndexes
	obj.recordedAnswerNumber = numberOfAnswers
	return obj
}

func (obj *DataBaseOfFacts) record(value, index int) {
	if index >= len(obj.answers) {
		return
	}
	obj.answers[index] = value
	obj.answeredIndexes[index] = IsAnswered
	obj.recordedAnswerNumber++
}

func (obj *DataBaseOfFacts) isAsked(index int) bool {
	if index >= len(obj.answers) || obj.recordedAnswerNumber >= len(obj.answers) {
		return true
	}
	return obj.answeredIndexes[index] == IsAnswered
}

func (obj *DataBaseOfFacts) hasBeenAskedEveryQuestion() bool {
	return obj.recordedAnswerNumber >= len(obj.answeredIndexes)
}

//func (obj *DataBaseOfFacts) saveInputToKnowledgeBase(record *Record) {
//
//}
