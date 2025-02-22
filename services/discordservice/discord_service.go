package discordservice

import (
	"doko/gvn-ultimate-bot/models"
	discordrepos "doko/gvn-ultimate-bot/repositories/discord_repos"
)

type DiscordService interface {
	// ################# For DiscordRole #################
	// Listing the roles as DiscordRole model array
	ListRoles() ([]*models.DiscordRole, error)
	CreateRole(*models.DiscordRole) (*models.DiscordRole, error)
	EditRole(*models.DiscordRole) (*models.DiscordRole, error)

	// This is actually just mark the role as "deleted" (IsDeleted = 1), not actually delete it
	// UnassignRole(id uint) (*models.DiscordRole, error)
}

type discordService struct {
	RoleRepo              discordrepos.DiscordRoleRepo
	RoleReactionEmbedRepo discordrepos.DiscordRoleReactionEmbedRepo
}

func NewDiscordRoleService(repo discordrepos.DiscordRoleRepo, embedRepo discordrepos.DiscordRoleReactionEmbedRepo) DiscordService {
	return &discordService{
		RoleRepo:              repo,
		RoleReactionEmbedRepo: embedRepo,
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

// func (dr *discordService) UnassignRole(id uint) (*models.DiscordRole, error) {
// 	return dr.RoleRepo.UnassignRole(id)
// }

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
	return d.RoleReactionRepo.Upsert(m)
}

func (d *discordRoleReactionEmbedService) GetSingleEmbed(id uint) (*models.DiscordRoleReactionEmbed, error) {
	return d.RoleReactionRepo.GetByID(id)
}
