package src

type Record struct {
	name       string
	statistics Statistic
}

type Statistic struct {
	positive   int
	negative   int
	irrelevant int
}
