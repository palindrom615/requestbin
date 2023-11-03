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

func TestIdentityHandler_SetOutboundHandler(t *testing.T) {
	// arange
	i := &IdentityHandler{}
	i2 := &IdentityHandler{}
	input := &Ctx{"key": "value"}

	// act
	i.SetOutboundHandler(i2)
	output := i.Handle(input)

	// assert
	if input != output {
		t.Errorf("IdentityHandler.Handle() = %v, want %v", output, input)
	}
}

func TestIdentityHandler_SetOutboundHandler_self(t *testing.T) {
	// arrange
	i := &IdentityHandler{}

	// assert
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("IdentityHandler.SetOutboundHandler() did not panic")
		}
	}()

	// act
	i.SetOutboundHandler(i)
}
