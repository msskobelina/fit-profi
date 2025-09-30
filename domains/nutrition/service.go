package nutrition

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	Create(ctx context.Context, userID int, req *CreateDiaryEntryRequest) (*DiaryEntry, error)
	Update(ctx context.Context, userID, id int, req *UpdateDiaryEntryRequest) (*DiaryEntry, error)
	Get(ctx context.Context, userID, id int, role string) (*DiaryEntry, error)
	List(ctx context.Context, userID int, date string, from, to string) ([]DiaryEntry, error)
	Delete(ctx context.Context, userID, id int, role string) error
	Summary(ctx context.Context, userID int, from, to string) (DiarySummary, error)
}

type service struct{ repo Repository }

func NewService(r Repository) Service { return &service{r} }

func parseDateYYYYMMDD(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("date is required")
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, errors.New("invalid date format (want YYYY-MM-DD)")
	}
	return t, nil
}

func calcTotals(items []CreateDiaryItem) (cals, prot, fat, carb float32) {
	for _, it := range items {
		cals += it.Calories
		prot += it.Protein
		fat += it.Fat
		carb += it.Carbs
	}
	return
}

func (s *service) Create(ctx context.Context, userID int, req *CreateDiaryEntryRequest) (*DiaryEntry, error) {
	if req == nil {
		return nil, errors.New("empty body")
	}
	d, err := parseDateYYYYMMDD(req.Date)
	if err != nil {
		return nil, err
	}
	cals, prot, fat, carb := calcTotals(req.Items)
	entry := &DiaryEntry{
		UserID:        userID,
		Date:          d,
		Meal:          req.Meal,
		Notes:         req.Notes,
		TotalCalories: cals,
		TotalProtein:  prot,
		TotalFat:      fat,
		TotalCarbs:    carb,
	}
	for _, it := range req.Items {
		entry.Items = append(entry.Items, DiaryItem{
			Name:     it.Name,
			Quantity: it.Quantity,
			Unit:     it.Unit,
			Calories: it.Calories,
			Protein:  it.Protein,
			Fat:      it.Fat,
			Carbs:    it.Carbs,
		})
	}
	return s.repo.CreateEntry(ctx, entry)
}

func (s *service) Update(ctx context.Context, userID, id int, req *UpdateDiaryEntryRequest) (*DiaryEntry, error) {
	if req == nil {
		return nil, errors.New("empty body")
	}
	cur, err := s.repo.GetEntryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cur.UserID != userID {
		return nil, errors.New("forbidden")
	}
	d, err := parseDateYYYYMMDD(req.Date)
	if err != nil {
		return nil, err
	}
	cals, prot, fat, carb := calcTotals(req.Items)
	upd := &DiaryEntry{
		ID:            id,
		UserID:        userID,
		Date:          d,
		Meal:          req.Meal,
		Notes:         req.Notes,
		TotalCalories: cals,
		TotalProtein:  prot,
		TotalFat:      fat,
		TotalCarbs:    carb,
	}
	for _, it := range req.Items {
		upd.Items = append(upd.Items, DiaryItem{
			Name:     it.Name,
			Quantity: it.Quantity,
			Unit:     it.Unit,
			Calories: it.Calories,
			Protein:  it.Protein,
			Fat:      it.Fat,
			Carbs:    it.Carbs,
		})
	}
	return s.repo.UpdateEntry(ctx, upd)
}

func (s *service) Get(ctx context.Context, userID, id int, role string) (*DiaryEntry, error) {
	e, err := s.repo.GetEntryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role != "admin" && e.UserID != userID {
		return nil, errors.New("forbidden")
	}
	return e, nil
}

func (s *service) List(ctx context.Context, userID int, date string, from, to string) ([]DiaryEntry, error) {
	if date != "" {
		d, err := parseDateYYYYMMDD(date)
		if err != nil {
			return nil, err
		}
		return s.repo.ListEntriesByDate(ctx, userID, d)
	}
	df, err := parseDateYYYYMMDD(from)
	if err != nil {
		return nil, err
	}
	dt, err := parseDateYYYYMMDD(to)
	if err != nil {
		return nil, err
	}
	return s.repo.ListEntriesByRange(ctx, userID, df, dt)
}

func (s *service) Delete(ctx context.Context, userID, id int, role string) error {
	e, err := s.repo.GetEntryByID(ctx, id)
	if err != nil {
		return err
	}
	if role != "admin" && e.UserID != userID {
		return errors.New("forbidden")
	}
	return s.repo.DeleteEntry(ctx, id)
}

func (s *service) Summary(ctx context.Context, userID int, from, to string) (DiarySummary, error) {
	df, err := parseDateYYYYMMDD(from)
	if err != nil {
		return DiarySummary{}, err
	}
	dt, err := parseDateYYYYMMDD(to)
	if err != nil {
		return DiarySummary{}, err
	}
	return s.repo.SummaryByRange(ctx, userID, df, dt)
}
