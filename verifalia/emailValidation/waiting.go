package emailValidation

import (
	"context"
	"time"
)

// WaitingOptions allows to customize how the waiting logic behaves.
type WaitingOptions struct {
	// A context.Context that can cancel the waiting process. Useful if you wish to implement a timeout logic that abort the
	// request if it takes too long.
	Context context.Context

	// A waiting function which pauses the current execution until the time of the next polling.
	WaitForNextPoll func(overview Overview, ctx context.Context) error

	// TODO: Progress reporting
}

// WaitForCompletion sleeps until the e-mail verification job completes.
func (client *Client) WaitForCompletion(validation *Job) (result *Job, err error) {
	return client.WaitForCompletionWithOptions(validation, nil)
}

// WaitForCompletionWithOptions sleeps until the e-mail verification job completes.
func (client *Client) WaitForCompletionWithOptions(validation *Job, options *WaitingOptions) (current *Job, err error) {
	var ctx context.Context
	current = validation

	if options != nil {
		ctx = options.Context
	}

	if ctx == nil {
		ctx = context.TODO()
	}

	for {
		if current.Overview.Status != JobStatus.InProgress {
			break
		}

		// Wait for the polling interval (or the context cancellation)

		if options != nil && options.WaitForNextPoll != nil {
			err = options.WaitForNextPoll(current.Overview, ctx)
		} else {
			err = defaultWaitForNextPoll(current.Overview, ctx)
		}

		if err != nil {
			return nil, err
		}

		// Retrieve the updated validation

		current, err = client.Get(current.Overview.Id)

		if err != nil {
			return nil, err
		}
	}

	return current, nil
}

func defaultWaitForNextPoll(overview Overview, ctx context.Context) error {
	// TODO: observe the job ETA if we have one

	// Either wait for the polling interval or until the context is cancelled

	select {
	case <-time.After(5 * time.Second):
	case <-ctx.Done():
	}

	// If the context is cancelled, let the caller know

	return ctx.Err()
}
