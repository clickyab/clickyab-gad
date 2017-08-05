package mr

// Capping is the structure for capping
type capping struct {
	View      int
	Ads       map[int64]int
	Sizes     map[int64]int
	Frequency int
	Selected  bool
	Target    bool
}

// CappingInterface interface capping
type CappingInterface interface {
	// GetView return the view of this campaign for this user
	GetView() int
	// GetView return the view of this campaign for this user
	GetAdView(int64) int
	// GetFrequency return the frequency for this user
	GetFrequency() int
	// GetCapping return the frequency capping value, the view/frequency
	GetCapping() int
	// GetCapping return the frequency capping value, the view/frequency
	GetAdCapping(int64) int
	// IncView increase the vie
	IncView(int64, int, bool)
	// GetSelected return if this campaign is already selected in this batch
	GetSelected() bool
	// IsTargeted return if targeted
	IsTargeted() bool
}

// CappingContext is the type used to handle capping locker
type CappingContext map[int64]CappingInterface

// NewCapping create new capping
func (caps CappingContext) NewCapping(cpID int64, view, freq int, target bool) CappingInterface {
	if _, ok := caps[cpID]; !ok {
		caps[cpID] = &capping{
			View:      view,
			Frequency: freq,
			Target:    target,
			Ads:       make(map[int64]int),
			Sizes:     make(map[int64]int),
		}
	}

	return caps[cpID]
}

func (c *capping) GetView() int {
	return c.View
}

func (c *capping) GetAdView(ad int64) int {
	return c.Ads[ad]
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

func (c *capping) GetAdCapping(ad int64) int {
	return c.Ads[ad] / c.Frequency
}

func (c *capping) IncView(ad int64, a int, sel bool) {
	c.View += a
	c.Ads[ad] += a
	if sel {
		c.Selected = true
	}
}

func (c *capping) GetSelected() bool {
	return c.Selected
}

func (c *capping) IsTargeted() bool {
	return c.Target
}
