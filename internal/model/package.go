package model

import (
	"fmt"
	pb "github.com/arslanovdi/logistic-package-api/pkg/logistic-package-api"
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

// Package сущность пакета
type Package struct {
	ID        uint64
	Title     string
	Weight    uint64
	CreatedAt time.Time
}

/*type Package struct {
	ID        uint64
	Title     string
	CreatedAt time.Time
}*/

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

func (c *Package) String() string {
	return fmt.Sprintf("ID: %d, Title: %s, Weight: %d, CreatedAt: %s", c.ID, c.Title, c.Weight, c.CreatedAt)
}

func (c *Package) ToProto() *pb.Package {
	return &pb.Package{
		Title:  c.Title,
		Weight: &c.Weight,
		Created: &timestamp.Timestamp{
			Seconds: c.CreatedAt.Unix(),
			Nanos:   int32(c.CreatedAt.Nanosecond()),
		},
	}
}

func (c *Package) FromProto(pkg *pb.Package) {
	c.ID = pkg.Id
	c.Title = pkg.Title
	c.Weight = *pkg.Weight
	c.CreatedAt = time.Unix(pkg.Created.Seconds, int64(pkg.Created.Nanos))
}
