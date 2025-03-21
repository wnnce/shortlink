package sqlbuild

import (
	"strings"
	"sync"
)

// SqlBuilder SQL构造器接口
type SqlBuilder interface {
	table(tableName string) SqlBuilder
	Sql() string
	StatementParameterBuffer
}

func NewSelectBuilder(table string) SelectBuilder {
	selectBuilder := &PostgresSelectBuilder{
		tableName: table,
		limit:     -1,
		offset:    -1,
	}
	selectBuilder.PostgresCondition.builder = selectBuilder
	return selectBuilder
}

func NewUpdateBuilder(table string) UpdateBuilder {
	updateBuilder := &PostgresUpdateBuilder{
		tableName: table,
	}
	updateBuilder.PostgresCondition.builder = updateBuilder
	return updateBuilder
}

func NewDeleteBuilder(table string) DeleteBuilder {
	deleteBuilder := &PostgresDeleteBuilder{
		tableName: table,
	}
	deleteBuilder.PostgresCondition.builder = deleteBuilder
	return deleteBuilder
}

func NewInsertBuilder(table string) InsertBuilder {
	return &PostgresInsertBuilder{
		tableName: table,
	}
}

var defaultPool *sqlBuilderPool

func init() {
	defaultPool = newSqlBuilderPool()
}

// 字段条件和stringBuilder对象池
type sqlBuilderPool struct {
	fieldPool  *sync.Pool
	bufferPool *sync.Pool
}

func newSqlBuilderPool() *sqlBuilderPool {
	return &sqlBuilderPool{
		fieldPool: &sync.Pool{
			New: func() any {
				return &Field{}
			},
		},
		bufferPool: &sync.Pool{
			New: func() any {
				return &strings.Builder{}
			},
		},
	}
}

func (self *sqlBuilderPool) GetField() *Field {
	return self.fieldPool.Get().(*Field)
}
func (self *sqlBuilderPool) GetStringBuilder() *strings.Builder {
	return self.bufferPool.Get().(*strings.Builder)
}

// RecycleField 回收Field
func (self *sqlBuilderPool) RecycleField(field *Field) {
	field.Recycle()
	self.fieldPool.Put(field)
}

// RecycleStringBuilder 回收stringBuilder
func (self *sqlBuilderPool) RecycleStringBuilder(builder *strings.Builder) {
	builder.Reset()
	self.bufferPool.Put(builder)
}
