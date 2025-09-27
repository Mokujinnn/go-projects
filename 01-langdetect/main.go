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
	profileFile := flag.String("profile", "", "")
	trainDir := flag.String("train", "", "")
	verbose := flag.Bool("verbose", false, "")
	flag.Parse()

	if *profileFile == "" {
		log.Fatal("Profile file path is required")
	}

	if *trainDir != "" {
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

	profiles, err := LoadProfiles(*profileFile)
	if err != nil {
		log.Fatalf("Failed to load profiles: %v", err)
	}

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

	textProfile := BuildProfileFromText(text)

	results := make([]LanguageResult, 0, len(profiles))
	for lang, profile := range profiles {
		distance := CalculateDistance(profile, textProfile)
		results = append(results, LanguageResult{
			Language: lang,
			Distance: distance,
		})
	}

	SortResults(results)

	if *verbose {
		fmt.Println("Language detection results (sorted by distance):")
		for _, result := range results {
			fmt.Printf("%s: %d\n", result.Language, result.Distance)
		}
	} else {
		if len(results) > 0 {
			fmt.Println(results[0].Language)
		} else {
			fmt.Println("No language profiles available")
		}
	}
}
