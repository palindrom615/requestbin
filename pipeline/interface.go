package pipeline

type CtxKey string
type CtxValue interface{}
type Ctx map[CtxKey]CtxValue

func NewCtx() *Ctx {
	ctx := make(Ctx)
	return &ctx
}

type Handler interface {
	Handle(input *Ctx) (output *Ctx)
}

type IdentityHandler struct {
}

func (h *IdentityHandler) Handle(input *Ctx) (output *Ctx) {
	return input
}
