package models

import "time"

// Word represents a complete dictionary entry
type Word struct {
	ID          int64     `json:"id" db:"id"`
	Word        string    `json:"word" db:"word"`
	Phonetic    string    `json:"phonetic,omitempty" db:"phonetic"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Meanings    []Meaning `json:"meanings,omitempty"`
	Phonetics   []Phonetic `json:"phonetics,omitempty"`
	SourceUrls  []string  `json:"sourceUrls,omitempty"`
}

// Meaning represents a part of speech with its definitions
type Meaning struct {
	ID           int64        `json:"id,omitempty" db:"id"`
	WordID       int64        `json:"-" db:"word_id"`
	PartOfSpeech string       `json:"partOfSpeech" db:"part_of_speech"`
	Definitions  []Definition `json:"definitions,omitempty"`
	Synonyms     []string     `json:"synonyms,omitempty"`
	Antonyms     []string     `json:"antonyms,omitempty"`
}

// Definition represents a single definition
type Definition struct {
	ID         int64    `json:"id,omitempty" db:"id"`
	MeaningID  int64    `json:"-" db:"meaning_id"`
	Definition string   `json:"definition" db:"definition"`
	Example    string   `json:"example,omitempty" db:"example"`
	Synonyms   []string `json:"synonyms,omitempty"`
	Antonyms   []string `json:"antonyms,omitempty"`
}

// Phonetic represents pronunciation information
type Phonetic struct {
	ID       int64  `json:"id,omitempty" db:"id"`
	WordID   int64  `json:"-" db:"word_id"`
	Text     string `json:"text" db:"text"`
	Audio    string `json:"audio,omitempty" db:"audio"`
}

// DictionaryAPIResponse matches the structure from dictionaryapi.dev
type DictionaryAPIResponse []struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic"`
	Phonetics []struct {
		Text  string `json:"text"`
		Audio string `json:"audio"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string   `json:"definition"`
			Example    string   `json:"example"`
			Synonyms   []string `json:"synonyms"`
			Antonyms   []string `json:"antonyms"`
		} `json:"definitions"`
		Synonyms []string `json:"synonyms"`
		Antonyms []string `json:"antonyms"`
	} `json:"meanings"`
	License struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"license"`
	SourceUrls []string `json:"sourceUrls"`
}
