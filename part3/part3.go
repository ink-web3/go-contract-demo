package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

func main() {
	client, err := ethclient.Dial("")
	if err != nil {
		fmt.Println("ethclient.Dial error : ", err)
		os.Exit(0)
	}
	token, err := NewToken(common.HexToAddress(contractAddress), client)
	if err != nil {
		fmt.Println("NewToken error : ", err)
	}
	//totalSupply, err := token.TotalSupply(nil)
	//if err != nil {
	//	fmt.Println("token.TotalSupply error : ", err)
	//}
	//fmt.Println("totalSupply is : ", totalSupply)

	// 获取当前区块链的ChainID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		fmt.Println("获取ChainID失败:", err)
		return
	}
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		fmt.Println("crypto.HexToECDSA error ,", err)
		return
	}

	gasTipCap, _ := client.SuggestGasTipCap(context.Background())

	//构建参数对象
	opts, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, chainID)
	if err != nil {
		fmt.Println("bind.NewKeyedTransactorWithChainID error ,", err)
		return
	}
	//设置参数
	opts.GasFeeCap = big.NewInt(108694000460)
	opts.GasLimit = uint64(100000)
	opts.GasTipCap = gasTipCap

	amount, _ := new(big.Int).SetString("100000000000000000000", 10)
	//调用合约transfer方法
	tx, err := token.Transfer(opts, common.HexToAddress(toAddress), amount)
	if err != nil {
		fmt.Println("token.Transfer error ,", err)
		return
	}

	fmt.Println("使用go调用智能合约第三讲：transfer tx : ", tx.Hash().Hex())

}
