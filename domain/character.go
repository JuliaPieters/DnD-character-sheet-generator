package domain

type Character struct {
	ID                 int
	Name               string
	PlayerName         string
	Race               string
	Class              string
	Level              int
	Background         string
	Alignment          string
	ProficiencyBonus   int
	Abilities          AbilityScores
	SkillProficiencies []string
	Skills             map[string]int

	StrengthMod     int
	DexterityMod    int
	ConstitutionMod int
	IntelligenceMod int
	WisdomMod       int
	CharismaMod     int

	Equipment Equipment

	Spells     []Spell
	SpellSlots map[int]int

	ArmorClass        int
	Initiative        int
	PassivePerception int

	SpellcastingAbility string
	SpellSaveDC         int
	SpellAttackBonus    int

	CanPrepareSpells bool

	ExperiencePoints   int
	Speed              int
	MaxHitPoints       int
	CurrentHitPoints   int
	TemporaryHitPoints int
	HitDiceTotal       string
	HitDiceRemaining   string
	DeathSaveSuccesses int
	DeathSaveFailures  int

	CopperPieces   int
	SilverPieces   int
	ElectrumPieces int
	GoldPieces     int
	PlatinumPieces int
	EquipmentText  string
	Personality    string
	Ideals         string
	Bonds          string
	Flaws          string
	Features       string
}
