package www

type Endpoints struct {
	Index  string
	Id     string
	Data   string
	Png    string
	Svg    string
	Spr    string
	Search string
}

type ErrorVars struct {
	Error     error
	Endpoints *Endpoints
}

type NotFoundVars struct {
	Endpoints *Endpoints
}
