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
	ID      uint64    `db:"id" json:"ID,omitempty"`
	Title   string    `db:"title" json:"title,omitempty"`
	Weight  uint64    `db:"weight" json:"weight,omitempty"`
	Created time.Time `db:"created" json:"created"`
	Updated time.Time `db:"updated" json:"updated"`
	Removed bool      `db:"removed" json:"removed"`
}

type EventType uint8

type EventStatus uint8

const (
	_ EventType = iota
	Created
	Updated
	Removed
)

const (
	_ EventStatus = iota
	Locked
	Unlocked
)

type PackageEvent struct {
	ID        uint64      `db:"id"`
	PackageID uint64      `db:"package_id"`
	Type      EventType   `db:"type"`
	Status    EventStatus `db:"status"`
	Payload   []byte      `db:"payload"`
	Updated   time.Time   `db:"updated"`
}

// String implements fmt.Stringer
func (c *Package) String() string {
	return fmt.Sprintf("ID: %d, Title: %s, Weight: %d, Created: %s, UpdatedAt: %s, Removed: %t", c.ID, c.Title, c.Weight, c.Created, c.Updated, c.Removed)
}

// LogValue implements slog.LogValuer interface
func (c *Package) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Uint64("ID", c.ID),
		slog.String("Title", c.Title),
		slog.Uint64("Weight", c.Weight),
		slog.Time("Created", c.Created),
		slog.Time("Updated", c.Updated),
		slog.Bool("Removed", c.Removed),
	)
}

// ToProto converts model.Package to pb.Package
func (c *Package) ToProto() *pb.Package {
	return &pb.Package{
		Id:     c.ID,
		Title:  c.Title,
		Weight: &c.Weight,
		Created: &timestamp.Timestamp{
			Seconds: c.Created.Unix(),
			Nanos:   int32(c.Created.Nanosecond()),
		},
		Updated: &timestamp.Timestamp{
			Seconds: c.Updated.Unix(),
			Nanos:   int32(c.Updated.Nanosecond()),
		},
	}
}

// FromProto converts pb.Package to model.Package
func (c *Package) FromProto(pkg *pb.Package) {
	c.ID = pkg.Id
	c.Title = pkg.Title
	c.Weight = *pkg.Weight
	c.Created = time.Unix(pkg.Created.Seconds, int64(pkg.Created.Nanos))
	c.Updated = time.Unix(pkg.Updated.Seconds, int64(pkg.Updated.Nanos))
}
