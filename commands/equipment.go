package commands

import (
	"dnd-character-sheet/application"
	"dnd-character-sheet/domain"
	"dnd-character-sheet/storage"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

var Armors = map[string]domain.Armor{}
var Shields = map[string]domain.Shield{}
var Weapons = map[string]domain.Weapon{}

var DefaultArmorStats = map[string]domain.Armor{
	// Light
	"padded":          {ArmorClass: 11, DexBonus: true, MaxDexBonus: 0},
	"leather":         {ArmorClass: 11, DexBonus: true, MaxDexBonus: 0},
	"studded leather": {ArmorClass: 12, DexBonus: true, MaxDexBonus: 0},

	// Medium
	"hide":        {ArmorClass: 12, DexBonus: true, MaxDexBonus: 2},
	"chain shirt": {ArmorClass: 13, DexBonus: true, MaxDexBonus: 2},
	"scale mail":  {ArmorClass: 14, DexBonus: true, MaxDexBonus: 2},
	"breastplate": {ArmorClass: 14, DexBonus: true, MaxDexBonus: 2},
	"half plate":  {ArmorClass: 15, DexBonus: true, MaxDexBonus: 2},

	// Heavy
	"ring mail":  {ArmorClass: 14, DexBonus: false},
	"chain mail": {ArmorClass: 16, DexBonus: false},
	"splint":     {ArmorClass: 17, DexBonus: false},
	"plate":      {ArmorClass: 18, DexBonus: false},

	// Shields
	"shield": {ArmorClass: 2, DexBonus: false},
}

func normalizeName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.TrimSuffix(name, " armor")
	name = strings.TrimSpace(name)
	return name
}

// ------------------------
// CSV Loading
// ------------------------
func LoadEquipmentCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open equipment CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("could not read equipment CSV: %w", err)
	}

	for i, record := range records {
		if i == 0 {
			continue
		}

		originalName := strings.TrimSpace(record[0])
		eqType := strings.ToLower(strings.TrimSpace(record[1]))
		key := normalizeName(originalName)

		switch eqType {
		case "armor":
			if key == "shield" {
				stats := DefaultArmorStats["shield"]
				shield := domain.Shield{
					Name:       strings.ToLower(originalName),
					ArmorClass: stats.ArmorClass,
				}
				Shields[key] = shield
				Shields[strings.ToLower(originalName)] = shield
			} else {
				stats, ok := DefaultArmorStats[key]
				if !ok {
					stats = domain.Armor{ArmorClass: 10, DexBonus: true}
				}
				armor := domain.Armor{
					Name:        key,
					ArmorClass:  stats.ArmorClass,
					DexBonus:    stats.DexBonus,
					MaxDexBonus: stats.MaxDexBonus,
				}
				Armors[key] = armor
				Armors[strings.ToLower(originalName)] = armor
			}
		case "weapon":
			weapon := domain.Weapon{Name: strings.ToLower(originalName)}
			Weapons[key] = weapon
			Weapons[strings.ToLower(originalName)] = weapon
		}
	}

	return nil
}

// ------------------------
// Weapons
// ------------------------
func AddWeapon(characterName string, newWeapon domain.Weapon) (string, error) {
	return AddWeaponToSlot(characterName, newWeapon, "")
}

func AddWeaponToSlot(characterName string, newWeapon domain.Weapon, slot string) (string, error) {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return "", fmt.Errorf("could not load characters: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return "", fmt.Errorf("character '%s' not found", characterName)
	}
	characterPtr := &character

	newWeapon.Name = strings.ToLower(strings.TrimSpace(newWeapon.Name))
	var hand string

	switch slot {
	case "":
		if characterPtr.Equipment.MainHand == nil {
			characterPtr.Equipment.MainHand = &newWeapon
			hand = "main hand"
		} else if characterPtr.Equipment.OffHand == nil {
			characterPtr.Equipment.OffHand = &newWeapon
			hand = "off hand"
		} else {
			return "", fmt.Errorf("both hands already occupied")
		}
	case "main hand":
		if characterPtr.Equipment.MainHand != nil {
			return "", fmt.Errorf("main hand already occupied")
		}
		characterPtr.Equipment.MainHand = &newWeapon
		hand = "main hand"
	case "off hand":
		if characterPtr.Equipment.OffHand != nil {
			return "", fmt.Errorf("off hand already occupied")
		}
		characterPtr.Equipment.OffHand = &newWeapon
		hand = "off hand"
	default:
		return "", fmt.Errorf("invalid slot: must be 'main hand' or 'off hand'")
	}

	cs := application.CharacterService{}
	cs.CalculateCombatStats(characterPtr)

	if err := storage.SaveCharacter(characterPtr); err != nil {
		return "", fmt.Errorf("could not save character: %w", err)
	}

	return hand, nil
}

