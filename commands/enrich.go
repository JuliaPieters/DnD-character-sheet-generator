package commands

import (
	"dnd-character-sheet/api"
	"dnd-character-sheet/application"
	"dnd-character-sheet/domain"
	"dnd-character-sheet/storage"
	"fmt"
	"log"
)

func EnrichCharacter(name string) error {
	char, err := storage.GetCharacterByName(name)
	if err != nil {
		return fmt.Errorf("failed to load character: %w", err)
	}

	addSpells(char)
	addEquipment(char)

	if err := storage.SaveCharacter(char); err != nil {
		return fmt.Errorf("failed to save enriched character: %w", err)
	}

	fmt.Println("Character enriched successfully!")
	return nil
}

func addSpells(char *domain.Character) {
	if char.Level == 0 || len(char.SpellSlots) == 0 {
		return
	}

	spells, err := api.GetSpellsForClass(char.Class, char.SpellSlots)
	if err != nil {
		log.Println("failed to get spells:", err)
		return
	}

	char.Spells = spells
}

func addEquipment(char *domain.Character) {
	mainHand, offHand, armor, shield, err := api.GetEquipment()
	if err != nil {
		log.Println("failed to get equipment:", err)
		return
	}

	equipService := application.EquipmentService{}

	if mainHand != nil {
		char.Equipment.MainHand = mainHand
		char.Equipment.MainHand.Damage = equipService.CalculateWeaponDamage(char, mainHand)
	}
	if offHand != nil {
		char.Equipment.OffHand = offHand
		char.Equipment.OffHand.Damage = equipService.CalculateWeaponDamage(char, offHand)
	}
	if armor != nil {
		char.Equipment.Armor = armor
	}
	if shield != nil {
		char.Equipment.Shield = shield
	}
}
