package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type R1 struct {
	Param1 string
	Param2 []string
}

type R2 struct {
	Param1 string
	Param2 string
}

type KV struct {
	Key   string
	Value string
}

func main() {
	fmt.Println("=== Struct to String Lines ===")
	s := R1{
		Param1: "Val1",
		Param2: []string{"Val2", "Val3"},
	}
	fmt.Printf("Input:\n%+v\n", s)
	r1 := convToStr(s)
	fmt.Println("Output:")
	fmt.Println(r1)

	s2 := R2{
		Param1: "Val4",
		Param2: "Val5",
	}
	fmt.Printf("Input:\n%+v\n", s2)
	r2 := convToStr(s2)
	fmt.Println("Output:")
	fmt.Println(r2)

	fmt.Println("=== String Lines to JSON ===")
	fmt.Printf("Input:\n%+v\n", r1)
	var rj1 R1
	if err := decodeJSON(r1, &rj1); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Output:")
	RenderJson(rj1)

	fmt.Println()

	fmt.Printf("Input:\n%+v\n", r2)
	var rj2 R2
	if err := decodeJSON(r2, &rj2); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Output:")
	RenderJson(rj2)
}

func convToStr(s interface{}) string {
	var res string
	r := reflect.ValueOf(s)
	for i := 0; i < r.NumField(); i++ {
		t := r.Type().Field(i)
		f := r.Field(i)
		switch f.Kind() {
		case reflect.String:
			res += fmt.Sprintf("%s:%s\n", t.Name, f)
		case reflect.Slice:
			for i := 0; i < f.Len(); i++ {
				res += fmt.Sprintf("%s:%s\n", t.Name, f.Index(i))
			}
		default:
			fmt.Printf("not support kind %s", f.Kind())
		}
	}
	return res
}

func RenderJson(v any) error {
	e := json.NewEncoder(os.Stdout)
	// &, <, >などをエスケープ処理したい
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	return e.Encode(v)
}

func decodeJSON(input string, payload interface{}) error {
	params := make(map[string]interface{})
	responseLines := strings.Split(input, "\n")

	var kv []KV
	for _, l := range responseLines {
		// ":"を区切り文字としてKey, Valueを抽出
		idx := strings.Index(l, ":")
		if idx == -1 {
			continue
		}
		k := l[:idx]
		v := chop(l[idx+1:])
		kv = append(kv, KV{Key: k, Value: v})
	}
	fmt.Printf("KV: %+v\n", kv)

	r := reflect.TypeOf(payload).Elem()
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		vs := collectValues(kv, f.Name)
		if len(vs) == 0 {
			continue
		}
		if f.Type == reflect.TypeOf([]string{}) {
			params[f.Name] = vs
		} else {
			params[f.Name] = vs[0]
		}
	}
	js, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return json.Unmarshal(js, payload)
}

// 文末の改行コード削除
func chop(s string) string {
	s = strings.TrimRight(s, "\n")
	if strings.HasSuffix(s, "\r") {
		s = strings.TrimRight(s, "\r")
	}

	return s
}

// KVから同じKeyの値を集めて返す
func collectValues(kv []KV, key string) []string {
	var v []string
	for _, e := range kv {
		if e.Key == key {
			v = append(v, e.Value)
		}
	}
	return v
}
