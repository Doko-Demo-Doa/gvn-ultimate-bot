package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ReactionMode controls how emoji reactions toggle roles.
// Buttons and dropdowns are one-shot interactions, so Reverse simply inverts
// their effect (remove role instead of granting it).
type ReactionMode string

const (
	// ReactionModeDefault grants the mapped role when a user reacts / clicks.
	ReactionModeDefault ReactionMode = "default"
	// ReactionModeReverse removes the mapped role when a user reacts / clicks,
	// and grants it back when the reaction is removed.
	ReactionModeReverse ReactionMode = "reverse"
)

// InteractionType describes the kind of UI element that grants a role.
type InteractionType string

const (
	InteractionTypeEmoji    InteractionType = "emoji"
	InteractionTypeButton   InteractionType = "button"
	InteractionTypeDropdown InteractionType = "dropdown"
)

// ReactionRoleMessagePayload is the structured JSON document stored in
// DiscordRoleReactionEmbed.Payload. It is stored as a plain string so the
// model stays database-agnostic (works on Postgres, MySQL, SQLite, etc.).
// If you ever need Postgres-specific jsonb queries/indexes, migrate the
// Payload column to jsonb and switch the model tag to "type:jsonb".
type ReactionRoleMessagePayload struct {
	ChannelID   string                   `json:"channel_id"`
	Message     string                   `json:"message,omitempty"`
	Mode        ReactionMode             `json:"mode"`
	Embed       *ReactionRoleEmbed       `json:"embed,omitempty"`
	Interactions []ReactionInteraction   `json:"interactions"`
}

// ReactionRoleEmbed mirrors Discord's rich embed fields.
type ReactionRoleEmbed struct {
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	Color       int               `json:"color,omitempty"`
	ImageURL    string            `json:"image_url,omitempty"`
	ThumbnailURL string           `json:"thumbnail_url,omitempty"`
	Footer      string            `json:"footer,omitempty"`
	Author      string            `json:"author,omitempty"`
	Fields      []ReactionEmbedField `json:"fields,omitempty"`
}

// ReactionEmbedField is a single field inside a rich embed.
type ReactionEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// ReactionInteraction maps one UI element to a Discord role.
// For emoji/button interactions the role is stored directly on this struct.
// For dropdown interactions the role is stored inside each DropdownOption.
type ReactionInteraction struct {
	ID        string             `json:"id"`
	Type      InteractionType    `json:"type"`
	Emoji     string             `json:"emoji,omitempty"`
	Label     string             `json:"label,omitempty"`
	Style     ButtonStyle        `json:"style,omitempty"`
	RoleNativeID string          `json:"role_native_id,omitempty"`
	Placeholder  string          `json:"placeholder,omitempty"`
	Options      []DropdownOption `json:"options,omitempty"`
}

// ButtonStyle matches Discord's button styles.
type ButtonStyle string

const (
	ButtonStylePrimary   ButtonStyle = "primary"
	ButtonStyleSecondary ButtonStyle = "secondary"
	ButtonStyleSuccess   ButtonStyle = "success"
	ButtonStyleDanger    ButtonStyle = "danger"
)

// DropdownOption is a single option inside a dropdown/select menu.
type DropdownOption struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	Emoji        string `json:"emoji,omitempty"`
	Description  string `json:"description,omitempty"`
	RoleNativeID string `json:"role_native_id"`
}

// Validate checks that the payload is well-formed and every interaction maps
// to a non-empty role native id.
func (p *ReactionRoleMessagePayload) Validate() error {
	if p.ChannelID == "" {
		return errors.New("channel_id is required")
	}

	if p.Mode != "" && p.Mode != ReactionModeDefault && p.Mode != ReactionModeReverse {
		return fmt.Errorf("invalid reaction mode: %s", p.Mode)
	}

	seenIDs := make(map[string]struct{})
	for i, it := range p.Interactions {
		if it.ID == "" {
			return fmt.Errorf("interaction %d is missing an id", i)
		}
		if _, exists := seenIDs[it.ID]; exists {
			return fmt.Errorf("duplicate interaction id: %s", it.ID)
		}
		seenIDs[it.ID] = struct{}{}

		switch it.Type {
		case InteractionTypeEmoji, InteractionTypeButton:
			if strings.TrimSpace(it.RoleNativeID) == "" {
				return fmt.Errorf("interaction %s must map to a role", it.ID)
			}
		case InteractionTypeDropdown:
			if len(it.Options) == 0 {
				return fmt.Errorf("dropdown %s must have at least one option", it.ID)
			}
			for j, opt := range it.Options {
				if strings.TrimSpace(opt.RoleNativeID) == "" {
					return fmt.Errorf("dropdown %s option %d must map to a role", it.ID, j)
				}
			}
		default:
			return fmt.Errorf("interaction %s has invalid type: %s", it.ID, it.Type)
		}
	}

	return nil
}

// FindRoleByEmoji returns the role native id mapped to the given emoji,
// or an empty string if none matches.
func (p *ReactionRoleMessagePayload) FindRoleByEmoji(emoji string) string {
	for _, it := range p.Interactions {
		if it.Type == InteractionTypeEmoji && it.Emoji == emoji {
			return it.RoleNativeID
		}
	}
	return ""
}

// FindInteractionByID returns the interaction with the matching custom id.
func (p *ReactionRoleMessagePayload) FindInteractionByID(id string) *ReactionInteraction {
	for i := range p.Interactions {
		if p.Interactions[i].ID == id {
			return &p.Interactions[i]
		}
	}
	return nil
}

// FindDropdownOption returns the option with the matching id.
func (it *ReactionInteraction) FindDropdownOption(id string) *DropdownOption {
	for i := range it.Options {
		if it.Options[i].ID == id {
			return &it.Options[i]
		}
	}
	return nil
}

// ToJSON serialises the payload to a compact JSON string.
func (p *ReactionRoleMessagePayload) ToJSON() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ParseReactionRolePayload parses a JSON string into a structured payload.
func ParseReactionRolePayload(s string) (*ReactionRoleMessagePayload, error) {
	var p ReactionRoleMessagePayload
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return nil, err
	}
	if p.Mode == "" {
		p.Mode = ReactionModeDefault
	}
	return &p, nil
}
