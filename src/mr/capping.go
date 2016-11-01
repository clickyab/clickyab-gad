package mr

import "github.com/labstack/echo"

// Capping is the structure for capping
type capping struct {
	View      int
	Frequency int
}

type CappingInterface interface {
	GetView() int
	GetFrequency() int
	GetCapping() int
	IncView()
}

const (
	cappingCtx = "__capping_context__"
)

func NewCapping(ctx echo.Context, cpID int64, view, freq int) CappingInterface {
	var caps map[int64]*capping
	var ok bool
	if caps, ok = ctx.Get(cappingCtx).(map[int64]*capping); !ok {
		caps = make(map[int64]*capping)
		ctx.Set(cappingCtx, caps)
	}
	if _, ok := caps[cpID]; !ok {
		caps[cpID] = &capping{
			View:      view,
			Frequency: freq,
		}
	}

	return caps[cpID]
}

func (c *capping) GetView() int {
	return c.View
}

func (c *capping) GetFrequency() int {
	return c.Frequency
}

func (c *capping) GetCapping() int {
	return c.View / c.Frequency
}

func (c *capping) IncView() {
	c.View++
}
