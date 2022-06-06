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

	s.AddHandler(func(m *gateway.MessageCreateEvent) {
		log.Printf("%s: %s", m.Author.Username, m.Content)
	})

	_, err := s.CurrentApplication()
	if err != nil {
		log.Fatalln("Failed to get application ID: ", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := s.Open(ctx); err != nil {
		log.Println("Cannot close: ", err)
	}

	select {}
}
