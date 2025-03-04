// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"fmt"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// cleanBuild is a helper function to kill the build
// without execution. This will kill all resources,
// like steps and services, for the build in the
// configured backend.
func CleanBuild(ctx context.Context, database database.Interface, b *library.Build, services []*library.Service, steps []*library.Step, e error) {
	// update fields in build object
	b.SetError(fmt.Sprintf("unable to publish to queue: %s", e.Error()))
	b.SetStatus(constants.StatusError)
	b.SetFinished(time.Now().UTC().Unix())

	// send API call to update the build
	b, err := database.UpdateBuild(ctx, b)
	if err != nil {
		logrus.Errorf("unable to kill build %d: %v", b.GetNumber(), err)
	}

	for _, s := range services {
		// update fields in service object
		s.SetStatus(constants.StatusKilled)
		s.SetFinished(time.Now().UTC().Unix())

		// send API call to update the service
		_, err := database.UpdateService(ctx, s)
		if err != nil {
			logrus.Errorf("unable to kill service %s for build %d: %v", s.GetName(), b.GetNumber(), err)
		}
	}

	for _, s := range steps {
		// update fields in step object
		s.SetStatus(constants.StatusKilled)
		s.SetFinished(time.Now().UTC().Unix())

		// send API call to update the step
		_, err := database.UpdateStep(s)
		if err != nil {
			logrus.Errorf("unable to kill step %s for build %d: %v", s.GetName(), b.GetNumber(), err)
		}
	}
}
