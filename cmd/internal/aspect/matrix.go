package aspect

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/xorvercom/camo/pkg/camo"
	"github.com/xorvercom/util/pkg/json"
)

var (
	InvalidMatrixJson = errors.New("invalid Matrix json")
	InvalidParam      = errors.New("invalid param")
)

// 乱数列マトリックス
type Mat struct {
	// 配列長
	length int
	// 乱数列
	matrix []*Lane
}

// 新規の空の乱数列マトリックスを生成
func MakeMatrix() *Mat {
	return &Mat{length: 0, matrix: []*Lane{}}
}

// レーン追加
func (m *Mat) AppendLane(bytes []byte) {
	m.matrix = append(m.matrix, MakeLane(bytes))
	m.length++
}

// 乱数列マトリックスのレーン数を返す
func (m *Mat) Length() int {
	return m.length
}

// 乱数列マトリックスのレーンを返す
func (m *Mat) Lane(no int) *Lane {
	return m.matrix[no%m.length]
}

const (
	// type property tag
	TYPE_TAG = "type"
	// type array value
	TYPE_VAL_ARRAY = "array"
	// lanes property tag
	LANES_TAG = "lanes"
)

// 乱数列マトリックスから JSON 要素を生成
func (m *Mat) ToJSON() json.Element {
	elem := json.NewElemObject()
	elem.Put(TYPE_TAG, json.NewElemString(TYPE_VAL_ARRAY))
	lanes := json.NewElemArray()
	elem.Put(LANES_TAG, lanes)
	for laneNo := 0; laneNo < m.Length(); laneNo++ {
		le := json.NewElemArray()
		lanes.Append(le)
		lane := m.Lane(laneNo)
		for pos := 0; pos < int(lane.Length()); pos++ {
			b := lane.GetByte(uint(pos))
			hex := hex.EncodeToString([]byte{b}) //fmt.Sprintf("%02x", b)
			le.Append(json.NewElemString(hex))
		}
	}
	return elem
}

// JSON 要素から乱数列マトリックスを生成
func MatrixFromJSON(json json.Element, params ...interface{}) (*Mat, error) {
	if root, ok := json.AsObject(); false == ok {
		return nil, errors.Wrap(InvalidMatrixJson, "json no object")
	} else {
		mat := MakeMatrix()
		// root = {"type":"array", "lane":[["FF", ...], ...]}
		ch := root.Child(TYPE_TAG)
		if false == ch.IsString() {
			return nil, errors.Wrap(InvalidMatrixJson, "property type must be string")
		}
		str, _ := ch.AsString()
		if false == strings.EqualFold(str.Text(), TYPE_VAL_ARRAY) {
			return nil, errors.Wrap(InvalidMatrixJson, "unknown property type: "+str.String())
		}
		ch = root.Child(LANES_TAG)
		if matrixarr, ok := ch.AsArray(); false == ok {
			return nil, errors.Wrap(InvalidMatrixJson, "lane not array")
		} else {
			// matrixarr = [["FF", ...], ...]
			sz := matrixarr.Size()
			for i := 0; i < sz; i++ {
				if lanearr, ok := matrixarr.Child(i).AsArray(); ok {
					// lanearr = ["FF", ...]
					lanesize := lanearr.Size()
					lane := make([]byte, lanesize)
					mat.AppendLane(lane)
					for j := 0; j < lanesize; j++ {
						if hex, ok := lanearr.Child(j).AsString(); ok {
							// hex = "FF"
							hexstr := hex.Value().(string)
							if hex, err := strconv.ParseUint(hexstr, 16, 8); err == nil {
								lane[j] = byte(hex % 0xff)
							} else {
								app := fmt.Sprintf("invalid hex: [%v, %v] %v", i, j, hexstr)
								return nil, errors.Wrap(err, app)
							}
						} else {
							app := fmt.Sprintf("invalid hex: [%v, %v]", i, j)
							return nil, errors.Wrap(InvalidMatrixJson, app)
						}
					}
				}
			}
			return mat, nil
		}
	}
}

// 乱数列マトリックスを生成するファクトリ
func MatrixFactory(params ...interface{}) (*Mat, error) {
	if len(params) == 0 {
		mat := &Mat{length: 0, matrix: []*Lane{}}
		p := camo.MakeRandomPattern(16/2, []int{53, 59, 61, 67, 71, 73, 79, 83})
		for i := 0; i < len(p); i++ {
			mat.AppendLane(p[i])
		}
		return mat, nil
	}
	if jsonname, ok := params[0].(string); ok {
		if j, err := json.LoadFromJSONFile(jsonname); err == nil {
			return MatrixFromJSON(j, params)
		} else {
			return nil, errors.Wrap(err, "invalid jsonfile: "+jsonname)
		}
	} else {
		return nil, InvalidParam
	}
}
