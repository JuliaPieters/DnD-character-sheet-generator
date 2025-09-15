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
			fmt.Printf("Skill proficiencies: %s\n", formatSkillProficiencies(c.SkillProficiencies))

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

// formatSkillProficiencies maakt een comma-separated lijst van skill proficiencies
func formatSkillProficiencies(skills []string) string {
	for i := range skills {
		skills[i] = lowerCase(skills[i])
	}
	return join(skills, ", ")
}

// lowerCase zet een string naar lowercase
func lowerCase(s string) string {
	return strings.ToLower(s)
}

// join voegt een slice van strings samen met separator
func join(items []string, sep string) string {
	return strings.Join(items, sep)
}
