package service

import (
	"context"
	"testing"

	"personal-finance/core/domain"
)

// --- mocks ---

type mockTransactionRepo struct {
	getByIdFunc func(ctx context.Context, id string) (*domain.Transaction, error)
	updateFunc  func(ctx context.Context, id string, tx *domain.Transaction) (*domain.Transaction, error)
}

func (m *mockTransactionRepo) GetTransactionsByUserId(ctx context.Context, page, limit uint64, userId string) ([]domain.Transaction, int64, int, error) {
	return nil, 0, 0, nil
}

func (m *mockTransactionRepo) GetTransactionsByDate(ctx context.Context, userId string, page, limit uint64, year int, month int) ([]domain.Transaction, int64, int, error) {
	return nil, 0, 0, nil
}

func (m *mockTransactionRepo) GetTransactionsBySubject(ctx context.Context, userId string, page, limit uint64, subject string, personOrBusiness string) ([]domain.Transaction, int64, int, error) {
	return nil, 0, 0, nil
}

func (m *mockTransactionRepo) GetTransactionById(ctx context.Context, id string) (*domain.Transaction, error) {
	return m.getByIdFunc(ctx, id)
}

func (m *mockTransactionRepo) CreateTransaction(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	return tx, nil
}

func (m *mockTransactionRepo) UpdateTransaction(ctx context.Context, id string, tx *domain.Transaction) (*domain.Transaction, error) {
	return m.updateFunc(ctx, id, tx)
}

func (m *mockTransactionRepo) DeleteTransaction(ctx context.Context, id string) error {
	return nil
}

func (m *mockTransactionRepo) DeleteTransactionsByUserId(ctx context.Context, id string) error {
	return nil
}

type originUpdateCall struct {
	id     string
	typ    string
	amount float64
	total  float64 // resulting Total after applying the delta
}

type mockOriginRepo struct {
	origins   map[string]*domain.Origin
	getErr    error
	updateLog []originUpdateCall
}

func newMockOriginRepo(origins map[string]*domain.Origin) *mockOriginRepo {
	return &mockOriginRepo{origins: origins}
}

func (m *mockOriginRepo) GetOriginsByUserId(ctx context.Context, userId string) ([]domain.Origin, error) {
	return nil, nil
}

func (m *mockOriginRepo) GetOriginById(ctx context.Context, id string) (*domain.Origin, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	o, ok := m.origins[id]
	if !ok {
		return nil, domain.ErrDataNotFound
	}
	copy := *o
	return &copy, nil
}

func (m *mockOriginRepo) CreateOrigin(ctx context.Context, origin *domain.Origin) (*domain.Origin, error) {
	return origin, nil
}

func (m *mockOriginRepo) UpdateOrigin(ctx context.Context, id string, updated *domain.Origin) (*domain.Origin, error) {
	m.origins[id] = updated
	return updated, nil
}

func (m *mockOriginRepo) DeleteOrigin(ctx context.Context, id string) error {
	return nil
}

// --- helpers ---

func strPtr(s string) *string { return &s }

func newTransactionService(tRepo *mockTransactionRepo, oRepo *mockOriginRepo) *TransactionService {
	return NewTransactionService(tRepo, oRepo)
}

// --- UpdateTotalOrigin ---

func TestUpdateTotalOrigin_Income(t *testing.T) {
	oRepo := newMockOriginRepo(map[string]*domain.Origin{
		"o1": {ID: "o1", Total: 100},
	})
	ts := newTransactionService(&mockTransactionRepo{}, oRepo)

	if err := ts.UpdateTotalOrigin(context.Background(), "o1", "Income", 50); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := oRepo.origins["o1"].Total; got != 150 {
		t.Errorf("expected total 150, got %v", got)
	}
}

