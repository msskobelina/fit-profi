package programs

import (
	"context"

	"github.com/msskobelina/fit-profi/internal/domain/model"
	"github.com/msskobelina/fit-profi/internal/domain/repository"
)

type TrackProgressHandler interface {
	TrackProgress(context.Context, TrackProgressCommand) (*model.ExerciseProgress, error)
}

type trackProgressService struct {
	repo repository.ProgramsRepository
}

func NewTrackProgressService(repo repository.ProgramsRepository) TrackProgressHandler {
	return &trackProgressService{repo: repo}
}

func (s *trackProgressService) TrackProgress(ctx context.Context, cmd TrackProgressCommand) (*model.ExerciseProgress, error) {
	return s.repo.TrackProgress(ctx, model.ExerciseProgress{
		UserID:     cmd.UserID,
		ExerciseID: cmd.ExerciseID,
		Sets:       cmd.Sets,
		Reps:       cmd.Reps,
		WeightKg:   cmd.WeightKg,
		Notes:      cmd.Notes,
	})
}
