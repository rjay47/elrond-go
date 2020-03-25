package transaction_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/coordinator"
	"github.com/ElrondNetwork/elrond-go/process/mock"
	txproc "github.com/ElrondNetwork/elrond-go/process/transaction"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-go/vm/factory"
	"github.com/stretchr/testify/assert"
)

func generateRandomByteSlice(size int) []byte {
	buff := make([]byte, size)
	_, _ = rand.Reader.Read(buff)

	return buff
}

func feeHandlerMock() *mock.FeeHandlerStub {
	return &mock.FeeHandlerStub{
		CheckValidityTxValuesCalled: func(tx process.TransactionWithFeeHandler) error {
			return nil
		},
		ComputeFeeCalled: func(tx process.TransactionWithFeeHandler) *big.Int {
			return big.NewInt(0)
		},
	}
}

func createAccountStub(sndAddr, rcvAddr []byte,
	acntSrc, acntDst *state.Account,
) *mock.AccountsStub {
	accounts := mock.AccountsStub{}

	accounts.GetAccountWithJournalCalled = func(addressContainer state.AddressContainer) (state.AccountHandler, error) {
		if bytes.Equal(addressContainer.Bytes(), sndAddr) {
			return acntSrc, nil
		}

		if bytes.Equal(addressContainer.Bytes(), rcvAddr) {
			return acntDst, nil
		}

		return nil, errors.New("failure")
	}

	return &accounts
}

