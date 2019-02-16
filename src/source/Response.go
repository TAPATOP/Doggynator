package source

import (
	"strings"
)

type Response int

const (
	Yes Response = iota
	No
	DontKnowOrIrrelevant
	ProbablyYes
	ProbablyNo
	IncorrectResponse
)

func (resp Response) Integer() int {
	return [...]int{0, 1, 2, 3, 4}[resp]
}

func stringToResponse(forConverting string) (value Response) {
	switch strings.ToLower(forConverting) {
	case "yes", "y", "yeah", "yup", "ya", "ye", "yea":
		return Response(Yes)
	case "no", "n", "nah", "nope", "nein", "nay":
		return Response(No)
	case "irrelevant", "don't know", "no idea", "irr", "dk":
		return Response(DontKnowOrIrrelevant)
	case "probably", "p", "prob":
		return Response(ProbablyYes)
	case "probably not", "pn", "prob no":
		return Response(ProbablyNo)
	default:
		return Response(IncorrectResponse)
	}
}

func (resp Response) toString() string {
	switch resp {
	case Response(Yes):
		return "yes"
	case Response(No):
		return "no"
	case Response(DontKnowOrIrrelevant):
		return "don't know/ irrelevant"
	default:
		return "o fug"
	}
}
