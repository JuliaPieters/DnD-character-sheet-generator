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
  %s create -name CHARACTER_NAME -player PLAYER_NAME -race RACE -class CLASS -background BACKGROUND -alignment ALIGNMENT -level N -str N -dex N -con N -int N -wis N -cha N -skills "Skill1,Skill2"
  %s view -name CHARACTER_NAME
  %s list
  %s delete -name CHARACTER_NAME
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0])
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

		// Character basic info
		name := createCmd.String("name", "", "Character Name (required)")
		playerName := createCmd.String("player", "", "Player Name")
		race := createCmd.String("race", "", "Race")
		className := createCmd.String("class", "", "Class")
		background := createCmd.String("background", "", "Background")
		alignment := createCmd.String("alignment", "", "Alignment")
		level := createCmd.Int("level", 1, "Level")

		// Ability Scores
		str := createCmd.Int("str", 10, "Strength")
		dex := createCmd.Int("dex", 10, "Dexterity")
		con := createCmd.Int("con", 10, "Constitution")
		intt := createCmd.Int("int", 10, "Intelligence")
		wis := createCmd.Int("wis", 10, "Wisdom")
		cha := createCmd.Int("cha", 10, "Charisma")

		// Skills
		skillsFlag := createCmd.String("skills", "", "Comma-separated skill list")

		err := createCmd.Parse(os.Args[2:])
		if err != nil || *name == "" {
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

		// Call your existing CreateCharacter function
		if err := commands.CreateCharacter(
			*name,
			*playerName,
			*race,
			*className,
			*background,
			*alignment,
			*level,
			*str,
			*dex,
			*con,
			*intt,
			*wis,
			*cha,
			skills,
		); err != nil {
			fmt.Println("Error creating character:", err)
			os.Exit(1)
		}

	case "view":
		viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
		name := viewCmd.String("name", "", "Character Name (required)")
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
		name := deleteCmd.String("name", "", "Character Name (required)")
		deleteCmd.Parse(os.Args[2:])
		if *name == "" {
			fmt.Println("❌ name is required")
			os.Exit(2)
		}
		if err := commands.DeleteCharacter(*name); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	default:
		usage()
		os.Exit(2)
	}
}
