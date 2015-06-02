package todo

//Todo model
type Todo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}
