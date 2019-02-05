package source

import (
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
			return nil, err
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
	return obj.data[0] + obj.data[1] + obj.data[2]
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

func (obj *Statistic) mostProbableAnswerToAttribute() Response {
	mostProbableIndex := 0
	for i := range obj.data[1:] {
		if obj.data[mostProbableIndex] < obj.data[i] {
			mostProbableIndex = i
		}
	}
	return Response(mostProbableIndex)
}

//TODO::
// Uncapitalize methods
// Do some renaming of the functions
