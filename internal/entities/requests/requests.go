package requests

import (
	"time"

	m "github.com/risetyll/finuslugi/internal/entities/material"
)

type CreateMaterialRequest struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type GetMaterialByIdRequest struct {
	UUID int `json:"uuid"`
}

type UpdateMaterialRequest struct {
	UUID    int    `json:"uuid"`
	Status  string `json:"status"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type GetMaterialsRequest struct {
	Type        string    `json:"type"`
	Page        int       `json:"page"`
	PageSize    int       `json:"page_size"`
	CreatedFrom time.Time `json:"created_from"`
	CreatedTo   time.Time `json:"created_to"`
}

type GetMaterialsResponse struct {
	Materials []*m.Material `json:"materials"`
	Total     int           `json:"total"`
	Page      int           `json:"page"`
	PageSize  int           `json:"page_size"`
}
