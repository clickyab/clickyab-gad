package config

import (
	"assert"
	"runtime"

	"time"

	"transport"

	"fmt"

	"github.com/fzerorubigd/expand"
	"gopkg.in/fzerorubigd/onion.v2"
	_ "gopkg.in/fzerorubigd/onion.v2/yamlloader" // config need this to load yaml file
)

const (
	organization = "clickyab"
	appName      = "gad"
)

//Config is the global application config instance
var Config AppConfig

var o = onion.New()

// AppConfig is the application config
type AppConfig struct {
	DevelMode       bool   `onion:"devel_mode"`
	CORS            bool   `onion:"cors"`
	MaxCPUAvailable int    `onion:"max_cpu_available"`
	MountPoint      string `onion:"mount_point"`

	Site  string
	Proto string

	Port       string
	StaticRoot string `onion:"static_root"`

	TimeZone string `onion:"time_zone"`

	Redis struct {
		Size     int
		Network  string
		Address  string
		Password string
		Days     int //Daily Statistic TimeOut Expiration TODO : the worst position for this
	}

	Mysql struct {
		RDSN              string `onion:"rdsn"`
		WDSN              string `onion:"wdsn"`
		MaxConnection     int    `onion:"max_connection"`
		MaxIdleConnection int    `onion:"max_idle_connection"`
	}

	AMQP struct {
		DSN        string
		Exchange   string
		Publisher  int
		ConfirmLen int
	}

	Select struct {
		Date    int `onion:"date"`
		Hour    int `onion:"hour"`
		Balance int `onion:"Balance"`
	}

	Page struct {
		PerPage    int `onion:"per_page"`
		MaxPerPage int `onion:"max_per_page"`
		MinPerPage int `onion:"min_per_page"`
	}

	Redmine struct {
		APIKey         string
		URL            string
		ProjectID      int `onion:"project_id"`
		Active         bool
		NewIssueTypeID int `onion:"new_issue_type_id"`
	}

	Slack struct {
		Channel    string
		Username   string
		WebHookURL string
		Active     bool
	}

	Clickyab struct {
		MaxLoadFail      int           `onion:"max_load_fail"`
		DefaultCTR       float64       `onion:"default_ctr"`
		CTRConst         []string      `onion:"ctr_const"`
		MinImp           int64         `onion:"min_imp"`
		MinFrequency     int           `onion:"min_frequency"`
		DailyImpExpire   time.Duration `onion:"daily_imp_expire"`
		DailyClickExpire time.Duration `onion:"daily_click_expire"`
		DailyCapExpire   time.Duration `onion:"daily_cap_expire"`
		MinCPMFloor      int64         `onion:"min_cpm_floor"`
		CopLen           int           `onion:"cop_len"`
	}
}

func defaultLayer() onion.Layer {
	d := onion.NewDefaultLayer()
	assert.Nil(d.SetDefault("site", "gad.loc"))
	assert.Nil(d.SetDefault("mount_point", "/"))
	assert.Nil(d.SetDefault("devel_mode", true))
	assert.Nil(d.SetDefault("cors", true))
	assert.Nil(d.SetDefault("max_cpu_available", runtime.NumCPU()))
	assert.Nil(d.SetDefault("proto", "http"))
	assert.Nil(d.SetDefault("port", ":80"))
	assert.Nil(d.SetDefault("time_zone", "Asia/Tehran"))
	p, err := expand.Path("$HOME/gad/statics")
	assert.Nil(err)
	assert.Nil(d.SetDefault("static_root", p))
	fmt.Println(p)

	assert.Nil(d.SetDefault("redis.size", 10))
	assert.Nil(d.SetDefault("redis.network", "tcp"))
	assert.Nil(d.SetDefault("redis.address", ":6379"))

	// TODO : move it to clickyab section
	assert.Nil(d.SetDefault("redis.days", 2))

	// TODO :  make sure ?parseTime=true is always set!
	assert.Nil(d.SetDefault("mysql.rdsn", "dev:cH3M7Z7I4sY8QP&ll130U&73&6KS$o@tcp(db-1.clickyab.ae:3306)/clickyab?charset=utf8&parseTime=true"))
	assert.Nil(d.SetDefault("mysql.wdsn", "root:bita123@tcp(127.0.0.1:3306)/clickyab?charset=utf8&parseTime=true"))
	assert.Nil(d.SetDefault("mysql.max_connection", 30))
	assert.Nil(d.SetDefault("mysql.max_idle_connection", 5))

	assert.Nil(d.SetDefault("amqp.publisher", 30))
	assert.Nil(d.SetDefault("amqp.exchange", "cy"))
	assert.Nil(d.SetDefault("amqp.dsn", "amqp://server:bita123@127.0.0.1:5672/"))
	assert.Nil(d.SetDefault("amqp.confirmlen", 50))

	assert.Nil(d.SetDefault("page.per_page", 10))
	assert.Nil(d.SetDefault("page.max_per_page", 100))
	assert.Nil(d.SetDefault("page.min_per_page", 1))

	assert.Nil(d.SetDefault("select.date", 0))
	assert.Nil(d.SetDefault("select.hour", 1))
	assert.Nil(d.SetDefault("select.balance", 50000))

	assert.Nil(d.SetDefault("clickyab.default_ctr", 0.1))
	assert.Nil(d.SetDefault("clickyab.ctr_const", []string{transport.AD_SLOT, transport.AD_WEBSITE, transport.CAMPAIGN, transport.CAMPAIGN_SLOT, transport.SLOT}))
	assert.Nil(d.SetDefault("clickyab.min_imp", 1000))
	assert.Nil(d.SetDefault("clickyab.min_frequency", 2))
	assert.Nil(d.SetDefault("clickyab.daily_imp_expire", "72h"))
	assert.Nil(d.SetDefault("clickyab.daily_click_expire", "72h"))
	assert.Nil(d.SetDefault("clickyab.daily_cap_expire", "72h"))
	assert.Nil(d.SetDefault("clickyab.min_cpm_floor", 150))
	assert.Nil(d.SetDefault("clickyab.max_load_fail", 3))
	assert.Nil(d.SetDefault("clickyab.cop_len", 10))

	return d
}
