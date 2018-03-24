package api

// An Equation is a representation of an equation from the database.
type Equation struct {
	ID          int64    `json:"id"`
	Source      string   `json:"source"`
	Description string   `json:"description"`
	Author      int64    `json:"author"`
	Categories  []string `json:"categories"`
	Score       int      `json:"score"`
	Confirmed   bool     `json:"confirmed"`
	Timestamp   int64    `json:"added"`
}
