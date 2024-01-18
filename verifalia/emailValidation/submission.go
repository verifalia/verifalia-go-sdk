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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ericlagergren/decimal"
	"github.com/verifalia/verifalia-go-sdk/verifalia/common"
	"github.com/verifalia/verifalia-go-sdk/verifalia/rest"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"time"
)

type ValidationRequestEntry struct {
	InputData string `json:"inputData"`
	Custom    string `json:"custom,omitempty"`
}

// SubmissionOptions allows to define generic submission options for an e-mail verification job.
type SubmissionOptions struct {
	// A context.Context that can cancel the validation. Useful if you wish to implement a timeout logic that abort the
	// request if it takes too long.
	Context context.Context

	// An optional user-defined name for the validation job, for your own reference.
	Name string

	// A reference to the quality level this job was validated against. The Quality enum-like object contains
	// the supported values, for example: Quality.High
	Quality string

	// The deduplication mode for the validation job. The Deduplication enum-like object contains
	// the supported values, for example: Deduplication.Relaxed
	Deduplication string

	// The eventual priority (speed) of the validation job, relative to the parent Verifalia account. In the event of an account
	// with many concurrent validation jobs, this value allows to increase the processing speed of a job with respect to the others.
	// The enum-like struct Priority contains some useful values you may want to use, should you need to specify a priority
	// for your job.
	Priority *uint8

	// The maximum data retention period Verifalia observes for this verification job, after which the job will be
	// automatically deleted.
	// A verification job can be deleted anytime prior to its retention period through the Delete() function.
	Retention time.Duration

	// An optional URL which Verifalia will invoke once the results for this job are ready.
	CompletionCallback url.URL

	// TODO: Introduce a new CompletionCallback object which includes the additional options available in Verifalia API v2.4+ (breaking change)

	// Defines how much time to ask the Verifalia API to wait for the completion of the job on the server side, during the
	// initial job submission request.
	SubmissionWaitTime time.Duration
}

// FileSubmissionOptions allows to define file-specific submission options for an e-mail verification job.
type FileSubmissionOptions struct {
	// The MIME Content-Type of the file data. The ContentType enum-like object contains the supported values, for
	// example: ContentType.ExcelXlsx
	ContentType string

	// The zero-based index of the first row to import and process.
	StartingRow int

	// An optional, zero-based index of the last row to import and process. If not specified, Verifalia will process
	// rows until the end of the file.
	EndingRow *int

	// The zero-based index of the column to import; applies to comma-separated (.csv), tab-separated (.tsv)
	// and other delimiter-separated values files, and Excel files.
	Column int

	// The zero-based index of the worksheet to import; applies to Excel files only.
	Sheet int

	// Allows to specify the line ending sequence of the provided file; applies to plain-text files, comma-separated
	// (.csv), tab-separated (.tsv) and other delimiter-separated values files. The LineEnding enum-like object contains
	// the supported values, for example: LineEnding.CrLf
	LineEnding string

	// An optional string with the column delimiter sequence of the file; applies to comma-separated (.csv),
	// tab-separated (.tsv) and other delimiter-separated values files. If not specified, Verifalia will use the `,`
	// (comma) symbol for CSV files and the `\t` (tab) symbol for TSV files.
	Delimiter string
}

// Internal struct used to serialize the validation request

type validationRequestBase struct {
	Name          *string `json:"name,omitempty"`
	Quality       *string `json:"quality,omitempty"`
	Deduplication *string `json:"deduplication,omitempty"`
	Priority      *uint8  `json:"priority,omitempty"`
	Retention     *string `json:"retention,omitempty"`
	Callback      *struct {
		Url string `json:"url"`
	} `json:"callback,omitempty"`
}

type validationRequest struct {
	validationRequestBase
	Entries []ValidationRequestEntry `json:"entries"`
}

type fileValidationRequest struct {
	validationRequestBase
	StartingRow int    `json:"startingRow,omitempty"`
	EndingRow   *int   `json:"endingRowRow,omitempty"`
	Column      int    `json:"column,omitempty"`
	Sheet       int    `json:"sheet,omitempty"`
	LineEnding  string `json:"lineEnding,omitempty"`
	Delimiter   string `json:"delimiter,omitempty"`
}

