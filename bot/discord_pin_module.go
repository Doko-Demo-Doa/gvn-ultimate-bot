package bot

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

const THRESHOLD = 1

func RegisterPinModule(s *state.State) {
	s.AddHandler(func(m *gateway.MessageReactionAddEvent) {
		msg, err := s.Message(m.ChannelID, m.MessageID)
		if err != nil {
			return
		}
		ProcessMessage(s, msg)
	})

	s.AddHandler(func(m *gateway.MessageReactionRemoveEvent) {
		msg, err := s.Message(m.ChannelID, m.MessageID)
		if err != nil {
			return
		}
		ProcessMessage(s, msg)
	})
}

func ProcessMessage(s *state.State, msg *discord.Message) {
	pin_count := 0
	pin_symbol := "📌"

	for i := 0; i < len(msg.Reactions); i++ {
		if msg.Reactions[i].Emoji.Name == pin_symbol {
			pin_count += 1
		}
	}

	if pin_count >= THRESHOLD {
		s.PinMessage(msg.ChannelID, msg.ID, "Pin message threshold reached")
	} else {
		s.UnpinMessage(msg.ChannelID, msg.ID, "Unpin message threshold reached")
	}
}
