package nutrition

import (
	"time"

	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type CreateEntryCommand struct {
	UserID   int
	Date     time.Time
	MealType string
	Items    []model.DiaryItem
}
