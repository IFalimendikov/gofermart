package models

type User struct {
	ID       string `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	ID         string `json:"user_id,omitempty"`
	Order      string `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual,omitempty"`
	UploadedAt string `json:"uploaded_at"`
}

type Balance struct {
	ID        string `json:"user_id"`
	Current   int    `json:"current"`
	Withdrawn int    `json:"withdrawn"`
}

type Withdrawal struct {
	ID          string `json:"user_id"`
	Order       string `json:"order"`
	Sum         int    `json:"sum"`
	ProcessedAt string `json:"processed_at"`
}
