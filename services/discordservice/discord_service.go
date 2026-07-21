package discordservice

import (
	"doko/gvn-ultimate-bot/models"
	discordrepos "doko/gvn-ultimate-bot/repositories/discord_repos"
	"time"
)

type DiscordService interface {
	// ################# For DiscordRole #################
	ListRoles() ([]*models.DiscordRole, error)
	CreateRole(*models.DiscordRole) (*models.DiscordRole, error)
	EditRole(*models.DiscordRole) (*models.DiscordRole, error)
	UnassignRole(nativeUserId string, roleId uint) (*models.DiscordRole, error)

	// ################# For DiscordUserRole (timed assignments) #################
	AssignRoleToUser(userNativeID string, roleNativeID string, duration time.Duration) (*models.DiscordUserRole, error)
	GetExpiredRoleAssignments() ([]*models.DiscordUserRole, error)
	RevokeRoleAssignment(assignmentID uint) error
	GetActiveAssignmentsForUser(nativeUserID string) ([]*models.DiscordUserRole, error)
}

type discordService struct {
	RoleRepo              discordrepos.DiscordRoleRepo
	RoleReactionEmbedRepo discordrepos.DiscordRoleReactionEmbedRepo
	UserRoleRepo          discordrepos.DiscordUserRoleRepo
}

func NewDiscordRoleService(
	repo discordrepos.DiscordRoleRepo,
	embedRepo discordrepos.DiscordRoleReactionEmbedRepo,
	userRoleRepo discordrepos.DiscordUserRoleRepo,
) DiscordService {
	return &discordService{
		RoleRepo:              repo,
		RoleReactionEmbedRepo: embedRepo,
		UserRoleRepo:          userRoleRepo,
	}
}

func (dr *discordService) CreateRole(r *models.DiscordRole) (*models.DiscordRole, error) {
	return dr.RoleRepo.CreateRole(r)
}

func (dr *discordService) EditRole(r *models.DiscordRole) (*models.DiscordRole, error) {
	return dr.RoleRepo.EditRole(r)
}

func (dr *discordService) ListRoles() ([]*models.DiscordRole, error) {
	return dr.RoleRepo.ListRoles()
}

func (dr *discordService) UnassignRole(nativeUserId string, roleId uint) (*models.DiscordRole, error) {
	panic("unimplemented")
}

func (dr *discordService) AssignRoleToUser(userNativeID string, roleNativeID string, duration time.Duration) (*models.DiscordUserRole, error) {
	now := time.Now()
	expiration := now.Add(duration)

	assignment := &models.DiscordUserRole{
		UserNativeID:   userNativeID,
		RoleNativeID:   roleNativeID,
		GrantedDate:    now,
		ExpirationDate: expiration,
	}

	return dr.UserRoleRepo.CreateAssignment(assignment)
}

func (dr *discordService) GetExpiredRoleAssignments() ([]*models.DiscordUserRole, error) {
	return dr.UserRoleRepo.GetExpiredAssignments()
}

func (dr *discordService) RevokeRoleAssignment(assignmentID uint) error {
	return dr.UserRoleRepo.RevokeAssignment(assignmentID)
}

func (dr *discordService) GetActiveAssignmentsForUser(nativeUserID string) ([]*models.DiscordUserRole, error) {
	return dr.UserRoleRepo.GetActiveAssignmentsByUser(nativeUserID)
}

// ################# For DiscordRoleReactionEmbed #################

type DiscordRoleReactionEmbedService interface {
	ListEmbeds() ([]*models.DiscordRoleReactionEmbed, error)
	UpsertEmbed(*models.DiscordRoleReactionEmbed) (*models.DiscordRoleReactionEmbed, error)
	GetSingleEmbed(id uint) (*models.DiscordRoleReactionEmbed, error)
}

type discordRoleReactionEmbedService struct {
	RoleReactionRepo discordrepos.DiscordRoleReactionEmbedRepo
}

func NewDiscordRoleReactionEmbedService(repo discordrepos.DiscordRoleReactionEmbedRepo) DiscordRoleReactionEmbedService {
	return &discordRoleReactionEmbedService{
		RoleReactionRepo: repo,
	}
}

func (d *discordRoleReactionEmbedService) ListEmbeds() ([]*models.DiscordRoleReactionEmbed, error) {
	return d.RoleReactionRepo.ListRoleReactionEmbeds()
}

func (d *discordRoleReactionEmbedService) UpsertEmbed(m *models.DiscordRoleReactionEmbed) (*models.DiscordRoleReactionEmbed, error) {
	data, err := d.RoleReactionRepo.GetByNativeID(m.NativeMessageId)
	if err != nil || data == nil {
		return d.RoleReactionRepo.Create(m)
	}
	return d.RoleReactionRepo.Update(m.NativeMessageId, m)
}

func (d *discordRoleReactionEmbedService) GetSingleEmbed(id uint) (*models.DiscordRoleReactionEmbed, error) {
	return d.RoleReactionRepo.GetByID(id)
}
