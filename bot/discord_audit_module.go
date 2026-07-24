package bot

import (
	"doko/gvn-ultimate-bot/services/discordservice"
	"log"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

// cachedMessage holds the content we need to log when a message is deleted or edited.
type cachedMessage struct {
	Content     string
	AuthorID    string
	AuthorName  string
	GuildId     string
	Attachments []string
	Timestamp   time.Time
}

// messageCache keeps the last N messages in memory so we can audit log deletions
// and edits. It is a simple thread-safe map with TTL.
type messageCache struct {
	mu      sync.RWMutex
	items   map[discord.MessageID]cachedMessage
	maxAge  time.Duration
	maxSize int
}

func newMessageCache(maxSize int, maxAge time.Duration) *messageCache {
	mc := &messageCache{
		items:   make(map[discord.MessageID]cachedMessage),
		maxAge:  maxAge,
		maxSize: maxSize,
	}
	go mc.gc()
	return mc
}

func (mc *messageCache) set(id discord.MessageID, msg cachedMessage) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	// Simple eviction: if over maxSize, clear half the map
	if len(mc.items) >= mc.maxSize {
		mc.evictHalf()
	}
	mc.items[id] = msg
}

func (mc *messageCache) get(id discord.MessageID) (cachedMessage, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	msg, ok := mc.items[id]
	if !ok {
		return cachedMessage{}, false
	}
	if time.Since(msg.Timestamp) > mc.maxAge {
		return cachedMessage{}, false
	}
	return msg, true
}

func (mc *messageCache) evictHalf() {
	cut := len(mc.items) / 2
	i := 0
	for k := range mc.items {
		if i >= cut {
			break
		}
		delete(mc.items, k)
		i++
	}
}

func (mc *messageCache) gc() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		mc.mu.Lock()
		now := time.Now()
		for id, msg := range mc.items {
			if now.Sub(msg.Timestamp) > mc.maxAge {
				delete(mc.items, id)
			}
		}
		mc.mu.Unlock()
	}
}

func attachmentURLs(m *gateway.MessageCreateEvent) []string {
	urls := make([]string, 0, len(m.Attachments))
	for _, a := range m.Attachments {
		urls = append(urls, a.URL)
	}
	return urls
}

func attachmentURLsFromMessage(m *gateway.MessageUpdateEvent) []string {
	urls := make([]string, 0, len(m.Attachments))
	for _, a := range m.Attachments {
		urls = append(urls, a.URL)
	}
	return urls
}

// RegisterAuditModule hooks message create / delete / update events and
// writes audit logs when messages are deleted or edited.
func RegisterAuditModule(s *state.State, svc discordservice.DiscordAuditLogService, guildID discord.GuildID) {
	mc := newMessageCache(5000, 24*time.Hour)

	// Cache every incoming message
	s.AddHandler(func(e *gateway.MessageCreateEvent) {
		if e.GuildID != guildID {
			return
		}
		authorName := ""
		if e.Member != nil && e.Member.Nick != "" {
			authorName = e.Member.Nick
		} else if e.Author.ID != 0 {
			authorName = e.Author.Username
		}
		mc.set(e.ID, cachedMessage{
			Content:     e.Content,
			AuthorID:    e.Author.ID.String(),
			AuthorName:  authorName,
			GuildId:     e.GuildID.String(),
			Attachments: attachmentURLs(e),
			Timestamp:   time.Now(),
		})
	})

	// Log deletions
	s.AddHandler(func(e *gateway.MessageDeleteEvent) {
		if e.GuildID != guildID {
			return
		}
		msg, ok := mc.get(e.ID)
		if !ok {
			return
		}
		if err := svc.LogMessageDelete(
			e.ID.String(),
			e.ChannelID.String(),
			e.GuildID.String(),
			msg.AuthorID,
			msg.AuthorName,
			msg.Content,
			msg.Attachments,
		); err != nil {
			log.Printf("[audit_module] failed to log delete: %v", err)
		}
	})

	// Log edits
	s.AddHandler(func(e *gateway.MessageUpdateEvent) {
		old, ok := mc.get(e.ID)
		if !ok {
			return
		}
		// Use the best available content: partial updates may leave e.Content empty
		newContent := e.Content
		if newContent == "" {
			newContent = old.Content
		}
		// Skip if nothing meaningful changed
		if old.Content == newContent && len(old.Attachments) == len(attachmentURLsFromMessage(e)) {
			return
		}
		authorName := old.AuthorName
		if e.Member != nil && e.Member.Nick != "" {
			authorName = e.Member.Nick
		} else if e.Author.ID != 0 {
			authorName = e.Author.Username
		}
		if err := svc.LogMessageEdit(
			e.ID.String(),
			e.ChannelID.String(),
			old.GuildId,
			old.AuthorID,
			authorName,
			old.Content,
			newContent,
			attachmentURLsFromMessage(e),
		); err != nil {
			log.Printf("[audit_module] failed to log edit: %v", err)
		}
		// Update cache with new content
		mc.set(e.ID, cachedMessage{
			Content:     newContent,
			AuthorID:    old.AuthorID,
			AuthorName:  authorName,
			GuildId:     old.GuildId,
			Attachments: attachmentURLsFromMessage(e),
			Timestamp:   time.Now(),
		})
	})
}
