package application

import (
	"strings"

	"dnd-character-sheet/domain"
)

type EquipmentService struct{}

func (s *CharacterService) CalculateCombatStats(c *domain.Character) {
	c.Initiative = c.Abilities.Modifier("Dexterity")

	c.PassivePerception = 10 + c.Abilities.Modifier("Wisdom")
	if p, ok := c.Skills["Perception"]; ok {
		c.PassivePerception = 10 + p
	}

	ac := 10

	if c.Equipment.Armor != nil {
		ac = c.Equipment.Armor.ArmorClass
		if c.Equipment.Armor.DexBonus {
			dexMod := c.Abilities.Modifier("Dexterity")
			if c.Equipment.Armor.MaxDexBonus > 0 && dexMod > c.Equipment.Armor.MaxDexBonus {
				dexMod = c.Equipment.Armor.MaxDexBonus
			}
			ac += dexMod
		}
	} else {
		switch strings.ToLower(c.Class) {
		case "barbarian":
			ac = 10 + c.Abilities.Modifier("Dexterity") + c.Abilities.Modifier("Constitution")
		case "monk":
			ac = 10 + c.Abilities.Modifier("Dexterity") + c.Abilities.Modifier("Wisdom")
		default:
			ac = 10 + c.Abilities.Modifier("Dexterity")
		}
	}

	if c.Equipment.Shield != nil {
		ac += c.Equipment.Shield.ArmorClass
	}

	c.ArmorClass = ac
}
