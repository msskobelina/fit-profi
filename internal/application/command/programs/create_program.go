package programs

import "github.com/msskobelina/fit-profi/internal/domain/model"

type CreateProgramCommand struct {
	UserID      int
	Title       string
	Description string
	Days        []model.ProgramDay
}
