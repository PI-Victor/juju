// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver

import (
	"fmt"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/params"
)

type adminApiV3 struct {
	*admin
}

func newAdminApiV3(srv *Server, root *apiHandler, reqNotifier *requestNotifier) interface{} {
	return &adminApiV3{
		&admin{
			srv:         srv,
			root:        root,
			reqNotifier: reqNotifier,
		},
	}
}

// Admin returns an object that provides API access to methods that can be
// called even when not authenticated.
func (r *adminApiV3) Admin(id string) (*adminApiV3, error) {
	if id != "" {
		// Safeguard id for possible future use.
		return nil, common.ErrBadId
	}
	return r, nil
}

// Login logs in with the provided credentials.  All subsequent requests on the
// connection will act as the authenticated user.
func (a *adminApiV3) Login(req params.LoginRequest) (params.LoginResultV1, error) {
	return a.doLogin(req, 3)
}

// RedirectInfo returns redirected host information for the model.
// In Juju it always returns an error because the Juju controller
// does not multiplex controllers.
func (a *adminApiV3) RedirectInfo() (params.RedirectInfoResult, error) {
	return params.RedirectInfoResult{}, fmt.Errorf("not redirected")
}
