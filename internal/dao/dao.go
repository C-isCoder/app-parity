package dao

import (
	"context"
	"fmt"
	"time"

	pb "app-parity/api"
	"github.com/bilibili/kratos/pkg/cache/memcache"
	"github.com/bilibili/kratos/pkg/cache/redis"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
	"github.com/bilibili/kratos/pkg/log"
	xtime "github.com/bilibili/kratos/pkg/time"
)

//go:generate kratos tool genmc
type _mc interface {
	// mc: -key=cacheId
	CacheDrugs(c context.Context, sum int32) ([]*pb.Drug, error)
	// mc: -key=noneKey
	CacheNone(c context.Context) (*pb.Drug, error)
	// mc: -key=cacheKey
	CacheString(c context.Context, key string) (string, error)
	// mc: -key=cacheId -expire=d.mcExpire -encode=json
	AddCacheDrugs(c context.Context, sum int32, values []*pb.Drug) error
	// mc: -key=noneKey
	AddCacheNone(c context.Context, value *pb.DrugsResp) error
	// mc: -key=cacheKey -expire=d.mcExpire
	AddCacheString(c context.Context, key string, value string) error

	// mc: -key=cacheId
	DelCacheDrugs(c context.Context, sum int32) error
	// mc: -key=noneKey
	DelCacheNone(c context.Context) error
}

func cacheId(sum int32) string {
	return fmt.Sprintf("num_%d", sum)
}
func cacheKey(key string) string {
	return fmt.Sprintf("key_%s", key)
}
func noneKey() string {
	return "none"
}

// Dao dao.
type Dao struct {
	db          *sql.DB
	redis       *redis.Pool
	redisExpire int32
	mc          *memcache.Memcache
	mcExpire    int32
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// New new a dao and return.
func New() (dao *Dao) {
	var (
		dc struct {
			MySql *sql.Config
		}
		rc struct {
			Redis       *redis.Config
			RedisExpire xtime.Duration
		}
		mc struct {
			Mem       *memcache.Config
			MemExpire xtime.Duration
		}
	)
	checkErr(paladin.Get("mysql.toml").UnmarshalTOML(&dc))
	checkErr(paladin.Get("redis.toml").UnmarshalTOML(&rc))
	checkErr(paladin.Get("memcache.toml").UnmarshalTOML(&mc))
	dao = &Dao{
		// mysql
		db: sql.NewMySQL(dc.MySql),
		// redis
		redis:       redis.NewPool(rc.Redis),
		redisExpire: int32(time.Duration(rc.RedisExpire) / time.Second),
		// memcache
		mc:       memcache.New(mc.Mem),
		mcExpire: int32(time.Duration(mc.MemExpire) / time.Second),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.redis.Close()
	d.db.Close()
}

// Ping ping the resource.
func (d *Dao) Ping(ctx context.Context) (err error) {
	if err = d.pingMC(ctx); err != nil {
		return
	}
	if err = d.pingRedis(ctx); err != nil {
		return
	}
	return d.db.Ping(ctx)
}

func (d *Dao) pingMC(ctx context.Context) (err error) {
	if err = d.mc.Set(ctx, &memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}

func (d *Dao) pingRedis(ctx context.Context) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("SET", "ping", "pong"); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}

// format SQL LIKE
func like(str string) string {
	return fmt.Sprintf("%%%s%%", str)
}
