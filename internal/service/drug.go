package service

import (
	"context"
	pb "app-parity/api"
)

// 药品关键字搜索
func (s *Service) DrugSearch(ctx context.Context, req *pb.SearchReq) (resp *pb.SearchResp,
	err error) {
	resp = new(pb.SearchResp)

	resp.Count, err = s.dao.DrugCount(ctx, req.Key)
	if err != nil {
		return
	}
	resp.Drugs, err = s.dao.DrugQuery(ctx, req.Key, req.PageSize, req.PageNum)
	if err != nil {
		return
	}
	return
}

// 获取所有药品列表
func (s *Service) GetDrugs(ctx context.Context, req *pb.DrugsReq) (resp *pb.DrugsResp, err error) {
	resp = new(pb.DrugsResp)
	resp.Drugs, err = s.dao.DrugAllQuery(ctx, req.PageSize, req.PageNum)
	if err != nil {
		return
	}
	return
}
