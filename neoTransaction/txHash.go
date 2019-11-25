package neoTransaction

import (
	"encoding/hex"
	"errors"
	"github.com/blocktree/go-owcrypt"
)

type NormalTx struct {
	Address string
	SigType byte
	SigPub  SignaturePubkey
}

type MultiTx struct {
	Pubkey  string
	SigType byte
	SigPub  SignaturePubkey
}

type TxHash struct {
	Hash      string
	NRequired byte
	Normal    *NormalTx
	Multi     []MultiTx
}

func (tx TxHash) IsMultisig() bool {
	if tx.NRequired == 0 {
		return false
	}
	return true
}

func (tx TxHash) GetTxHashHex() string {
	return tx.Hash
}

func (tx TxHash) GetNormalTxAddress() string {
	return tx.Normal.Address
}

func (tx TxHash) GetMultiTxPubkeys() []string {
	var ret []string
	for _, p := range tx.Multi {
		ret = append(ret, p.Pubkey)
	}
	return ret
}

func newTxHash(hash, lockscript, redeem []byte, inType, sigType byte, addressPrefix AddressPrefix) (*TxHash, error) {
	//var prefixStr string
	//var p2pkhPrefixByte []byte
	//var p2wpkhPrefixByte []byte
	//
	//prefixStr = addressPrefix.Bech32Prefix
	//
	//p2pkhPrefixByte = addressPrefix.P2PKHPrefix
	//p2wpkhPrefixByte = addressPrefix.P2WPKHPrefix
	//
	//if inType == TypeP2PKH {
	//	return &TxHash{hex.EncodeToString(hash), 0, &NormalTx{nil, sigType, SignaturePubkey{nil, nil}}, nil}, nil
	//} else if inType == TypeP2WPKH {
	//	return &TxHash{hex.EncodeToString(hash), 0, &NormalTx{EncodeCheck(p2wpkhPrefixByte, lockscript[2:22]), sigType, SignaturePubkey{nil, nil}}, nil}, nil
	//} else if inType == TypeBech32 {
	//	return &TxHash{hex.EncodeToString(hash), 0, &NormalTx{Bech32Encode(prefixStr, BTCBech32Alphabet, lockscript[2:]), sigType, SignaturePubkey{nil, nil}}, nil}, nil
	//} else if inType == TypeMultiSig {
	//	nRequired, pubkeys, err := getMultiDetails(redeem)
	//	if err != nil {
	//		return nil, err
	//	}
	//	var multiTx []MultiTx
	//	for _, p := range pubkeys {
	//		multiTx = append(multiTx, MultiTx{p, sigType, SignaturePubkey{nil, nil}})
	//	}
	//	return &TxHash{hex.EncodeToString(hash), nRequired, nil, multiTx}, nil
	//}
	return nil, nil
}

func checkScriptType(scriptPubkey, redeemScript string) ([]byte, []byte, byte, error) {
	//script, err := hex.DecodeString(scriptPubkey)
	//if err != nil {
	//	return nil, nil, 0, errors.New("Invalid scriptPubkey data!")
	//}
	//if len(script) == 25 && script[0] == OpCodeDup && script[1] == OpCodeHash160 && script[2] == 0x14 && script[23] == OpCodeEqualVerify && script[24] == OpCodeCheckSig {
	//	if redeemScript != "" {
	//		return nil, nil, 0, errors.New("Found redeemScript when unlock a p2pkh input!")
	//	}
	//	return script, nil, TypeP2PKH, nil
	//} else if len(script) == 23 && script[0] == OpCodeHash160 && script[1] == 0x14 && script[22] == OpCodeEqual {
	//	redeem, err := hex.DecodeString(redeemScript)
	//	if err != nil {
	//		return nil, nil, 0, errors.New("Invalid redeemScript for a P2SH type input!")
	//	}
	//	if len(redeem) == 22 && redeem[0] == 0x00 && redeem[1] == 0x14 {
	//		return script, redeem, TypeP2WPKH, nil
	//	}
	//	if len(redeem) >= 37 && redeem[len(redeem)-1] == OpCheckMultiSig {
	//		return script, redeem, TypeMultiSig, nil
	//	}
	//} else if len(script) == 22 && script[0] == 0x00 && script[1] == 0x14 {
	//	if redeemScript != "" {
	//		return nil, nil, 0, errors.New("Found redeemScript when unlock a bech32 input!")
	//	}
	//	return script, nil, TypeBech32, nil
	//}
	return nil, nil, 0, errors.New("Unknown type of lockScript!")
}

