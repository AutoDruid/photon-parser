package wireshark

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
)

type WiresharkParser struct{
	Requests []WireSharkRequest
}

func LoadFromWiresharkExport(path string) (*WiresharkParser, error){
	var res []WireSharkRequest

	data, err := os.ReadFile(path)
	if err != nil{
		return nil, err
	}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return &WiresharkParser{
		Requests: res,
	}, err
}

func (wp *WiresharkParser) RequestsDataToBytes() ([][]byte, error){
	res := make([][]byte, len(wp.Requests))
	for i := range wp.Requests{
		data := wp.Requests[i].Source.Layers.Data.DataData
		tmp, err := hex.DecodeString(strings.ReplaceAll(data, ":",""))
		if err != nil{
			return nil, err
		}
		res[i] =tmp
	}

	return res, nil
}
