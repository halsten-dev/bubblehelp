// Package bubblehelp is a manager to help render, contextualize and manage keybinds for bubbletea.
package bubblehelp

import (
	"log"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

var (
	// CurrentContext holds the current context identifier.
	CurrentContext KeymapContext

	// Contexts map of single every registered contexts.
	Contexts map[KeymapContext]*Keymap

	// ShowAll is the flag to define whether it's the full or short help that will be rendered.
	ShowAll bool

	// previousContext allows to reset the previous context.
	previousContext KeymapContext

	// defaultStyle is used by defautl by all keymaps.
	defaultStyle Style
)

// Init is the first function that needs to be called in the start of the app.
func Init() {
	Contexts = make(map[KeymapContext]*Keymap)
	defaultStyle = Style{
		EssentialKey: lipgloss.NewStyle().
			Bold(true),
		EssentialKeyDescription: lipgloss.NewStyle().
			Italic(true),
		EssentialKeySeparator: lipgloss.NewStyle().
			Italic(true),
		EssentialKeySeparatorValue: " - ",
		EssentialColSeparator: lipgloss.NewStyle().
			Bold(true),
		EssentialColSeparatorValue: " • ",
		FullKey: lipgloss.NewStyle().
			Bold(true),
		FullKeyDescription: lipgloss.NewStyle().
			Italic(true),
		FullKeySeparator: lipgloss.NewStyle().
			Italic(true),
		FullKeySeparatorValue: " - ",
		FullColSeparator: lipgloss.NewStyle().
			Italic(true),
		FullColSeparatorValue: "   ",
	}
}

// RegisterContext allows to register a new Keymap context and link it with the given identifier.
func RegisterContext(context KeymapContext, keymap *Keymap) {
	Contexts[context] = keymap
}

// GetCurrentContextKeymap returns a pointer on the current Keymap context.
func GetCurrentContextKeymap() *Keymap {
	ctx, ok := Contexts[CurrentContext]

	if !ok {
		return nil
	}

	return ctx
}

// SwitchContext take care of properly changing context. Resets the current context before switching.
func SwitchContext(context KeymapContext) {
	_, ok := Contexts[context]

	if !ok {
		log.Println("bubblehelp: context not found")
	}

	keymap := GetCurrentContextKeymap()

	if keymap != nil {
		keymap.Reset()
	}

	previousContext = CurrentContext
	CurrentContext = context
	ShowAll = false
}

func SwitchToPreviousContext() {
	SwitchContext(previousContext)
}

// UpdateKeybindHelpDesc allows to temporary change the help description for a keybind in the current Keymap context.
func UpdateKeybindHelpDesc(keybind key.Binding, desc string) {
	keymap := GetCurrentContextKeymap()

	if keymap == nil {
		return
	}

	keymap.updateHelpDesc(keybind, desc)
}

// SetKeybindVisible allows to set keybinds visibility in the current Keymap context.
func SetKeybindVisible(keybind key.Binding, visible bool) {
	keymap := GetCurrentContextKeymap()

	if keymap == nil {
		return
	}

	keymap.setVisible(keybind, visible)
}

// IsKeybindVisible returns the current visibility of a Keybind.
// If the Keybind do not exist in the current context, returns false.
func IsKeybindVisible(keybind key.Binding) bool {
	keymap := GetCurrentContextKeymap()

	if keymap == nil {
		return false
	}

	return keymap.isVisible(keybind)
}

// View is the main render function of bubblehelp.
// Take the ShowAll flag into account to render full or short help.
func View(width int) string {
	keymap := GetCurrentContextKeymap()

	if keymap == nil {
		return "ERROR : UNKNOWN KEYMAP CONTEXT"
	}

	if ShowAll {
		return ViewAll(keymap, width)
	} else {
		return ViewEssential(keymap, width)
	}
}

// ViewAll is the full help render function, can be called directly.
func ViewAll(keymap *Keymap, width int) string {
	var keys []Key
	var columns []string
	var keyStr, sepStr, descStr strings.Builder

	keys = keymap.AllBindings()

	colCount := keymap.ShowAllColumnCount
	rowCount := int(math.Ceil(float64(len(keys)) / float64(colCount)))

	columns = make([]string, 0)

	for i, key := range keys {
		remainingCount := len(keys) - i
		notLastCol := len(columns)+1 < colCount

		if i%rowCount > 0 || (remainingCount == 1 && notLastCol) {
			keyStr.WriteString("\n")
			sepStr.WriteString("\n")
			descStr.WriteString("\n")
		}

		keyStr.WriteString(keymap.Style.FullKey.
			Render(key.Binding.Help().Key))
		sepStr.WriteString(keymap.Style.FullKeySeparator.
			Render(keymap.Style.FullKeySeparatorValue))
		descStr.WriteString(keymap.Style.FullKeyDescription.
			Render(key.GetHelpDesc()))

		if ((i+1)%rowCount == 0 && i != 0 && notLastCol) || (remainingCount == 1) {
			columns = append(columns, lipgloss.
				JoinHorizontal(lipgloss.Center,
					keyStr.String(),
					sepStr.String(),
					descStr.String()))

			if i < len(keys)-1 {
				columns = append(columns, keymap.Style.FullColSeparatorValue)
				colCount++
			}

			keyStr.Reset()
			sepStr.Reset()
			descStr.Reset()
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}

// ViewEssential is the short help render function, can be called directly
func ViewEssential(keymap *Keymap, width int) string {
	var b strings.Builder
	var keys []Key

	keys = keymap.EssentialBindings()

	for i, key := range keys {
		if i > 0 {
			b.WriteString(keymap.Style.EssentialColSeparator.
				Render(keymap.Style.EssentialColSeparatorValue))
		}

		b.WriteString(keymap.Style.EssentialKey.
			Render(key.Binding.Help().Key))

		b.WriteString(keymap.Style.EssentialKeySeparator.
			Render(keymap.Style.EssentialKeySeparatorValue))

		b.WriteString(keymap.Style.EssentialKeyDescription.
			Render(key.GetHelpDesc()))
	}

	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Width(width).Render(b.String())
}

func SetDefaultStyle(style Style) {
	defaultStyle = style
}
