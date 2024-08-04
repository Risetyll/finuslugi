package material

import (
	"time"
)

// Material представляет собой материал школы
type Material struct {
	UUID      int       `json:"uuid"`
	Type      string    `json:"type"`   // статья, видеоролик, презентация
	Status    string    `json:"status"` // архивный, активный
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
