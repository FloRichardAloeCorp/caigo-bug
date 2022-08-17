package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/types"
)

const (
	name                      = "dev"
	compiledOZAccount         = "./OZAccount.json"
	compiledERC20Contract     = "./erc20_custom_compiled.json"
	maxPoll               int = 15
	pollInterval          int = 5
)

func main() {
	// init starknet gateway client
	gw := gateway.NewProvider(gateway.WithChain(name))
	erc20Response, err := gw.Deploy(context.Background(), compiledERC20Contract, types.DeployRequest{
		Type:                gateway.DEPLOY,
		ContractAddressSalt: "0x30a9f5025c1321d78b4027cf28d02014a0319eddcffc7e63319e3ac0b1f0bb1",
		ConstructorCalldata: []string{
			caigo.HexToBN("0x2e4e4e1f6f85efeeed1d545e7b840387f2f79eec6c05d7569959722ebf7c00a").String(), // owner
			"2000", // initial supply
			"0",    // Uint256 additionnal parameter
		},
	})
	if err != nil {
		fmt.Println("can't deploy erc20 contract :", err)
		os.Exit(1)
	}

	// poll until the desired transaction status
	if err := waitForTransaction(gw, erc20Response); err != nil {
		fmt.Println("ERC20 deployement transaction failure :", err)
		os.Exit(1)
	}
}

// Utils function to wait for transaction to be accepted on L2 and print tx status
func waitForTransaction(gw *gateway.GatewayProvider, txResp types.AddTxResponse) error {
	acceptedOnL2 := false
	var receipt *types.TransactionReceipt
	var err error
	fmt.Println("Polling until transaction is accepted on L2...")
	for !acceptedOnL2 {
		_, receipt, err = gw.PollTx(context.Background(), txResp.TransactionHash, types.ACCEPTED_ON_L2, pollInterval, maxPoll)
		if err != nil {
			fmt.Println(receipt.Status, receipt.StatusData)
			return fmt.Errorf("Transaction Failure (%s) : can't poll to desired status : %s", txResp.TransactionHash, err.Error())
		}
		fmt.Println("Current status : ", receipt.Status)
		if receipt.Status == types.ACCEPTED_ON_L2.String() {
			acceptedOnL2 = true
		}
	}
	//fmt.Println("Transaction accepted on L2. Transaction hash : ", receipt.TransactionHash)
	return nil
}
