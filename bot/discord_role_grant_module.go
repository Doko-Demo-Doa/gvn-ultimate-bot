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
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{

		{
			Name:        "followups",
			Description: "Followup messages",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"followups": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
		},
	}
)

var registeredCommands []*discordgo.ApplicationCommand

func RegisterGrantRoleModule(s *discordgo.Session) {
	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))

	lst, err := s.ApplicationCommandBulkOverwrite(DiscordAppID, DiscordGuildID, commands)
	if err != nil {
		println(err.Error())
		return
	}

	registeredCommands = lst

	// for i, v := range commands {
	// 	cmd, err := s.ApplicationCommandCreate(DiscordAppID, DiscordGuildID, v)
	// 	if err != nil {
	// 		log.Panicf("Cannot create '%v' command: %v", v.Name, err)
	// 	}
	// 	registeredCommands[i] = cmd
	// }
}

func UnregisterGrantRoleModule(s *discordgo.Session) {
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(DiscordAppID, DiscordGuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Shutting shown...")
}