func createTxProcessor() txproc.TxProcessor {
	txProc, _ := txproc.NewTxProcessor(
		&mock.AccountsStub{},
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	return txProc
}

//------- NewTxProcessor

func TestNewTxProcessor_NilAccountsShouldErr(t *testing.T) {
	t.Parallel()

	txProc, err := txproc.NewTxProcessor(
		nil,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	assert.Equal(t, process.ErrNilAccountsAdapter, err)
	assert.Nil(t, txProc)
}

func TestNewTxProcessor_NilHasherShouldErr(t *testing.T) {
	t.Parallel()

	txProc, err := txproc.NewTxProcessor(
		&mock.AccountsStub{},
		nil,
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	assert.Equal(t, process.ErrNilHasher, err)
	assert.Nil(t, txProc)
}

func TestNewTxProcessor_NilAddressConverterMockShouldErr(t *testing.T) {
	t.Parallel()

	txProc, err := txproc.NewTxProcessor(
		&mock.AccountsStub{},
		mock.HasherMock{},
		nil,
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	assert.Equal(t, process.ErrNilAddressConverter, err)
	assert.Nil(t, txProc)
}

func TestNewTxProcessor_NilMarshalizerMockShouldErr(t *testing.T) {
	t.Parallel()

	txProc, err := txproc.NewTxProcessor(
		&mock.AccountsStub{},
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		nil,
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	assert.Equal(t, process.ErrNilMarshalizer, err)
	assert.Nil(t, txProc)
}

func TestNewTxProcessor_NilShardCoordinatorMockShouldErr(t *testing.T) {
	t.Parallel()

	txProc, err := txproc.NewTxProcessor(
		&mock.AccountsStub{},
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		nil,
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	assert.Equal(t, process.ErrNilShardCoordinator, err)
	assert.Nil(t, txProc)
}

func TestNewTxProcessor_NilSCProcessorShouldErr(t *testing.T) {
	t.Parallel()

	txProc, err := txproc.NewTxProcessor(
		&mock.AccountsStub{},
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		nil,
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	assert.Equal(t, process.ErrNilSmartContractProcessor, err)
	assert.Nil(t, txProc)
}

func TestNewTxProcessor_NilTxFeeHandlerShouldErr(t *testing.T) {
	t.Parallel()

	txProc, err := txproc.NewTxProcessor(
		&mock.AccountsStub{},
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		nil,
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	assert.Equal(t, process.ErrNilUnsignedTxHandler, err)
	assert.Nil(t, txProc)
}

func TestNewTxProcessor_OkValsShouldWork(t *testing.T) {
	t.Parallel()

	txProc, err := txproc.NewTxProcessor(
		&mock.AccountsStub{},
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	assert.Nil(t, err)
	assert.NotNil(t, txProc)
}

//------- getAddresses

func TestTxProcessor_GetAddressErrAddressConvShouldErr(t *testing.T) {
	t.Parallel()

	addressConv := &mock.AddressConverterMock{}

	execTx, _ := txproc.NewTxProcessor(
		&mock.AccountsStub{},
		mock.HasherMock{},
		addressConv,
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	addressConv.Fail = true

	tx := transaction.Transaction{}

	_, _, err := execTx.GetAddresses(&tx)
	assert.NotNil(t, err)
}

func TestTxProcessor_GetAddressOkValsShouldWork(t *testing.T) {
	t.Parallel()

	execTx := *createTxProcessor()

	tx := transaction.Transaction{}
	tx.RcvAddr = []byte{65, 66, 67}
	tx.SndAddr = []byte{32, 33, 34}

	adrSnd, adrRcv, err := execTx.GetAddresses(&tx)
	assert.Nil(t, err)
	assert.Equal(t, []byte{65, 66, 67}, adrRcv.Bytes())
	assert.Equal(t, []byte{32, 33, 34}, adrSnd.Bytes())
}

//------- getAccounts

func TestTxProcessor_GetAccountsShouldErrNilAddressContainer(t *testing.T) {
	t.Parallel()

	accounts := createAccountStub(nil, nil, nil, nil)

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	adr1 := mock.NewAddressMock([]byte{65})
	adr2 := mock.NewAddressMock([]byte{67})

	_, _, err := execTx.GetAccounts(nil, adr2)
	assert.Equal(t, process.ErrNilAddressContainer, err)

	_, _, err = execTx.GetAccounts(adr1, nil)
	assert.Equal(t, process.ErrNilAddressContainer, err)
}

func TestTxProcessor_GetAccountsMalfunctionAccountsShouldErr(t *testing.T) {
	t.Parallel()

	accounts := createAccountStub(nil, nil, nil, nil)

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	adr1 := mock.NewAddressMock([]byte{65})
	adr2 := mock.NewAddressMock([]byte{67})

	_, _, err := execTx.GetAccounts(adr1, adr2)
	assert.NotNil(t, err)
}

func TestTxProcessor_GetAccountsOkValsSrcShouldWork(t *testing.T) {
	t.Parallel()

	accounts := mock.AccountsStub{}

	adr1 := mock.NewAddressMock([]byte{65})
	adr2 := mock.NewAddressMock([]byte{67})

	acnt1, _ := state.NewAccount(adr1, &mock.AccountTrackerStub{})
	acnt2, _ := state.NewAccount(adr2, &mock.AccountTrackerStub{})

	accounts.GetAccountWithJournalCalled = func(addressContainer state.AddressContainer) (state.AccountHandler, error) {
		if addressContainer == adr1 {
			return acnt1, nil
		}

		if addressContainer == adr2 {
			return nil, errors.New("failure on destination")
		}

		return nil, errors.New("failure")
	}

	shardCoordinator := mock.NewOneShardCoordinatorMock()

	execTx, _ := txproc.NewTxProcessor(
		&accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		shardCoordinator,
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	shardCoordinator.ComputeIdCalled = func(container state.AddressContainer) uint32 {
		if bytes.Equal(container.Bytes(), adr2.Bytes()) {
			return 1
		}

		return 0
	}

	a1, a2, err := execTx.GetAccounts(adr1, adr2)

	assert.Nil(t, err)
	assert.Equal(t, acnt1, a1)
	assert.NotEqual(t, acnt2, a2)
	assert.Nil(t, a2)
}

func TestTxProcessor_GetAccountsOkValsDsthouldWork(t *testing.T) {
	t.Parallel()

	accounts := mock.AccountsStub{}

	adr1 := mock.NewAddressMock([]byte{65})
	adr2 := mock.NewAddressMock([]byte{67})

	acnt1, _ := state.NewAccount(adr1, &mock.AccountTrackerStub{})
	acnt2, _ := state.NewAccount(adr2, &mock.AccountTrackerStub{})

	accounts.GetAccountWithJournalCalled = func(addressContainer state.AddressContainer) (state.AccountHandler, error) {
		if addressContainer == adr1 {
			return nil, errors.New("failure on source")
		}

		if addressContainer == adr2 {
			return acnt2, nil
		}

		return nil, errors.New("failure")
	}

	shardCoordinator := mock.NewOneShardCoordinatorMock()

	execTx, _ := txproc.NewTxProcessor(
		&accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		shardCoordinator,
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	shardCoordinator.ComputeIdCalled = func(container state.AddressContainer) uint32 {
		if bytes.Equal(container.Bytes(), adr1.Bytes()) {
			return 1
		}

		return 0
	}

	a1, a2, err := execTx.GetAccounts(adr1, adr2)
	assert.Nil(t, err)
	assert.NotEqual(t, acnt1, a1)
	assert.Nil(t, a1)
	assert.Equal(t, acnt2, a2)
}

func TestTxProcessor_GetAccountsOkValsShouldWork(t *testing.T) {
	t.Parallel()

	adr1 := mock.NewAddressMock([]byte{65})
	adr2 := mock.NewAddressMock([]byte{67})

	acnt1, _ := state.NewAccount(adr1, &mock.AccountTrackerStub{})
	acnt2, _ := state.NewAccount(adr2, &mock.AccountTrackerStub{})

	accounts := createAccountStub(adr1.Bytes(), adr2.Bytes(), acnt1, acnt2)

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	a1, a2, err := execTx.GetAccounts(adr1, adr2)
	assert.Nil(t, err)
	assert.Equal(t, acnt1, a1)
	assert.Equal(t, acnt2, a2)
}

func TestTxProcessor_GetSameAccountShouldWork(t *testing.T) {
	t.Parallel()

	adr1 := mock.NewAddressMock([]byte{65})
	adr2 := mock.NewAddressMock([]byte{65})

	acnt1, _ := state.NewAccount(adr1, &mock.AccountTrackerStub{})
	acnt2, _ := state.NewAccount(adr2, &mock.AccountTrackerStub{})

	accounts := createAccountStub(adr1.Bytes(), adr2.Bytes(), acnt1, acnt2)

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	a1, a2, err := execTx.GetAccounts(adr1, adr1)
	assert.Nil(t, err)
	assert.True(t, a1 == a2)
}

//------- checkTxValues

func TestTxProcessor_CheckTxValuesHigherNonceShouldErr(t *testing.T) {
	t.Parallel()

	adr1 := mock.NewAddressMock([]byte{65})
	acnt1, err := state.NewAccount(adr1, &mock.AccountTrackerStub{})
	assert.Nil(t, err)

	execTx := *createTxProcessor()

	acnt1.Nonce = 6

	err = execTx.CheckTxValues(&transaction.Transaction{Nonce: 7}, acnt1)
	assert.Equal(t, process.ErrHigherNonceInTransaction, err)
}

func TestTxProcessor_CheckTxValuesLowerNonceShouldErr(t *testing.T) {
	t.Parallel()

	adr1 := mock.NewAddressMock([]byte{65})
	acnt1, err := state.NewAccount(adr1, &mock.AccountTrackerStub{})
	assert.Nil(t, err)

	execTx := *createTxProcessor()

	acnt1.Nonce = 6

	err = execTx.CheckTxValues(&transaction.Transaction{Nonce: 5}, acnt1)
	assert.Equal(t, process.ErrLowerNonceInTransaction, err)
}

func TestTxProcessor_CheckTxValuesInsufficientFundsShouldErr(t *testing.T) {
	t.Parallel()

	adr1 := mock.NewAddressMock([]byte{65})
	acnt1, err := state.NewAccount(adr1, &mock.AccountTrackerStub{})
	assert.Nil(t, err)

	execTx := *createTxProcessor()

	acnt1.Balance = big.NewInt(67)

	err = execTx.CheckTxValues(&transaction.Transaction{Value: big.NewInt(68)}, acnt1)
	assert.Equal(t, process.ErrInsufficientFunds, err)
}

func TestTxProcessor_CheckTxValuesOkValsShouldErr(t *testing.T) {
	t.Parallel()

	adr1 := mock.NewAddressMock([]byte{65})
	acnt1, err := state.NewAccount(adr1, &mock.AccountTrackerStub{})
	assert.Nil(t, err)

	execTx := *createTxProcessor()

	acnt1.Balance = big.NewInt(67)

	err = execTx.CheckTxValues(&transaction.Transaction{Value: big.NewInt(67)}, acnt1)
	assert.Nil(t, err)
}

//------- moveBalances
func TestTxProcessor_MoveBalancesShouldNotFailWhenAcntSrcIsNotInNodeShard(t *testing.T) {
	t.Parallel()

	adrDst := mock.NewAddressMock([]byte{67})
	journalizeCalled := false
	saveAccountCalled := false
	acntDst, _ := state.NewAccount(adrDst, &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled = true
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled = true
			return nil
		},
	})

	execTx := *createTxProcessor()
	err := execTx.MoveBalances(nil, acntDst, big.NewInt(0))

	assert.True(t, journalizeCalled && saveAccountCalled)
	assert.Nil(t, err)
}

func TestTxProcessor_MoveBalancesShouldNotFailWhenAcntDstIsNotInNodeShard(t *testing.T) {
	t.Parallel()

	adrSrc := mock.NewAddressMock([]byte{65})
	journalizeCalled := false
	saveAccountCalled := false
	acntSrc, _ := state.NewAccount(adrSrc, &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled = true
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled = true
			return nil
		},
	})

	execTx := *createTxProcessor()
	err := execTx.MoveBalances(acntSrc, nil, big.NewInt(0))

	assert.True(t, journalizeCalled && saveAccountCalled)
	assert.Nil(t, err)
}

func TestTxProcessor_MoveBalancesOkValsShouldWork(t *testing.T) {
	t.Parallel()

	journalizeCalled := 0
	saveAccountCalled := 0
	tracker := &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled++
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled++
			return nil
		},
	}

	adrSrc := mock.NewAddressMock([]byte{65})
	acntSrc, err := state.NewAccount(adrSrc, tracker)
	assert.Nil(t, err)

	adrDst := mock.NewAddressMock([]byte{67})
	acntDst, err := state.NewAccount(adrDst, tracker)
	assert.Nil(t, err)

	execTx := *createTxProcessor()

	acntSrc.Balance = big.NewInt(64)
	acntDst.Balance = big.NewInt(31)
	err = execTx.MoveBalances(acntSrc, acntDst, big.NewInt(14))

	assert.Nil(t, err)
	assert.Equal(t, uint64(50), acntSrc.Balance.Uint64())
	assert.Equal(t, uint64(45), acntDst.Balance.Uint64())
	assert.Equal(t, 2, journalizeCalled)
	assert.Equal(t, 2, saveAccountCalled)
}

func TestTxProcessor_MoveBalancesToSelfOkValsShouldWork(t *testing.T) {
	t.Parallel()

	journalizeCalled := 0
	saveAccountCalled := 0
	tracker := &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled++
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled++
			return nil
		},
	}

	adrSrc := mock.NewAddressMock([]byte{65})
	acntSrc, err := state.NewAccount(adrSrc, tracker)
	assert.Nil(t, err)

	acntDst := acntSrc

	execTx := *createTxProcessor()

	acntSrc.Balance = big.NewInt(64)

	err = execTx.MoveBalances(acntSrc, acntDst, big.NewInt(1))
	assert.Nil(t, err)
	assert.Equal(t, uint64(64), acntSrc.Balance.Uint64())
	assert.Equal(t, uint64(64), acntDst.Balance.Uint64())
	assert.Equal(t, 2, journalizeCalled)
	assert.Equal(t, 2, saveAccountCalled)
}

//------- increaseNonce

func TestTxProcessor_IncreaseNonceOkValsShouldWork(t *testing.T) {
	t.Parallel()

	journalizeCalled := 0
	saveAccountCalled := 0
	tracker := &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled++
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled++
			return nil
		},
	}

	adrSrc := mock.NewAddressMock([]byte{65})
	acntSrc, err := state.NewAccount(adrSrc, tracker)
	assert.Nil(t, err)

	execTx := *createTxProcessor()

	acntSrc.Nonce = 45

	err = execTx.IncreaseNonce(acntSrc)
	assert.Nil(t, err)
	assert.Equal(t, uint64(46), acntSrc.Nonce)
	assert.Equal(t, 1, journalizeCalled)
	assert.Equal(t, 1, saveAccountCalled)
}

//------- ProcessTransaction

func TestTxProcessor_ProcessTransactionNilTxShouldErr(t *testing.T) {
	t.Parallel()

	execTx := *createTxProcessor()

	err := execTx.ProcessTransaction(nil)
	assert.Equal(t, process.ErrNilTransaction, err)
}

func TestTxProcessor_ProcessTransactionErrAddressConvShouldErr(t *testing.T) {
	t.Parallel()

	addressConv := &mock.AddressConverterMock{}

	execTx, _ := txproc.NewTxProcessor(
		&mock.AccountsStub{},
		mock.HasherMock{},
		addressConv,
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	addressConv.Fail = true

	err := execTx.ProcessTransaction(&transaction.Transaction{})
	assert.NotNil(t, err)
}

func TestTxProcessor_ProcessTransactionMalfunctionAccountsShouldErr(t *testing.T) {
	t.Parallel()

	accounts := createAccountStub(nil, nil, nil, nil)

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	tx := transaction.Transaction{}
	tx.Nonce = 1
	tx.SndAddr = []byte("SRC")
	tx.RcvAddr = []byte("DST")
	tx.Value = big.NewInt(45)

	err := execTx.ProcessTransaction(&tx)
	assert.NotNil(t, err)
}

func TestTxProcessor_ProcessCheckNotPassShouldErr(t *testing.T) {
	t.Parallel()

	//these values will trigger ErrHigherNonceInTransaction
	tx := transaction.Transaction{}
	tx.Nonce = 1
	tx.SndAddr = []byte("SRC")
	tx.RcvAddr = []byte("DST")
	tx.Value = big.NewInt(45)

	acntSrc, err := state.NewAccount(mock.NewAddressMock(tx.SndAddr), &mock.AccountTrackerStub{})
	assert.Nil(t, err)
	acntDst, err := state.NewAccount(mock.NewAddressMock(tx.RcvAddr), &mock.AccountTrackerStub{})
	assert.Nil(t, err)

	accounts := createAccountStub(tx.SndAddr, tx.RcvAddr, acntSrc, acntDst)

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	err = execTx.ProcessTransaction(&tx)
	assert.Equal(t, process.ErrHigherNonceInTransaction, err)
}

func TestTxProcessor_ProcessCheckShouldPassWhenAdrSrcIsNotInNodeShard(t *testing.T) {
	t.Parallel()

	journalizeCalled := 0
	saveAccountCalled := 0
	tracker := &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled++
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled++
			return nil
		},
	}

	shardCoordinator := mock.NewOneShardCoordinatorMock()

	tx := transaction.Transaction{}
	tx.Nonce = 1
	tx.SndAddr = []byte("SRC")
	tx.RcvAddr = []byte("DST")
	tx.Value = big.NewInt(45)

	shardCoordinator.ComputeIdCalled = func(container state.AddressContainer) uint32 {
		if bytes.Equal(container.Bytes(), tx.SndAddr) {
			return 1
		}

		return 0
	}

	acntSrc, err := state.NewAccount(mock.NewAddressMock(tx.SndAddr), tracker)
	assert.Nil(t, err)
	acntDst, err := state.NewAccount(mock.NewAddressMock(tx.RcvAddr), tracker)
	assert.Nil(t, err)

	accounts := createAccountStub(tx.SndAddr, tx.RcvAddr, acntSrc, acntDst)

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		shardCoordinator,
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	err = execTx.ProcessTransaction(&tx)
	assert.Nil(t, err)
	assert.Equal(t, 1, journalizeCalled)
	assert.Equal(t, 1, saveAccountCalled)
}

func TestTxProcessor_ProcessMoveBalancesShouldWork(t *testing.T) {
	t.Parallel()

	journalizeCalled := 0
	saveAccountCalled := 0
	tracker := &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled++
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled++
			return nil
		},
	}

	tx := transaction.Transaction{}
	tx.Nonce = 0
	tx.SndAddr = []byte("SRC")
	tx.RcvAddr = []byte("DST")
	tx.Value = big.NewInt(0)

	acntSrc, err := state.NewAccount(mock.NewAddressMock(tx.SndAddr), tracker)
	assert.Nil(t, err)
	acntDst, err := state.NewAccount(mock.NewAddressMock(tx.RcvAddr), tracker)
	assert.Nil(t, err)

	accounts := createAccountStub(tx.SndAddr, tx.RcvAddr, acntSrc, acntDst)

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	err = execTx.ProcessTransaction(&tx)
	assert.Nil(t, err)
	assert.Equal(t, 4, journalizeCalled)
	assert.Equal(t, 4, saveAccountCalled)
}

