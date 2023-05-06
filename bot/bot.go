package bot

import (
	"context"
	"doko/gvn-ultimate-bot/services/discordservice"
	"doko/gvn-ultimate-bot/services/moduleservice"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"gorm.io/gorm"
)

func Bootstrap(db *gorm.DB, ds discordservice.DiscordService, ms moduleservice.ModuleService) {
	var (
		AppID    = discord.AppID(mustSnowflakeEnv("DISCORD_APP_ID"))
		GuildID  = discord.GuildID(mustSnowflakeEnv("DISCORD_GUILD_ID"))
		BotToken = os.Getenv("DISCORD_TOKEN")
	)

	s := state.New("Bot " + BotToken)
	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentGuildMessages)
	s.AddIntents(gateway.IntentGuildMessageReactions)

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
	// Will be disabled when enough data is provided
	StartRoleSync(s, ds)

	// Individual modules
	availableModules, err := ms.ListModules()
	if err == nil {
		for _, module := range availableModules {
			if module.ModuleName == "pin_module" && module.IsActivated == 1 {
				fmt.Println("Registering pin module...")
				RegisterPinModule(s)
			}
			if module.ModuleName == "grant_role_module" {
				RegisterRoleReactModule(s)
			}
			// TODO: Wip
			if module.ModuleName == "grant_role_command" {
				RegisterGrantRoleModule(s)
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
