package source

import "strconv"

type Statistic struct {
	positive   int
	negative   int
	irrelevant int
}

//func StatisticConstructor(positive, negative, irrelevant int) *Statistic {
//	return &Statistic{positive:positive, negative:negative, irrelevant:irrelevant}
//}

func (stat *Statistic) ToString() string {
	return string(strconv.Itoa(stat.positive) + " " + strconv.Itoa(stat.negative) + " " + strconv.Itoa(stat.irrelevant))
}
