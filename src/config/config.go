package config

import (
	"assert"
	"fmt"
	"runtime"
	"time"

	"os"
	"strconv"

	"regexp"

	"github.com/fzerorubigd/expand"
	onion "gopkg.in/fzerorubigd/onion.v2"
	_ "gopkg.in/fzerorubigd/onion.v2/yamlloader" // config need this to load yaml file
)

const (
	organization = "clickyab"
	appName      = "gad"
)

var redisPattern = regexp.MustCompile("^redis://([^:]+):([^@]+)@([^:]+):([0-9]+)$")

//Config is the global application config instance
var Config AppConfig

var o = onion.New()

// AppConfig is the application config
type AppConfig struct {
	DevelMode       bool          `onion:"devel_mode"`
	CORS            bool          `onion:"cors"`
	MaxCPUAvailable int           `onion:"max_cpu_available"`
	MountPoint      string        `onion:"mount_point"`
	ServerID        string        `onion:"server_id"`
	Retry           time.Duration `onion:"retry"`

	Site  string
	Proto string

	MachineName string `onion:"machine_name"`
	Port        int
	StaticRoot  string `onion:"static_root"`

	TimeZone string `onion:"time_zone"`

	Redis struct {
		Size     int
		Network  string
		Address  string
		Password string
		Databse  int
		Days     int //Daily Statistic TimeOut Expiration TODO : the worst position for this
	}

	Mysql struct {
		RDSNSlice         string `onion:"rdsn"`
		WDSN              string `onion:"wdsn"`
		MaxConnection     int    `onion:"max_connection"`
		MaxIdleConnection int    `onion:"max_idle_connection"`
	}

	AMQP struct {
		DSN        string
		Exchange   string
		Publisher  int
		ConfirmLen int
		Debug      bool
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
		MaxLoadFail int     `onion:"max_load_fail"`
		DefaultCTR  float64 `onion:"default_ctr"`
		//CTRConst         []string      `onion:"ctr_const"`
		MinImp           int64         `onion:"min_imp"`
		MinFrequency     int           `onion:"min_frequency"`
		DailyImpExpire   time.Duration `onion:"daily_imp_expire"`
		DailyClickExpire time.Duration `onion:"daily_click_expire"`
		DailyCapExpire   time.Duration `onion:"daily_cap_expire"`
		MegaImpExpire    time.Duration `onion:"mega_imp_expire"`
		MinCPMFloorWeb   int64         `onion:"min_cpm_floor_web"`
		MinCPMFloorApp   int64         `onion:"min_cpm_floor_app"`
		CopLen           int           `onion:"cop_len"`
		FastClick        int64         `onion:"fast_click"`
		ConvDelay        time.Duration `onion:"conv_delay"`
		ConvRetry        int64         `onion:"conv_retry"`
		Native           struct {
			MaxCount int `onion:"max_count"`
		}
		Vast struct {
			DefaultDuration string `onion:"default_duration"`
			DefaultSkipOff  string `onion:"default_skipoff"`
		}
		AdCTREffect   int64 `onion:"ad_ctr_effect"`
		SlotCTREffect int64 `onion:"slot_ctr_effect"`
		UnderFloor    bool  `onion:"under_floor"`

		RetargettingHour int `onion:"retargetting_hour"`

		MinCPCWeb    int64 `onion:"min_cpc_web"`
		MinCPCNative int64 `onion:"min_cpc_native"`
		MinCPCApp    int64 `onion:"min_cpc_app"`
		MinCPCVast   int64 `onion:"min_cpc_vast"`

		FakeSupplier string `onion:"fake_supplier"`

		FloorDiv struct {
			Web    int64 `onion:"web"`
			App    int64 `onion:"app"`
			Native int64 `onion:"native"`
			Vast   int64 `onion:"vast"`
			Demand int64 `onion:"demand"`
		} `onion:"floor_div"`
	}

	Fluentd struct {
		Host       string
		Port       int64
		Enable     bool
		Tag        string
		All_levels bool
	}

	PHPCode struct {
		Root string
		FPM  string
	} `onion:"php_code"`
}

