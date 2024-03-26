package model

import (
	"fmt"
	pb "github.com/arslanovdi/logistic-package-api/pkg/logistic-package-api"
	"github.com/golang/protobuf/ptypes/timestamp"
	"log/slog"
	"time"
)

// Package сущность пакета
type Package struct {
	ID        uint64    `db:"id"`
	Title     string    `db:"title"`
	Weight    uint64    `db:"weight"`
	CreatedAt time.Time `db:"createdAt"`
}

type EventType uint8

type EventStatus uint8

const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

type PackageEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Package
}

// String implements fmt.Stringer
func (c *Package) String() string {
	return fmt.Sprintf("ID: %d, Title: %s, Weight: %d, CreatedAt: %s", c.ID, c.Title, c.Weight, c.CreatedAt)
}

// LogValue implements slog.LogValuer interface
func (c *Package) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Uint64("ID", c.ID),
		slog.String("Title", c.Title),
		slog.Uint64("Weight", c.Weight),
		slog.Time("CreatedAt", c.CreatedAt),
	)
}

// ToProto converts model.Package to pb.Package
func (c *Package) ToProto() *pb.Package {
	return &pb.Package{
		Id:     c.ID,
		Title:  c.Title,
		Weight: &c.Weight,
		Created: &timestamp.Timestamp{
			Seconds: c.CreatedAt.Unix(),
			Nanos:   int32(c.CreatedAt.Nanosecond()),
		},
	}
}

// FromProto converts pb.Package to model.Package
func (c *Package) FromProto(pkg *pb.Package) {
	c.ID = pkg.Id
	c.Title = pkg.Title
	c.Weight = *pkg.Weight
	c.CreatedAt = time.Unix(pkg.Created.Seconds, int64(pkg.Created.Nanos))
}
