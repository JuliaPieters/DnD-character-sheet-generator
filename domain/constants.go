package domain

// ------------------------
// Constants for Skills
// ------------------------
const (
	AnimalHandling = "Animal Handling"
	Athletics      = "Athletics"
	Insight        = "Insight"
	Religion       = "Religion"
	Acrobatics     = "Acrobatics"
	Deception      = "Deception"
	History        = "History"
	Investigation  = "Investigation"
	Persuasion     = "Persuasion"
	SleightOfHand  = "Sleight of Hand"
	Arcana         = "Arcana"
	Medicine       = "Medicine"
	Intimidation   = "Intimidation"
	Stealth        = "Stealth"
	Survival       = "Survival"
)

// ------------------------
// Standard Array
// ------------------------
var StandardArray = []int{15, 14, 13, 12, 10, 8}

// ------------------------
// Race Modifiers
// ------------------------
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
	"hill dwarf":         {"Constitution": 2, "Wisdom": 1},
	"half orc":           {"Strength": 2, "Constitution": 1},
}

// ------------------------
// Class Skills
// ------------------------
var ClassSkills = map[string][]string{
	"barbarian": {AnimalHandling, Athletics, Insight, Religion},
	"bard":      {Deception, History, Investigation, Persuasion, SleightOfHand},
	"cleric":    {History, Insight, Insight, Religion},
	"druid":     {Arcana, AnimalHandling, Insight, Medicine},
	"fighter":   {Acrobatics, AnimalHandling, Insight, Religion},
	"monk":      {Acrobatics, Athletics, Insight, Religion},
	"paladin":   {Athletics, Insight, Insight, Religion},
	"ranger":    {AnimalHandling, Athletics, Insight, Investigation},
	"rogue":     {Acrobatics, Athletics, Deception, Insight, Insight, Religion},
	"sorcerer":  {Arcana, Deception, Insight, Intimidation, Persuasion, Religion},
	"warlock":   {Arcana, Deception, Insight, Religion},
	"wizard":    {Arcana, History, Insight, Religion},
}

// ------------------------
// Skill Abilities
// ------------------------
var SkillAbilities = map[string]string{
	Acrobatics:     "Dexterity",
	AnimalHandling: "Wisdom",
	Arcana:         "Intelligence",
	Athletics:      "Strength",
	Deception:      "Charisma",
	History:        "Intelligence",
	Insight:        "Wisdom",
	Intimidation:   "Charisma",
	Investigation:  "Intelligence",
	Medicine:       "Wisdom",
	"Nature":       "Intelligence",
	"Perception":   "Wisdom",
	"Performance":  "Charisma",
	Persuasion:     "Charisma",
	Religion:       "Intelligence",
	SleightOfHand:  "Dexterity",
	Stealth:        "Dexterity",
	Survival:       "Wisdom",
}

// ------------------------
// Spellcasting Classes
// ------------------------
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
// Full Caster Spell Slots
// ------------------------
var FullCasterSpellSlots = [][]int{
	{2},
	{3},
	{4},
	{4, 2},
	{4, 3},
	{4, 3, 2},
	{4, 3, 3},
	{4, 3, 3, 1},
	{4, 3, 3, 2},
	{4, 3, 3, 3},
	{4, 3, 3, 3, 1},
	{4, 3, 3, 3, 2},
	{4, 3, 3, 3, 2, 1},
	{4, 3, 3, 3, 2, 1},
	{4, 3, 3, 3, 2, 1, 1},
	{4, 3, 3, 3, 2, 1, 1},
	{4, 3, 3, 3, 2, 1, 1, 1},
	{4, 3, 3, 3, 2, 1, 1, 1},
	{4, 3, 3, 3, 2, 1, 1, 1, 1},
	{4, 3, 3, 3, 2, 1, 1, 1, 1},
}
