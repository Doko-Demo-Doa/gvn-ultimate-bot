package bot

import (
	"os"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

var (
	GuildID  = discord.GuildID(mustSnowflakeEnv("DISCORD_GUILD_ID"))
	BotToken = os.Getenv("DISCORD_TOKEN")
)

var (
	DiscordAppID   = os.Getenv("DISCORD_APP_ID")
	DiscordGuildID = os.Getenv("DISCORD_GUILD_ID")
)

func RegisterGrantRoleModule(s *state.State) {

}
