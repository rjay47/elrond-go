package metachain

import (
	"testing"

	arwenConfig "github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/process/economics"
	"github.com/ElrondNetwork/elrond-go/process/factory"
	"github.com/ElrondNetwork/elrond-go/process/mock"
	"github.com/ElrondNetwork/elrond-go/process/smartContract/builtInFunctions"
	"github.com/ElrondNetwork/elrond-go/process/smartContract/hooks"
	"github.com/ElrondNetwork/elrond-go/testscommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockVMAccountsArguments() hooks.ArgBlockChainHook {
	datapool := testscommon.NewPoolsHolderMock()
	arguments := hooks.ArgBlockChainHook{
		Accounts: &mock.AccountsStub{
			GetExistingAccountCalled: func(address []byte) (handler state.AccountHandler, e error) {
				return &mock.AccountWrapMock{}, nil
			},
		},
		PubkeyConv:         mock.NewPubkeyConverterMock(32),
		StorageService:     &mock.ChainStorerMock{},
		BlockChain:         &mock.BlockChainMock{},
		ShardCoordinator:   mock.NewOneShardCoordinatorMock(),
		Marshalizer:        &mock.MarshalizerMock{},
		Uint64Converter:    &mock.Uint64ByteSliceConverterMock{},
		BuiltInFunctions:   builtInFunctions.NewBuiltInFunctionContainer(),
		DataPool:           datapool,
		CompiledSCPool:     datapool.SmartContracts(),
		NilCompiledSCStore: true,
	}
	return arguments
}

func TestNewVMContainerFactory_OkValues(t *testing.T) {
	t.Parallel()

	gasSchedule := makeGasSchedule()
	vmf, err := NewVMContainerFactory(
		createMockVMAccountsArguments(),
		&economics.EconomicsData{},
		&mock.MessageSignVerifierMock{},
		gasSchedule,
		&mock.NodesConfigProviderStub{},
		&mock.HasherMock{},
		&mock.MarshalizerMock{},
		&config.SystemSmartContractsConfig{
			ESDTSystemSCConfig: config.ESDTSystemSCConfig{
				BaseIssuingCost: "100000000",
				OwnerAddress:    "aaaaaa",
			},
			GovernanceSystemSCConfig: config.GovernanceSystemSCConfig{
				ProposalCost:     "500",
				NumNodes:         100,
				MinQuorum:        50,
				MinPassThreshold: 50,
				MinVetoThreshold: 50,
			},
			StakingSystemSCConfig: config.StakingSystemSCConfig{
				GenesisNodePrice:                     "1000",
				UnJailValue:                          "10",
				MinStepValue:                         "10",
				MinStakeValue:                        "1",
				UnBondPeriod:                         1,
				AuctionEnableEpoch:                   0,
				StakeEnableEpoch:                     0,
				NumRoundsWithoutBleed:                1,
				MaximumPercentageToBleed:             1,
				BleedPercentagePerRound:              1,
				MaxNumberOfNodesForStake:             1,
				NodesToSelectInAuction:               100,
				ActivateBLSPubKeyMessageVerification: false,
			},
		},
		&mock.AccountsStub{},
		&mock.RaterMock{},
		&mock.EpochNotifierStub{},
	)

	assert.NotNil(t, vmf)
	assert.Nil(t, err)
	assert.False(t, vmf.IsInterfaceNil())
}

