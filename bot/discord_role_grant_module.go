package bot

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

var (
	TARGET       = "target"
	BAN_DURATION = "ban-duration"
	REASON       = "reason"
)

func RegisterGrantRoleModule(s *state.State) {
	var (
		AppID   = discord.AppID(mustSnowflakeEnv("DISCORD_APP_ID"))
		GuildID = discord.GuildID(mustSnowflakeEnv("DISCORD_GUILD_ID"))
	)

	s.AddHandler(func(e *gateway.InteractionCreateEvent) {
		// e.InteractionEvent.Data.InteractionType() = 2 => Command interaction
		var targetOpt, reasonOpt discord.CommandInteractionOption
		var user discord.User

		switch data := e.InteractionEvent.Data.(type) {
		case *discord.CommandInteraction:
			parsedOptions := data.Options

			targetOpt = parsedOptions.Find(TARGET)
			reasonOpt = parsedOptions.Find(REASON)

			snow, err := targetOpt.SnowflakeValue()
			if err != nil {
				return
			}

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

		default:
			log.Println("Not supported")
		}

		data := api.InteractionResponse{
			Type: api.MessageInteractionWithSource,
			Data: &api.InteractionResponseData{
				Embeds: &[]discord.Embed{
					{
						Title:       "User Banned",
						Color:       discord.Color(15418782), // Fuchsia
						Description: fmt.Sprintf("Đã ban %s, lý do: %s", user.Username, reasonOpt.Value),
						Footer: &discord.EmbedFooter{
							Text: "Hạn ban: 7 ngày",
						},
					},
				},
			},
		}

		if err := s.RespondInteraction(e.ID, e.Token, data); err != nil {
			log.Println("Failed to send interaction callback:", err)
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
						{
							Name:  "1 ngày",
							Value: 1,
						},
						{
							Name:  "2 ngày",
							Value: 2,
						},
						{
							Name:  "3 ngày",
							Value: 3,
						},
						{
							Name:  "4 ngày",
							Value: 4,
						},
						{
							Name:  "5 ngày",
							Value: 5,
						},
						{
							Name:  "6 ngày",
							Value: 6,
						},
						{
							Name:  "7 ngày",
							Value: 7,
						},
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
					OptionName:  "Tên role",
					Description: "Tên của role, dạng game-name",
					Required:    true,
					MinLength:   option.NewInt(2),
				},
				&discord.StringOption{
					OptionName:  "Thời hạn",
					Description: "Thời hạn của role, dạng YYYY-MM-DD, bỏ trống để không có thời hạn",
					MinLength:   option.NewInt(10),
				},
			},
		},
		{
			Name:        "grant-role",
			Description: "Gán role cho ai đó",
			Options: discord.CommandOptions{
				&discord.UserOption{
					OptionName:  "target",
					Description: "Người được gán role",
					Required:    true,
				},
				&discord.RoleOption{
					OptionName:  "Role",
					Description: "Role mà member sẽ được gán",
					Required:    true,
				},
				&discord.StringOption{
					OptionName:  "Thời hạn",
					Description: "Thời hạn của role, dạng YYYY-MM-DD, bỏ trống để không có thời hạn",
					MinLength:   option.NewInt(10),
				},
			},
		},
	}

	if _, err := s.BulkOverwriteGuildCommands(AppID, GuildID, newCommands); err != nil {
		log.Fatalln("Failed to create guild commands", err)
	}
}
