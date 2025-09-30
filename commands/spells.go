package commands

import (
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
	"fmt"
)

// LearnSpell adds a spell to a caster character's known spells
func LearnSpell(characterName string, newSpell models.Spell) error {
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

	fmt.Printf("Learned spell %s\n", newSpell.Name)
	return nil
}

// PrepareSpell marks a known spell as prepared
func PrepareSpell(characterName string, spellName string) error {
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

	for i, spell := range character.Spells {
		if spell.Name == spellName {
			character.Spells[i].Prepared = true
			if saveErr := storage.SaveCharacter(character); saveErr != nil {
				return fmt.Errorf("could not save character: %w", saveErr)
			}
			fmt.Printf("Prepared spell %s\n", spellName)
			return nil
		}
	}

	return fmt.Errorf("spell '%s' not known by character '%s'", spellName, characterName)
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
	removed := false
	for _, existingSpell := range character.Spells {
		if existingSpell.Name != spellName {
			newSpellList = append(newSpellList, existingSpell)
		} else {
			removed = true
		}
	}
	character.Spells = newSpellList

	if saveErr := storage.SaveCharacter(character); saveErr != nil {
		return fmt.Errorf("could not save character: %w", saveErr)
	}

	if removed {
		fmt.Printf("Removed spell %s\n", spellName)
	} else {
		fmt.Printf("Spell '%s' not found for character '%s'\n", spellName, characterName)
	}

	return nil
}
