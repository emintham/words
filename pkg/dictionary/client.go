package dictionary

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/words-api/words/internal/models"
)

const dictionaryAPIURL = "https://api.dictionaryapi.dev/api/v2/entries/en"

// Client handles fetching word definitions from the external API
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new dictionary API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchWord fetches word definition from the external API
func (c *Client) FetchWord(word string) (*models.Word, error) {
	url := fmt.Sprintf("%s/%s", dictionaryAPIURL, strings.ToLower(word))

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch word: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("word not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp models.DictionaryAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResp) == 0 {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert API response to our Word model
	return convertAPIResponse(apiResp[0]), nil
}

func convertAPIResponse(apiWord struct {
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
}) *models.Word {
	word := &models.Word{
		Word:       apiWord.Word,
		Phonetic:   apiWord.Phonetic,
		SourceUrls: apiWord.SourceUrls,
	}

	// Convert phonetics
	for _, p := range apiWord.Phonetics {
		word.Phonetics = append(word.Phonetics, models.Phonetic{
			Text:  p.Text,
			Audio: p.Audio,
		})
	}

	// Convert meanings
	for _, m := range apiWord.Meanings {
		meaning := models.Meaning{
			PartOfSpeech: m.PartOfSpeech,
			Synonyms:     m.Synonyms,
			Antonyms:     m.Antonyms,
		}

		// Convert definitions
		for _, d := range m.Definitions {
			definition := models.Definition{
				Definition: d.Definition,
				Example:    d.Example,
				Synonyms:   d.Synonyms,
				Antonyms:   d.Antonyms,
			}
			meaning.Definitions = append(meaning.Definitions, definition)
		}

		word.Meanings = append(word.Meanings, meaning)
	}

	return word
}
