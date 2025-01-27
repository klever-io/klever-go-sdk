package nonceHandlerV2

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/klever-io/klever-go-sdk/interactors"
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/testsCommon"
	"github.com/klever-io/klever-go/data/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNonceTransactionHandlerV2(t *testing.T) {
	t.Parallel()

	t.Run("nil proxy", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNonceTransactionsHandlerV2()
		args.Proxy = nil
		nth, err := NewNonceTransactionHandlerV2(args)
		require.Nil(t, nth)
		assert.Equal(t, interactors.ErrNilProxy, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNonceTransactionsHandlerV2()
		nth, err := NewNonceTransactionHandlerV2(args)
		require.NotNil(t, nth)
		require.Nil(t, err)

		require.Nil(t, nth.Close())
	})
}

func TestNonceTransactionsHandlerV2_GetNonce(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)
	numCalls := 0

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address address.Address) (*models.Account, error) {
			addressAsBech32String := address.Bech32()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			numCalls++

			return &models.Account{
				AccountInfo: &models.AccountInfo{
					Nonce: currentNonce,
				},
			}, nil
		},
	}

	nth, _ := NewNonceTransactionHandlerV2(args)
	err := nth.ApplyNonceAndGasPrice(context.Background(), nil, nil)
	assert.Equal(t, interactors.ErrNilAddress, err)

	tx := transaction.Transaction{}
	err = nth.ApplyNonceAndGasPrice(context.Background(), testAddress, &tx)
	assert.Nil(t, err)
	assert.Equal(t, currentNonce, tx.RawData.Nonce)

	err = nth.ApplyNonceAndGasPrice(context.Background(), testAddress, &tx)
	assert.Nil(t, err)
	assert.Equal(t, currentNonce+1, tx.RawData.Nonce)

	assert.Equal(t, 2, numCalls)

	require.Nil(t, nth.Close())
}

func TestNonceTransactionsHandlerV2_SendMultipleTransactionsResendingEliminatingOne(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)

	mutSentTransactions := sync.Mutex{}
	numCalls := 0
	sentTransactions := make(map[int][]*transaction.Transaction)

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address address.Address) (*models.Account, error) {
			addressAsBech32String := address.Bech32()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &models.Account{
				AccountInfo: &models.AccountInfo{
					Nonce: atomic.LoadUint64(&currentNonce),
				},
			}, nil
		},
		SendTransactionsCalled: func(txs []*transaction.Transaction) ([]string, error) {
			mutSentTransactions.Lock()
			defer mutSentTransactions.Unlock()

			sentTransactions[numCalls] = txs
			numCalls++
			hashes := make([]string, len(txs))

			return hashes, nil
		},
		SendTransactionCalled: func(tx *transaction.Transaction) (string, error) {
			mutSentTransactions.Lock()
			defer mutSentTransactions.Unlock()

			sentTransactions[numCalls] = []*transaction.Transaction{tx}
			numCalls++

			return "", nil
		},
	}
	nth, _ := NewNonceTransactionHandlerV2(args)

	numTxs := 5
	txs := createMockTransactions(testAddress, numTxs, atomic.LoadUint64(&currentNonce))
	for i := 0; i < numTxs; i++ {
		_, err := nth.SendTransaction(context.TODO(), txs[i])
		require.Nil(t, err)
	}

	time.Sleep(time.Second * 3)
	_ = nth.Close()

	mutSentTransactions.Lock()
	defer mutSentTransactions.Unlock()

	numSentTransaction := 5
	numSentTransactions := 1
	assert.Equal(t, numSentTransaction+numSentTransactions, len(sentTransactions))
	for i := 0; i < numSentTransaction; i++ {
		assert.Equal(t, 1, len(sentTransactions[i]))
	}
	assert.Equal(t, numTxs-1, len(sentTransactions[numSentTransaction])) // resend
}

func TestNonceTransactionsHandlerV2_SendMultipleTransactionsResendingEliminatingAll(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)

	mutSentTransactions := sync.Mutex{}
	numCalls := 0
	sentTransactions := make(map[int][]*transaction.Transaction)

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address address.Address) (*models.Account, error) {
			addressAsBech32String := address.Bech32()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &models.Account{
				AccountInfo: &models.AccountInfo{
					Nonce: atomic.LoadUint64(&currentNonce),
				},
			}, nil
		},
		SendTransactionCalled: func(tx *transaction.Transaction) (string, error) {
			mutSentTransactions.Lock()
			defer mutSentTransactions.Unlock()

			sentTransactions[numCalls] = []*transaction.Transaction{tx}
			numCalls++

			return "", nil
		},
	}
	numTxs := 5
	nth, _ := NewNonceTransactionHandlerV2(args)
	txs := createMockTransactions(testAddress, numTxs, atomic.LoadUint64(&currentNonce))
	for i := 0; i < numTxs; i++ {
		_, err := nth.SendTransaction(context.Background(), txs[i])
		require.Nil(t, err)
	}

	atomic.AddUint64(&currentNonce, uint64(numTxs))
	time.Sleep(time.Second * 3)
	_ = nth.Close()

	mutSentTransactions.Lock()
	defer mutSentTransactions.Unlock()

	//no resend operation was made because all transactions were executed (nonce was incremented)
	assert.Equal(t, 5, len(sentTransactions))
	assert.Equal(t, 1, len(sentTransactions[0]))
}

