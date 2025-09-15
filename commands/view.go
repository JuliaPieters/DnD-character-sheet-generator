package commands

import (
	"dnd-character-sheet/storage"
	"fmt"
	"strings"
)

// ViewCharacter toont een character volledig volgens testoutput
func ViewCharacter(name string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return err
	}

	for _, c := range characters {
		if c.Name == name {
			fmt.Printf("Name: %s\n", c.Name)
			fmt.Printf("Class: %s\n", c.Class)
			fmt.Printf("Race: %s\n", c.Race)
			fmt.Printf("Background: %s\n", c.Background)
			fmt.Printf("Level: %d\n", c.Level)

			fmt.Printf("Ability scores:\n")
			fmt.Printf("  STR: %d (%+d)\n", c.Abilities.Strength, abilityModifier(c.Abilities.Strength))
			fmt.Printf("  DEX: %d (%+d)\n", c.Abilities.Dexterity, abilityModifier(c.Abilities.Dexterity))
			fmt.Printf("  CON: %d (%+d)\n", c.Abilities.Constitution, abilityModifier(c.Abilities.Constitution))
			fmt.Printf("  INT: %d (%+d)\n", c.Abilities.Intelligence, abilityModifier(c.Abilities.Intelligence))
			fmt.Printf("  WIS: %d (%+d)\n", c.Abilities.Wisdom, abilityModifier(c.Abilities.Wisdom))
			fmt.Printf("  CHA: %d (%+d)\n", c.Abilities.Charisma, abilityModifier(c.Abilities.Charisma))

			fmt.Printf("Proficiency bonus: %+d\n", proficiencyBonus(c.Level))
			fmt.Printf("Skill proficiencies: %s\n", formatSkills(c.Skills))

			return nil
		}
	}

	return fmt.Errorf("character not found: %s", name)
}

// abilityModifier berekent de modifier van een ability score
func abilityModifier(score int) int {
	return (score - 10) / 2
}

// proficiencyBonus berekent de D&D 5e proficiency bonus op basis van level
func proficiencyBonus(level int) int {
	return 2 + (level-1)/4
}

// formatSkills zet de skills map om naar een comma-separated lijst van skill names
func formatSkills(skills map[string]int) string {
	names := []string{}
	for skill := range skills {
		names = append(names, skill)
	}
	return joinLower(names)
}

// joinLower join met comma en lowercase
func joinLower(items []string) string {
	for i := range items {
		items[i] = strings.ToLower(items[i])
	}
	return strings.Join(items, ", ")
}
