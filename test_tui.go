package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dylan/spcstr/internal/tui"
)

func main() {
	m := tui.NewModel()

	// Test initial state
	fmt.Println("Initial view:", m.View())

	// Test window size message
	sizeMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updated, _ := m.Update(sizeMsg)
	m = updated.(tui.Model)
	fmt.Println("\nAfter window size update:")
	fmt.Println(m.View())

	// Test navigation to observe
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}}
	updated, _ = m.Update(keyMsg)
	m = updated.(tui.Model)
	fmt.Println("\nAfter pressing 'o':")
	fmt.Println(m.View())

	// Test navigation back to plan
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}}
	updated, _ = m.Update(keyMsg)
	m = updated.(tui.Model)
	fmt.Println("\nAfter pressing 'p':")
	fmt.Println(m.View())

	fmt.Println("\nTUI structure test completed successfully!")
}
