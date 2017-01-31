package rtb

import (
	"eav/redis"
	"entity"
	"fmt"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	userCapKey string = "CAP"
	adCapKey          = "AD"
)

// CappingContext is the type used to handle capping locker
type cappingContext map[int64]entity.Capping

// Capping is the structure for capping
type capping struct {
	View      int
	Ads       map[int64]int
	Sizes     map[int64]int
	Frequency int
	Selected  bool
	Target    bool
}

// NewCapping create new capping
func (caps cappingContext) NewCapping(cpID int64, view, freq int, target bool) entity.Capping {
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

func getCappingKey(copID int64) string {
	return fmt.Sprintf(
		"%s_%d_%s",
		userCapKey,
		copID,
		time.Now().Format("060102"),
	)
}

func getCapping(clientID int64, ads map[int][]entity.Advertise, slots []entity.Slot) entity.Capping {
	kiwi := redis.NewRedisEAVStore(getCappingKey(clientID))
	caps := kiwi.GetAllKeys()
	c := make(cappingContext)
	for s := range slots {
		found := false
		// First check for ads one by one
		keys := []string{}
		for a := range ads[slots[s].Size()] {
			capKey := fmt.Sprintf("%s_%d", adCapKey, ads[slots[s]][a].ID())
			// check if the cap for this ad is full
			view, _ := strconv.ParseInt(caps[capKey], 10, 0)
			cp := view / ads[slots[s]][a].Campaign().Frequency()
			if cp < 1 {
				// ok there is at least one that the capping is less than 1
				found = true
				break
			}
			// Set it to zero, but not save it yet
			keys = append(keys, capKey)
		}
		// No ad pass the cap
		if !found {
			logrus.Debugf("Removing key for size %d", slots[s].Size())
			for i := range keys {
				kiwi.SubKey(keys[i], "0")
				caps[keys[i]] = "0"
			}
		}

		for a := range ads[slots[s].Size()] {
			capKey := fmt.Sprintf("%s_%d", adCapKey, ads[slots[s]][a].ID())
			view, _ := strconv.ParseInt(caps[capKey], 10, 0)

			capp := c.NewCapping(
				ads[slots[s]][a].Campaign().ID(),
				0,
				ads[slots[s]][a].Campaign().ID(),
				false,
			)
			capp.IncView(ads[slots[s]][a].ID(), view, false)
			ads[slots[s]][a].SetCapping(capp)
		}
		// todo : sort by cap
	}
}
