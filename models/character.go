package models

import (
	"fmt"
	"math"
	"strings"
)

// ------------------------
// Ability Scores
// ------------------------
type AbilityScores struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
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

	mod := int(math.Floor(float64(score-10) / 2))
	return mod
}

// ------------------------
// Equipment
// ------------------------
type Weapon struct {
	Name      string `json:"name"`
	Category  string `json:"category,omitempty"`
	Range     string `json:"range,omitempty"`
	TwoHanded bool   `json:"two_handed,omitempty"`
}

type Armor struct {
	Name        string `json:"name"`
	ArmorClass  int    `json:"armor_class,omitempty"`
	DexBonus    bool   `json:"dex_bonus,omitempty"`
	MaxDexBonus int    `json:"max_dex_bonus,omitempty"`
}

type Shield struct {
	Name       string `json:"name"`
	ArmorClass int    `json:"armor_class"`
}

type Equipment struct {
	MainHand *Weapon `json:"main_hand,omitempty"`
	OffHand  *Weapon `json:"off_hand,omitempty"`
	Armor    *Armor  `json:"armor,omitempty"`
	Shield   *Shield `json:"shield,omitempty"`
}

// ------------------------
// Spells
// ------------------------
type Spell struct {
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Prepared bool   `json:"prepared"`
}

// ------------------------
// Character
// ------------------------
type Character struct {
	ID                 int            `json:"id"`
	Name               string         `json:"name"`
	PlayerName         string         `json:"player_name,omitempty"`
	Race               string         `json:"race"`
	Class              string         `json:"class"`
	Level              int            `json:"level"`
	Background         string         `json:"background"`
	Alignment          string         `json:"alignment,omitempty"`
	ProficiencyBonus   int            `json:"proficiency_bonus"`
	Abilities          AbilityScores  `json:"abilities"`
	SkillProficiencies []string       `json:"skill_proficiencies"`
	Skills             map[string]int `json:"skills"`

	Equipment Equipment `json:"equipment"`

	Spells     []Spell     `json:"spells,omitempty"`
	SpellSlots map[int]int `json:"spell_slots,omitempty"`

	ArmorClass        int `json:"armor_class"`
	Initiative        int `json:"initiative"`
	PassivePerception int `json:"passive_perception"`

	SpellcastingAbility string `json:"spellcasting_ability,omitempty"`
	SpellSaveDC         int    `json:"spell_save_dc,omitempty"`
	SpellAttackBonus    int    `json:"spell_attack_bonus,omitempty"`

	CanPrepareSpells bool `json:"can_prepare_spells"` // Nieuw veld

	ExperiencePoints   int    `json:"experience_points,omitempty"`
	Speed              int    `json:"speed,omitempty"`
	MaxHitPoints       int    `json:"max_hit_points,omitempty"`
	CurrentHitPoints   int    `json:"current_hit_points,omitempty"`
	TemporaryHitPoints int    `json:"temporary_hit_points,omitempty"`
	HitDiceTotal       string `json:"hit_dice_total,omitempty"`
	HitDiceRemaining   string `json:"hit_dice_remaining,omitempty"`
	DeathSaveSuccesses int    `json:"death_save_successes,omitempty"`
	DeathSaveFailures  int    `json:"death_save_failures,omitempty"`

	CopperPieces   int    `json:"copper_pieces,omitempty"`
	SilverPieces   int    `json:"silver_pieces,omitempty"`
	ElectrumPieces int    `json:"electrum_pieces,omitempty"`
	GoldPieces     int    `json:"gold_pieces,omitempty"`
	PlatinumPieces int    `json:"platinum_pieces,omitempty"`
	EquipmentText  string `json:"equipment_text,omitempty"`
	Personality    string `json:"personality,omitempty"`
	Ideals         string `json:"ideals,omitempty"`
	Bonds          string `json:"bonds,omitempty"`
	Flaws          string `json:"flaws,omitempty"`
	Features       string `json:"features,omitempty"`
}