func TestUpdateTotalOrigin_Output(t *testing.T) {
	oRepo := newMockOriginRepo(map[string]*domain.Origin{
		"o1": {ID: "o1", Total: 100},
	})
	ts := newTransactionService(&mockTransactionRepo{}, oRepo)

	if err := ts.UpdateTotalOrigin(context.Background(), "o1", "Output", 30); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := oRepo.origins["o1"].Total; got != 70 {
		t.Errorf("expected total 70, got %v", got)
	}
}

func TestUpdateTotalOrigin_OriginNotFound(t *testing.T) {
	oRepo := newMockOriginRepo(map[string]*domain.Origin{})
	ts := newTransactionService(&mockTransactionRepo{}, oRepo)

	err := ts.UpdateTotalOrigin(context.Background(), "missing", "Income", 10)
	if err != domain.ErrDataNotFound {
		t.Fatalf("expected ErrDataNotFound, got %v", err)
	}
}

// --- UpdateTransaction / reconcileOriginBalance ---

func TestUpdateTransaction_SameOrigin_TypeChanges(t *testing.T) {
	actual := &domain.Transaction{OriginId: strPtr("o1"), Type: "Income", Amount: 100}
	updated := &domain.Transaction{OriginId: strPtr("o1"), Type: "Output", Amount: 50}

	oRepo := newMockOriginRepo(map[string]*domain.Origin{
		"o1": {ID: "o1", Total: 200},
	})
	tRepo := &mockTransactionRepo{
		getByIdFunc: func(ctx context.Context, id string) (*domain.Transaction, error) { return actual, nil },
		updateFunc: func(ctx context.Context, id string, tx *domain.Transaction) (*domain.Transaction, error) {
			return tx, nil
		},
	}
	ts := newTransactionService(tRepo, oRepo)

	if _, err := ts.UpdateTransaction(context.Background(), "t1", updated); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// undo the +100 income, apply the -50 output => 200 - 100 - 50 = 50
	if got := oRepo.origins["o1"].Total; got != 50 {
		t.Errorf("expected total 50, got %v", got)
	}
}

func TestUpdateTransaction_SameOrigin_AmountIncreases(t *testing.T) {
	actual := &domain.Transaction{OriginId: strPtr("o1"), Type: "Income", Amount: 100}
	updated := &domain.Transaction{OriginId: strPtr("o1"), Type: "Income", Amount: 150}

	oRepo := newMockOriginRepo(map[string]*domain.Origin{
		"o1": {ID: "o1", Total: 500},
	})
	tRepo := &mockTransactionRepo{
		getByIdFunc: func(ctx context.Context, id string) (*domain.Transaction, error) { return actual, nil },
		updateFunc: func(ctx context.Context, id string, tx *domain.Transaction) (*domain.Transaction, error) {
			return tx, nil
		},
	}
	ts := newTransactionService(tRepo, oRepo)

	if _, err := ts.UpdateTransaction(context.Background(), "t1", updated); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// only the +50 delta should be applied => 500 + 50 = 550
	if got := oRepo.origins["o1"].Total; got != 550 {
		t.Errorf("expected total 550, got %v", got)
	}
}

func TestUpdateTransaction_SameOrigin_AmountDecreases(t *testing.T) {
	actual := &domain.Transaction{OriginId: strPtr("o1"), Type: "Income", Amount: 150}
	updated := &domain.Transaction{OriginId: strPtr("o1"), Type: "Income", Amount: 100}

	oRepo := newMockOriginRepo(map[string]*domain.Origin{
		"o1": {ID: "o1", Total: 500},
	})
	tRepo := &mockTransactionRepo{
		getByIdFunc: func(ctx context.Context, id string) (*domain.Transaction, error) { return actual, nil },
		updateFunc: func(ctx context.Context, id string, tx *domain.Transaction) (*domain.Transaction, error) {
			return tx, nil
		},
	}
	ts := newTransactionService(tRepo, oRepo)

	if _, err := ts.UpdateTransaction(context.Background(), "t1", updated); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// only the -50 delta should be applied => 500 - 50 = 450
	if got := oRepo.origins["o1"].Total; got != 450 {
		t.Errorf("expected total 450, got %v", got)
	}
}

