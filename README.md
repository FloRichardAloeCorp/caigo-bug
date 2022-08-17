# caigo-bug

This repository contains the example and files to reproduce the bug of this [issue](https://github.com/dontpanicdao/caigo/issues/95).

- erc20_contract.cairo : the contract that can't be deployed with caigo gateway
- OZAccount.cairo : an account used when the bug occurs. It occurs with default devnet account too.
- Intercepted payloads are:
  - deploy_caigo.json :  payload of the caigo.Deploy call
  - deploy_cli.json : payload of the deploy sent by the starknet cli.
  - to see the difference between the files : `diff deploy_cli.json deploy_caigo.json -y`
  
 # Run the program
 `go run main.go`
