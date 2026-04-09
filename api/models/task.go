package models

type Task struct {
	ID          int    `json:"id" gorm:"primary_key"`
	Assignee    string `json:"assignee"`
	Description string `json:"description"`
	ProductID   int    `json:"product_id"`
	ProductCode string `json:"product_code"`
	ProductName string `json:"product_name"`
	Link        string `json:"link"`
	Creator     string `json:"creator"`
	Approver    string `json:"approver"`
	Comments    string `json:"comments"`
	Status      string `json:"status"`
}

type Comment struct {
	ID     int    `json:"id"`
	TaskID int    `json:"task_id"`
	Text   string `json:"text"`
}
