package utils

import (
	"config"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"version"

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/go-redmine"
)

// RedmineDoError try to create an issue in redmine tracker
func RedmineDoError(err interface{}, stack []byte) {
	url := config.Config.Redmine.URL
	key := config.Config.Redmine.APIKey

	c := redmine.NewClient(url, key)
	v := version.GetVersion()
	short := v.Short
	commits := v.Count

	title := extractTitle(short, commits, err)

	// redmine can not accept more than 255 character title
	if len(title.Error()) > 200 {
		str := title.Error()
		title = errors.New(str[:200] + "...")
	}

	var filters []redmine.IssueFilter
	filters = append(filters, redmine.IssueFilter{Key: "limit", Value: "1"})
	filters = append(filters, redmine.IssueFilter{Key: "subject", Value: title.Error()})
	//filters = append(filters, redmine.IssueFilter{Key: "status_id", Value: "open"})

	issues, err := c.FilterIssues(filters...)
	if err != nil {
		logrus.Warn(err)
		return
	}
	var is *redmine.Issue
	if len(issues) > 0 {
		for i := range issues {
			if issues[i].Status.Id == config.Config.Redmine.NewIssueTypeID {
				is = &issues[i]
				break
			}
		}
	}

	if is != nil {
		is.Notes = string(stack)
		err := c.UpdateIssue(*is)
		if err != nil {
			logrus.Warn(err)
		}
	} else {
		is = &redmine.Issue{}
		is.Subject = title.Error()
		is.Description = string(stack)
		is.ProjectId = config.Config.Redmine.ProjectID

		_, err := c.CreateIssue(*is)
		if err != nil {
			logrus.Warn(err)
		}
	}
}

// SlackPayload the slack payload
type SlackPayload struct {
	Channel     string            `json:"channel"`
	Text        string            `json:"text"`
	Username    string            `json:"username"`
	IconURL     string            `json:"icon_url,omitempty"`
	IconEmoji   string            `json:"icon_emoji,omitempty"`
	Parse       string            `json:"parse"`
	Attachments []SlackAttachment `json:"attachments"`
}

// SlackAttachment the attachment
type SlackAttachment struct {
	Color   string `json:"color"`
	Text    string `json:"text"`
	PreText string `json:"pretext,omitempty"`
	Title   string `json:"title,omitempty"`
}

func extractTitle(short string, commits int64, err interface{}) error {
	title := fmt.Errorf("[%s, %d] cannot extract title, the type is %T, value is %v", short, commits, err, err)
	switch err.(type) {
	case string:
		title = fmt.Errorf("[%s, %d] %s", short, commits, err.(string))
	case error:
		title = fmt.Errorf("[%s, %d] %s", short, commits, err.(error).Error())
	case *logrus.Entry:
		title = fmt.Errorf("[%s, %d] %s", short, commits, err.(*logrus.Entry).Message)
	}

	return title
}

// SlackDoMessage Try to send message to configured slack channel
func SlackDoMessage(err interface{}, icon string, attachments ...SlackAttachment) {
	payload := &SlackPayload{}
	payload.Channel = config.Config.Slack.Channel

	v := version.GetVersion()
	short := v.Short
	commits := v.Count

	title := extractTitle(short, commits, err)
	payload.Text = title.Error()
	payload.Username = config.Config.Slack.Username
	payload.Parse = "full" // WTF?
	if icon != "" {
		if icon[0] == ':' {
			payload.IconEmoji = icon
		} else {
			payload.IconURL = icon
		}
	}

	payload.Attachments = attachments

	encoded, err := json.Marshal(payload)
	if err != nil {
		logrus.WithField("payload", payload).Warn(err)
		return
	}

	resp, err := http.PostForm(config.Config.Slack.WebHookURL, url.Values{"payload": {string(encoded)}})
	if err != nil {
		logrus.WithField("payload", payload).Warn(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		logrus.WithField("response", resp).Warn("sending payload to slack failed")
		return
	}
}
