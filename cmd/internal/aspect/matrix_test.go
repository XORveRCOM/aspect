package aspect

import (
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/xorvercom/util/pkg/fileutil"
	"github.com/xorvercom/util/pkg/json"
)

func TestMatrixFactory(t *testing.T) {
	// 存在しないJSONファイル
	if _, err := MatrixFactory("notexist.json"); err == nil {
		t.Error("MatrixFactory(\"notexist.json\") error = nil")
		return
	} else {
		t.Logf("MatrixFactory(\"notexist.json\") normal error = %v", err)
	}

	// 実装されていないパラメータ
	if _, err := MatrixFactory(-1); err == nil {
		t.Error("MatrixFactory(-1) error = nil")
		return
	} else {
		t.Logf("MatrixFactory(-1) normal error = %v", err)
	}

	tt := struct {
		name string
		file string
	}{name: "test 1", file: "out.json"}
	t.Run(tt.name, func(t *testing.T) {
		if src_mat, err := MatrixFactory(); err != nil {
			t.Errorf("MatrixFactory() error = %v", err)
			return
		} else {
			err := fileutil.TempSpace(func(tempdir string) error {
				// JSON ファイル化
				outJsonName := filepath.Join(tempdir, "out.json")
				_ = json.SaveToJSONFile(outJsonName, src_mat.ToJSON(), true)

				// JSON
				// import "github.com/xorvercom/util/pkg/text"
				// t.Log("[" + outJsonName + "]")
				// if textlines, err := text.LoadFrom(outJsonName); err != nil {
				// 	return err
				// } else {
				// 	for _, line := range textlines.Lines() {
				// 		t.Log(line)
				// 	}
				// }

				// JSON ファイルから読み込む
				if load_mat, err := MatrixFactory(outJsonName); err != nil {
					return err
				} else {
					if src_mat.Length() != load_mat.Length() {
						return errors.New("src_mat.Length()!=load_mat.Length()")
					}
				}
				return nil
			})
			if err != nil {
				t.Errorf("%v", err)
				return
			}
		}
	})
}

func TestMatrixFromJSON(t *testing.T) {
	var err error

	// json no object: invalid Matrix json
	_, err = MatrixFromJSON(json.NewElemNull())
	if err == nil {
		t.Error()
		return
	}
	t.Log(err)

	// property type must be string: invalid Matrix json
	if j, err := json.LoadFromJSONByte([]byte("{}")); err != nil {
		t.Log(err)
	} else {
		_, err = MatrixFromJSON(j)
		if err == nil {
			t.Error()
			return
		}
		t.Log(err)
	}

	// unknown property type: "": invalid Matrix json
	if j, err := json.LoadFromJSONByte([]byte("{\"type\":\"\"}")); err != nil {
		t.Log(err)
	} else {
		_, err = MatrixFromJSON(j)
		if err == nil {
			t.Error()
			return
		}
		t.Log(err)
	}

	// lane not array: invalid Matrix json
	if j, err := json.LoadFromJSONByte([]byte("{\"type\":\"array\", \"lanes\": 1}")); err != nil {
		t.Log(err)
	} else {
		_, err = MatrixFromJSON(j)
		if err == nil {
			t.Error()
			return
		}
		t.Log(err)
	}

	// &{0 []}
	if j, err := json.LoadFromJSONByte([]byte("{\"type\":\"array\", \"lanes\": []}")); err != nil {
		t.Log(err)
	} else {
		if m, err := MatrixFromJSON(j); err != nil {
			t.Error(err)
			return
		} else {
			t.Logf("%v", m)
		}
	}

	// &{0 []}
	if j, err := json.LoadFromJSONByte([]byte("{\"type\":\"array\", \"lanes\": [[]]}")); err != nil {
		t.Log(err)
	} else {
		if m, err := MatrixFromJSON(j); err != nil {
			t.Error(err)
			return
		} else {
			t.Logf("%v", m)
		}
	}

	// &{0 ["ff"]}
	if j, err := json.LoadFromJSONByte([]byte("{\"type\":\"array\", \"lanes\": [[\"ff\"]]}")); err != nil {
		t.Log(err)
	} else {
		if m, err := MatrixFromJSON(j); err != nil {
			t.Error(err)
			return
		} else {
			t.Logf("%v", m)
		}
	}

	// &{0 ["xx"]}
	if j, err := json.LoadFromJSONByte([]byte("{\"type\":\"array\", \"lanes\": [[\"xx\"]]}")); err != nil {
		t.Log(err)
	} else {
		if m, err := MatrixFromJSON(j); err != nil {
			//t.Error(err)
			//return
			t.Log(err)
		} else {
			//t.Logf("%v", m)
			t.Errorf("%v", m)
			return
		}
	}

	// &{0 ["ff"]}
	if j, err := json.LoadFromJSONByte([]byte("{\"type\":\"array\", \"lanes\": [[\"ff\", 1]]}")); err != nil {
		t.Log(err)
	} else {
		if m, err := MatrixFromJSON(j); err != nil {
			//t.Error(err)
			//return
			t.Log(err)
		} else {
			//t.Logf("%v", m)
			t.Errorf("%v", m)
			return
		}
	}
}
