package main

import "fmt"

type HelloInterface interface {
	Hello()
}

type Foo struct{}

func (f *Foo) Greeting(h HelloInterface) {
	h.Hello()
}

type TypeA struct{}
type TypeB struct{}

var _ HelloInterface = TypeA{}

//var _ HelloInterface = TypeB{}

func NewTypeA() HelloInterface {
	return TypeA{}
}

func (a TypeA) Hello() {}

type Person interface {
	// 敬称
	Title() string
	// 名前
	Name() string
}

type person struct {
	firstName string
	lastName  string
}

func (p *person) Name() string {
	return fmt.Sprintf("%s %s", p.firstName, p.lastName)
}

type Gender int

const (
	Female = iota
	Male
)

type female struct {
	*person
}

func (f *female) Title() string {
	return "Ms."
}

type male struct {
	*person
}

func (m *male) Title() string {
	return "Mr."
}

func NewPerson(gender Gender, firstName, lastName string) Person {
	p := &person{firstName, lastName}

	if gender == Female {
		return &female{p}
	}

	return &male{p}
}

func main() {
	p := NewPerson(Male, "Taro", "Yamada")
	fmt.Println(p.Title(), p.Name())
}
