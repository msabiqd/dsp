package entity

type GetUserByIdRequest struct {
	Id int
}

type GetUserByPhoneNumberRequest struct {
	PhoneNumber string
}

type CreateUserRequest struct {
	FullName    string `validate:"required,min=3,max=60"`
	PhoneNumber string `validate:"required,min=10,max=13,phonenumber"`
	Password    string `validate:"required,min=6,max=64,password"`
}

type UpdateUserRequest struct {
	Id          int
	FullName    string
	PhoneNumber string
}

type SuccessLoginIncrementRequest struct {
	Id              int
	SuccessfulLogin int
}

type UserResponse struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}

type User struct {
	Id              int
	FullName        string
	PhoneNumber     string
	Password        string
	SuccessfulLogin int
}
