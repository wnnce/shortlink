package data

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"shortlink/config"
)

// DbConfigureReader Db包的配置文件读取器
// 读取配置文件后完成数据库连接的初始化
func DbConfigureReader(bootstrap config.Bootstrap) (func(), error) {
	masterPool, err := makePgxPool(bootstrap.Data.Database.Master)
	if err != nil {
		return nil, err
	}
	replicas := bootstrap.Data.Database.Replicas
	replicaPools := make([]*pgxpool.Pool, 0, len(replicas))
	for _, replica := range replicas {
		replicaPool, err := makePgxPool(replica)
		if err != nil {
			return nil, err
		}
		replicaPools = append(replicaPools, replicaPool)
	}
	defaultDbData = &dbData{
		master:      masterPool,
		replicas:    replicaPools,
		loadBalance: NewPollLoadBalance(int32(len(replicaPools))),
	}
	return func() {
		masterPool.Close()
		for _, replica := range replicaPools {
			replica.Close()
		}
	}, nil
}

func makePgxPool(database config.Database) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", database.Host, database.Port, database.Username, database.Password, database.DbName)
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		slog.Error("数据库连接创建失败", slog.String("host", database.Host), slog.Int("port", database.Port),
			slog.String("username", database.Username), slog.String("password", database.Password),
			slog.String("dbName", database.DbName))
		return nil, err
	}
	slog.Info("数据库连接创建成功", slog.String("host", database.Host), slog.Int("port", database.Port),
		slog.String("username", database.Username), slog.String("password", database.Password),
		slog.String("dbName", database.DbName))
	return db, nil
}

var defaultDbData *dbData

type dbData struct {
	master      *pgxpool.Pool
	replicas    []*pgxpool.Pool
	loadBalance LoadBalance
}

func (self *dbData) MasterPool() *pgxpool.Pool {
	return self.master
}

func (self *dbData) ReplicaPools() []*pgxpool.Pool {
	return self.replicas
}

func (self *dbData) ReplicaPool() *pgxpool.Pool {
	if self.replicas == nil || len(self.replicas) == 0 {
		return self.master
	}
	index := self.loadBalance.GetIndex()
	return self.replicas[index]
}

func MasterDb() *pgxpool.Pool {
	if defaultDbData == nil {
		return nil
	}
	return defaultDbData.MasterPool()
}

func ReplicaDb() *pgxpool.Pool {
	if defaultDbData == nil {
		return nil
	}
	return defaultDbData.ReplicaPool()
}
