package bot

import (
	"doko/gvn-ultimate-bot/models"
	discordrepos "doko/gvn-ultimate-bot/repositories/discord_repos"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

// RegisterMemberSyncModule hooks guild member add / update / remove events and
// keeps the discord_user table in sync in real time, without requiring the
// manual "Sync Discord Users" backfill.
func RegisterMemberSyncModule(s *state.State, repo discordrepos.DiscordUserRepo, guildID discord.GuildID) {
	s.AddHandler(func(e *gateway.GuildMemberAddEvent) {
		if e.GuildID != guildID {
			return
		}
		upsertMember(repo, e.User, e.Nick)
	})

	s.AddHandler(func(e *gateway.GuildMemberUpdateEvent) {
		if e.GuildID != guildID {
			return
		}
		upsertMember(repo, e.User, e.Nick)
	})

	s.AddHandler(func(e *gateway.GuildMemberRemoveEvent) {
		if e.GuildID != guildID {
			return
		}
		if err := repo.DeleteByNativeID(e.User.ID.String()); err != nil {
			log.Printf("[member_sync_module] failed to delete user %s: %v", e.User.ID.String(), err)
		}
	})
}

func upsertMember(repo discordrepos.DiscordUserRepo, user discord.User, nick string) {
	_, err := repo.Upsert(&models.DiscordUser{
		NativeId:      user.ID.String(),
		Discriminator: user.Discriminator,
		Avatar:        user.AvatarURL(),
		Username:      user.Username,
		Nickname:      nick,
	})
	if err != nil {
		log.Printf("[member_sync_module] failed to upsert user %s: %v", user.ID.String(), err)
	}
}
