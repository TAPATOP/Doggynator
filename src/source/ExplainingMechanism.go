package source

type ExplainingMechanism struct {
	questions []string
	dbf       *DataBaseOfFacts
}

func ExplainingMechanismConstructor(questions []string, dbf *DataBaseOfFacts) *ExplainingMechanism {
	em := new(ExplainingMechanism)
	em.questions = questions
	em.dbf = dbf
	return em
}

func (obj *ExplainingMechanism) explain(record *Record) (*string, *string) {
	explanation := ""
	surprised := ""
	for i := range obj.dbf.answeredIndexes {
		if obj.dbf.answeredIndexes[i] == IsAnswered {
			response := Response(obj.dbf.answers[i])
			mostProbableAnswer := record.statistics[i].mostProbableAnswer()
			if response == mostProbableAnswer {
				explanation += generateExplanationLn(response, obj.questions[i])
			} else {
				surprised += generateSurprisedLn(response, mostProbableAnswer, obj.questions[i])
			}
		}
	}
	if len(explanation) > 0 {
		explanation = explanation[:(len(explanation) - 1)]
	}
	if len(surprised) > 0 {
		surprised = surprised[:(len(surprised) - 1)]
	}
	return &explanation, &surprised
}

func generateExplanation(response Response, question string) string {
	return "\"" + response.ToString() + "\" to question \"" + question + "\""
}

func generateExplanationLn(response Response, question string) string {
	return generateExplanation(response, question) + "\n"
}

func generateSurprised(response, expected Response, question string) string {
	return "\"" + response.ToString() + "\" to question \"" + question + "\"" +
		". I expected \"" + expected.ToString() + "\""
}

func generateSurprisedLn(response, expected Response, question string) string {
	return generateSurprised(response, expected, question) + "\n"
}
