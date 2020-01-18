package uid

import (
	"context"
	"time"
)

func init() {
	ctx := context.Background()
	pr := NewYMDProvider()
	RegisterProvider(ctx, "ymd", pr)
}

type YMDProvider struct {
	Provider
}

type YMDUID struct {
	UID
	date time.Time
}

func NewYMDProvider() Provider {

	pr := &YMDProvider{}
	return pr
}

func (pr *YMDProvider) Open(ctx context.Context, uri string) error {
	return nil
}

func (pr *YMDProvider) UID(args ...interface{}) (UID, error) {

	date := time.Now()

	if len(args) == 1 {

		str_date := args[0].(string)

		t, err := time.Parse("20060102", str_date)

		if err != nil {
			return nil, err
		}

		date = t
	}

	return NewYMDUID(date)
}

func NewYMDUID(date time.Time) (UID, error) {

	u := &YMDUID{
		date: date,
	}

	return u, nil
}

func (u *YMDUID) String() string {

	return u.date.Format("20060102")
}
