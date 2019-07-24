package service

import (
	pb "app-parity/api"
	"context"
	"github.com/bilibili/kratos/pkg/ecode"
	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/net/http/blademaster/middleware/auth"
)

func (s *Service) Login(ctx context.Context, req *pb.LoginReq) (resp *pb.LoginResp, err error) {
	if req.Account == "" {
		return nil, ecode.Error(2000, "请输入账号")
	}
	if req.Password == "" {
		return nil, ecode.Error(2000, "请输入密码")
	}
	user := s.dao.UserQuery(ctx, req.Account)
	if user == nil {
		return nil, ecode.Error(2000, "用户不存在")
	}
	if user.Password == req.Password {
		resp = new(pb.LoginResp)
		resp.Name = user.Name
		resp.Account = user.Account
		resp.Token = auth.NewToken(s.JWT.Secret, user.Uid, user.Name, user.IsAdmin).String()
		resp.IsAdmin = user.IsAdmin
		return resp, nil
	} else {
		return nil, ecode.Error(2000, "账号或密码错误")
	}
}

func (s *Service) Register(ctx context.Context, req *pb.RegisterReq) (resp *pb.LoginResp, err error) {
	if req.Account == "" {
		return nil, ecode.Error(2000, "请输入账号")
	}
	if req.Password == "" {
		return nil, ecode.Error(2000, "请输入密码")
	}
	user := s.dao.UserQuery(ctx, req.Account)
	log.Info("query user(%s) size:%v", user.String(), user.Size())
	if user.Size() != 0 {
		return nil, ecode.Error(2000, "用户已存在")
	}
	id, err := s.dao.UserInsert(ctx, req.Account, req.Password, req.Name)
	if err != nil {
		log.Error("user.service register error(%v)", err)
		return nil, err
	}
	resp = &pb.LoginResp{
		Name:    req.Name,
		Account: req.Account,
		Token:   auth.NewToken(s.JWT.Secret, id, req.Name, false, ).String(),
	}
	return resp, nil
}
