package main

import (
	. "fmt"
	"unsafe"
)

func main() {
	var u uint64

	for u = 0; u < 0x100000000; u++ {
		i:=int32(u)

		// first way:
		// Most languages (e.g. C/Go) encapsulates "value convertions" from int/uint to float
		// These convertion ether use cpu directives or (when cpu don't have such directives) use a well-knonwn way to do the same work.
		// The key part is `i` must be a signed integer type, `float64(u)/256` will return a wrong float64 value
		// Most developers should use this way unless the language you use don't encapsulates "value convertions" from int/uint to float
		// This way is much more readable and easy to understand
		fA:=float64(i)/256

		// second way:
		// This function use variant of a well-known way(mentioned in first way) to "value convert" int to float64
		// Like first way, the key part is arg `i` must be a signed integer type
		// I don't know how it works but I implemented a clumsy way to do the same work
		// Most developers should just use the first way, unless the language you use don't encapsulates "value convertions" from int/uint to float
		fB:=WlFixedToFloat64(i)
		if fA!=fB{
			Printf("i: %X\nfGoSimple: %X\nfGo: %X\n", i, fA, fB)
			break
		}
		
		// My way:
		// My clumsy way to "value convert" 24.8 fixed value to float64 value
		fC := MyFixedToFloat64(uint32(i))
		if fA!=fC{
			Printf("i: %X\nfGoSimple: %X\nfC: %X\n", i, fA, fC)
			break
		}

	}//for

} //main

// 
func WlFixedToFloat64(f int32) float64 {
	u_i := (1023+44)<<52 + (1 << 51) + int64(f)
	u_d := *(*float64)( unsafe.Pointer(&u_i) )
	return u_d - (3 << 43)
}


func MyFixedToFloat64(fix uint32) float64 {


	// fixed value have 1 sign bit and 31 value bits
	// decimal point between bit9 and bit8(right to left index order, start at 1)

	// get sign bit
	var sign uint64 = uint64( fix&0x80000000 )<<32

	// if fixed value is positive
	if sign==0{
		// if value bits is all `0`, return float64 value 0
		if fix==0{return 0}

	}else{
		// if fixed value is negtive
		// convert fix bits from complement form to true form
		// then put sign bit to zero
		fix=((^fix)+1)&0x7fffffff

		// there is a special case, when all fix value bits is `0`, (^fix)+1 will overflow
		// the really true form should be 0x80000000
		if fix==0{fix=0x80000000}
	}







	// shift out all successive `0`s on the left side and save number of zeros shifted out in `n0`
	// the most left `1` will be the hidden bit of Fraction, should be shifted out too
	// shift out n0 bits of `0` and 1 bit of `1` on left, pad in n0+1 bits of `0` on the right, so number of value bits is still 32
	// decimal point and value bits move together(n0+1bits left), don't change the value;
	// decimal point between bit(n0+1+9) and bit(n0+1+8) i.e bit(n0+10) and bit(n0+9)
	var n0 uint64
	for {
		if fix&0x80000000 != 0  {
			fix <<=1
			break
		}
		fix <<= 1
		n0++
	} //for

	// float64 have 52 Fraction bits, fix only have 32 bits
	// so fix should be convert to uint64 and shift left 20 bits(pad 20bits of `0` on the right side)
	// decimal point and value bits move together(20bits left), don't change the value;
	// decimal point between bit(20+n0+10) and bit(20+n0+9) i.e bit(30+n0) and bit(29+n0)
	var frac uint64 = uint64(fix) << 20

	// float64's decimal point between bit53 and bit 52
	// so float64's decimal point should move 52-(29+n0)bits right i.e (23-n0)bits right
	// move decimal point right equals to increase exponent, so exponent should increase by 22-n0
	// float64's exponent Biased up by 1023, so exponent should be 1023+23-n0 i.e 1046-n0
	// shift left 52bits put exponent at the correct place
	var exp uint64 = (1046 - n0) << 52

	// `or` sign, exponent and fraction together, get the float64bits
	fBits := sign | exp | frac

	// pointer convert float64bits to float64 value
	return *(*float64)(unsafe.Pointer(&fBits))
} // func