// LineEnding provides enumerated-like values for the line-ending modes for an input text file provided to the
// Verifalia API for verification.
var LineEnding = struct {
	Auto string
	CrLf string
	Cr   string
	Lf   string
}{
	// Automatic line-ending detection, attempts to guess the correct line ending from the first chunk of data.
	Auto: "",

	// CR + LF sequence (\r\n), commonly used in files generated on Windows.
	CrLf: "CrLf",

	// CR sequence (\r), commonly used in files generated on classic macOS.
	Cr: "Cr",

	// LF (\n), commonly used in files generated on Unix and Unix-like systems (including Linux and macOS).
	Lf: "Lf",
}

// Internal struct used to deserialize the validation response
type overview struct {
	Id            string     `json:"id"`
	SubmittedOn   time.Time  `json:"submittedOn"`
	CompletedOn   *time.Time `json:"completedOn"`
	Priority      *uint8     `json:"priority"`
	Name          string     `json:"name"`
	Owner         string     `json:"owner"`
	ClientIP      string     `json:"clientIP"`
	CreatedOn     time.Time  `json:"createdOn"`
	Quality       string     `json:"quality"`
	Retention     string     `json:"retention"`
	Deduplication string     `json:"deduplication"`
	Status        string     `json:"status"`
	NoOfEntries   uint       `json:"noOfEntries"`
	Progress      *struct {
		Percentage             decimal.Big `json:"percentage"`
		EstimatedTimeRemaining string      `json:"estimatedTimeRemaining"`
	} `json:"progress"`
}

type partialJob struct {
	Overview overview                      `json:"overview"`
	Entries  *common.ListingSegment[Entry] `json:"entries"`
}

// Submit starts processing a new single e-mail verification; this function does not wait for the completion of the email validation
// job: use the WaitForCompletion() function to do that.
func (client *Client) Submit(inputData string) (*Job, error) {
	return client.SubmitWithOptions(ValidationRequestEntry{
		InputData: inputData,
	}, nil)
}

// SubmitWithOptions starts processing a new single e-mail verification; this function does not wait for the completion of the email validation
// job: use the WaitForCompletion() function to do that.
func (client *Client) SubmitWithOptions(entry ValidationRequestEntry, options *SubmissionOptions) (*Job, error) {
	return client.SubmitManyWithOptions([]ValidationRequestEntry{
		entry,
	}, options)
}

// SubmitMany starts processing a new verification with multiple e-mail addresses; this function does not wait for the completion of the email validation
// job: use the WaitForCompletion() function to do that.
func (client *Client) SubmitMany(inputData []string) (*Job, error) {
	var entries = make([]ValidationRequestEntry, len(inputData))

	for i, item := range inputData {
		entries[i] = ValidationRequestEntry{
			InputData: item,
		}
	}

	return client.SubmitManyWithOptions(entries, nil)
}

// SubmitManyWithOptions starts processing a new verification with multiple e-mail addresses; this function does not wait for the completion of the email validation
// job: use the WaitForCompletion() function to do that.
func (client *Client) SubmitManyWithOptions(entries []ValidationRequestEntry, options *SubmissionOptions) (*Job, error) {
	var ctx context.Context

	request := validationRequest{
		Entries: entries,
	}

	fillSubmissionRequestOptions(&request.validationRequestBase, options)

	jsonData, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	// Invoke the API through the common submission code path

	var queryParams map[string][]string

	if options != nil {
		ctx = options.Context

		queryParams = make(map[string][]string)
		queryParams["waitTime"] = []string{fmt.Sprintf("%v", options.SubmissionWaitTime.Seconds())}
	}

	return client.submit(rest.InvocationOptions{
		Method:      http.MethodPost,
		Resource:    "email-validations",
		QueryParams: queryParams,
		Body:        bytes.NewReader(jsonData),
		Context:     ctx,
	})
}

// SubmitFile starts processing a new verification from a file; this function does not wait for the completion of the email validation
// job: use the WaitForCompletion() function to do that.
func (client *Client) SubmitFile(file os.File) (*Job, error) {

	return client.SubmitFileWithOptions(file, nil, nil)
}

// SubmitFileWithOptions starts processing a new verification from a file; this function does not wait for the completion of the email validation
// job: use the WaitForCompletion() function to do that.
func (client *Client) SubmitFileWithOptions(file os.File, fileOptions *FileSubmissionOptions, options *SubmissionOptions) (*Job, error) {
	// If the caller did not specify a content type we will try to guess it based on the file extension

	if fileOptions == nil {
		fileOptions = &FileSubmissionOptions{}
	}

	if fileOptions.ContentType == "" {
		fileOptions.ContentType = guessContentType(path.Ext(file.Name()))

		if fileOptions.ContentType == "" {
			return nil, fmt.Errorf("cannot guess the content type for the provided file %v, please specify it through the options arg", file.Name())
		}
	}

	return client.SubmitFileReaderWithOptions(&file, fileOptions, options)
}

