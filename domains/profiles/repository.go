package profiles

import (
	"context"
	"github.com/msskobelina/fit-profi/pkg/mysql"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUserProfile(ctx context.Context, p *UserProfile) (*UserProfile, error)
	UpdateUserProfile(ctx context.Context, userID int, patch *UpdateUserProfileRequest) (*UserProfile, error)
	GetUserProfileByUserID(ctx context.Context, userID int) (*UserProfile, error)
	DeleteUserProfile(ctx context.Context, userID int) error

	CreateCoachProfile(ctx context.Context, p *CoachProfile) (*CoachProfile, error)
	UpdateCoachProfile(ctx context.Context, userID int, patch *UpdateCoachProfileRequest) (*CoachProfile, error)
	GetCoachProfileByUserID(ctx context.Context, userID int) (*CoachProfile, error)
	DeleteCoachProfile(ctx context.Context, userID int) error
}

type gormRepo struct{ *mysql.MySQL }

func NewRepository(sql *mysql.MySQL) Repository {
	return &gormRepo{
		sql,
	}
}

func (r *gormRepo) CreateUserProfile(ctx context.Context, p *UserProfile) (*UserProfile, error) {
	if err := r.DB.WithContext(ctx).Create(p).Error; err != nil {
		return nil, err
	}

	return r.GetUserProfileByUserID(ctx, p.UserID)
}

func (r *gormRepo) UpdateUserProfile(ctx context.Context, userID int, patch *UpdateUserProfileRequest) (*UserProfile, error) {
	updates := map[string]any{}
	if patch.FullName != nil {
		updates["full_name"] = *patch.FullName
	}
	if patch.Age != nil {
		updates["age"] = *patch.Age
	}
	if patch.WeightKg != nil {
		updates["weight_kg"] = *patch.WeightKg
	}
	if patch.Goal != nil {
		updates["goal"] = *patch.Goal
	}
	if patch.Description != nil {
		updates["description"] = *patch.Description
	}

	tx := r.DB.WithContext(ctx).Model(&UserProfile{}).Where("user_id = ?", userID).Updates(updates)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return r.GetUserProfileByUserID(ctx, userID)
}

func (r *gormRepo) GetUserProfileByUserID(ctx context.Context, userID int) (*UserProfile, error) {
	var profile UserProfile
	if err := r.DB.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *gormRepo) DeleteUserProfile(ctx context.Context, userID int) error {
	return r.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&UserProfile{}).Error
}

func (r *gormRepo) CreateCoachProfile(ctx context.Context, p *CoachProfile) (*CoachProfile, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Achievements", "Education").Create(p).Error; err != nil {
			return err
		}

		if len(p.Achievements) > 0 {
			for i := range p.Achievements {
				p.Achievements[i].ID = 0
				p.Achievements[i].CoachProfileID = p.ID
			}
			if err := tx.Create(&p.Achievements).Error; err != nil {
				return err
			}
		}

		if len(p.Education) > 0 {
			for i := range p.Education {
				p.Education[i].ID = 0
				p.Education[i].CoachProfileID = p.ID
			}
			if err := tx.Create(&p.Education).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetCoachProfileByUserID(ctx, p.UserID)
}

func (r *gormRepo) UpdateCoachProfile(ctx context.Context, userID int, patch *UpdateCoachProfileRequest) (*CoachProfile, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		upd := map[string]any{}
		if patch.FullName != nil {
			upd["full_name"] = *patch.FullName
		}
		if patch.Category != nil {
			upd["category"] = *patch.Category
		}
		if patch.Info != nil {
			upd["info"] = *patch.Info
		}
		if len(upd) > 0 {
			res := tx.Model(&CoachProfile{}).Where("user_id = ?", userID).Updates(upd)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}

		var cp CoachProfile
		if err := tx.Where("user_id = ?", userID).First(&cp).Error; err != nil {
			return err
		}

		if patch.Achievements != nil {
			if err := tx.Where("coach_profile_id = ?", cp.ID).Delete(&CoachAchievement{}).Error; err != nil {
				return err
			}
			ac := *patch.Achievements
			for i := range ac {
				ac[i].CoachProfileID = cp.ID
			}
			if len(ac) > 0 {
				rows := make([]CoachAchievement, 0, len(ac))
				for _, a := range ac {
					rows = append(rows, CoachAchievement{
						CoachProfileID: cp.ID,
						StartPeriod:    a.StartPeriod,
						EndPeriod:      a.EndPeriod,
						Title:          a.Title,
						CertificateURL: a.CertificateURL,
					})
				}
				if err := tx.Create(&rows).Error; err != nil {
					return err
				}
			}
		}
		if patch.Education != nil {
			if err := tx.Where("coach_profile_id = ?", cp.ID).Delete(&CoachEducation{}).Error; err != nil {
				return err
			}
			ed := *patch.Education
			if len(ed) > 0 {
				rows := make([]CoachEducation, 0, len(ed))
				for _, e := range ed {
					rows = append(rows, CoachEducation{
						CoachProfileID: cp.ID,
						StartPeriod:    e.StartPeriod,
						EndPeriod:      e.EndPeriod,
						Place:          e.Place,
						Description:    e.Description,
					})
				}
				if err := tx.Create(&rows).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetCoachProfileByUserID(ctx, userID)
}

func (r *gormRepo) GetCoachProfileByUserID(ctx context.Context, userID int) (*CoachProfile, error) {
	var p CoachProfile
	if err := r.DB.WithContext(ctx).
		Preload("Achievements").
		Preload("Education").
		Where("user_id = ?", userID).
		First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *gormRepo) DeleteCoachProfile(ctx context.Context, userID int) error {
	return r.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&CoachProfile{}).Error
}
