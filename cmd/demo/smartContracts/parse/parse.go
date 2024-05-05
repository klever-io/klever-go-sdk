package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/klever-io/klever-go-sdk/cmd/demo"
	"github.com/klever-io/klever-go-sdk/provider"
)

type HexOutput struct {
	Data struct {
		Data string `json:"data"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

type ReturnData struct {
	ReturnData []string `json:"returnData"`
	ReturnCode string   `json:"returnCode"`
}

type QueryOutput struct {
	Data struct {
		Data ReturnData `json:"data"`
	} `json:"data"`
}

func scReq(endpoint, funcName, scAddress string, args ...string) ([]byte, error) {
	jsonArgs, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("\n\nerror marshaling arguments: %v\n\n", err)
	}

	var argsString string
	if len(args) == 0 {
		argsString = `[]`
	} else {
		argsString = string(jsonArgs)
	}

	bodyTemplate := `{
		"scAddress":"%s",
		"funcName":"%s",
		"args":%s
	}`

	bodyStr := fmt.Sprintf(bodyTemplate, scAddress, funcName, argsString)

	bodyReader := strings.NewReader(bodyStr)

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://node.testnet.klever.finance/vm/%s", endpoint),
		bodyReader,
	)
	if err != nil {
		return nil, fmt.Errorf("\n\nerror creating http request: %w\n\n", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("\n\nerror executing http request: %w\n\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("\n\nerror reading http response body: %w\n\n", err)
	}

	return body, nil
}

func parseHex(hexBytes []byte, abiPath, funcName string, kc provider.KleverChain) (interface{}, error) {
	var h HexOutput

	if err := json.Unmarshal(hexBytes, &h); err != nil {
		return nil, err
	}

	if h.Error != "" {
		return nil, fmt.Errorf("\n\nerror requesting kleverchain vm: %s\n\n", h.Error)
	}

	parser := kc.NewScOutputParser()

	abi, err := os.Open(abiPath)
	if err != nil {
		return nil, err
	}

	err = parser.LoadAbi(abi)
	if err != nil {
		return nil, err
	}

	parsedValue, err := parser.ParseHex(funcName, []string{h.Data.Data})
	if err != nil {
		return nil, err
	}

	return parsedValue, nil
}

func reqAndParseHex(abiPath, scAddress, endpoint string, kc provider.KleverChain) {
	bytes, err := scReq("hex", endpoint, scAddress)
	if err != nil {
		panic(err)
	}

	parsedValue, err := parseHex(
		bytes,
		abiPath,
		endpoint,
		kc,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\nParsed %s hex output:\n %+v", endpoint, parsedValue)
}

func reqAndParseQuery(abiPath, scAddress, endpoint string, kc provider.KleverChain, args ...string) {
	var hexArgs []string
	for _, arg := range args {
		hexArgs = append(hexArgs, hex.EncodeToString([]byte(arg)))
	}

	bytes, err := scReq("query", endpoint, scAddress, hexArgs...)
	if err != nil {
		panic(err)
	}

	var q QueryOutput

	if err := json.Unmarshal(bytes, &q); err != nil {
		panic(err)
	}

	if q.Data.Data.ReturnCode != "Ok" {
		panic(fmt.Errorf("\n\nvm return code isn't Ok\n\n"))
	}

	parser := kc.NewScOutputParser()

	abi, err := os.Open(abiPath)
	if err != nil {
		panic(err)
	}

	err = parser.LoadAbi(abi)
	if err != nil {
		fmt.Printf("Error %v", err)
		panic(err)
	}

	parsedValue, err := parser.ParseQuery(endpoint, q.Data.Data.ReturnData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\nParsed %s query output:\n %+v", endpoint, parsedValue)
}

func main() {
	_, _, kc, err := demo.InitWallets() // needs the pem file of your testnet wallet on the root of the project
	if err != nil {
		panic(err)
	}

	const ExampleScAddress string = "klv1qqqqqqqqqqqqqpgqz0ce8rdktkap33cnmup545543pxy4f67d20qpkcuus"
	const LotteryScAddress string = "klv1qqqqqqqqqqqqqpgqstur8rugf2k6dulnwl48frxwfgu6a7yyhtxsqf7j4f"
	const ExampleAbiPath string = "./cmd/demo/smartContracts/decode/example.abi.json"
	const LotteryAbiPath string = "./cmd/demo/smartContracts/scFiles/lottery-kda.abi.json"

	// nested big int list
	reqAndParseHex(ExampleAbiPath, ExampleScAddress, "list_list_list_big_int", kc)

	// nested tokens identifiers list
	reqAndParseHex(ExampleAbiPath, ExampleScAddress, "list_list_list_token", kc)

	// struct with multiple fields
	reqAndParseHex(ExampleAbiPath, ExampleScAddress, "struct_test", kc)

	// list of structs
	reqAndParseQuery(LotteryAbiPath, LotteryScAddress, "getWinnersInfo", kc, "SCLotteryDemo")

	// variadic list, multiple values, of int64 (fix size)
	reqAndParseQuery(ExampleAbiPath, ExampleScAddress, "variadic_i64", kc)

	// variadic list, multiple values, of big int (dynamic size)
	reqAndParseQuery(ExampleAbiPath, ExampleScAddress, "variadic_bign", kc)

	// multiple values and types output, is like a managed tuple
	reqAndParseQuery(ExampleAbiPath, ExampleScAddress, "multi_value_nested_list_struct", kc)
}
