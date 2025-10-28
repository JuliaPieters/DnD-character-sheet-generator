package application

import (
	"strings"

	"dnd-character-sheet/domain"
)

type CharacterService struct{}

// ------------------------
// Parameter struct voor NewCharacter
// ------------------------
type NewCharacterParams struct {
	ID            int
	Name          string
	Race          string
	Class         string
	Background    string
	Level         int
	AbilityScores []int
	SkillChoices  []string
	SpellService  *SpellService
}

// ------------------------
// Character creation
// ------------------------
func (s *CharacterService) NewCharacter(params NewCharacterParams) *domain.Character {
	raceKey := strings.ToLower(params.Race)
	classKey := strings.ToLower(params.Class)

	mod := domain.RaceModifiers[raceKey]

	abilities := s.prepareAbilities(params.AbilityScores, mod)

	char := &domain.Character{
		ID:                 params.ID,
		Name:               params.Name,
		Race:               raceKey,
		Class:              classKey,
		Level:              params.Level,
		Background:         params.Background,
		ProficiencyBonus:   s.CalculateProfBonus(params.Level),
		Abilities:          abilities,
		SkillProficiencies: params.SkillChoices,
		Skills:             make(map[string]int),
		ArmorClass:         10,
		Speed:              30,
		MaxHitPoints:       10,
		CurrentHitPoints:   10,
		Equipment:          domain.Equipment{},
	}

	s.UpdateModifiers(char)
	s.CalculateAllSkills(char)
	s.CalculateCombatStats(char)

	if params.SpellService != nil {
		params.SpellService.SetupSpellcasting(char)
	}

	return char
}

func (s *CharacterService) prepareAbilities(scores []int, mod map[string]int) domain.AbilityScores {
	if len(scores) == 6 {
		return domain.AbilityScores{
			Strength:     scores[0] + mod["Strength"],
			Dexterity:    scores[1] + mod["Dexterity"],
			Constitution: scores[2] + mod["Constitution"],
			Intelligence: scores[3] + mod["Intelligence"],
			Wisdom:       scores[4] + mod["Wisdom"],
			Charisma:     scores[5] + mod["Charisma"],
		}
	}
	return s.AssignAbilities(mod)
}

// ------------------------
// Leveling
// ------------------------
func (s *CharacterService) LevelUp(c *domain.Character, newLevel int, spellService *SpellService) {
	c.Level = newLevel
	c.ProficiencyBonus = s.CalculateProfBonus(newLevel)
	s.UpdateModifiers(c)
	s.CalculateAllSkills(c)
	s.CalculateCombatStats(c)
	if spellService != nil {
		spellService.SetupSpellcasting(c)
	}
}

// ------------------------
// Abilities
// ------------------------
func (s *CharacterService) AssignAbilities(mod map[string]int) domain.AbilityScores {
	abilities := domain.AbilityScores{}
	order := []string{"Strength", "Dexterity", "Constitution", "Intelligence", "Wisdom", "Charisma"}
	for i, name := range order {
		score := domain.StandardArray[i] + mod[name]
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

func (s *CharacterService) CalculateProfBonus(level int) int {
	return 2 + (level-1)/4
}

func (s *CharacterService) UpdateModifiers(c *domain.Character) {
	c.StrengthMod = c.Abilities.Modifier("Strength")
	c.DexterityMod = c.Abilities.Modifier("Dexterity")
	c.ConstitutionMod = c.Abilities.Modifier("Constitution")
	c.IntelligenceMod = c.Abilities.Modifier("Intelligence")
	c.WisdomMod = c.Abilities.Modifier("Wisdom")
	c.CharismaMod = c.Abilities.Modifier("Charisma")
}

// ------------------------
// Skills
// ------------------------
func (s *CharacterService) CalculateAllSkills(c *domain.Character) {
	c.Skills = make(map[string]int)
	for skill, ability := range domain.SkillAbilities {
		mod := c.Abilities.Modifier(ability)
		if s.contains(c.SkillProficiencies, skill) {
			mod += c.ProficiencyBonus
		}
		c.Skills[skill] = mod
	}
}

func (s *CharacterService) GetAvailableSkills(className string) []string {
	classKey := strings.ToLower(className)
	if skills, ok := domain.ClassSkills[classKey]; ok {
		return skills
	}
	return []string{}
}

// ------------------------
// Utility
// ------------------------
func (s *CharacterService) contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
