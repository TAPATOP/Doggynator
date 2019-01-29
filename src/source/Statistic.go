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

func (stat *Statistic) ToString() (output string) {
	for _, elem := range stat.data {
		output += strconv.Itoa(elem) + " "
	}
	return
}
