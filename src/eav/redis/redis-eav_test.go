package redis

import (
	"config"
	"services/redis"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	config.Initialize()
	aredis.Initialize()
	Convey("Test keyvalue store for redis", t, func() {
		store := newRedisEAVStore("test_key")
		So(store.Key(), ShouldEqual, "test_key")
		Convey("check set and get", func() {
			store.SetSubKey("test", "test_val")
			So(store.SubKey("test"), ShouldEqual, "test_val")
			store.SetSubKey("another", "2")
			So(store.SubKey("another"), ShouldEqual, "2")
			So(store.Save(time.Hour), ShouldBeNil)
			So(store.AllKeys(), ShouldResemble, map[string]string{"test": "test_val", "another": "2"})
		})

		another := newRedisEAVStore("test_key")
		So(another.SubKey("another"), ShouldEqual, "2")
		So(another.AllKeys(), ShouldResemble, map[string]string{"test": "test_val", "another": "2"})
	})
}
