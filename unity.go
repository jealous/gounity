/*
Package gounity provides the Golang API for Unity Storage Systems.

The library communicates with Unity with REST API.
*/
package gounity

import (
	"github.com/jealous/gounity/lib"
)

type Unity interface {
	// Ip returns the IP address of the connected system.
	Ip() string

	// Username returns the current user used for REST connection.
	Username() string
}

func New(ip, username, password string) (Unity, error) {
	unity := *lib.NewSystem(ip, username, password)
	return &unity, nil
}
