package errors

import "errors"

var HeaderSize = errors.New("header size to low")
var EncryptedPacket = errors.New("encrypted or unknown packet, signature")
