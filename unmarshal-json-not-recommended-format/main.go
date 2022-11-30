package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type S struct {
	List List `json:"list"`
}

type List struct {
	ListItem []Item `json:"listItem"`
}

type Item struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
	Buz string `json:"buz"`
}

func (l *List) UnmarshalJSON(b []byte) error {
	// 同一KeyのオブジェクトはUnmarshalできないのでKeyをユニークな値に置換する
	k := "listItem"
	c := strings.Count(string(b), k)
	rep := strings.Replace(string(b), k, "replaced00", 1)
	for i := 0; i < c; i++ {
		rep = strings.Replace(string(rep), k, fmt.Sprintf("replaced%02d", i+1), 1)
	}

	var itf map[string]interface{}
	if err := json.Unmarshal(b, &itf); err != nil {
		return err
	}
	fmt.Println(itf)

	var items []Item
	for _, v := range itf {
		rv := reflect.ValueOf(v)
		foo := rv.MapIndex(reflect.ValueOf("foo"))
		bar := rv.MapIndex(reflect.ValueOf("bar"))
		buz := rv.MapIndex(reflect.ValueOf("buz"))
		item := Item{
			Foo: fmt.Sprintf("%v", foo),
			Bar: fmt.Sprintf("%v", bar),
			Buz: fmt.Sprintf("%v", buz),
		}
		items = append(items, item)
	}
	l.ListItem = items

	return nil
}

func main() {
	input := []byte(`
  {
    "list": {
      "listItem": {
        "foo": "1",
        "bar": "2",
        "buz": "3"
      },
      "listItem": {
        "foo": "4",
        "bar": "5",
        "buz": "6"
      },
      "listItem": {
        "foo": "7",
        "bar": "8",
        "buz": "9"
      }
    }
  }`)

	var s S
	if err := json.Unmarshal(input, &s); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v", s)
}
