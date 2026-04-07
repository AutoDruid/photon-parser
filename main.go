package main

import (
	"encoding/binary"
	"fmt"
	"michelprogram/photon-parser/photon"
	"michelprogram/photon-parser/pkg/wireshark"
	"os"
)

func main(){
	data, err := os.ReadFile("./ressources/wireshark.json")
	if err != nil{
		fmt.Errorf("cant read wireshark file")
		return
	}

	requests, err := wireshark.InitWiresharkParser(data)
	if err != nil{
		fmt.Errorf("cant read wireshark file")
		return
	}

	res, err := requests.RequestsDataToBytes()
	if err != nil{
		fmt.Errorf("cant read wireshark file")
		return
	}

	first := photon.ParseFromByte(res[0])

	
	fmt.Println(first)

	fmt.Println(test, first[0:2])
	val16 := binary.BigEndian.Uint16(test[0:2])
    fmt.Printf("16-bit value: %d %x)\n", val16, val16)
}
