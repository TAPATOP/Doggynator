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
	stringifiedStats := rec.statistics[0].ToString() + "\n"
	for _, elem := range rec.statistics[1:] {
		stringifiedStats = stringifiedStats + elem.ToString() + "\n"
	}
	return string(rec.name + "\n" + stringifiedStats)
}
