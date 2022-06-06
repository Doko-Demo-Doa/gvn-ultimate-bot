package bot

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

func StartRoleScanner(s *state.State) {
	roles, err := s.Roles(discord.GuildID(mustSnowflakeEnv("DISCORD_GUILD_ID")))

	log.Println(len(roles), err)
	// Insert roles into DB
}
