# service

业务逻辑层，处于 server 层和 dao 层之间。service 只能通过 dao 层获取数据。

业务接口必须接受 `ctx context.Context` 对象，并向下传递。

## 错误日志
框架支持自动打印错误日志，使用方法:

```go
import "sniper/util/errors"

// ...

return errors.Wrap(err)
// 如果有附加信息，则可以
return errors.Wrap(err, "extram msg")
```

## eg:
```go
package drug

import (
	"context"
	"sniper/dao/drug"
	model "sniper/rpc/drug/v1"
	"sniper/util/errors"
)

func GetCount(ctx context.Context, key string) (int32, error) {
	rows, err := drug.QueryCount(ctx, key)
	if err != nil {
		return -1, errors.Wrap(err)
	}
	if rows.Next() {
		count := int32(-1)
		return count, rows.Scan(&count)
	}
	return 0, nil
}

func GetDrugs(ctx context.Context, key string, pageSize, pageNum int32) ([]*model.DrugData, error) {
	rows, err := drug.QueryDrug(ctx, key, pageSize, pageNum)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	drs := make([]*model.DrugData, 0)
	for rows.Next() {
		dr := &model.DrugData{}
		err := rows.Scan(
			&dr.WholesaleId, &dr.Level0, &dr.Level1, &dr.Level2, &dr.DrugName,
			&dr.ProviderId, &dr.ProviderName, &dr.Specification, &dr.Unit,
			&dr.Manufacturer, &dr.ValidDate, &dr.ChainPrice, &dr.DisPrice,
			&dr.MinPrice, &dr.MaxPrice, &dr.OldPrice, &dr.Price,
		)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		drs = append(drs, dr)
	}
	return drs, nil
}
```
