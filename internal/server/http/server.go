package http

import (
	"app-parity/internal/service"
	"github.com/bilibili/kratos/pkg/net/http/blademaster/middleware/auth"
	"net/http"

	pb "app-parity/api"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

var (
	svc   *service.Service
	_cors = []string{
		"*",
		"http://192.168.1.254:8080",
		"http://localhost:8080",
		"http://bijia.ybr1.cn",
	}
	_csrf = []string{
		"*",
		"http://192.168.1.254:8080",
		"http://localhost:8080",
		"http://bijia.ybr1.cn",
	}
)

// New new a bm server.
func New(s *service.Service) (e *bm.Engine) {
	var (
		hc struct {
			Server *bm.ServerConfig
		}
		fs = []string{"/user/v1/login"}
	)
	if err := paladin.Get("http.toml").UnmarshalTOML(&hc); err != nil {
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	svc = s
	e = bm.DefaultServer(hc.Server)
	e.Ping(ping)
	e.Use(bm.CORS(_cors))
	e.Use(bm.CSRF(_csrf, []string{}))
	e.Use(auth.New(&auth.Config{JwtSecret: s.JWT.Secret, Filters: fs}))
	//e.Inject("^[login]", auth.New(&auth.Config{JwtSecret: s.JwtSecret, Filters: fs}))
	pb.RegisterUserServiceBMServer(e, svc)
	pb.RegisterDrugServiceBMServer(e, svc)
	if err := e.Start(); err != nil {
		panic(err)
	}
	return
}

func ping(ctx *bm.Context) {
	if err := svc.Ping(ctx); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
