//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/ntquang98/go-rkinetics-service/internal/biz"
	"github.com/ntquang98/go-rkinetics-service/internal/conf"
	"github.com/ntquang98/go-rkinetics-service/internal/data"
	"github.com/ntquang98/go-rkinetics-service/internal/server"
	"github.com/ntquang98/go-rkinetics-service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		data.ProviderSet,
		biz.ProviderSet,
		server.ProviderSet,
		service.ProviderSet,
		newApp,
	))
}
