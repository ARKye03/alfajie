package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time
type spawnMsg time.Time

type Word struct {
	Text string
	X, Y int
}

type Game struct {
	ActiveWords  []Word
	Input        string
	Score        int
	Lives        int
	StartTime    time.Time
	GameOver     bool
	Paused       bool
	Width        int
	Height       int
	CorrectChars int
	CorrectWords int
	MissedWords  int
	Streak       int
	BestStreak   int
}

var (
	normalStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	greenStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	dangerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	hudStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	inputStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	inputErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))
	pauseStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))
	heartStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

func NewGame() Game {
	return Game{
		Lives:     3,
		StartTime: time.Now(),
	}
}

func (g Game) Init() tea.Cmd {
	return tea.Batch(tickCmd(g.Score), spawnCmd(g.Score))
}

func tickInterval(score int) time.Duration {
	switch {
	case score >= 30:
		return 120 * time.Millisecond
	case score >= 10:
		return 200 * time.Millisecond
	default:
		return 300 * time.Millisecond
	}
}

func spawnInterval(score int) time.Duration {
	switch {
	case score >= 30:
		return 900 * time.Millisecond
	case score >= 10:
		return 1400 * time.Millisecond
	default:
		return 2000 * time.Millisecond
	}
}

func tickCmd(score int) tea.Cmd {
	return tea.Tick(tickInterval(score), func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func spawnCmd(score int) tea.Cmd {
	return tea.Tick(spawnInterval(score), func(t time.Time) tea.Msg {
		return spawnMsg(t)
	})
}

func (g Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		g.Width = msg.Width
		g.Height = msg.Height
		return g, nil

	case tea.KeyMsg:
		if g.GameOver {
			switch msg.String() {
			case "r":
				ng := NewGame()
				ng.Width = g.Width
				ng.Height = g.Height
				return ng, tea.Batch(tickCmd(0), spawnCmd(0))
			case "q", "ctrl+c":
				return g, tea.Quit
			}
			return g, nil
		}
		switch msg.String() {
		case "ctrl+c":
			return g, tea.Quit
		case "p":
			g.Paused = !g.Paused
			if !g.Paused {
				return g, tea.Batch(tickCmd(g.Score), spawnCmd(g.Score))
			}
			return g, nil
		case "backspace":
			runes := []rune(g.Input)
			if len(runes) > 0 {
				g.Input = string(runes[:len(runes)-1])
			}
		default:
			if g.Paused {
				return g, nil
			}
			r := msg.String()
			if len([]rune(r)) == 1 {
				g.Input += r
				for i, w := range g.ActiveWords {
					if g.Input == w.Text {
						g.Score += len(w.Text) + g.Streak
						g.CorrectChars += len(w.Text)
						g.CorrectWords++
						g.Streak++
						if g.Streak > g.BestStreak {
							g.BestStreak = g.Streak
						}
						g.ActiveWords = append(g.ActiveWords[:i], g.ActiveWords[i+1:]...)
						g.Input = ""
						break
					}
				}
			}
		}

	case tickMsg:
		if g.Paused || g.GameOver {
			return g, nil
		}
		gameH := g.gameHeight()
		var newWords []Word
		for _, w := range g.ActiveWords {
			w.Y++
			if w.Y >= gameH {
				g.Lives--
				g.MissedWords++
				g.Streak = 0
				if g.Lives <= 0 {
					g.GameOver = true
					_ = saveHighScore(HighScore{
						Score:    g.Score,
						WPM:      calcWPM(g.CorrectChars, time.Since(g.StartTime)),
						Accuracy: calcAccuracy(g.CorrectWords, g.MissedWords),
						Date:     time.Now(),
					})
					return g, nil
				}
			} else {
				newWords = append(newWords, w)
			}
		}
		g.ActiveWords = newWords
		return g, tickCmd(g.Score)

	case spawnMsg:
		if g.Paused || g.GameOver {
			return g, nil
		}
		g.ActiveWords = append(g.ActiveWords, spawnWord(g.Score, g.Width))
		return g, spawnCmd(g.Score)
	}

	return g, nil
}

func (g Game) gameHeight() int {
	if g.Height < 5 {
		return 1
	}
	return g.Height - 3
}

func (g Game) View() string {
	if g.Width == 0 || g.Height == 0 {
		return "Loading..."
	}
	if g.GameOver {
		return g.gameOverView()
	}

	gameH := g.gameHeight()

	// Determine which words match the current input prefix
	matched := make(map[int]bool)
	hasMatch := false
	if g.Input != "" {
		for i, w := range g.ActiveWords {
			if strings.HasPrefix(w.Text, g.Input) {
				matched[i] = true
				hasMatch = true
			}
		}
	}

	// Build per-cell grid
	type cell struct {
		ch    rune
		color int // 0=normal 1=green 2=warn 3=danger
	}
	grid := make([][]cell, gameH)
	for i := range grid {
		grid[i] = make([]cell, g.Width)
		for j := range grid[i] {
			grid[i][j] = cell{ch: ' '}
		}
	}

	for i, w := range g.ActiveWords {
		if w.Y < 0 || w.Y >= gameH {
			continue
		}
		color := 0
		if matched[i] {
			color = 1
		} else {
			ratio := float64(w.Y+1) / float64(gameH)
			switch {
			case ratio > 0.8:
				color = 3
			case ratio > 0.5:
				color = 2
			}
		}
		for j, ch := range []rune(w.Text) {
			x := w.X + j
			if x >= 0 && x < g.Width {
				grid[w.Y][x] = cell{ch: ch, color: color}
			}
		}
	}

	// Render grid
	var sb strings.Builder
	styleFor := func(color int) lipgloss.Style {
		switch color {
		case 1:
			return greenStyle
		case 2:
			return warnStyle
		case 3:
			return dangerStyle
		default:
			return normalStyle
		}
	}

	for rowIdx, row := range grid {
		if rowIdx > 0 {
			sb.WriteByte('\n')
		}
		i := 0
		for i < len(row) {
			c := row[i]
			if c.ch == ' ' {
				sb.WriteByte(' ')
				i++
				continue
			}
			j := i + 1
			for j < len(row) && row[j].color == c.color && row[j].ch != ' ' {
				j++
			}
			var seg []rune
			for _, cc := range row[i:j] {
				seg = append(seg, cc.ch)
			}
			sb.WriteString(styleFor(c.color).Render(string(seg)))
			i = j
		}
	}

	// HUD
	elapsed := time.Since(g.StartTime)
	wpm := calcWPM(g.CorrectChars, elapsed)
	acc := calcAccuracy(g.CorrectWords, g.MissedWords)

	hearts := heartStyle.Render(strings.Repeat("♥ ", g.Lives)) +
		dimStyle.Render(strings.Repeat("♡ ", 3-g.Lives))
	hud := hudStyle.Render(fmt.Sprintf("Score:%-5d  %s  WPM:%-4d  Acc:%.0f%%  Streak:%d",
		g.Score, hearts, wpm, acc, g.Streak))

	sep := dimStyle.Render(strings.Repeat("─", g.Width))

	var inputLine string
	if g.Paused {
		inputLine = pauseStyle.Render("  PAUSED  p=resume  ctrl+c=quit")
	} else {
		display := g.Input
		if display == "" {
			display = "_"
		}
		iStyle := inputStyle
		if g.Input != "" && !hasMatch {
			iStyle = inputErrStyle
		}
		inputLine = "> " + iStyle.Render(display)
	}

	sb.WriteString("\n" + sep)
	sb.WriteString("\n" + hud)
	sb.WriteString("\n" + inputLine)

	return sb.String()
}

func (g Game) gameOverView() string {
	elapsed := time.Since(g.StartTime)
	wpm := calcWPM(g.CorrectChars, elapsed)
	acc := calcAccuracy(g.CorrectWords, g.MissedWords)

	scores := loadHighScores()
	best := 0
	for _, s := range scores {
		if s.Score > best {
			best = s.Score
		}
	}

	lines := []string{
		titleStyle.Render("╔═══ GAME OVER ═══╗"),
		"",
		fmt.Sprintf("  Score:    %d", g.Score),
		fmt.Sprintf("  Best:     %d", best),
		fmt.Sprintf("  WPM:      %d", wpm),
		fmt.Sprintf("  Accuracy: %.0f%%", acc),
		fmt.Sprintf("  Streak:   %d  (best: %d)", g.Streak, g.BestStreak),
		"",
		dimStyle.Render("  r — restart   q — quit"),
	}

	content := strings.Join(lines, "\n")
	return lipgloss.Place(g.Width, g.Height, lipgloss.Center, lipgloss.Center, content)
}
