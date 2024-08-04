package repository

import (
	"github.com/risetyll/finuslugi/internal/entities/material"
	"github.com/risetyll/finuslugi/internal/entities/requests"
)

type Repository interface {
	CreateMaterial(req *requests.CreateMaterialRequest) error
	GetMaterialById(req *requests.GetMaterialByIdRequest) (*material.Material, error)
	UpdateMaterial(req *requests.UpdateMaterialRequest) error
	GetMaterials(req *requests.GetMaterialsRequest) ([]*material.Material, error)
}
