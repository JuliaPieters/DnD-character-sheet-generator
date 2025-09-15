package commands

import (
	"dnd-character-sheet/storage"
	"fmt"
)

// ViewCharacter toont alle details van een character op basis van naam
func ViewCharacter(characterName string) error {
	allCharacters, err := storage.LoadCharacters()
	if err != nil {
		return err
	}

	for _, character := range allCharacters {
		if character.Name == characterName {
			fmt.Printf("Name: %s\n", character.Name)
			fmt.Printf("Player Name: %s\n", character.PlayerName)
			fmt.Printf("Level: %d | Race: %s | Class: %s\n", character.Level, character.Race, character.Class)
			fmt.Printf("Background: %s | Alignment: %s\n", character.Background, character.Alignment)
			fmt.Printf("Abilities:\n")
			fmt.Printf("  Strength: %d | Dexterity: %d | Constitution: %d\n",
				character.Abilities.Strength, character.Abilities.Dexterity, character.Abilities.Constitution)
			fmt.Printf("  Intelligence: %d | Wisdom: %d | Charisma: %d\n",
				character.Abilities.Intelligence, character.Abilities.Wisdom, character.Abilities.Charisma)
			fmt.Printf("Skill Modifiers: %v\n", character.Skills)
			fmt.Printf("Armor Class: %d | Initiative: %d | Passive Perception: %d\n",
				character.ArmorClass, character.Initiative, character.PassivePerception)
			fmt.Printf("Hit Points: %d/%d | Speed: %d\n", character.CurrentHitPoints, character.MaxHitPoints, character.Speed)
			return nil
		}
	}

	return fmt.Errorf("character not found: %s", characterName)
}
