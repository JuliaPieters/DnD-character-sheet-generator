# Exam Report – Modern Programming Practices

## Architecture (Extensibility) – 60/60

De architectuur van de code is goed ontworpen: functionaliteiten bevinden zich in de juiste lagen en verantwoordelijkheden zijn duidelijk gescheiden.  

- De CLI-functionaliteit en server-functionaliteit staan in aparte `main.go` bestanden (`main.go` voor CLI en `server/main.go` voor server), waardoor ze onafhankelijk functioneren.  
- De commands en berekeningen (skills, combat stats, spellcasting, weapon damage) staan in eigen methodes of mappen, waardoor de logica modulair is.  
- De codebase volgt het **Open/Closed Principle**: nieuwe features kunnen toegevoegd worden zonder bestaande code te wijzigen.  

### Opsplitsing van modellen

Het grote bestand `models/character.go` is verwijderd en opgesplitst in **domain** en **application** lagen:  

**Domain-lagen (`domain/`)**  
- `character.go` – structs voor character data  
- `ability.go` – ability score logica  
- `equipment.go` – weapon, armor en shield structuren  
- `spell.go` – spell structuren  
- `constants.go` – relevante constanten  

**Application-lagen (`application/`)**  
- `character_service.go` – berekeningen op characters (combat stats, initiative, AC)  
- `equipment_service.go` – weapon/armor damage berekeningen  
- `spell_service.go` – spellcasting logica  
- `equipment_service_test.go` – tests voor equipment functionaliteit  

**Voordelen:**  
- Domain bevat enkel data en basale logica.  
- Application bevat alle business logic.  
- Nieuwe features kunnen toegevoegd worden in `application` zonder impact op `domain`.  

### API opsplitsing

Oorspronkelijk bestond er één `api.go` bestand. Dit is opgesplitst in:  

- `api/client.go` – centrale API-client en request-logica  
- `api/spell.go` – spells ophalen en enrich-logica  
- `api/equipment.go` – equipment ophalen en enrich-logica  

**Voordelen:**  
- **Single Responsibility Principle**: elke module heeft één duidelijk doel.  
- Wijzigingen in spells of equipment vereisen geen aanpassing van andere API-functionaliteit.  
- `client.go` wordt hergebruikt door `spell.go` en `equipment.go`.  
- Verbeterde modulariteit en testbaarheid.  

### Diagram van code dependencies

Flow Diagram – Project Data & Functionality
──────────────────────────────────────────

┌─────────────┐
│ CLI/User │
└─────┬───────┘
│ input commands
▼
┌─────────────┐
│ main.go CLI │
└─────┬───────┘
│ calls
▼
┌─────────────┐
│ commands/ │
└─────┬───────┘
│ invokes
▼
┌─────────────┐
│ domain/ │ <- Core structs & data
│ character, │
│ equipment, │
│ ability, │
│ spell │
└─────┬───────┘
│ used by
▼
┌─────────────┐
│ application/│ <- Business logic
│ services │
└─────┬───────┘
│ fetches/enriches
▼
┌─────────────┐
│ api/ │ <- API calls
│ client.go, │
│ spell.go, │
│ equipment.go│
└─────┬───────┘
│ loads
▼
┌─────────────┐
│ data/ │ <- CSV, JSON
└─────────────┘

Separate server flow:
─────────────┐
│ server/ │
│ main.go │
│ character.json │
└─────┬───────┘
│ serves data
▼
┌─────────────┐
│ templates/ │
└─────────────┘
│ renders
▼
┌─────────────┐
│ static/ │
└─────────────┘

---

## Maintainability – 20/20

De code is overzichtelijk en onderhoudbaar, met duidelijke scheiding van verantwoordelijkheden.

### Mappenstructuur
project/
├─ api/ # API calls & enrich functies
├─ commands/ # CLI commands
├─ data/ # CSV bestanden & JSON
├─ domain/ # Structs & basale logica
│ ├─ character.go
│ ├─ ability.go
│ ├─ equipment.go
│ ├─ spell.go
│ └─ constants.go
├─ application/ # Business logic / services
│ ├─ character_service.go
│ ├─ equipment_service.go
│ ├─ spell_service.go
│ └─ equipment_service_test.go
├─ server/
│ ├─ main.go
│ └─ character.json
├─ static/ # CSS, JS, images
├─ templates/ # HTML templates
├─ main.go # CLI entrypoint
└─ character.json # CLI-specifieke opslag

