// Copyright 2026 PointerByte Contributors
// SPDX-License-Identifier: Apache-2.0

package local

import "errors"

// pkcs7Pad applies PKCS#7 padding to b for the given block size. PKCS#7 encodes
// the pad length in a single byte, so blockSize must be in the range 1..255.
func pkcs7Pad(b []byte, blockSize int) []byte {
	if blockSize <= 0 || blockSize > 255 {
		panic("blockSize must be in the range 1..255")
	}
	padLen := blockSize - (len(b) % blockSize)
	if padLen == 0 {
		padLen = blockSize
	}
	// padLen is bounded by blockSize (<=255), so it always fits in a byte.
	pad := bytesRepeat(byte(padLen), padLen) // #nosec G115 -- padLen <= blockSize <= 255
	return append(b, pad...)
}

// pkcs7Unpad removes PKCS#7 padding from b and returns an error when the
// padding is invalid.
func pkcs7Unpad(b []byte, blockSize int) ([]byte, error) {
	if len(b) == 0 || len(b)%blockSize != 0 {
		return nil, errors.New("invalid padding: size")
	}
	padByte := b[len(b)-1]
	padLen := int(padByte)
	if padLen == 0 || padLen > blockSize || padLen > len(b) {
		return nil, errors.New("invalid padding: length")
	}

	for i := range padLen {
		if b[len(b)-1-i] != padByte {
			return nil, errors.New("invalid padding: content")
		}
	}
	return b[:len(b)-padLen], nil
}

// bytesRepeat returns a slice made of count copies of v.
func bytesRepeat(v byte, count int) []byte {
	out := make([]byte, count)
	for i := range out {
		out[i] = v
	}
	return out
}
