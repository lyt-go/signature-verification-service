package store

import (
	"context"
	"fmt"
)

func StreamVerificationRecords(ctx context.Context, inputs []string) (<-chan string, <-chan error) {
	records := make(chan string)
	errs := make(chan error)
	go func() {
		for _, input := range inputs {
			if input == "" {
				errs <- fmt.Errorf("empty verification record")
				return
			}
			select {
			case records <- input:
			case <-ctx.Done():
				return
			}
		}
		close(records)
	}()
	return records, errs
}
