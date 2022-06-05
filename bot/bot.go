package bot

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func Bootstrap() {
	token := os.Getenv("DISCORD_TOKEN")
	// s = Discord session.
	s, _ := discordgo.New("Bot " + token)

	// Register individual modules
	RegisterPinModule(s)
	RegisterGrantRoleModule(s)

	// Other configs
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	log.Println("Graceful shutdown")

	UnregisterGrantRoleModule(s)
}
