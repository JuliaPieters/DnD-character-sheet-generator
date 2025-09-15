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
func NewCharacter(name string, race string, class string, background string, level int, strength int, dexterity int, constitution int, intelligence int, wisdom int, charisma int) *Character {
	newCharacter := &Character{
		Name:       name,
		Race:       race,
		Class:      class,
		Level:      level,
		Background: background,
		ProfBonus:  CalculateProfBonus(level),
		Abilities: AbilityScores{
			Strength:     strength,
			Dexterity:    dexterity,
			Constitution: constitution,
			Intelligence: intelligence,
			Wisdom:       wisdom,
			Charisma:     charisma,
		},
		Skills: make(map[string]int),
	}
	return newCharacter
}

// Bereken de proficiency bonus op basis van level
func CalculateProfBonus(level int) int {
	return 2 + (level-1)/4
}

// Update level en pas proficiency bonus aan
func (currentCharacter *Character) UpdateLevel(newLevel int) {
	currentCharacter.Level = newLevel
	currentCharacter.ProfBonus = CalculateProfBonus(newLevel)
}

// Bereken ability modifier voor een gegeven ability score
func (abilityScores AbilityScores) Modifier(score string) int {
	switch score {
	case "STR":
		return (abilityScores.Strength - 10) / 2
	case "DEX":
		return (abilityScores.Dexterity - 10) / 2
	case "CON":
		return (abilityScores.Constitution - 10) / 2
	case "INT":
		return (abilityScores.Intelligence - 10) / 2
	case "WIS":
		return (abilityScores.Wisdom - 10) / 2
	case "CHA":
		return (abilityScores.Charisma - 10) / 2
	}
	return 0
}

// Bereken skill modifier: ability modifier + profBonus (als character proficiency heeft)
func CalculateSkillModifier(abilityScore int, isProficient bool, proficiencyBonus int) int {
	modifier := (abilityScore - 10) / 2
	if isProficient {
		modifier += proficiencyBonus
	}
	return modifier
}

// Voeg skills toe en bereken hun modifiers
func (currentCharacter *Character) AddSkills(skillNames []string) {
	for _, skillName := range skillNames {
		var abilityScore int
		switch skillName {
		case "Athletics":
			abilityScore = currentCharacter.Abilities.Strength
		case "Acrobatics", "Sleight of Hand", "Stealth":
			abilityScore = currentCharacter.Abilities.Dexterity
		case "Arcana", "History", "Investigation", "Nature", "Religion":
			abilityScore = currentCharacter.Abilities.Intelligence
		case "Animal Handling", "Insight", "Medicine", "Perception", "Survival":
			abilityScore = currentCharacter.Abilities.Wisdom
		case "Deception", "Intimidation", "Performance", "Persuasion":
			abilityScore = currentCharacter.Abilities.Charisma
		}
		currentCharacter.Skills[skillName] = CalculateSkillModifier(abilityScore, true, currentCharacter.ProfBonus)
	}
}
