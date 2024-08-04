package usecase

import (
	"github.com/risetyll/finuslugi/internal/entities/material"
	"github.com/risetyll/finuslugi/internal/entities/requests"
	"github.com/risetyll/finuslugi/internal/storage/repository"
)

type DatabaseUsecase struct {
	repo repository.Repository
}

func New(repo repository.Repository) *DatabaseUsecase {
	return &DatabaseUsecase{repo}
}

func (uc *DatabaseUsecase) CreateMaterial(req *requests.CreateMaterialRequest) error {
	return uc.repo.CreateMaterial(req)
}

func (uc *DatabaseUsecase) GetMaterialById(req *requests.GetMaterialByIdRequest) (*material.Material, error) {
	return uc.repo.GetMaterialById(req)
}

func (uc *DatabaseUsecase) UpdateMaterial(req *requests.UpdateMaterialRequest) error {
	return uc.repo.UpdateMaterial(req)
}

func (uc *DatabaseUsecase) GetMaterials(req *requests.GetMaterialsRequest) ([]*material.Material, error) {
	return uc.repo.GetMaterials(req)
}
