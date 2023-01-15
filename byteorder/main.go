package main

import(
	. "fmt"
	"unsafe"
)


func main(){
	Println(NativeEndian)
	var u64 uint64 =0x1122334455667788

	Println("Your way:")
	Println("uint64 to []byte:")
	b:=make([]byte, 8)
	NativeEndian.PutUint64(b, u64)
	Printf("% X\n\n", b)
	Println("[]byte to uint64:")
	u64=NativeEndian.Uint64(b)
	Printf("%X\n\n\n", u64)

	Println("My way:")
	Println("uint64 to []byte:")
	pb:=(*byte)( unsafe.Pointer(&u64) )
	b1:=unsafe.Slice(pb, 8)
	Printf("% X\n\n", b1)
	Println("[]byte to uint64:")
	u64=*(*uint64)( unsafe.Pointer(&b1[0]) )
	Printf("%X\n",u64)


}//main