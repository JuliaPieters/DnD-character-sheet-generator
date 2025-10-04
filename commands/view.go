package commands

import (
	"dnd-character-sheet/storage"
	"fmt"
	"sort"
	"strings"
)

var fullCasters = map[string]bool{
	"wizard": true, "cleric": true, "druid": true, "bard": true, "sorcerer": true,
}

var pactCasters = map[string]bool{
	"warlock": true,
}

func ViewCharacter(name string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return err
	}

	for _, c := range characters {
		if c.Name != name {
			continue
		}

		c.CalculateCombatStats()

		fmt.Printf("Name: %s\n", c.Name)
		fmt.Printf("Class: %s\n", strings.ToLower(c.Class))
		fmt.Printf("Race: %s\n", strings.ToLower(c.Race))
		fmt.Printf("Background: %s\n", strings.ToLower(c.Background))
		fmt.Printf("Level: %d\n", c.Level)

		fmt.Println("Ability scores:")
		fmt.Printf("  STR: %d (%+d)\n", c.Abilities.Strength, c.Abilities.Modifier("Strength"))
		fmt.Printf("  DEX: %d (%+d)\n", c.Abilities.Dexterity, c.Abilities.Modifier("Dexterity"))
		fmt.Printf("  CON: %d (%+d)\n", c.Abilities.Constitution, c.Abilities.Modifier("Constitution"))
		fmt.Printf("  INT: %d (%+d)\n", c.Abilities.Intelligence, c.Abilities.Modifier("Intelligence"))
		fmt.Printf("  WIS: %d (%+d)\n", c.Abilities.Wisdom, c.Abilities.Modifier("Wisdom"))
		fmt.Printf("  CHA: %d (%+d)\n", c.Abilities.Charisma, c.Abilities.Modifier("Charisma"))

		fmt.Printf("Proficiency bonus: %+d\n", c.ProficiencyBonus)
		fmt.Printf("Skill proficiencies: %s\n", formatSkillProficiencies(c.SkillProficiencies))

		if c.Equipment.MainHand != nil {
			fmt.Printf("Main hand: %s\n", c.Equipment.MainHand.Name)
		}
		if c.Equipment.OffHand != nil {
			fmt.Printf("Off hand: %s\n", c.Equipment.OffHand.Name)
		}
		if c.Equipment.Armor != nil {
			fmt.Printf("Armor: %s\n", c.Equipment.Armor.Name)
		}
		if c.Equipment.Shield != nil {
			fmt.Printf("Shield: %s\n", c.Equipment.Shield.Name)
		}

		if len(c.SpellSlots) > 0 {
			fmt.Println("Spell slots:")
			levels := make([]int, 0, len(c.SpellSlots))
			for lvl := range c.SpellSlots {
				levels = append(levels, lvl)
			}
			sort.Ints(levels)
			for _, lvl := range levels {
				fmt.Printf("  Level %d: %d\n", lvl, c.SpellSlots[lvl])
			}
		}

		if fullCasters[strings.ToLower(c.Class)] || pactCasters[strings.ToLower(c.Class)] {
			if c.SpellcastingAbility != "" {
				fmt.Printf("Spellcasting ability: %s\n", strings.ToLower(c.SpellcastingAbility))
				fmt.Printf("Spell save DC: %d\n", c.SpellSaveDC)
				fmt.Printf("Spell attack bonus: %+d\n", c.SpellAttackBonus)
			}
		}

		fmt.Printf("Armor class: %d\n", c.ArmorClass)
		fmt.Printf("Initiative bonus: %d\n", c.Initiative)
		fmt.Printf("Passive perception: %d\n", c.PassivePerception)

		return nil
	}

	return fmt.Errorf(`character "%s" not found`, name)
}

func formatSkillProficiencies(skills []string) string {
	for i := range skills {
		skills[i] = strings.ToLower(skills[i])
	}
	return strings.Join(skills, ", ")
}
