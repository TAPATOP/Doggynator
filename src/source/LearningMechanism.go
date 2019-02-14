package source

const LearningIncrease = 1

type LearningMechanism struct {
	dbf *DataBaseOfFacts
}

func LearningMechanismConstructor(dbf *DataBaseOfFacts) *LearningMechanism {
	obj := new(LearningMechanism)
	obj.dbf = dbf
	return obj
}

func (obj *LearningMechanism) learn(record *Record) {
	for i := range obj.dbf.answeredIndexes {
		if obj.dbf.answeredIndexes[i] == IsAnswered {
			record.statistics[i].reduce()
			record.statistics[i].data[obj.dbf.answers[i]] += LearningIncrease
		}
	}
}
