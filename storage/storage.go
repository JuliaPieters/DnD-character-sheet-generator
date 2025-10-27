package storage

import (
	"dnd-character-sheet/domain"
	"encoding/json"
	"errors"
	"os"
)

var CharactersFilePath = "characters.json"

func SaveCharacter(character *domain.Character) error {
	allCharacters, _ := LoadCharacters()
	allCharacters[character.Name] = *character
	return SaveAllCharacters(allCharacters)
}

func LoadCharacters() (map[string]domain.Character, error) {
	characters := make(map[string]domain.Character)

	if _, err := os.Stat(CharactersFilePath); errors.Is(err, os.ErrNotExist) {
		return characters, nil
	}

	fileData, err := os.ReadFile(CharactersFilePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(fileData, &characters); err != nil {
		return nil, err
	}

	return characters, nil
}

func SaveAllCharacters(allCharacters map[string]domain.Character) error {
	data, err := json.MarshalIndent(allCharacters, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(CharactersFilePath, data, 0644)
}

func GetNextCharacterID() (int, error) {
	allCharacters, err := LoadCharacters()
	if err != nil {
		return 0, err
	}

	highestID := 0
	for _, c := range allCharacters {
		if c.ID > highestID {
			highestID = c.ID
		}
	}

	return highestID + 1, nil
}

func DeleteCharacter(characterName string) error {
	allCharacters, err := LoadCharacters()
	if err != nil {
		return err
	}

	if _, exists := allCharacters[characterName]; !exists {
		return errors.New("character not found")
	}

	delete(allCharacters, characterName)
	return SaveAllCharacters(allCharacters)
}

func GetCharacterByName(characterName string) (*domain.Character, error) {
	allCharacters, err := LoadCharacters()
	if err != nil {
		return nil, err
	}

	c, exists := allCharacters[characterName]
	if !exists {
		return nil, errors.New("character not found")
	}

	return &c, nil
}
