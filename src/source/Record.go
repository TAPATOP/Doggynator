package source

const StartingRecordValue = 2

type Record struct {
	name       string
	statistics []Statistic
}

func RecordConstructor(name string, stats []Statistic) *Record {
	return &Record{name: name, statistics: stats}
}

func EmptyRecordConstructor(name string, numberOfAttributes int) *Record {
	var stats []Statistic
	for i := 0; i < numberOfAttributes; i++ {
		stats = append(stats, *EmptyStatisticConstructor())
	}
	return RecordConstructor(name, stats)
}

func (obj *Record) AddField() {
	obj.statistics = append(obj.statistics, *StatisticConstructor([StatisticSize]int{StartingRecordValue, StartingRecordValue, StartingRecordValue}))
}

func (rec *Record) ToString() string {
	stringifiedStats := rec.statistics[0].ToString() + "\n"
	for _, elem := range rec.statistics[1:] {
		stringifiedStats = stringifiedStats + elem.ToString() + "\n"
	}
	return string(rec.name + "\n" + stringifiedStats)
}
