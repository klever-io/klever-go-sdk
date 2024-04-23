package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

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

func main() {
	_, _, kc, err := demo.InitWallets()
	if err != nil {
		panic(err)
	}

	const ScAddress string = "klv1qqqqqqqqqqqqqpgqz0ce8rdktkap33cnmup545543pxy4f67d20qpkcuus"

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()

		viewName := "list_list_list_big_int"

		bytes, err := scReq("hex", viewName, ScAddress)
		if err != nil {
			panic(err)
		}

		parsedValue, err := parseHex(
			bytes,
			"./cmd/demo/smartContracts/decode/example.abi.json",
			viewName,
			kc,
		)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\n\nParsed output of big int nested list:\n %+v", parsedValue)
	}()

	go func() {
		defer wg.Done()

		viewName := "struct_test"

		bytes, err := scReq("hex", viewName, ScAddress)
		if err != nil {
			panic(err)
		}

		parsedValue, err := parseHex(
			bytes,
			"./cmd/demo/smartContracts/decode/example.abi.json",
			viewName,
			kc,
		)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\n\nParsed output of custom struct type:\n %+v", parsedValue)
	}()

	go func() {
		defer wg.Done()

		viewName := "getWinnersInfo"

		lotteryNameHex := hex.EncodeToString([]byte("SCLotteryDemo"))

		bytes, err := scReq("query", viewName, "klv1qqqqqqqqqqqqqpgqstur8rugf2k6dulnwl48frxwfgu6a7yyhtxsqf7j4f", lotteryNameHex)
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

		abi, err := os.Open("./cmd/demo/smartContracts/scFiles/lottery-kda.abi.json")
		if err != nil {
			panic(err)
		}

		err = parser.LoadAbi(abi)
		if err != nil {
			fmt.Printf("Error %v", err)
			panic(err)
		}

		parsedValue, err := parser.DecodeQuery(viewName, q.Data.Data.ReturnData[0])
		if err != nil {
			panic(err)
		}

		fmt.Printf("\n\nParsed output of list of custom struct:\n %+v", parsedValue)
	}()

	wg.Wait()
}
