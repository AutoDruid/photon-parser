package errors

import "errors"

var HeaderSize = errors.New("header size too low")
var EncryptedPacket = errors.New("packet is encrypted or unknown: unexpected signature byte")  
