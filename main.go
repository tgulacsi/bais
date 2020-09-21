package main

import (
	"fmt"
	"github.com/jarrodhroberson/bais/encoding/bais"
)

func main() {
	arr := []byte{67, 97, 116, 128, 10, 69, 255, 65, 66, 67, 68}
	fmt.Println(string(arr[:]))
	//str := "Cat\b`@iE?tEB!CD"
	str2 := bais.Encode(&arr, false)
	fmt.Println(string(str2[:]))
}
