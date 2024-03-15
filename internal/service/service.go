package service

import "github.com/arslanovdi/logistic-package-api/internal/model"

type PackageService interface {
	Describe(PackageID uint64) (*model.Package, error) // описание
	List(cursor uint64, limit uint64) ([]model.Package, error)
	Get(cursor uint64) (model.Package, error)
	Create(model.Package) (uint64, error)
	Update(packageID uint64, pkg model.Package) error
	Remove(packageID uint64) (bool, error)
}
