package commands

import (
	"dnd-character-sheet/application"
	"dnd-character-sheet/storage"
	"fmt"
	"sort"
	"strings"
)

func ViewCharacter(name string) error {
	charPtr, err := storage.GetCharacterByName(name)
	if err != nil {
		return fmt.Errorf(`character "%s" not found`, name)
	}

	charService := application.CharacterService{}
	charService.CalculateCombatStats(charPtr)

	fmt.Printf("Name: %s\n", charPtr.Name)
	fmt.Printf("Class: %s\n", strings.ToLower(charPtr.Class))
	fmt.Printf("Race: %s\n", strings.ToLower(charPtr.Race))
	fmt.Printf("Background: %s\n", strings.ToLower(charPtr.Background))
	fmt.Printf("Level: %d\n", charPtr.Level)

	fmt.Println("Ability scores:")
	fmt.Printf("  STR: %d (%+d)\n", charPtr.Abilities.Strength, charPtr.Abilities.Modifier("Strength"))
	fmt.Printf("  DEX: %d (%+d)\n", charPtr.Abilities.Dexterity, charPtr.Abilities.Modifier("Dexterity"))
	fmt.Printf("  CON: %d (%+d)\n", charPtr.Abilities.Constitution, charPtr.Abilities.Modifier("Constitution"))
	fmt.Printf("  INT: %d (%+d)\n", charPtr.Abilities.Intelligence, charPtr.Abilities.Modifier("Intelligence"))
	fmt.Printf("  WIS: %d (%+d)\n", charPtr.Abilities.Wisdom, charPtr.Abilities.Modifier("Wisdom"))
	fmt.Printf("  CHA: %d (%+d)\n", charPtr.Abilities.Charisma, charPtr.Abilities.Modifier("Charisma"))

	fmt.Printf("Proficiency bonus: %+d\n", charPtr.ProficiencyBonus)
	fmt.Printf("Skill proficiencies: %s\n", formatSkillProficiencies(charPtr.SkillProficiencies))

	if charPtr.Equipment.MainHand != nil {
		fmt.Printf("Main hand: %s\n", charPtr.Equipment.MainHand.Name)
	}
	if charPtr.Equipment.OffHand != nil {
		fmt.Printf("Off hand: %s\n", charPtr.Equipment.OffHand.Name)
	}
	if charPtr.Equipment.Armor != nil {
		fmt.Printf("Armor: %s\n", charPtr.Equipment.Armor.Name)
	}
	if charPtr.Equipment.Shield != nil {
		fmt.Printf("Shield: %s\n", charPtr.Equipment.Shield.Name)
	}

	if len(charPtr.SpellSlots) > 0 {
		fmt.Println("Spell slots:")
		levels := make([]int, 0, len(charPtr.SpellSlots))
		for lvl := range charPtr.SpellSlots {
			levels = append(levels, lvl)
		}
		sort.Ints(levels)
		for _, lvl := range levels {
			fmt.Printf("  Level %d: %d\n", lvl, charPtr.SpellSlots[lvl])
		}
	}

	// Gebruik de geëxporteerde FullCasters en PactCasters uit spell.go
	if FullCasters[strings.ToLower(charPtr.Class)] || PactCasters[strings.ToLower(charPtr.Class)] {
		if charPtr.SpellcastingAbility != "" {
			fmt.Printf("Spellcasting ability: %s\n", strings.ToLower(charPtr.SpellcastingAbility))
			fmt.Printf("Spell save DC: %d\n", charPtr.SpellSaveDC)
			fmt.Printf("Spell attack bonus: %+d\n", charPtr.SpellAttackBonus)
		}
	}

	fmt.Printf("Armor class: %d\n", charPtr.ArmorClass)
	fmt.Printf("Initiative bonus: %d\n", charPtr.Initiative)
	fmt.Printf("Passive perception: %d\n", charPtr.PassivePerception)

	return nil
}

func formatSkillProficiencies(skills []string) string {
	for i := range skills {
		skills[i] = strings.ToLower(skills[i])
	}
	return strings.Join(skills, ", ")
}
