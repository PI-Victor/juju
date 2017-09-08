// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package firewall

import (
	"github.com/juju/errors"

	"github.com/juju/juju/network"
	"github.com/juju/juju/state"
	"github.com/juju/juju/status"
)

// State provides the subset of global state required by the
// remote firewaller facade.
type State interface {
	state.ModelMachinesWatcher
	state.ModelAccessor

	WatchSubnets(func(id interface{}) bool) state.StringsWatcher

	KeyRelation(string) (Relation, error)

	Unit(string) (Unit, error)

	Machine(string) (Machine, error)

	Application(string) (Application, error)
}

// TODO(wallyworld) - for tests, remove when remaining firewaller tests become unit tests.
func StateShim(st *state.State) stateShim {
	return stateShim{st}
}

type stateShim struct {
	*state.State
}

func (st stateShim) KeyRelation(key string) (Relation, error) {
	rel, err := st.State.KeyRelation(key)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return relationShim{rel}, nil
}

type Relation interface {
	status.StatusSetter
	Endpoints() []state.Endpoint
	WatchUnits(applicationName string) (state.RelationUnitsWatcher, error)
	UnitInScope(Unit) (bool, error)
	WatchRelationIngressNetworks() state.StringsWatcher
	WatchRelationEgressNetworks() state.StringsWatcher
}

type relationShim struct {
	*state.Relation
}

func (r relationShim) UnitInScope(u Unit) (bool, error) {
	ru, err := r.Relation.Unit(u.(*state.Unit))
	if err != nil {
		return false, errors.Trace(err)
	}
	return ru.InScope()
}

func (st stateShim) Application(name string) (Application, error) {
	app, err := st.State.Application(name)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return applicationShim{app}, nil
}

type Application interface {
	Name() string
	AllUnits() ([]Unit, error)
}

type applicationShim struct {
	*state.Application
}

func (a applicationShim) AllUnits() (results []Unit, err error) {
	units, err := a.Application.AllUnits()
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, unit := range units {
		results = append(results, unit)
	}
	return results, nil
}

type Unit interface {
	Name() string
	PublicAddress() (network.Address, error)
	AssignedMachineId() (string, error)
}

func (st stateShim) Unit(name string) (Unit, error) {
	return st.State.Unit(name)
}

type Machine interface {
	Id() string
	WatchAddresses() state.NotifyWatcher
}

func (st stateShim) Machine(id string) (Machine, error) {
	return st.State.Machine(id)
}
