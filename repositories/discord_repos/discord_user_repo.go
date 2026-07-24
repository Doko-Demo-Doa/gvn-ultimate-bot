package discordrepos

import (
	"doko/gvn-ultimate-bot/models"

	"gorm.io/gorm"
)

type DiscordUserRepo interface {
	GetByID(id uint) (*models.DiscordUser, error)
	GetByNativeID(nativeId string) (*models.DiscordUser, error)
	Upsert(user *models.DiscordUser) (*models.DiscordUser, error)
	ListAll() ([]*models.DiscordUser, error)
}

type discordUserRepo struct {
	db *gorm.DB
}

func NewDiscordUserRepo(db *gorm.DB) DiscordUserRepo {
	return &discordUserRepo{db: db}
}

func (r *discordUserRepo) GetByID(id uint) (*models.DiscordUser, error) {
	var user models.DiscordUser
	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *discordUserRepo) GetByNativeID(nativeId string) (*models.DiscordUser, error) {
	var user models.DiscordUser
	if err := r.db.Where("native_id = ?", nativeId).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *discordUserRepo) Upsert(user *models.DiscordUser) (*models.DiscordUser, error) {
	existing, err := r.GetByNativeID(user.NativeId)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		if err := r.db.Create(user).Error; err != nil {
			return nil, err
		}
		return user, nil
	}

	existing.Discriminator = user.Discriminator
	existing.Avatar = user.Avatar
	existing.PremiumType = user.PremiumType
	existing.Nickname = user.Nickname
	if err := r.db.Save(existing).Error; err != nil {
		return nil, err
	}
	return existing, nil
}

func (r *discordUserRepo) ListAll() ([]*models.DiscordUser, error) {
	var users []*models.DiscordUser
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