// ------------------------
// Constants
// ------------------------
var StandardArray = []int{15, 14, 13, 12, 10, 8}

var RaceModifiers = map[string]map[string]int{
	"human":              {"Strength": 1, "Dexterity": 1, "Constitution": 1, "Intelligence": 1, "Wisdom": 1, "Charisma": 1},
	"elf":                {"Dexterity": 2},
	"dwarf":              {"Constitution": 2},
	"lightfoot halfling": {"Dexterity": 2, "Charisma": 1},
	"dragonborn":         {"Strength": 2, "Charisma": 1},
	"gnome":              {"Intelligence": 2},
	"half-elf":           {"Charisma": 2},
	"half-orc":           {"Strength": 2, "Constitution": 1},
	"tiefling":           {"Intelligence": 1, "Charisma": 2},
}

var ClassSkills = map[string][]string{
	"barbarian": {"Animal Handling", "Athletics", "Insight", "Religion", "Perception", "Survival"},
	"bard":      {"Deception", "History", "Investigation", "Persuasion", "Sleight of Hand"},
	"cleric":    {"History", "Insight", "Insight", "Religion"},
	"druid":     {"Arcana", "Animal Handling", "Insight", "Medicine", "Nature", "Perception", "Religion", "Survival"},
	"fighter":   {"Acrobatics", "Animal Handling", "Insight", "Religion", "Intimidation", "Perception", "Survival"},
	"monk":      {"Acrobatics", "Athletics", "History", "Insight", "Religion", "Stealth"},
	"paladin":   {"Athletics", "Insight","Insight", "Religion"},
	"ranger":    {"Animal Handling", "Athletics", "Insight", "Investigation", "Nature", "Perception", "Stealth", "Survival"},
	"rogue":     {"Acrobatics", "Athletics", "Deception", "Insight", "Insight", "Religion"},
	"sorcerer":  {"Arcana", "Deception", "Insight", "Intimidation", "Persuasion", "Religion"},
	"warlock":   {"Arcana", "Deception", "Insight", "Religion"},
	"wizard":    {"Arcana", "History", "Insight", "Investigation", "Medicine", "Religion"},
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
	"bard":     "Charisma",
	"cleric":   "Wisdom",
	"druid":    "Wisdom",
	"paladin":  "Charisma",
	"ranger":   "Wisdom",
	"sorcerer": "Charisma",
	"warlock":  "Charisma",
	"wizard":   "Intelligence",
}

// ------------------------
// Constructor
// ------------------------
func NewCharacter(id int, name, race, class, background string, level int, abilityScores []int, skillChoices []string) *Character {
	raceKey := strings.ToLower(race)
	classKey := strings.ToLower(class)

	var abilities AbilityScores
	mod := RaceModifiers[raceKey]

	if len(abilityScores) == 6 {
		abilities = AbilityScores{
			Strength:     abilityScores[0] + mod["Strength"],
			Dexterity:    abilityScores[1] + mod["Dexterity"],
			Constitution: abilityScores[2] + mod["Constitution"],
			Intelligence: abilityScores[3] + mod["Intelligence"],
			Wisdom:       abilityScores[4] + mod["Wisdom"],
			Charisma:     abilityScores[5] + mod["Charisma"],
		}
	} else {
		abilities = AssignAbilities(mod)
	}

	character := &Character{
		ID:                 id,
		Name:               name,
		Race:               raceKey,
		Class:              classKey,
		Level:              level,
		Background:         background,
		ProficiencyBonus:   CalculateProfBonus(level),
		Abilities:          abilities,
		SkillProficiencies: skillChoices,
		Skills:             make(map[string]int),
		Equipment:          Equipment{},
		ArmorClass:         10,
		Speed:              30,
		MaxHitPoints:       10,
		CurrentHitPoints:   10,
		CanPrepareSpells:   isPreparedCaster(classKey),
	}

	character.CalculateAllSkills()
	character.CalculateCombatStats()
	character.SetupSpellcasting()

	return character
}

