package source

import (
	"strconv"
	"strings"
)

const StatisticSize = 3

type Statistic struct {
	data [StatisticSize]int
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

func NullStatisticConstructor() *Statistic {
	return StatisticConstructor([StatisticSize]int{0, 0, 0})
}

func (stat *Statistic) ToString() (output string) {
	output = strconv.Itoa(stat.data[0])
	for _, elem := range stat.data[1:] {
		output += " " + strconv.Itoa(elem)
	}
	return
}
