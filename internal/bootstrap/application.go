package bootstrap

import (
	"context"

	cmdAuthorize "github.com/msskobelina/fit-profi/internal/application/command/authorize"
	cmdCalendar "github.com/msskobelina/fit-profi/internal/application/command/calendar"
	cmdIntegrations "github.com/msskobelina/fit-profi/internal/application/command/integrations"
	cmdNutrition "github.com/msskobelina/fit-profi/internal/application/command/nutrition"
	cmdProfiles "github.com/msskobelina/fit-profi/internal/application/command/profiles"
	cmdPrograms "github.com/msskobelina/fit-profi/internal/application/command/programs"
	qryAuthorize "github.com/msskobelina/fit-profi/internal/application/query/authorize"
	qryCalendar "github.com/msskobelina/fit-profi/internal/application/query/calendar"
	qryNutrition "github.com/msskobelina/fit-profi/internal/application/query/nutrition"
	qryProfiles "github.com/msskobelina/fit-profi/internal/application/query/profiles"
	qryPrograms "github.com/msskobelina/fit-profi/internal/application/query/programs"
	"github.com/msskobelina/fit-profi/internal/domain/model"
)

type application struct {
	// authorize
	registerUser   cmdAuthorize.RegisterUserHandler
	loginUser      cmdAuthorize.LoginUserHandler
	logoutUser     cmdAuthorize.LogoutUserHandler
	sendResetEmail cmdAuthorize.SendResetEmailHandler
	resetPassword  cmdAuthorize.ResetPasswordHandler
	verifyToken    qryAuthorize.VerifyTokenHandler
	// profiles
	createUserProfile  cmdProfiles.CreateUserProfileHandler
	updateUserProfile  cmdProfiles.UpdateUserProfileHandler
	getUserProfile     qryProfiles.GetUserProfileHandler
	createCoachProfile cmdProfiles.CreateCoachProfileHandler
	updateCoachProfile cmdProfiles.UpdateCoachProfileHandler
	getCoachProfile    qryProfiles.GetCoachProfileHandler
	// programs
	createProgram cmdPrograms.CreateProgramHandler
	deleteProgram cmdPrograms.DeleteProgramHandler
	trackProgress cmdPrograms.TrackProgressHandler
	getProgram    qryPrograms.GetProgramHandler
	listPrograms  qryPrograms.ListProgramsHandler
	// nutrition
	createEntry cmdNutrition.CreateEntryHandler
	updateEntry cmdNutrition.UpdateEntryHandler
	deleteEntry cmdNutrition.DeleteEntryHandler
	listEntries qryNutrition.ListEntriesHandler
	getEntry    qryNutrition.GetEntryHandler
	// integrations
	connectGoogle    cmdIntegrations.ConnectGoogleHandler
	exchangeCallback cmdIntegrations.ExchangeCallbackHandler
	// calendar
	listCalendars qryCalendar.ListCalendarsHandler
	createEvent   cmdCalendar.CreateEventHandler
}

// authorize

func (a *application) Register(ctx context.Context, cmd cmdAuthorize.RegisterUserCommand) (*cmdAuthorize.RegisterUserResult, error) {
	return a.registerUser.Register(ctx, cmd)
}

func (a *application) Login(ctx context.Context, cmd cmdAuthorize.LoginUserCommand) (*cmdAuthorize.LoginUserResult, error) {
	return a.loginUser.Login(ctx, cmd)
}

func (a *application) Logout(ctx context.Context, cmd cmdAuthorize.LogoutUserCommand) error {
	return a.logoutUser.Logout(ctx, cmd)
}

func (a *application) SendResetEmail(ctx context.Context, cmd cmdAuthorize.SendResetEmailCommand) error {
	return a.sendResetEmail.SendResetEmail(ctx, cmd)
}

func (a *application) ResetPassword(ctx context.Context, cmd cmdAuthorize.ResetPasswordCommand) error {
	return a.resetPassword.ResetPassword(ctx, cmd)
}

func (a *application) VerifyToken(ctx context.Context, q qryAuthorize.VerifyTokenQuery) (*qryAuthorize.VerifyTokenResult, error) {
	return a.verifyToken.VerifyToken(ctx, q)
}

// profiles

