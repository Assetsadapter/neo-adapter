/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package neocoin

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/blocktree/openwallet/crypto"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/gjson"
	"strings"
)

//BlockchainInfo 本地节点区块链信息
type BlockchainInfo struct {
	Chain                string `json:"chain"`
	Blocks               uint64 `json:"blocks"`
	Headers              uint64 `json:"headers"`
	Bestblockhash        string `json:"bestblockhash"`
	Difficulty           string `json:"difficulty"`
	Mediantime           uint64 `json:"mediantime"`
	Verificationprogress string `json:"verificationprogress"`
	Chainwork            string `json:"chainwork"`
	Pruned               bool   `json:"pruned"`
}

func NewBlockchainInfo(json *gjson.Result) *BlockchainInfo {
	b := &BlockchainInfo{}
	//解析json
	b.Chain = gjson.Get(json.Raw, "chain").String()
	b.Blocks = gjson.Get(json.Raw, "blocks").Uint()
	b.Headers = gjson.Get(json.Raw, "headers").Uint()
	b.Bestblockhash = gjson.Get(json.Raw, "bestblockhash").String()
	b.Difficulty = gjson.Get(json.Raw, "difficulty").String()
	b.Mediantime = gjson.Get(json.Raw, "mediantime").Uint()
	b.Verificationprogress = gjson.Get(json.Raw, "verificationprogress").String()
	b.Chainwork = gjson.Get(json.Raw, "chainwork").String()
	b.Pruned = gjson.Get(json.Raw, "pruned").Bool()
	return b
}

// 账户余额 包含 NEO 主币 与 交易费用 GAS
type UnspentBalance struct {
	/*
		"balance": [
			{
				"unspent": [
					{
						"txid": "bd454059e58da4221aaf4effa3278660b231e9af7cea97912f4ac5c4995bb7e4",
						"n": 0,
						"value": 600.41014479
					}
				],
				"asset_hash": "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7",
				"asset": "GAS",
				"asset_symbol": "GAS",
				"amount": 29060.02316479
			},
			{
				"unspent": [
					{
						"txid": "c3182952855314b3f4b1ecf01a03b891d4627d19426ce841275f6d4c186e729a",
						"n": 0,
						"value": 800000
					}
				],
				"asset_hash": "c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b",
				"asset": "NEO",
				"asset_symbol": "NEO",
				"amount": 800000
			}
		],
		"address": "AGofsxAUDwt52KjaB664GYsqVAkULYvKNt"
	*/
	Key        string   `storm:"id"`
	NEOUnspent *Unspent `json:"neo_unspent"` // 未花费的 NEO
	GASUnspent *Unspent `json:"gas_unspent"` // 未花费的 GAS
	AccountID  string   `json:"account_id" storm:"index"`
	Address    string   `json:"address"` // 地址
	HDAddress  openwallet.Address
}

//Unspent 未花记录
type Unspent struct {
	/*
		{
			"unspent": [
				{
					"txid": "bd454059e58da4221aaf4effa3278660b231e9af7cea97912f4ac5c4995bb7e4",
					"n": 0,
					"value": 600.41014479
				}
			],
			"asset_hash": "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7",
			"asset": "GAS",
			"asset_symbol": "GAS",
			"amount": 29060.02316479
		}
	*/

	UnspentTxs  *[]UnspentTx `json:"unspent_txs"`
	AssetHash   string       `json:"asset_hash"`
	Asset       string       `json:"asset"`
	AssetSymbol string       `json:"asset_symbol"`
	Amount      string       `json:"amount"`
}

// 未花费交易信息
type UnspentTx struct {
	/*
		{
			"txid": "c3182952855314b3f4b1ecf01a03b891d4627d19426ce841275f6d4c186e729a",
			"n": 0,
			"value": 800000
		}
	*/
	TxID  string `json:"tx_id"`
	N     uint64 `json:"n"`
	Value string `json:"value"`
}

func NewUnspentBalance(json *gjson.Result) (*UnspentBalance, error) {
	obj := &UnspentBalance{}
	//解析json
	arr := json.Get("balance").Array()
	if len(arr) == 0 {
		return nil, errors.New("Balance is nil!")
	}
	for _, a := range arr {
		unspent := NewUnspent(&a)
		if unspent.AssetSymbol == AssetSymbolGAS {
			obj.GASUnspent = unspent
		} else {
			obj.NEOUnspent = unspent
		}
	}
	obj.Address = gjson.Get(json.Raw, "address").String()

	return obj, nil
}

