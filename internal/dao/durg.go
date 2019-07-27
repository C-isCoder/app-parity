package dao

import (
	pb "app-parity/api/data"
	"context"
	"github.com/bilibili/kratos/pkg/log"
)

const (
	_drug = "SELECT a.*, b.approval FROM `drug` as a, " +
		"detail as b WHERE (a.Drug_name LIKE ? OR a. manufacturer LIKE ? OR a." +
		"provider_name LIKE ? OR b. approval LIKE ?) AND a.wholesale_id = b.wholesale_id LIMIT ? OFFSET ?"
	_count = "SELECT COUNT(*) FROM `drug` as a, " +
		"detail as b WHERE (a.Drug_name LIKE ? OR a. manufacturer LIKE ? OR a." +
		"provider_name LIKE ? OR b. approval LIKE ?) AND a.wholesale_id = b.wholesale_id"
	_drugs = "SELECT * FROM `drug` LIMIT ? OFFSET ?"
)

func (d *Dao) DrugCount(ctx context.Context, key string) (count int32, err error) {
	err = d.db.QueryRow(ctx, _count, like(key), like(key), like(key), like(key)).Scan(&count)
	log.Info("count: %d", count)
	if err != nil {
		log.Error("d.DrugCount() error(%v)", err)
		return
	}
	return
}

func (d *Dao) DrugQuery(ctx context.Context, key string, pageSize, pageNum int32) (ds []*pb.Drug,
	err error) {
	rows, err := d.db.Query(ctx, _drug, like(key), like(key), like(key), like(key), pageSize,
		pageSize*(pageNum-1))
	if err != nil {
		log.Error("d.DrugQuery() error(%v)", err)
		return
	}
	defer rows.Close()
	ds = make([]*pb.Drug, 0)
	for rows.Next() {
		d := new(pb.Drug)
		rows.Scan(&d.WholesaleId, &d.Level0, &d.Level1, &d.Level2, &d.DrugName,
			&d.ProviderId, &d.ProviderName, &d.Specification, &d.Unit,
			&d.Manufacturer, &d.ValidDate, &d.ChainPrice, &d.DisPrice,
			&d.MinPrice, &d.MaxPrice, &d.OldPrice, &d.Price, &d.ApprovalNumber, )
		ds = append(ds, d)
	}
	return
}

func (d *Dao) DrugAllQuery(ctx context.Context, pageSize, pageNum int32) (ds []*pb.Drug, err error) {
	rows, err := d.db.Query(ctx, _drugs, pageSize, pageSize*(pageNum-1))
	if err != nil {
		log.Error("d.DrugAllQuery() error(%v)", err, pageSize, pageSize*(pageNum-1))
		return
	}
	defer rows.Close()
	ds = make([]*pb.Drug, 0)
	for rows.Next() {
		d := new(pb.Drug)
		rows.Scan(&d.WholesaleId, &d.Level0, &d.Level1, &d.Level2, &d.DrugName,
			&d.ProviderId, &d.ProviderName, &d.Specification, &d.Unit,
			&d.Manufacturer, &d.ValidDate, &d.ChainPrice, &d.DisPrice,
			&d.MinPrice, &d.MaxPrice, &d.OldPrice, &d.Price, &d.ApprovalNumber, )
		ds = append(ds, d)
	}
	return
}
