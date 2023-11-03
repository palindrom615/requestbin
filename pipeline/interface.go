package pipeline

type CtxKey string
type CtxValue interface{}
type Ctx map[CtxKey]CtxValue

func NewCtx() *Ctx {
	ctx := make(Ctx)
	return &ctx
}

type Process *func(input Ctx) (output Ctx)

type Handler interface {
	Handle(input *Ctx) (output *Ctx)
	SetOutboundHandler(Handler)
}

type IdentityHandler struct {
	outboundHandler Handler
}

func (h *IdentityHandler) Handle(input *Ctx) (output *Ctx) {
	if h.outboundHandler != nil {
		return h.outboundHandler.Handle(input)
	} else {
		return input
	}
}

func (h *IdentityHandler) SetOutboundHandler(outboundHandler Handler) {
	if h == outboundHandler {
		panic("Cannot set outbound handler to self")
	}
	h.outboundHandler = outboundHandler
}
