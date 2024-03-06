package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"finapp/bootstrap"
)

//	@title		Finapp API
//	@version	0.1

//	@BasePath	/api/v1/

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Bearer

func main() {
	_ = godotenv.Load()
	bootstrap.RootApp.Execute()
}
