package commands

import (
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
	"fmt"
)

var SpellcastingClasses = map[string]bool{
	"bard":     true,
	"cleric":   true,
	"druid":    true,
	"paladin":  true,
	"ranger":   true,
	"sorcerer": true,
	"warlock":  true,
	"wizard":   true,
}

var PreparedCasters = map[string]bool{
	"cleric":  true,
	"druid":   true,
	"paladin": true,
	"wizard":  true,
}

var StartingSpells = map[string][]string{
	"wizard":   {"burning hands", "disguise self"},
	"cleric":   {"guidance", "sacred flame"},
	"druid":    {"shillelagh", "thorn whip"},
	"bard":     {"vicious mockery", "dancing lights"},
	"sorcerer": {"fire bolt", "light"},
	"warlock":  {"eldritch blast", "mage hand"},
	"paladin":  {"divine sense", "lay on hands"},
	"ranger":   {"hunter's mark", "cure wounds"},
}

func GiveStartingSpells(character *models.Character) error {
	spells, ok := StartingSpells[character.Class]
	if !ok {
		return nil
	}

	for _, name := range spells {
		level := 0 // cantrips
		spell := models.Spell{
			Name:     name,
			Level:    level,
			Prepared: false,
		}
		character.Spells = append(character.Spells, spell)
	}

	// Gebruik de Character-logica voor spell slots
	character.SetupSpellcasting()

	if err := storage.SaveCharacter(*character); err != nil {
		return fmt.Errorf("could not save character: %w", err)
	}

	return nil
}

func LearnSpell(characterName, spellName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return err
	}
	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character \"%s\" not found", characterName)
	}

	if !SpellcastingClasses[character.Class] {
		return fmt.Errorf("this class can't cast spells")
	}

	if character.CanPrepareSpells {
		return fmt.Errorf("this class prepares spells and can't learn them")
	}

	for _, s := range character.Spells {
		if s.Name == spellName {
			return fmt.Errorf("character '%s' already knows spell '%s'", characterName, spellName)
		}
	}

	character.Spells = append(character.Spells, models.Spell{
		Name:     spellName,
		Level:    1,
		Prepared: false,
	})

	if err := storage.SaveCharacter(character); err != nil {
		return err
	}

	fmt.Printf("Learned spell %s\n", spellName)
	return nil
}

func PrepareSpell(characterName, spellName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return err
	}
	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character \"%s\" not found", characterName)
	}

	if !SpellcastingClasses[character.Class] {
		return fmt.Errorf("this class can't cast spells")
	}

	if !character.CanPrepareSpells {
		return fmt.Errorf("this class learns spells and can't prepare them")
	}

	for i, spell := range character.Spells {
		if spell.Name == spellName {
			character.Spells[i].Prepared = true
			if err := storage.SaveCharacter(character); err != nil {
				return err
			}
			fmt.Printf("Prepared spell %s\n", spellName)
			return nil
		}
	}

	return fmt.Errorf("spell '%s' not known by character '%s'", spellName, characterName)
}

func RemoveSpell(characterName, spellName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return err
	}
	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character \"%s\" not found", characterName)
	}

	newSpells := []models.Spell{}
	removed := false
	for _, s := range character.Spells {
		if s.Name != spellName {
			newSpells = append(newSpells, s)
		} else {
			removed = true
		}
	}
	character.Spells = newSpells

	if err := storage.SaveCharacter(character); err != nil {
		return err
	}

	if removed {
		fmt.Printf("Removed spell %s\n", spellName)
	} else {
		fmt.Printf("Spell '%s' not found for character '%s'\n", spellName, characterName)
	}

	return nil
}
