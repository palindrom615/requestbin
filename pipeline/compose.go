package pipeline

type ComposeHandler struct {
	handlers []Handler
}

func NewComposeHandler(handlers ...Handler) *ComposeHandler {
	return &ComposeHandler{handlers}
}

func (h *ComposeHandler) Handle(input *Ctx) (output *Ctx) {
	for _, handler := range h.handlers {
		input = handler.Handle(input)
	}
	return input
}
