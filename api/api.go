package api

import (
	"dnd-character-sheet/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type APIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type APIListResponse struct {
	Results []APIResource `json:"results"`
}

// Voor individuele spell/equipment details:
type APISpell struct {
	Name   string `json:"name"`
	Level  int    `json:"level"`
	School struct {
		Name string `json:"name"`
	} `json:"school"`
	Range string `json:"range"`
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
	TwoHanded bool   `json:"two_handed,omitempty"`
	Range     string `json:"range,omitempty"`
}

func getJSON(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}

func GetSpellsForClass(className string, slots map[int]int) ([]models.Spell, error) {
	var list APIListResponse
	err := getJSON("https://www.dnd5eapi.co/api/spells", &list)
	if err != nil {
		return nil, err
	}

	selected := []models.Spell{}

	for _, res := range list.Results {
		var spell APISpell
		err := getJSON("https://www.dnd5eapi.co"+res.URL, &spell)
		if err != nil {
			log.Println("Error fetching spell:", res.Name, err)
			continue
		}

		// Alleen spells die passen bij een beschikbare level
		if _, ok := slots[spell.Level]; ok {
			selected = append(selected, models.Spell{
				Name:   spell.Name,
				Level:  spell.Level,
			})
		}
	}

	// Willekeurig kiezen uit geselecteerde spells per slot
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

// GetEquipment haalt equipment op en verdeelt weapons over main hand en off hand
func GetEquipment() (*models.Weapon, *models.Weapon, *models.Armor, *models.Shield, error) {
	var list APIListResponse
	err := getJSON("https://www.dnd5eapi.co/api/equipment", &list)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var mainHand *models.Weapon
	var offHand *models.Weapon
	var armor *models.Armor
	var shield *models.Shield

	for _, res := range list.Results {
		var eq APIEquipment
		err := getJSON("https://www.dnd5eapi.co"+res.URL, &eq)
		if err != nil {
			log.Println("Error fetching equipment:", res.Name, err)
			continue
		}

		switch eq.EquipmentCategory.Name {
		case "Weapon":
			weapon := &models.Weapon{
				Name:      eq.Name,
				TwoHanded: eq.TwoHanded,
				Range:     eq.Range,
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
