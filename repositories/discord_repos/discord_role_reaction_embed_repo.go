package discordrepos

import (
	"doko/gvn-ultimate-bot/models"
	"errors"

	"gorm.io/gorm"
)

type DiscordRoleReactionEmbedRepo interface {
	GetByID(id uint) (*models.DiscordRoleReactionEmbed, error)
	GetByNativeID(nativeId string) (*models.DiscordRoleReactionEmbed, error)
	Create(role *models.DiscordRoleReactionEmbed) (*models.DiscordRoleReactionEmbed, error) // Actually upsert
	Edit(id uint, role *models.DiscordRoleReactionEmbed) (*models.DiscordRoleReactionEmbed, error)
	ListRoleReactionEmbeds() ([]*models.DiscordRoleReactionEmbed, error)
	Delete(id uint) error
}

type discordRoleReactionEmbedRepo struct {
	db *gorm.DB
}

func NewDiscordRoleReactionEmbedRepo(db *gorm.DB) DiscordRoleReactionEmbedRepo {
	return &discordRoleReactionEmbedRepo{
		db: db,
	}
}

func (dr *discordRoleReactionEmbedRepo) Create(payload *models.DiscordRoleReactionEmbed) (*models.DiscordRoleReactionEmbed, error) {
	var r models.DiscordRoleReactionEmbed
	if err := dr.db.Where(&models.DiscordRoleReactionEmbed{NativeMessageId: payload.NativeMessageId}).First(&r).Error; err != nil {
		dr.db.Create(&payload)
		return payload, err
	}

	return &r, nil
}

func (dr *discordRoleReactionEmbedRepo) Edit(id uint, payload *models.DiscordRoleReactionEmbed) (*models.DiscordRoleReactionEmbed, error) {
	var r models.DiscordRoleReactionEmbed

	// Query it first
	if err := dr.db.Where("id = ?", payload.ID).First(&r).Error; err != nil {
		return payload, err
	}

	dr.db.Save(&r)
	return &r, nil
}

func (dr *discordRoleReactionEmbedRepo) Delete(id uint) error {
	return dr.db.Where("id = ?", id).Delete(&models.DiscordRoleReactionEmbed{}).Error
}

func (dr *discordRoleReactionEmbedRepo) GetByID(id uint) (*models.DiscordRoleReactionEmbed, error) {
	var module *models.DiscordRoleReactionEmbed
	if err := dr.db.Where("id = ?", id).First(&module).Error; err == nil {
		return module, err
	}

	return nil, nil
}

func (dr *discordRoleReactionEmbedRepo) GetByNativeID(nativeMessageId string) (*models.DiscordRoleReactionEmbed, error) {
	var embed *models.DiscordRoleReactionEmbed
	err := dr.db.First(&embed, "native_message_id = ?", nativeMessageId).Error

	if err != nil {
		return nil, err
	}

	// If not found, err is nil but no data
	if embed == nil {
		return nil, errors.New("not found")
	}

	return embed, err
}

func (dr *discordRoleReactionEmbedRepo) ListRoleReactionEmbeds() ([]*models.DiscordRoleReactionEmbed, error) {
	var embeds []*models.DiscordRoleReactionEmbed
	if err := dr.db.Find(&embeds).Error; err != nil {
		return embeds, err
	}

	return embeds, nil
}