func RemoveWeapon(characterName, weaponName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("could not load characters: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}
	characterPtr := &character

	weaponName = normalizeName(weaponName)
	removed := false
	if characterPtr.Equipment.MainHand != nil && normalizeName(characterPtr.Equipment.MainHand.Name) == weaponName {
		characterPtr.Equipment.MainHand = nil
		removed = true
	}
	if characterPtr.Equipment.OffHand != nil && normalizeName(characterPtr.Equipment.OffHand.Name) == weaponName {
		characterPtr.Equipment.OffHand = nil
		removed = true
	}

	if !removed {
		return fmt.Errorf("weapon '%s' not found on character '%s'", weaponName, characterName)
	}

	cs := application.CharacterService{}
	cs.CalculateCombatStats(characterPtr)

	if err := storage.SaveCharacter(characterPtr); err != nil {
		return fmt.Errorf("could not save character: %w", err)
	}

	return nil
}

// ------------------------
// Armor
// ------------------------
func AddArmor(characterName, armorName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("could not load characters: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}
	characterPtr := &character

	key := strings.ToLower(strings.TrimSpace(armorName))
	armor, ok := Armors[key]
	if !ok {
		return fmt.Errorf("armor '%s' not found", armorName)
	}

	displayArmor := armor
	switch key {
	case "padded", "leather", "studded leather", "plate":
		displayArmor.Name = key + " armor"
	default:
		displayArmor.Name = key
	}

	characterPtr.Equipment.Armor = &displayArmor

	cs := application.CharacterService{}
	cs.CalculateCombatStats(characterPtr)

	if err := storage.SaveCharacter(characterPtr); err != nil {
		return fmt.Errorf("could not save character: %w", err)
	}

	fmt.Printf("Equipped armor %s\n", displayArmor.Name)
	return nil
}

func RemoveArmor(characterName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("could not load characters: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}
	characterPtr := &character

	characterPtr.Equipment.Armor = nil

	cs := application.CharacterService{}
	cs.CalculateCombatStats(characterPtr)

	if err := storage.SaveCharacter(characterPtr); err != nil {
		return fmt.Errorf("could not save character: %w", err)
	}

	return nil
}

// ------------------------
// Shield
// ------------------------
func AddShield(characterName, shieldName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("could not load characters: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}
	characterPtr := &character

	key := normalizeName(shieldName)
	shield, ok := Shields[key]
	if !ok {
		shield, ok = Shields[strings.ToLower(shieldName)]
		if !ok {
			return fmt.Errorf("shield '%s' not found", shieldName)
		}
	}

	displayShield := shield
	displayShield.Name = key
	characterPtr.Equipment.Shield = &displayShield

	cs := application.CharacterService{}
	cs.CalculateCombatStats(characterPtr)

	if err := storage.SaveCharacter(characterPtr); err != nil {
		return fmt.Errorf("could not save character: %w", err)
	}

	return nil
}

func RemoveShield(characterName string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return fmt.Errorf("could not load characters: %w", err)
	}

	character, exists := characters[characterName]
	if !exists {
		return fmt.Errorf("character '%s' not found", characterName)
	}
	characterPtr := &character

	characterPtr.Equipment.Shield = nil

	cs := application.CharacterService{}
	cs.CalculateCombatStats(characterPtr)

	if err := storage.SaveCharacter(characterPtr); err != nil {
		return fmt.Errorf("could not save character: %w", err)
	}

	return nil
}
