package artisanalinteger

type Engine interface {
	NextInt() (int64, error)
	LastInt() (int64, error)
	SetLastInt(int64) error
	SetKey(string) error
	SetOffset(int64) error
	SetIncrement(int64) error
	Close() error
}

type Service interface {
	NextInt() (int64, error)
	LastInt() (int64, error)
}

type Server interface {
	ListenAndServe(Service) error
	Address() string
}

type Client interface {
	NextInt() (int64, error)
}

type Integer struct {
	Integer int64 `json:"integer"`
}
