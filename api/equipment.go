package api

import (
	"dnd-character-sheet/models"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// ---------- Structs ----------

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

// ---------- Helpers ----------

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

// ---------- Logic ----------

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
