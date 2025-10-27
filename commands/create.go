package commands

import (
	"dnd-character-sheet/application"
	"dnd-character-sheet/storage"
	"errors"
	"fmt"
)

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

	existingCharacters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("failed to load characters: %w", err)
	}

	if _, exists := existingCharacters[characterName]; exists {
		return errors.New("character met deze naam bestaat al")
	}

	if len(abilityScores) != 6 {
		abilityScores = nil
	}

	characterService := &application.CharacterService{}
	spellService := &application.SpellService{}

	if len(skillProficiencies) == 0 {
		availableSkills := characterService.GetAvailableSkills(characterClass)
		uniqueSkills := []string{}
		skillSet := map[string]bool{}
		for _, s := range availableSkills {
			if !skillSet[s] {
				skillSet[s] = true
				uniqueSkills = append(uniqueSkills, s)
			}
			if len(uniqueSkills) == 4 {
				break
			}
		}
		skillProficiencies = uniqueSkills
	}

	newCharacter := characterService.NewCharacter(
		len(existingCharacters)+1,
		characterName,
		characterRace,
		characterClass,
		characterBackground,
		characterLevel,
		abilityScores,
		skillProficiencies,
		spellService, 
	)

	if err := GiveStartingSpells(newCharacter); err != nil {
		return fmt.Errorf("failed to give starting spells: %w", err)
	}

	if err := storage.SaveCharacter(newCharacter); err != nil {
		return fmt.Errorf("failed to save character: %w", err)
	}

	return nil
}
