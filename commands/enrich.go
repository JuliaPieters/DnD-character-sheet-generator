package commands

import (
	"dnd-character-sheet/api"
	"dnd-character-sheet/application"
	"dnd-character-sheet/domain"
	"dnd-character-sheet/storage"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

func EnrichCharacter(name string) error {
	char, err := storage.GetCharacterByName(name)
	if err != nil {
		return fmt.Errorf("failed to load character: %w", err)
	}

	addSpells(char)
	mergeEquipment(char)

	if err := storage.SaveCharacter(char); err != nil {
		return fmt.Errorf("failed to save enriched character: %w", err)
	}

	fmt.Printf("Enriched character %s with API data\n", char.Name)
	return nil
}

func addSpells(char *domain.Character) {
	if char.Level == 0 || len(char.SpellSlots) == 0 {
		return
	}
	if len(char.Spells) > 0 {
		return
	}

	spells, err := api.GetSpellsForClass(char.Class, char.SpellSlots)
	if err != nil {
		log.Println("failed to get spells:", err)
		return
	}

	char.Spells = spells
}

func mergeEquipment(char *domain.Character) {
	equipService := application.EquipmentService{}

	allWeapons, armorAPI, shieldAPI, err := api.GetAllEquipment()
	if err != nil {
		log.Println("failed to get equipment:", err)
		return
	}

	rand.Seed(time.Now().UnixNano())

	if char.Equipment.MainHand == nil && len(allWeapons) > 0 {
		randomIndex := rand.Intn(len(allWeapons))
		char.Equipment.MainHand = allWeapons[randomIndex]
	} else if char.Equipment.MainHand != nil {
		fillWeaponData(char.Equipment.MainHand, allWeapons)
	}
	if char.Equipment.MainHand != nil {
		char.Equipment.MainHand.Damage = equipService.CalculateWeaponDamage(char, char.Equipment.MainHand)
	}

	if char.Equipment.OffHand == nil && len(allWeapons) > 1 {
		randomIndex := rand.Intn(len(allWeapons))
		char.Equipment.OffHand = allWeapons[randomIndex]
	} else if char.Equipment.OffHand != nil {
		fillWeaponData(char.Equipment.OffHand, allWeapons)
	}
	if char.Equipment.OffHand != nil {
		char.Equipment.OffHand.Damage = equipService.CalculateWeaponDamage(char, char.Equipment.OffHand)
	}

	if char.Equipment.Armor == nil && armorAPI != nil {
		char.Equipment.Armor = armorAPI
	}

	if char.Equipment.Shield == nil && shieldAPI != nil {
		char.Equipment.Shield = shieldAPI
	}
}


func fillWeaponData(existing *domain.Weapon, allWeapons []*domain.Weapon) {
	for _, w := range allWeapons {
		if strings.EqualFold(w.Name, existing.Name) {
			existing.Category = w.Category

			if existing.DamageDie == "" && w.DamageDie != "" {
				existing.DamageDie = w.DamageDie
			}
			if existing.Range == "" && w.Range != "" {
				existing.Range = w.Range
			}
			if !existing.TwoHanded && w.TwoHanded {
				existing.TwoHanded = true
			}
			if !existing.IsFinesse && w.IsFinesse {
				existing.IsFinesse = true
			}
			break
		}
	}
}

