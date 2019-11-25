package neoTransaction

import (
	"errors"
	"fmt"
)

type Transaction struct {
	/*
		type	uint8	交易类型
		version	uint8	兼容版本
		attributes	array	交易其他属性
		outputs	array	资产的接收地址
		inputs	array	交易的资产输入
		scripts	array	用于验证交易的脚本
	*/

	Type       byte
	Version    byte
	Attributes []TxAttribute
	Vouts      []TxOut
	Vins       []TxIn
	Scripts    []TxScript
}

func newEmptyTransaction(txType TransactionType, vins []Vin, vouts []Vout, attributes []Attribute) (*Transaction, error) {
	txtype := txType.hexValue
	txIn, err := newTxInForEmptyTrans(vins)
	if err != nil {
		return nil, err
	}

	txOut, err := newTxOutForEmptyTrans(vouts)
	if err != nil {
		return nil, err
	}

	txAttributes, err := newTxAttributeForEmptyTrans(attributes)
	if err != nil {
		return nil, err
	}

	version := byte(DefaultTxVersion)

	return &Transaction{txtype, version, txAttributes, txOut, txIn, nil}, nil
}

func (t Transaction) encodeToBytes() (ret []byte, err error) {
	ret = append(ret, t.Type)
	ret = append(ret, t.Version)
	ret = append(ret, byte(len(t.Attributes)))
	for _, attr := range t.Attributes {
		attrBytes, err := attr.toBytes()
		if err != nil {
			return nil, err
		}
		ret = append(ret, attrBytes...)
	}

	ret = append(ret, byte(len(t.Vins)))
	for _, vin := range t.Vins {
		inBytes, err := vin.toBytes()
		if err != nil {
			return nil, err
		}
		ret = append(ret, inBytes...)
	}

	ret = append(ret, byte(len(t.Vouts)))
	for _, vout := range t.Vouts {
		outBytes, err := vout.toBytes()
		if err != nil {
			return nil, err
		}
		ret = append(ret, outBytes...)
	}

	if t.Scripts == nil {
		return ret, nil
	}

	ret = append(ret, byte(len(t.Scripts)))
	for _, script := range t.Scripts {
		scriptBytes, err := script.toBytes()
		if err != nil {
			return nil, err
		}
		ret = append(ret, scriptBytes...)
	}

	return ret, err
}

func DecodeRawTransaction(txBytes []byte) (*Transaction, error) {
	limit := len(txBytes)

	if limit == 0 {
		return nil, errors.New("Invalid transaction data!")
	}

	var rawTx Transaction

	index := 0

	if index+4 > limit {
		return nil, errors.New("Invalid transaction data length!")
	}

	if index > limit {
		return nil, errors.New("Invalid transaction data length!")
	}

	rawTx.Type = txBytes[index]
	index++

	if index > limit {
		return nil, errors.New("Invalid transaction data length!")
	}

	rawTx.Version = txBytes[index]
	index++

	attrs, newIndex, err := decodeTxAttributeFromRawTrans(txBytes, index)
	if err != nil {
		return nil, err
	}
	index = newIndex
	rawTx.Attributes = attrs

	vins, newIndex, err := decodeTxInFromRawTrans(txBytes, index)
	if err != nil {
		return nil, err
	}
	index = newIndex
	rawTx.Vins = vins

	vouts, newIndex, err := decodeTxOutFromRawTrans(txBytes, index)
	if err != nil {
		return nil, err
	}
	index = newIndex
	rawTx.Vouts = vouts

	if index == limit {
		fmt.Println(rawTx.String())
		return &rawTx, nil
	}
	scrips, newIndex, err := decodeTxScriptVerificationFromRawTrans(txBytes, index)
	if err != nil {
		return nil, err
	}
	rawTx.Scripts = scrips

	return &rawTx, nil
}

func (t Transaction) cloneEmpty() Transaction {
	var ret Transaction
	ret.Type = t.Type
	ret.Version = t.Version
	ret.Attributes = append(ret.Attributes, t.Attributes...)
	ret.Vouts = append(ret.Vouts, t.Vouts...)
	ret.Vins = append(ret.Vins, t.Vins...)
	return ret
}

func (t *Transaction) String() string {
	fmtStr := "{ Transaction : { Type : %x, version : %x, "
	fmtParams := []interface{}{t.Type, t.Version}

	fmtStr += "Attribute : ["
	for _, v := range t.Attributes {
		fmtStr += v.String()
	}
	fmtStr += "],"

	fmtStr += "Vout : ["
	for _, vout := range t.Vouts {
		fmtStr += vout.String()
	}
	fmtStr += "],"
	fmtStr += "Vins : ["
	for _, vin := range t.Vins {
		fmtStr += vin.String()
	}
	fmtStr += "],"
	fmtStr += "Scripts : ["
	for _, script := range t.Scripts {
		fmtStr += script.String()
	}
	fmtStr += "]}"

	return fmt.Sprintf(fmtStr, fmtParams...)
}
