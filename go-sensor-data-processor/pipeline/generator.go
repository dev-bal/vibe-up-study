package pipeline

import (
	"log"
	"os"
	"path/filepath"

	"golang.org/x/exp/rand"
)

func GetFiles(inputDir string) <-chan string {
	fileChan := make(chan string, 10)

	go func() {
		entries, err := os.ReadDir(inputDir)
		if err != nil {
			log.Printf("Error reading directory: %v", err)
			close(fileChan)
			return
		}

		var files []string

		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".gz" {
				files = append(files, filepath.Join(inputDir, entry.Name()))
			}
		}

		rand.Seed(22) // Use current time for randomness
		rand.Shuffle(len(files), func(i, j int) {
			files[i], files[j] = files[j], files[i]
		})

		for _, file := range files {
			fileChan <- file
		}

		close(fileChan)
	}()

	return fileChan

}
