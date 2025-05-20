# Klever Go SDK
The Klever Go SDK provides tools and libraries for interacting with the Klever blockchain, including asset management, smart contract deployment, staking, governance, and more.

## Using the Klever Go SDK in Your Project

To use the Klever Go SDK in your own Go project, follow these steps:

1. **Install the SDK**

   Add the Klever Go SDK to your project using `go get`:
   ```sh
   go get github.com/klever-io/klever-go-sdk
   ```

2. **Import the SDK Packages**

   Import the relevant packages in your Go files. For example:
   ```go
   import (
       "github.com/klever-io/klever-go-sdk/core"
       "github.com/klever-io/klever-go-sdk/provider"
   )
   ```

3. **Initialize and Use SDK Features**

   You can now use the SDK to interact with the Klever blockchain. For example, to create a wallet and send a transaction:
   ```go
   // Create a new wallet
   wallet := core.NewWallet()
   account, err := wallet.NewAccount()
   if err != nil {
       panic(err)
   }

   // Send a transaction (example)
   tx := provider.NewTransaction(account.Address, "recipient_address", 1000)
   err = provider.SendTransaction(account, tx)
   if err != nil {
       panic(err)
   }
   ```

## Examples

The `cmd/demo` directory provides a variety of runnable examples that demonstrate how to use the Klever Go SDK for different blockchain operations. Here are the types of examples included:

- **Asset Creation:**  
  Shows how to create a new asset (KDA) with minimal configuration, including setting supply and properties.

- **Deposit:**  
  Demonstrates how to deposit tokens into an account or contract.

- **Smart Contract Deployment:**  
  Provides an example of deploying a smart contract to the Klever blockchain.

- **Smart Contract Output Decoding:**  
  Illustrates how to decode smart contract outputs using ABI files, including parsing complex return types.

- **Multi-Contract Transactions:**  
  Shows how to build and send transactions that interact with multiple smart contracts in a single operation.

Each example is self-contained and can be run independently to learn about specific features of the SDK.