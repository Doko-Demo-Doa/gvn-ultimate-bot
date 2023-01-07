package bot

import (
	"doko/gvn-ultimate-bot/models"
	"doko/gvn-ultimate-bot/services/discordservice"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

func StartRoleSync(state *state.State, ds discordservice.DiscordService) {
	roles, err := state.Roles(discord.GuildID(mustSnowflakeEnv("DISCORD_GUILD_ID")))

	log.Println(len(roles), err)
	// Insert roles into DB
	for _, elem := range roles {
		var newRole = roleAdapter(elem)
		ds.CreateRole(&newRole)
	}
}

func roleAdapter(r1 discord.Role) models.DiscordRole {
	var mentionableNum uint = 0
	if r1.Mentionable {
		mentionableNum = 1
	} else {
		mentionableNum = 0
	}

	var isHoist uint = 0
	if r1.Hoist {
		isHoist = 1
	} else {
		isHoist = 0
	}

	return models.DiscordRole{
		NativeId:     r1.ID.String(),
		Name:         r1.Name,
		Mentionable:  mentionableNum,
		Hoist:        isHoist,
		Color:        uint(r1.Color),
		ImplicitType: 0,
	}
}