func TestTxProcessor_ProcessMoveBalancesShouldPassWhenAdrSrcIsNotInNodeShard(t *testing.T) {
	t.Parallel()

	journalizeCalled := 0
	saveAccountCalled := 0
	tracker := &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled++
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled++
			return nil
		},
	}

	shardCoordinator := mock.NewOneShardCoordinatorMock()

	tx := transaction.Transaction{}
	tx.Nonce = 0
	tx.SndAddr = []byte("SRC")
	tx.RcvAddr = []byte("DST")
	tx.Value = big.NewInt(0)

	shardCoordinator.ComputeIdCalled = func(container state.AddressContainer) uint32 {
		if bytes.Equal(container.Bytes(), tx.SndAddr) {
			return 1
		}

		return 0
	}

	acntSrc, err := state.NewAccount(mock.NewAddressMock(tx.SndAddr), tracker)
	assert.Nil(t, err)
	acntDst, err := state.NewAccount(mock.NewAddressMock(tx.RcvAddr), tracker)
	assert.Nil(t, err)

	accounts := createAccountStub(tx.SndAddr, tx.RcvAddr, acntSrc, acntDst)

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		shardCoordinator,
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	err = execTx.ProcessTransaction(&tx)
	assert.Nil(t, err)
	assert.Equal(t, 1, journalizeCalled)
	assert.Equal(t, 1, saveAccountCalled)
}

