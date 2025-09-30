// FitProfi API
//
//  Schemes: http
//  Host: localhost:8080
//  BasePath: /api/v1
//  Version: 0.1.0
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
//
//  SecurityDefinitions:
//  Bearer:
//    type: apiKey
//    name: Authorization
//    in: header
//
// swagger:meta

package main

import "github.com/msskobelina/fit-profi/api"

func main() {
	api.Run()
}
