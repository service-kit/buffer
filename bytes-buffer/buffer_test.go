package bytes_buffer

import (
	"fmt"
	"testing"
)

func Test_BytesBuffer(t *testing.T) {
	buffer := NewBytesBuffer(8)
	fmt.Println(buffer.Cap(), buffer.Len())
	testBytes := make([]byte, 10000)
	for i := 0; i < 10000; i++ {
		testBytes[i] = byte(i)
	}
	for i := 0; i < 32; i++ {
		buffer.Write(testBytes)
	}
	fmt.Println("write pos id ", buffer.writePos.id, " read pos id ", buffer.readPos.id)
	for i := 0; i < 16; i++ {
		buffer.Read(testBytes)
	}
	fmt.Println(buffer.Cap(), buffer.Len())
	fmt.Println("write pos id ", buffer.writePos.id, " read pos id ", buffer.readPos.id)
	for i := 0; i < 48; i++ {
		buffer.Write(testBytes)
	}
	fmt.Println(buffer.Cap(), buffer.Len())
	fmt.Println("write pos id ", buffer.writePos.id, " read pos id ", buffer.readPos.id)
	for i := 0; i < 64; i++ {
		buffer.Read(testBytes)
	}
	fmt.Println(buffer.Cap(), buffer.Len())
	fmt.Println("write pos id ", buffer.writePos.id, " read pos id ", buffer.readPos.id)
}