func TestTxProcessor_ProcessIncreaseNonceShouldPassWhenAdrSrcIsNotInNodeShard(t *testing.T) {
	t.Parallel()

	journalizeCalled := 0
	saveAccountCalled := 0
	tracker := &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled++
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled++
			return nil
		},
	}

	shardCoordinator := mock.NewOneShardCoordinatorMock()

	tx := transaction.Transaction{}
	tx.Nonce = 0
	tx.SndAddr = []byte("SRC")
	tx.RcvAddr = []byte("DST")
	tx.Value = big.NewInt(0)

	shardCoordinator.ComputeIdCalled = func(container state.AddressContainer) uint32 {
		if bytes.Equal(container.Bytes(), tx.SndAddr) {
			return 1
		}

		return 0
	}

	acntSrc, err := state.NewAccount(mock.NewAddressMock(tx.SndAddr), tracker)
	assert.Nil(t, err)
	acntDst, err := state.NewAccount(mock.NewAddressMock(tx.RcvAddr), tracker)
	assert.Nil(t, err)

	accounts := createAccountStub(tx.SndAddr, tx.RcvAddr, acntSrc, acntDst)

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		shardCoordinator,
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	err = execTx.ProcessTransaction(&tx)
	assert.Nil(t, err)
	assert.Equal(t, 1, journalizeCalled)
	assert.Equal(t, 1, saveAccountCalled)
}

