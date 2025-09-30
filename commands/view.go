package commands

import (
	"dnd-character-sheet/storage"
	"fmt"
	"math"
	"strings"
)

func ViewCharacter(name string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return err
	}

	for _, c := range characters {
		if c.Name == name {
			fmt.Printf("Name: %s\n", c.Name)
			fmt.Printf("Class: %s\n", strings.ToLower(c.Class))
			fmt.Printf("Race: %s\n", strings.ToLower(c.Race))
			fmt.Printf("Background: %s\n", strings.ToLower(c.Background))
			fmt.Printf("Level: %d\n", c.Level)

			fmt.Printf("Ability scores:\n")
			fmt.Printf("  STR: %d (%+d)\n", c.Abilities.Strength, abilityModifier(c.Abilities.Strength))
			fmt.Printf("  DEX: %d (%+d)\n", c.Abilities.Dexterity, abilityModifier(c.Abilities.Dexterity))
			fmt.Printf("  CON: %d (%+d)\n", c.Abilities.Constitution, abilityModifier(c.Abilities.Constitution))
			fmt.Printf("  INT: %d (%+d)\n", c.Abilities.Intelligence, abilityModifier(c.Abilities.Intelligence))
			fmt.Printf("  WIS: %d (%+d)\n", c.Abilities.Wisdom, abilityModifier(c.Abilities.Wisdom))
			fmt.Printf("  CHA: %d (%+d)\n", c.Abilities.Charisma, abilityModifier(c.Abilities.Charisma))

			fmt.Printf("Proficiency bonus: %+d\n", proficiencyBonus(c.Level))
			fmt.Printf("Skill proficiencies: %s\n", formatSkillProficiencies(c.SkillProficiencies))

			// 🔑 Equipment toevoegen
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

			return nil
		}
	}

	return fmt.Errorf("character not found: %s", name)
}

func abilityModifier(score int) int {
	return int(math.Floor(float64(score-10) / 2))
}

func proficiencyBonus(level int) int {
	return 2 + (level-1)/4
}

func formatSkillProficiencies(skills []string) string {
	for i := range skills {
		skills[i] = lowerCase(skills[i])
	}
	return join(skills, ", ")
}

func lowerCase(s string) string {
	return strings.ToLower(s)
}

func join(items []string, sep string) string {
	return strings.Join(items, sep)
}
