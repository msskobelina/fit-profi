package programs

import (
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, coachID int, role string, in *CreateProgramRequest) (*TrainingProgram, error)
	Get(ctx context.Context, requesterID int, role string, id int) (*TrainingProgram, error)
	ListByUser(ctx context.Context, uid int) ([]TrainingProgram, error)
	Delete(ctx context.Context, role string, id int) error
	AddProgress(ctx context.Context, uid int, in *AddProgressRequest) (*ExerciseProgress, error)
}

type service struct{ repo Repository }

func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) Create(ctx context.Context, coachID int, role string, in *CreateProgramRequest) (*TrainingProgram, error) {
	if role != "admin" {
		return nil, errors.New("forbidden")
	}
	if in == nil || in.UserID == 0 || in.Title == "" || in.Weeks <= 0 {
		return nil, errors.New("invalid payload")
	}
	p := &TrainingProgram{
		UserID: in.UserID, CoachUserID: &coachID, Title: in.Title, Notes: in.Notes, Weeks: in.Weeks,
	}
	for _, d := range in.Days {
		day := ProgramDay{WeekNumber: d.WeekNumber, DayNumber: d.DayNumber, Title: d.Title}
		for _, ex := range d.Exercises {
			day.Exercises = append(day.Exercises, ProgramExercise{
				Name: ex.Name, TargetSets: ex.TargetSets, TargetReps: ex.TargetReps, TargetWeight: ex.TargetWeight,
			})
		}
		p.Days = append(p.Days, day)
	}

	return s.repo.CreateProgram(ctx, p)
}

func (s *service) Get(ctx context.Context, requesterID int, role string, id int) (*TrainingProgram, error) {
	p, err := s.repo.GetProgram(ctx, id)
	if err != nil {
		return nil, err
	}
	if role != "admin" && p.UserID != requesterID {
		return nil, errors.New("forbidden")
	}

	return p, nil
}

func (s *service) ListByUser(ctx context.Context, uid int) ([]TrainingProgram, error) {
	return s.repo.ListProgramsByUser(ctx, uid)
}

func (s *service) Delete(ctx context.Context, role string, id int) error {
	if role != "admin" {
		return errors.New("forbidden")
	}

	return s.repo.DeleteProgram(ctx, id)
}

func (s *service) AddProgress(ctx context.Context, uid int, in *AddProgressRequest) (*ExerciseProgress, error) {
	if in == nil || in.ProgramExerciseID == 0 || in.WeekNumber <= 0 || in.Reps == "" {
		return nil, errors.New("invalid payload")
	}
	p, err := s.repo.GetProgramByExerciseID(ctx, in.ProgramExerciseID)
	if err != nil {
		return nil, err
	}
	if p.UserID != uid {
		return nil, errors.New("forbidden")
	}
	rec := &ExerciseProgress{
		ProgramExerciseID: in.ProgramExerciseID,
		WeekNumber:        in.WeekNumber,
		SetNumber:         in.SetNumber,
		Weight:            in.Weight,
		Reps:              in.Reps,
	}

	return s.repo.AddProgress(ctx, rec)
}
