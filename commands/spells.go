package commands

import (
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
	"fmt"
)

// Welke classes spellcasting gebruiken
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

// Welke classes spells kunnen preparen
var PreparedCasters = map[string]bool{
	"cleric":  true,
	"druid":   true,
	"paladin": true,
	"wizard":  true,
}

// Startspells per class
var StartingSpells = map[string][]string{
	"wizard":   {"burning hands", "disguise self"},
	"cleric":   {"guidance", "sacred flame", "etherealness"},
	"druid":    {"shillelagh", "thorn whip"},
	"bard":     {"vicious mockery", "dancing lights"},
	"sorcerer": {"fire bolt", "light"},
	"warlock":  {"eldritch blast", "mage hand"},
	"paladin":  {"divine sense", "lay on hands"},
}

var SpellLevels = map[string]int{
	"guidance":        0,
	"sacred flame":    0,
	"etherealness":    7,
	"burning hands":   1,
	"disguise self":   1,
	"shillelagh":      0,
	"thorn whip":      0,
	"vicious mockery": 0,
	"dancing lights":  0,
	"fire bolt":       0,
	"light":           0,
	"eldritch blast":  0,
	"mage hand":       0,
	"divine sense":    0,
	"lay on hands":    0,
}

func GiveStartingSpells(character *models.Character) error {
	spells, ok := StartingSpells[character.Class]
	if !ok {
		return nil
	}

	for _, name := range spells {
		level := 0
		spell := models.Spell{
			Name:     name,
			Level:    level,
			Prepared: false,
		}
		character.Spells = append(character.Spells, spell)
	}

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

func PrepareSpell(characterName, spellName string, spellLevel int) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return err
	}
	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf(`character "%s" not found`, characterName)
	}

	if !SpellcastingClasses[character.Class] {
		return fmt.Errorf("this class can't cast spells")
	}

	if !character.CanPrepareSpells {
		return fmt.Errorf("this class learns spells and can't prepare them")
	}

	var spellIndex int = -1
	for i, s := range character.Spells {
		if s.Name == spellName {
			spellIndex = i
			break
		}
	}
	if spellIndex == -1 {
		return fmt.Errorf("spell '%s' not known by character '%s'", spellName, characterName)
	}

	requiredLevel := character.Spells[spellIndex].Level
	if lvl, ok := SpellLevels[spellName]; ok {
		requiredLevel = lvl
	}

	if spellLevel < requiredLevel {
		return fmt.Errorf("the spell has higher level than the available spell slots")
	}

	if slots, ok := character.SpellSlots[spellLevel]; !ok || slots == 0 {
		return fmt.Errorf("no available spell slots of level %d", spellLevel)
	}

	character.Spells[spellIndex].Prepared = true
	character.Spells[spellIndex].Level = spellLevel

	if err := storage.SaveCharacter(character); err != nil {
		return err
	}

	fmt.Printf("Prepared spell %s\n", spellName)
	return nil
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
