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

func (obj *ExplainingMechanism) explain(record *Record) *string {
	explanation := ""
	for i := range obj.dbf.answeredIndexes {
		if obj.dbf.answeredIndexes[i] == IsAsked {
			response := Response(obj.dbf.answers[i])
			if response == record.statistics[i].mostProbableAnswerToAttribute() {
				explanation += "\"" + response.toString() + "\" to question \"" + obj.questions[i] + "\"\n"
			}
		}
	}
	explanation = explanation[:(len(explanation) - 1)]
	return &explanation
}
