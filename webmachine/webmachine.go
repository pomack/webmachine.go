package webmachine

import (
    "log"
    "net/http"
)

func NewWebMachine() WebMachine {
    return new(webMachine)
}

func (p *webMachine) AddRouteHandler(handler RouteHandler) {
    p.routeHandlers = append(p.routeHandlers, handler)
}

func (p *webMachine) RemoveRouteHandler(handler RouteHandler) {
    for i, h := range p.routeHandlers {
        if h == handler {
            handlers := make([]RouteHandler, 0, cap(p.routeHandlers))
            copy(handlers, p.routeHandlers[0:i])
            copy(handlers, p.routeHandlers[i+1:])
            p.routeHandlers = handlers
            break
        }
    }
}

func (p *webMachine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
    r := NewRequestFromHttpRequest(req)
    rs := NewResponseWriter(resp)
    log.Print("running URL: ", r.URL().Path, "\n")
    for _, h := range p.routeHandlers {
        rh := h.(RouteHandler)

        if handler := rh.HandlerFor(r, rs); handler != nil {
            log.Print("found route handler for: ", r.URL().Path, " ", handler, "\n")
            handleRequest(handler, r, rs)
            return
        }
    }
    log.Print("no route handlers matched: ", r.URL().Path, "\n")
    resp.WriteHeader(http.StatusBadRequest)
}
