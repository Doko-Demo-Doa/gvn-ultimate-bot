package discordservice

import (
	"doko/gvn-ultimate-bot/models"
	discordrepos "doko/gvn-ultimate-bot/repositories/discord_repos"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

type DiscordService interface {
	// ################# For DiscordRole #################
	ListRoles() ([]*models.DiscordRole, error)
	CreateRole(*models.DiscordRole) (*models.DiscordRole, error)
	EditRole(*models.DiscordRole) (*models.DiscordRole, error)
	UnassignRole(nativeUserId string, roleId uint) (*models.DiscordRole, error)

	// ################# For DiscordUserRole (timed assignments) #################
	AssignRoleToUser(userNativeID string, roleNativeID string, duration time.Duration) (*models.DiscordUserRole, error)
	GetExpiredRoleAssignments() ([]*models.DiscordUserRole, error)
	GetAllActiveAssignments() ([]*models.DiscordUserRole, error)
	RevokeRoleAssignment(assignmentID uint) error
	GetAssignmentByID(id uint) (*models.DiscordUserRole, error)
	GetActiveAssignmentsForUser(nativeUserID string) ([]*models.DiscordUserRole, error)
}

type discordService struct {
	RoleRepo              discordrepos.DiscordRoleRepo
	RoleReactionEmbedRepo discordrepos.DiscordRoleReactionEmbedRepo
	UserRoleRepo          discordrepos.DiscordUserRoleRepo
}

func NewDiscordRoleService(
	repo discordrepos.DiscordRoleRepo,
	embedRepo discordrepos.DiscordRoleReactionEmbedRepo,
	userRoleRepo discordrepos.DiscordUserRoleRepo,
) DiscordService {
	return &discordService{
		RoleRepo:              repo,
		RoleReactionEmbedRepo: embedRepo,
		UserRoleRepo:          userRoleRepo,
	}
}

func (dr *discordService) CreateRole(r *models.DiscordRole) (*models.DiscordRole, error) {
	return dr.RoleRepo.CreateRole(r)
}

func (dr *discordService) EditRole(r *models.DiscordRole) (*models.DiscordRole, error) {
	return dr.RoleRepo.EditRole(r)
}

func (dr *discordService) ListRoles() ([]*models.DiscordRole, error) {
	return dr.RoleRepo.ListRoles()
}

func (dr *discordService) UnassignRole(nativeUserId string, roleId uint) (*models.DiscordRole, error) {
	panic("unimplemented")
}

func (dr *discordService) AssignRoleToUser(userNativeID string, roleNativeID string, duration time.Duration) (*models.DiscordUserRole, error) {
	now := time.Now()
	expiration := now.Add(duration)

	assignment := &models.DiscordUserRole{
		UserNativeID:   userNativeID,
		RoleNativeID:   roleNativeID,
		GrantedDate:    now,
		ExpirationDate: expiration,
	}

	return dr.UserRoleRepo.CreateAssignment(assignment)
}

func (dr *discordService) GetExpiredRoleAssignments() ([]*models.DiscordUserRole, error) {
	return dr.UserRoleRepo.GetExpiredAssignments()
}

func (dr *discordService) GetAllActiveAssignments() ([]*models.DiscordUserRole, error) {
	return dr.UserRoleRepo.GetAllActiveAssignments()
}

func (dr *discordService) RevokeRoleAssignment(assignmentID uint) error {
	return dr.UserRoleRepo.RevokeAssignment(assignmentID)
}

func (dr *discordService) GetAssignmentByID(id uint) (*models.DiscordUserRole, error) {
	return dr.UserRoleRepo.GetByID(id)
}

func (dr *discordService) GetActiveAssignmentsForUser(nativeUserID string) ([]*models.DiscordUserRole, error) {
	return dr.UserRoleRepo.GetActiveAssignmentsByUser(nativeUserID)
}

// ################# For DiscordRoleReactionEmbed #################

type DiscordRoleReactionEmbedService interface {
	ListEmbeds() ([]*models.DiscordRoleReactionEmbed, error)
	UpsertEmbed(*models.DiscordRoleReactionEmbed, *models.ReactionRoleMessagePayload) (*models.DiscordRoleReactionEmbed, error)
	GetSingleEmbed(id uint) (*models.DiscordRoleReactionEmbed, error)
	GetEmbedByNativeMessageID(nativeMessageID string) (*models.DiscordRoleReactionEmbed, error)
	DeleteEmbed(id uint) error
	EditEmbed(nativeMessageID string, payload *models.ReactionRoleMessagePayload) (*models.DiscordRoleReactionEmbed, error)
	PublishEmbed(*models.ReactionRoleMessagePayload) (*models.DiscordRoleReactionEmbed, error)
}

type discordRoleReactionEmbedService struct {
	RoleReactionRepo discordrepos.DiscordRoleReactionEmbedRepo
	state            *state.State
	guildID          discord.GuildID
}

func NewDiscordRoleReactionEmbedService(repo discordrepos.DiscordRoleReactionEmbedRepo, s *state.State, guildID discord.GuildID) DiscordRoleReactionEmbedService {
	return &discordRoleReactionEmbedService{
		RoleReactionRepo: repo,
		state:            s,
		guildID:          guildID,
	}
}

func (d *discordRoleReactionEmbedService) ListEmbeds() ([]*models.DiscordRoleReactionEmbed, error) {
	return d.RoleReactionRepo.ListRoleReactionEmbeds()
}

func (d *discordRoleReactionEmbedService) UpsertEmbed(m *models.DiscordRoleReactionEmbed, payload *models.ReactionRoleMessagePayload) (*models.DiscordRoleReactionEmbed, error) {
	if payload != nil {
		if err := payload.Validate(); err != nil {
			return nil, fmt.Errorf("invalid payload: %w", err)
		}
		jsonStr, err := payload.ToJSON()
		if err != nil {
			return nil, err
		}
		m.Payload = jsonStr
		m.Mode = string(payload.Mode)
	}

	data, err := d.RoleReactionRepo.GetByNativeID(m.NativeMessageId)
	if err != nil || data == nil {
		return d.RoleReactionRepo.Create(m)
	}
	return d.RoleReactionRepo.Update(m.NativeMessageId, m)
}

func (d *discordRoleReactionEmbedService) GetSingleEmbed(id uint) (*models.DiscordRoleReactionEmbed, error) {
	return d.RoleReactionRepo.GetByID(id)
}

func (d *discordRoleReactionEmbedService) GetEmbedByNativeMessageID(nativeMessageID string) (*models.DiscordRoleReactionEmbed, error) {
	return d.RoleReactionRepo.GetByNativeID(nativeMessageID)
}

func (d *discordRoleReactionEmbedService) DeleteEmbed(id uint) error {
	embed, err := d.RoleReactionRepo.GetByID(id)
	if err != nil || embed == nil {
		return fmt.Errorf("embed not found: %w", err)
	}

	payload, err := embed.ParsedPayload()
	if err == nil && payload != nil && payload.ChannelID != "" && embed.NativeMessageId != "" {
		channelID, err := discord.ParseSnowflake(payload.ChannelID)
		if err == nil {
			msgID, parseErr := discord.ParseSnowflake(embed.NativeMessageId)
			if parseErr == nil {
				if delErr := d.state.DeleteMessage(discord.ChannelID(channelID), discord.MessageID(msgID), api.AuditLogReason("")); delErr != nil {
					log.Printf("[delete_embed] failed to delete Discord message %s: %v", embed.NativeMessageId, delErr)
				}
			}
		}
	}

	return d.RoleReactionRepo.Delete(id)
}

// EditEmbed updates the Discord message content / components / embed and
// persists the new configuration in the database.
func (d *discordRoleReactionEmbedService) EditEmbed(nativeMessageID string, payload *models.ReactionRoleMessagePayload) (*models.DiscordRoleReactionEmbed, error) {
	if payload == nil {
		return nil, errors.New("payload is required")
	}
	if err := payload.Validate(); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	channelID, err := discord.ParseSnowflake(payload.ChannelID)
	if err != nil {
		return nil, fmt.Errorf("invalid channel_id: %w", err)
	}

	msgID, err := discord.ParseSnowflake(nativeMessageID)
	if err != nil {
		return nil, fmt.Errorf("invalid native_message_id: %w", err)
	}

	embed, components, err := buildDiscordMessage(payload)
	if err != nil {
		return nil, err
	}

	var editData api.EditMessageData
	editData.Content = option.NewNullableString(payload.Message)
	if embed != nil {
		editData.Embeds = embed
	}
	if components != nil {
		editData.Components = components
	}

	if _, err := d.state.EditMessageComplex(discord.ChannelID(channelID), discord.MessageID(msgID), editData); err != nil {
		return nil, fmt.Errorf("failed to edit message: %w", err)
	}

	// Clear old emoji reactions and re-add current ones.
	if err := d.state.DeleteAllReactions(discord.ChannelID(channelID), discord.MessageID(msgID)); err != nil {
		log.Printf("[edit_embed] failed to clear reactions on message %s: %v", nativeMessageID, err)
	}
	for _, it := range payload.Interactions {
		if it.Type == models.InteractionTypeEmoji && it.Emoji != "" {
			if err := d.state.React(discord.ChannelID(channelID), discord.MessageID(msgID), discord.APIEmoji(it.Emoji)); err != nil {
				log.Printf("[edit_embed] failed to add reaction %s: %v", it.Emoji, err)
			}
		}
	}

	payloadJSON, err := payload.ToJSON()
	if err != nil {
		return nil, err
	}

	embedModel := &models.DiscordRoleReactionEmbed{
		NativeMessageId: nativeMessageID,
		Name:            payload.Message,
		Payload:         payloadJSON,
		Mode:            string(payload.Mode),
	}

	data, err := d.RoleReactionRepo.GetByNativeID(nativeMessageID)
	if err != nil || data == nil {
		return d.RoleReactionRepo.Create(embedModel)
	}
	return d.RoleReactionRepo.Update(nativeMessageID, embedModel)
}

// PublishEmbed sends the composed message to Discord, stores the configuration
// in the database keyed by the returned native message id, and adds emoji
// reactions for emoji interactions.
func (d *discordRoleReactionEmbedService) PublishEmbed(payload *models.ReactionRoleMessagePayload) (*models.DiscordRoleReactionEmbed, error) {
	if payload == nil {
		return nil, errors.New("payload is required")
	}
	if err := payload.Validate(); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	channelID, err := discord.ParseSnowflake(payload.ChannelID)
	if err != nil {
		return nil, fmt.Errorf("invalid channel_id: %w", err)
	}

	embed, components, err := buildDiscordMessage(payload)
	if err != nil {
		return nil, err
	}

	var msgData api.SendMessageData
	msgData.Content = payload.Message
	if embed != nil {
		msgData.Embeds = *embed
	}
	if components != nil {
		msgData.Components = *components
	}
	msg, err := d.state.SendMessageComplex(discord.ChannelID(channelID), msgData)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Add emoji reactions for emoji interactions.
	for _, it := range payload.Interactions {
		if it.Type == models.InteractionTypeEmoji && it.Emoji != "" {
			if err := d.state.React(discord.ChannelID(channelID), msg.ID, discord.APIEmoji(it.Emoji)); err != nil {
				log.Printf("[publish_embed] failed to add reaction %s: %v", it.Emoji, err)
			}
		}
	}

	payloadJSON, err := payload.ToJSON()
	if err != nil {
		return nil, err
	}

	embedModel := &models.DiscordRoleReactionEmbed{
		NativeMessageId: msg.ID.String(),
		Name:            payload.Message,
		Payload:         payloadJSON,
		Mode:            string(payload.Mode),
		Version:         1,
	}

	return d.UpsertEmbed(embedModel, nil)
}

func buildDiscordMessage(payload *models.ReactionRoleMessagePayload) (*[]discord.Embed, *discord.ContainerComponents, error) {
	var embeds *[]discord.Embed
	if payload.Embed != nil {
		e := discord.Embed{
			Title:       payload.Embed.Title,
			Description: payload.Embed.Description,
			Color:       discord.Color(payload.Embed.Color),
		}
		if payload.Embed.ImageURL != "" {
			e.Image = &discord.EmbedImage{URL: payload.Embed.ImageURL}
		}
		if payload.Embed.ThumbnailURL != "" {
			e.Thumbnail = &discord.EmbedThumbnail{URL: payload.Embed.ThumbnailURL}
		}
		if payload.Embed.Footer != "" {
			e.Footer = &discord.EmbedFooter{Text: payload.Embed.Footer}
		}
		if payload.Embed.Author != "" {
			e.Author = &discord.EmbedAuthor{Name: payload.Embed.Author}
		}
		for _, f := range payload.Embed.Fields {
			e.Fields = append(e.Fields, discord.EmbedField{
				Name:   f.Name,
				Value:  f.Value,
				Inline: f.Inline,
			})
		}
		embeds = &[]discord.Embed{e}
	}

	var rows discord.ContainerComponents
	for _, it := range payload.Interactions {
		switch it.Type {
		case models.InteractionTypeButton:
			btn := discord.ButtonComponent{
				Label:    it.Label,
				CustomID: discord.ComponentID(it.ID),
				Style:    buttonStyleToDiscord(it.Style),
			}
			if it.Emoji != "" {
				btn.Emoji = &discord.ComponentEmoji{Name: it.Emoji}
			}
			rows = append(rows, &discord.ActionRowComponent{&btn})
		case models.InteractionTypeDropdown:
			var opts []discord.SelectOption
			for _, opt := range it.Options {
				o := discord.SelectOption{
					Label: opt.Label,
					Value: opt.ID,
				}
				if opt.Description != "" {
					o.Description = opt.Description
				}
				if opt.Emoji != "" {
					o.Emoji = &discord.ComponentEmoji{Name: opt.Emoji}
				}
				opts = append(opts, o)
			}
			selectMenu := discord.StringSelectComponent{
				CustomID:    discord.ComponentID(it.ID),
				Placeholder: it.Placeholder,
				Options:     opts,
			}
			rows = append(rows, &discord.ActionRowComponent{&selectMenu})
		}
	}

	if len(rows) == 0 {
		return embeds, nil, nil
	}
	return embeds, &rows, nil
}

func buttonStyleToDiscord(s models.ButtonStyle) discord.ButtonComponentStyle {
	switch s {
	case models.ButtonStylePrimary:
		return discord.PrimaryButtonStyle()
	case models.ButtonStyleSuccess:
		return discord.SuccessButtonStyle()
	case models.ButtonStyleDanger:
		return discord.DangerButtonStyle()
	default:
		return discord.SecondaryButtonStyle()
	}
}

// PrettyPrintPayload is a small helper used by controllers/tests to return a
// readable version of the stored JSON payload.
func PrettyPrintPayload(payload string) (string, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		return "", err
	}
	b, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
