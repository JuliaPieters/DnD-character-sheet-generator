package application

import "dnd-character-sheet/domain"

type SpellService struct{}

func (s *SpellService) SetupSpellcasting(c *domain.Character) {
	if !s.CanCastSpells(c.Class) {
		c.SpellcastingAbility = ""
		c.SpellSaveDC = 0
		c.SpellAttackBonus = 0
		c.SpellSlots = nil
		c.CanPrepareSpells = false
		return
	}

	if ability, ok := domain.SpellcastingClasses[c.Class]; ok {
		mod := c.Abilities.Modifier(ability)
		c.SpellcastingAbility = ability
		c.SpellSaveDC = 8 + c.ProficiencyBonus + mod
		c.SpellAttackBonus = c.ProficiencyBonus + mod
	}

	c.CanPrepareSpells = s.IsPreparedCaster(c.Class)
	s.UpdateSpellSlots(c)
}

func (s *SpellService) CanCastSpells(className string) bool {
	_, ok := domain.SpellcastingClasses[className]
	return ok
}

func (s *SpellService) IsPreparedCaster(className string) bool {
	preparedClasses := map[string]bool{
		"cleric":  true,
		"druid":   true,
		"paladin": true,
		"wizard":  true,
	}
	return preparedClasses[className]
}

func (s *SpellService) UpdateSpellSlots(c *domain.Character) {
	if !s.CanCastSpells(c.Class) {
		c.SpellSlots = nil
		return
	}

	c.SpellSlots = make(map[int]int)

	switch c.Class {
	case "wizard", "cleric", "druid", "bard", "sorcerer":
		if c.Level <= len(domain.FullCasterSpellSlots) {
			for lvl, slots := range domain.FullCasterSpellSlots[c.Level-1] {
				c.SpellSlots[lvl+1] = slots
			}
		}
	case "paladin":
		c.SpellSlots = calculatePaladinSlots(c.Level)
	case "ranger":
		c.SpellSlots = calculateRangerSlots(c.Level)
	case "warlock":
		c.SpellSlots = calculateWarlockSlots(c.Level)
	}
}

func calculatePaladinSlots(level int) map[int]int {
	slots := map[int]int{}
	if level >= 1 {
		slots[1] = 4
	}
	if level >= 2 {
		slots[2] = 3
	}
	if level >= 3 {
		slots[3] = 3
	}
	if level >= 4 {
		slots[4] = 3
	}
	if level >= 5 {
		slots[5] = 2
	}
	return slots
}

func calculateRangerSlots(level int) map[int]int {
	slots := map[int]int{}
	slots[1] = (level + 1) / 2
	if level >= 4 {
		slots[2] = level / 2
	}
	return slots
}

func calculateWarlockSlots(level int) map[int]int {
	slots := map[int]int{}
	slots[0] = 4
	var pactLevel, numSlots int

	switch {
	case level >= 1 && level <= 1:
		numSlots = 1
		pactLevel = 1
	case level >= 2 && level <= 8:
		numSlots = 2
		pactLevel = 1
	case level >= 9 && level <= 11:
		numSlots = 3
		pactLevel = 2
	case level >= 12 && level <= 16:
		numSlots = 3
		pactLevel = 3
	case level >= 17 && level <= 20:
		numSlots = 4
		pactLevel = 5
	}

	slots[pactLevel] = numSlots
	return slots
}
