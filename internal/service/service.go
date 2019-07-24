package service

import (
	"app-parity/internal/dao"
	"context"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao *dao.Dao
	JWT struct{ Secret string }
}

// New new a service and return.
func New() (s *Service) {
	var ac = new(paladin.TOML)
	if err := paladin.Watch("application.toml", ac); err != nil {
		panic(err)
	}
	s = &Service{
		ac:  ac,
		dao: dao.New(),
	}
	if err := paladin.Get("application.toml").UnmarshalTOML(&s); err != nil {
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	log.Info("JWT secret-> %s", s.JWT.Secret)
	return s
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}
