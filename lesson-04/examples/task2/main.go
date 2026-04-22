package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// task2.go
// 使用通用 ABI 调用 Counter 合约的方法，包括：
// 1. number: 查询number值
// 2. setNumber: 初始化number值 （需要设置 SENDER_PRIVATE_KEY 环境变量）
// 3. increment: number 自增 （需要设置 SENDER_PRIVATE_KEY 环境变量）
//
// 执行示例：
//
// 1. 查询number值：
//    export ETH_RPC_URL="http://127.0.0.1:8545"
//    go run main.go --mode number --contract 0xa513E6E4b8f2a923D98304ec87F64353C4D5C853
//
// 2. 初始化number值：
//    export ETH_RPC_URL="http://127.0.0.1:8545"
//    export SENDER_PRIVATE_KEY="your_private_key_hex"
//    go run main.go --mode setNumber --contract 0xa513E6E4b8f2a923D98304ec87F64353C4D5C853 --number 10
//
// 3. number 自增：
//    export ETH_RPC_URL="http://127.0.0.1:8545"
//    export SENDER_PRIVATE_KEY="your_private_key_hex"
//    go run main.go --mode increment --contract 0xa513E6E4b8f2a923D98304ec87F64353C4D5C853

// 注意事项：
// - 所有示例中的地址和交易哈希都是示例，请替换为实际值
// - setNumber和 increment 模式需要设置 SENDER_PRIVATE_KEY 环境变量（私钥十六进制，可带或不带 0x 前缀）
// - 仅在测试网或本地开发链上使用，不要在主网使用包含真实资产的私钥

const counterABIJSON = `[
	{
		"type": "function",
		"name": "increment",
		"inputs": [],
		"outputs": [],
		"stateMutability": "nonpayable"
	},
	{
		"type": "function",
		"name": "number",
		"inputs": [],
		"outputs": [
			{
				"name": "",
				"type": "uint256",
				"internalType": "uint256"
			}
		],
		"stateMutability": "view"
	},
	{
		"type": "function",
		"name": "setNumber",
		"inputs": [
			{
				"name": "newNumber",
				"type": "uint256",
				"internalType": "uint256"
			}
		],
		"outputs": [],
		"stateMutability": "nonpayable"
	}
]`

