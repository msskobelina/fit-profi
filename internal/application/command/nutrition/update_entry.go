package nutrition

import "github.com/msskobelina/fit-profi/internal/domain/model"

type UpdateEntryCommand struct {
	EntryID  int
	UserID   int
	MealType string
	Items    []model.DiaryItem
}
