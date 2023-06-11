package forms

type SMSForm struct {
	Mobile string `json:"mobile"`
	Type   uint   `json:"type"` // 1 register or 2 login
}
