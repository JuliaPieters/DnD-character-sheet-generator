package main

import (
	"dnd-character-sheet/commands"
	"dnd-character-sheet/models"
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
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func main() {
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
			if skills, ok := models.ClassSkills[classKey]; ok {
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
			fmt.Printf(`character "%s" not found`+"\n", *characterName)
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
			fmt.Printf(`character "%s" not found`+"\n", *characterName)
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
		category := equipCmd.String("category", "", "Weapon Category")
		weaponRange := equipCmd.String("range", "", "Weapon Range")
		twoHanded := equipCmd.Bool("two-handed", false, "Two-Handed Weapon")
		armorClass := equipCmd.Int("armor-class", 0, "Armor Class")
		dexBonus := equipCmd.Bool("dex-bonus", false, "Dexterity Bonus applies")
		maxDexBonus := equipCmd.Int("max-dex", 0, "Maximum Dexterity Bonus")
		_ = equipCmd.Parse(os.Args[2:])

		if *characterName == "" {
			fmt.Println("character name is required")
			os.Exit(2)
		}

		if *weaponName != "" {
			newWeapon := models.Weapon{
				Name:      *weaponName,
				Category:  *category,
				Range:     *weaponRange,
				TwoHanded: *twoHanded,
			}
			var hand string
			var err error
			if *slot == "" {
				hand, err = commands.AddWeapon(*characterName, newWeapon)
			} else {
				hand, err = commands.AddWeaponToSlot(*characterName, newWeapon, *slot)
			}
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("Equipped weapon %s to %s\n", *weaponName, hand)
			return
		}

		if *armorName != "" {
			newArmor := models.Armor{
				Name:        *armorName,
				ArmorClass:  *armorClass,
				DexBonus:    *dexBonus,
				MaxDexBonus: *maxDexBonus,
			}
			if err := commands.AddArmor(*characterName, newArmor); err != nil {
				fmt.Printf(`character "%s" not found`+"\n", *characterName)
				os.Exit(1)
			}
			fmt.Printf("Equipped armor %s\n", *armorName)
			return
		}

		if *shieldName != "" {
			newShield := models.Shield{
				Name:       *shieldName,
				ArmorClass: *armorClass,
			}
			if err := commands.AddShield(*characterName, newShield); err != nil {
				fmt.Printf(`character "%s" not found`+"\n", *characterName)
				os.Exit(1)
			}
			fmt.Printf("Equipped shield %s\n", *shieldName)
			return
		}

		fmt.Println("You must provide either -weapon, -armor or -shield")
		os.Exit(2)

	// ---------------- LEARN SPELL ----------------
	case "learn-spell":
		learnCmd := flag.NewFlagSet("learn-spell", flag.ExitOnError)
		characterName := learnCmd.String("name", "", "Character Name")
		spellName := learnCmd.String("spell", "", "Spell Name")
		spellLevel := learnCmd.Int("level", 1, "Spell Level")
		_ = learnCmd.Parse(os.Args[2:])

		if *characterName == "" || *spellName == "" {
			fmt.Println("character name and spell name are required")
			os.Exit(2)
		}

		newSpell := models.Spell{
			Name:  *spellName,
			Level: *spellLevel,
		}
		if err := commands.LearnSpell(*characterName, newSpell); err != nil {
			fmt.Printf(`character "%s" not found`+"\n", *characterName)
			os.Exit(1)
		}

	// ---------------- PREPARE SPELL ----------------
	case "prepare-spell":
		prepareCmd := flag.NewFlagSet("prepare-spell", flag.ExitOnError)
		characterName := prepareCmd.String("name", "", "Character Name")
		spellName := prepareCmd.String("spell", "", "Spell Name")
		_ = prepareCmd.Parse(os.Args[2:])

		if *characterName == "" || *spellName == "" {
			fmt.Println("character name and spell name are required")
			os.Exit(2)
		}

		if err := commands.PrepareSpell(*characterName, *spellName); err != nil {
			fmt.Printf(`could not prepare spell "%s" for character "%s"`+"\n", *spellName, *characterName)
			os.Exit(1)
		}
		fmt.Printf("prepared spell %s for %s\n", *spellName, *characterName)

	// ---------------- DEFAULT ----------------
	default:
		printUsage()
		os.Exit(2)
	}
}
