package discordrepos

import (
	"doko/gvn-ultimate-bot/models"

	"gorm.io/gorm"
)

type DiscordRoleReactionEmbedRepo interface {
	GetByID(id uint) (*models.DiscordRoleReactionEmbed, error)
	GetByNativeID(nativeId string) (*models.DiscordRoleReactionEmbed, error)
	Create(role *models.DiscordRoleReactionEmbed) (*models.DiscordRoleReactionEmbed, error) // Actually upsert
	Update(nativeMessageId string, role *models.DiscordRoleReactionEmbed) (*models.DiscordRoleReactionEmbed, error)
	Upsert(role *models.DiscordRoleReactionEmbed) (*models.DiscordRoleReactionEmbed, error)
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
	dr.db.FirstOrInit(&r, payload)
	if err := dr.db.Where(&models.DiscordRoleReactionEmbed{NativeMessageId: payload.NativeMessageId}).FirstOrCreate(&r).Error; err != nil {
		return &r, err
	}

	return &r, nil
}

func (dr *discordRoleReactionEmbedRepo) Update(nativeMessageId string, payload *models.DiscordRoleReactionEmbed) (*models.DiscordRoleReactionEmbed, error) {
	var r models.DiscordRoleReactionEmbed
	result := dr.db.Where("native_message_id = ?", nativeMessageId).First(&r)

	if result.Error != nil {
		return nil, result.Error
	}

	r.Name = payload.Name
	r.NativeMessageId = payload.NativeMessageId
	r.Tags = payload.Tags
	r.Version = payload.Version
	r.Payload = payload.Payload

	dr.db.Save(&r)
	return &r, nil
}

func (dr *discordRoleReactionEmbedRepo) Upsert(role *models.DiscordRoleReactionEmbed) (*models.DiscordRoleReactionEmbed, error) {
	panic("unimplemented")
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
	err := dr.db.Where("native_message_id = ?", nativeMessageId).First(&embed).Error
	if err != nil {
		return nil, err
	}

	return embed, nil
}

func (dr *discordRoleReactionEmbedRepo) ListRoleReactionEmbeds() ([]*models.DiscordRoleReactionEmbed, error) {
	var embeds []*models.DiscordRoleReactionEmbed
	if err := dr.db.Find(&embeds).Error; err != nil {
		return embeds, err
	}

	return embeds, nil
}
