package service

import (
	"context"
	"signatureservice/internal/store"
)

// CollectVerificationBatch 消费流式产出的验签记录，并在出错或取消时返回已收集的部分结果。
//
// 它先排空 records 通道（这样即使中途遇到空记录，前面已读到的记录也能被收集），
// 再读取 errs 通道（生产者总会关闭两者，因此这里不会挂起）。
// 出错时返回已收集到的部分结果与错误，调用方既能看到错误，也不会丢失已读记录。
func CollectVerificationBatch(ctx context.Context, inputs []string) ([]string, error) {
	records, errs := store.StreamVerificationRecords(ctx, inputs)
	result := make([]string, 0, len(inputs))
	for record := range records {
		result = append(result, record)
	}
	if err := <-errs; err != nil {
		return result, err
	}
	return result, nil
}
