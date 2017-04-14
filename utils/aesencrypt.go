package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func zeroPad(buf []byte) ([]byte, error) {
	bufLen := len(buf)
	padLen := aes.BlockSize - (bufLen % aes.BlockSize)
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(buf, padText...), nil
}

//秘钥key长度需要时AES-128(16bytes)或者AES-256(32bytes)
//原文sourcedata必须填充至blocksize的整数倍
func AESCBCEncrypter(key, sourcedata []byte) string {
	if len(sourcedata)%aes.BlockSize != 0 { //块大小在aes.BlockSize中定义
		LogDebugf("%s\n", "plaintext is not a multiple of the block size")
		sourcedata, _ = zeroPad(sourcedata)
	}

	block, err := aes.NewCipher(key) //生成加密用的block
	if err != nil {
		panic(err)
		return ""
	}

	// 对IV有随机性要求，但没有保密性要求，所以常见的做法是将IV包含在加密文本当中
	ciphertext := make([]byte, aes.BlockSize+len(sourcedata))
	//随机一个block大小作为IV
	//采用不同的IV时相同的秘钥将会产生不同的密文，可以理解为一次加密的session
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], sourcedata)

	ExampleNewCFBDecrypter(key, base64.StdEncoding.EncodeToString(ciphertext))
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func ExampleNewCFBDecrypter(key []byte, encrypteddata string) {
	ciphertext, _ := base64.StdEncoding.DecodeString(encrypteddata)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	LogDebugf("%s\n", hex.EncodeToString(iv))
	LogDebugf("%s\n", hex.EncodeToString(ciphertext))

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)
}
