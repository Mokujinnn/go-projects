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

type Bigram struct {
	Text      string
	Frequency int
}

type Profile struct {
	Bigrams []string
}

type LanguageResult struct {
	Language string
	Distance int
}

func TrainProfiles(trainDir string) (map[string]Profile, error) {
	profiles := make(map[string]Profile)

	err := filepath.WalkDir(trainDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".txt" {
			return nil
		}

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

func BuildProfileFromText(text string) Profile {
	bigramFreq := make(map[string]int)
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := normalizeWord(scanner.Text())
		if word == "" {
			continue
		}

		wordWithMarkers := "_" + word + "_"
		runes := []rune(wordWithMarkers)

		for i := 0; i < len(runes)-1; i++ {
			bigram := string(runes[i : i+2])
			bigramFreq[bigram]++
		}
	}

	bigrams := make([]Bigram, 0, len(bigramFreq))
	for bigram, freq := range bigramFreq {
		bigrams = append(bigrams, Bigram{Text: bigram, Frequency: freq})
	}

	sort.Slice(bigrams, func(i, j int) bool {
		if bigrams[i].Frequency == bigrams[j].Frequency {
			return bigrams[i].Text < bigrams[j].Text
		}
		return bigrams[i].Frequency > bigrams[j].Frequency
	})

	profile := Profile{Bigrams: make([]string, len(bigrams))}
	for i, bg := range bigrams {
		profile.Bigrams[i] = bg.Text
	}

	return profile
}

func normalizeWord(word string) string {
	word = strings.ToLower(word)

	word = strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	return word
}

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

func findBigramIndex(bigrams []string, bigram string) int {
	for i, bg := range bigrams {
		if bg == bigram {
			return i
		}
	}
	return -1
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func SortResults(results []LanguageResult) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})
}

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
