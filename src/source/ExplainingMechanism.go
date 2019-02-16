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
				explanation += "\"" + response.toString() + "\" to question \"" + obj.questions[i] + "\"\n"
			} else {
				surprised += "\"" + response.toString() + "\" to question \"" + obj.questions[i] + "\""
				surprised += ". I expected \"" + mostProbableAnswer.toString() + "\"\n"
			}
		}
	}
	explanation = explanation[:(len(explanation) - 1)]
	return &explanation, &surprised
}
