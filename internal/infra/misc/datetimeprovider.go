package misc

import "time"

type DatetimeProviderProd struct{}

func (d *DatetimeProviderProd) Provide() time.Time {
	return time.Now()
}