func (t Transaction) calcSegwitSerializationHashes() ([]byte, []byte, []byte) {
	//hashPrevouts := []byte{}
	//hashSequence := []byte{}
	//hashOutputs := []byte{}
	//
	//for _, vin := range t.Vins {
	//	hashPrevouts = append(hashPrevouts, vin.txID...)
	//	hashPrevouts = append(hashPrevouts, vin.vout...)
	//
	//	hashSequence = append(hashSequence, vin.sequence...)
	//}
	//for _, vout := range t.Vouts {
	//hashOutputs = append(hashOutputs, vout.amount...)
	//hashOutputs = append(hashOutputs, byte(len(vout.lockScript)))
	//hashOutputs = append(hashOutputs, vout.lockScript...)
	//}
	//return owcrypt.Hash(hashPrevouts, 0, owcrypt.HASh_ALG_DOUBLE_SHA256),
	//	owcrypt.Hash(hashSequence, 0, owcrypt.HASh_ALG_DOUBLE_SHA256),
	//	owcrypt.Hash(hashOutputs, 0, owcrypt.HASh_ALG_DOUBLE_SHA256)
	return nil, nil, nil
}

func genScriptCodeFromRedeemScript(redeemBytes []byte) ([]byte, error) {

	ret := []byte{}
	if redeemBytes[0] == 0x00 && redeemBytes[1] == 0x14 {
		ret = redeemBytes[2:]

		if len(ret) != 0x14 {
			return nil, errors.New("Invalid redeem script!")
		}
		ret = append([]byte{byte(len(ret))}, ret...)
		//ret = append([]byte{OpCodeDup, OpCodeHash160}, ret...)
		//ret = append(ret, OpCodeEqualVerify, OpCodeCheckSig)
	} else {
		ret = redeemBytes
	}
	return ret, nil
}

func (t Transaction) getSegwitBytesForSig(reddemBytes, txid, vout, sequence []byte, sigType byte, amount uint64) ([]byte, error) {
	sigBytes := []byte{}

	sigBytes = append(sigBytes, t.Version)

	hashPrevouts, hashSequence, hashOutputs := t.calcSegwitSerializationHashes()

	sigBytes = append(sigBytes, hashPrevouts...)
	sigBytes = append(sigBytes, hashSequence...)

	sigBytes = append(sigBytes, txid...)
	sigBytes = append(sigBytes, vout...)

	scriptCode, err := genScriptCodeFromRedeemScript(reddemBytes)
	if err != nil {
		return nil, err
	}

	sigBytes = append(sigBytes, byte(len(scriptCode)))
	sigBytes = append(sigBytes, scriptCode...)

	if amount == 0 {
		return nil, errors.New("Invalid amount of input!")
	}

	sigBytes = append(sigBytes, uint64ToLittleEndianBytes(amount)...)
	sigBytes = append(sigBytes, sequence...)

	sigBytes = append(sigBytes, hashOutputs...)
	//sigBytes = append(sigBytes, t.LockTime...)

	return sigBytes, nil
}

func (t Transaction) getBytesForSig() ([]byte, error) {
	sigBytes := []byte{}
	var err error
	sigBytes, err = t.encodeToBytes()
	if err != nil {
		return nil, err
	}

	return sigBytes, nil
}

