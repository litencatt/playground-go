package main

import (
	"encoding/json"
	"log"
)

const inputA = `
{
	"key": "typeA",
	"value": "bbb"
}
`
const inputB = `
{
	"key": "typeB",
	"list": [
		"bbb",
		"ccc"
	]
}
`
const inputC = `
{
	"key": "typeC",
	"list": [
		{
			"key1": "value1",
			"key2": "value2"
		}
	]
}
`
const input = `
{
  "list": [
    {
      "key": "A",
      "value": "aaa"
    },
    {
      "key": "B",
      "variable": [
        "b-1",
        "b-2"
      ]
    },
    {
      "key": "C",
      "variable": [
        {
          "foo": "1",
          "bar": "2"
        },
        {
          "foo": "3",
          "bar": "4"
        }
      ]
    }
  ]
}
`

type Input struct {
	List []Element `json:"list"`
}

type Element struct {
	Key          string      `json:"key"`
	Value        string      `json:"value,omitempty"`
	VariablePart interface{} `json:"variable,omitempty"`
}

type ObjectVariable struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

// 原則としてBはjson内に1つしか含まれない
func (i *Input) GetVariableB() []string {
	as := []string{}
	for _, e := range i.List {
		if e.Key != "B" {
			continue
		}
		switch e.VariablePart.(type) {
		case *[]string:
			as = *e.VariablePart.(*[]string)
		}
		break
	}
	return as
}

// 原則としてCはjson内に1つしか含まれない
func (i *Input) GetVariableC() []ObjectVariable {
	ov := []ObjectVariable{}
	for _, e := range i.List {
		if e.Key != "C" {
			continue
		}
		switch e.VariablePart.(type) {
		case *[]ObjectVariable:
			ov = *e.VariablePart.(*[]ObjectVariable)
		}
		break
	}
	return ov
}

func (e *Element) UnmarshalJSON(b []byte) error {
	// log.Println("unmarshal json is called")

	type Alias Element
	a := &struct {
		// 可変部分をjson.RawMessageにしておく
		VariablePart json.RawMessage `json:"variable"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	// 一旦全体をUnmarshal
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}

	// Unmarshal後のKeyに応じてVariablePartのUnmarshalする構造体を条件分岐
	switch e.Key {
	case "B":
		var s []string
		if err := json.Unmarshal(a.VariablePart, &s); err != nil {
			log.Fatal(err)
		}
		e.VariablePart = &s
		// log.Printf("typeB List: %+v", l.VariablePart)
	case "C":
		var s []ObjectVariable
		if err := json.Unmarshal(a.VariablePart, &s); err != nil {
			log.Fatal(err)
		}
		e.VariablePart = &s
		// log.Printf("typeC List: %+v", l.VariablePart)
	default:
		return nil
	}

	return nil
}

func main() {
	// 行数表示
	log.SetFlags(log.Lshortfile)

	log.Println(string(input))

	// inputのUnmarshal
	var i Input
	if err := json.Unmarshal([]byte(input), &i); err != nil {
		log.Fatal(err)
	}

	// Unmarshalした可変部の値を取得
	ss := i.GetVariableB()
	for _, s := range ss {
		log.Printf("%q", s)
	}

	lo := i.GetVariableC()
	for _, o := range lo {
		log.Printf("%+v", o)
	}
}
