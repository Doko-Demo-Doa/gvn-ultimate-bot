package bot

import (
	"doko/gvn-ultimate-bot/models"
	"doko/gvn-ultimate-bot/scheduler"
	"doko/gvn-ultimate-bot/services/discordservice"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func RegisterRoleReactModule(
	s *state.State,
	ds discordservice.DiscordService,
	dre discordservice.DiscordRoleReactionEmbedService,
	rs *scheduler.RoleScheduler,
	guildID discord.GuildID,
) {
	// Listen for raw emoji reactions.
	s.AddHandler(func(e *gateway.MessageReactionAddEvent) {
		handleReaction(s, dre, rs, e.MessageID, e.ChannelID, string(e.Emoji.APIString()), e.UserID.String(), true)
	})

	s.AddHandler(func(e *gateway.MessageReactionRemoveEvent) {
		handleReaction(s, dre, rs, e.MessageID, e.ChannelID, string(e.Emoji.APIString()), e.UserID.String(), false)
	})

	// Listen for button clicks and dropdown selections.
	s.AddHandler(func(e *gateway.InteractionCreateEvent) {
		handleInteraction(s, dre, rs, e)
	})
}

func handleReaction(
	s *state.State,
	dre discordservice.DiscordRoleReactionEmbedService,
	rs *scheduler.RoleScheduler,
	messageID discord.MessageID,
	channelID discord.ChannelID,
	emoji string,
	userID string,
	added bool,
) {
	// Ignore bot's own reactions.
	me, err := s.Me()
	if err == nil && me.ID.String() == userID {
		return
	}

	embed, err := dre.GetEmbedByNativeMessageID(messageID.String())
	if err != nil || embed == nil {
		return
	}

	payload, err := embed.ParsedPayload()
	if err != nil {
		log.Printf("[role_react] failed to parse payload for message %s: %v", messageID, err)
		return
	}

	roleID := payload.FindRoleByEmoji(emoji)
	if roleID == "" {
		return
	}

	applyReactionRole(rs, payload.Mode, userID, roleID, added)
}

func handleInteraction(
	s *state.State,
	dre discordservice.DiscordRoleReactionEmbedService,
	rs *scheduler.RoleScheduler,
	e *gateway.InteractionCreateEvent,
) {
	var messageID discord.MessageID
	var userID string
	var customID string
	var selectedValues []string

	switch data := e.InteractionEvent.Data.(type) {
	case *discord.ButtonInteraction:
		if e.Message != nil {
			messageID = e.Message.ID
		}
		if sender := e.Sender(); sender != nil {
			userID = sender.ID.String()
		}
		customID = string(data.CustomID)
	case *discord.StringSelectInteraction:
		if e.Message != nil {
			messageID = e.Message.ID
		}
		if sender := e.Sender(); sender != nil {
			userID = sender.ID.String()
		}
		customID = string(data.CustomID)
		selectedValues = data.Values
	default:
		return
	}

	if messageID == 0 || userID == "" {
		return
	}

	embed, err := dre.GetEmbedByNativeMessageID(messageID.String())
	if err != nil || embed == nil {
		return
	}

	payload, err := embed.ParsedPayload()
	if err != nil {
		log.Printf("[role_react] failed to parse payload for message %s: %v", messageID, err)
		return
	}

	interaction := payload.FindInteractionByID(customID)
	if interaction == nil {
		return
	}

	var roleID string
	if interaction.Type == models.InteractionTypeButton {
		roleID = interaction.RoleNativeID
	} else if interaction.Type == models.InteractionTypeDropdown && len(selectedValues) > 0 {
		opt := interaction.FindDropdownOption(selectedValues[0])
		if opt != nil {
			roleID = opt.RoleNativeID
		}
	}

	if roleID == "" {
		respondInteraction(s, e, "Không tìm thấy role cho tương tác này.")
		return
	}

	// For one-shot interactions, "added" is always true.
	applyReactionRole(rs, payload.Mode, userID, roleID, true)
	respondInteraction(s, e, "Đã cập nhật role của bạn.")
}

func applyReactionRole(rs *scheduler.RoleScheduler, mode models.ReactionMode, userID, roleID string, added bool) {
	grant := added
	if mode == models.ReactionModeReverse {
		grant = !added
	}

	if grant {
		if err := rs.AddRole(userID, roleID); err != nil {
			log.Printf("[role_react] failed to add role %s to user %s: %v", roleID, userID, err)
		}
	} else {
		if err := rs.RemoveRole(userID, roleID); err != nil {
			log.Printf("[role_react] failed to remove role %s from user %s: %v", roleID, userID, err)
		}
	}
}

func respondInteraction(s *state.State, e *gateway.InteractionCreateEvent, msg string) {
	data := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: option.NewNullableString(msg),
			Flags:   discord.EphemeralMessage,
		},
	}
	if err := s.RespondInteraction(e.ID, e.Token, data); err != nil {
		log.Println("[role_react] failed to respond to interaction:", err)
	}
}
