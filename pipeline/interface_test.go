package pipeline

import "testing"

func TestIdentityHandler_Handle(t *testing.T) {
	// arange
	i := IdentityHandler{}
	input := NewCtx()
	key := CtxKey("key")
	value := CtxValue("value")
	(*input)[key] = value

	// act
	output := i.Handle(input)

	// assert
	v, e := (*input)[key]
	if !e {
		t.Fatal("Context key not found")
	}
	vo, eo := (*output)[key]
	if !eo || v != vo {
		t.Errorf("IdentityHandler.Handle() = %v, want %v", v, vo)
	}
}
