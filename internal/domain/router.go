package domain

import (
	"errors"
	"strings"
)

type Router struct {
	routes map[string]CommandHandler
}

func NewRouter(routes map[string]CommandHandler) Router {
	return Router{
		routes: routes,
	}
}

func (r *Router) Route(msg *Message) (CommandHandler, error) {
	for signature, handler := range r.routes {
		if strings.HasPrefix(msg.Text, signature) {
			msg.Command = signature
			return handler, nil
		}
	}
	return nil, errors.New("handler not found")
}
