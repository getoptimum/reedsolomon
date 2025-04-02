package reedsolomon

func GalMulSlice(c byte, in, out []byte) {
	galMulSlice(c, in, out, &defaultOptions)
}

func GalMulSliceXor(c byte, in, out []byte) {
	galMulSliceXor(c, in, out, &defaultOptions)
}

func Inv(e byte) byte {
	return invTable[e]
}
