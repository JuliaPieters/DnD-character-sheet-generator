package api

import (
	"dnd-character-sheet/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type APIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type APIListResponse struct {
	Results []APIResource `json:"results"`
}

type APISpell struct {
	Name   string `json:"name"`
	Level  int    `json:"level"`
	School struct {
		Name string `json:"name"`
	} `json:"school"`
	Range   string `json:"range"`
	Classes []struct {
		Name string `json:"name"`
	} `json:"classes"`
}

type EquipmentRange struct {
	Normal int `json:"normal,omitempty"`
	Long   int `json:"long,omitempty"`
}

type APIEquipment struct {
	Name              string `json:"name"`
	EquipmentCategory struct {
		Name string `json:"name"`
	} `json:"equipment_category"`
	ArmorClass struct {
		Base     int  `json:"base"`
		DexBonus bool `json:"dex_bonus"`
		MaxDex   int  `json:"max_bonus"`
	} `json:"armor_class,omitempty"`
	TwoHanded bool            `json:"two_handed,omitempty"`
	Range     json.RawMessage `json:"range,omitempty"`
}

func getJSON(url string, target interface{}) error {
	url = strings.ToLower(strings.ReplaceAll(url, " ", "-"))

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("fout %d bij ophalen van %s: %s", resp.StatusCode, url, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !json.Valid(body) {
		return fmt.Errorf("ongeldige JSON van %s: %s", url, string(body[:min(100, len(body))]))
	}

	return json.Unmarshal(body, target)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GetSpellsForClass(className string, slots map[int]int) ([]models.Spell, error) {
	var list APIListResponse
	if err := getJSON("https://www.dnd5eapi.co/api/spells", &list); err != nil {
		return nil, err
	}

	type SpellResult struct {
		Spell models.Spell
		Err   error
	}

	results := make(chan SpellResult)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	classNameLower := strings.ToLower(className)

	for _, res := range list.Results {
		<-ticker.C
		go func(res APIResource) {
			url := "https://www.dnd5eapi.co" + strings.ToLower(strings.ReplaceAll(res.URL, " ", "-"))
			var spell APISpell
			err := getJSON(url, &spell)
			if err != nil {
				results <- SpellResult{Err: fmt.Errorf("error fetching %s: %w", res.Name, err)}
				return
			}

			isForClass := false
			for _, c := range spell.Classes {
				if strings.ToLower(c.Name) == classNameLower {
					isForClass = true
					break
				}
			}
			if !isForClass {
				results <- SpellResult{Err: nil}
				return
			}

			results <- SpellResult{
				Spell: models.Spell{
					Name:   spell.Name,
					Level:  spell.Level,
					School: spell.School.Name,
					Range:  spell.Range,
				},
			}
		}(res)
	}

	selected := []models.Spell{}
	for i := 0; i < len(list.Results); i++ {
		result := <-results
		if result.Err != nil {
			log.Println(result.Err)
			continue
		}
		if result.Spell.Name == "" {
			continue
		}
		if _, ok := slots[result.Spell.Level]; ok {
			selected = append(selected, result.Spell)
		}
	}

	rand.Seed(time.Now().UnixNano())
	final := []models.Spell{}
	for lvl, count := range slots {
		lvlSpells := []models.Spell{}
		for _, s := range selected {
			if s.Level == lvl {
				lvlSpells = append(lvlSpells, s)
			}
		}
		if len(lvlSpells) > count {
			rand.Shuffle(len(lvlSpells), func(i, j int) { lvlSpells[i], lvlSpells[j] = lvlSpells[j], lvlSpells[i] })
			final = append(final, lvlSpells[:count]...)
		} else {
			final = append(final, lvlSpells...)
		}
	}

	return final, nil
}

func parseRange(raw json.RawMessage) string {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}

	var r EquipmentRange
	if err := json.Unmarshal(raw, &r); err == nil {
		if r.Normal > 0 {
			return fmt.Sprintf("%d", r.Normal)
		}
	}
	return ""
}

func GetEquipment() (*models.Weapon, *models.Weapon, *models.Armor, *models.Shield, error) {
	var list APIListResponse
	if err := getJSON("https://www.dnd5eapi.co/api/equipment", &list); err != nil {
		return nil, nil, nil, nil, err
	}

	var mainHand *models.Weapon
	var offHand *models.Weapon
	var armor *models.Armor
	var shield *models.Shield

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for _, res := range list.Results {
		<-ticker.C
		var eq APIEquipment
		url := "https://www.dnd5eapi.co" + strings.ToLower(strings.ReplaceAll(res.URL, " ", "-"))
		if err := getJSON(url, &eq); err != nil {
			log.Println("Error fetching equipment:", res.Name, err)
			continue
		}

		switch eq.EquipmentCategory.Name {
		case "Weapon":
			weapon := &models.Weapon{
				Name:      eq.Name,
				TwoHanded: eq.TwoHanded,
				Range:     parseRange(eq.Range),
			}
			if mainHand == nil {
				mainHand = weapon
			} else if offHand == nil {
				offHand = weapon
			}
		case "Armor":
			armor = &models.Armor{
				Name:        eq.Name,
				ArmorClass:  eq.ArmorClass.Base,
				DexBonus:    eq.ArmorClass.DexBonus,
				MaxDexBonus: eq.ArmorClass.MaxDex,
			}
		case "Shield":
			shield = &models.Shield{
				Name:       eq.Name,
				ArmorClass: eq.ArmorClass.Base,
			}
		}
	}

	return mainHand, offHand, armor, shield, nil
}
