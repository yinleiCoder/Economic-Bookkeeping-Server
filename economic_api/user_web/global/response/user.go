package response

type UserResponse struct {
	Id       int32  `json:"id"`
	Mobile   string `json:"mobile"`
	NickName string `json:"name"`
	Birthday string `json:"birthday"`
	Gender   string `json:"gender"`
}
