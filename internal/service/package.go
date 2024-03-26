package service

import (
	"context"
	"fmt"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repo interface for work with database
type Repo interface {
	Create(ctx context.Context, pkg model.Package) (*uint64, error)
	DeletePackage(ctx context.Context, id uint64) (bool, error)
	GetPackage(ctx context.Context, id uint64) (*model.Package, error)
	ListPackages(ctx context.Context, offset uint64, limit uint64) ([]model.Package, error)
	UpdatePackage(ctx context.Context, pkg model.Package) (bool, error)
}

// PackageService is service for Package
type PackageService struct {
	dbpool *pgxpool.Pool // для транзакций
	repo   Repo
}

func (p *PackageService) Create(ctx context.Context, pkg model.Package) (*uint64, error) {

	id, err := p.repo.Create(ctx, pkg)
	if err != nil {
		return nil, fmt.Errorf("service.PackageService.Create: %w", err)
	}

	return id, nil
}

func (p *PackageService) DeletePackage(ctx context.Context, id uint64) (bool, error) {

	ok, err := p.repo.DeletePackage(ctx, id)
	if err != nil {
		return false, fmt.Errorf("service.PackageService.DeletePackage: %w", err)
	}

	return ok, nil
}

func (p *PackageService) GetPackage(ctx context.Context, id uint64) (*model.Package, error) {

	pkg, err := p.repo.GetPackage(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service.PackageService.GetPackage: %w", err)
	}
	return pkg, nil
}

func (p *PackageService) ListPackages(ctx context.Context, offset uint64, limit uint64) ([]model.Package, error) {

	packages, err := p.repo.ListPackages(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("service.PackageService.ListPackages: %w", err)
	}
	return packages, nil
}

func (p *PackageService) UpdatePackage(ctx context.Context, pkg model.Package) (bool, error) {

	ok, err := p.repo.UpdatePackage(ctx, pkg)
	if err != nil {
		return false, fmt.Errorf("service.PackageService.UpdatePackage: %w", err)
	}
	return ok, nil
}

func NewPackageService(dbpool *pgxpool.Pool, repo Repo) *PackageService {
	return &PackageService{
		dbpool: dbpool,
		repo:   repo,
	}
}
