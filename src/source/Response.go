package source

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

func toResponse(forConverting string) (value Response) {
	switch forConverting {
	case "yes", "y":
		return Response(Yes)
	case "no", "n":
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
