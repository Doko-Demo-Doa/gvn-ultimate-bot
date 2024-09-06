package discordservice

import (
	"doko/gvn-ultimate-bot/models"
	discordrepos "doko/gvn-ultimate-bot/repositories/discord_repos"
)

// ################# For DiscordRole #################

type DiscordService interface {
	// Listing the roles as DiscordRole model array
	ListRoles() ([]*models.DiscordRole, error)
	CreateRole(*models.DiscordRole) (*models.DiscordRole, error)
	EditRole(*models.DiscordRole) (*models.DiscordRole, error)

	// This is actually just mark the role as "deleted" (IsDeleted = 1), not actually delete it
	RemoveRole(email uint) (*models.DiscordRole, error)
}

type discordService struct {
	Repo discordrepos.DiscordRoleRepo
}

func NewDiscordRoleService(repo discordrepos.DiscordRoleRepo) DiscordService {
	return &discordService{
		Repo: repo,
	}
}

func (dr *discordService) CreateRole(r *models.DiscordRole) (*models.DiscordRole, error) {
	return dr.Repo.CreateRole(r)
}

func (dr *discordService) EditRole(r *models.DiscordRole) (*models.DiscordRole, error) {
	return dr.Repo.EditRole(r)
}

func (dr *discordService) ListRoles() ([]*models.DiscordRole, error) {
	return dr.Repo.ListRoles()
}

func (*discordService) RemoveRole(email uint) (*models.DiscordRole, error) {
	panic("unimplemented")
}

// ################# For DiscordRoleReactionEmbed #################

type DiscordRoleReactionEmbedService interface {
	ListEmbeds() ([]*models.DiscordRoleReactionEmbed, error)
}

type discordRoleReactionEmbedService struct {
	Repo discordrepos.DiscordRoleReactionEmbedRepo
}

func NewDiscordRoleReactionEmbedService(repo discordrepos.DiscordRoleReactionEmbedRepo) DiscordRoleReactionEmbedService {
	return &discordRoleReactionEmbedService{
		Repo: repo,
	}
}

// ListEmbeds implements DiscordRoleReactionEmbedService.
func (d *discordRoleReactionEmbedService) ListEmbeds() ([]*models.DiscordRoleReactionEmbed, error) {
	return d.Repo.ListRoleReactionEmbeds()
}

func (dr *discordService) UpsertEmbed(r *models.DiscordRole) (*models.DiscordRole, error) {
	return dr.Repo.CreateRole(r)
}
