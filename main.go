package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

type Request struct {
	RequestId            [4]uint8       "json:\"requestId\""
	RequesterAddress     common.Address "json:\"requesterAddress\""
	TargetAddress        common.Address "json:\"targetAddress\""
	TargetChainDomain    *big.Int       "json:\"targetChainDomain\""
	RequestCreationBlock *big.Int       "json:\"requestCreationBlock\""
	RequestCreationTime  *big.Int       "json:\"requestCreationTime\""
	Message              []uint8        "json:\"message\""
}

func main() {
	client, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	contractAddress := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")
	parsedABI, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[{"name":"requestId","type":"bytes4"}],"name":"getRequest","outputs":[{"components":[{"name":"requestId","type":"bytes4"},{"name":"requesterAddress","type":"address"},{"name":"targetAddress","type":"address"},{"name":"targetChainDomain","type":"uint256"},{"name":"requestCreationBlock","type":"uint256"},{"name":"requestCreationTime","type":"uint256"},{"name":"message","type":"bytes"}],"name":"","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"}]`))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	// Generate request IDs
	requestIds := make([][4]byte, 3)

	for i := 0; i < 3; i++ {
		requestIds[i] = generateRequestId(uint64(i))
	}
	requests, err := getRequests(client, contractAddress, parsedABI, requestIds)
	if err != nil {
		log.Fatalf("Failed to get requests: %v", err)
	}
	log.Printf("Requests: %+v", requests)
}

func getRequests(client *rpc.Client, contractAddress common.Address, parsedABI abi.ABI, requestIds [][4]byte) ([]Request, error) {
	// Prepare batch calls
	var (
		calls    []rpc.BatchElem
		requests []Request
	)
	for i, requestId := range requestIds {
		data, err := parsedABI.Pack("getRequest", requestId)
		if err != nil {
			log.Fatalf("Failed to pack data for request %d: %v", i, err)
		}

		calls = append(calls, rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"to":   contractAddress.Hex(),
					"data": fmt.Sprintf("%x", data),
				},
				"latest",
			},
			Result: new(string),
		})
	}

	// Perform batch call
	err := client.BatchCallContext(context.Background(), calls)
	if err != nil {
		log.Fatalf("Failed to perform batch call: %v", err)
	}

	// Process responses
	for i, call := range calls {
		if call.Error != nil {
			return nil, call.Error
		}
		data := common.FromHex(*call.Result.(*string))
		unpackedData, err := parsedABI.Unpack("getRequest", data)
		if err != nil {
			log.Printf("Failed to unpack result for request %d: %v", i, err)
			return nil, err
		}
		request := unpackedData[0].(struct {
			RequestId            [4]byte        "json:\"requestId\""
			RequesterAddress     common.Address "json:\"requesterAddress\""
			TargetAddress        common.Address "json:\"targetAddress\""
			TargetChainDomain    *big.Int       "json:\"targetChainDomain\""
			RequestCreationBlock *big.Int       "json:\"requestCreationBlock\""
			RequestCreationTime  *big.Int       "json:\"requestCreationTime\""
			Message              []byte         "json:\"message\""
		})
		requests = append(requests, request)
	}
	return requests, nil
}

func generateRequestId(i uint64) [4]byte {
	data := make([]byte, 8)
	for j := 0; j < 8; j++ {
		data[j] = byte(i >> (8 * (7 - j)))
	}

	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	var id [4]byte
	copy(id[:], hash.Sum(nil)[:4])
	return id
}
