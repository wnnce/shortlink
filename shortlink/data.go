package shortlink

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"shortlink/data"
	"shortlink/pb"
	"shortlink/pkg/sqlbuild"
)

const (
	TablePrefix = "tb_record_"
)

type ShortLinkData struct{}

func NewShortLinkData() *ShortLinkData {
	return &ShortLinkData{}
}

// 根据数据的base62值选择具体的表进行操作
func (self ShortLinkData) chooseTableName(baseKey string) string {
	return "\"\"\"" + TablePrefix + string(baseKey[len(baseKey)-1]) + "\"\"\""
}

// 根据不同的sql语句选择不同的数据源
func (self ShortLinkData) chooseDbPool(build sqlbuild.SqlBuilder) *pgxpool.Pool {
	if _, ok := build.(sqlbuild.SelectBuilder); ok {
		return data.ReplicaDb()
	}
	return data.MasterDb()
}

func (self ShortLinkData) Add(record *pb.LinkRecord) error {
	builder := sqlbuild.NewInsertBuilder(self.chooseTableName(record.BaseValue)).
		InsertByCondition(!record.IsLasting, "valid_hour", record.ValidHour).
		InsertByCondition(!record.IsLasting, "expire_time", record.ExpireTime).
		InsertByMap(map[string]any{
			"unique_id":   record.UniqueId,
			"base_value":  record.BaseValue,
			"origin_url":  record.OriginUrl,
			"is_lasting":  record.IsLasting,
			"create_time": record.CreateTime,
			"expire_mode": record.ExpireMode,
			"client_ip":   record.ClientIp,
			"user_agent":  record.UserAgent,
		}).Returning("id")
	row := self.chooseDbPool(builder).QueryRow(context.TODO(), builder.Sql(), builder.Args()...)
	return row.Scan(&record.Id)
}

func (self ShortLinkData) SelectByKey(baseKey string) (*pb.LinkRecord, error) {
	builder := sqlbuild.NewSelectBuilder(self.chooseTableName(baseKey)).
		Select("id", "unique_id", "base_value", "origin_url", "create_time", "expire_time", "status").
		Where("base_value").Eq(baseKey).BuildAsSelect()
	rows, err := self.chooseDbPool(builder).Query(context.TODO(), builder.Sql(), builder.Args()...)
	if rows.Next() {
		return pgx.RowToAddrOfStructByNameLax[pb.LinkRecord](rows)
	}
	return nil, err
}

func (self ShortLinkData) SelectInfoByKey(baseKey string) (*pb.LinkRecord, error) {
	builder := sqlbuild.NewSelectBuilder(self.chooseTableName(baseKey)).
		Select("id", "unique_id", "base_value", "origin_url", "valid_hour", "is_lasting", "create_time", "expire_time", "expire_mode", "client_ip", "user_agent", "status").
		Where("base_value").Eq(baseKey).BuildAsSelect()
	rows, err := self.chooseDbPool(builder).Query(context.TODO(), builder.Sql(), builder.Args()...)
	if rows.Next() {
		return pgx.RowToAddrOfStructByName[pb.LinkRecord](rows)
	}
	return nil, err
}

func (self ShortLinkData) DeleteByBaseKey(baseKey string) (int64, error) {
	builder := sqlbuild.NewDeleteBuilder(self.chooseTableName(baseKey)).
		Where("base_value").Eq(baseKey).BuildAsDelete()
	result, err := self.chooseDbPool(builder).Exec(context.TODO(), builder.Sql(), builder.Args()...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (self ShortLinkData) ListExpireLinkRecord(tableName string, expireTime, lastId, limit int64) ([]*pb.LinkRecord, error) {
	builder := sqlbuild.NewSelectBuilder(tableName).Select("id", "base_value", "expire_time", "expire_mode").
		Where("expire_time").Gt(expireTime).And("id").Gt(lastId).BuildAsSelect().
		OrderByAsc("id").
		Limit(limit)
	rows, err := self.chooseDbPool(builder).Query(context.TODO(), builder.Sql(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*pb.LinkRecord, error) {
		return pgx.RowToAddrOfStructByNameLax[pb.LinkRecord](row)
	})
}
