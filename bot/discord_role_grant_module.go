package bot

import (
	"doko/gvn-ultimate-bot/services/discordservice"
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

var (
	grantRoleTarget   = "target"
	grantRoleName     = "role-name"
	grantRoleDuration = "duration"
)

func RegisterGrantRoleModule(s *state.State, ds discordservice.DiscordService) {
	var (
		AppID   = discord.AppID(mustSnowflakeEnv("DISCORD_APP_ID"))
		GuildID = discord.GuildID(mustSnowflakeEnv("DISCORD_GUILD_ID"))
	)

	s.AddHandler(func(e *gateway.InteractionCreateEvent) {
		switch data := e.InteractionEvent.Data.(type) {
		case *discord.CommandInteraction:
			switch data.Name {
			case "grant-role":
				handleGrantRole(s, ds, e, data, GuildID)
			case "quick-generate-role":
				handleQuickGenerateRole(s, e, data, GuildID)
			case "vkbp":
				handleVkbp(s, e, data, GuildID)
			}
		default:
			log.Println("Not supported interaction type")
		}
	})

	newCommands := []api.CreateCommandData{
		{
			Name:        "vkbp",
			Description: "Bạn kiếp bất phục",
			Options: discord.CommandOptions{
				&discord.UserOption{
					OptionName:  "target",
					Description: "Chém, cho thành Covid",
					Required:    true,
				},
				&discord.NumberOption{
					OptionName:  "ban-duration",
					Description: "Số ngày bị ban",
					Required:    true,
					Choices: []discord.NumberChoice{
						{Name: "1 ngày", Value: 1},
						{Name: "2 ngày", Value: 2},
						{Name: "3 ngày", Value: 3},
						{Name: "4 ngày", Value: 4},
						{Name: "5 ngày", Value: 5},
						{Name: "6 ngày", Value: 6},
						{Name: "7 ngày", Value: 7},
					},
				},
				&discord.StringOption{
					OptionName:  "reason",
					Description: "Lý do covid",
					Required:    true,
				},
			},
		},
		{
			Name:        "quick-generate-role",
			Description: "Tạo role, có thể có thời hạn",
			Options: discord.CommandOptions{
				&discord.StringOption{
					OptionName:  "role-name",
					Description: "Tên của role, dạng game-name",
					Required:    true,
					MinLength:   option.NewInt(2),
				},
				&discord.StringOption{
					OptionName:  "duration",
					Description: "Thời hạn của role, dạng 30m, 2h, 7d, 1w. Bỏ trống để không có thời hạn",
				},
			},
		},
		{
			Name:        "grant-role",
			Description: "Gán role cho ai đó (có thể có thời hạn)",
			Options: discord.CommandOptions{
				&discord.UserOption{
					OptionName:  "target",
					Description: "Người được gán role",
					Required:    true,
				},
				&discord.RoleOption{
					OptionName:  "role-name",
					Description: "Role mà member sẽ được gán",
					Required:    true,
				},
				&discord.StringOption{
					OptionName:  "duration",
					Description: "Thời hạn của role, dạng 30m, 2h, 7d, 1w. Bỏ trống để không có thời hạn",
				},
			},
		},
	}

	if _, err := s.BulkOverwriteGuildCommands(AppID, GuildID, newCommands); err != nil {
		log.Fatalln("Failed to create guild commands", err)
	}
}

func handleGrantRole(s *state.State, ds discordservice.DiscordService, e *gateway.InteractionCreateEvent, data *discord.CommandInteraction, guildID discord.GuildID) {
	targetOpt := data.Options.Find(grantRoleTarget)
	roleOpt := data.Options.Find(grantRoleName)
	durationOpt := data.Options.Find(grantRoleDuration)

	if targetOpt.Name == "" || roleOpt.Name == "" {
		respondError(s, e, "Thiếu thông tin user hoặc role")
		return
	}

	userSnow, err := targetOpt.SnowflakeValue()
	if err != nil {
		respondError(s, e, "Không thể đọc ID user")
		return
	}
	userID := discord.UserID(userSnow)

	roleSnow, err := roleOpt.SnowflakeValue()
	if err != nil {
		respondError(s, e, "Không thể đọc ID role")
		return
	}
	roleID := discord.RoleID(roleSnow)

	usr, err := s.User(userID)
	if err != nil {
		respondError(s, e, "Không thể lấy thông tin user")
		return
	}

	// Default duration = permanent (0)
	var duration time.Duration
	var durationText = "Vĩnh viễn"

	if durationOpt.Name != "" && len(durationOpt.Value) > 0 {
		d, err := ParseDuration(durationOpt.String())
		if err != nil {
			respondError(s, e, fmt.Sprintf("Định dạng thời hạn không hợp lệ: %s", err.Error()))
			return
		}
		duration = d
		durationText = durationOpt.String()
	}

	// Add role on Discord
	if err := s.AddRole(guildID, userID, roleID, api.AddRoleData{}); err != nil {
		respondError(s, e, fmt.Sprintf("Không thể gán role trên Discord: %s", err.Error()))
		return
	}

	// Track in database (even for permanent, so we have history)
	if duration > 0 {
		_, err = ds.AssignRoleToUser(userID.String(), roleID.String(), duration)
		if err != nil {
			log.Println("[grant-role] failed to track assignment in DB:", err)
			// Non-fatal: role is already granted on Discord
		}
	}

	embed := discord.Embed{
		Title:       "Đã gán role",
		Color:       discord.Color(3066993), // Green
		Description: fmt.Sprintf("**User:** %s\n**Role:** <@&%s>\n**Thời hạn:** %s", usr.Username, roleID.String(), durationText),
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %s", usr.ID.String()),
		},
	}

	respondEmbed(s, e, embed)
}

func handleQuickGenerateRole(s *state.State, e *gateway.InteractionCreateEvent, data *discord.CommandInteraction, guildID discord.GuildID) {
	respondError(s, e, "Lệnh này đang được phát triển")
}

func handleVkbp(s *state.State, e *gateway.InteractionCreateEvent, data *discord.CommandInteraction, guildID discord.GuildID) {
	var targetOpt, reasonOpt discord.CommandInteractionOption
	var user discord.User

	parsedOptions := data.Options
	targetOpt = parsedOptions.Find("target")
	reasonOpt = parsedOptions.Find("reason")

	snow, err := targetOpt.SnowflakeValue()
	if err != nil {
		log.Fatalf("Cannot parse Snowflake")
		return
	}
	usr, err := s.User(discord.UserID(snow))
	if err != nil {
		log.Fatalf("Cannot get user")
		return
	}
	user = *usr

	respondEmbed(s, e, discord.Embed{
		Title:       "User Banned",
		Color:       discord.Color(15418782), // Fuchsia
		Description: fmt.Sprintf("Đã ban %s, lý do: %s", user.Username, reasonOpt.Value),
		Footer: &discord.EmbedFooter{
			Text: "Hạn ban: 7 ngày",
		},
	})
}

func respondError(s *state.State, e *gateway.InteractionCreateEvent, msg string) {
	respondEmbed(s, e, discord.Embed{
		Title:       "Lỗi",
		Color:       discord.Color(15158332), // Red
		Description: msg,
	})
}

func respondEmbed(s *state.State, e *gateway.InteractionCreateEvent, embed discord.Embed) {
	data := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Embeds: &[]discord.Embed{embed},
		},
	}
	if err := s.RespondInteraction(e.ID, e.Token, data); err != nil {
		log.Println("Failed to send interaction callback:", err)
	}
}