func TestTxProcessor_ProcessOkValsShouldWork(t *testing.T) {
	t.Parallel()

	journalizeCalled := 0
	saveAccountCalled := 0
	tracker := &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled++
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled++
			return nil
		},
	}

	tx := transaction.Transaction{}
	tx.Nonce = 4
	tx.SndAddr = []byte("SRC")
	tx.RcvAddr = []byte("DST")
	tx.Value = big.NewInt(61)

	acntSrc, err := state.NewAccount(mock.NewAddressMock(tx.SndAddr), tracker)
	assert.Nil(t, err)
	acntDst, err := state.NewAccount(mock.NewAddressMock(tx.RcvAddr), tracker)
	assert.Nil(t, err)

	acntSrc.Nonce = 4
	acntSrc.Balance = big.NewInt(90)
	acntDst.Balance = big.NewInt(10)

	accounts := createAccountStub(tx.SndAddr, tx.RcvAddr, acntSrc, acntDst)

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	err = execTx.ProcessTransaction(&tx)
	assert.Nil(t, err)
	assert.Equal(t, uint64(5), acntSrc.Nonce)
	assert.Equal(t, uint64(29), acntSrc.Balance.Uint64())
	assert.Equal(t, uint64(71), acntDst.Balance.Uint64())
	assert.Equal(t, 4, journalizeCalled)
	assert.Equal(t, 4, saveAccountCalled)
}

