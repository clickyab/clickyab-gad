package rtb

import (
	"eav"
	"entity"
	"fmt"
	"sort"
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
	view      int
	ads       map[int64]int
	sizes     map[int64]int
	frequency int
	selected  bool
	target    bool
}

// NewCapping create new capping
func (caps cappingContext) NewCapping(cpID int64, view, freq int, target bool) entity.Capping {
	if _, ok := caps[cpID]; !ok {
		caps[cpID] = &capping{
			view:      view,
			frequency: freq,
			target:    target,
			ads:       make(map[int64]int),
			sizes:     make(map[int64]int),
		}
	}

	return caps[cpID]
}

func (c *capping) View() int {
	return c.view
}

func (c *capping) AdView(ad int64) int {
	return c.ads[ad]
}

func (c *capping) Frequency() int {
	return c.frequency
}

func (c *capping) Capping() int {
	if c.target {
		return 0
	}
	return c.view / c.frequency
}

func (c *capping) AdCapping(ad int64) int {
	return c.ads[ad] / c.frequency
}

func (c *capping) IncView(ad int64, a int, sel bool) {
	c.view += a
	c.ads[ad] += a
	if sel {
		c.selected = true
	}
}

func (c *capping) Selected() bool {
	return c.selected
}

func (c *capping) IsTargeted() bool {
	return c.target
}

func getCappingKey(copID int64) string {
	return fmt.Sprintf(
		"%s_%d_%s",
		userCapKey,
		copID,
		time.Now().Format("060102"),
	)
}

func getCapping(clientID int64, ads map[int][]entity.Advertise, slots []entity.Slot) map[int][]entity.Advertise {
	kiwi := eav.NewEavStore(getCappingKey(clientID))
	caps := kiwi.AllKeys()
	c := make(cappingContext)
	for s := range slots {
		size := slots[s].Size()
		found := false
		// First check for ads one by one
		keys := []string{}
		for a := range ads[slots[s].Size()] {
			capKey := fmt.Sprintf("%s_%d", adCapKey, ads[size][a].ID())
			// check if the cap for this ad is full
			view, _ := strconv.ParseInt(caps[capKey], 10, 0)
			cp := int(view) / ads[size][a].Campaign().Frequency()
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
				kiwi.SetSubKey(keys[i], "0")
				caps[keys[i]] = "0"
			}
		}

		for a := range ads[slots[s].Size()] {
			capKey := fmt.Sprintf("%s_%d", adCapKey, ads[size][a].ID())
			view, _ := strconv.ParseInt(caps[capKey], 10, 0)

			capp := c.NewCapping(
				ads[size][a].Campaign().ID(),
				0,
				ads[size][a].Campaign().Frequency(),
				false,
			)
			capp.IncView(ads[size][a].ID(), int(view), false)
			ads[size][a].SetCapping(capp)
		}
		// todo : sort by cap
		sortCap := entity.SortByCap(ads[size])
		sort.Sort(sortCap)
		ads[size] = []entity.Advertise(sortCap)
	}

	return ads
}
