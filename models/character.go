package models

import "fmt"

// AbilityScores bevat de zes ability scores van een character
type AbilityScores struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

// Equipment structures volgens requirements
type Weapon struct {
	Name      string `json:"name"`
	Category  string `json:"category,omitempty"`   // Van API
	Range     string `json:"range,omitempty"`      // Van API
	TwoHanded bool   `json:"two_handed,omitempty"` // Van API
}

type Armor struct {
	Name        string `json:"name"`
	ArmorClass  int    `json:"armor_class,omitempty"`  // Van API
	DexBonus    bool   `json:"dex_bonus,omitempty"`    // Van API
	MaxDexBonus int    `json:"max_dex_bonus,omitempty"` // Van API
}

type Shield struct {
	Name       string `json:"name"`
	ArmorClass int    `json:"armor_class"` // Meestal +2
}

type Equipment struct {
	Weapons []Weapon `json:"weapons"`
	Armor   *Armor   `json:"armor,omitempty"`
	Shield  *Shield  `json:"shield,omitempty"`
}

// Spellcasting volgens requirements
type Spell struct {
	Name   string `json:"name"`
	Level  int    `json:"level"`
	School string `json:"school,omitempty"` // Van API
	Range  string `json:"range,omitempty"`  // Van API
}

// Character is het volledige model van een D&D character
type Character struct {
	ID                   int                `json:"id"`                     // Nieuw: voor CRU(D)
	Name                 string             `json:"name"`
	PlayerName           string             `json:"player_name,omitempty"`  // Niet in requirements
	Race                 string             `json:"race"`
	Class                string             `json:"class"`
	Level                int                `json:"level"`
	Background           string             `json:"background"`
	Alignment            string             `json:"alignment,omitempty"`    // Niet in requirements
	ProficiencyBonus     int                `json:"proficiency_bonus"`
	Abilities            AbilityScores      `json:"abilities"`
	SkillProficiencies   []string           `json:"skill_proficiencies"`   // Nieuw: welke skills character heeft
	Skills               map[string]int     `json:"skills"`                // Berekende skill modifiers

	// Equipment volgens requirements
	Equipment            Equipment          `json:"equipment"`

	// Spellcasting volgens requirements (alleen voor casters)
	Spells               []Spell            `json:"spells,omitempty"`
	SpellSlots           map[int]int        `json:"spell_slots,omitempty"` // level -> aantal slots

	// Combat stats volgens requirements
	ArmorClass           int                `json:"armor_class"`
	Initiative           int                `json:"initiative"`
	PassivePerception    int                `json:"passive_perception"`    // "passive wisdom"
	
	// Spellcasting stats (alleen voor casters)
	SpellcastingAbility  string             `json:"spellcasting_ability,omitempty"`
	SpellSaveDC          int                `json:"spell_save_dc,omitempty"`
	SpellAttackBonus     int                `json:"spell_attack_bonus,omitempty"`

	// Extra velden die je had (niet in requirements, maar kunnen blijven)
	ExperiencePoints     int                `json:"experience_points,omitempty"`
	Speed                int                `json:"speed,omitempty"`
	MaxHitPoints         int                `json:"max_hit_points,omitempty"`
	CurrentHitPoints     int                `json:"current_hit_points,omitempty"`
	TemporaryHitPoints   int                `json:"temporary_hit_points,omitempty"`
	HitDiceTotal         string             `json:"hit_dice_total,omitempty"`
	HitDiceRemaining     string             `json:"hit_dice_remaining,omitempty"`
	DeathSaveSuccesses   int                `json:"death_save_successes,omitempty"`
	DeathSaveFailures    int                `json:"death_save_failures,omitempty"`
	
	// Geld en flavor text
	CopperPieces         int                `json:"copper_pieces,omitempty"`
	SilverPieces         int                `json:"silver_pieces,omitempty"`
	ElectrumPieces       int                `json:"electrum_pieces,omitempty"`
	GoldPieces           int                `json:"gold_pieces,omitempty"`
	PlatinumPieces       int                `json:"platinum_pieces,omitempty"`
	EquipmentText        string             `json:"equipment_text,omitempty"`
	Personality          string             `json:"personality,omitempty"`
	Ideals               string             `json:"ideals,omitempty"`
	Bonds                string             `json:"bonds,omitempty"`
	Flaws                string             `json:"flaws,omitempty"`
	Features             string             `json:"features,omitempty"`
}

// Game data volgens requirements
var StandardArray = []int{15, 14, 13, 12, 10, 8}

var RaceModifiers = map[string]map[string]int{
	"Human":      {"Strength": 1, "Dexterity": 1, "Constitution": 1, "Intelligence": 1, "Wisdom": 1, "Charisma": 1},
	"Elf":        {"Dexterity": 2},
	"Dwarf":      {"Constitution": 2},
	"Halfling":   {"Dexterity": 2},
	"Dragonborn": {"Strength": 2, "Charisma": 1},
	"Gnome":      {"Intelligence": 2},
	"Half-Elf":   {"Charisma": 2},
	"Half-Orc":   {"Strength": 2, "Constitution": 1},
	"Tiefling":   {"Intelligence": 1, "Charisma": 2},
}

