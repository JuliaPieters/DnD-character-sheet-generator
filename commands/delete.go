package commands

import (
	"dnd-character-sheet/storage"
	"fmt"
)

func DeleteCharacter(characterName string) error {
	allCharacters, err := storage.LoadCharacters()
	if err != nil {
		return err
	}

	for key, character := range allCharacters {
		if character.Name == characterName {
			delete(allCharacters, key)
			if err := storage.SaveAllCharacters(allCharacters); err != nil {
				return fmt.Errorf("failed to delete character: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("character not found: %s", characterName)
}
