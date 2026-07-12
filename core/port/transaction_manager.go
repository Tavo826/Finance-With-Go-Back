package port

import "context"

// TransactionManager runs fn inside a single atomic unit of work.
// Implementations must propagate the transactional context to fn so that
// repository calls made with it participate in the same transaction.
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