var ClassSkills = map[string][]string{
	"Barbarian": {"Animal Handling", "Athletics", "Intimidation", "Nature", "Perception", "Survival"},
	"Bard":      {"Deception", "History", "Investigation", "Persuasion", "Sleight of Hand"},
	"Cleric":    {"History", "Insight", "Medicine", "Persuasion", "Religion"},
	"Druid":     {"Arcana", "Animal Handling", "Insight", "Medicine", "Nature", "Perception", "Religion", "Survival"},
	"Fighter":   {"Acrobatics", "Animal Handling", "Athletics", "History", "Insight", "Intimidation", "Perception", "Survival"},
	"Monk":      {"Acrobatics", "Athletics", "History", "Insight", "Religion", "Stealth"},
	"Paladin":   {"Athletics", "Insight", "Intimidation", "Medicine", "Persuasion", "Religion"},
	"Ranger":    {"Animal Handling", "Athletics", "Insight", "Investigation", "Nature", "Perception", "Stealth", "Survival"},
	"Rogue":     {"Acrobatics", "Athletics", "Deception", "Insight", "Intimidation", "Investigation", "Perception", "Performance", "Persuasion", "Sleight of Hand", "Stealth"},
	"Sorcerer":  {"Arcana", "Deception", "Insight", "Intimidation", "Persuasion", "Religion"},
	"Warlock":   {"Arcana", "Deception", "History", "Intimidation", "Investigation", "Nature", "Religion"},
	"Wizard":    {"Arcana", "History", "Insight", "Investigation", "Medicine", "Religion"},
}

var SkillAbilities = map[string]string{
	"Acrobatics":      "Dexterity",
	"Animal Handling": "Wisdom",
	"Arcana":          "Intelligence",
	"Athletics":       "Strength",
	"Deception":       "Charisma",
	"History":         "Intelligence",
	"Insight":         "Wisdom",
	"Intimidation":    "Charisma",
	"Investigation":   "Intelligence",
	"Medicine":        "Wisdom",
	"Nature":          "Intelligence",
	"Perception":      "Wisdom",
	"Performance":     "Charisma",
	"Persuasion":      "Charisma",
	"Religion":        "Intelligence",
	"Sleight of Hand": "Dexterity",
	"Stealth":         "Dexterity",
	"Survival":        "Wisdom",
}

var SpellcastingClasses = map[string]string{
	"Bard":     "Charisma",
	"Cleric":   "Wisdom", 
	"Druid":    "Wisdom",
	"Paladin":  "Charisma",
	"Ranger":   "Wisdom",
	"Sorcerer": "Charisma",
	"Warlock":  "Charisma",
	"Wizard":   "Intelligence",
}

// Spell slots per level per class (vereenvoudigd)
var SpellSlotsByLevel = map[string]map[int]map[int]int{
	"Bard": {
		1: {1: 2},
		2: {1: 3},
		3: {1: 4, 2: 2},
		4: {1: 4, 2: 3},
		5: {1: 4, 2: 3, 3: 2},
	},
	"Cleric": {
		1: {1: 2},
		2: {1: 3},
		3: {1: 4, 2: 2},
		4: {1: 4, 2: 3},
		5: {1: 4, 2: 3, 3: 2},
	},
	"Wizard": {
		1: {1: 2},
		2: {1: 3},
		3: {1: 4, 2: 2},
		4: {1: 4, 2: 3},
		5: {1: 4, 2: 3, 3: 2},
	},
}

// Nieuwe character aanmaken MET Standard Array + Race modifiers
func NewCharacter(id int, name, race, class, background string, level int, abilityAssignment []string, skillChoices []string) *Character {
	// Standard Array toewijzen volgens user keuze
	abilities := AssignAbilities(abilityAssignment, race)
	
	character := &Character{
		ID:               id,
		Name:             name,
		Race:             race,
		Class:            class,
		Level:            level,
		Background:       background,
		ProficiencyBonus: CalculateProfBonus(level),
		Abilities:        abilities,
		SkillProficiencies: skillChoices,
		Skills:           make(map[string]int),
		Equipment:        Equipment{},
		ArmorClass:       10, // Base AC
		Speed:            30,
		MaxHitPoints:     10,
		CurrentHitPoints: 10,
	}

	// Skills berekenen
	character.CalculateAllSkills()
	
	// Combat stats berekenen
	character.CalculateCombatStats()
	
	// Spellcasting setup als van toepassing
	character.SetupSpellcasting()

	return character
}

