package store

import (
	"context"
	"fmt"
)

// StreamVerificationRecords 以流式方式逐条产出验签记录。
//
// records 在所有有效记录产出完毕、出错或上下文取消后关闭；errs 至多承载一个错误后关闭。
// 遇到空记录会立即停止并发送错误，已产出的有效记录仍可被消费方读取——消费方先排空
// records 再读取 errs，因此不会因单条空记录而挂起。
func StreamVerificationRecords(ctx context.Context, inputs []string) (<-chan string, <-chan error) {
	records := make(chan string)
	// errs 缓冲为 1：发送错误时无需等待消费方即可立即返回并关闭 records，从而避免死锁。
	errs := make(chan error, 1)
	go func() {
		defer close(records)
		defer close(errs)
		for _, input := range inputs {
			if input == "" {
				errs <- fmt.Errorf("empty verification record")
				return
			}
			select {
			case records <- input:
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			}
		}
	}()
	return records, errs
}