func (t Transaction) getHashesForSig() ([]TxHash, error) {
	hashes := []TxHash{}
	if t.Vouts == nil || len(t.Vouts) == 0 {
		return nil, errors.New("No output found!")
	}

	emptyTrans := t.cloneEmpty()

	emptyTransBytes, err := emptyTrans.encodeToBytes()
	if err != nil {
		return nil, err
	}

	hash := owcrypt.Hash(emptyTransBytes, 0, owcrypt.HASH_ALG_SHA256)

	for _, script := range t.Scripts {
		pubKey, err := script.GetPubKeyByVerificationScript()
		if err != nil {
			return nil, err
		}
		sign, err := script.GetSignatureByInvocationScript()
		hashes = append(hashes, TxHash{
			hex.EncodeToString(hash),
			0,
			&NormalTx{
				"",
				0,
				SignaturePubkey{
					sign,
					pubKey,
				},
			},
			nil,
		})
	}
	return hashes, nil
}

func (t TxHash) encodeToScript(redeem []byte, SegwitON bool) ([]byte, error) {
	//if t.NRequired == 0 {
	//	if t.Normal.SigPub.Signature == nil || len(t.Normal.SigPub.Signature) != 64 {
	//		return nil, errors.New("Invalid signature data!")
	//	}
	//	if t.Normal.SigPub.Pubkey == nil || len(t.Normal.SigPub.Pubkey) != 33 {
	//		return nil, errors.New("Invalid pubkey data!")
	//	}
	//	return t.Normal.SigPub.encodeToScript(t.Normal.SigType), nil
	//}
	//
	//count := byte(0)
	//for _, s := range t.Multi {
	//	if s.SigPub.Pubkey == nil && s.SigPub.Signature == nil {
	//		continue
	//	}
	//	if len(s.SigPub.Pubkey) != 33 || len(s.SigPub.Signature) != 64 {
	//		return nil, errors.New("Invalid signature or pubkey data for multisig!")
	//	}
	//	count++
	//}
	//
	//if count < t.NRequired {
	//	return nil, errors.New("The multisig transaction is not complete signed yet!")
	//}

	var ret []byte
	//var sigs []byte
	//
	//if redeem == nil {
	//	return nil, errors.New("Missing redeem for multisig!")
	//}
	//redeemLen := len(redeem)
	//ret = append(ret, redeem...)
	//
	//if !SegwitON {
	//	if redeemLen < 0x4C {
	//		ret = append([]byte{byte(redeemLen)}, ret...)
	//	} else if redeemLen <= 0xFF {
	//		ret = append([]byte{byte(redeemLen)}, ret...)
	//		ret = append([]byte{OpPushData1}, ret...)
	//	} else if redeemLen <= 0xFFFF {
	//		ret = append(uint16ToLittleEndianBytes(uint16(redeemLen)), ret...)
	//		ret = append([]byte{OpPushData2}, ret...)
	//	} else {
	//		return nil, errors.New("MultiSig redeem data is too long!")
	//	}
	//} else {
	//	ret = append([]byte{byte(redeemLen)}, ret...)
	//}
	//
	//count = 0
	//for _, s := range t.Multi {
	//	if s.SigPub.Signature == nil {
	//		continue
	//	}
	//	sigs = append(sigs, s.SigPub.encodeSignatureToScript(s.SigType)...)
	//	count++
	//
	//	if count == t.NRequired {
	//		break
	//	}
	//}
	//
	//ret = append(sigs, ret...)
	//ret = append([]byte{0x00}, ret...)
	//if !SegwitON {
	//	length := len(ret)
	//	if length < 0xFD {
	//		ret = append([]byte{byte(length)}, ret...)
	//	} else {
	//		ret = append(uint16ToLittleEndianBytes(uint16(length)), ret...)
	//		ret = append([]byte{0xFD}, ret...)
	//	}
	//}

	return ret, nil
}
