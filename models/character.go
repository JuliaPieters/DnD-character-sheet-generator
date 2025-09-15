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

// Attack bevat de informatie voor een aanval of spellcasting entry
type Attack struct {
	Name       string `json:"name"`
	AttackBonus string `json:"attack_bonus"`
	Damage     string `json:"damage"`
}

// Character is het volledige model van een D&D character
type Character struct {
	Name              string         `json:"name"`
	PlayerName        string         `json:"player_name"`
	Race              string         `json:"race"`
	Class             string         `json:"class"`
	Level             int            `json:"level"`
	Background        string         `json:"background"`
	Alignment         string         `json:"alignment"`
	ExperiencePoints  int            `json:"experience_points"`
	ProfBonus         int            `json:"prof_bonus"`
	Abilities         AbilityScores  `json:"abilities"`
	Skills            map[string]int `json:"skills"`

	// Combat
	ArmorClass        int    `json:"armor_class"`
	Initiative        int    `json:"initiative"`
	Speed             int    `json:"speed"`
	MaxHitPoints      int    `json:"max_hit_points"`
	CurrentHitPoints  int    `json:"current_hit_points"`
	TemporaryHitPoints int   `json:"temporary_hit_points"`
	HitDiceTotal      string `json:"hit_dice_total"`
	HitDiceRemaining  string `json:"hit_dice_remaining"`
	DeathSaveSuccesses int   `json:"death_save_successes"`
	DeathSaveFailures  int   `json:"death_save_failures"`

	// Attacks & Spellcasting
	Attacks           []Attack `json:"attacks"`

	// Equipment
	CopperPieces      int    `json:"copper_pieces"`
	SilverPieces      int    `json:"silver_pieces"`
	ElectrumPieces    int    `json:"electrum_pieces"`
	GoldPieces        int    `json:"gold_pieces"`
	PlatinumPieces    int    `json:"platinum_pieces"`
	EquipmentText     string `json:"equipment_text"`

	// Personality / Flavor
	Personality       string `json:"personality"`
	Ideals            string `json:"ideals"`
	Bonds             string `json:"bonds"`
	Flaws             string `json:"flaws"`
	Features          string `json:"features"`
}

// StandardArray voor het toewijzen van ability scores
var StandardArray = []int{15, 14, 13, 12, 10, 8}

// Nieuwe character aanmaken
func NewCharacter(name, playerName, race, class, background, alignment string, level int, strength, dexterity, constitution, intelligence, wisdom, charisma int) *Character {
	return &Character{
		Name:             name,
		PlayerName:       playerName,
		Race:             race,
		Class:            class,
		Level:            level,
		Background:       background,
		Alignment:        alignment,
		ExperiencePoints: 0,
		ProfBonus:        CalculateProfBonus(level),
		Abilities: AbilityScores{
			Strength:     strength,
			Dexterity:    dexterity,
			Constitution: constitution,
			Intelligence: intelligence,
			Wisdom:       wisdom,
			Charisma:     charisma,
		},
		Skills:           make(map[string]int),
		ArmorClass:       10,
		Initiative:       0,
		Speed:            30,
		MaxHitPoints:     10,
		CurrentHitPoints: 10,
		TemporaryHitPoints: 0,
		HitDiceTotal:     "1d8",
		HitDiceRemaining: "1d8",
		DeathSaveSuccesses: 0,
		DeathSaveFailures: 0,
		Attacks: []Attack{
			{"", "", ""},
			{"", "", ""},
			{"", "", ""},
		},
		CopperPieces:   0,
		SilverPieces:   0,
		ElectrumPieces: 0,
		GoldPieces:     0,
		PlatinumPieces: 0,
		EquipmentText:  "",
		Personality:    "",
		Ideals:         "",
		Bonds:          "",
		Flaws:          "",
		Features:       "",
	}
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
	case "Strength":
		return (abilityScores.Strength - 10) / 2
	case "Dexterity":
		return (abilityScores.Dexterity - 10) / 2
	case "Constitution":
		return (abilityScores.Constitution - 10) / 2
	case "Intelligence":
		return (abilityScores.Intelligence - 10) / 2
	case "Wisdom":
		return (abilityScores.Wisdom - 10) / 2
	case "Charisma":
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
