// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"fmt"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

// PlanSteps is a helper function to plan all steps
// in the build for execution. This creates the steps
// for the build in the configured backend.
func PlanSteps(ctx context.Context, database database.Interface, p *pipeline.Build, b *library.Build) ([]*library.Step, error) {
	// variable to store planned steps
	steps := []*library.Step{}

	// iterate through all pipeline stages
	for _, stage := range p.Stages {
		// iterate through all steps for each pipeline stage
		for _, step := range stage.Steps {
			// create the step object
			s, err := planStep(ctx, database, b, step, stage.Name)
			if err != nil {
				return steps, err
			}

			steps = append(steps, s)
		}
	}

	// iterate through all pipeline steps
	for _, step := range p.Steps {
		s, err := planStep(ctx, database, b, step, "")
		if err != nil {
			return steps, err
		}

		steps = append(steps, s)
	}

	return steps, nil
}

func planStep(ctx context.Context, database database.Interface, b *library.Build, c *pipeline.Container, stage string) (*library.Step, error) {
	// create the step object
	s := new(library.Step)
	s.SetBuildID(b.GetID())
	s.SetRepoID(b.GetRepoID())
	s.SetNumber(c.Number)
	s.SetName(c.Name)
	s.SetImage(c.Image)
	s.SetStage(stage)
	s.SetStatus(constants.StatusPending)
	s.SetCreated(time.Now().UTC().Unix())

	// send API call to create the step
	s, err := database.CreateStep(s)
	if err != nil {
		return nil, fmt.Errorf("unable to create step %s: %w", s.GetName(), err)
	}

	// populate environment variables from step library
	//
	// https://pkg.go.dev/github.com/go-vela/types/library#step.Environment
	err = c.MergeEnv(s.Environment())
	if err != nil {
		return nil, err
	}

	// create the log object
	l := new(library.Log)
	l.SetStepID(s.GetID())
	l.SetBuildID(b.GetID())
	l.SetRepoID(b.GetRepoID())
	l.SetData([]byte{})

	// send API call to create the step logs
	err = database.CreateLog(ctx, l)
	if err != nil {
		return nil, fmt.Errorf("unable to create logs for step %s: %w", s.GetName(), err)
	}

	return s, nil
}
