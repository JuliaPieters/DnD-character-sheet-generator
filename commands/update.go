package commands

import (
	"dnd-character-sheet/storage"
	"fmt"
)

func UpdateCharacterLevel(characterName string, nieuwLevel int) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("kan characters niet laden: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' bestaat niet", characterName)
	}

	character.Level = nieuwLevel

	if err := storage.SaveCharacter(character); err != nil {
		return fmt.Errorf("kan character niet opslaan: %w", err)
	}
	return nil
}
