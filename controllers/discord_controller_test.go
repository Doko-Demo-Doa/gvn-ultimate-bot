package controllers

import (
	"testing"
	"time"
)

// Regression test: inputToDiscordRole used to swap Name and NativeID when
// mapping the create-role request into the model, silently corrupting every
// role created through the admin UI.
func TestInputToDiscordRole_MapsNameAndNativeIDCorrectly(t *testing.T) {
	ctl := &discordController{}

	input := DiscordRoleInput{
		Name:         "Moderator",
		NativeID:     "123456789012345678",
		Mentionable:  1,
		Hoist:        1,
		Color:        0xff0000,
		Expiry:       time.Now(),
		ImplicitType: 2,
	}

	role := ctl.inputToDiscordRole(input)

	if role.Name != input.Name {
		t.Errorf("expected Name %q, got %q", input.Name, role.Name)
	}
	if role.NativeID != input.NativeID {
		t.Errorf("expected NativeID %q, got %q", input.NativeID, role.NativeID)
	}
	if role.Mentionable != input.Mentionable {
		t.Errorf("expected Mentionable %d, got %d", input.Mentionable, role.Mentionable)
	}
	if role.Hoist != input.Hoist {
		t.Errorf("expected Hoist %d, got %d", input.Hoist, role.Hoist)
	}
	if role.Color != input.Color {
		t.Errorf("expected Color %d, got %d", input.Color, role.Color)
	}
	if role.ImplicitType != input.ImplicitType {
		t.Errorf("expected ImplicitType %d, got %d", input.ImplicitType, role.ImplicitType)
	}
}
