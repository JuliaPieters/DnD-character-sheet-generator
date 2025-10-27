package commands

import (
	"dnd-character-sheet/api"
	"dnd-character-sheet/storage"
	"fmt"
	"log"
)

func EnrichCharacter(name string) error {
	char, err := storage.GetCharacterByName(name)
	if err != nil {
		return fmt.Errorf("failed to load character: %w", err)
	}

	if char.Level > 0 && len(char.SpellSlots) > 0 {
		spells, err := api.GetSpellsForClass(char.Class, char.SpellSlots)
		if err != nil {
			log.Println("failed to get spells:", err)
		} else {
			char.Spells = spells
		}
	}

	mainHand, offHand, armor, shield, err := api.GetEquipment()
	if err != nil {
		log.Println("failed to get equipment:", err)
	} else {
		if mainHand != nil {
			char.Equipment.MainHand = mainHand
		}
		if offHand != nil {
			char.Equipment.OffHand = offHand
		}
		if armor != nil {
			char.Equipment.Armor = armor
		}
		if shield != nil {
			char.Equipment.Shield = shield
		}
	}

	if err := storage.SaveCharacter(char); err != nil {
		return fmt.Errorf("failed to save enriched character: %w", err)
	}

	fmt.Println("Character enriched successfully!")
	return nil
}
