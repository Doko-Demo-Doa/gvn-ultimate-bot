package bot

import (
	"context"
	discordrepos "doko/gvn-ultimate-bot/repositories/discord_repos"
	"doko/gvn-ultimate-bot/scheduler"
	"doko/gvn-ultimate-bot/services/discordservice"
	"doko/gvn-ultimate-bot/services/moduleservice"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

var (
	AppID     = discord.AppID(mustSnowflakeEnv("DISCORD_APP_ID"))
	GuildID   = discord.GuildID(mustSnowflakeEnv("DISCORD_GUILD_ID"))
	BotToken  = os.Getenv("DISCORD_TOKEN")
	IsWorking = false
)

func Bootstrap(s *state.State, ds discordservice.DiscordService, dre discordservice.DiscordRoleReactionEmbedService, ms moduleservice.ModuleService, rs *scheduler.RoleScheduler, als discordservice.DiscordAuditLogService, userRepo discordrepos.DiscordUserRepo) {
	// Reset all commands
	commands, err := s.GuildCommands(AppID, GuildID)
	if err != nil {
		log.Fatalf("Cannot get guild commands")
	}

	log.Printf("Found %d command(s), unregistering if needed...", len(commands))
	for _, c := range commands {
		log.Println("Command: ", c.Name)
		s.DeleteGuildCommand(AppID, GuildID, c.ID)
	}

	app, err := s.CurrentApplication()
	if err != nil {
		log.Fatalln("Failed to get application ID: ", err)
	}
	log.Println("App ID", app.ID)

	// Sync the roles into database
	StartRoleSync(s, ds)

	// Start background scheduler to remove expired timed roles
	rs.Start()

	// Mark the bot as "working"
	IsWorking = true

	// Individual modules
	availableModules, err := ms.ListModules()
	if err == nil {
		for _, module := range availableModules {
			if module.ModuleName == "pin_module" && module.IsActivated == 1 {
				fmt.Println("Registering pin module...")
				RegisterPinModule(s)
			}
			if module.ModuleName == "grant_role_module" {
				RegisterRoleReactModule(s, ds, dre, rs, GuildID)
			}
			if module.ModuleName == "grant_role_command" {
				RegisterGrantRoleModule(s, rs)
			}
			if module.ModuleName == "message_audit_module" && module.IsActivated == 1 {
				fmt.Println("Registering message audit module...")
				RegisterAuditModule(s, als, GuildID)
			}
			if module.ModuleName == "member_sync_module" && module.IsActivated == 1 {
				fmt.Println("Registering member sync module...")
				RegisterMemberSyncModule(s, userRepo, GuildID)
			}
		}
	}

	// Setup app context and interrupt channel
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := s.Open(ctx); err != nil {
		log.Println("Cannot close: ", err)
	}
}

func mustSnowflakeEnv(env string) discord.Snowflake {
	s, err := discord.ParseSnowflake(os.Getenv(env))
	if err != nil {
		log.Fatalf("Invalid snowflake for $%s: %v", env, err)
	}
	return s
}
