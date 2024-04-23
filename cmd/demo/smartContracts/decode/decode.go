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
		return nil, fmt.Errorf("error marshaling arguments: %v", err)
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
		return nil, fmt.Errorf("error creating http request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading http response body: %w", err)
	}

	return body, nil
}

func parseHex(hexBytes []byte, abiPath, funcName string, kc provider.KleverChain) (interface{}, error) {
	var h HexOutput

	if err := json.Unmarshal(hexBytes, &h); err != nil {
		return nil, err
	}

	if h.Error != "" {
		return nil, fmt.Errorf("error requesting kleverchain vm: %s", h.Error)
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

	parsedValue, err := parser.DecodeHex(funcName, h.Data.Data)
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

func reqAndParseQuery(abiPath, scAddress, endpoint string, kc provider.KleverChain) {
	lotteryNameHex := hex.EncodeToString([]byte("SCLotteryDemo"))

	bytes, err := scReq("query", endpoint, scAddress, lotteryNameHex)
	if err != nil {
		panic(err)
	}

	var q QueryOutput

	if err := json.Unmarshal(bytes, &q); err != nil {
		panic(err)
	}

	if q.Data.Data.ReturnCode != "Ok" {
		panic(fmt.Errorf("vm return code isn't Ok"))
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

	parsedValue, err := parser.DecodeQuery(endpoint, q.Data.Data.ReturnData[0])
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\nParsed %s query output:\n %+v", endpoint, parsedValue)
}

func main() {
	_, _, kc, err := demo.InitWallets()
	if err != nil {
		panic(err)
	}

	const ExampleScAddress string = "klv1qqqqqqqqqqqqqpgqz0ce8rdktkap33cnmup545543pxy4f67d20qpkcuus"
	const LotteryScAddress string = "klv1qqqqqqqqqqqqqpgqstur8rugf2k6dulnwl48frxwfgu6a7yyhtxsqf7j4f"
	const ExampleAbiPath string = "./cmd/demo/smartContracts/decode/example.abi.json"
	const LotteryAbiPath string = "./cmd/demo/smartContracts/scFiles/lottery-kda.abi.json"

	reqAndParseHex(ExampleAbiPath, ExampleScAddress, "list_list_list_big_int", kc)
	reqAndParseHex(ExampleAbiPath, ExampleScAddress, "struct_test", kc)
	reqAndParseQuery(LotteryAbiPath, LotteryScAddress, "getWinnersInfo", kc)
}
