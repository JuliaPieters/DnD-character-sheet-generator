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
  %s create -name CHARACTER_NAME -player PLAYER_NAME -race RACE -class CLASS -background BACKGROUND -level N -str N -dex N -con N -int N -wis N -cha N -skills "Skill1,Skill2"
  %s view -name CHARACTER_NAME
  %s list
  %s delete -name CHARACTER_NAME
  %s update-level -name CHARACTER_NAME -level NEW_LEVEL
  %s add-weapon -name CHARACTER_NAME -weapon "Weapon Name" -category "Category" -range "Range" -two-handed true/false
  %s remove-weapon -name CHARACTER_NAME -weapon "Weapon Name"
  %s add-armor -name CHARACTER_NAME -armor "Armor Name" -armor-class N -dex-bonus true/false -max-dex N
  %s remove-armor -name CHARACTER_NAME
  %s add-shield -name CHARACTER_NAME -shield "Shield Name" -armor-class N
  %s remove-shield -name CHARACTER_NAME
  %s add-spell -name CHARACTER_NAME -spell "Spell Name" -level N -school "School" -range "Range"
  %s remove-spell -name CHARACTER_NAME -spell "Spell Name"
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0],
		os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0],
		os.Args[0], os.Args[0], os.Args[0], os.Args[0])
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
		background := createCmd.String("background", "", "Background")
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
			fmt.Println("❌ Character Name is required")
			createCmd.Usage()
			os.Exit(2)
		}

		var skillProficiencies []string
		if *skillsFlag != "" {
			skillProficiencies = strings.Split(*skillsFlag, ",")
			for i := range skillProficiencies {
				skillProficiencies[i] = strings.TrimSpace(skillProficiencies[i])
			}
		}

		abilityScores := []int{*strength, *dexterity, *constitution, *intelligence, *wisdom, *charisma}

		if err := commands.CreateCharacter(*characterName, *playerName, *characterRace, *characterClass, *background, *level, abilityScores, skillProficiencies); err != nil {
			fmt.Println("❌ Error creating character:", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Character '%s' created successfully!\n", *characterName)

	// ---------------- VIEW CHARACTER ----------------
	case "view":
		viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
		characterName := viewCmd.String("name", "", "Character Name (required)")
		_ = viewCmd.Parse(os.Args[2:])
		if *characterName == "" {
			fmt.Println("❌ Character Name is required")
			os.Exit(2)
		}
		if err := commands.ViewCharacter(*characterName); err != nil {
			fmt.Println("❌ Error viewing character:", err)
			os.Exit(1)
		}

	// ---------------- LIST CHARACTERS ----------------
	case "list":
		if err := commands.ListCharacters(); err != nil {
			fmt.Println("❌ Error listing characters:", err)
			os.Exit(1)
		}

	// ---------------- DELETE CHARACTER ----------------
	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		characterName := deleteCmd.String("name", "", "Character Name (required)")
		_ = deleteCmd.Parse(os.Args[2:])
		if *characterName == "" {
			fmt.Println("❌ Character Name is required")
			os.Exit(2)
		}
		if err := commands.DeleteCharacter(*characterName); err != nil {
			fmt.Println("❌ Error deleting character:", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Character '%s' deleted successfully!\n", *characterName)

	// ---------------- UPDATE LEVEL ----------------
	case "update-level":
		updateLevelCmd := flag.NewFlagSet("update-level", flag.ExitOnError)
		characterName := updateLevelCmd.String("name", "", "Character Name (required)")
		newLevel := updateLevelCmd.Int("level", 1, "New Level")
		_ = updateLevelCmd.Parse(os.Args[2:])
		if *characterName == "" {
			fmt.Println("❌ Character Name is required")
			os.Exit(2)
		}
		if err := commands.UpdateCharacterLevel(*characterName, *newLevel); err != nil {
			fmt.Println("❌ Error updating character level:", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Character '%s' updated to level %d successfully!\n", *characterName, *newLevel)

	// ---------------- ADD / REMOVE WEAPON ----------------
	case "add-weapon":
		addWeaponCmd := flag.NewFlagSet("add-weapon", flag.ExitOnError)
		characterName := addWeaponCmd.String("name", "", "Character Name")
		weaponName := addWeaponCmd.String("weapon", "", "Weapon Name")
		category := addWeaponCmd.String("category", "", "Weapon Category")
		weaponRange := addWeaponCmd.String("range", "", "Weapon Range")
		twoHanded := addWeaponCmd.Bool("two-handed", false, "Two-Handed Weapon")
		_ = addWeaponCmd.Parse(os.Args[2:])
		if *characterName == "" || *weaponName == "" {
			fmt.Println("❌ Character Name and Weapon Name are required")
			os.Exit(2)
		}
		newWeapon := models.Weapon{
			Name:      *weaponName,
			Category:  *category,
			Range:     *weaponRange,
			TwoHanded: *twoHanded,
		}
		if err := commands.VoegWapenToe(*characterName, newWeapon); err != nil {
			fmt.Println("❌ Error adding weapon:", err)
			os.Exit(1)
		}

	case "remove-weapon":
		removeWeaponCmd := flag.NewFlagSet("remove-weapon", flag.ExitOnError)
		characterName := removeWeaponCmd.String("name", "", "Character Name")
		weaponName := removeWeaponCmd.String("weapon", "", "Weapon Name")
		_ = removeWeaponCmd.Parse(os.Args[2:])
		if *characterName == "" || *weaponName == "" {
			fmt.Println("❌ Character Name and Weapon Name are required")
			os.Exit(2)
		}
		if err := commands.VerwijderWapen(*characterName, *weaponName); err != nil {
			fmt.Println("❌ Error removing weapon:", err)
			os.Exit(1)
		}

	// ---------------- ADD / REMOVE ARMOR ----------------
	case "add-armor":
		addArmorCmd := flag.NewFlagSet("add-armor", flag.ExitOnError)
		characterName := addArmorCmd.String("name", "", "Character Name")
		armorName := addArmorCmd.String("armor", "", "Armor Name")
		armorClass := addArmorCmd.Int("armor-class", 0, "Armor Class")
		dexBonus := addArmorCmd.Bool("dex-bonus", false, "Dexterity Bonus applies")
		maxDexBonus := addArmorCmd.Int("max-dex", 0, "Maximum Dexterity Bonus")
		_ = addArmorCmd.Parse(os.Args[2:])
		if *characterName == "" || *armorName == "" {
			fmt.Println("❌ Character Name and Armor Name are required")
			os.Exit(2)
		}
		newArmor := models.Armor{
			Name:        *armorName,
			ArmorClass:  *armorClass,
			DexBonus:    *dexBonus,
			MaxDexBonus: *maxDexBonus,
		}
		if err := commands.VoegArmorToe(*characterName, newArmor); err != nil {
			fmt.Println("❌ Error adding armor:", err)
			os.Exit(1)
		}

	case "remove-armor":
		removeArmorCmd := flag.NewFlagSet("remove-armor", flag.ExitOnError)
		characterName := removeArmorCmd.String("name", "", "Character Name")
		_ = removeArmorCmd.Parse(os.Args[2:])
		if *characterName == "" {
			fmt.Println("❌ Character Name is required")
			os.Exit(2)
		}
		if err := commands.VerwijderArmor(*characterName); err != nil {
			fmt.Println("❌ Error removing armor:", err)
			os.Exit(1)
		}

	// ---------------- ADD / REMOVE SHIELD ----------------
	case "add-shield":
		addShieldCmd := flag.NewFlagSet("add-shield", flag.ExitOnError)
		characterName := addShieldCmd.String("name", "", "Character Name")
		shieldName := addShieldCmd.String("shield", "", "Shield Name")
		armorClass := addShieldCmd.Int("armor-class", 0, "Armor Class")
		_ = addShieldCmd.Parse(os.Args[2:])
		if *characterName == "" || *shieldName == "" {
			fmt.Println("❌ Character Name and Shield Name are required")
			os.Exit(2)
		}
		newShield := models.Shield{
			Name:       *shieldName,
			ArmorClass: *armorClass,
		}
		if err := commands.VoegShieldToe(*characterName, newShield); err != nil {
			fmt.Println("❌ Error adding shield:", err)
			os.Exit(1)
		}

	case "remove-shield":
		removeShieldCmd := flag.NewFlagSet("remove-shield", flag.ExitOnError)
		characterName := removeShieldCmd.String("name", "", "Character Name")
		_ = removeShieldCmd.Parse(os.Args[2:])
		if *characterName == "" {
			fmt.Println("❌ Character Name is required")
			os.Exit(2)
		}
		if err := commands.VerwijderShield(*characterName); err != nil {
			fmt.Println("❌ Error removing shield:", err)
			os.Exit(1)
		}

	// ---------------- ADD / REMOVE SPELL ----------------
	case "add-spell":
		addSpellCmd := flag.NewFlagSet("add-spell", flag.ExitOnError)
		characterName := addSpellCmd.String("name", "", "Character Name")
		spellName := addSpellCmd.String("spell", "", "Spell Name")
		spellLevel := addSpellCmd.Int("level", 1, "Spell Level")
		spellSchool := addSpellCmd.String("school", "", "Spell School")
		spellRange := addSpellCmd.String("range", "", "Spell Range")
		_ = addSpellCmd.Parse(os.Args[2:])
		if *characterName == "" || *spellName == "" {
			fmt.Println("❌ Character Name and Spell Name are required")
			os.Exit(2)
		}
		newSpell := models.Spell{
			Name:   *spellName,
			Level:  *spellLevel,
			School: *spellSchool,
			Range:  *spellRange,
		}
		if err := commands.VoegSpellToe(*characterName, newSpell); err != nil {
			fmt.Println("❌ Error adding spell:", err)
			os.Exit(1)
		}

	case "remove-spell":
		removeSpellCmd := flag.NewFlagSet("remove-spell", flag.ExitOnError)
		characterName := removeSpellCmd.String("name", "", "Character Name")
		spellName := removeSpellCmd.String("spell", "", "Spell Name")
		_ = removeSpellCmd.Parse(os.Args[2:])
		if *characterName == "" || *spellName == "" {
			fmt.Println("❌ Character Name and Spell Name are required")
			os.Exit(2)
		}
		if err := commands.VerwijderSpell(*characterName, *spellName); err != nil {
			fmt.Println("❌ Error removing spell:", err)
			os.Exit(1)
		}

	// ---------------- DEFAULT ----------------
	default:
		printUsage()
		os.Exit(2)
	}
}
