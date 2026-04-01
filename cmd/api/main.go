//	@title			FitProfi API
//	@version		0.1.0
//	@description	Fitness platform REST API — training programs, nutrition diary, Google Calendar integration.
//	@host			localhost:8086
//	@BasePath		/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter: Bearer <token>

package main

import (
	_ "github.com/msskobelina/fit-profi/docs"
	"github.com/msskobelina/fit-profi/internal/bootstrap"
)

func main() {
	bootstrap.Run()
}
