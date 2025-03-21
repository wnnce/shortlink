package sqlbuild

import (
	"fmt"
	"testing"
	"time"
)

func TestPostgresSelectBuilder_Select(t *testing.T) {
	builder := NewSelectBuilder("example").
		Select("id", "name", "age")
	fmt.Println(builder.Sql())
}
func TestSelectBuilder(t *testing.T) {
	fmt.Println(time.Now().UnixMilli())
	builder := NewSelectBuilder("t_blog_article as ba").
		Select("article_id", "title", "sort").
		LeftJoin("t_blog_comment as bc").On("bc.article_id").EqRaw("ba.article_id").
		And("bc.delete_at").EqRaw("0").BuildAsSelect().
		Select("count(bc.*) as comment_num").
		Where("ba.article_id").In(1, 2, 3).
		And("ba.delete_at").EqRaw("0").
		And("ba.status").EqRaw("0").
		AndByCondition(1 < 0, "article_id").LtRaw("0").BuildAsSelect().
		GroupBy("ba.article_id", "ba.sort").
		OrderBy("ba.article_id", "ba.sort")
	fmt.Println(builder.Sql())
	fmt.Println(builder.CountSql())
	fmt.Println(builder.Args())
	fmt.Println(time.Now().UnixMilli())
}

func BenchmarkSelectBuilder(b *testing.B) {
	builder := NewSelectBuilder("t_blog_article as ba").
		Select("article_id", "title", "sort").
		LeftJoin("t_blog_comment as bc").On("bc.article_id").EqRaw("ba.article_id").
		And("bc.delete_at").EqRaw("0").BuildAsSelect().
		Select("count(bc.*) as comment_num").
		Where("ba.article_id").In(1, 2, 3).
		And("ba.delete_at").EqRaw("0").
		And("ba.status").EqRaw("0").
		AndByCondition(1 < 0, "article_id").LtRaw("0").BuildAsSelect().
		GroupBy("ba.article_id", "ba.sort").
		OrderBy("ba.article_id", "ba.sort")
	fmt.Println(builder.Sql())
	fmt.Println(builder.CountSql())
	fmt.Println(builder.Args())
}
