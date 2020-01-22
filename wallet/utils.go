package wallet

import "github.com/btcsuite/btcutil/base58"

// Base58Encode encode and returns the base58 encoding
// of the supplied byte
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

// Base58Encode decode and returns the raw bytes
// of the supplied byte
func Base58Decode(input []byte) []byte {
	decode := base58.Decode(string(input))

	return decode
}