func main() {
	mode := flag.String("mode", "number", "operation mode: number, setNumber, or increment")
	contractHex := flag.String("contract", "", "Counter contract address")
	number := flag.Uint64("number", 0, "number value (Reset the number value)")
	flag.Parse()

	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		log.Fatal("ETH_RPC_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	// 读取 ABI 文件
	abiBytes, err := os.ReadFile("Counter.abi.json")
	if err != nil {
		log.Fatal("Failed to read ABI file:", err)
	}

	// 解析 ABI
	parsedABI, err := abi.JSON(strings.NewReader(string(abiBytes)))

	// parsedABI, err := abi.JSON(strings.NewReader(counterABIJSON))
	if err != nil {
		log.Fatalf("failed to parse ABI: %v", err)
	}

	switch *mode {
	case "number":
		handleNumberValue(ctx, client, parsedABI, *contractHex)
	case "setNumber":
		handleSetNumber(ctx, client, parsedABI, *contractHex, *number)
	case "increment":
		handleIncrement(ctx, client, parsedABI, *contractHex)

	default:
		log.Fatalf("unknown mode: %s (use: value, inc, or incBy)", *mode)
	}
}

// handleNumberValue 查询 number 值
func handleNumberValue(ctx context.Context, client *ethclient.Client, parsedABI abi.ABI, contractHex string) {
	if contractHex == "" {
		log.Fatal("missing --contract ")
	}

	contractAddr := common.HexToAddress(contractHex)

	// 编码 number 调用数据
	data, err := parsedABI.Pack("number")
	if err != nil {
		log.Fatalf("failed to pack data: %v", err)
	}

	callMsg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	// 执行只读调用
	output, err := client.CallContract(ctx, callMsg, nil)
	if err != nil {
		log.Fatalf("CallContract error: %v", err)
	}

	// 解码返回值
	var number *big.Int
	err = parsedABI.UnpackIntoInterface(&number, "number", output)
	if err != nil {
		log.Fatalf("failed to unpack output: %v", err)
	}

	fmt.Printf("Contract : %s\n", contractAddr.Hex())
	fmt.Printf("number  : %s (raw uint256)\n", number.String())
}

// handleSetNumber 发送 设置 number 值
func handleSetNumber(ctx context.Context, client *ethclient.Client, parsedABI abi.ABI, contractHex string, number uint64) {
	if contractHex == "" {
		log.Fatal("missing --contract")
	}

	// 检查私钥环境变量
	privKeyHex := os.Getenv("SENDER_PRIVATE_KEY")
	if privKeyHex == "" {
		log.Fatal("SENDER_PRIVATE_KEY is not set (required for transfer mode)")
	}

	// 解析私钥
	privKey, err := crypto.HexToECDSA(trim0x(privKeyHex))
	if err != nil {
		log.Fatalf("invalid private key: %v", err)
	}

	// 获取发送方地址
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddr := crypto.PubkeyToAddress(*publicKeyECDSA)

	contractAddr := common.HexToAddress(contractHex)

	// 获取链 ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("failed to get chain id: %v", err)
	}

	// 获取 nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		log.Fatalf("failed to get nonce: %v", err)
	}

	// 编码 setNumber 调用数据
	// setNumber(uint256 newNumber)
	newNumber := new(big.Int).SetUint64(number)
	callData, err := parsedABI.Pack("setNumber", newNumber)
	if err != nil {
		log.Fatalf("failed to pack setNumber data: %v", err)
	}

	// 估算 Gas Limit（合约调用需要更多 Gas）
	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From: fromAddr,
		To:   &contractAddr,
		Data: callData,
	})
	if err != nil {
		log.Fatalf("failed to estimate gas: %v", err)
	}
	// 增加 20% 的缓冲，避免 Gas 不足
	gasLimit = gasLimit * 120 / 100

	// 获取建议的 Gas 价格（使用 EIP-1559 动态费用）
	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		log.Fatalf("failed to get gas tip cap: %v", err)
	}

	// 获取 base fee，计算 fee cap
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("failed to get header: %v", err)
	}

	baseFee := header.BaseFee
	if baseFee == nil {
		// 如果不支持 EIP-1559，使用传统 gas price
		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			log.Fatalf("failed to get gas price: %v", err)
		}
		baseFee = gasPrice
	}

	// fee cap = base fee * 2 + tip cap（简单策略）
	gasFeeCap := new(big.Int).Add(
		new(big.Int).Mul(baseFee, big.NewInt(2)),
		gasTipCap,
	)

	// 检查 ETH 余额是否足够支付 Gas 费用
	balance, err := client.BalanceAt(ctx, fromAddr, nil)
	if err != nil {
		log.Fatalf("failed to get balance: %v", err)
	}

	// 计算总费用：gasFeeCap * gasLimit（Counter 转账不需要发送 ETH，只需要支付 Gas）
	totalGasCost := new(big.Int).Mul(gasFeeCap, big.NewInt(int64(gasLimit)))

	if balance.Cmp(totalGasCost) < 0 {
		log.Fatalf("insufficient ETH balance for gas: have %s wei, need %s wei", balance.String(), totalGasCost.String())
	}

	// 构造交易（EIP-1559 动态费用交易）
	// 注意：Counter setNumber 的 value 为 0，调用数据在 Data 字段中
	txData := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &contractAddr, // 合约地址
		Value:     big.NewInt(0), // Counter 转账不需要发送 ETH
		Data:      callData,      // setNumber 调用数据
	}
	tx := types.NewTx(txData)

	// 签名交易
	signer := types.NewLondonSigner(chainID)
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		log.Fatalf("failed to sign transaction: %v", err)
	}

	// 发送交易
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		log.Fatalf("failed to send transaction: %v", err)
	}

	// 输出交易信息
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Counter setNumber Transaction Sent\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Contract      : %s\n", contractAddr.Hex())
	fmt.Printf("setNumber: %d\n", number)
	fmt.Printf("Gas Limit     : %d\n", gasLimit)
	fmt.Printf("Gas Tip Cap   : %s Wei\n", gasTipCap.String())
	fmt.Printf("Gas Fee Cap   : %s Wei\n", gasFeeCap.String())
	fmt.Printf("Estimated Cost: %s Wei\n", totalGasCost.String())
	fmt.Printf("Nonce         : %d\n", nonce)
	fmt.Printf("Tx Hash       : %s\n", signedTx.Hash().Hex())
	fmt.Printf("\n")
	fmt.Printf("Transaction is pending. Waiting for confirmation...\n")
	fmt.Printf("\n")

	// 等待交易确认
	waitForTransaction(ctx, client, signedTx.Hash())
}

