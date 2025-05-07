package dto

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY2OTE0MjQsImlhdCI6MTc0NjYwNTAyNCwic3ViIjoiMDFKVEU1VE5URVZZREVQUzhHRUUwWUc0Qk0ifQ.yDYZBaEAgeAshF7zLSGNxVL8Q5P70YheQCXO68pEzHc"`
}
