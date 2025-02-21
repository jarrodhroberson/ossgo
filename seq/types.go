package seq

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type Decimal interface {
	float32 | float64
}

type Number interface {
	Integer | Decimal
}
