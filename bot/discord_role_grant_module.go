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
			Name:        "tiaraping",
			Description: "Tiaramisu ping",
		},
	}

	if _, err := s.BulkOverwriteGuildCommands(AppID, GuildID, newCommands); err != nil {
		log.Fatalln("Failed to create guild commands", err)
	}
}
