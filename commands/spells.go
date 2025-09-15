package commands

import (
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
	"fmt"
)

// AddSpell adds a spell to a caster character
func AddSpell(characterName string, newSpell models.Spell) error {
	characters, loadErr := storage.LoadCharacters()
	if loadErr != nil {
		return fmt.Errorf("could not load characters: %w", loadErr)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	if _, isCaster := models.SpellcastingClasses[character.Class]; !isCaster {
		return fmt.Errorf("character '%s' is not a spellcaster", characterName)
	}

	character.Spells = append(character.Spells, newSpell)

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return fmt.Errorf("could not save character: %w", saveErr)
	}

	fmt.Printf("✅ Spell '%s' added to character '%s'\n", newSpell.Name, characterName)
	return nil
}

// RemoveSpell removes a spell from a caster character
func RemoveSpell(characterName string, spellName string) error {
	characters, loadErr := storage.LoadCharacters()
	if loadErr != nil {
		return fmt.Errorf("could not load characters: %w", loadErr)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}

	if _, isCaster := models.SpellcastingClasses[character.Class]; !isCaster {
		return fmt.Errorf("character '%s' is not a spellcaster", characterName)
	}

	newSpellList := []models.Spell{}
	for _, existingSpell := range character.Spells {
		if existingSpell.Name != spellName {
			newSpellList = append(newSpellList, existingSpell)
		}
	}
	character.Spells = newSpellList

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return fmt.Errorf("could not save character: %w", saveErr)
	}

	fmt.Printf("✅ Spell '%s' removed from character '%s'\n", spellName, characterName)
	return nil
}
