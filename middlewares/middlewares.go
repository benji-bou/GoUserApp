package middlewares

import "net/http"

//MiddlewareHandler structure base for Middlewares managements
type MiddlewareHandler interface {
	Middle(w http.ResponseWriter, r *http.Request, next func())
}

//MiddlewareHandlerFunc type is an adapter to allow the use of
// ordinary functions as Middleware handlers.  If f is a function
// with the appropriate signature, MiddlewareHandlerFunc(f) is a
// Handler object that calls f.
type MiddlewareHandlerFunc func(w http.ResponseWriter, r *http.Request, next func())

//Middle call  f(w, r, next)
func (f MiddlewareHandlerFunc) Middle(w http.ResponseWriter, r *http.Request, next func()) {
	f(w, r, next)
}

//NewMiddlewares create pointer to Middlewares
func NewMiddlewares(hs ...MiddlewareHandler) *Middlewares {
	mdl := &Middlewares{}
	mdl.Use(hs...)
	return mdl
}

//NewMiddlewaresFunc create new Middleware with first Middleware
func NewMiddlewaresFunc(fs ...MiddlewareHandlerFunc) *Middlewares {
	mdl := &Middlewares{}
	for _, f := range fs {
		mdl.UseFunc(f)
	}
	return mdl
}

//Middlewares manage array of handler
type Middlewares struct {
	middles []MiddlewareHandler
}

//Use add a middleWareHandler to the queue
func (m *Middlewares) Use(newMiddles ...MiddlewareHandler) *Middlewares {
	m.middles = append(m.middles, newMiddles...)
	return m
}

//UseFunc add a Middleware as a function
func (m *Middlewares) UseFunc(newMiddleFunc MiddlewareHandlerFunc) *Middlewares {
	m.middles = append(m.middles, newMiddleFunc)
	return m
}

//UseHandler wrapp a basic http.Handler to use it as a middleware
func (m *Middlewares) UseHandler(handle http.Handler) *Middlewares {
	var wrapper MiddlewareHandlerFunc = func(w http.ResponseWriter, r *http.Request, next func()) {
		handle.ServeHTTP(w, r)
		next()
	}
	m.middles = append(m.middles, wrapper)
	return m
}

func (m *Middlewares) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.runHandler(w, r, 0)
}

func (m *Middlewares) runHandler(w http.ResponseWriter, r *http.Request, index int) {
	if index < len(m.middles) {
		// log.Println(tools.GetFunctionName(m.middles[index]))
		m.middles[index].Middle(w, r, func() {
			index++
			m.runHandler(w, r, index)
		})
	}
}
