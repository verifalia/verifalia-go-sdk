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
	"io"
	"os"
)

// Run verifies a new single e-mail address; this function automatically waits for the completion of the email validation
// job: should you need to handle the waiting process manually, use a combination of Submit() and WaitForCompletion().
func (client *Client) Run(inputData string) (*Job, error) {
	// Submission

	validation, err := client.Submit(inputData)

	if err != nil {
		return validation, err
	}

	// Waiting

	return client.WaitForCompletion(validation)
}

// RunWithOptions verifies a new single e-mail address; this function automatically waits for the completion of the email validation
// job: should you need to handle the waiting process manually, use a combination of SubmitWithOptions() and WaitForCompletionWithOptions().
func (client *Client) RunWithOptions(entry ValidationRequestEntry, options *SubmissionOptions, waitingOptions *WaitingOptions) (*Job, error) {
	// Submission

	validation, err := client.SubmitWithOptions(entry, options)

	if err != nil {
		return validation, err
	}

	// Waiting

	return client.WaitForCompletionWithOptions(validation, waitingOptions)
}

// RunMany verifies multiple e-mail addresses; this function automatically waits for the completion of the email validation
// job: should you need to handle the waiting process manually, use a combination of SubmitMany() and WaitForCompletion().
func (client *Client) RunMany(inputData []string) (*Job, error) {
	// Submission

	validation, err := client.SubmitMany(inputData)

	if err != nil {
		return validation, err
	}

	// Waiting

	return client.WaitForCompletion(validation)
}

// RunManyWithOptions verifies multiple e-mail addresses; this function automatically waits for the completion of the email validation
// job: should you need to handle the waiting process manually, use a combination of SubmitManyWithOptions() and WaitForCompletionWithOptions().
func (client *Client) RunManyWithOptions(entries []ValidationRequestEntry, options *SubmissionOptions, waitingOptions *WaitingOptions) (*Job, error) {
	// Submission

	validation, err := client.SubmitManyWithOptions(entries, options)

	if err != nil {
		return validation, err
	}

	// Waiting

	return client.WaitForCompletionWithOptions(validation, waitingOptions)
}

// RunFile verifies a file containing email addresses; this function automatically waits for the completion of the email validation
// job: should you need to handle the waiting process manually, use a combination of SubmitFile() and WaitForCompletion().
func (client *Client) RunFile(file os.File) (*Job, error) {
	// Submission

	validation, err := client.SubmitFile(file)

	if err != nil {
		return validation, err
	}

	// Waiting

	return client.WaitForCompletion(validation)
}

// RunFileWithOptions verifies a file containing email addresses; this function automatically waits for the completion of the email validation
// job: should you need to handle the waiting process manually, use a combination of SubmitFileWithOptions() and WaitForCompletionWithOptions().
func (client *Client) RunFileWithOptions(file os.File, fileOptions *FileSubmissionOptions, options *SubmissionOptions, waitingOptions *WaitingOptions) (*Job, error) {
	// Submission

	validation, err := client.SubmitFileWithOptions(file, fileOptions, options)

	if err != nil {
		return validation, err
	}

	// Waiting

	return client.WaitForCompletionWithOptions(validation, waitingOptions)
}

// RunFileReaderWithOptions verifies a file containing email addresses, using an io.Reader; this function automatically waits for the completion of the email validation
// job: should you need to handle the waiting process manually, use a combination of RunFileReaderWithOptions() and WaitForCompletionWithOptions().
func (client *Client) RunFileReaderWithOptions(reader io.Reader, fileOptions *FileSubmissionOptions, options *SubmissionOptions, waitingOptions *WaitingOptions) (*Job, error) {
	// Submission

	validation, err := client.SubmitFileReaderWithOptions(reader, fileOptions, options)

	if err != nil {
		return validation, err
	}

	// Waiting

	return client.WaitForCompletionWithOptions(validation, waitingOptions)
}
