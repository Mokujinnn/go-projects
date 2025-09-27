package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	profileFile := flag.String("profile", "", "Path to language profile file")
	trainDir := flag.String("train", "", "Path to training directory")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	if *profileFile == "" {
		log.Fatal("Profile file path is required")
	}

	if *trainDir != "" {
		// Training mode
		profiles, err := TrainProfiles(*trainDir)
		if err != nil {
			log.Fatalf("Training failed: %v", err)
		}

		if err := SaveProfiles(*profileFile, profiles); err != nil {
			log.Fatalf("Failed to save profiles: %v", err)
		}

		fmt.Printf("Profiles saved to %s\n", *profileFile)
		return
	}

	// Detection mode
	profiles, err := LoadProfiles(*profileFile)
	if err != nil {
		log.Fatalf("Failed to load profiles: %v", err)
	}

	// Read text from stdin
	reader := bufio.NewReader(os.Stdin)
	var textBuilder strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Failed to read input: %v", err)
		}
		textBuilder.WriteString(line)
	}

	text := textBuilder.String()
	if text == "" {
		log.Fatal("No text provided for detection")
	}

	// Build profile for input text
	textProfile := BuildProfileFromText(text)

	// Calculate distances for all languages
	results := make([]LanguageResult, 0, len(profiles))
	for lang, profile := range profiles {
		distance := CalculateDistance(profile, textProfile)
		results = append(results, LanguageResult{
			Language: lang,
			Distance: distance,
		})
	}

	// Sort by distance (ascending)
	SortResults(results)

	if *verbose {
		// Verbose output
		fmt.Println("Language detection results (sorted by distance):")
		for _, result := range results {
			fmt.Printf("%s: %d\n", result.Language, result.Distance)
		}
	} else {
		// Normal output - best match
		if len(results) > 0 {
			fmt.Println(results[0].Language)
		} else {
			fmt.Println("No language profiles available")
		}
	}
}
