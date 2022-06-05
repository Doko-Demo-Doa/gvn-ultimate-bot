package bot

import (
	"flag"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID  = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken = flag.String("token", "", "Bot access token")
)

var (
	DiscordAppID   = os.Getenv("DISCORD_APP_ID")
	DiscordGuildID = os.Getenv("DISCORD_GUILD_ID")
)

var (
	commands = []*discordgo.ApplicationCommand{

		{
			Name:        "birdup",
			Description: "Followup messages",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"birdup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
		},
	}
)

func RegisterGrantRoleModule(s *discordgo.Session) {
	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(DiscordAppID, DiscordGuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

func UnregisterGrantRoleModule(s *discordgo.Session) {
	cmdList, er := s.ApplicationCommands(DiscordAppID, DiscordGuildID)
	if er != nil {
		log.Panicf("Cannot get command list")
		return
	}

	for _, v := range cmdList {
		err := s.ApplicationCommandDelete(DiscordAppID, DiscordGuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Shutting shown...")
}