func NewUnspent(json *gjson.Result) *Unspent {
	return &Unspent{
		UnspentTxs:  NewUnspentTxs(json.Get("unspent").Array()),
		AssetHash:   gjson.Get(json.Raw, "asset_hash").String(),
		Asset:       gjson.Get(json.Raw, "asset").String(),
		AssetSymbol: gjson.Get(json.Raw, "asset_symbol").String(),
		Amount:      gjson.Get(json.Raw, "amount").String(),
	}
}

func NewUnspentTxs(json []gjson.Result) *[]UnspentTx {
	unspentTxs := new([]UnspentTx)
	for _, j := range json {
		unspentTx := UnspentTx{
			TxID:  gjson.Get(j.Raw, "txid").String(),
			N:     gjson.Get(j.Raw, "n").Uint(),
			Value: gjson.Get(j.Raw, "value").String(),
		}
		*unspentTxs = append(*unspentTxs, unspentTx)
	}
	return unspentTxs
}

type UnspentSort struct {
	Values     []*UnspentBalance
	Comparator func(a, b *UnspentBalance) int
}

func (s UnspentSort) Len() int {
	return len(s.Values)
}
func (s UnspentSort) Swap(i, j int) {
	s.Values[i], s.Values[j] = s.Values[j], s.Values[i]
}
func (s UnspentSort) Less(i, j int) bool {
	return s.Comparator(s.Values[i], s.Values[j]) < 0
}

//type Address struct {
//	Address   string `json:"address" storm:"id"`
//	Account   string `json:"account" storm:"index"`
//	HDPath    string `json:"hdpath"`
//	CreatedAt time.Time
//}

type User struct {
	UserKey string `storm:"id"`     // primary key
	Group   string `storm:"index"`  // this field will be indexed
	Email   string `storm:"unique"` // this field will be indexed with a unique constraint
	Name    string // this field will not be indexed
	Age     int    `storm:"index"`
}

type Block struct {
	/*

		"hash": "000000000000000127454a8c91e74cf93ad76752cceb7eb3bcff0c398ba84b1f",
		"confirmations": 2,
		"strippedsize": 191875,
		"size": 199561,
		"weight": 775186,
		"height": 1354760,
		"version": 536870912,
		"versionHex": "20000000",
		"merkleroot": "48239e76f8b37d9c8824fef93d42ac3d7c433029c1e9fa23b6416dd0356f3e57",
		"tx": ["c1e12febeb58aefb0b01c04360262138f4ee0faeb207276e79ea3866608ed84f"]
		"time": 1532143012,
		"mediantime": 1532140298,
		"nonce": 3410287696,
		"bits": "19499855",
		"difficulty": 58358570.79038175,
		"chainwork": "00000000000000000000000000000000000000000000006f68c43926cd6c2d1f",
		"previousblockhash": "00000000000000292d142fcc1ddbd9dafd4518310009f152bdca2a66cc589f97",
		"nextblockhash": "0000000000004a50ef5733ab333f718e6ef5c1995e2cfd5a7caa0875f118cd30"

	*/

	Hash              string
	Confirmations     uint64
	Merkleroot        string
	tx                []string
	Previousblockhash string
	Height            uint64 `storm:"id"`
	Version           uint64
	Time              uint64
	Fork              bool
	txDetails         []*Transaction
	isVerbose         bool
}

func (wm *WalletManager) NewBlock(json *gjson.Result) *Block {
	obj := &Block{}
	//解析json
	obj.Height = gjson.Get(json.Raw, "index").Uint()
	obj.Hash = gjson.Get(json.Raw, "hash").String()
	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	obj.Merkleroot = gjson.Get(json.Raw, "merkleroot").String()
	obj.Previousblockhash = gjson.Get(json.Raw, "previousblockhash").String()
	obj.Version = gjson.Get(json.Raw, "version").Uint()
	obj.Time = gjson.Get(json.Raw, "time").Uint()

	txs := make([]string, 0)
	txDetails := make([]*Transaction, 0)
	for _, tx := range gjson.Get(json.Raw, "tx").Array() {
		if tx.IsObject() {
			obj.isVerbose = true
			txObj := wm.newTxByCore(&tx)
			txDetails = append(txDetails, txObj)
			txs = append(txs, txObj.TxID)
		} else {
			obj.isVerbose = false
			txs = append(txs, tx.String())
		}

	}

	obj.tx = txs
	obj.txDetails = txDetails

	return obj
}

