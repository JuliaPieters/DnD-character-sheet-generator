package commands

import (
	"dnd-character-sheet/models"
	"dnd-character-sheet/storage"
	"fmt"
	"math"
	"strings"
)

// Spell slot tabel voor full casters (level 1-10)
var FullCasterSlots = map[int]map[int]int{
	1:  {1: 2},
	2:  {1: 3},
	3:  {1: 4, 2: 2},
	4:  {1: 4, 2: 3},
	5:  {1: 4, 2: 3, 3: 2},
	6:  {1: 4, 2: 3, 3: 3},
	7:  {1: 4, 2: 3, 3: 3, 4: 1},
	8:  {1: 4, 2: 3, 3: 3, 4: 2},
	9:  {1: 4, 2: 3, 3: 3, 4: 3, 5: 1},
	10: {1: 4, 2: 3, 3: 3, 4: 3, 5: 2},
}

// HalfCasters: paladin & ranger
var HalfCasters = map[string]bool{
	"paladin": true,
	"ranger":  true,
}

func ViewCharacter(name string) error {
	characters, err := storage.LoadCharacters()
	if err != nil {
		return err
	}

	for _, c := range characters {
		if c.Name != name {
			continue
		}

		fmt.Printf("Name: %s\n", c.Name)
		fmt.Printf("Class: %s\n", strings.ToLower(c.Class))
		fmt.Printf("Race: %s\n", strings.ToLower(c.Race))
		fmt.Printf("Background: %s\n", strings.ToLower(c.Background))
		fmt.Printf("Level: %d\n", c.Level)

		fmt.Printf("Ability scores:\n")
		fmt.Printf("  STR: %d (%+d)\n", c.Abilities.Strength, abilityModifier(c.Abilities.Strength))
		fmt.Printf("  DEX: %d (%+d)\n", c.Abilities.Dexterity, abilityModifier(c.Abilities.Dexterity))
		fmt.Printf("  CON: %d (%+d)\n", c.Abilities.Constitution, abilityModifier(c.Abilities.Constitution))
		fmt.Printf("  INT: %d (%+d)\n", c.Abilities.Intelligence, abilityModifier(c.Abilities.Intelligence))
		fmt.Printf("  WIS: %d (%+d)\n", c.Abilities.Wisdom, abilityModifier(c.Abilities.Wisdom))
		fmt.Printf("  CHA: %d (%+d)\n", c.Abilities.Charisma, abilityModifier(c.Abilities.Charisma))

		fmt.Printf("Proficiency bonus: %+d\n", proficiencyBonus(c.Level))
		fmt.Printf("Skill proficiencies: %s\n", formatSkillProficiencies(c.SkillProficiencies))

		if c.Equipment.MainHand != nil {
			fmt.Printf("Main hand: %s\n", c.Equipment.MainHand.Name)
		}
		if c.Equipment.OffHand != nil {
			fmt.Printf("Off hand: %s\n", c.Equipment.OffHand.Name)
		}
		if c.Equipment.Armor != nil {
			fmt.Printf("Armor: %s\n", c.Equipment.Armor.Name)
		}
		if c.Equipment.Shield != nil {
			fmt.Printf("Shield: %s\n", c.Equipment.Shield.Name)
		}

		// Spell slots voor full casters
		if _, isCaster := models.SpellcastingClasses[c.Class]; isCaster {
			fmt.Println("Spell slots:")

			// Check half casters
			if HalfCasters[c.Class] {
				if slots, ok := FullCasterSlots[c.Level]; ok {
					for lvl := 1; lvl <= 9; lvl++ {
						if n, exists := slots[lvl]; exists && n > 0 {
							half := (n + 1) / 2 // helft van full caster slots, afgerond omhoog
							fmt.Printf("  Level %d: %d\n", lvl, half)
						}
					}
				}
			} else {
				// Full caster
				if slots, ok := FullCasterSlots[c.Level]; ok {
					for lvl := 1; lvl <= 9; lvl++ {
						if n, exists := slots[lvl]; exists && n > 0 {
							fmt.Printf("  Level %d: %d\n", lvl, n)
						}
					}
				}
			}
		}

		return nil
	}

	return fmt.Errorf(`character "%s" not found`, name)
}

func abilityModifier(score int) int {
	return int(math.Floor(float64(score-10) / 2))
}

func proficiencyBonus(level int) int {
	return 2 + (level-1)/4
}

func formatSkillProficiencies(skills []string) string {
	for i := range skills {
		skills[i] = strings.ToLower(skills[i])
	}
	return strings.Join(skills, ", ")
}
