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
	// DeleteNotIn removes all users whose NativeId is not in nativeIds.
	// If nativeIds is empty, all users are removed.
	DeleteNotIn(nativeIds []string) (int64, error)
	DeleteByNativeID(nativeId string) error
	// Search matches against native_id (exact) or username/nickname (substring,
	// case-insensitive). An empty query returns the first `limit` users.
	Search(query string, limit int) ([]*models.DiscordUser, error)
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
	existing.Username = user.Username
	existing.Nickname = user.Nickname
	if err := r.db.Save(existing).Error; err != nil {
		return nil, err
	}
	return existing, nil
}

func (r *discordUserRepo) Search(query string, limit int) ([]*models.DiscordUser, error) {
	var users []*models.DiscordUser
	q := r.db.Model(&models.DiscordUser{})
	if query != "" {
		q = q.Where("native_id = ? OR username ILIKE ? OR nickname ILIKE ?", query, "%"+query+"%", "%"+query+"%")
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *discordUserRepo) ListAll() ([]*models.DiscordUser, error) {
	var users []*models.DiscordUser
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *discordUserRepo) DeleteByNativeID(nativeId string) error {
	return r.db.Where("native_id = ?", nativeId).Delete(&models.DiscordUser{}).Error
}

func (r *discordUserRepo) DeleteNotIn(nativeIds []string) (int64, error) {
	q := r.db
	if len(nativeIds) > 0 {
		q = q.Where("native_id NOT IN ?", nativeIds)
	}
	result := q.Delete(&models.DiscordUser{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