// Assign abilities volgens Standard Array + Race modifiers
func AssignAbilities(assignment []string, race string) AbilityScores {
	if len(assignment) != 6 {
		// Fallback: default assignment
		assignment = []string{"Strength", "Dexterity", "Constitution", "Intelligence", "Wisdom", "Charisma"}
	}

	abilities := AbilityScores{}
	
	for i, abilityName := range assignment {
		baseScore := StandardArray[i]
		
		// Race modifier toepassen
		if raceModifiers, exists := RaceModifiers[race]; exists {
			if modifier, hasModifier := raceModifiers[abilityName]; hasModifier {
				baseScore += modifier
			}
		}

		switch abilityName {
		case "Strength":
			abilities.Strength = baseScore
		case "Dexterity":
			abilities.Dexterity = baseScore
		case "Constitution":
			abilities.Constitution = baseScore
		case "Intelligence":
			abilities.Intelligence = baseScore
		case "Wisdom":
			abilities.Wisdom = baseScore
		case "Charisma":
			abilities.Charisma = baseScore
		}
	}

	return abilities
}

// Bereken de proficiency bonus op basis van level
func CalculateProfBonus(level int) int {
	return 2 + (level-1)/4
}

// Update level en pas proficiency bonus aan
func (c *Character) UpdateLevel(newLevel int) {
	c.Level = newLevel
	c.ProficiencyBonus = CalculateProfBonus(newLevel)
	c.CalculateAllSkills()
	c.CalculateCombatStats()
	c.UpdateSpellSlots()
}

// Bereken ability modifier voor een gegeven ability score
func (abilities AbilityScores) Modifier(abilityName string) int {
	var score int
	switch abilityName {
	case "Strength":
		score = abilities.Strength
	case "Dexterity":
		score = abilities.Dexterity
	case "Constitution":
		score = abilities.Constitution
	case "Intelligence":
		score = abilities.Intelligence
	case "Wisdom":
		score = abilities.Wisdom
	case "Charisma":
		score = abilities.Charisma
	default:
		return 0
	}
	return (score - 10) / 2
}

// Bereken alle skill modifiers
func (c *Character) CalculateAllSkills() {
	c.Skills = make(map[string]int)
	
	for skillName, abilityName := range SkillAbilities {
		isProficient := false
		for _, profSkill := range c.SkillProficiencies {
			if profSkill == skillName {
				isProficient = true
				break
			}
		}
		
		abilityModifier := c.Abilities.Modifier(abilityName)
		skillModifier := abilityModifier
		if isProficient {
			skillModifier += c.ProficiencyBonus
		}
		
		c.Skills[skillName] = skillModifier
	}
}

// Bereken combat statistics
func (c *Character) CalculateCombatStats() {
	// Initiative = Dex modifier
	c.Initiative = c.Abilities.Modifier("Dexterity")
	
	// Passive Perception = 10 + Perception skill
	if perceptionMod, exists := c.Skills["Perception"]; exists {
		c.PassivePerception = 10 + perceptionMod
	} else {
		c.PassivePerception = 10 + c.Abilities.Modifier("Wisdom")
	}
	
	// Armor Class wordt later berekend op basis van armor
	c.CalculateArmorClass()
}

// Bereken Armor Class op basis van equipment
func (c *Character) CalculateArmorClass() {
	baseAC := 10 + c.Abilities.Modifier("Dexterity") // Unarmored
	
	if c.Equipment.Armor != nil {
		baseAC = c.Equipment.Armor.ArmorClass
		if c.Equipment.Armor.DexBonus {
			dexMod := c.Abilities.Modifier("Dexterity")
			if c.Equipment.Armor.MaxDexBonus > 0 && dexMod > c.Equipment.Armor.MaxDexBonus {
				dexMod = c.Equipment.Armor.MaxDexBonus
			}
			baseAC += dexMod
		}
	}
	
	if c.Equipment.Shield != nil {
		baseAC += c.Equipment.Shield.ArmorClass
	}
	
	c.ArmorClass = baseAC
}

// Setup spellcasting als character een caster is
func (c *Character) SetupSpellcasting() {
	if spellAbility, isCaster := SpellcastingClasses[c.Class]; isCaster {
		c.SpellcastingAbility = spellAbility
		
		// Spell save DC = 8 + prof bonus + ability modifier
		abilityMod := c.Abilities.Modifier(spellAbility)
		c.SpellSaveDC = 8 + c.ProficiencyBonus + abilityMod
		c.SpellAttackBonus = c.ProficiencyBonus + abilityMod
		
		c.UpdateSpellSlots()
	}
}

// Update spell slots op basis van level
func (c *Character) UpdateSpellSlots() {
	if _, isCaster := SpellcastingClasses[c.Class]; !isCaster {
		return
	}
	
	if classSlots, exists := SpellSlotsByLevel[c.Class]; exists {
		if levelSlots, hasLevel := classSlots[c.Level]; hasLevel {
			c.SpellSlots = make(map[int]int)
			for level, slots := range levelSlots {
				c.SpellSlots[level] = slots
			}
		}
	}
}

// Helper: Get available skills for a class
func GetAvailableSkills(className string) []string {
	if skills, exists := ClassSkills[className]; exists {
		return skills
	}
	return []string{}
}

// Helper: Format modifier voor display
func FormatModifier(modifier int) string {
	if modifier >= 0 {
		return fmt.Sprintf("+%d", modifier)
	}
	return fmt.Sprintf("%d", modifier)
}