package bot

import (
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func RegisterGrantRoleModule(s *state.State) {
	var (
		AppID   = discord.AppID(mustSnowflakeEnv("DISCORD_APP_ID"))
		GuildID = discord.GuildID(mustSnowflakeEnv("DISCORD_GUILD_ID"))
	)

	s.AddHandler(func(e *gateway.InteractionCreateEvent) {
		// e.InteractionEvent.Data.InteractionType() = 2 => Command interaction
		log.Println("222", e.InteractionEvent.Message)

		switch data := e.InteractionEvent.Data.(type) {
		case *discord.CommandInteraction:
			log.Println("", data.Name)

		default:
			log.Println("22")
		}

		data := api.InteractionResponse{
			Type: api.MessageInteractionWithSource,
			Data: &api.InteractionResponseData{
				Content: option.NewNullableString("Pong"),
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
					Description: "Number of days to covid",
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
				},
			},
		},
	}

	if _, err := s.BulkOverwriteGuildCommands(AppID, GuildID, newCommands); err != nil {
		log.Fatalln("Failed to create guild commands", err)
	}
}
