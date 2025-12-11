package comm

func MakeBytes(count int, defVal byte) []byte {
	ret := make([]byte, count)
	for i := 0; i < len(ret); i++ {
		ret[i] = defVal
	}

	return ret
}
