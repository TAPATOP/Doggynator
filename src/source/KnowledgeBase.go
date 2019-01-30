package source

const IsAsked = 1

type KnowledgeBase struct {
	answers              []int
	answeredIndexes      []int
	recordedAnswerNumber int
}

func KnowledgeBaseConstructor(size int) *KnowledgeBase {
	obj := new(KnowledgeBase)
	obj.answers = make([]int, size)
	obj.answeredIndexes = make([]int, size)
	return obj
}

func (obj *KnowledgeBase) record(value, index int) {
	if index >= len(obj.answers) {
		return
	}
	obj.answers[index] = value
	obj.answeredIndexes[index] = IsAsked
	obj.recordedAnswerNumber++
}

func (obj *KnowledgeBase) isAsked(index int) bool {
	if index >= len(obj.answers) || obj.recordedAnswerNumber >= len(obj.answers) {
		return true
	}
	return obj.answeredIndexes[index] == IsAsked
}

func (obj *KnowledgeBase) hasBeenAskedEveryQuestion() bool {
	return obj.recordedAnswerNumber >= len(obj.answeredIndexes)
}
