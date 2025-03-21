package sqlbuild

import (
	"fmt"
	"testing"
)

func TestPostgresDeleteBuilder(t *testing.T) {
	builder := NewDeleteBuilder("t_blog_article").
		Where("article_id").Le(10).
		And("delete_at").EqRaw("0").BuildAsDelete()
	fmt.Println(builder.Sql())
	fmt.Println(builder.Args())
}
