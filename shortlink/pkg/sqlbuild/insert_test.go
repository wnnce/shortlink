package sqlbuild

import (
	"fmt"
	"strconv"
	"testing"
)

type Book struct {
	ID   int
	Name string
}

func TestBatchInsertBuilder(t *testing.T) {
	books := make([]*Book, 0)
	for i := 0; i < 10; i++ {
		book := &Book{
			ID:   i,
			Name: "name" + strconv.Itoa(i),
		}
		books = append(books, book)
	}
	sql, args := BatchInsertBuilder[*Book]("t_book", books, func(book *Book) []any {
		return []any{book.ID, book.Name}
	}, "id", "name")
	fmt.Println(sql)
	fmt.Println(args)
}

func TestPostgresInsertBuilder_Insert(t *testing.T) {
	builder := NewInsertBuilder("t_blog_article").
		Insert("title", "测试标题").
		InsertRaw("create_time", "now()").
		InsertBySlice([]string{"sort", "status"}, []any{0, 0}).
		InsertByMap(map[string]any{
			"is_hot": false,
			"is_top": true,
		}).
		Returning("article_id")
	fmt.Println(builder.Sql())
	fmt.Println(builder.Args())
}
