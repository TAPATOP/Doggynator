package source

import "math"

const ConclusionFactor = 20
const MinimumAnsweredQuestions = 3
const MinimumIntervalBetweenAnswers = 3
const MaximumIntervalBetweenAnswers = 5
const MentionReductionFactor = 0.0001

type InferenceEngine struct {
	records                  []Record
	dbf                      *DataBaseOfFacts
	enquiriesSinceLastAnswer int
}

func InferenceEngineConstructor(records []Record, dbf *DataBaseOfFacts) *InferenceEngine {
	obj := new(InferenceEngine)
	obj.records = records
	obj.dbf = dbf
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
		for i := range obj.dbf.recordProbability {
			if i == candidateIndex {
				continue
			}
			if math.Abs(obj.dbf.recordProbability[i]/obj.dbf.recordProbability[candidateIndex]) < ConclusionFactor {
				return nil, -1
			}
		}
	}
	obj.enquiriesSinceLastAnswer = 0
	return &obj.records[candidateIndex], candidateIndex
}

func (obj *InferenceEngine) getBestGuessIndex() int {
	candidateIndex := 0
	for i := range obj.dbf.recordProbability[1:] {
		if obj.dbf.recordProbability[candidateIndex] < obj.dbf.recordProbability[i+1] {
			candidateIndex = i + 1
		}
	}
	return candidateIndex
}

func (obj *InferenceEngine) getBestGuess() *Record {
	return &obj.records[obj.getBestGuessIndex()]
}

func (obj *InferenceEngine) reduceProbability(index int) {
	obj.dbf.recordProbability[index] *= MentionReductionFactor
}