func TestVmContainerFactory_Create(t *testing.T) {
	t.Parallel()

	argsNewEconomicsData := economics.ArgsNewEconomicsData{
		Economics: &config.EconomicsConfig{
			GlobalSettings: config.GlobalSettings{
				GenesisTotalSupply: "2000000000000000000000",
				MinimumInflation:   0,
				YearSettings: []*config.YearSetting{
					{
						Year:             0,
						MaximumInflation: 0.01,
					},
				},
			},
			RewardsSettings: config.RewardsSettings{
				LeaderPercentage:                 0.1,
				ProtocolSustainabilityPercentage: 0.1,
				ProtocolSustainabilityAddress:    "erd1932eft30w753xyvme8d49qejgkjc09n5e49w4mwdjtm0neld797su0dlxp",
			},
			FeeSettings: config.FeeSettings{
				MaxGasLimitPerBlock:     "10000000000",
				MaxGasLimitPerMetaBlock: "10000000000",
				MinGasPrice:             "10",
				MinGasLimit:             "10",
				GasPerDataByte:          "1",
				DataLimitForBaseCalc:    "10000",
			},
		},
		PenalizedTooMuchGasEnableEpoch: 0,
		EpochNotifier:                  &mock.EpochNotifierStub{},
	}
	economicsData, _ := economics.NewEconomicsData(argsNewEconomicsData)

	vmf, err := NewVMContainerFactory(
		createMockVMAccountsArguments(),
		economicsData,
		&mock.MessageSignVerifierMock{},
		makeGasSchedule(),
		&mock.NodesConfigProviderStub{},
		&mock.HasherMock{},
		&mock.MarshalizerMock{},
		&config.SystemSmartContractsConfig{
			ESDTSystemSCConfig: config.ESDTSystemSCConfig{
				BaseIssuingCost: "100000000",
				OwnerAddress:    "aaaaaa",
			},
			GovernanceSystemSCConfig: config.GovernanceSystemSCConfig{
				ProposalCost:     "500",
				NumNodes:         100,
				MinQuorum:        50,
				MinPassThreshold: 50,
				MinVetoThreshold: 50,
			},
			StakingSystemSCConfig: config.StakingSystemSCConfig{
				GenesisNodePrice:                     "1000",
				UnJailValue:                          "100",
				MinStepValue:                         "100",
				MinStakeValue:                        "1",
				UnBondPeriod:                         1,
				AuctionEnableEpoch:                   1,
				StakeEnableEpoch:                     1,
				NumRoundsWithoutBleed:                1,
				MaximumPercentageToBleed:             1,
				BleedPercentagePerRound:              1,
				MaxNumberOfNodesForStake:             100,
				NodesToSelectInAuction:               100,
				ActivateBLSPubKeyMessageVerification: false,
			},
		},
		&mock.AccountsStub{},
		&mock.RaterMock{},
		&mock.EpochNotifierStub{},
	)
	assert.NotNil(t, vmf)
	assert.Nil(t, err)

	container, err := vmf.Create()
	require.Nil(t, err)
	require.NotNil(t, container)
	defer func() {
		_ = container.Close()
	}()

	assert.Nil(t, err)
	assert.NotNil(t, container)

	vm, err := container.Get(factory.SystemVirtualMachine)
	assert.Nil(t, err)
	assert.NotNil(t, vm)

	acc := vmf.BlockChainHookImpl()
	assert.NotNil(t, acc)
}

func makeGasSchedule() map[string]map[string]uint64 {
	gasSchedule := arwenConfig.MakeGasMapForTests()
	FillGasMapInternal(gasSchedule, 1)
	return gasSchedule
}

func FillGasMapInternal(gasMap map[string]map[string]uint64, value uint64) map[string]map[string]uint64 {
	gasMap[core.BaseOperationCost] = FillGasMapBaseOperationCosts(value)
	gasMap[core.MetaChainSystemSCsCost] = FillGasMapMetaChainSystemSCsCosts(value)

	return gasMap
}

func FillGasMapBaseOperationCosts(value uint64) map[string]uint64 {
	gasMap := make(map[string]uint64)
	gasMap["StorePerByte"] = value
	gasMap["DataCopyPerByte"] = value
	gasMap["ReleasePerByte"] = value
	gasMap["PersistPerByte"] = value
	gasMap["CompilePerByte"] = value
	gasMap["AoTPreparePerByte"] = value

	return gasMap
}

func FillGasMapMetaChainSystemSCsCosts(value uint64) map[string]uint64 {
	gasMap := make(map[string]uint64)
	gasMap["Stake"] = value
	gasMap["UnStake"] = value
	gasMap["UnBond"] = value
	gasMap["Claim"] = value
	gasMap["Get"] = value
	gasMap["ChangeRewardAddress"] = value
	gasMap["ChangeValidatorKeys"] = value
	gasMap["UnJail"] = value
	gasMap["ESDTIssue"] = value
	gasMap["ESDTOperations"] = value
	gasMap["Proposal"] = value
	gasMap["Vote"] = value
	gasMap["DelegateVote"] = value
	gasMap["RevokeVote"] = value
	gasMap["CloseProposal"] = value

	return gasMap
}
