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
	abilityScores []int,        // ints in plaats van []string
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

	newCharacterID := len(existingCharacters) + 1

	// Standaard volgorde voor ability names
	abilityOrder := []string{"Strength", "Dexterity", "Constitution", "Intelligence", "Wisdom", "Charisma"}

	// Maak nieuw character aan
	newCharacter := models.NewCharacter(
		newCharacterID,
		characterName,
		characterRace,
		characterClass,
		characterBackground,
		characterLevel,
		abilityOrder,       // names blijven nodig voor NewCharacter
		skillProficiencies,
	)

	// Overschrijf ability scores met de ints die doorgegeven zijn
	if len(abilityScores) != 6 {
		return errors.New("abilityScores moet exact 6 waarden bevatten")
	}
	newCharacter.Abilities.Strength = abilityScores[0]
	newCharacter.Abilities.Dexterity = abilityScores[1]
	newCharacter.Abilities.Constitution = abilityScores[2]
	newCharacter.Abilities.Intelligence = abilityScores[3]
	newCharacter.Abilities.Wisdom = abilityScores[4]
	newCharacter.Abilities.Charisma = abilityScores[5]

	// Herbereken skills en stats
	newCharacter.CalculateAllSkills()
	newCharacter.CalculateCombatStats()
	newCharacter.SetupSpellcasting()

	// Opslaan
	if err := storage.SaveCharacter(*newCharacter); err != nil {
		return fmt.Errorf("failed to save character: %w", err)
	}
	return nil
}
