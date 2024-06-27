package probe

import (
	"bytes"
	"chall2/exporter-sd/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
)

type RPCRequest struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      int      `json:"id"`
}

type RPCResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

var (
	ankrURL   = "https://rpc.ankr.com/eth"
	infuraURL = "https://mainnet.infura.io/v3/"
)

func checkDiffBlockNumbers(id int, blockNumberDiffGauge prometheus.Gauge) bool {

	err := godotenv.Load("config.env")
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("INFURA_API_KEY")

	infuraURL := infuraURL + apiKey

	rpcClient, err := ethclient.Dial(ankrURL)

	if err != nil {
		panic(err)
	}

	ankrBlockNumber, err := rpcClient.BlockNumber(context.Background())

	if err != nil {
		panic(err)
	}

	infuraReq := RPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_blockNumber",
		Params:  []string{},
		ID:      id,
	}

	infuraResp, err := sendRequest(infuraURL, infuraReq)
	if err != nil {
		panic(err)
	}

	infuraBlockNumber, err := strconv.ParseUint(infuraResp.Result, 0, 64)
	if err != nil {
		panic(err)
	}

	diff := ankrBlockNumber - infuraBlockNumber

	blockNumberDiffGauge.Set(float64(diff))

	return diff < 5

}

func sendRequest(url string, req RPCRequest) (*RPCResponse, error) {
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResp RPCResponse

	err = json.Unmarshal(body, &rpcResp)
	if err != nil {
		return nil, err
	}

	return &rpcResp, nil
}

func Handler(w http.ResponseWriter, r *http.Request, c *config.Config, timeoutOffset float64, params url.Values, moduleUnknownCounter prometheus.Counter) {
	fmt.Printf("%v\n", params)

	if params == nil {
		params = r.URL.Query()
	}

	id, err := strconv.Atoi(params.Get("id"))

	if err != nil {
		panic(err)
	}

	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", http.StatusBadRequest)
		return
	}

	blockNumberDiffGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "block_number_diff",
			Help: "Difference in block numbers between Ankr and Infura",
		},
	)
	successCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "success",
			Help: "Number of successful checks",
		},
	)
	failureCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "failure",
			Help: "Number of failed checks",
		},
	)

	registry := prometheus.NewRegistry()

	registry.MustRegister(blockNumberDiffGauge)
	registry.MustRegister(successCounter)
	registry.MustRegister(failureCounter)

	success := checkDiffBlockNumbers(id, blockNumberDiffGauge)

	if success {
		successCounter.Add(1)
	} else {
		failureCounter.Add(1)
	}

}

// func getTimeout(r *http.Request, module config.Module, offset float64) (timeoutSeconds float64, err error) {
// 	// If a timeout is configured via the Prometheus header, add it to the request.
// 	if v := r.Header.Get("X-Prometheus-Scrape-Timeout-Seconds"); v != "" {
// 		var err error
// 		timeoutSeconds, err = strconv.ParseFloat(v, 64)
// 		if err != nil {
// 			return 0, err
// 		}
// 	}
// 	if timeoutSeconds == 0 {
// 		timeoutSeconds = 120
// 	}

// 	var maxTimeoutSeconds = timeoutSeconds - offset
// 	if module.Timeout.Seconds() < maxTimeoutSeconds && module.Timeout.Seconds() > 0 || maxTimeoutSeconds < 0 {
// 		timeoutSeconds = module.Timeout.Seconds()
// 	} else {
// 		timeoutSeconds = maxTimeoutSeconds
// 	}

// 	return timeoutSeconds, nil
// }