func TestTxProcessor_MoveBalanceWithFeesShouldWork(t *testing.T) {
	journalizeCalled := 0
	saveAccountCalled := 0
	tracker := &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled++
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled++
			return nil
		},
	}

	tx := transaction.Transaction{}
	tx.Nonce = 4
	tx.SndAddr = []byte("SRC")
	tx.RcvAddr = []byte("DST")
	tx.Value = big.NewInt(61)
	tx.GasPrice = 2
	tx.GasLimit = 2

	acntSrc, err := state.NewAccount(mock.NewAddressMock(tx.SndAddr), tracker)
	assert.Nil(t, err)
	acntDst, err := state.NewAccount(mock.NewAddressMock(tx.RcvAddr), tracker)
	assert.Nil(t, err)

	acntSrc.Nonce = 4
	acntSrc.Balance = big.NewInt(90)
	acntDst.Balance = big.NewInt(10)

	accounts := createAccountStub(tx.SndAddr, tx.RcvAddr, acntSrc, acntDst)

	txCost := big.NewInt(16)
	feeHandler := &mock.FeeHandlerStub{
		CheckValidityTxValuesCalled: func(tx process.TransactionWithFeeHandler) error {
			return nil
		},
		ComputeFeeCalled: func(tx process.TransactionWithFeeHandler) *big.Int {
			return txCost
		},
	}

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		feeHandler,
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	err = execTx.ProcessTransaction(&tx)
	assert.Nil(t, err)
	assert.Equal(t, uint64(5), acntSrc.Nonce)
	assert.Equal(t, uint64(13), acntSrc.Balance.Uint64())
	assert.Equal(t, uint64(71), acntDst.Balance.Uint64())
	assert.Equal(t, 4, journalizeCalled)
	assert.Equal(t, 4, saveAccountCalled)
}

func TestTxProcessor_ProcessTransactionScTxShouldWork(t *testing.T) {
	t.Parallel()

	journalizeCalled := 0
	saveAccountCalled := 0
	tracker := &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled++
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled++
			return nil
		},
	}

	addrConverter := &mock.AddressConverterMock{}

	tx := transaction.Transaction{}
	tx.Nonce = 0
	tx.SndAddr = []byte("SRC")
	tx.RcvAddr = generateRandomByteSlice(addrConverter.AddressLen())
	tx.Value = big.NewInt(45)
	tx.GasPrice = 1
	tx.GasLimit = 1

	acntSrc, err := state.NewAccount(mock.NewAddressMock(tx.SndAddr), tracker)
	assert.Nil(t, err)

	acntDst, err := state.NewAccount(mock.NewAddressMock(tx.RcvAddr), tracker)
	assert.Nil(t, err)

	acntSrc.Balance = big.NewInt(46)
	acntDst.SetCode([]byte{65})

	accounts := createAccountStub(tx.SndAddr, tx.RcvAddr, acntSrc, acntDst)
	scProcessorMock := &mock.SCProcessorMock{}

	wasCalled := false
	scProcessorMock.ExecuteSmartContractTransactionCalled = func(tx data.TransactionHandler, acntSrc, acntDst state.UserAccountHandler) error {
		wasCalled = true
		return nil
	}

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		scProcessorMock,
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{
			ComputeTransactionTypeCalled: func(tx data.TransactionHandler) (transactionType process.TransactionType, e error) {
				return process.SCInvoking, nil
			},
		},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	err = execTx.ProcessTransaction(&tx)
	assert.Nil(t, err)
	assert.True(t, wasCalled)
	assert.Equal(t, 0, journalizeCalled)
	assert.Equal(t, 0, saveAccountCalled)
}

func TestTxProcessor_ProcessTransactionScTxShouldReturnErrWhenExecutionFails(t *testing.T) {
	t.Parallel()

	journalizeCalled := 0
	saveAccountCalled := 0
	tracker := &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled++
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled++
			return nil
		},
	}

	addrConverter := &mock.AddressConverterMock{}

	tx := transaction.Transaction{}
	tx.Nonce = 0
	tx.SndAddr = []byte("SRC")
	tx.RcvAddr = generateRandomByteSlice(addrConverter.AddressLen())
	tx.Value = big.NewInt(45)

	acntSrc, err := state.NewAccount(mock.NewAddressMock(tx.SndAddr), tracker)
	assert.Nil(t, err)
	acntSrc.Balance = big.NewInt(45)
	acntDst, err := state.NewAccount(mock.NewAddressMock(tx.RcvAddr), tracker)
	assert.Nil(t, err)
	acntDst.SetCode([]byte{65})

	accounts := createAccountStub(tx.SndAddr, tx.RcvAddr, acntSrc, acntDst)

	scProcessorMock := &mock.SCProcessorMock{}

	wasCalled := false
	scProcessorMock.ExecuteSmartContractTransactionCalled = func(tx data.TransactionHandler, acntSrc, acntDst state.UserAccountHandler) error {
		wasCalled = true
		return process.ErrNoVM
	}

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewOneShardCoordinatorMock(),
		scProcessorMock,
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{ComputeTransactionTypeCalled: func(tx data.TransactionHandler) (transactionType process.TransactionType, e error) {
			return process.SCInvoking, nil
		}},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	err = execTx.ProcessTransaction(&tx)
	assert.Equal(t, process.ErrNoVM, err)
	assert.True(t, wasCalled)
	assert.Equal(t, 0, journalizeCalled)
	assert.Equal(t, 0, saveAccountCalled)
}

