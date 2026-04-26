package v18

import "michelprogram/photon-parser/internal/reader"

func scanDictionary(r *reader.Reader, value *Value) error {
    keyType, err := r.ReadUInt8()
    if err != nil {
        return err
    }
    valueType, err := r.ReadUInt8()
    if err != nil {
        return err
    }
    count, err := r.ReadVarintUInt32()
    if err != nil {
        return err
    }
    start := r.Cursor
    for i := uint32(0); i < count; i++ {
        if _, err := scanPayload(r, ParameterType(keyType)); err != nil {
            return err
        }
        if _, err := scanPayload(r, ParameterType(valueType)); err != nil {
            return err
        }
    }
    value.Blob = r.Buffer[start:r.Cursor]
    value.Num = uint64(count)
    value.KeyType = ParameterType(keyType)
    value.ValType = ParameterType(valueType)
    return nil
}