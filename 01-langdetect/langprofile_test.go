package main

import (
	"os"
	"testing"
)

func TestNormalizeWord(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello!", "hello"},
		{"'world'", "world"},
		{"Тест.", "тест"},
		{"123abc", "123abc"},
		{"", ""},
	}

	for _, test := range tests {
		result := normalizeWord(test.input)
		if result != test.expected {
			t.Errorf("normalizeWord(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestBuildProfileFromText(t *testing.T) {
	text := "Deep sleep"
	profile := BuildProfileFromText(text)

	expectedBigrams := []string{"ee", "ep", "p_", "_d", "_s", "de", "le", "sl"}
	if len(profile.Bigrams) != len(expectedBigrams) {
		t.Errorf("Expected %d bigrams, got %d", len(expectedBigrams), len(profile.Bigrams))
	}

	for i, expected := range expectedBigrams {
		if i < len(profile.Bigrams) && profile.Bigrams[i] != expected {
			t.Errorf("Bigram at position %d: expected %q, got %q", i, expected, profile.Bigrams[i])
		}
	}
}

func TestCalculateDistance(t *testing.T) {
	langProfile := Profile{Bigrams: []string{"th", "in", "on", "er", "ed"}}
	textProfile := Profile{Bigrams: []string{"th", "er", "on", "le", "in"}}

	distance := CalculateDistance(langProfile, textProfile)
	expected := 1005

	if distance != expected {
		t.Errorf("Expected distance %d, got %d", expected, distance)
	}
}

func TestFindBigramIndex(t *testing.T) {
	bigrams := []string{"aa", "bb", "cc"}

	if findBigramIndex(bigrams, "bb") != 1 {
		t.Error("Failed to find existing bigram")
	}

	if findBigramIndex(bigrams, "dd") != -1 {
		t.Error("Found non-existing bigram")
	}
}

func TestProfileSerialization(t *testing.T) {
	profiles := map[string]Profile{
		"en": {Bigrams: []string{"th", "he", "in"}},
		"ru": {Bigrams: []string{"то", "на", "ов"}},
	}

	tempFile := "test_profiles.json"
	defer os.Remove(tempFile)

	if err := SaveProfiles(tempFile, profiles); err != nil {
		t.Fatalf("Failed to save profiles: %v", err)
	}

	loadedProfiles, err := LoadProfiles(tempFile)
	if err != nil {
		t.Fatalf("Failed to load profiles: %v", err)
	}

	if len(loadedProfiles) != len(profiles) {
		t.Errorf("Loaded %d profiles, expected %d", len(loadedProfiles), len(profiles))
	}
}
