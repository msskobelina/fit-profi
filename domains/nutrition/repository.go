package nutrition

import (
	"context"
	"time"

	"github.com/msskobelina/fit-profi/pkg/mysql"
	"gorm.io/gorm"
)

type Repository interface {
	CreateEntry(ctx context.Context, e *DiaryEntry) (*DiaryEntry, error)
	GetEntryByID(ctx context.Context, id int) (*DiaryEntry, error)
	ListEntriesByDate(ctx context.Context, userID int, date time.Time) ([]DiaryEntry, error)
	ListEntriesByRange(ctx context.Context, userID int, from, to time.Time) ([]DiaryEntry, error)
	UpdateEntry(ctx context.Context, e *DiaryEntry) (*DiaryEntry, error)
	DeleteEntry(ctx context.Context, id int) error
	SummaryByRange(ctx context.Context, userID int, from, to time.Time) (DiarySummary, error)
}

type gormRepo struct{ *mysql.MySQL }

func NewRepository(sql *mysql.MySQL) Repository { return &gormRepo{sql} }

func (r *gormRepo) CreateEntry(ctx context.Context, e *DiaryEntry) (*DiaryEntry, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Items").Create(e).Error; err != nil {
			return err
		}
		for i := range e.Items {
			e.Items[i].ID = 0
			e.Items[i].DiaryEntryID = e.ID
		}
		if len(e.Items) > 0 {
			if err := tx.Create(&e.Items).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetEntryByID(ctx, e.ID)
}

func (r *gormRepo) GetEntryByID(ctx context.Context, id int) (*DiaryEntry, error) {
	var e DiaryEntry
	if err := r.DB.WithContext(ctx).
		Preload("Items").
		First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *gormRepo) ListEntriesByDate(ctx context.Context, userID int, date time.Time) ([]DiaryEntry, error) {
	var list []DiaryEntry
	if err := r.DB.WithContext(ctx).
		Preload("Items").
		Where("user_id = ? AND date = ?", userID, date.Format("2006-01-02")).
		Order("meal asc, id asc").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *gormRepo) ListEntriesByRange(ctx context.Context, userID int, from, to time.Time) ([]DiaryEntry, error) {
	var list []DiaryEntry
	if err := r.DB.WithContext(ctx).
		Preload("Items").
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, from.Format("2006-01-02"), to.Format("2006-01-02")).
		Order("date asc, meal asc, id asc").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *gormRepo) UpdateEntry(ctx context.Context, e *DiaryEntry) (*DiaryEntry, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&DiaryEntry{}).
			Where("id = ?", e.ID).
			Updates(map[string]any{
				"date":           e.Date,
				"meal":           e.Meal,
				"notes":          e.Notes,
				"total_calories": e.TotalCalories,
				"total_protein":  e.TotalProtein,
				"total_fat":      e.TotalFat,
				"total_carbs":    e.TotalCarbs,
			}).Error; err != nil {
			return err
		}
		if err := tx.Where("diary_entry_id = ?", e.ID).Delete(&DiaryItem{}).Error; err != nil {
			return err
		}
		for i := range e.Items {
			e.Items[i].ID = 0
			e.Items[i].DiaryEntryID = e.ID
		}
		if len(e.Items) > 0 {
			if err := tx.Create(&e.Items).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetEntryByID(ctx, e.ID)
}

func (r *gormRepo) DeleteEntry(ctx context.Context, id int) error {
	return r.DB.WithContext(ctx).Delete(&DiaryEntry{}, id).Error
}

func (r *gormRepo) SummaryByRange(ctx context.Context, userID int, from, to time.Time) (DiarySummary, error) {
	type agg struct {
		Cals float32
		Prot float32
		Fat  float32
		Carb float32
	}
	var a agg
	err := r.DB.WithContext(ctx).
		Model(&DiaryEntry{}).
		Select("COALESCE(SUM(total_calories),0) AS cals, COALESCE(SUM(total_protein),0) AS prot, COALESCE(SUM(total_fat),0) AS fat, COALESCE(SUM(total_carbs),0) AS carb").
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, from.Format("2006-01-02"), to.Format("2006-01-02")).
		Scan(&a).Error
	if err != nil {
		return DiarySummary{}, err
	}
	return DiarySummary{
		From:          from.Format("2006-01-02"),
		To:            to.Format("2006-01-02"),
		TotalCalories: a.Cals,
		TotalProtein:  a.Prot,
		TotalFat:      a.Fat,
		TotalCarbs:    a.Carb,
	}, nil
}
