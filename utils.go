package main

import (
	"strings"
	"math/big"
	"fmt"
)

func AddressToHex(address string) string {
	return strings.Replace(address, "0x", "", -1)
}

// HexToBigInt converts a hex string to a big.Int
func HexToBigInt(hexStr string) (*big.Int, error) {
	str := strings.Replace(hexStr, "0x", "", -1)
	value, ok := new(big.Int).SetString(str, 16)
	if !ok {
		return nil, fmt.Errorf("invalid hex string: %s", hexStr)
	}
	return value, nil
}
