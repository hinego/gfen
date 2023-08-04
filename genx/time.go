package genx

import "github.com/gogf/gf/v2/os/gtime"

type Time struct {
	*gtime.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	return t.Time.UnmarshalJSON(b)
}
