package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// Bigram represents a bigram with its frequency
type Bigram struct {
	Text      string
	Frequency int
}

// Profile represents a language profile
type Profile struct {
	Bigrams []string
}

// LanguageResult stores detection result
type LanguageResult struct {
	Language string
	Distance int
}

// TrainProfiles builds language profiles from training directory
func TrainProfiles(trainDir string) (map[string]Profile, error) {
	profiles := make(map[string]Profile)

	err := filepath.WalkDir(trainDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".txt" {
			return nil
		}

		// Use filename without extension as language name
		lang := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %v", path, err)
		}

		profile := BuildProfileFromText(string(content))
		profiles[lang] = profile

		log.Printf("Trained profile for %s (%d bigrams)", lang, len(profile.Bigrams))
		return nil
	})

	return profiles, err
}

// BuildProfileFromText builds a language profile from text
func BuildProfileFromText(text string) Profile {
	// Count bigram frequencies
	bigramFreq := make(map[string]int)
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := normalizeWord(scanner.Text())
		if word == "" {
			continue
		}

		// Add start/end markers and convert to runes for proper Unicode handling
		wordWithMarkers := "_" + word + "_"
		runes := []rune(wordWithMarkers)

		// Generate bigrams
		for i := 0; i < len(runes)-1; i++ {
			bigram := string(runes[i : i+2])
			bigramFreq[bigram]++
		}
	}

	// Convert to slice and sort
	bigrams := make([]Bigram, 0, len(bigramFreq))
	for bigram, freq := range bigramFreq {
		bigrams = append(bigrams, Bigram{Text: bigram, Frequency: freq})
	}

	// Sort by frequency (descending) and then alphabetically
	sort.Slice(bigrams, func(i, j int) bool {
		if bigrams[i].Frequency == bigrams[j].Frequency {
			return bigrams[i].Text < bigrams[j].Text
		}
		return bigrams[i].Frequency > bigrams[j].Frequency
	})

	// Extract just the bigram texts
	profile := Profile{Bigrams: make([]string, len(bigrams))}
	for i, bg := range bigrams {
		profile.Bigrams[i] = bg.Text
	}

	return profile
}

// normalizeWord normalizes a word by converting to lowercase and removing punctuation
func normalizeWord(word string) string {
	// Convert to lowercase
	word = strings.ToLower(word)

	// Remove punctuation from start and end
	word = strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	return word
}

// CalculateDistance calculates the distance between two profiles
func CalculateDistance(langProfile, textProfile Profile) int {
	distance := 0

	for textIdx, textBigram := range textProfile.Bigrams {
		langIdx := findBigramIndex(langProfile.Bigrams, textBigram)
		if langIdx == -1 {
			distance += 1000
		} else {
			distance += abs(langIdx - textIdx)
		}
	}

	return distance
}

// findBigramIndex finds the index of a bigram in a profile
func findBigramIndex(bigrams []string, bigram string) int {
	for i, bg := range bigrams {
		if bg == bigram {
			return i
		}
	}
	return -1
}

// abs returns absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// SortResults sorts language results by distance
func SortResults(results []LanguageResult) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})
}

// SaveProfiles saves language profiles to a file
func SaveProfiles(filename string, profiles map[string]Profile) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(profiles)
}

// LoadProfiles loads language profiles from a file
func LoadProfiles(filename string) (map[string]Profile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var profiles map[string]Profile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&profiles)
	return profiles, err
}
