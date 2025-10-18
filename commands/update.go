package commands

import (
	"dnd-character-sheet/storage"
	"fmt"
)

func UpdateCharacterLevel(characterName string, newLevel int) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("cannot load characters: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' does not exist", characterName)
	}

	character.UpdateLevel(newLevel)

	if err := storage.SaveCharacter(character); err != nil {
		return fmt.Errorf("cannot save character: %w", err)
	}

	return nil
}
