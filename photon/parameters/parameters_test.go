package parameters_test

import (
	"bytes"
	"fmt"
	"log"
	"michelprogram/photon-parser/photon/parameters"
	"testing"
)

func TestReliableInt32(t *testing.T) {

	payload := []byte{0x0, 0x69, 0x0, 0x4, 0x6c, 0xcb}
	reader := bytes.NewReader(payload)
	res, err := parameters.Parse(reader)

	if err != nil {
		t.Fatalf("LoadFromWiresharkExport() failed: %v", err)
	}

	if res.Value != int32(289995){
		t.Fatalf("should be int 289995: %v", err)
	}

}


func TestReliableString(t *testing.T) {

	payload := []byte{0x2, 0x73, 0x0, 0xd, 0x4c, 0x6f, 0x75, 0x74, 0x72, 0x65, 0x56, 0x65, 0x52, 0x74, 0x45, 0x61, 0x7a}
	reader := bytes.NewReader(payload)
	res, err := parameters.Parse(reader)

	if err != nil {
		t.Fatalf("LoadFromWiresharkExport() failed: %v", err)
	}

	if res.Value != "LoutreVeRtEaz"{
		t.Fatalf("should be int 289995: %v", err)
	}

}

func TestReliableFloat32(t *testing.T) {

	payload := []byte{0xa, 0x66, 0x43, 0x34, 0x3b, 0xdd}
	reader := bytes.NewReader(payload)
	res, err := parameters.Parse(reader)

	if err != nil {
		t.Fatalf("LoadFromWiresharkExport() failed: %v", err)
	}

	if res.Value != float32(180.23384){
		t.Fatalf("should be int 180.23384: %v", err)
	}

}

func TestReliableBoolean(t *testing.T) {

	payload := []byte{0x31, 0x6f, 0x01}
	reader := bytes.NewReader(payload)
	res, err := parameters.Parse(reader)

	if err != nil {
		t.Fatalf("LoadFromWiresharkExport() failed: %v", err)
	}

	if res.Value != float32(180.23384){
		t.Fatalf("should be int 180.23384: %v", err)
	}

}

func TestReliableArrayInt8(t *testing.T) {

	payload := []byte{0x5, 0x78, 0x0, 0x0, 0x0, 0x5, 0x1, 0x0, 0x3, 0x2, 0x5}
	reader := bytes.NewReader(payload)
	res, err := parameters.Parse(reader)

	if err != nil {
		t.Fatalf("LoadFromWiresharkExport() failed: %v", err)
	}

	log.Println(res)
	if res.Value != float32(180.23384){
		t.Fatalf("should be int 180.23384: %v", err)
	}

}

func TestReliableArrayInt32(t *testing.T) {

	payload := []byte{0x5, 0x78, 0x0, 0x0, 0x0, 0x5, 0x1, 0x0, 0x3, 0x2, 0x5}
	reader := bytes.NewReader(payload)
	res, err := parameters.Parse(reader)

	if err != nil {
		t.Fatalf("LoadFromWiresharkExport() failed: %v", err)
	}

	log.Println(res)
	if res.Value != float32(180.23384){
		t.Fatalf("should be int 180.23384: %v", err)
	}

}

func TestReliableArrayInt(t *testing.T) {

	payload := []byte{0x0, 0x79, 0x0, 0x2, 0x6b, 0x2, 0x84, 0x2, 0x85}
	reader := bytes.NewReader(payload)
	res, err := parameters.Parse(reader)

	if err != nil {
		t.Fatalf("LoadFromWiresharkExport() failed: %v", err)
	}

	fmt.Println(res.Value)
}


func TestReliableArrayFloat(t *testing.T) {

	payload := []byte{0x1, 0x79, 0x0, 0x2, 0x66, 0xc2, 0x97, 0xfe, 0x49, 0xc3, 0xba, 0xf3, 0xeb}
	reader := bytes.NewReader(payload)
	res, err := parameters.Parse(reader)

	if err != nil {
		t.Fatalf("LoadFromWiresharkExport() failed: %v", err)
	}

	fmt.Println(res.Value)
}
