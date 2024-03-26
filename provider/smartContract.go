package provider

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/models/proto"
)

const lengthOfCodeMetadata = 2

// Const group for the first byte of the metadata
const (
	// MetadataUpgradeable is the bit for upgradable flag
	MetadataUpgradeable = 1
	// MetadataReadable is the bit for readable flag
	MetadataReadable = 4
)

// Const group for the second byte of the metadata
const (
	// MetadataPayable is the bit for payable flag
	MetadataPayable = 2
	// MetadataPayableBySC is the bit for payable flag
	MetadataPayableBySC = 4
)

const defaultVMType = "0500"

type SCType int32

const (
	SmartContractInvoke SCType = iota
	SmartContractDeploy
)

type scMetadata struct {
	Payable     bool
	PayableBySC bool
	Upgradeable bool
	Readable    bool
}

// ToBytes converts the metadata to bytes
func (metadata *scMetadata) ToBytes() []byte {
	bytes := make([]byte, lengthOfCodeMetadata)

	if metadata.Upgradeable {
		bytes[0] |= MetadataUpgradeable
	}
	if metadata.Readable {
		bytes[0] |= MetadataReadable
	}
	if metadata.Payable {
		bytes[1] |= MetadataPayable
	}
	if metadata.PayableBySC {
		bytes[1] |= MetadataPayableBySC
	}

	return bytes
}

func convertArguments(args []string) (string, error) {
	var result string
	for _, arg := range args {
		// split the argument into key and value
		kv := strings.SplitN(arg, ":", 2)

		// return hex encoded if has 1 part
		if len(kv) == 1 {
			result += "@" + hex.EncodeToString([]byte(kv[0]))
			continue
		}

		if len(kv) != 2 {
			return "", fmt.Errorf("invalid argument: %s", arg)
		}

		isOption := false
		// check if it is an option argument
		if strings.HasPrefix(kv[0], "option") {
			isOption = true
			kv[0] = kv[0][6:]
		}

		var value string

		// check type of value
		switch kv[0] {
		case "bi", "BI", "n", "N": // BigNumber
			var v *big.Int
			var ok bool
			// check 0x
			if strings.HasPrefix(kv[1], "0x") {
				// remove 0x and convert from hex string
				v, ok = new(big.Int).SetString(kv[1][2:], 16)
			} else {
				// convert from int string
				v, ok = new(big.Int).SetString(kv[1], 10)
			}
			if !ok {
				return "", fmt.Errorf("invalid value: %s", kv[1])
			}
			value = fmt.Sprintf("%X", v)
			// check padding
			if len(value)%2 != 0 {
				value = "0" + value
			}
		case "i", "I", "i64", "I64": // int64
			// string to int64
			v, err := strconv.ParseInt(kv[1], 10, 64)
			if err != nil {
				return "", fmt.Errorf("invalid value: %w", err)
			}
			value = fmt.Sprintf("%016X", v)
		case "u", "U", "u64", "U64": // uint64
			// string to uint64
			v, err := strconv.ParseUint(kv[1], 10, 64)
			if err != nil {
				return "", fmt.Errorf("invalid value: %w", err)
			}
			value = fmt.Sprintf("%016X", v)
		case "i32", "I32": // int32
			// string to int32
			str, err := strconv.ParseInt(kv[1], 10, 32)
			if err != nil {
				return "", fmt.Errorf("invalid value: %w", err)
			}
			value = fmt.Sprintf("%08X", str)
		case "u32", "U32": // uint32
			// string to uint32
			str, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil {
				return "", fmt.Errorf("invalid value: %w", err)
			}
			value = fmt.Sprintf("%08X", str)
		case "s", "S": // String
			// convert to string hex
			value = hex.EncodeToString([]byte(kv[1]))
		case "x", "X": // Hex Value
			// remove 0x if exists
			kv[1] = strings.TrimPrefix(kv[1], "0x")
			// validate hex string
			_, err := hex.DecodeString(kv[1])
			if err != nil {
				return "", fmt.Errorf("invalid hex value: %s", kv[1])
			}
			value = kv[1]
		case "a", "A":
			// validate address
			address, err := address.NewAddress(kv[1])
			if err != nil {
				return "", fmt.Errorf("invalid address: %s", kv[1])
			}
			value = address.Hex()
		case "0", "e", "E": // empty param
		default:
			return "", fmt.Errorf("invalid type: %s", kv[0])
		}

		if isOption {
			// append option
			value = "01" + value
		}

		// append param
		result += "@" + value
	}

	return result, nil
}

// `payable`, `payableBySC`, `upgradeable`, `readable` are the metadata to the smart contract.
// `wasm` is a string witj the path to your compiled wasm
// Pass an empty string to `vmType`, only pass a customized value if you are very
// sure of what you are doing.
// If your contract requires it, `arguments` are the parameters needed to deploy
// your smart contract.
// If the argument is a ManagedBuffer pass as a simple string, e.g "SCName".
// If the argument isn't a ManagedBuffer pass a string with the type of argument following
// by a colon and then the value of the argument, e.g "u64:123456".
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
		vmType = defaultVMType
	}

	message := []string{fmt.Sprintf("%s@%s@%s", hex.EncodeToString(wasm), vmType, metadataHex)}

	base.Message = append(base.Message, message...)

	return kc.HandleSmartContracts(
		base,
		int32(SmartContractDeploy),
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
// If the argument is a ManagedBuffer pass as a simple string, e.g "SCName".
// If the argument isn't a ManagedBuffer pass a string with the type of argument following
// by a colon and then the value of the argument, e.g "u64:123456".
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
		int32(SmartContractInvoke),
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

	argsParsed, err := convertArguments(arguments)

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
