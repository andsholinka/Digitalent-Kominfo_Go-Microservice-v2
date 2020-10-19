package entity

type AuthResponse struct {
	Code         int    `json:"code"`
	Status       string `json:"status,omitempty"`
	ErrorDetails string `json:"error_details,omitempty"`
	ErrorType    string `json:"error_type,omitempty"`
	Data         Data
}

type Data struct {
	Username string `json:"username"`
	Role     int    `json:"role"`
}
