// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package logsender

import (
	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/logsender"
	"github.com/juju/juju/worker"
	"github.com/juju/juju/worker/dependency"
	"github.com/juju/juju/worker/util"
)

// ManifoldConfig defines the names of the manifolds on which a Manifold will
// depend.
type ManifoldConfig struct {
	util.PostUpgradeManifoldConfig
	LogSource LogRecordCh
}

// Manifold returns a dependency manifold that runs a logger
// worker, using the resource names defined in the supplied config.
func Manifold(config ManifoldConfig) dependency.Manifold {
	return dependency.Manifold{
		Inputs: []string{
			config.APICallerName,
		},
		Start: func(getResource dependency.GetResourceFunc) (worker.Worker, error) {
			if config.UpgradeWaiterName != util.UpgradeWaitNotRequired {
				var upgradesDone bool
				if err := getResource(config.UpgradeWaiterName, &upgradesDone); err != nil {
					return nil, err
				}
				if !upgradesDone {
					return nil, dependency.ErrMissing
				}
			}

			var apiCaller base.APICaller
			if err := getResource(config.APICallerName, &apiCaller); err != nil {
				return nil, err
			}

			return New(config.LogSource, logsender.NewAPI(apiCaller)), nil
		},
	}
}
