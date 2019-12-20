package uid

type Provider interface {
	UID() (int64, error)
}