func TestNonceTransactionsHandlerV2_SendTransactionResendingEliminatingAll(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)

	mutSentTransactions := sync.Mutex{}
	numCalls := 0
	sentTransactions := make(map[int][]*transaction.Transaction)

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address address.Address) (*models.Account, error) {
			addressAsBech32String := address.Bech32()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &models.Account{
				AccountInfo: &models.AccountInfo{
					Nonce: atomic.LoadUint64(&currentNonce),
				},
			}, nil
		},
		SendTransactionCalled: func(tx *transaction.Transaction) (string, error) {
			mutSentTransactions.Lock()
			defer mutSentTransactions.Unlock()

			sentTransactions[numCalls] = []*transaction.Transaction{tx}
			numCalls++

			return "", nil
		},
	}

	numTxs := 1
	nth, _ := NewNonceTransactionHandlerV2(args)
	txs := createMockTransactions(testAddress, numTxs, atomic.LoadUint64(&currentNonce))

	hash, err := nth.SendTransaction(context.Background(), txs[0])
	require.Nil(t, err)
	require.Equal(t, "", hash)

	atomic.AddUint64(&currentNonce, uint64(numTxs))
	time.Sleep(time.Second * 3)
	_ = nth.Close()

	mutSentTransactions.Lock()
	defer mutSentTransactions.Unlock()

	//no resend operation was made because all transactions were executed (nonce was incremented)
	assert.Equal(t, 1, len(sentTransactions))
	assert.Equal(t, numTxs, len(sentTransactions[0]))
}

func TestNonceTransactionsHandlerV2_SendTransactionErrors(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)

	var errSent error

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address address.Address) (*models.Account, error) {
			addressAsBech32String := address.Bech32()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &models.Account{
				AccountInfo: &models.AccountInfo{
					Nonce: atomic.LoadUint64(&currentNonce),
				},
			}, nil
		},
		SendTransactionCalled: func(tx *transaction.Transaction) (string, error) {
			return "", errSent
		},
	}

	numTxs := 1
	nth, _ := NewNonceTransactionHandlerV2(args)
	txs := createMockTransactions(testAddress, numTxs, atomic.LoadUint64(&currentNonce))

	hash, err := nth.SendTransaction(context.Background(), nil)
	require.Equal(t, interactors.ErrNilTransaction, err)
	require.Equal(t, "", hash)

	errSent = errors.New("expected error")

	hash, err = nth.SendTransaction(context.Background(), txs[0])
	require.True(t, errors.Is(err, errSent))
	require.Equal(t, "", hash)
}

func createMockTransactions(addr address.Address, numTxs int, startNonce uint64) []*transaction.Transaction {
	txs := make([]*transaction.Transaction, 0, numTxs)
	for i := 0; i < numTxs; i++ {
		tx := transaction.NewBaseTransaction(addr.Bytes(), startNonce, nil, 0, 0)
		tx.GasLimit = 50000
		tx.GasMultiplier = 100000
		tx.AddSignature([]byte("sig"))
		tx.SetChainID([]byte{3})

		// tx := &transaction.FrontendTransaction{
		// 	Value:    "1",
		// 	Receiver: addrAsBech32String,
		// 	Sender:   addrAsBech32String,
		// }

		txs = append(txs, tx)
		startNonce++
	}

	return txs
}

func TestNonceTransactionsHandlerV2_SendTransactionsWithGetNonce(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)

	mutSentTransactions := sync.Mutex{}
	numCalls := 0
	sentTransactions := make(map[int][]*transaction.Transaction)

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address address.Address) (*models.Account, error) {
			addressAsBech32String := address.Bech32()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &models.Account{
				AccountInfo: &models.AccountInfo{
					Nonce: atomic.LoadUint64(&currentNonce),
				},
			}, nil
		},
		SendTransactionCalled: func(tx *transaction.Transaction) (string, error) {
			mutSentTransactions.Lock()
			defer mutSentTransactions.Unlock()

			sentTransactions[numCalls] = []*transaction.Transaction{tx}
			numCalls++

			return "", nil
		},
	}

	numTxs := 5
	nth, _ := NewNonceTransactionHandlerV2(args)
	txs := createMockTransactionsWithGetNonce(t, testAddress, 5, nth)
	for i := 0; i < numTxs; i++ {
		_, err := nth.SendTransaction(context.Background(), txs[i])
		require.Nil(t, err)
	}

	atomic.AddUint64(&currentNonce, uint64(numTxs))
	time.Sleep(time.Second * 3)
	_ = nth.Close()

	mutSentTransactions.Lock()
	defer mutSentTransactions.Unlock()

	//no resend operation was made because all transactions were executed (nonce was incremented)
	assert.Equal(t, numTxs, len(sentTransactions))
	assert.Equal(t, 1, len(sentTransactions[0]))
}

