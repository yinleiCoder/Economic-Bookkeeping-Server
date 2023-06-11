package forms

// form validate.
type LoginPasswordForm struct {
	Mobile    string `json:"mobile"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

type RegisterForm struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	SMSCode  string `json:"smscode"`
}
