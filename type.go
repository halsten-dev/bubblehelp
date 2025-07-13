package bubblehelp

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// KeymapContext is a special type used inside bubblehelp to serve as a unique identifier for different contexts.
type KeymapContext string

// Style describes every rendering aspect of bubblehelp.
type Style struct {
	EssentialKey               lipgloss.Style
	EssentialKeyDescription    lipgloss.Style
	EssentialKeySeparator      lipgloss.Style
	EssentialKeySeparatorValue string
	EssentialColSeparator      lipgloss.Style
	EssentialColSeparatorValue string
	FullKey                    lipgloss.Style
	FullKeyDescription         lipgloss.Style
	FullKeySeparator           lipgloss.Style
	FullKeySeparatorValue      string
	FullColSeparator           lipgloss.Style
	FullColSeparatorValue      string
}

// Key describes a single keybind in bubblehelp.
type Key struct {
	Binding        key.Binding
	Essential      bool
	CustomHelpDesc string
	Visible        bool
}

// GetHelpDesc returns the current help description for the Key.
// Returns the temporary CustomHelpDesc first.
func (k *Key) GetHelpDesc() string {
	if k.CustomHelpDesc == "" {
		return k.Binding.Help().Desc
	}

	return k.CustomHelpDesc
}

// Keymap describes a keymap context inside bubblehelp. Each context can have its own Style definition
type Keymap struct {
	Bindings           []Key
	ShowAllColumnCount int
	Style              Style
}

// NewKeymap helps initialize a new context with default style and given column count for full help rendering.
// showAllColCount parameters allows definition of the columns of the full help rendering.
func NewKeymap(showAllColCount int) *Keymap {
	return &Keymap{
		Bindings:           make([]Key, 0),
		ShowAllColumnCount: showAllColCount,
		Style: Style{
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
		},
	}
}

// NewKeyBinding registers a new key.Binding to the Keymap context.
// essential parameter allows to define if the key should appear on the short help or only in full help.
func (k *Keymap) NewKeyBinding(binding key.Binding, essential bool) {
	k.Bindings = append(k.Bindings, Key{
		Binding:        binding,
		Essential:      essential,
		CustomHelpDesc: "",
		Visible:        true,
	})
}

// EssentialBindings returns the slice of all essential binding of the Keymap context.
func (k *Keymap) EssentialBindings() []Key {
	var essentials []Key

	for _, k := range k.Bindings {
		if !k.Essential || !k.Visible {
			continue
		}

		essentials = append(essentials, k)
	}

	return essentials
}

// AllBindings returns every single binding of the Keymap context.
func (k *Keymap) AllBindings() []Key {
	var all []Key

	for _, k := range k.Bindings {
		if !k.Visible {
			continue
		}

		all = append(all, k)
	}

	return all
}

// Reset is called when the manager switches contexts.
// Every Key are set back to visible and temporary CustomHelpDesc is cleared.
func (k *Keymap) Reset() {
	for i := 0; i < len(k.Bindings); i++ {
		k.Bindings[i].CustomHelpDesc = ""
		k.Bindings[i].Visible = true
	}
}

// SetHelpDesc allows for permanently change the default help text of the keybind.
// Only in the Keymap context
func (k *Keymap) SetHelpDesc(keybind key.Binding, desc string) {
	for i := 0; i < len(k.Bindings); i++ {
		if keybind.Help().Key == k.Bindings[i].Binding.Help().Key {
			k.Bindings[i].Binding.SetHelp(k.Bindings[i].Binding.Help().Key, desc)
			return
		}
	}
}

// updateHelpDesc helper method to set a CustomHelpDesc on the given Key.
// Identified by the key.Binding directly.
// Call this function with an empty string to return to the default help text.
func (k *Keymap) updateHelpDesc(keybind key.Binding, desc string) {
	for i := 0; i < len(k.Bindings); i++ {
		if keybind.Help().Key == k.Bindings[i].Binding.Help().Key {
			k.Bindings[i].CustomHelpDesc = desc
			return
		}
	}
}

// setVisible allows to set the keybind visibility. Resets when contexts are switched.
func (k *Keymap) setVisible(keybind key.Binding, visible bool) {
	for i := 0; i < len(k.Bindings); i++ {
		if keybind.Help().Key == k.Bindings[i].Binding.Help().Key {
			k.Bindings[i].Visible = visible
			return
		}
	}
}

// isVisible returns the current visibility of the keybind.
func (k *Keymap) isVisible(keybind key.Binding) bool {
	for i := 0; i < len(k.Bindings); i++ {
		if keybind.Help().Key == k.Bindings[i].Binding.Help().Key {
			return k.Bindings[i].Visible
		}
	}

	return false
}
