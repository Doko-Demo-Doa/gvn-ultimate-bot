package bot

import (
	"log"
	"os"

	"github.com/diamondburned/arikawa/v3/discord"
)

var (
	AppID    = discord.AppID(mustSnowflakeEnv("DISCORD_GUILD_ID"))
	GuildID  = discord.GuildID(mustSnowflakeEnv("DISCORD_GUILD_ID"))
	BotToken = os.Getenv("DISCORD_TOKEN")
)

func mustSnowflakeEnv(env string) discord.Snowflake {
	s, err := discord.ParseSnowflake(os.Getenv(env))
	if err != nil {
		log.Fatalf("Invalid snowflake for $%s: %v", env, err)
	}
	return s
}
