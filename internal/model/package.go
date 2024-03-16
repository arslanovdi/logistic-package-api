package model

import (
	"errors"
	"time"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

// Package сущность пакета
type Package struct {
	ID        uint64
	Title     string
	CreatedAt time.Time
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
