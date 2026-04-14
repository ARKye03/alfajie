package main

import "time"

func calcWPM(correctChars int, elapsed time.Duration) int {
	if elapsed < time.Second || correctChars == 0 {
		return 0
	}
	return int(float64(correctChars) / 5.0 / elapsed.Minutes())
}

func calcAccuracy(correctWords, missedWords int) float64 {
	total := correctWords + missedWords
	if total == 0 {
		return 100.0
	}
	return float64(correctWords) / float64(total) * 100.0
}
