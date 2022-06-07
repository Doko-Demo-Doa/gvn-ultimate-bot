package bot

import (
	"os"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

var (
	DiscordAppID   = os.Getenv("DISCORD_APP_ID")
	DiscordGuildID = os.Getenv("DISCORD_GUILD_ID")
)

func RegisterGrantRoleModule(s *state.State) {
	s.AddHandler(func(e *gateway.InteractionCreateEvent) {

	})
}
