package pipeline

import "testing"

type AddIntHandler struct {
	a int
}

const key = CtxKey("key")

func (a *AddIntHandler) Handle(input *Ctx) (output *Ctx) {
	value := (*input)[key]
	switch value.(type) {
	case int:
		(*input)[key] = value.(int) + a.a
	default:
		panic("value is not int")
	}
	return input
}

type MultipleIntHandler struct {
	m int
}

func (m *MultipleIntHandler) Handle(input *Ctx) (output *Ctx) {
	value := (*input)[key]
	switch value.(type) {
	case int:
		(*input)[key] = value.(int) * m.m
	default:
		panic("value is not int")
	}
	return input
}

func TestComposeHandler_Handle(t *testing.T) {
	// arange
	plus1 := AddIntHandler{1}
	times2 := MultipleIntHandler{2}
	compose := NewComposeHandler(&plus1, &times2)
	input := NewCtx()
	(*input)[key] = 1

	// act
	res := compose.Handle(input)

	// assert
	if (*res)[key] != 4 {
		t.Errorf("ComposeHandler.Handle() = %v, want %v", (*res)[key], 4)
	}

	compose = NewComposeHandler(&times2, &plus1)
	input = NewCtx()
	(*input)[key] = 1
	res = compose.Handle(input)
	if (*res)[key] != 3 {
		t.Errorf("ComposeHandler.Handle() = %v, want %v", (*res)[key], 3)
	}
}
