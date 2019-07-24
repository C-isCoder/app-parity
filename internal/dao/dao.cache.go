package dao

import (
	"context"
	"github.com/bilibili/kratos/pkg/stat/prom"
	pb "app-parity/api"
)

func (d *Dao) GetDrugs(c context.Context, pageSize, pageNum int32) (res []*pb.Drug, err error) {
	addCache := true
	res, err = d.CacheDrugs(c, pageSize+pageNum)
	if err != nil {
		addCache = false
		err = nil
	}
	if res != nil {
		prom.CacheHit.Incr("Drugs")
		return
	}
	prom.CacheMiss.Incr("Drugs")
	res, err = d.DrugAllQuery(c, pageSize, pageNum)
	if err != nil {
		return
	}
	miss := res
	if !addCache {
		return
	}
	d.AddCacheDrugs(c, pageSize+pageNum, miss)
	return
}