func defaultLayer() onion.Layer {
	d := onion.NewDefaultLayer()
	assert.Nil(d.SetDefault("site", "gad.loc"))
	assert.Nil(d.SetDefault("mount_point", "/"))
	assert.Nil(d.SetDefault("retry", 1*time.Minute))
	assert.Nil(d.SetDefault("devel_mode", true))
	assert.Nil(d.SetDefault("cors", true))
	assert.Nil(d.SetDefault("max_cpu_available", runtime.NumCPU()))
	assert.Nil(d.SetDefault("proto", "http"))
	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 0)
	if err != nil {
		port = 80
	}

	assert.Nil(d.SetDefault("port", port))
	assert.Nil(d.SetDefault("ip", "127.0.0.1"))
	assert.Nil(d.SetDefault("time_zone", "Asia/Tehran"))
	assert.Nil(d.SetDefault("machine_name", "m1"))

	p, err := expand.Path("$HOME/gad/statics")
	assert.Nil(err)
	assert.Nil(d.SetDefault("static_root", p))
	var (
		rport = "6379"
		rhost = "127.0.0.1"
		rpass string
	)

	redisURL := os.Getenv("REDIS_URL")
	if all := redisPattern.FindStringSubmatch(redisURL); len(all) == 5 {
		rport = all[4]
		rhost = all[3]
		rpass = all[2]
	}

	assert.Nil(d.SetDefault("redis.size", 200))
	assert.Nil(d.SetDefault("redis.network", "tcp"))
	assert.Nil(d.SetDefault("redis.address", fmt.Sprintf("%s:%s", rhost, rport)))
	assert.Nil(d.SetDefault("redis.password", rpass))

	// TODO : move it to clickyab section
	assert.Nil(d.SetDefault("redis.days", 2))

	// TODO :  make sure ?parseTime=true is always set!
	assert.Nil(
		d.SetDefault(
			"mysql.rdsn",
			"dev:cH3M7Z7I4sY8QP&ll130U&73&6KS$o@tcp(db-1.clickyab.ae:3306)/clickyab?charset=utf8&parseTime=true&charset=utf8",
		),
	)
	assert.Nil(d.SetDefault("mysql.wdsn", "root:bita123@tcp(127.0.0.1:3306)/clickyab?charset=utf8&parseTime=true&charset=utf8"))
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
	//assert.Nil(d.SetDefault(
	//	"clickyab.ctr_const",
	//	[]string{
	//		transport.AD_SLOT,
	//		transport.AD_WEBSITE,
	//		transport.CAMPAIGN,
	//		transport.CAMPAIGN_SLOT,
	//		transport.SLOT,
	//	},
	//))
	assert.Nil(d.SetDefault("clickyab.min_imp", 1000))
	assert.Nil(d.SetDefault("clickyab.min_frequency", 2))
	assert.Nil(d.SetDefault("clickyab.daily_imp_expire", 7*24*time.Hour))
	assert.Nil(d.SetDefault("clickyab.daily_click_expire", 7*24*time.Hour))
	assert.Nil(d.SetDefault("clickyab.daily_cap_expire", 72*time.Hour))
	assert.Nil(d.SetDefault("clickyab.mega_imp_expire", 2*time.Hour))
	assert.Nil(d.SetDefault("clickyab.conv_delay", time.Second*10))
	assert.Nil(d.SetDefault("clickyab.conv_retry", 8))
	assert.Nil(d.SetDefault("clickyab.min_cpm_floor_web", 1000))
	assert.Nil(d.SetDefault("clickyab.min_cpm_floor_app", 100))
	assert.Nil(d.SetDefault("clickyab.max_load_fail", 3))
	assert.Nil(d.SetDefault("clickyab.cop_len", 10))
	assert.Nil(d.SetDefault("clickyab.fast_click", 2))
	assert.Nil(d.SetDefault("clickyab.ad_ctr_effect", 70))
	assert.Nil(d.SetDefault("clickyab.slot_ctr_effect", 30))
	assert.Nil(d.SetDefault("clickyab.native.max_count", 12))
	assert.Nil(d.SetDefault("clickyab.vast.default_duration", "00:00:05"))
	assert.Nil(d.SetDefault("clickyab.vast.default_skipoff", "00:00:03"))
	assert.Nil(d.SetDefault("clickyab.under_floor", true))
	assert.Nil(d.SetDefault("clickyab.web_min_bid", 2000))
	assert.Nil(d.SetDefault("clickyab.app_min_bid", 700))

	assert.Nil(d.SetDefault("clickyab.min_cpc_vast", 2000))
	assert.Nil(d.SetDefault("clickyab.min_cpc_app", 700))
	assert.Nil(d.SetDefault("clickyab.min_cpc_web", 2500))
	assert.Nil(d.SetDefault("clickyab.min_cpc_native", 1500))
	assert.Nil(d.SetDefault("clickyab.floor_div.native", 3))
	assert.Nil(d.SetDefault("clickyab.floor_div.app", 1))
	assert.Nil(d.SetDefault("clickyab.floor_div.vast", 3))

	assert.Nil(d.SetDefault("services.fluentd.host", "fluentd.monitoring"))
	assert.Nil(d.SetDefault("services.fluentd.port", 24224))
	assert.Nil(d.SetDefault("services.fluentd.enable", false))
	assert.Nil(d.SetDefault("services.fluentd.tag", "log.gad"))
	assert.Nil(d.SetDefault("services.fluentd.all_levels", false))

	p, err = expand.Path("$HOME/gad/clickyab-server/a")
	assert.Nil(err)
	assert.Nil(d.SetDefault("php_code.root", p))
	assert.Nil(d.SetDefault("php_code.fpm", "127.0.0.1:9999"))

	assert.Nil(d.SetDefault("slack.channel", "notifications"))
	assert.Nil(d.SetDefault("slack.username", "LilBro"))
	assert.Nil(d.SetDefault("slack.webhookurl", "https://hooks.slack.com/services/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"))
	assert.Nil(d.SetDefault("slack.active", false))

	return d
}
