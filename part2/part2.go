package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
)

const (
	privateKey      = ""
	contractAddress = ""
	toAddress       = ""
)

func transfer(client *ethclient.Client, privateKey, toAddress, contract string) (string, error) {

	//从私钥推导出 公钥
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		fmt.Println("crypto.HexToECDSA error ,", err)
		return "", err
	}
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("publicKeyECDSA error ,", err)
		return "", err
	}
	//从公钥推导出钱包地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("钱包地址：", fromAddress.Hex())
	//构造请求参数
	//var data []byte
	//methodName := crypto.Keccak256([]byte("transfer(address,uint256)"))[:4]
	//paddedToAddress := common.LeftPadBytes(common.HexToAddress(toAddress).Bytes(), 32)
	//amount, _ := new(big.Int).SetString("100000000000000000000", 10)
	//paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	//data = append(data, methodName...)
	//data = append(data, paddedToAddress...)
	//data = append(data, paddedAmount...)

	//读取abi文件
	abiData, err := os.ReadFile("./part2/abi.json")
	if err != nil {
		fmt.Println("os.ReadFile error ,", err)
		return "", err
	}
	//将abi数据转成合约abi对象
	contractAbi, err := abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		fmt.Println("abi.JSON error ,", err)
		return "", err
	}
	amount, _ := new(big.Int).SetString("100000000000000000000", 10)
	data, err := contractAbi.Pack("transfer", common.HexToAddress(toAddress), amount)
	if err != nil {
		fmt.Println("contractAbi.Pack error ,", err)
		return "", err
	}

	//获取nonce
	nonce, err := client.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return "", err
	}
	fmt.Println("当前nonce:", nonce)
	//获取小费
	gasTipCap, _ := client.SuggestGasTipCap(context.Background())
	gas := uint64(100000)
	//最大gas fee
	gasFeeCap := big.NewInt(108694000460)

	contractAddress := common.HexToAddress(contract)
	//创建交易
	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gas,
		To:        &contractAddress,
		Value:     big.NewInt(0),
		Data:      data,
	})
	// 获取当前区块链的ChainID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		fmt.Println("获取ChainID失败:", err)
		return "", err
	}

	fmt.Println("当前区块链的ChainID:", chainID)
	//创建签名者
	signer := types.NewLondonSigner(chainID)
	//对交易进行签名
	signTx, err := types.SignTx(tx, signer, privateKeyECDSA)
	if err != nil {
		return "", err
	}
	//发送交易
	err = client.SendTransaction(context.Background(), signTx)
	if err != nil {
		return "", err
	}
	//返回交易哈希
	return signTx.Hash().Hex(), err

}

func main() {
	client, err := ethclient.Dial("https://goerli.infura.io/v3/3214cac49d354e48ad196cdfcefae1f8")
	if err != nil {
		fmt.Println("ethclient.Dial error : ", err)
		os.Exit(0)
	}
	tx, err := transfer(client, privateKey, toAddress, contractAddress)
	if err != nil {
		fmt.Println("transfer error : ", err)
		os.Exit(0)
	}

	fmt.Println("使用go调用智能合约第二讲：transfer tx : ", tx)

}
