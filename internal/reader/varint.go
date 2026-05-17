package reader

// ReadVarintUInt32 reads a 32-bit unsigned integer from the reader in varint format.
func (r *Reader) ReadVarintUInt32() (uint32, error) {

	var res uint32
	var shift uint8
	for {

		buff, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		res |= uint32(buff&VARINT_MASK) << shift

		if buff&VARINT_MSB_MASK == 0 {
			break
		}
		shift += VARINT_SHIFT
	}

	return res, nil
}

// ReadVarintInt32 reads a 32-bit signed integer from the reader in varint format.
func (r *Reader) ReadVarintInt32() (int32, error) {

	res, err := r.ReadVarintUInt32()
	if err != nil {
		return 0, err
	}

	//ZigZag decode
	return int32((res >> 1) ^ uint32(-(int32(res & 1)))), nil
}

// ReadVarintUInt64 reads a 64-bit unsigned integer from the reader in varint format.
func (r *Reader) ReadVarintUInt64() (uint64, error) {

	var res uint64
	var shift uint8
	for {

		buff, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		res |= uint64(buff&VARINT_MASK) << shift

		if buff&VARINT_MSB_MASK == 0 {
			break
		}
		shift += VARINT_SHIFT
	}

	return res, nil
}

// ReadVarintInt64 reads a 64-bit signed integer from the reader in varint format.
func (r *Reader) ReadVarintInt64() (int64, error) {

	res, err := r.ReadVarintUInt64()
	if err != nil {
		return 0, err
	}

	//ZigZag decode
	return int64((res >> 1) ^ uint64(-(int64(res & 1)))), nil
}
