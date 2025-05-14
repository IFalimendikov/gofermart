package models

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	ID         string  `json:"login,omitempty"`
	Order      string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

type Balance struct {
	ID        string `json:"login"`
	Current   float64    `json:"current"`
	Withdrawn float64    `json:"withdrawn"`
}

type Withdrawal struct {
	ID          string `json:"login"`
	Order       string `json:"order"`
	Sum         float64    `json:"sum"`
	ProcessedAt string `json:"processed_at"`
}