//BlockHeader 区块链头
func (b *Block) BlockHeader(symbol string) *openwallet.BlockHeader {

	obj := openwallet.BlockHeader{}
	//解析json
	obj.Hash = b.Hash
	obj.Confirmations = b.Confirmations
	obj.Merkleroot = b.Merkleroot
	obj.Previousblockhash = b.Previousblockhash
	obj.Height = b.Height
	obj.Version = b.Version
	obj.Time = b.Time
	obj.Symbol = symbol

	return &obj
}

//UnscanRecords 扫描失败的区块及交易
type UnscanRecord struct {
	ID          string `storm:"id"` // primary key
	BlockHeight uint64
	TxID        string
	Reason      string
}

func NewUnscanRecord(height uint64, txID, reason string) *UnscanRecord {
	obj := UnscanRecord{}
	obj.BlockHeight = height
	obj.TxID = txID
	obj.Reason = reason
	obj.ID = common.Bytes2Hex(crypto.SHA256([]byte(fmt.Sprintf("%d_%s", height, txID))))
	return &obj
}

type Transaction struct {
	TxID          string
	Size          uint64
	Type          string
	Version       uint64
	Attributes    *[]Attribute
	Vins          []*Vin
	Vouts         []*Vout
	SysFee        string // 系统交易费 每笔交易都有10GAS的免费额度
	NetFee        string // 网络交易费 交易大小<1024 byte时网络费是可选的，最低为0.001GAS，>1024 byte时需要支付0.001GAS作为基础费用，且额外收取每字节 0.00001 GAS 的网络费
	BlockHash     string
	BlockHeight   uint64
	Confirmations uint64
	Blocktime     int64
}

type Attribute struct {
	/*
		usage	uint8	使用类型
		length	uint8	数据长度 (如有需要)
		data	uint8	使用类型相关的外部数据
	*/

	/*
		以下使用类型可以包括在交易的属性中
		0	ContractHash	合约脚本哈希	32
		2	ECDH02	用于ECDH密钥交换的公钥	32
		3	ECDH03	用于ECDH密钥交换的公钥	32
		32	Script	交易额外的验证	20
		48	Vote	投票payload	应指定(最多255个字节)
		129	DescriptionUrl	描述说明的URL	应指定 (最多255个字节)
		144	Description	说明	应指定 (最多255个字节
		161 - 175	Hash1-Hash15	自定义的存储哈希	32
		240 - 255	Remark-Remark15	自定义的一般备注	应指定 (最多65535个字节)
	*/

	/*
	   {
	      "usage":144,
	      "data":"5473685f323031372f322f382031363a31383a353931373033323132343035"
	   }
	*/

	Usage  uint64 // 使用类型
	Length uint8  // 数据长度 (如有需要, 可选)
	Data   string // 使用类型相关的外部数据
}

// 交易输入
type Vin struct {
	Coinbase string
	TxID     string
	Vout     uint64
	N        uint64
	Addr     string
	Value    string
}

// 交易输出
type Vout struct {
	N            uint64
	Addr         string
	Value        string
	Asset        string
	ScriptPubKey string
	Type         string
}

