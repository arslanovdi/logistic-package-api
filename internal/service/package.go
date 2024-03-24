package service

import (
	"context"
	"github.com/arslanovdi/logistic-package-api/internal/database"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/jmoiron/sqlx"
)

// PackageService is service for Package
type PackageService struct {
	db   *sqlx.DB // для транзакций
	repo database.Repo
}

func (p *PackageService) Create(ctx context.Context, pkg model.Package) (*uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PackageService) DeletePackage(ctx context.Context, id uint64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PackageService) GetPackage(ctx context.Context, id uint64) (*model.Package, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PackageService) ListPackages(ctx context.Context, offset uint64, limit uint64) ([]model.Package, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PackageService) UpdatePackage(ctx context.Context, cursor uint64, pkg model.Package) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewPackageService(db *sqlx.DB, repo database.Repo) *PackageService {
	return &PackageService{
		db:   db,
		repo: repo,
	}
}
