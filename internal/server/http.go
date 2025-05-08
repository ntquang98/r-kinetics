package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	v1 "github.com/ntquang98/go-rkinetics-service/api/helloworld/v1"
	"github.com/ntquang98/go-rkinetics-service/internal/conf"
	"github.com/ntquang98/go-rkinetics-service/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, file *service.FileService, logger log.Logger) *transhttp.Server {
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

	// ============= ROUTER ====================
	router := mux.NewRouter()
	router.HandleFunc("/upload", UploadFileHandler(file, logger)).Methods("POST")
	srv.HandlePrefix("/upload", router)

	h := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", h)
	return srv
}

// UploadFileHandler handles the file upload endpoint
func UploadFileHandler(file *service.FileService, logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		fileURL, err := file.UploadFile(r.Context(), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"fileURL": fileURL})

		return
	}
}
