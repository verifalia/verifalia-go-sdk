package emailValidation

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
