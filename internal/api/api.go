package api

import (
	"context"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/arslanovdi/logistic-package-api/internal/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"

	pb "github.com/arslanovdi/logistic-package-api/pkg/logistic-package-api"
)

var (
	totalTemplateNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "omp_template_api_template_not_found_total",
		Help: "Total number of templates that were not found",
	})
)

// packageAPI имплементация grpc сервера
type packageAPI struct {
	pb.UnimplementedLogisticPackageApiServiceServer
	packageService *service.PackageService
}

func (p *packageAPI) CreatePackageV1(ctx context.Context, req *pb.CreatePackageRequestV1) (*pb.CreatePackageResponseV1, error) {

	log := slog.With("func", "api.CreatePackageV1")

	if err := req.Validate(); err != nil {
		log.Error("CreatePackageV1 - invalid argument", slog.Any("error", err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pkg := model.Package{}
	pkg.FromProto(req.Value)

	id, err := p.packageService.Create(ctx, pkg)
	if err != nil {
		log.Error("CreatePackageV1 - failed", slog.Any("error", err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug("CreatePackageV1 - created", slog.Uint64("id", *id))

	return &pb.CreatePackageResponseV1{PackageId: *id},
		status.New(codes.OK, "").Err()
}

func (p *packageAPI) DeletePackageV1(ctx context.Context, req *pb.DeletePackageV1Request) (*pb.DeletePackageV1Response, error) {

	log := slog.With("func", "api.DeletePackageV1")

	if err := req.Validate(); err != nil {
		log.Error("DeletePackageV1 - invalid argument", slog.Any("error", err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ok, err := p.packageService.DeletePackage(ctx, req.PackageId)
	if err != nil {
		log.Error("DeletePackageV1 - failed", slog.Any("error", err))
		return nil, status.Error(codes.Internal, err.Error()) // TODO переделать коды, с возвратом ошибок из сервисного слоя -> из репо
	}

	if !ok { // TODO возможно нужно избавиться от возврата ok в проекте, возвращать ошибки
		log.Debug("DeletePackageV1 - not found", slog.Uint64("id", req.PackageId))
		return nil, status.Error(codes.NotFound, "")
	}

	log.Debug("DeletePackageV1 - deleted", slog.Uint64("id", req.PackageId))
	return &pb.DeletePackageV1Response{
			Deleted: ok,
		},
		status.New(codes.OK, "").Err()

}

func (p *packageAPI) GetPackageV1(ctx context.Context, req *pb.GetPackageV1Request) (*pb.GetPackageV1Response, error) {

	log := slog.With("func", "api.GetPackageV1")

	if err := req.Validate(); err != nil {
		log.Error("GetPackageV1 - invalid argument", slog.Any("error", err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pkg, err := p.packageService.GetPackage(ctx, req.PackageId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if pkg == nil { // TODO возможно нужно возвращать ошибки NotFound
		log.Debug("GetPackageV1 - not found", slog.Uint64("id", req.PackageId))
		return nil, status.Error(codes.NotFound, "")
	}

	log.Debug("GetPackageV1 - found", slog.Uint64("id", req.PackageId))

	return &pb.GetPackageV1Response{
			Value: pkg.ToProto(),
		},
		status.New(codes.OK, "").Err()
}

func (p *packageAPI) ListPackagesV1(ctx context.Context, req *pb.ListPackagesV1Request) (*pb.ListPackagesV1Response, error) {

	log := slog.With("func", "api.ListPackagesV1")

	if err := req.Validate(); err != nil {
		log.Error("ListPackagesV1 - invalid argument", slog.Any("error", err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	packages, err := p.packageService.ListPackages(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(packages) == 0 {
		log.Debug("ListPackagesV1 - empty")
		return &pb.ListPackagesV1Response{}, status.Error(codes.NotFound, "")
	}

	log.Debug("ListPackagesV1 - found", slog.Uint64("count", uint64(len(packages))))

	resp := make([]*pb.Package, len(packages))
	for i := 0; i < len(packages); i++ {
		resp[i] = packages[i].ToProto()
	}

	return &pb.ListPackagesV1Response{
			Packages: resp,
		},
		status.New(codes.OK, "").Err()
}

func (p *packageAPI) UpdatePackageV1(ctx context.Context, req *pb.UpdatePackageV1Request) (*pb.UpdatePackageV1Response, error) {

	log := slog.With("func", "api.UpdatePackageV1")

	if err := req.Validate(); err != nil {
		log.Error("UpdatePackageV1 - invalid argument", slog.Any("error", err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pkg := model.Package{}
	pkg.FromProto(req.Value)

	ok, err := p.packageService.UpdatePackage(ctx, req.PackageId, pkg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !ok {
		log.Debug("UpdatePackageV1 - not found", slog.Uint64("id", req.PackageId))
		return nil, status.Error(codes.NotFound, "")
	}

	log.Debug("UpdatePackageV1 - updated", slog.Uint64("id", req.PackageId))

	return &pb.UpdatePackageV1Response{
			Updated: ok,
		},
		status.New(codes.OK, "").Err()
}

// NewPackageAPI returns api of logistic-package-api service
func NewPackageAPI(packageService *service.PackageService) *packageAPI {
	return &packageAPI{
		packageService: packageService,
	}
}
