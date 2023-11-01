package entity

type Question struct {
	ID              uint
	Title           string
	PossibleAnswer  []PossibleAnswer
	CorrectAnswerID uint
	Dificulity      QuestionDificulity
	Category        uint
}

type PossibleAnswer struct {
	ID     uint
	Text   string
	Choice PossibleAnswerChoice
}

type PossibleAnswerChoice uint8

func (p PossibleAnswerChoice) IsValid() bool {

	if p >= PossibleAnswerA && p <= PossibleAnswerD {
		return true
	}

	return false
}

const (
	PossibleAnswerA PossibleAnswerChoice = iota + 1
	PossibleAnswerB
	PossibleAnswerC
	PossibleAnswerD
)

type QuestionDificulity uint8

const (
	QuestionDificulityEasy QuestionDificulity = iota + 1
	QuestionDificulityMedium
	QuestionDificulityHard
)

func (q QuestionDificulity) IsValid() bool {
	if q >= QuestionDificulityEasy && q <= QuestionDificulityHard {
		return true
	}

	return false
}
