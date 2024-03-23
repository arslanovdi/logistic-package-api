package api

import (
	"context"
	"github.com/arslanovdi/logistic-package-api/internal/app/repo"
	"github.com/arslanovdi/logistic-package-api/internal/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log/slog"

	pb "github.com/arslanovdi/logistic-package-api/pkg/logistic-package-api"
)

var (
	totalTemplateNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "omp_template_api_template_not_found_total",
		Help: "Total number of templates that were not found",
	})
)

type packageAPI struct {
	pb.UnimplementedLogisticPackageApiServiceServer
	repo repo.Repo
}

func (p *packageAPI) CreatePackageV1(ctx context.Context, req *pb.CreatePackageRequestV1) (*pb.CreatePackageResponseV1, error) {
	log := slog.With("func", "api.CreatePackageV1")
	//TODO implement me
	log.Debug("CreatePackageV1 - not implemented")
	return nil, model.ErrNotImplemented
}

func (p *packageAPI) DeletePackageV1(ctx context.Context, req *pb.DeletePackageV1Request) (*pb.DeletePackageV1Response, error) {
	log := slog.With("func", "api.DeletePackageV1")
	//TODO implement me
	log.Debug("DeletePackageV1 - not implemented")
	return nil, model.ErrNotImplemented
}

func (p *packageAPI) GetPackageV1(ctx context.Context, req *pb.GetPackageV1Request) (*pb.GetPackageV1Response, error) {
	log := slog.With("func", "api.GetPackageV1")
	//TODO implement me
	log.Debug("GetPackageV1 - not implemented")
	return nil, model.ErrNotImplemented

	/*if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeTemplateV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	template, err := o.repo.Describe(ctx, req.PackageId)
	if err != nil {
		log.Error().Err(err).Msg("DescribeTemplateV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if template == nil {
		log.Debug().Uint64("templateId", req.PackageId).Msg("template not found")
		totalTemplateNotFound.Inc()

		return nil, status.Error(codes.NotFound, "template not found")
	}

	log.Debug().Msg("DescribeTemplateV1 - success")

	return &pb.DescribePackageV1Response{}, nil*/
}

func (p *packageAPI) ListPackagesV1(ctx context.Context, req *pb.ListPackagesV1Request) (*pb.ListPackagesV1Response, error) {
	log := slog.With("func", "api.ListPackagesV1")
	//TODO implement me
	log.Debug("ListPackagesV1 - not implemented")
	return nil, model.ErrNotImplemented
}

func (p *packageAPI) UpdatePackageV1(ctx context.Context, req *pb.UpdatePackageV1Request) (*pb.UpdatePackageV1Response, error) {
	log := slog.With("func", "api.UpdatePackageV1")
	//TODO implement me
	log.Debug("UpdatePackageV1 - not implemented")
	return nil, model.ErrNotImplemented
}

// NewPackageAPI returns api of logistic-package-api service
func NewPackageAPI(r repo.Repo) *packageAPI {
	return &packageAPI{repo: r}
}
