// package www implements HTTP handlers for the whosonfirst-browser web application.
package www

type ErrorVars struct {
	Error error
	Paths *Paths
}

type NotFoundVars struct {
	Paths *Paths
}
