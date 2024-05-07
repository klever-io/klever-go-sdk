package provider

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/models/proto"
	"github.com/klever-io/klever-go-sdk/provider/utils"
)

type scMetadata struct {
	Payable     bool
	PayableBySC bool
	Upgradeable bool
	Readable    bool
}

// ToBytes converts the metadata to bytes
func (metadata *scMetadata) ToBytes() []byte {
	bytes := make([]byte, utils.LengthOfCodeMetadata)

	if metadata.Upgradeable {
		bytes[0] |= utils.MetadataUpgradeable
	}
	if metadata.Readable {
		bytes[0] |= utils.MetadataReadable
	}
	if metadata.Payable {
		bytes[1] |= utils.MetadataPayable
	}
	if metadata.PayableBySC {
		bytes[1] |= utils.MetadataPayableBySC
	}

	return bytes
}

// `payable`, `payableBySC`, `upgradeable`, `readable` are the metadata to the smart contract.
// `wasm` is a string with the path to your compiled wasm
// Pass an empty string to `vmType`, only pass a customized value if you are very
// sure of what you are doing.
// If your contract requires it, `arguments` are the parameters needed to deploy
// your smart contract.
// Checks the docs to see how to pass arguments with their correct types
func (kc *kleverChain) DeploySmartContract(
	base *models.BaseTX,
	wasmPath string,
	payable, payableBySC, upgradeable, readable bool,
	vmType string,
	arguments ...string,
) (*proto.Transaction, error) {
	if wasmPath == "" {
		return nil, fmt.Errorf("invalid file path provided: %s", wasmPath)
	}

	wasm, err := os.ReadFile(wasmPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	metadata := scMetadata{
		Payable:     payable,
		PayableBySC: payableBySC,
		Upgradeable: upgradeable,
		Readable:    readable,
	}

	metadataHex := fmt.Sprintf("%04X", metadata.ToBytes())

	if len(vmType) == 0 {
		vmType = utils.DefaultVMType
	}

	message := []string{fmt.Sprintf("%s@%s@%s", hex.EncodeToString(wasm), vmType, metadataHex)}

	base.Message = append(base.Message, message...)

	return kc.HandleSmartContracts(
		base,
		int32(utils.SmartContractDeploy),
		"",
		map[string]int64{},
		arguments...,
	)
}

// `scAddress` is the address of the smart contract that you want to invoke
// `functionToCall` is the smart contract function that you want to call
// `callvalue` are the tokens that you will send to the smart contract if the
// called function requires it, the map key is the token ticker and the value
// is the quantity of the token to be sent.
// If your contract requires it, `arguments` are the parameters needed to deploy
// your smart contract.
// `arguments` are the parameters of the function to be invoked, if it requires it
// Checks the docs to see how to pass arguments with their correct types
func (kc *kleverChain) InvokeSmartContract(
	base *models.BaseTX,
	scAddress string,
	functionToCall string,
	callValue map[string]int64,
	arguments ...string,
) (*proto.Transaction, error) {
	_, err := address.NewAddress(scAddress)

	if err != nil {
		return nil, err
	}

	base.Message = append(base.Message, functionToCall)

	return kc.HandleSmartContracts(
		base,
		int32(utils.SmartContractInvoke),
		scAddress,
		callValue,
		arguments...,
	)
}

func (kc *kleverChain) HandleSmartContracts(
	base *models.BaseTX,
	scType int32,
	scAddress string,
	callValue map[string]int64,
	arguments ...string,
) (*proto.Transaction, error) {
	contract := models.SmartContractRequest{
		SCType:    int32(scType),
		CallValue: callValue,
	}

	if len(scAddress) != 0 {
		contract.Address = scAddress
	}

	argsParsed, err := EncodeInput(arguments)

	if err != nil {
		return nil, err
	}

	base.Message[0] += argsParsed

	data, err := kc.buildRequest(
		proto.TXContract_SmartContractType,
		base,
		[]interface{}{contract},
	)

	if err != nil {
		return nil, err
	}

	return kc.PrepareTransaction(data)
}
