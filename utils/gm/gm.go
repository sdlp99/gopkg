package gm

import (
	"bytes"
	rd "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"github.com/tjfoc/gmsm/sm4"
	"github.com/tjfoc/gmsm/x509"
	"math/rand"
	"strings"
	"time"
)

func RandStringBytes(l int) []byte {
	b := make([]byte, l)
	rand.Seed(time.Now().Unix())
	for i := 0; i < l; i++ {
		b[i] = uint8(rand.Intn(255))
	}
	return b
}

func BytesJoin(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}
func IntToBytes(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}
func EnSM4(str string) string {
	if strings.HasPrefix(str, "SM4(") &&
		strings.HasSuffix(str, ")") {
		return str
	}
	key := RandStringBytes(16)
	data := []byte(str)
	//iv := []byte("0000000000000000")
	//err = SetIV(iv)//设置SM4算法实现的IV值,不设置则使用默认值
	ecbMsg, err := sm4.Sm4Ecb(key, data, true) //sm4Ecb模式pksc7填充加密
	if err != nil {
		return ""
	}
	return "SM4(" + base64.StdEncoding.EncodeToString(BytesJoin(BytesJoin(key[8:16], key[0:8]), ecbMsg)) + ")"
}

func DeSM4(str string) string {
	if strings.HasPrefix(str, "SM4(") &&
		strings.HasSuffix(str, ")") {
		encStr := str[4 : len(str)-1]
		ecbMsg, err := base64.StdEncoding.DecodeString(encStr)
		if err != nil {
			return str
		}
		ecbDec, err := sm4.Sm4Ecb(BytesJoin(ecbMsg[8:16], ecbMsg[0:8]), ecbMsg[16:], false) //sm4Ecb模式pksc7填充解密
		if err != nil {
			return str
		}
		return string(ecbDec)
	} else {
		return str
	}

}

func CheckSm2(pubkey string, data string, sign string) bool {
	msg := []byte(data)
	sign2, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false
	}

	pub, err := x509.ReadPublicKeyFromHex(pubkey)
	if err != nil {
		return false
	}
	isok := pub.Verify(msg, sign2) //sm2验签
	return isok
}

func Sm2Sign(appid string, data string, key string) string {
	priv, err := x509.ReadPrivateKeyFromHex(key)
	if err != nil {
		return ""
	}
	msg := []byte(appid + data)
	sign, err := priv.Sign(rd.Reader, msg, nil) //sm2签名
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(sign)
}
