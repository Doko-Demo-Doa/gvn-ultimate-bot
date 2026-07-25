package discordservice

import (
	"doko/gvn-ultimate-bot/models"
	discordrepos "doko/gvn-ultimate-bot/repositories/discord_repos"
	"doko/gvn-ultimate-bot/services/systemservice"
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

type ChannelInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     uint   `json:"type"`
	Position int    `json:"position"`
}

type EmojiInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Animated bool   `json:"animated"`
	ImageURL string `json:"image_url"`
	APIName  string `json:"api_name"`
}

type MemberInfo struct {
	NativeId string `json:"native_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type DiscordRoleReactionEmbedService interface {
	ListEmbeds() ([]*models.DiscordRoleReactionEmbed, error)
	UpsertEmbed(*models.DiscordRoleReactionEmbed, *models.ReactionRoleMessagePayload) (*models.DiscordRoleReactionEmbed, error)
	GetSingleEmbed(id uint) (*models.DiscordRoleReactionEmbed, error)
	GetEmbedByNativeMessageID(nativeMessageID string) (*models.DiscordRoleReactionEmbed, error)
	DeleteEmbed(id uint) error
	EditEmbed(nativeMessageID string, payload *models.ReactionRoleMessagePayload) (*models.DiscordRoleReactionEmbed, error)
	PublishEmbed(*models.ReactionRoleMessagePayload) (*models.DiscordRoleReactionEmbed, error)
	ListChannels() ([]ChannelInfo, error)
	ListEmojis() ([]EmojiInfo, error)
	SearchGuildMembers(query string) ([]MemberInfo, error)
	SyncGuildMembers() (*UserSyncResult, error)
	GetLastUserSync() (*models.SystemEventLog, error)
	SyncGuildRoles() (*RoleSyncResult, error)
	GetLastRoleSync() (*models.SystemEventLog, error)
}

type UserSyncResult struct {
	SyncedCount  int `json:"synced_count"`
	RemovedCount int `json:"removed_count"`
}

type RoleSyncResult struct {
	SyncedCount  int `json:"synced_count"`
	RemovedCount int `json:"removed_count"`
}

type discordRoleReactionEmbedService struct {
	RoleReactionRepo discordrepos.DiscordRoleReactionEmbedRepo
	UserRepo         discordrepos.DiscordUserRepo
	RoleRepo         discordrepos.DiscordRoleRepo
	EventLogService  systemservice.SystemEventLogService
	state            *state.State
	guildID          discord.GuildID
}

func NewDiscordRoleReactionEmbedService(
	repo discordrepos.DiscordRoleReactionEmbedRepo,
	userRepo discordrepos.DiscordUserRepo,
	roleRepo discordrepos.DiscordRoleRepo,
	eventLogService systemservice.SystemEventLogService,
	s *state.State,
	guildID discord.GuildID,
) DiscordRoleReactionEmbedService {
	return &discordRoleReactionEmbedService{
		RoleReactionRepo: repo,
		UserRepo:         userRepo,
		RoleRepo:         roleRepo,
		EventLogService:  eventLogService,
		state:            s,
		guildID:          guildID,
	}
}

func (d *discordRoleReactionEmbedService) ListEmbeds() ([]*models.DiscordRoleReactionEmbed, error) {
	return d.RoleReactionRepo.ListRoleReactionEmbeds()
}

func (d *discordRoleReactionEmbedService) ListChannels() ([]ChannelInfo, error) {
	channels, err := d.state.Channels(d.guildID)
	if err != nil {
		return nil, err
	}
	var result []ChannelInfo
	for _, ch := range channels {
		// Only include text channels (Type 0) and announcement channels (Type 5)
		if ch.Type == discord.GuildText || ch.Type == discord.GuildNews {
			result = append(result, ChannelInfo{
				ID:       ch.ID.String(),
				Name:     ch.Name,
				Type:     uint(ch.Type),
				Position: ch.Position,
			})
		}
	}
	return result, nil
}

func (d *discordRoleReactionEmbedService) ListEmojis() ([]EmojiInfo, error) {
	emojis, err := d.state.Emojis(d.guildID)
	if err != nil {
		return nil, err
	}
	var result []EmojiInfo
	for _, em := range emojis {
		// Only custom emojis (ID != 0)
		if em.ID == 0 {
			continue
		}
		ext := "png"
		if em.Animated {
			ext = "gif"
		}
		apiName := em.Name + ":" + em.ID.String()
		result = append(result, EmojiInfo{
			ID:       em.ID.String(),
			Name:     em.Name,
			Animated: em.Animated,
			ImageURL: fmt.Sprintf("https://cdn.discordapp.com/emojis/%s.%s", em.ID.String(), ext),
			APIName:  apiName,
		})
	}
	return result, nil
}

// SearchGuildMembers searches the locally-synced discord_user table (kept up
// to date by SyncGuildMembers and the member-sync bot module), rather than
// Discord's gateway member cache, which only contains members observed
// through live events and is incomplete on large guilds.
func (d *discordRoleReactionEmbedService) SearchGuildMembers(query string) ([]MemberInfo, error) {
	users, err := d.UserRepo.Search(query, 25)
	if err != nil {
		return nil, err
	}
	result := make([]MemberInfo, 0, len(users))
	for _, u := range users {
		result = append(result, MemberInfo{
			NativeId: u.NativeID,
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
		})
	}
	return result, nil
}

// SyncGuildMembers fetches the full guild roster via the REST API (the
// gateway member cache only contains members observed through live events,
// which is incomplete on large guilds), upserts them into the discord_user
// table, removes users no longer in the guild, and logs a SystemEventLog
// entry (USER_SYNC) with the result.
func (d *discordRoleReactionEmbedService) SyncGuildMembers() (*UserSyncResult, error) {
	members, err := d.state.AllMembers(d.guildID)
	if err != nil {
		d.logUserSync("failure", err.Error(), 0, 0)
		return nil, err
	}

	nativeIds := make([]string, 0, len(members))
	synced := 0
	for _, m := range members {
		nativeIds = append(nativeIds, m.User.ID.String())
		_, err := d.UserRepo.Upsert(&models.DiscordUser{
			NativeID:      m.User.ID.String(),
			Discriminator: m.User.Discriminator,
			Avatar:        m.User.AvatarURL(),
			Username:      m.User.Username,
			Nickname:      m.Nick,
		})
		if err != nil {
			log.Printf("[sync_guild_members] failed to upsert user %s: %v", m.User.ID.String(), err)
			continue
		}
		synced++
	}

	removed, err := d.UserRepo.DeleteNotIn(nativeIds)
	if err != nil {
		log.Printf("[sync_guild_members] failed to remove stale users: %v", err)
	}

	d.logUserSync("success", fmt.Sprintf("synced %d/%d members, removed %d stale", synced, len(members), removed), synced, int(removed))
	return &UserSyncResult{SyncedCount: synced, RemovedCount: int(removed)}, nil
}

func (d *discordRoleReactionEmbedService) logUserSync(status, message string, syncedCount, removedCount int) {
	if d.EventLogService == nil {
		return
	}
	metadata := models.MarshalJSONColumn(map[string]int{"synced_count": syncedCount, "removed_count": removedCount})
	if err := d.EventLogService.LogEvent(models.SystemEventTypeUserSync, status, message, metadata); err != nil {
		log.Printf("[sync_guild_members] failed to write event log: %v", err)
	}
}

func (d *discordRoleReactionEmbedService) GetLastUserSync() (*models.SystemEventLog, error) {
	if d.EventLogService == nil {
		return nil, nil
	}
	return d.EventLogService.GetLatestByEventType(models.SystemEventTypeUserSync)
}

// SyncGuildRoles fetches the guild's full role list (always complete in
// arikawa's cache, unlike members) and reconciles the discord_role table
// against it: upserting Discord-sourced fields for every existing role and
// removing any role no longer present in the guild.
func (d *discordRoleReactionEmbedService) SyncGuildRoles() (*RoleSyncResult, error) {
	roles, err := d.state.Roles(d.guildID)
	if err != nil {
		d.logRoleSync("failure", err.Error(), 0, 0)
		return nil, err
	}

	nativeIds := make([]string, 0, len(roles))
	synced := 0
	for _, r := range roles {
		nativeIds = append(nativeIds, r.ID.String())
		_, err := d.RoleRepo.Upsert(&models.DiscordRole{
			NativeID:    r.ID.String(),
			Name:        r.Name,
			Mentionable: boolToUint(r.Mentionable),
			Hoist:       boolToUint(r.Hoist),
			Color:       uint(r.Color),
		})
		if err != nil {
			log.Printf("[sync_guild_roles] failed to upsert role %s: %v", r.ID.String(), err)
			continue
		}
		synced++
	}

	removed, err := d.RoleRepo.DeleteNotIn(nativeIds)
	if err != nil {
		log.Printf("[sync_guild_roles] failed to remove stale roles: %v", err)
	}

	d.logRoleSync("success", fmt.Sprintf("synced %d/%d roles, removed %d stale", synced, len(roles), removed), synced, int(removed))
	return &RoleSyncResult{SyncedCount: synced, RemovedCount: int(removed)}, nil
}

func boolToUint(b bool) uint {
	if b {
		return 1
	}
	return 0
}

func (d *discordRoleReactionEmbedService) logRoleSync(status, message string, syncedCount, removedCount int) {
	if d.EventLogService == nil {
		return
	}
	metadata := models.MarshalJSONColumn(map[string]int{"synced_count": syncedCount, "removed_count": removedCount})
	if err := d.EventLogService.LogEvent(models.SystemEventTypeRoleSync, status, message, metadata); err != nil {
		log.Printf("[sync_guild_roles] failed to write event log: %v", err)
	}
}

func (d *discordRoleReactionEmbedService) GetLastRoleSync() (*models.SystemEventLog, error) {
	if d.EventLogService == nil {
		return nil, nil
	}
	return d.EventLogService.GetLatestByEventType(models.SystemEventTypeRoleSync)
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

	data, err := d.RoleReactionRepo.GetByNativeID(m.NativeMessageID)
	if err != nil || data == nil {
		return d.RoleReactionRepo.Create(m)
	}
	return d.RoleReactionRepo.Update(m.NativeMessageID, m)
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
	if err == nil && payload != nil && payload.ChannelID != "" && embed.NativeMessageID != "" {
		channelID, err := discord.ParseSnowflake(payload.ChannelID)
		if err == nil {
			msgID, parseErr := discord.ParseSnowflake(embed.NativeMessageID)
			if parseErr == nil {
				if delErr := d.state.DeleteMessage(discord.ChannelID(channelID), discord.MessageID(msgID), api.AuditLogReason("")); delErr != nil {
					log.Printf("[delete_embed] failed to delete Discord message %s: %v", embed.NativeMessageID, delErr)
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

	// Snapshot which emojis were reactable before this edit, so afterwards we
	// only touch the ones that actually changed rather than wiping every
	// reaction on the message.
	oldEmojis := map[string]bool{}
	if existing, err := d.RoleReactionRepo.GetByNativeID(nativeMessageID); err == nil && existing != nil {
		if oldPayload, err := existing.ParsedPayload(); err == nil {
			for _, it := range oldPayload.Interactions {
				if it.Type == models.InteractionTypeEmoji && it.Emoji != "" {
					oldEmojis[it.Emoji] = true
				}
			}
		}
	}

	newEmojis := map[string]bool{}
	for _, it := range payload.Interactions {
		if it.Type == models.InteractionTypeEmoji && it.Emoji != "" {
			newEmojis[it.Emoji] = true
		}
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

	payloadJSON, err := payload.ToJSON()
	if err != nil {
		return nil, err
	}

	embedModel := &models.DiscordRoleReactionEmbed{
		NativeMessageID: nativeMessageID,
		Name:            payload.Message,
		Payload:         payloadJSON,
		Mode:            string(payload.Mode),
	}

	// Persist before touching Discord reactions so the reaction listener can
	// always find the row once the message becomes reactable again.
	data, err := d.RoleReactionRepo.GetByNativeID(nativeMessageID)
	var saved *models.DiscordRoleReactionEmbed
	if err != nil || data == nil {
		saved, err = d.RoleReactionRepo.Create(embedModel)
	} else {
		saved, err = d.RoleReactionRepo.Update(nativeMessageID, embedModel)
	}
	if err != nil {
		return nil, err
	}

	// Remove reactions only for emojis dropped from the payload. This fires a
	// "reaction remove emoji" event, which the role-react listener doesn't
	// subscribe to, so it does NOT revoke roles from users who'd already
	// reacted — it only cleans up the stale reaction icon on Discord.
	for emoji := range oldEmojis {
		if newEmojis[emoji] {
			continue
		}
		if err := d.state.DeleteReactions(discord.ChannelID(channelID), discord.MessageID(msgID), discord.APIEmoji(emoji)); err != nil {
			log.Printf("[edit_embed] failed to remove stale reaction %s: %v", emoji, err)
		}
	}

	// Add the bot's own reaction for newly introduced emojis. Emojis present
	// both before and after are left completely untouched, preserving
	// existing user reactions.
	for emoji := range newEmojis {
		if oldEmojis[emoji] {
			continue
		}
		if err := d.state.React(discord.ChannelID(channelID), discord.MessageID(msgID), discord.APIEmoji(emoji)); err != nil {
			log.Printf("[edit_embed] failed to add reaction %s: %v", emoji, err)
		}
	}

	return saved, nil
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

	payloadJSON, err := payload.ToJSON()
	if err != nil {
		return nil, err
	}

	embedModel := &models.DiscordRoleReactionEmbed{
		NativeMessageID: msg.ID.String(),
		Name:            payload.Message,
		Payload:         payloadJSON,
		Mode:            string(payload.Mode),
		Version:         1,
	}

	// Persist before touching Discord reactions so the reaction listener can
	// always find the row once the message becomes reactable.
	saved, err := d.UpsertEmbed(embedModel, nil)
	if err != nil {
		return nil, err
	}

	// Add emoji reactions for emoji interactions.
	for _, it := range payload.Interactions {
		if it.Type == models.InteractionTypeEmoji && it.Emoji != "" {
			if err := d.state.React(discord.ChannelID(channelID), msg.ID, discord.APIEmoji(it.Emoji)); err != nil {
				log.Printf("[publish_embed] failed to add reaction %s: %v", it.Emoji, err)
			}
		}
	}

	return saved, nil
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
