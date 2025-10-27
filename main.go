package main

import (
	"dnd-character-sheet/commands"
	"dnd-character-sheet/domain"
	"flag"
	"fmt"
	"os"
	"strings"
)

func printUsage() {
	fmt.Printf(`Usage:
		 %s create -name CHARACTER_NAME -race RACE -class CLASS -level N -str N -dex N -con N -int N -wis N -cha N
		 %s view -name CHARACTER_NAME
		 %s list
		 %s delete -name CHARACTER_NAME
		 %s equip -name CHARACTER_NAME -weapon WEAPON_NAME -slot SLOT
		 %s equip -name CHARACTER_NAME -armor ARMOR_NAME
		 %s equip -name CHARACTER_NAME -shield SHIELD_NAME
		 %s learn-spell -name CHARACTER_NAME -spell SPELL_NAME
		 %s prepare-spell -name CHARACTER_NAME -spell SPELL_NAME
		 %s enrich -name CHARACTER_NAME
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func main() {
	if err := commands.LoadSpellsFromCSV("data/spells.csv"); err != nil {
		fmt.Println("failed to load spells:", err)
		os.Exit(1)
	}

	if err := commands.LoadEquipmentCSV("data/equipment.csv"); err != nil {
		fmt.Println("failed to load equipment:", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {

	// ---------------- CREATE CHARACTER ----------------
	case "create":
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		characterName := createCmd.String("name", "", "Character Name (required)")
		playerName := createCmd.String("player", "", "Player Name")
		characterRace := createCmd.String("race", "", "Race")
		characterClass := createCmd.String("class", "", "Class")
		background := createCmd.String("background", "acolyte", "Background")
		level := createCmd.Int("level", 1, "Level")
		strength := createCmd.Int("str", 10, "Strength")
		dexterity := createCmd.Int("dex", 10, "Dexterity")
		constitution := createCmd.Int("con", 10, "Constitution")
		intelligence := createCmd.Int("int", 10, "Intelligence")
		wisdom := createCmd.Int("wis", 10, "Wisdom")
		charisma := createCmd.Int("cha", 10, "Charisma")
		skillsFlag := createCmd.String("skills", "", "Comma-separated skill list")
		_ = createCmd.Parse(os.Args[2:])

		if *characterName == "" {
			fmt.Println("character name is required")
			createCmd.Usage()
			os.Exit(2)
		}

		var skillProficiencies []string
		if *skillsFlag != "" {
			skillProficiencies = strings.Split(*skillsFlag, ",")
			for i := range skillProficiencies {
				skillProficiencies[i] = strings.TrimSpace(skillProficiencies[i])
			}
		} else {
			classKey := strings.ToLower(*characterClass)
			if skills, ok := domain.ClassSkills[classKey]; ok {
				skillProficiencies = append(skillProficiencies, skills...)
			}
		}

		abilityScores := []int{*strength, *dexterity, *constitution, *intelligence, *wisdom, *charisma}

		if err := commands.CreateCharacter(*characterName, *playerName, *characterRace, *characterClass, *background, *level, abilityScores, skillProficiencies); err != nil {
			fmt.Printf(`failed to save character "%s"`+"\n", *characterName)
			os.Exit(1)
		}
		fmt.Printf("saved character %s\n", *characterName)

	// ---------------- VIEW CHARACTER ----------------
	case "view":
		viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
		characterName := viewCmd.String("name", "", "Character Name (required)")
		_ = viewCmd.Parse(os.Args[2:])
		if *characterName == "" {
			fmt.Println("character name is required")
			os.Exit(2)
		}
		if err := commands.ViewCharacter(*characterName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	// ---------------- LIST CHARACTERS ----------------
	case "list":
		if err := commands.ListCharacters(); err != nil {
			fmt.Println("failed to list characters")
			os.Exit(1)
		}

	// ---------------- DELETE CHARACTER ----------------
	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		characterName := deleteCmd.String("name", "", "Character Name (required)")
		_ = deleteCmd.Parse(os.Args[2:])
		if *characterName == "" {
			fmt.Println("character name is required")
			os.Exit(2)
		}
		if err := commands.DeleteCharacter(*characterName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("deleted %s\n", *characterName)

	// ---------------- EQUIP ----------------
	case "equip":
		equipCmd := flag.NewFlagSet("equip", flag.ExitOnError)
		characterName := equipCmd.String("name", "", "Character Name (required)")
		weaponName := equipCmd.String("weapon", "", "Weapon Name")
		armorName := equipCmd.String("armor", "", "Armor Name")
		shieldName := equipCmd.String("shield", "", "Shield Name")
		slot := equipCmd.String("slot", "", "Weapon Slot (main hand / off hand)")
		_ = equipCmd.Parse(os.Args[2:])

		if *characterName == "" {
			fmt.Println("character name is required")
			os.Exit(2)
		}

		if *weaponName != "" {
			weapon, ok := commands.Weapons[strings.ToLower(*weaponName)]
			if !ok {
				fmt.Printf("Weapon '%s' not found in CSV\n", *weaponName)
				os.Exit(1)
			}
			var hand string
			var err error
			if *slot == "" {
				hand, err = commands.AddWeapon(*characterName, weapon)
			} else {
				hand, err = commands.AddWeaponToSlot(*characterName, weapon, *slot)
			}
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("Equipped weapon %s to %s\n", weapon.Name, hand)
			return
		}

		if *armorName != "" {
			armor, ok := commands.Armors[strings.ToLower(*armorName)]
			if !ok {
				fmt.Printf("Armor '%s' not found in CSV\n", *armorName)
				os.Exit(1)
			}
			if err := commands.AddArmor(*characterName, armor.Name); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}

		if *shieldName != "" {
			shield, ok := commands.Shields[strings.ToLower(*shieldName)]
			if !ok {
				fmt.Printf("Shield '%s' not found in CSV\n", *shieldName)
				os.Exit(1)
			}
			if err := commands.AddShield(*characterName, shield.Name); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("Equipped shield %s\n", shield.Name)
			return
		}

		fmt.Println("You must provide either -weapon, -armor or -shield")
		os.Exit(2)

	// ---------------- LEARN SPELL ----------------
	case "learn-spell":
		learnCmd := flag.NewFlagSet("learn-spell", flag.ExitOnError)
		characterName := learnCmd.String("name", "", "Character Name")
		spellName := learnCmd.String("spell", "", "Spell Name")
		_ = learnCmd.Parse(os.Args[2:])
		if *characterName == "" || *spellName == "" {
			fmt.Println("character name and spell name are required")
			os.Exit(2)
		}
		if err := commands.LearnSpell(*characterName, *spellName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	// ---------------- PREPARE SPELL ----------------
	case "prepare-spell":
		prepareCmd := flag.NewFlagSet("prepare-spell", flag.ExitOnError)
		characterName := prepareCmd.String("name", "", "Character Name")
		spellName := prepareCmd.String("spell", "", "Spell Name")
		level := prepareCmd.Int("level", 1, "Spell Level")
		_ = prepareCmd.Parse(os.Args[2:])
		if *characterName == "" || *spellName == "" {
			fmt.Println("character name and spell name are required")
			os.Exit(2)
		}
		if err := commands.PrepareSpell(*characterName, *spellName, *level); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	// ---------------- ENRICH CHARACTER ----------------
	case "enrich":
		enrichCmd := flag.NewFlagSet("enrich", flag.ExitOnError)
		characterName := enrichCmd.String("name", "", "Character Name (required)")
		_ = enrichCmd.Parse(os.Args[2:])

		if *characterName == "" {
			fmt.Println("character name is required")
			enrichCmd.Usage()
			os.Exit(2)
		}

		if err := commands.EnrichCharacter(*characterName); err != nil {
			fmt.Println("failed to enrich character:", err)
			os.Exit(1)
		}

		fmt.Printf("Enriched character %s with API data\n", *characterName)

	// ---------------- DEFAULT ----------------
	default:
		printUsage()
		os.Exit(2)
	}
}
