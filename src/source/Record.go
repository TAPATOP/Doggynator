package source

type Record struct {
	name       string
	statistics []Statistic
}

func RecordConstructor(name string, stats []Statistic) *Record {
	return &Record{name: name, statistics: stats}
}

func (obj *Record) AddField() {
	obj.statistics = append(obj.statistics, *NullStatisticConstructor())
}

func (rec *Record) ToString() string {
	stringifiedStats := "\n"
	for _, elem := range rec.statistics {
		stringifiedStats = stringifiedStats + elem.ToString() + "\n"
	}
	return string(rec.name + stringifiedStats)
}
