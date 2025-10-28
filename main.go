package main

import (
	"dnd-character-sheet/commands"
	"dnd-character-sheet/domain"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	CharacterNameFlag     = "Character Name (required)"
	CharacterNameRequired = "character name is required"
	CharacterAndSpellReq  = "character name and spell name are required"
)

// -------------------- HELPERS --------------------

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

func parseNameFlag(cmd *flag.FlagSet) *string {
	name := cmd.String("name", "", CharacterNameFlag)
	_ = cmd.Parse(os.Args[2:])
	if *name == "" {
		fmt.Println(CharacterNameRequired)
		cmd.Usage()
		os.Exit(2)
	}
	return name
}

func parseSpellFlags(cmd *flag.FlagSet) (name, spell *string) {
	name = cmd.String("name", "", CharacterNameFlag)
	spell = cmd.String("spell", "", "Spell Name")
	_ = cmd.Parse(os.Args[2:])
	if *name == "" || *spell == "" {
		fmt.Println(CharacterAndSpellReq)
		os.Exit(2)
	}
	return
}

func parseEquipFlags(cmd *flag.FlagSet) (name, weapon, armor, shield, slot *string) {
	name = cmd.String("name", "", CharacterNameFlag)
	weapon = cmd.String("weapon", "", "Weapon Name")
	armor = cmd.String("armor", "", "Armor Name")
	shield = cmd.String("shield", "", "Shield Name")
	slot = cmd.String("slot", "", "Weapon Slot (main hand / off hand)")
	_ = cmd.Parse(os.Args[2:])
	return
}

// -------------------- HANDLERS --------------------

func handleCreate() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	characterName := createCmd.String("name", "", CharacterNameFlag)
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
		fmt.Println(CharacterNameRequired)
		createCmd.Usage()
		os.Exit(2)
	}

	var skillProficiencies []string
	if *skillsFlag != "" {
		for _, s := range strings.Split(*skillsFlag, ",") {
			skillProficiencies = append(skillProficiencies, strings.TrimSpace(s))
		}
	} else if skills, ok := domain.ClassSkills[strings.ToLower(*characterClass)]; ok {
		skillProficiencies = append(skillProficiencies, skills...)
	}

	abilityScores := []int{*strength, *dexterity, *constitution, *intelligence, *wisdom, *charisma}

	params := commands.CreateCharacterParams{
		CharacterName:      *characterName,
		PlayerName:         *playerName,
		Race:               *characterRace,
		Class:              *characterClass,
		Background:         *background,
		Level:              *level,
		AbilityScores:      abilityScores,
		SkillProficiencies: skillProficiencies,
	}

	if err := commands.CreateCharacter(params); err != nil {
		fmt.Printf(`failed to save character "%s"`+"\n", *characterName)
		os.Exit(1)
	}
	fmt.Printf("saved character %s\n", *characterName)
}

func handleView() {
	viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
	name := parseNameFlag(viewCmd)
	if err := commands.ViewCharacter(*name); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleList() {
	if err := commands.ListCharacters(); err != nil {
		fmt.Println("failed to list characters")
		os.Exit(1)
	}
}

func handleDelete() {
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	name := parseNameFlag(deleteCmd)
	if err := commands.DeleteCharacter(*name); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("deleted %s\n", *name)
}

func handleEquip() {
	equipCmd := flag.NewFlagSet("equip", flag.ExitOnError)
	name, weapon, armor, shield, slot := parseEquipFlags(equipCmd)

	if *weapon != "" {
		w, ok := commands.Weapons[strings.ToLower(*weapon)]
		if !ok {
			fmt.Printf("Weapon '%s' not found in CSV\n", *weapon)
			os.Exit(1)
		}
		var hand string
		var err error
		if *slot == "" {
			hand, err = commands.AddWeapon(*name, w)
		} else {
			hand, err = commands.AddWeaponToSlot(*name, w, *slot)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Equipped weapon %s to %s\n", w.Name, hand)
		return
	}

	if *armor != "" {
		a, ok := commands.Armors[strings.ToLower(*armor)]
		if !ok {
			fmt.Printf("Armor '%s' not found in CSV\n", *armor)
			os.Exit(1)
		}
		if err := commands.AddArmor(*name, a.Name); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}

	if *shield != "" {
		s, ok := commands.Shields[strings.ToLower(*shield)]
		if !ok {
			fmt.Printf("Shield '%s' not found in CSV\n", *shield)
			os.Exit(1)
		}
		if err := commands.AddShield(*name, s.Name); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Equipped shield %s\n", s.Name)
		return
	}

	fmt.Println("You must provide either -weapon, -armor or -shield")
	os.Exit(2)
}

func handleLearnSpell() {
	learnCmd := flag.NewFlagSet("learn-spell", flag.ExitOnError)
	name, spell := parseSpellFlags(learnCmd)
	if err := commands.LearnSpell(*name, *spell); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handlePrepareSpell() {
	prepareCmd := flag.NewFlagSet("prepare-spell", flag.ExitOnError)
	name, spell := parseSpellFlags(prepareCmd)
	level := prepareCmd.Int("level", 1, "Spell Level")
	_ = prepareCmd.Parse(os.Args[2:])
	if err := commands.PrepareSpell(*name, *spell, *level); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleEnrich() {
	enrichCmd := flag.NewFlagSet("enrich", flag.ExitOnError)
	name := parseNameFlag(enrichCmd)
	if err := commands.EnrichCharacter(*name); err != nil {
		fmt.Println("failed to enrich character:", err)
		os.Exit(1)
	}
	fmt.Printf("Enriched character %s with API data\n", *name)
}

// -------------------- MAIN --------------------

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

	switch os.Args[1] {
	case "create":
		handleCreate()
	case "view":
		handleView()
	case "list":
		handleList()
	case "delete":
		handleDelete()
	case "equip":
		handleEquip()
	case "learn-spell":
		handleLearnSpell()
	case "prepare-spell":
		handlePrepareSpell()
	case "enrich":
		handleEnrich()
	default:
		printUsage()
		os.Exit(2)
	}
}
