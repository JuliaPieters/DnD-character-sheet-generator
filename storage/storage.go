package storage

import (
	"dnd-character-sheet/models"
	"encoding/json"
	"errors"
	"os"
)

var CharactersFilePath = "characters.json"

func SaveCharacter(character models.Character) error {
	allCharacters, _ := LoadCharacters()
	allCharacters[character.Name] = character

	data, err := json.MarshalIndent(allCharacters, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(CharactersFilePath, data, 0644)
}

func LoadCharacters() (map[string]models.Character, error) {
	characters := make(map[string]models.Character)

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

func SaveAllCharacters(allCharacters map[string]models.Character) error {
	fileData, err := json.MarshalIndent(allCharacters, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(CharactersFilePath, fileData, 0644)
}

func GetNextCharacterID() (int, error) {
	allCharacters, err := LoadCharacters()
	if err != nil {
		return 0, err
	}

	highestID := 0
	for _, character := range allCharacters {
		if character.ID > highestID {
			highestID = character.ID
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

func GetCharacterByName(characterName string) (models.Character, error) {
	allCharacters, err := LoadCharacters()
	if err != nil {
		return models.Character{}, err
	}

	character, exists := allCharacters[characterName]
	if !exists {
		return models.Character{}, errors.New("character not found")
	}

	return character, nil
}
