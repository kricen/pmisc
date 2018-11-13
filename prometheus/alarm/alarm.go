package alarm

// AlarmManagement : alarm management model,contains the alarm channel and  person in charge of
type AlarmManagement struct {
	JobName     string
	Email       EmailDTO
	Wechat      []string
	WechatGroup []string
	Handlers    []string
}

// EmailDTO : email data transfer object,contains necessary items
type EmailDTO struct {
	Recipients []string
	CCList     []string
}

func NewAlarmManagement() {

}