func TestNonceTransactionsHandlerV2_SendDuplicateTransactions(t *testing.T) {
	initialNonce := uint64(664)
	currentNonce := initialNonce

	numCalls := 0

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address address.Address) (*models.Account, error) {
			addressAsBech32String := address.Bech32()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &models.Account{
				AccountInfo: &models.AccountInfo{
					Nonce: atomic.LoadUint64(&currentNonce),
				},
			}, nil
		},
		SendTransactionCalled: func(tx *transaction.Transaction) (string, error) {
			require.LessOrEqual(t, numCalls, 1)
			atomic.AddUint64(&currentNonce, 1)
			return "", nil
		},
	}
	nth, _ := NewNonceTransactionHandlerV2(args)

	tx := transaction.NewBaseTransaction(testAddress.Bytes(), 0, nil, 0, 0)
	tx.GasLimit = 50000
	tx.GasMultiplier = 100000
	tx.AddSignature([]byte("sig"))
	tx.SetChainID([]byte{3})

	// tx := &transaction.Transaction{
	// 	RawData: &transaction.Transaction_Raw{
	// 		ChainID: 3,
	// 		Version: 1,
	// 		Data:    nil,
	// 	},
	// 	Value:         "1",
	// 	Receiver:      testAddressAsBech32String,
	// 	Sender:        testAddressAsBech32String,
	// 	GasMultiplier: 100000,
	// 	GasLimit:      50000,
	// }
	err := nth.ApplyNonceAndGasPrice(context.Background(), testAddress, tx)
	require.Nil(t, err)

	_, err = nth.SendTransaction(context.Background(), tx)
	require.Nil(t, err)
	acc, _ := nth.getOrCreateAddressNonceHandler(testAddress)
	accWithPrivateAccess, ok := acc.(*addressNonceHandler)
	require.True(t, ok)

	// after sending first tx, nonce shall increase
	require.Equal(t, atomic.LoadUint64(&currentNonce), accWithPrivateAccess.computedNonce+1)

	// trying to apply nonce for the same tx, NonceTransactionHandler shall return ErrTxAlreadySent
	// and computedNonce shall not increase
	tx.RawData.Nonce = initialNonce
	err = nth.ApplyNonceAndGasPrice(context.Background(), testAddress, tx)
	require.Equal(t, interactors.ErrTxWithSameNonceAndGasPriceAlreadySent, err)
	require.Equal(t, initialNonce, tx.RawData.Nonce)
	require.Equal(t, currentNonce, accWithPrivateAccess.computedNonce+1)
}

func createMockTransactionsWithGetNonce(
	tb testing.TB,
	addr address.Address,
	numTxs int,
	nth *nonceTransactionsHandlerV2,
) []*transaction.Transaction {
	txs := make([]*transaction.Transaction, 0, numTxs)
	for i := 0; i < numTxs; i++ {
		tx := &transaction.Transaction{}
		err := nth.ApplyNonceAndGasPrice(context.Background(), addr, tx)
		require.Nil(tb, err)

		//tx.Value = "1"
		//tx.Receiver = addrAsBech32String
		tx.RawData.Sender = addr.Bytes()
		tx.GasLimit = 50000
		tx.RawData.Data = nil
		tx.Signature = [][]byte{[]byte("sig")}
		tx.SetChainID([]byte{3})
		tx.RawData.Version = 1

		txs = append(txs, tx)
	}

	return txs
}

func TestNonceTransactionsHandlerV2_ForceNonceReFetch(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address address.Address) (*models.Account, error) {
			addressAsBech32String := address.Bech32()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &models.Account{
				AccountInfo: &models.AccountInfo{
					Nonce: atomic.LoadUint64(&currentNonce),
				},
			}, nil
		},
	}

	nth, _ := NewNonceTransactionHandlerV2(args)
	tx := &transaction.Transaction{}
	_ = nth.ApplyNonceAndGasPrice(context.Background(), testAddress, tx)
	_ = nth.ApplyNonceAndGasPrice(context.Background(), testAddress, tx)
	err := nth.ApplyNonceAndGasPrice(context.Background(), testAddress, tx)
	require.Nil(t, err)
	assert.Equal(t, atomic.LoadUint64(&currentNonce)+2, tx.RawData.Nonce)

	err = nth.DropTransactions(nil)
	assert.Equal(t, interactors.ErrNilAddress, err)

	err = nth.DropTransactions(testAddress)
	assert.Nil(t, err)

	err = nth.ApplyNonceAndGasPrice(context.Background(), testAddress, tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, atomic.LoadUint64(&currentNonce), tx.RawData.Nonce)
}

func createMockArgsNonceTransactionsHandlerV2() ArgsNonceTransactionsHandlerV2 {
	return ArgsNonceTransactionsHandlerV2{
		Proxy:            &testsCommon.ProxyStub{},
		IntervalToResend: time.Second * 2,
	}
}