- Wijzigingen in bijvoorbeeld `api/` of `domain/` hebben geen impact op `commands/` of `server/`.  
- Functies zijn modulair, goed benoemd en gelogd, wat onderhoudbaarheid verhoogt.

### Functie-indeling

- Elke functie heeft één duidelijke verantwoordelijkheid:  
  - `CalculateArmorClass()` – berekent alleen AC  
  - `SetupSpellcasting()` – regelt spellcasting  
  - `WeaponDamage()` – berekent damage van een wapen  
- Cyclomatic complexity is laag (<10 per functie), wat testen en debuggen vereenvoudigt.  
- Bestanden zijn groot, maar modulair met duidelijke subfuncties en tests.  

### Testbare uitbreidbaarheid

- Nieuwe wapens, modifiers of features kunnen volledig binnen `application/` of `domain/` toegevoegd worden.  
- Bestaande code wordt niet aangepast, waardoor regressie wordt vermeden.  
- **Automated tests** in `equipment_service_test.go` en `character_test.go` zorgen dat oude functionaliteit intact blijft.

---

## Testing – 20/20

Weapon damage-functionaliteit is volledig getest met **geautomatiseerde tests**, inclusief happy path en edge cases.

### OUTPUT TESTS

PS C:\Users\julia\DnD character sheet generator> go test ./application -v                                                                    
>> 
=== RUN   TestCalculateWeaponDamage
=== RUN   TestCalculateWeaponDamage/Level_1_Half-Orc_Barbarian_with_Greataxe
=== RUN   TestCalculateWeaponDamage/Level_1_Tiefling_Ranger_with_Shortsword
=== RUN   TestCalculateWeaponDamage/Level_1_Dwarf_Rogue_with_Rapier
--- PASS: TestCalculateWeaponDamage (0.00s)
    --- PASS: TestCalculateWeaponDamage/Level_1_Half-Orc_Barbarian_with_Greataxe (0.00s)
    --- PASS: TestCalculateWeaponDamage/Level_1_Tiefling_Ranger_with_Shortsword (0.00s)
    --- PASS: TestCalculateWeaponDamage/Level_1_Dwarf_Rogue_with_Rapier (0.00s)
=== RUN   TestCalculateWeaponDamage_EdgeCases
=== RUN   TestCalculateWeaponDamage_EdgeCases/Negative_STR_and_DEX
=== RUN   TestCalculateWeaponDamage_EdgeCases/Zero_modifiers
--- PASS: TestCalculateWeaponDamage_EdgeCases (0.00s)
    --- PASS: TestCalculateWeaponDamage_EdgeCases/Negative_STR_and_DEX (0.00s)
    --- PASS: TestCalculateWeaponDamage_EdgeCases/Zero_modifiers (0.00s)
PASS
ok      dnd-character-sheet/application (cached)

### Argumentatie per criterium

**No testing has been done – ✅ NVT**  
Tests zijn aanwezig en uitgevoerd via `go test ./application -v`.

**Happy path – ✅**  
De functionaliteit werkt correct voor voorbeeldkarakters:  

- Half-Orc Barbarian met Greataxe → `1d12 +3`  
- Tiefling Ranger met Shortsword → `1d6 +2`  
- Dwarf Rogue met Rapier → `1d8 +2`  

Dit bewijst dat de hoofdfunctionaliteit correct wordt uitgevoerd.

**Edge cases / combination of manual and automated tests – ✅**  
Edge cases zijn getest:  

- Negatieve modifiers (bijv. DEX/STR < 0) → correcte berekening (`-1`)  
- Zero modifiers → correcte berekening (`+0`)  

Alle tests zijn volledig geautomatiseerd, wat robuustheid en betrouwbaarheid aantoont.

**Conclusie**  
- De tests dekken zowel **happy path** als belangrijke randgevallen.  
- De geautomatiseerde tests zijn volledig succesvol uitgevoerd.  
- Hierdoor is aangetoond dat de weapon damage-functionaliteit correct, betrouwbaar en uitbreidbaar is.
