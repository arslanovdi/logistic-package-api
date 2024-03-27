package api

import (
	"context"
	"errors"
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

func (p *packageAPI) CreateV1(ctx context.Context, req *pb.CreateRequestV1) (*pb.CreateResponseV1, error) {

	log := slog.With("func", "api.CreateV1")

	if err1 := req.Validate(); err1 != nil {
		log.Error("CreateV1 - invalid argument", slog.String("error", err1.Error()))
		return nil, status.Error(codes.InvalidArgument, err1.Error())
	}

	pkg := model.Package{}
	pkg.FromProto(req.Value)

	id, err2 := p.packageService.Create(ctx, pkg)
	if err2 != nil {
		log.Error("CreateV1 - failed", slog.String("error", err2.Error()))
		return nil, status.Error(codes.Internal, err2.Error())
	}

	log.Debug("CreateV1 - created", slog.Uint64("id", *id))

	return &pb.CreateResponseV1{PackageId: *id},
		status.New(codes.OK, "").Err()
}

func (p *packageAPI) DeleteV1(ctx context.Context, req *pb.DeleteV1Request) (*pb.DeleteV1Response, error) {

	log := slog.With("func", "api.DeleteV1")

	if err1 := req.Validate(); err1 != nil {
		log.Error("DeleteV1 - invalid argument", slog.String("error", err1.Error()))
		return nil, status.Error(codes.InvalidArgument, err1.Error())
	}

	err2 := p.packageService.Delete(ctx, req.PackageId)
	if err2 != nil {
		log.Error("DeleteV1 - failed", slog.String("error", err2.Error()))

		if errors.Is(err2, model.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err2.Error())
		}
		return nil, status.Error(codes.Internal, err2.Error())
	}

	log.Debug("DeleteV1 - deleted", slog.Uint64("id", req.PackageId))
	return &pb.DeleteV1Response{},
		status.New(codes.OK, "").Err()

}

func (p *packageAPI) GetV1(ctx context.Context, req *pb.GetV1Request) (*pb.GetV1Response, error) {

	log := slog.With("func", "api.GetV1")

	if err1 := req.Validate(); err1 != nil {
		log.Error("GetV1 - invalid argument", slog.String("error", err1.Error()))
		return nil, status.Error(codes.InvalidArgument, err1.Error())
	}

	pkg, err2 := p.packageService.Get(ctx, req.PackageId)
	if err2 != nil {
		log.Error("GetV1 - failed", slog.String("error", err2.Error()))
		return nil, status.Error(codes.Internal, err2.Error())
	}

	if pkg == nil { // TODO возможно нужно возвращать ошибки NotFound
		log.Debug("GetV1 - not found", slog.Uint64("id", req.PackageId))
		return nil, status.Error(codes.NotFound, "")
	}

	log.Debug("package found", slog.Any("package", pkg)) // TODO проверить

	return &pb.GetV1Response{
			Value: pkg.ToProto(),
		},
		status.New(codes.OK, "").Err()
}

func (p *packageAPI) ListV1(ctx context.Context, req *pb.ListV1Request) (*pb.ListV1Response, error) {

	log := slog.With("func", "api.ListV1")

	if err1 := req.Validate(); err1 != nil {
		log.Error("ListV1 - invalid argument", slog.String("error", err1.Error()))
		return nil, status.Error(codes.InvalidArgument, err1.Error())
	}

	packages, err2 := p.packageService.List(ctx, req.Offset, req.Limit)
	if err2 != nil {
		log.Error("ListV1 - failed", slog.String("error", err2.Error()))
		return nil, status.Error(codes.Internal, err2.Error())
	}

	if len(packages) == 0 {
		log.Debug("ListV1 - empty")
		return &pb.ListV1Response{}, status.Error(codes.NotFound, "")
	}

	log.Debug("ListV1 - found", slog.Uint64("count", uint64(len(packages))))

	resp := make([]*pb.Package, len(packages))
	for i := 0; i < len(packages); i++ {
		resp[i] = packages[i].ToProto()
	}

	return &pb.ListV1Response{
			Packages: resp,
		},
		status.New(codes.OK, "").Err()
}

func (p *packageAPI) UpdateV1(ctx context.Context, req *pb.UpdateV1Request) (*pb.UpdateV1Response, error) {

	log := slog.With("func", "api.UpdateV1")

	if err1 := req.Validate(); err1 != nil {
		log.Error("UpdateV1 - invalid argument", slog.String("error", err1.Error()))
		return nil, status.Error(codes.InvalidArgument, err1.Error())
	}

	pkg := model.Package{}
	pkg.FromProto(req.Value)

	ok, err2 := p.packageService.Update(ctx, pkg)
	if err2 != nil {
		log.Error("UpdateV1 - failed", slog.String("error", err2.Error()))
		return nil, status.Error(codes.Internal, err2.Error())
	}
	if !ok {
		log.Debug("UpdateV1 - not found", slog.Uint64("id", pkg.ID))
		return nil, status.Error(codes.NotFound, "")
	}

	log.Debug("UpdateV1 - updated", slog.Uint64("id", pkg.ID))

	return &pb.UpdateV1Response{},
		status.New(codes.OK, "").Err()
}

// NewPackageAPI returns api of logistic-package-api service
func NewPackageAPI(packageService *service.PackageService) *packageAPI {
	return &packageAPI{
		packageService: packageService,
	}
}
