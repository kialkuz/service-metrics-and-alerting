package collectalgo

import (
	"context"
)

type CollectAlgo interface {
	CollectAndSend(ctx context.Context) error
}
