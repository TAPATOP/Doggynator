package source

import (
	"math"
	"strconv"
	"strings"
)

const StatisticSize = 3
const ReductionFactor = 2
const DataMaximumBeforeReduce = 100

type Statistic struct {
	data [StatisticSize]int // {positive, negative, irrelevant}
}

func StatisticConstructor(input [StatisticSize]int) *Statistic {
	return &Statistic{data: input}
}

func RawStatisticConstructor(rawData string) (*Statistic, error) {
	splitString := strings.Split(rawData, " ")
	readyData := [StatisticSize]int{}
	for index, elem := range splitString {
		newNumber, err := strconv.Atoi(elem)
		if err != nil {
			continue
		}
		readyData[index] = newNumber
	}
	returnObj := StatisticConstructor(readyData)
	return returnObj, nil
}

func EmptyStatisticConstructor() *Statistic {
	return StatisticConstructor([StatisticSize]int{0, 0, 0})
}

func (obj *Statistic) getProbability(index int) float64 {
	totalCount := (float64)(obj.getTotalCountOf())
	if totalCount == 0 {
		totalCount = 1
	}
	return (float64)(obj.data[index]) / totalCount
}

func (obj *Statistic) getTotalCountOf() int {
	return obj.getPositive() + obj.getNegative() + obj.getUnknown()
}

func (obj *Statistic) getPositive() int {
	return obj.data[0]
}

func (obj *Statistic) getNegative() int {
	return obj.data[1]
}

func (obj *Statistic) getUnknown() int {
	return obj.data[2]
}

func (obj *Statistic) setPositive(val int) {
	obj.data[0] = val
}

func (obj *Statistic) setNegative(val int) {
	obj.data[1] = val
}

func (obj *Statistic) setUnknown(val int) {
	obj.data[2] = val
}

func (obj *Statistic) sumWith(stat *Statistic) {
	obj.setPositive(obj.getPositive() + stat.getPositive())
	obj.setNegative(obj.getNegative() + stat.getNegative())
	obj.setUnknown(obj.getUnknown() + stat.getUnknown())
}

func (obj *Statistic) reduce() {
	if obj.getTotalCountOf() > DataMaximumBeforeReduce {
		for i := range obj.data {
			obj.data[i] = obj.data[i] / ReductionFactor
		}
	}
}

func (stat *Statistic) ToString() (output string) {
	output = strconv.Itoa(stat.data[0])
	for _, elem := range stat.data[1:] {
		output += " " + strconv.Itoa(elem)
	}
	return
}

func (obj *Statistic) mostProbableAnswer() Response {
	mostProbableIndex := 0
	for i := 1; i < len(obj.data); i++ {
		if obj.data[mostProbableIndex] < obj.data[i] {
			mostProbableIndex = i
		}
	}
	return Response(mostProbableIndex)
}

func (obj *Statistic) entropy() float64 {
	positiveProbability := (float64)(obj.getPositive()) / (float64)(obj.getTotalCountOf())
	negativeProbability := (float64)(obj.getNegative()) / (float64)(obj.getTotalCountOf())
	unknownProbability := (float64)(obj.getUnknown()) / (float64)(obj.getTotalCountOf())

	return -positiveProbability*math.Log2(positiveProbability) -
		negativeProbability*math.Log2(negativeProbability) -
		unknownProbability*math.Log2(unknownProbability)
}

//TODO::
// Uncapitalize methods
// Do some renaming of the functions
