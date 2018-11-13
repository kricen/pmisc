package alarm

import (
	"errors"
	"fmt"

	"github.com/pmisc/lib"
)

// AlarmManagement : alarm management model,contains the alarm channel and  person in charge of
type AlarmManagement struct {
	JobName     string   `json:"job_name"`
	Email       EmailDTO `json:"email"`
	Wechat      []string `json:"wechat"`
	WechatGroup []string `json:"wechat_group"`
	Handlers    []string `json:"handlers"`
}

// EmailDTO : email data transfer object,contains necessary items
type EmailDTO struct {
	Recipients []string `json:"recipients"`
	CCList     []string `json:"cc_list"`
}

// NewAlarmManagement : a func to return a AlarmManagement model
func NewAlarmManagement() (am *AlarmManagement) {

	return
}

// vaildate the email address before set email info
func (am *AlarmManagement) SetEmail(ed EmailDTO) (*AlarmManagement, error) {
	for _, addr := range ed.Recipients {
		if !lib.ValidateEmailAddress(addr) {
			return am, errors.New(fmt.Sprintf("%s is not a illedge email address", addr))
		}
	}
	for _, addr := range ed.CCList {
		if !lib.ValidateEmailAddress(addr) {
			return am, errors.New(fmt.Sprintf("%s is not a illedge email address", addr))
		}
	}
	am.Email = ed
	return am, nil
}

// Notification users of wechat
func (am *AlarmManagement) SetWechat(param []string) *AlarmManagement {
	am.Wechat = param
	return am
}

// Notification groups of wechat
func (am *AlarmManagement) SetWechatGroup(param []string) *AlarmManagement {
	am.WechatGroup = param
	return am
}

// person list  who in charge of this project
func (am *AlarmManagement) SetHandlers(param []string) *AlarmManagement {
	am.Handlers = param
	return am
}

// a distinct field for record project name
func (am *AlarmManagement) SetJobName() *AlarmManagement {

	return am
}
