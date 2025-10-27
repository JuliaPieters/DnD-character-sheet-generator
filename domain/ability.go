package domain

import "math"

type AbilityScores struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

func (a AbilityScores) Modifier(name string) int {
	var score int
	switch name {
	case "Strength":
		score = a.Strength
	case "Dexterity":
		score = a.Dexterity
	case "Constitution":
		score = a.Constitution
	case "Intelligence":
		score = a.Intelligence
	case "Wisdom":
		score = a.Wisdom
	case "Charisma":
		score = a.Charisma
	default:
		return 0
	}
	return int(math.Floor(float64(score-10) / 2))
}
