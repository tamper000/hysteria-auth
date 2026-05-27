package handler

type User struct {
	Body struct {
		ID       string `json:"id" doc:"user ID" example:"123456" example:"johnDoe"`
		Auth     string `json:"auth" doc:"auth data" example:"hardpassword1234"`
		Optional string `json:"optional" doc:"optional data" example:"telegramID:12345"` // Additional field (for example, for Telegram user ID)
	}
}

type UserID struct {
	Body struct {
		ID string `json:"id" doc:"user ID" example:"123456" example:"johnDoe"`
	}
}

type UserAuth struct {
	Body struct {
		Auth string `json:"auth" doc:"auth data" example:"hardpassword1234"`
	}
}

type ResultOutput struct {
	Status int
}

type AuthOutput struct {
	Body struct {
		Ok bool   `json:"ok" doc:"status" example:"true"`
		ID string `json:"id" doc:"user ID" example:"123456" example:"johnDoe"`
	}
}
