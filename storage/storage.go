package storage

import (
    "encoding/json"
    "errors"
    "os"
    "dnd-character-sheet/models"
)

const filePath = "characters.json"

// SaveCharacter adds or updates a character
func SaveCharacter(c models.Character) error {
    chars, _ := LoadCharacters()
    chars[c.Name] = c
    data, err := json.MarshalIndent(chars, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(filePath, data, 0644)
}

// LoadCharacters loads all characters from file
func LoadCharacters() (map[string]models.Character, error) {
    chars := make(map[string]models.Character)

    if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
        return chars, nil
    }

    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    if err := json.Unmarshal(data, &chars); err != nil {
        return nil, err
    }

    return chars, nil
}

// SaveAllCharacters overwrites all characters in the JSON file
func SaveAllCharacters(chars map[string]models.Character) error {
	data, err := json.MarshalIndent(chars, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
