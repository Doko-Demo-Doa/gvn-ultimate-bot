package bot

import "github.com/bwmarrin/discordgo"

const THRESHOLD = 1

func RegisterPinModule(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, mr *discordgo.MessageReactionAdd) {
		msg, _ := s.ChannelMessage(mr.ChannelID, mr.MessageID)
		ProcessMessage(s, msg)
	})
}

func ProcessMessage(s *discordgo.Session, msg *discordgo.Message) {
	pin_count := 0

	for i := 0; i < len(msg.Reactions); i++ {
		if msg.Reactions[i].Emoji.Name == "📌" {
			pin_count += 1
		}
	}

	if pin_count >= THRESHOLD {
		s.ChannelMessagePin(msg.ChannelID, msg.ID)
	} else {
		s.ChannelMessageUnpin(msg.ChannelID, msg.ID)
	}
}