func TestTxProcessor_ProcessTransactionScTxShouldNotBeCalledWhenAdrDstIsNotInNodeShard(t *testing.T) {
	t.Parallel()

	shardCoordinator := mock.NewOneShardCoordinatorMock()

	journalizeCalled := 0
	saveAccountCalled := 0
	tracker := &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
			journalizeCalled++
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			saveAccountCalled++
			return nil
		},
	}

	addrConverter := &mock.AddressConverterMock{}

	tx := transaction.Transaction{}
	tx.Nonce = 0
	tx.SndAddr = []byte("SRC")
	tx.RcvAddr = generateRandomByteSlice(addrConverter.AddressLen())
	tx.Value = big.NewInt(45)

	shardCoordinator.ComputeIdCalled = func(container state.AddressContainer) uint32 {
		if bytes.Equal(container.Bytes(), tx.RcvAddr) {
			return 1
		}

		return 0
	}

	acntSrc, err := state.NewAccount(mock.NewAddressMock(tx.SndAddr), tracker)
	assert.Nil(t, err)
	acntSrc.Balance = big.NewInt(45)
	acntDst, err := state.NewAccount(mock.NewAddressMock(tx.RcvAddr), tracker)
	assert.Nil(t, err)
	acntDst.SetCode([]byte{65})

	accounts := createAccountStub(tx.SndAddr, tx.RcvAddr, acntSrc, acntDst)

	scProcessorMock := &mock.SCProcessorMock{}
	wasCalled := false
	scProcessorMock.ExecuteSmartContractTransactionCalled = func(tx data.TransactionHandler, acntSrc, acntDst state.UserAccountHandler) error {
		wasCalled = true
		return process.ErrNoVM
	}

	computeType, _ := coordinator.NewTxTypeHandler(
		&mock.AddressConverterMock{},
		shardCoordinator,
		accounts)

	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		shardCoordinator,
		scProcessorMock,
		&mock.FeeAccumulatorStub{},
		computeType,
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	err = execTx.ProcessTransaction(&tx)
	assert.Nil(t, err)
	assert.False(t, wasCalled)
	assert.Equal(t, 3, journalizeCalled)
	assert.Equal(t, 3, saveAccountCalled)
}

func TestTxProcessor_ProcessTxFeeIntraShard(t *testing.T) {
	t.Parallel()

	moveBalanceFee := big.NewInt(50)
	negMoveBalanceFee := big.NewInt(0).Neg(moveBalanceFee)
	execTx, _ := txproc.NewTxProcessor(
		&mock.AccountsStub{},
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewMultiShardsCoordinatorMock(2),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		&mock.FeeHandlerStub{
			ComputeFeeCalled: func(tx process.TransactionWithFeeHandler) *big.Int {
				return moveBalanceFee
			},
		},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)
	tx := &transaction.Transaction{
		RcvAddr:  []byte("aaa"),
		SndAddr:  []byte("bbb"),
		GasPrice: moveBalanceFee.Uint64(),
		GasLimit: moveBalanceFee.Uint64(),
	}

	acntSnd := &mock.UserAccountStub{AddToBalanceCalled: func(value *big.Int) error {
		assert.True(t, value.Cmp(negMoveBalanceFee) == 0)
		return nil
	}}
	acntDst := &mock.UserAccountStub{}

	cost, err := execTx.ProcessTxFee(tx, acntSnd, acntDst)
	assert.Nil(t, err)
	assert.True(t, cost.Cmp(moveBalanceFee) == 0)
}

