package bot

import (
	"log"
	"os"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

var (
	DiscordAppID   = os.Getenv("DISCORD_APP_ID")
	DiscordGuildID = os.Getenv("DISCORD_GUILD_ID")
)

func RegisterGrantRoleModule(s *state.State) {
	s.AddHandler(func(e *gateway.InteractionCreateEvent) {
		data := api.InteractionResponse{
			Type: api.MessageInteractionWithSource,
			Data: &api.InteractionResponseData{
				Content: option.NewNullableString("Pong"),
			},
		}

		log.Println(data.Type)
	})
}
