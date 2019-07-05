# dao

数据访问层，负责访问 DB、MC、外部 HTTP 等接口，对上层屏蔽数据访问细节。

具体职责有：
- SQL 拼接和 DB 访问逻辑
- DB 的拆库折表逻辑
- DB 的缓存读写逻辑
- HTTP 接口调用逻辑

# eg

```go
package drug

import (
	"context"
	"database/sql"
	"sniper/util/db"
)

func QueryDrug(ctx context.Context, key string, pageSize, pageNum int32) (rows *sql.Rows, err error) {
	s := "SELECT * FROM drug WHERE drug_name LIKE ? OR manufacturer LIKE ? OR provider_name" +
		" LIKE ? LIMIT ? OFFSET ?"
	query := db.SQLSelect("drug", s)
	c := db.Get(ctx, "default")
	return c.QueryContext(
		ctx, query, "%"+key+"%", "%"+key+"%",
		"%"+key+"%", pageSize, pageSize*(pageNum-1),
	)
}

func QueryCount(ctx context.Context, key string) (rows *sql.Rows, err error) {
	s := "SELECT COUNT(*) FROM drug WHERE drug_name LIKE ? OR manufacturer LIKE ? OR provider_name LIKE ? "
	query := db.SQLSelect("drug", s)
	c := db.Get(ctx, "default")
	return c.QueryContext(ctx, query, "%"+key+"%", "%"+key+"%", "%"+key+"%")
}
```