func TestTxProcessor_ProcessTxFeeCrossShardMoveBalance(t *testing.T) {
	t.Parallel()

	moveBalanceFee := big.NewInt(50)
	negMoveBalanceFee := big.NewInt(0).Neg(moveBalanceFee)
	execTx, _ := txproc.NewTxProcessor(
		&mock.AccountsStub{},
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewMultiShardsCoordinatorMock(2),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		&mock.FeeHandlerStub{
			ComputeFeeCalled: func(tx process.TransactionWithFeeHandler) *big.Int {
				return moveBalanceFee
			},
		},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)
	tx := &transaction.Transaction{
		RcvAddr:  []byte("aaa"),
		SndAddr:  []byte("bbb"),
		GasPrice: moveBalanceFee.Uint64(),
		GasLimit: moveBalanceFee.Uint64(),
	}

	acntSnd := &mock.UserAccountStub{AddToBalanceCalled: func(value *big.Int) error {
		assert.True(t, value.Cmp(negMoveBalanceFee) == 0)
		return nil
	}}

	cost, err := execTx.ProcessTxFee(tx, acntSnd, nil)
	assert.Nil(t, err)
	assert.True(t, cost.Cmp(moveBalanceFee) == 0)

	tx = &transaction.Transaction{
		RcvAddr:  []byte("aaa"),
		SndAddr:  []byte("bbb"),
		GasPrice: moveBalanceFee.Uint64(),
		GasLimit: moveBalanceFee.Uint64(),
		Data:     []byte("data"),
	}

	cost, err = execTx.ProcessTxFee(tx, acntSnd, nil)
	assert.Nil(t, err)
	assert.True(t, cost.Cmp(moveBalanceFee) == 0)

	scAddress, _ := hex.DecodeString("000000000000000000005fed9c659422cd8429ce92f8973bba2a9fb51e0eb3a1")
	tx = &transaction.Transaction{
		RcvAddr:  scAddress,
		SndAddr:  []byte("bbb"),
		GasPrice: moveBalanceFee.Uint64(),
		GasLimit: moveBalanceFee.Uint64(),
	}

	cost, err = execTx.ProcessTxFee(tx, acntSnd, nil)
	assert.Nil(t, err)
	assert.True(t, cost.Cmp(moveBalanceFee) == 0)
}

func TestTxProcessor_ProcessTxFeeCrossShardSCCall(t *testing.T) {
	t.Parallel()

	moveBalanceFee := big.NewInt(50)
	execTx, _ := txproc.NewTxProcessor(
		&mock.AccountsStub{},
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		mock.NewMultiShardsCoordinatorMock(2),
		&mock.SCProcessorMock{},
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{},
		&mock.FeeHandlerStub{
			ComputeFeeCalled: func(tx process.TransactionWithFeeHandler) *big.Int {
				return moveBalanceFee
			},
		},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	scAddress, _ := hex.DecodeString("000000000000000000005fed9c659422cd8429ce92f8973bba2a9fb51e0eb3a1")
	tx := &transaction.Transaction{
		RcvAddr:  scAddress,
		SndAddr:  []byte("bbb"),
		GasPrice: moveBalanceFee.Uint64(),
		GasLimit: moveBalanceFee.Uint64(),
		Data:     []byte("data"),
	}

	totalCost := big.NewInt(0).Mul(big.NewInt(0).SetUint64(tx.GetGasPrice()), big.NewInt(0).SetUint64(tx.GetGasLimit()))
	negTotalCost := big.NewInt(0).Neg(totalCost)
	acntSnd := &mock.UserAccountStub{AddToBalanceCalled: func(value *big.Int) error {
		assert.True(t, value.Cmp(negTotalCost) == 0)
		return nil
	}}

	cost, err := execTx.ProcessTxFee(tx, acntSnd, nil)
	assert.Nil(t, err)
	assert.True(t, cost.Cmp(moveBalanceFee) == 0)
}

func TestTxProcessor_ProcessTransactionShouldReturnErrForInvalidMetaTx(t *testing.T) {
	t.Parallel()

	tracker := &mock.AccountTrackerStub{
		JournalizeCalled: func(entry state.JournalEntry) {
		},
		SaveAccountCalled: func(accountHandler state.AccountHandler) error {
			return nil
		},
	}

	tx := transaction.Transaction{}
	tx.Nonce = 0
	tx.SndAddr = []byte("SRC")
	tx.RcvAddr = factory.StakingSCAddress
	tx.Value = big.NewInt(45)
	tx.GasPrice = 1
	tx.GasLimit = 1

	acntSrc, err := state.NewAccount(mock.NewAddressMock(tx.SndAddr), tracker)
	assert.Nil(t, err)

	acntSrc.Balance = big.NewInt(46)
	accounts := createAccountStub(tx.SndAddr, tx.RcvAddr, acntSrc, nil)
	scProcessorMock := &mock.SCProcessorMock{}
	shardC, _ := sharding.NewMultiShardCoordinator(5, 3)
	execTx, _ := txproc.NewTxProcessor(
		accounts,
		mock.HasherMock{},
		&mock.AddressConverterMock{},
		&mock.MarshalizerMock{},
		shardC,
		scProcessorMock,
		&mock.FeeAccumulatorStub{},
		&mock.TxTypeHandlerMock{
			ComputeTransactionTypeCalled: func(tx data.TransactionHandler) (transactionType process.TransactionType, e error) {
				return process.MoveBalance, nil
			},
		},
		feeHandlerMock(),
		&mock.IntermediateTransactionHandlerMock{},
		&mock.IntermediateTransactionHandlerMock{},
	)

	err = execTx.ProcessTransaction(&tx)
	assert.Equal(t, err, process.ErrFailedTransaction)
}
