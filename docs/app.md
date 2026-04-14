# CLAUDE.md

Guidance for Claude Code (claude.ai/code) in this repo.

## Project

Terminal typing speed game in Go — Bubble Tea + Lip Gloss. Words fall down screen; type before bottom.

## Build & Run

```bash
go run .
go build -o alfajie .
go test ./...
go test ./... -run TestName
```

## Architecture

Bubble Tea model-based app. Single package or thin multi-package layout.

**Core types:**

```go
type Word struct {
    Text      string
    X, Y      int
    SpawnedAt time.Time
    Deadline  time.Time
}

type Game struct {
    ActiveWords []Word
    Input       string
    Score       int
    Lives       int
    StartTime   time.Time
    GameOver    bool
}
```

**Bubble Tea lifecycle:**
- `Init()` — start tick cmd
- `Update(msg)` — tick, keypress, resize
- `View()` — render via Lip Gloss

**Planned file split:**
- `model.go` — Game struct, Update, View
- `words.go` — word lists (easy/medium/hard), spawn logic
- `scoring.go` — WPM (`correctChars / 5 / minutesPlayed`), streak, accuracy
- `highscores.go` — JSON persistence to `~/.alfajie/scores.json`

## Game Rules

- 3 lives; word reaches bottom → lose 1 life
- Spawn 1 word every 2s; fall 1 row per tick (100–300ms)
- Exact match removes word instantly — no Enter
- Every 10 pts → spawn rate increases
- Difficulty: easy (<10pts), medium (10–30), hard (30+)
- Word score: `len(word)` pts + streak bonus

## Input & Matching

Shared input buffer. First exact match removed, buffer cleared. Backspace supported.

## UI

- Green highlight on typed prefix match
- Red on mismatch
- Progress bar per word for remaining time
- HUD: `Score | Lives | WPM | Accuracy | Streak`
- Game Over: final stats + restart (`r`) / quit (`q`)
- Pause: `p`

## Build Milestones (order)

1. Single word + input + Enter submit
2. Timer per word + score/lives
3. Embedded word list, random selection
4. WPM + accuracy + difficulty scaling
5. Multiple falling words with X/Y positions
6. High score JSON persistence