package commands

import (
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
	"fmt"
	"strings"
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
	if spells, ok := StartingSpells[character.Class]; ok {
		for _, name := range spells {
			level := SpellLevels[name]
			character.Spells = append(character.Spells, models.Spell{
				Name:     name,
				Level:    level,
				Prepared: false,
			})
		}
	}
	SetupSpellcasting(character)
	return storage.SaveCharacter(*character)
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
	level := 1
	if lvl, ok := SpellLevels[spellName]; ok {
		level = lvl
	}
	character.Spells = append(character.Spells, models.Spell{
		Name:     spellName,
		Level:    level,
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

	spellIndex := -1
	for i, s := range character.Spells {
		if s.Name == spellName {
			spellIndex = i
			break
		}
	}
	if spellIndex == -1 {
		return fmt.Errorf("spell '%s' not known by character '%s'", spellName, characterName)
	}

	requiredLevel := SpellLevels[spellName]
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

func SetupSpellcasting(c *models.Character) {
	class := strings.ToLower(c.Class)
	if !SpellcastingClasses[class] {
		c.SpellSlots = nil
		c.CanPrepareSpells = false
		return
	}
	c.CanPrepareSpells = PreparedCasters[class]
	c.SpellSlots = GenerateSpellSlots(class, c.Level)
	if class == "warlock" {
		c.SpellSlots[0] = 4
	}
}

func GenerateSpellSlots(class string, level int) map[int]int {
	slots := make(map[int]int)
	class = strings.ToLower(class)

	fullCasters := map[string]bool{
		"wizard": true, "cleric": true, "druid": true, "bard": true, "sorcerer": true,
	}

	if fullCasters[class] {
		switch {
		case level <= 3:
			slots[0] = 3
		case level <= 9:
			slots[0] = 4
		default:
			slots[0] = 5
		}

		slotTable := map[int][]int{
			1:  {2, 0, 0, 0, 0, 0, 0, 0, 0},
			2:  {3, 0, 0, 0, 0, 0, 0, 0, 0},
			3:  {4, 2, 0, 0, 0, 0, 0, 0, 0},
			4:  {4, 3, 0, 0, 0, 0, 0, 0, 0},
			5:  {4, 3, 2, 0, 0, 0, 0, 0, 0},
			6:  {4, 3, 3, 0, 0, 0, 0, 0, 0},
			7:  {4, 3, 3, 1, 0, 0, 0, 0, 0},
			8:  {4, 3, 3, 2, 0, 0, 0, 0, 0},
			9:  {4, 3, 3, 3, 1, 0, 0, 0, 0},
			10: {4, 3, 3, 3, 2, 0, 0, 0, 0},
			11: {4, 3, 3, 3, 2, 1, 0, 0, 0},
			12: {4, 3, 3, 3, 2, 1, 1, 0, 0},
			13: {4, 3, 3, 3, 2, 1, 1, 1, 0},
			14: {4, 3, 3, 3, 2, 1, 1, 1, 1},
			15: {4, 3, 3, 3, 2, 1, 1, 1, 1},
			16: {4, 3, 3, 3, 3, 1, 1, 1, 1},
			17: {4, 3, 3, 3, 3, 2, 1, 1, 1},
			18: {4, 3, 3, 3, 3, 2, 2, 1, 1},
			19: {4, 3, 3, 3, 3, 2, 2, 2, 1},
			20: {4, 3, 3, 3, 3, 2, 2, 1, 1},
		}

		if l, ok := slotTable[level]; ok {
			for i, count := range l {
				if count > 0 {
					slots[i+1] = count
				}
			}
		}
		return slots
	}

	if class == "paladin" || class == "ranger" {
		halfCaster := []struct {
			slotLevel int
			reqLevel  int
			slots     int
		}{
			{1, 2, 4}, {2, 5, 3}, {3, 9, 3}, {4, 13, 3}, {5, 17, 2},
		}
		for _, h := range halfCaster {
			if level >= h.reqLevel {
				slots[h.slotLevel] = h.slots
			}
		}
		return slots
	}

	if class == "warlock" {
		switch {
		case level == 1:
			slots[1] = 1
		case level == 2:
			slots[1] = 2
		case level >= 3 && level <= 4:
			slots[2] = 2
		case level >= 5 && level <= 6:
			slots[3] = 2
		case level >= 7 && level <= 8:
			slots[4] = 2
		case level >= 9 && level <= 10:
			slots[5] = 2
		case level >= 11 && level <= 16:
			slots[5] = 3
		case level >= 17 && level <= 19:
			slots[5] = 4
		case level == 20:
			slots[5] = 4
			slots[6] = 1
			slots[7] = 1
			slots[8] = 1
			slots[9] = 1
		}
		return slots
	}

	return slots
}
