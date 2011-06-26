package webmachine

import (
  "http"
  "log"
)


func NewWebMachine() WebMachine {
  return new(webMachine)
}

func (p *webMachine) AddRouteHandler(handler RouteHandler) {
  p.routeHandlers.Push(handler)
}

func (p *webMachine) RemoveRouteHandler(handler RouteHandler) {
  for i, h := range p.routeHandlers {
    if h == handler {
      p.routeHandlers.Delete(i)
      break
    }
  }
}

func (p *webMachine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
  r := NewRequestFromHttpRequest(req)
  log.Print("running URL: ", r.URL().Path, "\n")
  for _, h := range p.routeHandlers {
    rh := h.(RouteHandler)
    if handler := rh.HandlerFor(r, resp); handler != nil {
      log.Print("found route handler for: ", r.URL().Path, " ", handler, "\n")
      handleRequest(handler, r, NewResponseWriter(resp))
      return
    }
  }
  log.Print("no route handlers matched: ", r.URL().Path, "\n")
  resp.WriteHeader(http.StatusBadRequest)
}
