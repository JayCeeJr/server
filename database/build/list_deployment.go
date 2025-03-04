// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// ListBuildsForDeployment gets a list of builds by deployment url from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *engine) ListBuildsForDeployment(ctx context.Context, d *library.Deployment, filters map[string]interface{}, page, perPage int) ([]*library.Build, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"deployment": d.GetURL(),
	}).Tracef("listing builds for deployment %s from the database", d.GetURL())

	// variables to store query results and return values
	count := int64(0)
	b := new([]database.Build)
	builds := []*library.Build{}

	// count the results
	count, err := e.CountBuildsForDeployment(ctx, d, filters)
	if err != nil {
		return builds, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return builds, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	err = e.client.
		Table(constants.TableBuild).
		Where("source = ?", d.GetURL()).
		Where(filters).
		Order("number DESC").
		Limit(perPage).
		Offset(offset).
		Find(&b).
		Error
	if err != nil {
		return nil, count, err
	}

	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Build.ToLibrary
		builds = append(builds, tmp.ToLibrary())
	}

	return builds, count, nil
}