// ------------------------
// Helper voor prepared casters
// ------------------------
func isPreparedCaster(className string) bool {
	preparedClasses := map[string]bool{
		"cleric":  true,
		"druid":   true,
		"paladin": true,
		"wizard":  true,
	}
	return preparedClasses[className]
}

// ------------------------
// Helpers (berekeningen voor combat, skills, AC, spellcasting)
// ------------------------
func AssignAbilities(mod map[string]int) AbilityScores {
	abilities := AbilityScores{}
	order := []string{"Strength", "Dexterity", "Constitution", "Intelligence", "Wisdom", "Charisma"}

	for i, name := range order {
		score := StandardArray[i] + mod[name]
		switch name {
		case "Strength":
			abilities.Strength = score
		case "Dexterity":
			abilities.Dexterity = score
		case "Constitution":
			abilities.Constitution = score
		case "Intelligence":
			abilities.Intelligence = score
		case "Wisdom":
			abilities.Wisdom = score
		case "Charisma":
			abilities.Charisma = score
		}
	}

	return abilities
}

func CalculateProfBonus(level int) int {
	return 2 + (level-1)/4
}

func (c *Character) UpdateLevel(newLevel int) {
	c.Level = newLevel
	c.ProficiencyBonus = CalculateProfBonus(newLevel)
	c.CalculateAllSkills()
	c.CalculateCombatStats()
	c.UpdateSpellSlots()
}

func (c *Character) CalculateAllSkills() {
	c.Skills = make(map[string]int)
	for skill, ability := range SkillAbilities {
		mod := c.Abilities.Modifier(ability)
		if contains(c.SkillProficiencies, skill) {
			mod += c.ProficiencyBonus
		}
		c.Skills[skill] = mod
	}
}

func (c *Character) CalculateCombatStats() {
	c.Initiative = c.Abilities.Modifier("Dexterity")
	c.PassivePerception = 10 + c.Abilities.Modifier("Wisdom")
	if p, ok := c.Skills["Perception"]; ok {
		c.PassivePerception = 10 + p
	}
	c.CalculateArmorClass()
}

func (c *Character) CalculateArmorClass() {
	ac := 10 + c.Abilities.Modifier("Dexterity")
	if c.Equipment.Armor != nil {
		ac = c.Equipment.Armor.ArmorClass
		if c.Equipment.Armor.DexBonus {
			dexMod := c.Abilities.Modifier("Dexterity")
			if c.Equipment.Armor.MaxDexBonus > 0 && dexMod > c.Equipment.Armor.MaxDexBonus {
				dexMod = c.Equipment.Armor.MaxDexBonus
			}
			ac += dexMod
		}
	}
	if c.Equipment.Shield != nil {
		ac += c.Equipment.Shield.ArmorClass
	}
	c.ArmorClass = ac
}

func (c *Character) SetupSpellcasting() {
	if ability, ok := SpellcastingClasses[c.Class]; ok {
		c.SpellcastingAbility = ability
		mod := c.Abilities.Modifier(ability)
		c.SpellSaveDC = 8 + c.ProficiencyBonus + mod
		c.SpellAttackBonus = c.ProficiencyBonus + mod
		c.UpdateSpellSlots()
	}
}

func (c *Character) UpdateSpellSlots() {
	classKey := c.Class
	if _, ok := SpellcastingClasses[classKey]; !ok {
		return
	}
	c.SpellSlots = make(map[int]int)
	// Voeg hier je SpellSlotsByLevel implementatie toe zoals in je originele file
}

// ------------------------
// Utility
// ------------------------
func GetAvailableSkills(className string) []string {
	if skills, ok := ClassSkills[className]; ok {
		return skills
	}
	return []string{}
}

func FormatModifier(mod int) string {
	if mod >= 0 {
		return fmt.Sprintf("+%d", mod)
	}
	return fmt.Sprintf("%d", mod)
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
