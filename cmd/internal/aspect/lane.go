package aspect

// 乱数列
type Lane struct {
	length int
	stream []byte
}

// 新規の乱数列を生成
func MakeLane(stream []byte) *Lane {
	return &Lane{length: len(stream), stream: stream}
}

// 乱数列の長さを返す
func (l *Lane) Length() int {
	return l.length
}

// 乱数列のインデックス位置のバイトを返す
func (l *Lane) GetByte(index uint) byte {
	pos := index % uint(l.length)
	return l.stream[pos]
}
