# alfajie

Terminal typing speed game. Words fall from the top, you type them before they hit the bottom.

## Getting Started

```bash
git clone https://github.com/ARKye03/alfajie
cd alfajie
go run .
```

Requires Go 1.21+.

## How to Play

- Type the falling word exactly — no Enter needed, matches instantly
- Miss a word reaching the bottom → lose a life (3 lives total)
- Consecutive correct words build a **streak** for bonus points
- Speed and difficulty scale as your score climbs

| Score | Words              | Tick speed |
| ----- | ------------------ | ---------- |
| 0–9   | easy (2–3 chars)   | 300ms      |
| 10–29 | medium (5–6 chars) | 200ms      |
| 30+   | hard (9–10 chars)  | 120ms      |

**Keys:** `esc` pause · `r` restart · `ctrl+c` quit

## HUD

```sh
Score:42   ♥ ♥ ♡   WPM:68   Acc:91%   Streak:5
```

High scores saved to `~/.alfajie/scores.json`.

## Tech Stack

- **TUI**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) — event-driven model/update/view loop
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss) — terminal colors and layout

## Future Ideas

- [ ] Per-word progress bar (color ramp already in place)
- [ ] Custom word lists / import from file
- [ ] Online leaderboard??

### What's with the name?

- "alfajie" is a play on "alfajores", a type of cookie from Uruguay. I guess I'm a simple man 🤷🏻

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
