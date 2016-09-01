/*
Package gounity provides the Golang API for Unity Storage Systems.

The library communicates with Unity with REST API.
*/
package gounity

import (
	"github.com/jealous/gounity/rsc"
)

type Unity interface {
	// Ip returns the IP address of the connected system.
	Ip() string

	// Username returns the current user used for REST connection.
	Username() string

	// Serial returns the unique ID for the unity instance
	Serial() string

	// Authenticate checks the credentials and connectivity to the system.
	Authenticate() error

	// GetPoolList returns all pools on the system
	GetPoolList() *rsc.PoolList

	// GetPoolByName retrieves the storage pool instance by name.
	// Return nil if not found.
	GetPoolByName(name string) *rsc.Pool

	// GetPoolById retrieves the storage pool instance by id.
	// Returns nil if not found.
	GetPoolById(id string) *rsc.Pool

	// GetLunList returns all LUNs on the system
	GetLunList() *rsc.LunList

	// GetLunByName retrieves the LUN instance by name.
	// Returns nil if not found.
	GetLunByName(name string) *rsc.Lun

	// GetLunById retrieves the LUN instance by id.
	// Returns nil if not found.
	GetLunById(id string) *rsc.Lun

	// GetHostList retrieves all hosts available on the system.
	GetHostList() *rsc.HostList

	// GetHostById retrieves the host by id
	// Returns nil if not found.
	GetHostById(id string) *rsc.Host

	// GetHostByName retrieves the host by name
	// Returns nil if not found.
	GetHostByName(name string) *rsc.Host

	// CreateHost creates a host instance on the system
	CreateHost(name string) (*rsc.Host, error)
}

func Update(r rsc.Rscer) rsc.Rscer {
	return rsc.Update(r)
}

func UpdateList(r rsc.RscLister) rsc.RscLister {
	return rsc.UpdateList(r)
}

// New creates a new Unity storage system instance.
func New(ip, username, password string) (Unity, error) {
	unity := rsc.NewUnity(ip, username, password)
	return unity, nil
}

func NewWithConn(conn *rsc.Connection) (Unity, error) {
	unity := rsc.NewUnityByConn(conn)
	return unity, nil
}