func (wm *WalletManager) newTxByCore(json *gjson.Result) *Transaction {

	/*
		{
			"txid": "0x28975702b73450d0f466e5b931eafbc04c0ea6a732162c548ff3d569fa627d9d",
			"size": 262,
			"type": "ContractTransaction",
			"version": 0,
			"attributes": [],
			"vin": [
				{
					"txid": "0x9e6b682209f778a1246202524be785633e03129b6877040ad05134cc96336fcb",
					"vout": 1
				}
			],
			"vout": [
				{
					"n": 0,
					"asset": "0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b",
					"value": "100",
					"address": "AGVziqTEhJJTQckrUuTQcyHNGV4ksKPPUT"
				},
				{
					"n": 1,
					"asset": "0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b",
					"value": "99999590",
					"address": "AXXYzk1kn9Bj8PHeqha921gqCpwJNRmuHC"
				}
			],
			"sys_fee": "0",
			"net_fee": "0",
			"scripts": [
				{
					"invocation": "40f96445be5bd95cdb0f45049d4f69792ca30bef05aa05ee1c6a82ae7884f96038509930a4b731242e1b6577e457e6f8330238b4c1da1460df3d1f88bb46f8052b",
					"verification": "21036943c02168ce22fb2e48a3f92dd72336d295e793a52633beba22ac46916dc201ac"
				}
			],
			"blockhash": "0xd87f1b76d89a158ed54a0cb88701e5d5ad86ce6f86399ecb50c589a65d709881",
			"confirmations": 4897,
			"blocktime": 1573037731
		}
	*/

	obj := Transaction{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "txid").String()
	obj.Size = gjson.Get(json.Raw, "size").Uint()
	obj.Type = gjson.Get(json.Raw, "type").String()
	obj.Version = gjson.Get(json.Raw, "version").Uint()
	obj.SysFee = gjson.Get(json.Raw, "sys_fee").String()
	obj.NetFee = gjson.Get(json.Raw, "net_fee").String()
	obj.BlockHash = gjson.Get(json.Raw, "blockhash").String()
	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	obj.Blocktime = gjson.Get(json.Raw, "blocktime").Int()

	obj.Attributes = new([]Attribute)
	if attributes := gjson.Get(json.Raw, "attributes"); attributes.IsArray() {
		for _, attr := range attributes.Array() {
			*(obj.Attributes) = append(*(obj.Attributes), newAttributeByCore(&attr))
		}
	}

	obj.Vins = make([]*Vin, 0)
	if vins := gjson.Get(json.Raw, "vin"); vins.IsArray() {
		for _, vin := range vins.Array() {
			obj.Vins = append(obj.Vins, newTxVinByCore(&vin))
		}
	}

	obj.Vouts = make([]*Vout, 0)
	if vouts := gjson.Get(json.Raw, "vout"); vouts.IsArray() {
		for _, vout := range vouts.Array() {
			obj.Vouts = append(obj.Vouts, newTxVoutByCore(&vout))
		}
	}

	return &obj
}

func newAttributeByCore(json *gjson.Result) Attribute {
	/*
	   {
	      "usage":144,
	      "data":"5473685f323031372f322f382031363a31383a353931373033323132343035"
	   }
	*/

	return Attribute{
		Usage: gjson.Get(json.Raw, "usage").Uint(),
		Data:  gjson.Get(json.Raw, "data").String(),
	}
}

func newTxVinByCore(json *gjson.Result) *Vin {
	/*
	   {
	      "txid":"0x3631f66024ca6f5b033d7e0809eb993443374830025af904fb51b0334f127cda",
	      "vout":0
	   }
	*/

	return &Vin{
		TxID: gjson.Get(json.Raw, "txid").String(),
		Vout: gjson.Get(json.Raw, "vout").Uint(),
	}
}

func newTxVoutByCore(json *gjson.Result) *Vout {

	/*
	   {
	      "n":1,
	      "asset":"0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b",
	      "value":"99999000",
	      "address":"AWHX6wX5mEJ4Vwg7uBcqESeq3NggtNFhzD"
	   }
	*/
	return &Vout{
		N:     gjson.Get(json.Raw, "n").Uint(),
		Asset: gjson.Get(json.Raw, "asset").String(),
		Value: gjson.Get(json.Raw, "value").String(),
		Addr:  gjson.Get(json.Raw, "address").String(),
	}
}

func DecodeScript(script string) ([]byte, error) {
	opcodes := strings.Split(script, " ")
	scriptBuilder := txscript.NewScriptBuilder()
	for _, codeName := range opcodes {
		code, ok := txscript.OpcodeByName[codeName]
		if ok {
			scriptBuilder.AddOp(code)
		} else {
			if len(codeName)%2 != 0 {
				codeName = "0" + codeName
			}
			data, err := hex.DecodeString(codeName)
			if err != nil {
				return nil, err
			}
			scriptBuilder.AddData(data)
		}
	}
	return scriptBuilder.Script()
}
