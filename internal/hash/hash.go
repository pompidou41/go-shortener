package hash

import (
	"crypto/sha256"
	"math/big"
)

const base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func saltToInt(salt string) *big.Int {
	hash := sha256.Sum256([]byte(salt))
	return new(big.Int).SetBytes(hash[26:])
}

func obfuscateId(id, saltInt *big.Int) *big.Int {
	return new(big.Int).Xor(id, saltInt)
}

func encodeBase62(num *big.Int) string {
	base := big.NewInt(62)
	zero := big.NewInt(0)

	if num.Cmp(zero) == 0 {
		return string(base62Alphabet[0])
	}

	result := ""

	n := new(big.Int).Set(num)
	mod := new(big.Int)

	for n.Cmp(zero) > 0 {
		n.DivMod(n, base, mod)
		result = string(base62Alphabet[mod.Int64()]) + result
	}

	return result
}

func decodeBase62(str string) *big.Int {
	result := big.NewInt(0)
	base := big.NewInt(62)

	for _, c := range str {
		// TODO: optimize to map O(n^2) to O(n)
		index := int64(indexOf(string(c)))
		result.Mul(result, base)
		result.Add(result, big.NewInt(index))
	}

	return result
}

func indexOf(char string) int {
	for i, c := range base62Alphabet {
		if string(c) == char {
			return i
		}
	}
	return -1
}

func EncodeId(id int64, salt string) string {
	idBigInt := big.NewInt(id)
	saltInt := saltToInt(salt)

	obfuscated := obfuscateId(idBigInt, saltInt)
	return encodeBase62(obfuscated)
}

func DecodeCode(code, salt string) int64 {
	saltInt := saltToInt(salt)

	obfuscated := decodeBase62(code)
	originalId := obfuscateId(obfuscated, saltInt)

	return originalId.Int64()
}
