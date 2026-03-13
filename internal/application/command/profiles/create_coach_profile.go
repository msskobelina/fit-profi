package profiles

import "github.com/msskobelina/fit-profi/internal/domain/model"

type CreateCoachProfileCommand struct {
	UserID       int
	FullName     string
	Category     string
	Info         string
	Achievements []model.CoachAchievement
	Education    []model.CoachEducation
}
