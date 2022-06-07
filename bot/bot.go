package bot

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"gorm.io/gorm"
)

func Bootstrap(db *gorm.DB) {
	token := os.Getenv("DISCORD_TOKEN")
	s := state.New("Bot " + token)
	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentGuildMessages)

	// Reset all commands
	commands, err := s.GuildCommands(AppID, GuildID)
	if err != nil {
		log.Fatalf("Cannot get guild commands")
	}

	log.Println("Found these commands, unregistering...")
	for _, c := range commands {
		log.Println(c.Name)
		s.DeleteGuildCommand(AppID, GuildID, c.ID)
	}

	app, err := s.CurrentApplication()
	if err != nil {
		log.Fatalln("Failed to get application ID: ", err)
	}
	log.Println("App ID", app.ID)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := s.Open(ctx); err != nil {
		log.Println("Cannot close: ", err)
	}

	select {}
}
