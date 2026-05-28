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
		Addr string `json:"addr,omitempty" doc:"client address" example:"1.1.1.1"`
		Auth string `json:"auth" doc:"auth data" example:"hardpassword1234"`
		Tx   int    `json:"tx,omitempty"`
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
