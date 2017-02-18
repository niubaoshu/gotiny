package gotiny

// func encBoolSlice(e *Encoder, v []bool) {
// 	l := len(v)
// 	e.encUint(uint64(l))
// }

// func encBytes(e *Encoder, v []byte) {
// 	l := len(v)
// 	e.encUint(uint64(l))
// 	e.append(l)
// 	copy(e.buf[e.index:], v)
// 	e.index += l
// }
