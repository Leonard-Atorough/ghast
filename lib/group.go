package ghast

type routeGroup struct {
	prefix      string
	middlewares []Middleware
	router      Router
}
