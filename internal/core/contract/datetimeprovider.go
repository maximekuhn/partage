package contract

import "time"

type DatetimeProvider interface {
	Provide() time.Time
}
