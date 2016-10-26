package config

import (
	"runtime"

	"time"

	"github.com/Sirupsen/logrus"
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
		Size                 int
		Network              string
		Address              string
		Password             string //Daily Statistic TimeOut Expiration
		DailyImpExpireTime   time.Duration
		DailyClickExpireTime time.Duration
		DailyCapExpireTime   time.Duration
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

	DefaultCTR float64
}

func init() {
	var err error

	Config.Site = "gad.loc"
	Config.MountPoint = "/"
	Config.DevelMode = true
	Config.CORS = true
	Config.MaxCPUAvailable = runtime.NumCPU()
	Config.Proto = "http"
	Config.Port = ":80"
	Config.StaticRoot, err = expand.Path("/statics")
	if err != nil {
		logrus.Panic(err)
	}

	Config.Redis.Size = 10
	Config.Redis.Network = "tcp"
	Config.Redis.Address = ":6379"
	Config.Redis.DailyClickExpireTime = 72 * time.Hour
	//Config.Redis.Password = ""
	Config.Redis.DailyClickExpireTime = 72 * time.Hour
	Config.Redis.DailyCapExpireTime = 24 * time.Hour

	// TODO : make sure ?parseTime=true is always set!
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	//Config.Mysql.DSN = "novid:x4WT4a2o86oR1lup@tcp(5.9.150.114:3312)/clickyab?parseTime=true"

	//db, err := sql.Open("mysql", "clickyab_test:760f5bad06b18134ef6@tcp(46.4.116.104:3306)/clickyab?charset=utf8")
	Config.Mysql.DSN = "clickyab_test:760f5bad06b18134ef6@tcp(46.4.116.104:3306)/clickyab?charset=utf8&parseTime=true"

	Config.Mysql.MaxConnection = 100
	Config.Mysql.MaxIdleConnection = 10
	Config.Page.PerPage = 10
	Config.Page.MaxPerPage = 100
	Config.Page.MinPerPage = 1

	Config.Select.Date = 0
	Config.Select.Hour = 1
	Config.Select.Balance = 50000

	Config.TimeZone = "Asia/Tehran"

	Config.AMQP.Publisher = 30
	Config.AMQP.Exchange = "cy"
	Config.AMQP.DSN = "amqp://server:bita123@127.0.0.1:5672/"
	Config.AMQP.ConfirmLen = 100

	Config.Redmine.APIKey = "5d29e2039762e19fbfe3db72b013bf356b3ed072"
	Config.Redmine.URL = "https://redmine.azmoona.com/"
	Config.Redmine.ProjectID = 1
	Config.Redmine.Active = false
	Config.Redmine.NewIssueTypeID = 4

	Config.Slack.Channel = "#app"
	Config.Slack.Username = "azmoona"
	Config.Slack.WebHookURL = "https://hooks.slack.com/services/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"
	Config.Slack.Active = false

	Config.DefaultCTR = 0.1

}
