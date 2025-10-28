package commands

import (
	"dnd-character-sheet/application"
	"dnd-character-sheet/storage"
	"errors"
	"fmt"
)

type CreateCharacterParams struct {
	CharacterName      string
	PlayerName         string
	Race               string
	Class              string
	Background         string
	Level              int
	AbilityScores      []int
	SkillProficiencies []string
}

func CreateCharacter(params CreateCharacterParams) error {
	existingCharacters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("failed to load characters: %w", err)
	}

	if _, exists := existingCharacters[params.CharacterName]; exists {
		return errors.New("character met deze naam bestaat al")
	}

	if len(params.AbilityScores) != 6 {
		params.AbilityScores = nil
	}

	characterService := &application.CharacterService{}
	spellService := &application.SpellService{}

	if len(params.SkillProficiencies) == 0 {
		availableSkills := characterService.GetAvailableSkills(params.Class)
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
		params.SkillProficiencies = uniqueSkills
	}

	newCharacter := characterService.NewCharacter(application.NewCharacterParams{
		ID:            len(existingCharacters) + 1,
		Name:          params.CharacterName,
		Race:          params.Race,
		Class:         params.Class,
		Background:    params.Background,
		Level:         params.Level,
		AbilityScores: params.AbilityScores,
		SkillChoices:  params.SkillProficiencies,
		SpellService:  spellService,
	})

	newCharacter.PlayerName = params.PlayerName

	spellService.SetupSpellcasting(newCharacter)

	if err := GiveStartingSpells(newCharacter); err != nil {
		return fmt.Errorf("failed to give starting spells: %w", err)
	}

	if err := storage.SaveCharacter(newCharacter); err != nil {
		return fmt.Errorf("failed to save character: %w", err)
	}

	return nil
}
