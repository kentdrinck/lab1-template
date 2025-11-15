package api

type CreatePersonDTO struct {
	Name    string `json:"name" binding:"required"`
	Age     *int   `json:"age,omitempty"`
	Address string `json:"address,omitempty"`
	Work    string `json:"work,omitempty"`
}

type UpdatePersonDTO struct {
	CreatePersonDTO
}