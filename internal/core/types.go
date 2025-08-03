package core

type Config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, args ...string) error
}

type LocationAreas struct {
	Count    int      `json:"count"`
	Next     any      `json:"next"`     // NOTE: This can be null. If you are at the last page, the value is null since there are no "next" pages
	Previous any      `json:"previous"` // NOTE: This can be null. If you are at the first page, the value is null since there are no "previous" pages
	Results  []Detail `json:"results"`
}

type LocationEncounterDetails struct {
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	GameIndex            int                   `json:"game_index"`
	ID                   int                   `json:"id"`
	Location             Detail                `json:"location"`
	Name                 string                `json:"name"`
	Names                []NameAndLanguage     `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
}

type EncounterMethodRate struct {
	EncounterMethod Detail                             `json:"encounter_method"`
	VersionDetails  []EncounterMethodRateVersionDetail `json:"version_details"`
}

type EncounterMethodRateVersionDetail struct {
	Rate    int    `json:"rate"`
	Version Detail `json:"version"`
}

type NameAndLanguage struct {
	Language Detail `json:"language"`
	Name     string `json:"name"`
}

type PokemonEncounter struct {
	Pokemon        Detail                          `json:"pokemon"`
	VersionDetails []PokemonEncounterVersionDetail `json:"version_details"`
}

type PokemonEncounterVersionDetail struct {
	EncounterDetails []EncounterDetail `json:"encounter_details"`
	MaxChance        int               `json:"max_chance"`
	Version          Detail            `json:"version"`
}

type EncounterDetail struct {
	Chance          int    `json:"chance"`
	ConditionValues []any  `json:"condition_values"`
	MaxLevel        int    `json:"max_level"`
	Method          Detail `json:"method"`
	MinLevel        int    `json:"min_level"`
}

type Detail struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
