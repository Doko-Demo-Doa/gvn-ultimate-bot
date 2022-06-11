package bot

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"gorm.io/gorm"
)

func Bootstrap(db *gorm.DB) {
	var (
		AppID    = discord.AppID(mustSnowflakeEnv("DISCORD_APP_ID"))
		GuildID  = discord.GuildID(mustSnowflakeEnv("DISCORD_GUILD_ID"))
		BotToken = os.Getenv("DISCORD_TOKEN")
	)

	s := state.New("Bot " + BotToken)
	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentGuildMessages)

	// Reset all commands
	commands, err := s.GuildCommands(AppID, GuildID)
	if err != nil {
		log.Fatalf("Cannot get guild commands")
	}

	log.Println("Found these commands, unregistering if needed...")
	for _, c := range commands {
		log.Println(c.Name)
		s.DeleteGuildCommand(AppID, GuildID, c.ID)
	}

	app, err := s.CurrentApplication()
	if err != nil {
		log.Fatalln("Failed to get application ID: ", err)
	}
	log.Println("App ID", app.ID)

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
