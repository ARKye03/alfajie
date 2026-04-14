package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type HighScore struct {
	Score    int
	WPM      int
	Accuracy float64
	Date     time.Time
}

func scoresPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".alfajie", "scores.json")
}

func loadHighScores() []HighScore {
	data, err := os.ReadFile(scoresPath())
	if err != nil {
		return nil
	}
	var scores []HighScore
	_ = json.Unmarshal(data, &scores)
	return scores
}

func saveHighScore(s HighScore) error {
	scores := loadHighScores()
	scores = append(scores, s)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})
	if len(scores) > 10 {
		scores = scores[:10]
	}

	path := scoresPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.Marshal(scores)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
