package adminaccessservice

import (
	"doko/gvn-ultimate-bot/models"
	discordrepos "doko/gvn-ultimate-bot/repositories/discord_repos"
	"doko/gvn-ultimate-bot/statics"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

type AccessCheckResult struct {
	Allowed     bool   `json:"allowed"`
	AccessLevel string `json:"access_level"`
}

type AdminAccessService interface {
	ListWhitelistedRoles() ([]*models.AdminWhitelistedRole, error)
	UpsertWhitelistedRole(roleNativeID string, accessLevel string, label string) (*models.AdminWhitelistedRole, error)
	DeleteWhitelistedRole(id uint) error
	CheckAccess(discordUserNativeID string) (*AccessCheckResult, error)
}

type adminAccessService struct {
	repo    discordrepos.AdminWhitelistedRoleRepo
	state   *state.State
	guildID discord.GuildID
}

func NewAdminAccessService(repo discordrepos.AdminWhitelistedRoleRepo, s *state.State, guildID discord.GuildID) AdminAccessService {
	return &adminAccessService{
		repo:    repo,
		state:   s,
		guildID: guildID,
	}
}

func (a *adminAccessService) ListWhitelistedRoles() ([]*models.AdminWhitelistedRole, error) {
	return a.repo.List()
}

func (a *adminAccessService) UpsertWhitelistedRole(roleNativeID string, accessLevel string, label string) (*models.AdminWhitelistedRole, error) {
	if accessLevel == "" {
		accessLevel = statics.Standard
	}
	return a.repo.Upsert(&models.AdminWhitelistedRole{
		RoleNativeID: roleNativeID,
		AccessLevel:  accessLevel,
		Label:        label,
	})
}

func (a *adminAccessService) DeleteWhitelistedRole(id uint) error {
	return a.repo.Delete(id)
}

// CheckAccess looks up the Discord guild member's roles and checks whether any of
// them is whitelisted for admin dashboard access. When multiple whitelisted roles
// match, the highest access level found is returned (order: admin > standard).
func (a *adminAccessService) CheckAccess(discordUserNativeID string) (*AccessCheckResult, error) {
	userID, err := discord.ParseSnowflake(discordUserNativeID)
	if err != nil {
		return nil, fmt.Errorf("invalid discord user id: %w", err)
	}

	member, err := a.state.Member(a.guildID, discord.UserID(userID))
	if err != nil {
		return &AccessCheckResult{Allowed: false}, nil
	}

	if len(member.RoleIDs) == 0 {
		return &AccessCheckResult{Allowed: false}, nil
	}

	memberRoleNativeIDs := make([]string, 0, len(member.RoleIDs))
	for _, roleID := range member.RoleIDs {
		memberRoleNativeIDs = append(memberRoleNativeIDs, roleID.String())
	}

	whitelisted, err := a.repo.GetByRoleNativeIDs(memberRoleNativeIDs)
	if err != nil {
		return nil, err
	}

	if len(whitelisted) == 0 {
		return &AccessCheckResult{Allowed: false}, nil
	}

	accessLevel := statics.Standard
	for _, w := range whitelisted {
		if w.AccessLevel == statics.AdminRole {
			accessLevel = statics.AdminRole
			break
		}
	}

	return &AccessCheckResult{Allowed: true, AccessLevel: accessLevel}, nil
}
