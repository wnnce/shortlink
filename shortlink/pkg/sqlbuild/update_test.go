package sqlbuild

import (
	"fmt"
	"testing"
)

func TestPostgresUpdateBuilder_Set(t *testing.T) {
	builder := NewUpdateBuilder("t_blog_article").
		Set("title", "测试标题").
		Where("article_id").Eq(1).And("delete_at").EqRaw("0").BuildAsUpdate()
	fmt.Println(builder.Sql())
}

func TestPostgresUpdateBuilder_SetByCondition(t *testing.T) {
	builder := NewUpdateBuilder("t_blog_article").
		Set("title", "测试标题").
		SetByCondition(1 < 0, "status", 1).
		Where("article_id").Eq(1).BuildAsUpdate()
	fmt.Println(builder.Sql())
}

func TestPostgresUpdateBuilder_SetRow(t *testing.T) {
	builder := NewUpdateBuilder("t_blog_article").
		SetRaw("update_time", "now()").
		Set("title", "测试标题").
		Where("article_id").Between(1, 10).BuildAsUpdate()
	fmt.Println(builder.Sql())
	fmt.Println(builder.Args())
}

func TestPostgresUpdateBuilder_SetBySlice(t *testing.T) {
	builder := NewUpdateBuilder("t_blog_article").
		SetBySlice([]string{"title", "status"}, []any{"测试标题", 0}).
		Where("article_id").Eq(1).BuildAsUpdate()
	fmt.Println(builder.Sql())
	fmt.Println(builder.Args())
}

func TestPostgresUpdateBuilder_SetBySliceFail(t *testing.T) {
	builder := NewUpdateBuilder("t_blog_article").
		SetBySlice([]string{"title", "status"}, []any{"测试标题"}).
		Where("article_id").Eq(1).BuildAsUpdate()
	fmt.Println(builder.Sql())
	fmt.Println(builder.Args())
}

func TestPostgresUpdateBuilder_SetByMap(t *testing.T) {
	builder := NewUpdateBuilder("t_blog_article").
		SetByMap(map[string]any{
			"title":  "测试标题",
			"status": 1,
		}).Where("article_id").Eq(1).BuildAsUpdate()
	fmt.Println(builder.Sql())
	fmt.Println(builder.Args())
}
