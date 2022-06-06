package bot

import (
	"flag"
	"os"

	"github.com/diamondburned/arikawa/v3/state"
)

var (
	GuildID  = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken = flag.String("token", "", "Bot access token")
)

var (
	DiscordAppID   = os.Getenv("DISCORD_APP_ID")
	DiscordGuildID = os.Getenv("DISCORD_GUILD_ID")
)

func RegisterGrantRoleModule(s *state.State) {

}
