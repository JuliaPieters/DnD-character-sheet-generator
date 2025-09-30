package commands

import (
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
	"errors"
	"fmt"
)

// CreateCharacter maakt een nieuw character aan en slaat deze op
func CreateCharacter(
	characterName string,
	playerName string,
	characterRace string,
	characterClass string,
	characterBackground string,
	characterLevel int,
	abilityScores []int,
	skillProficiencies []string,
) error {
	// Laad bestaande characters
	existingCharacters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("failed to load characters: %w", err)
	}

	if _, exists := existingCharacters[characterName]; exists {
		return errors.New("character met deze naam bestaat al")
	}

	// Als geen scores meegegeven, gebruik standaard array
	if len(abilityScores) != 6 {
		abilityScores = nil
	}

	// Skill proficiencies beperken en uniek maken
	if len(skillProficiencies) == 0 {
		availableSkills := models.GetAvailableSkills(characterClass)
		skillSet := map[string]bool{}
		uniqueSkills := []string{}
		for _, s := range availableSkills {
			if !skillSet[s] {
				skillSet[s] = true
				uniqueSkills = append(uniqueSkills, s)
			}
			if len(uniqueSkills) == 4 { // maximaal 4 skills
				break
			}
		}
		skillProficiencies = uniqueSkills
	}

	newCharacterID := len(existingCharacters) + 1

	// Maak nieuw character aan met de juiste abilityScores en skills
	newCharacter := models.NewCharacter(
		newCharacterID,
		characterName,
		characterRace,
		characterClass,
		characterBackground,
		characterLevel,
		abilityScores,
		skillProficiencies,
	)

	if err := storage.SaveCharacter(*newCharacter); err != nil {
		return fmt.Errorf("failed to save character: %w", err)
	}
	return nil
}
