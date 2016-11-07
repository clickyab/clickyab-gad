package config

import (
	"assert"
	"runtime"

	"time"

	"transport"

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
		DSN               string `onion:"dsn"`
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
		DefaultCTR       float64       `onion:"default_ctr"`
		CTRConst         []string      `onion:"ctr_const"`
		MinImp           int64         `onion:"min_imp"`
		MinFrequency     int           `onion:"min_frequency"`
		DailyImpExpire   time.Duration `onion:"daily_imp_expire"`
		DailyClickExpire time.Duration `onion:"daily_click_expire"`
		DailyCapExpire   time.Duration `onion:"daily_cap_expire"`
		MinCPMFloor      int64         `onion:"min_cpm_floor"`
	}
}

func defaultLayer() onion.Layer {
	d := onion.NewDefaultLayer()
	d.SetDefault("site", "gad.loc")
	d.SetDefault("mount_point", "/")
	d.SetDefault("devel_mode", true)
	d.SetDefault("cors", true)
	d.SetDefault("max_cpu_available", runtime.NumCPU())
	d.SetDefault("proto", "http")
	d.SetDefault("port", ":80")
	d.SetDefault("time_zone", "Asia/Tehran")
	p, err := expand.Path("/statics")
	assert.Nil(err)
	d.SetDefault("static_root", p)

	d.SetDefault("redis.size", 10)
	d.SetDefault("redis.network", "tcp")
	d.SetDefault("redis.address", ":6379")

	// TODO : move it to clickyab section
	d.SetDefault("redis.days", 2)

	// TODO :  make sure ?parseTime=true is always set!
	d.SetDefault("mysql.dsn", "dev:cH3M7Z7I4sY8QP&ll130U&73&6KS$o@tcp(db-2.clickyab.ae:3306)/clickyab?charset=utf8&parseTime=true")
	d.SetDefault("mysql.max_connection", 30)
	d.SetDefault("mysql.max_idle_connection", 5)

	d.SetDefault("amqp.publisher", 30)
	d.SetDefault("amqp.exchange", "cy")
	d.SetDefault("amqp.dsn", "amqp://server:bita123@127.0.0.1:5672/")
	d.SetDefault("amqp.confirmlen", 50)

	d.SetDefault("page.per_page", 10)
	d.SetDefault("page.max_per_page", 100)
	d.SetDefault("page.min_per_page", 1)

	d.SetDefault("select.date", 0)
	d.SetDefault("select.hour", 1)
	d.SetDefault("select.balance", 50000)

	d.SetDefault("clickyab.default_ctr", 0.1)
	d.SetDefault("clickyab.ctr_const", []string{transport.AD_SLOT, transport.AD_WEBSITE, transport.CAMPAIGN, transport.CAMPAIGN_SLOT, transport.SLOT})
	d.SetDefault("clickyab.min_imp", 1000)
	d.SetDefault("clickyab.min_frequency", 2)
	d.SetDefault("clickyab.daily_imp_expire", "72h")
	d.SetDefault("clickyab.daily_click_expire", "72h")
	d.SetDefault("clickyab.daily_cap_expire", "72h")
	d.SetDefault("clickyab.min_cpm_floor", 150)

	return d
}