func TestUpdateTransaction_SameOrigin_NoChange_SkipsUpdate(t *testing.T) {
	actual := &domain.Transaction{OriginId: strPtr("o1"), Type: "Income", Amount: 100}
	updated := &domain.Transaction{OriginId: strPtr("o1"), Type: "Income", Amount: 100}

	// GetOriginById would return ErrDataNotFound if called; origins map is empty
	// on purpose so the test fails loudly if reconciliation wrongly touches the origin.
	oRepo := newMockOriginRepo(map[string]*domain.Origin{})
	tRepo := &mockTransactionRepo{
		getByIdFunc: func(ctx context.Context, id string) (*domain.Transaction, error) { return actual, nil },
		updateFunc: func(ctx context.Context, id string, tx *domain.Transaction) (*domain.Transaction, error) {
			return tx, nil
		},
	}
	ts := newTransactionService(tRepo, oRepo)

	if _, err := ts.UpdateTransaction(context.Background(), "t1", updated); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateTransaction_OriginChanges(t *testing.T) {
	actual := &domain.Transaction{OriginId: strPtr("o1"), Type: "Income", Amount: 100}
	updated := &domain.Transaction{OriginId: strPtr("o2"), Type: "Output", Amount: 40}

	oRepo := newMockOriginRepo(map[string]*domain.Origin{
		"o1": {ID: "o1", Total: 300},
		"o2": {ID: "o2", Total: 500},
	})
	tRepo := &mockTransactionRepo{
		getByIdFunc: func(ctx context.Context, id string) (*domain.Transaction, error) { return actual, nil },
		updateFunc: func(ctx context.Context, id string, tx *domain.Transaction) (*domain.Transaction, error) {
			return tx, nil
		},
	}
	ts := newTransactionService(tRepo, oRepo)

	if _, err := ts.UpdateTransaction(context.Background(), "t1", updated); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// o1 loses the reverted +100 income => 300 - 100 = 200
	if got := oRepo.origins["o1"].Total; got != 200 {
		t.Errorf("expected o1 total 200, got %v", got)
	}
	// o2 gets the new -40 output applied => 500 - 40 = 460
	if got := oRepo.origins["o2"].Total; got != 460 {
		t.Errorf("expected o2 total 460, got %v", got)
	}
}

// TestUpdateTransaction_NilActualOriginId_Panics documents a bug in
// reconcileOriginBalance: it dereferences actualTransaction.OriginId
// (line "*actualTransaction.OriginId != *updatedTransaction.OriginId")
// before the subsequent nil check, so a transaction loaded from storage
// without an origin (OriginId == nil) crashes UpdateTransaction instead
// of being handled by the "actual.OriginId == nil" branch below it.
func TestUpdateTransaction_NilActualOriginId_Panics(t *testing.T) {
	actual := &domain.Transaction{OriginId: nil, Type: "Income", Amount: 100}
	updated := &domain.Transaction{OriginId: strPtr("o1"), Type: "Income", Amount: 100}

	oRepo := newMockOriginRepo(map[string]*domain.Origin{
		"o1": {ID: "o1", Total: 0},
	})
	tRepo := &mockTransactionRepo{
		getByIdFunc: func(ctx context.Context, id string) (*domain.Transaction, error) { return actual, nil },
		updateFunc: func(ctx context.Context, id string, tx *domain.Transaction) (*domain.Transaction, error) {
			return tx, nil
		},
	}
	ts := newTransactionService(tRepo, oRepo)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected a nil-pointer panic (known bug in reconcileOriginBalance); UpdateTransaction returned normally instead - the bug may have been fixed, update this test")
		}
	}()

	ts.UpdateTransaction(context.Background(), "t1", updated)
}