func (a *application) CreateUserProfile(ctx context.Context, cmd cmdProfiles.CreateUserProfileCommand) (*model.UserProfile, error) {
	return a.createUserProfile.CreateUserProfile(ctx, cmd)
}

func (a *application) UpdateUserProfile(ctx context.Context, cmd cmdProfiles.UpdateUserProfileCommand) (*model.UserProfile, error) {
	return a.updateUserProfile.UpdateUserProfile(ctx, cmd)
}

func (a *application) GetUserProfile(ctx context.Context, q qryProfiles.GetUserProfileQuery) (*model.UserProfile, error) {
	return a.getUserProfile.GetUserProfile(ctx, q)
}

func (a *application) CreateCoachProfile(ctx context.Context, cmd cmdProfiles.CreateCoachProfileCommand) (*model.CoachProfile, error) {
	return a.createCoachProfile.CreateCoachProfile(ctx, cmd)
}

func (a *application) UpdateCoachProfile(ctx context.Context, cmd cmdProfiles.UpdateCoachProfileCommand) (*model.CoachProfile, error) {
	return a.updateCoachProfile.UpdateCoachProfile(ctx, cmd)
}

func (a *application) GetCoachProfile(ctx context.Context, q qryProfiles.GetCoachProfileQuery) (*model.CoachProfile, error) {
	return a.getCoachProfile.GetCoachProfile(ctx, q)
}

// programs

func (a *application) CreateProgram(ctx context.Context, cmd cmdPrograms.CreateProgramCommand) (*model.TrainingProgram, error) {
	return a.createProgram.CreateProgram(ctx, cmd)
}

func (a *application) DeleteProgram(ctx context.Context, cmd cmdPrograms.DeleteProgramCommand) error {
	return a.deleteProgram.DeleteProgram(ctx, cmd)
}

func (a *application) TrackProgress(ctx context.Context, cmd cmdPrograms.TrackProgressCommand) (*model.ExerciseProgress, error) {
	return a.trackProgress.TrackProgress(ctx, cmd)
}

func (a *application) GetProgram(ctx context.Context, q qryPrograms.GetProgramQuery) (*model.TrainingProgram, error) {
	return a.getProgram.GetProgram(ctx, q)
}

func (a *application) ListPrograms(ctx context.Context, q qryPrograms.ListProgramsQuery) ([]model.TrainingProgram, error) {
	return a.listPrograms.ListPrograms(ctx, q)
}

// nutrition

func (a *application) CreateEntry(ctx context.Context, cmd cmdNutrition.CreateEntryCommand) (*model.DiaryEntry, error) {
	return a.createEntry.CreateEntry(ctx, cmd)
}

func (a *application) UpdateEntry(ctx context.Context, cmd cmdNutrition.UpdateEntryCommand) (*model.DiaryEntry, error) {
	return a.updateEntry.UpdateEntry(ctx, cmd)
}

func (a *application) DeleteEntry(ctx context.Context, cmd cmdNutrition.DeleteEntryCommand) error {
	return a.deleteEntry.DeleteEntry(ctx, cmd)
}

func (a *application) ListEntries(ctx context.Context, q qryNutrition.ListEntriesQuery) ([]model.DiaryEntry, error) {
	return a.listEntries.ListEntries(ctx, q)
}

func (a *application) GetEntry(ctx context.Context, q qryNutrition.GetEntryQuery) (*model.DiaryEntry, error) {
	return a.getEntry.GetEntry(ctx, q)
}

// integrations

func (a *application) ConnectGoogle(ctx context.Context, cmd cmdIntegrations.ConnectGoogleCommand) (*cmdIntegrations.ConnectGoogleResult, error) {
	return a.connectGoogle.ConnectGoogle(ctx, cmd)
}

func (a *application) ExchangeCallback(ctx context.Context, cmd cmdIntegrations.ExchangeCallbackCommand) error {
	return a.exchangeCallback.ExchangeCallback(ctx, cmd)
}

// calendar

func (a *application) ListCalendars(ctx context.Context, q qryCalendar.ListCalendarsQuery) ([]qryCalendar.CalendarInfo, error) {
	return a.listCalendars.ListCalendars(ctx, q)
}

func (a *application) CreateEvent(ctx context.Context, cmd cmdCalendar.CreateEventCommand) (*cmdCalendar.CreateEventResult, error) {
	return a.createEvent.CreateEvent(ctx, cmd)
}