// SubmitFileReaderWithOptions starts processing a new verification from a file reader; this function does not wait for the completion of the email validation
// job: use the WaitForCompletion() function to do that.
func (client *Client) SubmitFileReaderWithOptions(reader io.Reader, fileOptions *FileSubmissionOptions, options *SubmissionOptions) (*Job, error) {
	contentType := rest.ContentType.TextPlain

	if fileOptions != nil && fileOptions.ContentType != "" {
		contentType = fileOptions.ContentType
	}

	// Create the multipart request body

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Input file part

	fileHeader := make(textproto.MIMEHeader)
	fileHeader.Set("Content-Type", contentType)
	fileHeader.Set("Content-Disposition", fmt.Sprintf("form-data; name=\"inputFile\"; filename=\"dummy\""))

	filePart, err := writer.CreatePart(fileHeader)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(filePart, reader)
	if err != nil {
		return nil, err
	}

	// Json settings part

	optionsHeader := make(textproto.MIMEHeader)
	optionsHeader.Set("Content-Type", rest.ContentType.ApplicationJson)
	optionsHeader.Set("Content-Disposition", fmt.Sprintf("form-data; name=\"settings\""))

	optionsPart, err := writer.CreatePart(optionsHeader)
	if err != nil {
		return nil, err
	}

	var ctx context.Context

	request := fileValidationRequest{
		StartingRow: fileOptions.StartingRow,
		EndingRow:   fileOptions.EndingRow,
		Column:      fileOptions.Column,
		Sheet:       fileOptions.Sheet,
		LineEnding:  fileOptions.LineEnding,
		Delimiter:   fileOptions.Delimiter,
	}

	fillSubmissionRequestOptions(&request.validationRequestBase, options)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	_, err = optionsPart.Write(jsonData)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Invoke the API through the common submission code path

	if options != nil {
		ctx = options.Context
	}

	queryParams := make(map[string][]string)
	queryParams["waitTime"] = []string{fmt.Sprintf("%v", options.SubmissionWaitTime.Seconds())}

	return client.submit(rest.InvocationOptions{
		Method: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
		Resource:    "email-validations",
		QueryParams: queryParams,
		Body:        body,
		Context:     ctx,
	})
}

func fillSubmissionRequestOptions(request *validationRequestBase, options *SubmissionOptions) {
	if options != nil {
		if options.Name != "" {
			request.Name = &options.Name
		}
		if options.Quality != "" {
			request.Quality = &options.Quality
		}
		if options.Deduplication != "" {
			request.Deduplication = &options.Deduplication
		}
		if options.Priority != nil {
			request.Priority = options.Priority
		}
		if options.Retention != 0 {
			retentionAsTimeSpan := common.DurationToTimeSpanString(options.Retention)
			request.Retention = &retentionAsTimeSpan
		}
		if options.CompletionCallback.Scheme != "" {
			request.Callback = &struct {
				Url string `json:"url"`
			}{
				options.CompletionCallback.String(),
			}
		}
	}
}

func guessContentType(extension string) string {
	switch extension {
	case ".txt":
		return rest.ContentType.TextPlain
	case ".csv":
		return rest.ContentType.TextCsv
	case ".tsv":
		return rest.ContentType.TextTsv
	case ".tab":
		return rest.ContentType.TextTsv
	case ".xls":
		return rest.ContentType.ExcelXls
	case ".xlsx":
		return rest.ContentType.ExcelXlsx
	default:
		return ""
	}
}

func (client *Client) submit(invocationOptions rest.InvocationOptions) (*Job, error) {
	response, err := client.RestClient.Invoke(invocationOptions)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusAccepted {
		responseData, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return nil, err
		}

		// Unmarshal a simplified projection of the original object, without segmentation/metadata info

		var partial partialJob

		if err := json.Unmarshal(responseData, &partial); err != nil {
			return nil, err
		}

		result, err := client.buildJob(partial, invocationOptions.Context)

		return result, err
	}

	return nil, fmt.Errorf("unexpected HTTP response: %d", response.StatusCode)
}
