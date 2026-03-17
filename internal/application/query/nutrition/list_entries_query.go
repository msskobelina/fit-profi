package nutrition

import "time"

type ListEntriesQuery struct {
	UserID int
	Date   time.Time
}
