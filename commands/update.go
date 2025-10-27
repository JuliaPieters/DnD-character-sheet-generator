package commands

import (
	"dnd-character-sheet/application"
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

	characterPtr := &character

	charService := application.CharacterService{}
	charService.LevelUp(characterPtr, newLevel, nil)

	characters[characterName] = *characterPtr

	if err := storage.SaveAllCharacters(characters); err != nil {
		return fmt.Errorf("cannot save character: %w", err)
	}

	return nil
}
