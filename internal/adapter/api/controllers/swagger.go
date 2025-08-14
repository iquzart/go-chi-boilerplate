package controllers

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func SwaggerHandler() http.Handler {
	return httpSwagger.Handler(
		httpSwagger.URL("/system/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.PersistAuthorization(true),
	)
}
