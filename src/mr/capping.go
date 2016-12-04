package mr

import (
	"sync"

	"github.com/labstack/echo"
)

// Capping is the structure for capping
type capping struct {
	View      int
	Frequency int
	Selected  bool
}

// CappingInterface interface capping
type CappingInterface interface {
	GetView() int
	GetFrequency() int
	GetCapping() int
	IncView(int)
	GetSelected() bool
}

// CappingLocker is the safe capping counter for min capping
type CappingLocker struct {
	cap  int
	lock sync.RWMutex
	data []*MinAdData
}

const (
	cappingCtx = "__capping_context__"
)

// NewCapping create new capping
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

func (c *capping) IncView(a int) {
	c.View += a
	c.Selected = true
}

func (c *capping) GetSelected() bool {
	return c.Selected
}

// Set the capping value
func (c *CappingLocker) Set(i int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.cap > 0 { // Do not allow set capping if already set
		return
	}
	c.cap = i
}

// Get the capping value
func (c *CappingLocker) Get() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.cap
}

// GetData return the slice
func (c *CappingLocker) GetData() []*MinAdData {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.data
}

// Len of the slice
func (c *CappingLocker) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return len(c.data)
}

// Append to slice
func (c *CappingLocker) Append(m *MinAdData) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data = append(c.data, m)
}
