package main

import (
	"dnd-character-sheet/commands"
	"flag"
	"fmt"
	"os"
	"strings"
)

func usage() {
	fmt.Printf(`Usage:
  %s create -name CHARACTER_NAME -race RACE -class CLASS -background BACKGROUND -level N -str N -dex N -con N -int N -wis N -cha N -skills "Skill1,Skill2"
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
		usage()
		os.Exit(1)
	}
	cmd := os.Args[1]

	switch cmd {
	case "create":
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		name := createCmd.String("name", "", "character name (required)")
		race := createCmd.String("race", "", "character race")
		class := createCmd.String("class", "", "character class")
		background := createCmd.String("background", "", "character background")
		level := createCmd.Int("level", 1, "character level")
		str := createCmd.Int("str", 10, "Strength")
		dex := createCmd.Int("dex", 10, "Dexterity")
		con := createCmd.Int("con", 10, "Constitution")
		intt := createCmd.Int("int", 10, "Intelligence")
		wis := createCmd.Int("wis", 10, "Wisdom")
		cha := createCmd.Int("cha", 10, "Charisma")
		skillsFlag := createCmd.String("skills", "", "comma-separated list of skill proficiencies")

		err := createCmd.Parse(os.Args[2:])
		if *name == "" || err != nil {
			fmt.Println("❌ name is required")
			createCmd.Usage()
			os.Exit(2)
		}

		// Parse skills
		var skills []string
		if *skillsFlag != "" {
			skills = strings.Split(*skillsFlag, ",")
			for i := range skills {
				skills[i] = strings.TrimSpace(skills[i])
			}
		}

		if err := commands.CreateCharacter(
			*name, *race, *class, *background, *level,
			*str, *dex, *con, *intt, *wis, *cha,
			skills,
		); err != nil {
			fmt.Println("Error creating character:", err)
			os.Exit(1)
		}

	case "view":
		viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
		name := viewCmd.String("name", "", "character name (required)")
		viewCmd.Parse(os.Args[2:])
		if *name == "" {
			fmt.Println("❌ name is required")
			os.Exit(2)
		}
		if err := commands.ViewCharacter(*name); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	case "list":
		if err := commands.ListCharacters(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		name := deleteCmd.String("name", "", "character name (required)")
		deleteCmd.Parse(os.Args[2:])
		if *name == "" {
			fmt.Println("❌ name is required")
			os.Exit(2)
		}
		if err := commands.DeleteCharacter(*name); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	case "equip":
		fmt.Println("⚔️ Equip command not yet implemented")

	case "learn-spell":
		fmt.Println("✨ Learn-spell command not yet implemented")

	case "prepare-spell":
		fmt.Println("📖 Prepare-spell command not yet implemented")

	default:
		usage()
		os.Exit(2)
	}
}
