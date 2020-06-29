package vo

type CreateCategory struct {
	Name string `json:"name" binding:"required"`
}
