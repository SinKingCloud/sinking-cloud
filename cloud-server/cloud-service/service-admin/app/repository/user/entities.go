package user

// SelectUser 用户查询
type SelectUser struct {
	Account         string
	Phone           string
	Email           string
	NickName        string
	Status          string
	CreateTimeStart string
	CreateTimeEnd   string
	UpdateTimeStart string
	UpdateTimeEnd   string
}

// UpdateUser 用户更新
type UpdateUser struct {
	Account  interface{}
	Phone    interface{}
	Email    interface{}
	Password interface{}
	Avatar   interface{}
	NickName interface{}
	Status   interface{}
}
