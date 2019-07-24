package dao

import (
	"context"
	"github.com/bilibili/kratos/pkg/database/sql"
	"github.com/bilibili/kratos/pkg/log"
	pb "app-parity/api"
	"time"
)

const (
	_user    = "SELECT `uid`,`account`,`password`,`name`,`is_admin` FROM `user` WHERE `account` = ?"
	_addUser = "INSERT INTO `user` (`account`, `password`, `name`, `is_admin`, `create`, `update`) VALUES (?, ?, ?, ?, ?, ?)"
)

func (d *Dao) UserQuery(ctx context.Context, account string) *pb.User {
	user := new(pb.User)
	err := d.db.QueryRow(ctx, _user, account).Scan(&user.Uid, &user.Account,
		&user.Password,
		&user.Name, &user.IsAdmin)
	if err != nil && err != sql.ErrNoRows {
		log.Error("db.UserQuery.Query error(%v)", err)
		return nil
	}
	return user
}

func (d *Dao) UserInsert(ctx context.Context, account, password, name string) (int64, error) {
	t := time.Now().Format("2006-01-02 15:04")
	res, err := d.db.Exec(ctx, _addUser, account, password, name, false, t, t)
	if err != nil {
		log.Error("db.UserInsert.Exec(%s) error(%v)", _addUser, err)
	}
	return res.LastInsertId()
}
