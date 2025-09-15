package storage

import (
	"dnd-character-sheet/models"
	"encoding/json"
	"errors"
	"os"
)

const charactersFilePath = "characters.json"

// SaveCharacter voegt een nieuw character toe of update een bestaand character
func SaveCharacter(character models.Character) error {
	allCharacters, _ := LoadCharacters()
	allCharacters[character.Name] = character

	data, err := json.MarshalIndent(allCharacters, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(charactersFilePath, data, 0644)
}

// LoadCharacters leest alle characters uit het JSON bestand
func LoadCharacters() (map[string]models.Character, error) {
	characters := make(map[string]models.Character)

	if _, err := os.Stat(charactersFilePath); errors.Is(err, os.ErrNotExist) {
		return characters, nil
	}

	fileData, err := os.ReadFile(charactersFilePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(fileData, &characters); err != nil {
		return nil, err
	}

	return characters, nil
}

// SaveAllCharacters overschrijft alle characters in het JSON bestand
func SaveAllCharacters(allCharacters map[string]models.Character) error {
	fileData, err := json.MarshalIndent(allCharacters, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(charactersFilePath, fileData, 0644)
}

// GetNextCharacterID geeft het volgende beschikbare ID terug
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

// DeleteCharacter verwijdert een character op basis van naam
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

// GetCharacterByName haalt een character op basis van naam
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
