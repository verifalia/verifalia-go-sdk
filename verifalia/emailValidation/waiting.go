package emailValidation

/*
* Verifalia - Email list cleaning and real-time email verification service
* https://verifalia.com/
* support@verifalia.com
*
* Copyright (c) 2005-2024 Cobisi Research
*
* Cobisi Research
* Via Della Costituzione, 31
* 35010 Vigonza
* Italy - European Union
*
* Permission is hereby granted, free of charge, to any person obtaining a copy
* of this software and associated documentation files (the "Software"), to deal
* in the Software without restriction, including without limitation the rights
* to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
* copies of the Software, and to permit persons to whom the Software is
* furnished to do so, subject to the following conditions:
*
* The above copyright notice and this permission notice shall be included in
* all copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
* THE SOFTWARE.
 */

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

	// Defines how much time to ask the Verifalia API to wait for the completion of the job on the server side, while polling
	// for the job.
	PollWaitTime time.Duration

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

	var retrievalOptions *RetrievalOptions

	if options != nil {
		retrievalOptions = &RetrievalOptions{RetrievalWaitTime: options.PollWaitTime}
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

		// Retrieve the updated job

		current, err = client.GetWithOptions(current.Overview.Id, retrievalOptions)

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