// handleIncrement 发送 number 自增
func handleIncrement(ctx context.Context, client *ethclient.Client, parsedABI abi.ABI, contractHex string) {
	if contractHex == "" {
		log.Fatal("missing --contract")
	}

	// 检查私钥环境变量
	privKeyHex := os.Getenv("SENDER_PRIVATE_KEY")
	if privKeyHex == "" {
		log.Fatal("SENDER_PRIVATE_KEY is not set (required for transfer mode)")
	}

	// 解析私钥
	privKey, err := crypto.HexToECDSA(trim0x(privKeyHex))
	if err != nil {
		log.Fatalf("invalid private key: %v", err)
	}

	// 获取发送方地址
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddr := crypto.PubkeyToAddress(*publicKeyECDSA)

	contractAddr := common.HexToAddress(contractHex)

	// 获取链 ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("failed to get chain id: %v", err)
	}

	// 获取 nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		log.Fatalf("failed to get nonce: %v", err)
	}

	// 编码 increment 调用数据
	// increment()
	callData, err := parsedABI.Pack("increment")
	if err != nil {
		log.Fatalf("failed to pack increment data: %v", err)
	}

	// 估算 Gas Limit（合约调用需要更多 Gas）
	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From: fromAddr,
		To:   &contractAddr,
		Data: callData,
	})
	if err != nil {
		log.Fatalf("failed to estimate gas: %v", err)
	}
	// 增加 20% 的缓冲，避免 Gas 不足
	gasLimit = gasLimit * 120 / 100

	// 获取建议的 Gas 价格（使用 EIP-1559 动态费用）
	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		log.Fatalf("failed to get gas tip cap: %v", err)
	}

	// 获取 base fee，计算 fee cap
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("failed to get header: %v", err)
	}

	baseFee := header.BaseFee
	if baseFee == nil {
		// 如果不支持 EIP-1559，使用传统 gas price
		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			log.Fatalf("failed to get gas price: %v", err)
		}
		baseFee = gasPrice
	}

	// fee cap = base fee * 2 + tip cap（简单策略）
	gasFeeCap := new(big.Int).Add(
		new(big.Int).Mul(baseFee, big.NewInt(2)),
		gasTipCap,
	)

	// 检查 ETH 余额是否足够支付 Gas 费用
	balance, err := client.BalanceAt(ctx, fromAddr, nil)
	if err != nil {
		log.Fatalf("failed to get balance: %v", err)
	}

	// 计算总费用：gasFeeCap * gasLimit（Counter 转账不需要发送 ETH，只需要支付 Gas）
	totalGasCost := new(big.Int).Mul(gasFeeCap, big.NewInt(int64(gasLimit)))

	if balance.Cmp(totalGasCost) < 0 {
		log.Fatalf("insufficient ETH balance for gas: have %s wei, need %s wei", balance.String(), totalGasCost.String())
	}

	// 构造交易（EIP-1559 动态费用交易）
	// 注意：Counter increment 的 value 为 0，调用数据在 Data 字段中
	txData := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &contractAddr, // 合约地址
		Value:     big.NewInt(0), // Counter 转账不需要发送 ETH
		Data:      callData,      // increment 调用数据
	}
	tx := types.NewTx(txData)

	// 签名交易
	signer := types.NewLondonSigner(chainID)
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		log.Fatalf("failed to sign transaction: %v", err)
	}

	// 发送交易
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		log.Fatalf("failed to send transaction: %v", err)
	}

	// 输出交易信息
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Counter Transfer Transaction Sent\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Contract      : %s\n", contractAddr.Hex())
	fmt.Printf("Gas Limit     : %d\n", gasLimit)
	fmt.Printf("Gas Tip Cap   : %s Wei\n", gasTipCap.String())
	fmt.Printf("Gas Fee Cap   : %s Wei\n", gasFeeCap.String())
	fmt.Printf("Estimated Cost: %s Wei\n", totalGasCost.String())
	fmt.Printf("Nonce         : %d\n", nonce)
	fmt.Printf("Tx Hash       : %s\n", signedTx.Hash().Hex())
	fmt.Printf("\n")
	fmt.Printf("Transaction is pending. Waiting for confirmation...\n")
	fmt.Printf("\n")

	// 等待交易确认
	waitForTransaction(ctx, client, signedTx.Hash())

	handleNumberValue(ctx, client, parsedABI, contractHex)
}

// waitForTransaction 等待交易确认并显示回执信息
func waitForTransaction(ctx context.Context, client *ethclient.Client, txHash common.Hash) {
	// 设置超时上下文（最多等待 2 分钟）
	waitCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	fmt.Printf("Polling for transaction receipt...\n")
	for {
		select {
		case <-waitCtx.Done():
			fmt.Printf("\nTimeout waiting for transaction confirmation.\n")
			fmt.Printf("You can check the transaction status later:\n")
			fmt.Printf("  go run main.go --mode parse-event --tx %s\n", txHash.Hex())
			return

		case <-ticker.C:
			receipt, err := client.TransactionReceipt(waitCtx, txHash)
			if err != nil {
				// 交易可能还在 pending
				continue
			}

			// 交易已确认
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			fmt.Printf("Transaction Confirmed!\n")
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			fmt.Printf("Status       : %d (1=success, 0=failed)\n", receipt.Status)
			fmt.Printf("Block Number : %d\n", receipt.BlockNumber.Uint64())
			fmt.Printf("Block Hash   : %s\n", receipt.BlockHash.Hex())
			fmt.Printf("Gas Used     : %d / %d\n", receipt.GasUsed, receipt.GasUsed)
			fmt.Printf("Logs Count   : %d\n", len(receipt.Logs))

			if receipt.Status == 0 {
				fmt.Printf("\n⚠️  Transaction failed! Check the transaction on block explorer.\n")
			} else {
				fmt.Printf("\n✅ Transaction successful!\n")
				if len(receipt.Logs) > 0 {
					fmt.Printf("\nTo parse Transfer event from this transaction:\n")
					fmt.Printf("  go run main.go --mode parse-event --tx %s\n", txHash.Hex())
				}
			}
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			return
		}
	}
}

// trim0x 移除十六进制字符串前缀 "0x"
func trim0x(s string) string {
	if len(s) >= 2 && s[0:2] == "0x" {
		return s[2:]
	}
	return s
}
