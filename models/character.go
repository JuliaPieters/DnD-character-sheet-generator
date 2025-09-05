package models

// AbilityScores bevat de zes ability scores van een character
type AbilityScores struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

// Character is het hoofdmodel van een D&D character
type Character struct {
	Name       string         `json:"name"`
	Race       string         `json:"race"`
	Class      string         `json:"class"`
	Level      int            `json:"level"`
	Background string         `json:"background"`
	ProfBonus  int            `json:"prof_bonus"`
	Abilities  AbilityScores  `json:"abilities"`
	Skills     map[string]int `json:"skills"`
}

// StandardArray voor het toewijzen van ability scores
var StandardArray = []int{15, 14, 13, 12, 10, 8}

// Nieuwe character aanmaken
func NewCharacter(name, race, class, background string, level int, str, dex, con, intt, wis, cha int) *Character {
	char := &Character{
		Name:       name,
		Race:       race,
		Class:      class,
		Level:      level,
		Background: background,
		ProfBonus:  CalculateProfBonus(level),
		Abilities: AbilityScores{
			Strength:     str,
			Dexterity:    dex,
			Constitution: con,
			Intelligence: intt,
			Wisdom:       wis,
			Charisma:     cha,
		},
		Skills: make(map[string]int),
	}
	return char
}

// Bereken de proficiency bonus op basis van level
func CalculateProfBonus(level int) int {
	return 2 + (level-1)/4
}

// Update level en pas proficiency bonus aan
func (c *Character) UpdateLevel(newLevel int) {
	c.Level = newLevel
	c.ProfBonus = CalculateProfBonus(newLevel)
}

// Bereken ability modifier voor een gegeven ability score
func (a AbilityScores) Modifier(score string) int {
	switch score {
	case "STR":
		return (a.Strength - 10) / 2
	case "DEX":
		return (a.Dexterity - 10) / 2
	case "CON":
		return (a.Constitution - 10) / 2
	case "INT":
		return (a.Intelligence - 10) / 2
	case "WIS":
		return (a.Wisdom - 10) / 2
	case "CHA":
		return (a.Charisma - 10) / 2
	}
	return 0
}

// Bereken skill modifier: ability modifier + profBonus (als character proficiency heeft)
func CalculateSkillModifier(abilityScore int, proficient bool, profBonus int) int {
	mod := (abilityScore - 10) / 2
	if proficient {
		mod += profBonus
	}
	return mod
}

// Voeg skills toe en bereken hun modifiers
func (c *Character) AddSkills(skills []string) {
	for _, skill := range skills {
		var abilityScore int
		switch skill {
		case "Athletics":
			abilityScore = c.Abilities.Strength
		case "Acrobatics", "Sleight of Hand", "Stealth":
			abilityScore = c.Abilities.Dexterity
		case "Arcana", "History", "Investigation", "Nature", "Religion":
			abilityScore = c.Abilities.Intelligence
		case "Animal Handling", "Insight", "Medicine", "Perception", "Survival":
			abilityScore = c.Abilities.Wisdom
		case "Deception", "Intimidation", "Performance", "Persuasion":
			abilityScore = c.Abilities.Charisma
		}
		c.Skills[skill] = CalculateSkillModifier(abilityScore, true, c.ProfBonus)
	}
}
