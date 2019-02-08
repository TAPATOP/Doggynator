package source

import (
	"math"
	"math/rand"
)

const ConclusionFactor = 1.5
const MinimumAnsweredQuestions = 3
const MinimumIntervalBetweenAnswers = 3
const MaximumIntervalBetweenAnswers = 5
const MentionReductionFactor = 3
const RandomQuestionProbability = 50

type InferenceEngine struct {
	records                  []Record
	mutt                     Record
	questions                []string
	dbf                      *DataBaseOfFacts
	recordProbability        []float64
	enquiriesSinceLastAnswer int
}

func InferenceEngineConstructor(records []Record, questions []string, dbf *DataBaseOfFacts) *InferenceEngine {
	obj := new(InferenceEngine)
	obj.records = records
	obj.questions = questions
	obj.dbf = dbf
	obj.recordProbability = make([]float64, len(records))

	obj.mutt = *EmptyRecordConstructor("mutt", len(obj.questions))
	for i := range obj.records {
		for j := range obj.questions {
			obj.mutt.statistics[j].sumWith(&obj.records[i].statistics[j])
		}
	}

	return obj
}

func (obj *InferenceEngine) concludeAnAnswer() (*Record, int) {
	obj.enquiriesSinceLastAnswer++
	if obj.dbf.recordedAnswerNumber < MinimumAnsweredQuestions ||
		obj.enquiriesSinceLastAnswer < MinimumIntervalBetweenAnswers {
		return nil, -1
	}
	candidateIndex := obj.getBestGuessIndex()
	if obj.enquiriesSinceLastAnswer <= MaximumIntervalBetweenAnswers {
		for i := range obj.recordProbability {
			if i == candidateIndex {
				continue
			}
			if math.Abs(obj.recordProbability[i]-obj.recordProbability[candidateIndex]) < ConclusionFactor {
				return nil, -1
			}
		}
	}
	obj.enquiriesSinceLastAnswer = 0
	return &obj.records[candidateIndex], candidateIndex
}

func (obj *InferenceEngine) getBestGuessIndex() int {
	candidateIndex := 0
	for i := range obj.recordProbability[1:] {
		if obj.recordProbability[candidateIndex] < obj.recordProbability[i+1] {
			candidateIndex = i + 1
		}
	}
	return candidateIndex
}

func (obj *InferenceEngine) processResponse(questionIndex int, response Response) {
	switch response {
	case Response(Yes):
		obj.calculateAllProbabilitiesOfAnswer(questionIndex, response.Integer(), 1)
	case Response(No):
		obj.calculateAllProbabilitiesOfAnswer(questionIndex, response.Integer(), 1)
	case Response(DontKnowOrIrrelevant):
		obj.calculateAllProbabilitiesOfAnswer(questionIndex, response.Integer(), 1)
	case Response(ProbablyYes):
		obj.calculateAllProbabilitiesOfAnswer(questionIndex, Response(Yes).Integer(), ProbablyModifier)
	case Response(ProbablyNo):
		obj.calculateAllProbabilitiesOfAnswer(questionIndex, Response(No).Integer(), ProbablyModifier)
	}
}

func (obj *InferenceEngine) calculateAllProbabilitiesOfAnswer(
	questionIndex int,
	answer int, // this is technically the array index of the answer in statistics' data...
	modifier float64,
) {
	obj.dbf.record(answer, questionIndex)
	for i := range obj.records {
		valueForMultiplication := obj.records[i].statistics[questionIndex].getProbability(answer)
		if valueForMultiplication < MinimumProbability {
			valueForMultiplication = MinimumProbability
		}
		obj.recordProbability[i] += math.Log2(valueForMultiplication) * modifier
	}
}

func (obj *InferenceEngine) getBestGuess() *Record {
	return &obj.records[obj.getBestGuessIndex()]
}

func (obj *InferenceEngine) reduceProbability(index int) {
	obj.recordProbability[index] -= MentionReductionFactor
}

func (obj *InferenceEngine) askQuestion() (index int) {
	if obj.dbf.hasBeenAskedEveryQuestion() {
		return -1
	}
	questionIndex := obj.chooseQuestionIndex()
	return questionIndex
}

func (obj *InferenceEngine) chooseQuestionIndex() int {
	randomNum := rand.Intn(100)
	if randomNum > RandomQuestionProbability {
		index := rand.Intn(len(obj.questions))
		for obj.dbf.isAsked(index) {
			index = rand.Intn(len(obj.questions))
		}
		return index
	}
	return obj.getHighestEntropyIndex()
}

func (obj *InferenceEngine) getHighestEntropyIndex() int {
	highestIndex := 0
	for i := 1; i < len(obj.questions); i++ {
		if !obj.dbf.isAsked(i) {
			if obj.dbf.isAsked(highestIndex) {
				highestIndex = i
				continue
			}
			if obj.mutt.statistics[i].entropy() > obj.mutt.statistics[highestIndex].entropy() {
				highestIndex = i
			}
		}
	}
	return highestIndex
}
