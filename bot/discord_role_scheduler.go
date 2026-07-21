package bot

import (
	"doko/gvn-ultimate-bot/services/discordservice"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

// StartRoleExpirationScheduler periodically checks for expired role assignments
// and removes the Discord roles from the affected users.
func StartRoleExpirationScheduler(s *state.State, ds discordservice.DiscordService) {
	guildID := discord.GuildID(mustSnowflakeEnv("DISCORD_GUILD_ID"))

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			expired, err := ds.GetExpiredRoleAssignments()
			if err != nil {
				log.Println("[role_scheduler] error fetching expired assignments:", err)
				continue
			}
			if len(expired) == 0 {
				continue
			}

			log.Printf("[role_scheduler] found %d expired assignment(s)", len(expired))

			for _, assignment := range expired {
				userID := discord.UserID(mustSnowflakeEnvOrDefault(assignment.UserNativeID, "0"))
				roleID := discord.RoleID(mustSnowflakeEnvOrDefault(assignment.RoleNativeID, "0"))

				if userID == 0 || roleID == 0 {
					log.Printf("[role_scheduler] skipping assignment %d: invalid IDs", assignment.ID)
					_ = ds.RevokeRoleAssignment(assignment.ID)
					continue
				}

				if err := s.RemoveRole(guildID, userID, roleID, api.AuditLogReason("")); err != nil {
					log.Printf("[role_scheduler] failed to remove role %s from user %s: %v", assignment.RoleNativeID, assignment.UserNativeID, err)
					// Do not delete the DB record so we can retry later
					continue
				}

				if err := ds.RevokeRoleAssignment(assignment.ID); err != nil {
					log.Printf("[role_scheduler] failed to revoke assignment %d in DB: %v", assignment.ID, err)
					continue
				}

				log.Printf("[role_scheduler] revoked role %s from user %s (assignment %d)", assignment.RoleNativeID, assignment.UserNativeID, assignment.ID)
			}
		}
	}()
}

func mustSnowflakeEnvOrDefault(val string, fallback string) discord.Snowflake {
	if val == "" {
		val = fallback
	}
	s, err := discord.ParseSnowflake(val)
	if err != nil {
		return 0
	}
	return s
}
