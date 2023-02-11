package pitstop

import (
	"github.com/valyala/fasthttp"
)

type Constructor func(fasthttp.RequestHandler) fasthttp.RequestHandler

type Chain struct {
	constructors []Constructor
}

func NewChain(constructors ...Constructor) Chain {
	return Chain{append(([]Constructor)(nil), constructors...)}
}

func (c Chain) Then(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	if h == nil {
		h = func(ctx *fasthttp.RequestCtx) {}
	}

	for i := range c.constructors {
		h = c.constructors[len(c.constructors)-1-i](h)
	}

	return h
}
func (c Chain) ThenFunc(fn func(ctx *fasthttp.RequestCtx)) fasthttp.RequestHandler {
	if fn == nil {
		return c.Then(nil)
	}
	return c.Then(fn)
}
func (c Chain) Append(constructors ...Constructor) Chain {
	newCons := make([]Constructor, 0, len(c.constructors)+len(constructors))
	newCons = append(newCons, c.constructors...)
	newCons = append(newCons, constructors...)

	return Chain{newCons}
}
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.constructors...)
}
