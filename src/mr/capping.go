package mr

import (
	"sync"

	"gopkg.in/labstack/echo.v3"
)

// Capping is the structure for capping
type capping struct {
	View      int
	Frequency int
	Selected  bool
	Target    bool
}

// CappingInterface interface capping
type CappingInterface interface {
	// GetView return the view of this campaign for this user
	GetView() int
	// GetFrequency return the frequency for this user
	GetFrequency() int
	// GetCapping return the frequency capping value, the view/frequency
	GetCapping() int
	// IncView increase the vie
	IncView(int)
	// GetSelected return if this campaign is already selected in this batch
	GetSelected() bool
	// IsTargeted return if the current campaign is targeted for this user?
	IsTargeted() bool
}

// CappingLocker is the safe capping counter for min capping
type CappingLocker struct {
	cap  int
	lock sync.RWMutex
	data []*AdData
}

const (
	cappingCtx = "__capping_context__"
)

// NewCapping create new capping
func NewCapping(ctx echo.Context, cpID int64, view, freq int, target bool) CappingInterface {
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
			Target:    target,
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
	if c.Target {
		return 0
	}
	return c.View / c.Frequency
}

func (c *capping) IncView(a int) {
	c.View += a
	c.Selected = true
}

func (c *capping) GetSelected() bool {
	return c.Selected
}

func (c *capping) IsTargeted() bool {
	return c.Target
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
func (c *CappingLocker) GetData() []*AdData {
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
func (c *CappingLocker) Append(m *AdData) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data = append(c.data, m)
}
