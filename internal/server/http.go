package server

import (
	v1App "github.com/ntquang98/go-rkinetics-service/api/app/v1"
	v1File "github.com/ntquang98/go-rkinetics-service/api/file/v1"
	v1 "github.com/ntquang98/go-rkinetics-service/api/helloworld/v1"
	"github.com/ntquang98/go-rkinetics-service/internal/conf"
	"github.com/ntquang98/go-rkinetics-service/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, file *service.FileService, analyticJob *service.AnalyticsJobService, logger log.Logger) *transhttp.Server {
	var opts = []transhttp.ServerOption{
		transhttp.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, transhttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, transhttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, transhttp.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := transhttp.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	v1App.RegisterAnalyticsJobHTTPServer(srv, analyticJob)
	v1Router := srv.Route("/v1")
	v1Router.POST("/file-upload", fileUploadHandler(file, logger))

	h := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", h)
	return srv
}

func fileUploadHandler(fileSrv *service.FileService, logger log.Logger) func(ctx transhttp.Context) error {
	return func(ctx transhttp.Context) error {
		req := ctx.Request()

		fileUrl, err := fileSrv.UploadFileHTTP(ctx, req)

		if err != nil {
			return err
		}

		return ctx.Result(200, &v1File.UploadFileReply{
			FileUrl: fileUrl,
		})
	}
}
