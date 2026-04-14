//go:build wireinject

package main

import (
	"superindo-test/internal/app"

	"github.com/google/wire"
)

func InitializeApp() (*app.App, func(), error) {
	wire.Build(app.ProviderSet)
	return nil, nil, nil
}

