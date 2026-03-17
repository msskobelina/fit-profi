package profiles

import "github.com/msskobelina/fit-profi/internal/domain/model"

type UpdateCoachProfileCommand struct {
	UserID       int
	FullName     string
	Category     string
	Info         string
	Achievements []model.CoachAchievement
	Education    []model.CoachEducation
}
