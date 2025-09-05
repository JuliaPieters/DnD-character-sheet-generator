package commands

import (
	"dnd-character-sheet/storage"
	"fmt"
)

// ViewCharacter toont alle details van een character in een nette layout
func ViewCharacter(name string) error {
	chars, err := storage.LoadCharacters()
	if err != nil {
		return err
	}

	c, ok := chars[name]
	if !ok {
		return fmt.Errorf("character %s not found", name)
	}

	fmt.Printf("📜 Character Sheet: %s\n", c.Name)
	fmt.Printf("Race: %s | Class: %s | Level: %d | Background: %s\n", c.Race, c.Class, c.Level, c.Background)
	fmt.Printf("Proficiency Bonus: +%d\n", c.ProfBonus)
	fmt.Println("Abilities:")
	fmt.Printf("  STR: %d (%+d)\n", c.Abilities.Strength, c.Abilities.Modifier("STR"))
	fmt.Printf("  DEX: %d (%+d)\n", c.Abilities.Dexterity, c.Abilities.Modifier("DEX"))
	fmt.Printf("  CON: %d (%+d)\n", c.Abilities.Constitution, c.Abilities.Modifier("CON"))
	fmt.Printf("  INT: %d (%+d)\n", c.Abilities.Intelligence, c.Abilities.Modifier("INT"))
	fmt.Printf("  WIS: %d (%+d)\n", c.Abilities.Wisdom, c.Abilities.Modifier("WIS"))
	fmt.Printf("  CHA: %d (%+d)\n", c.Abilities.Charisma, c.Abilities.Modifier("CHA"))

	if len(c.Skills) > 0 {
		fmt.Println("Skills:")
		for skill, mod := range c.Skills {
			fmt.Printf("  %s: %+d\n", skill, mod)
		}
	}

	return nil
}
