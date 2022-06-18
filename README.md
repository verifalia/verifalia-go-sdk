![Verifalia API](https://img.shields.io/badge/Verifalia%20API-v2.3-green)
![Go version](https://img.shields.io/badge/Go-%3E=1.18-green)

Verifalia API - Go SDK and helper library
=========================================

**The perfect starting point to integrate Verifalia within your Go project**

[Verifalia][0] provides a simple HTTPS-based API for **validating email addresses in real-time** and checking whether they are deliverable or not; this SDK library integrates with Verifalia and allows to [verify email addresses][0] under **Go v1.18 and higher**.

To learn more about Verifalia please see [https://verifalia.com][0]

## Getting started ##

First, add the Verifalia Go SDK as a new module to your Go project:

```bash
# First line is optional if your project is already defined as a Go module
go mod init <YOUR_PROJECT_NAME>
go get github.com/verifalia/verifalia-go-sdk/verifalia
```

To update the SDK use `go get -u` to retrieve the latest version of the SDK:

```bash
go get -u github.com/verifalia/verifalia-go-sdk/verifalia
```

### Authentication ###

First things first: authentication to the Verifalia API is performed by way of either
the credentials of your root Verifalia account or of one of your users (previously
known as sub-accounts): if you don't have a Verifalia account, just [register for a free one][4]. For security reasons, it is always advisable to [create and use a dedicated user][3] for accessing the API, as doing so will allow to assign only the specific needed permissions to it.

Learn more about authenticating to the Verifalia API at [https://verifalia.com/developers#authentication][2]

Once you have your Verifalia credentials at hand, use them to create a new `client` object, which will be the starting point to every other operation against the Verifalia API: the supplied credentials will be automatically provided to the API using the HTTP Basic Auth method.

```go
package main

import (
    "github.com/verifalia/verifalia-go-sdk/verifalia"
)

func main() {
    client := verifalia.NewClient("username", "password")

    // TODO: Use "client" as explained below
}
```

In addition to the HTTP Basic Auth method, this SDK also supports other different ways to authenticate to the Verifalia API, as explained in the subsequent paragraphs.

#### Authenticating via X.509 client certificate (TLS mutual authentication)

This authentication method uses a cryptographic X.509 client certificate to authenticate against the Verifalia API, through the TLS protocol. This method, also called mutual TLS authentication (mTLS) or two-way authentication, offers the highest degree of security, as only a cryptographically-derived key (and not the actual credentials) is sent over the wire on each request.

```go
package main

import (
    "crypto/tls"
    "github.com/verifalia/verifalia-go-sdk/verifalia"
)

func main() {
    cert, err := tls.LoadX509KeyPair("./mycertificate.pem", "./mycertificate.key")

    if err != nil {
        panic(err)
    }

    client := verifalia.NewClientWithCertificateAuth(&cert)
	
    // TODO: Use "client" as explained below
}
```

## Validating email addresses ##

Every operation related to verifying / validating email addresses is performed through the `EmailValidation` field exposed by the `client` instance you created above. The property exposes some useful functions: in the next few paragraphs we are looking at the most used ones, so it is strongly advisable to explore the library and look at the embedded help for other opportunities.

### How to validate / verify an email address ###

To verify an email address from a Go application you can call the `Run()` function exposed by
the `client.EmailValidation` field, as shown below:

```go
package main

import (
    "fmt"
    "github.com/verifalia/verifalia-go-sdk/verifalia"
)

func main() {
    client := verifalia.NewClient("<USERNAME>", "<PASSWORD>") // See above

    // Verifies an email address

    validation, err := client.EmailValidation.Run("batman@gmail.com")

    if err != nil {
        panic(err)
    }

    // Print some results

    entry := validation.Entries[0]
    fmt.Printf("%v => %v", entry.EmailAddress, entry.Classification)

    // Output:
    // batman@gmail.com => Deliverable
}
```

Once `Run()` completes successfully, the resulting verification job
is guaranteed to be completed and its results' data (e.g. its `Entries` field) to be available for use.

#### Advanced processing options

You can also specify additional processing options, like the results quality vs. processing time trade-off level Verifalia must consider while working on
your data, or the data retention policy Verifalia must obey for the verification job, a `context.Context` which can
limit the waiting time and several other processing details. To do that, call the `RunWithOptions()` function, a shown in the example below:

```go
package main

import (
    "fmt"
    "github.com/verifalia/verifalia-go-sdk/verifalia"
    "github.com/verifalia/verifalia-go-sdk/verifalia/emailValidation"
)

func main() {
    client := verifalia.NewClient("<USERNAME>", "<PASSWORD>") // See above

    // Configure the validation options 
  
    options := emailValidation.SubmissionOptions{
        // High quality results
        Quality:   emailValidation.Quality.High,
        // 30 minutes of data retention
        Retention: 30 * time.Minute,
    }

    // Verifies an email address, using the above options
    
    validation, err := client.EmailValidation.RunWithOptions("batman@gmail.com", options)

    if err != nil {
        panic(err)
    }

    // Print some results

    entry := validation.Entries[0]
    fmt.Printf("%v => %v", entry.EmailAddress, entry.Classification)

    // Output:
    // batman@gmail.com => Deliverable
}
```

### How to validate / verify a list of email addresses ###

To verify a list of email addresses you can call the `RunMany()` function, which accepts an array of strings with email addresses to verify:

```go
package main

import (
    "fmt"
    "github.com/verifalia/verifalia-go-sdk/verifalia"
)

func main() {
    client := verifalia.NewClient("<USERNAME>", "<PASSWORD>") // See above

    // Verifies the list of email addresses

    validation, err := client.EmailValidation.RunMany([]string{
        "batman@gmail.com",
        "steve.vai@best.music",
        "samantha42@yahoo.de",
    })

    if err != nil {
        panic(err)
    }

    // Print some results

    for _, entry := range validation.Entries {
        fmt.Printf("%v => %v\n",
            entry.EmailAddress,
            entry.Classification)
    }

    // Output:
    // batman@gmail.com => Deliverable
    // steve.vai@best.music => Undeliverable
    // samantha42@yahoo.de => Deliverable
}
```


#### Advanced processing options

Similarly to the `RunWithOptions()` function described above, you can use the `RunManyWithOptions()` function
to specify additional processing options, like the results quality vs. processing time trade-off level Verifalia must consider while working on
your data, or the data retention policy Verifalia must obey for the verification job, a `context.Context` which can
limit the waiting time and several other processing details.

### How to import and verify a file with a list of email addresses

This library also includes support for submitting and validating files with email addresses, including:

- **plain text files** (.txt), with one email address per line;
- **comma-separated values** (.csv), **tab-separated values** (.tsv) and other delimiter-separated values files;
- **Microsoft Excel spreadsheets** (.xls and .xlsx).

To import and verify such files, one can call the `RunFile()` function passing the desired
`os.File` instance (for `io.Reader` support please read below). Along with that, it is also possible to call the `RunFileWithOptions()` function to specify the eventual starting
and ending rows to process, the column, the sheet index, the line ending and the
delimiter - depending of course on the nature of the submitted file.

Here is how to verify an Excel file, for example:

```go
package main

import (
    "fmt"
    "github.com/verifalia/verifalia-go-sdk/verifalia"
    "os"
)

func main() {
    client := verifalia.NewClient("<USERNAME>", "<PASSWORD>") // See above
	
    // Verifies an Excel file with a list of email addresses

    thatFile, err := os.Open("that-file.xslx")
	
    if err != nil {
        panic(err)
    }

    validation, err = client.EmailValidation.RunFile(thatFile)

    if err != nil {
        panic(err)
    }

    // Print some results

    for _, entry := range validation.Entries {
        fmt.Printf("%v => %v\n",
            entry.EmailAddress,
            entry.Classification)
    }
}
```

#### Advanced processing options

Similarly to the aforementioned `RunWithOptions()` and `RunManyWithOptions()` functions, you can use the `RunFileWithOptions()` function
to specify additional processing options, like the results quality vs. processing time trade-off level Verifalia must consider while working on
your data, or the data retention policy Verifalia must obey for the verification job, a `context.Context` which can
limit the waiting time and several other processing details.
The function also accepts an additional `emailValidation.FileSubmissionOptions` instance which allows to further
define file-specific options, like the type of the file, the worksheet, row and column which
need to be imported, as well as the specific range of data to process.

Here is an example showing how to do that:

```go
package main

import (
    "fmt"
    "github.com/verifalia/verifalia-go-sdk/verifalia"
    "github.com/verifalia/verifalia-go-sdk/verifalia/emailValidation"
    "os"
)

func main() {
    client := verifalia.NewClient("<USERNAME>", "<PASSWORD>") // See above
	
    // Verifies an Excel file with a list of email addresses

    thatFile, err := os.Open("that-file.xslx")
	
    if err != nil {
        panic(err)
    }

    // Configure the validation options 

    options := emailValidation.SubmissionOptions{
        // High quality results
        Quality:   emailValidation.Quality.High,
        // One hour of data retention		
        Retention: 1 * time.Hour,
    }

    // Specify which data Verifalia should process through the file options 
	
    fileOptions := emailValidation.FileSubmissionOptions{
        // Second sheet
        Sheet:       1,
        // Ninth sheet      
        Column:      8,
        // Will start importing from the third row
        StartingRow: 2,
    }	

    validation, err = client.EmailValidation.RunFileWithOptions(thatFile,
        fileOptions,
        options,
        nil)

    if err != nil {
        panic(err)
    }

    // Print some results

    for _, entry := range validation.Entries {
        fmt.Printf("%v => %v\n",
            entry.EmailAddress,
            entry.Classification)
    }
}
```

For advanced usage scenarios, it is possible to import a file by way of a `io.Reader` instance. To do that, call the `RunFileReaderWithOptions()` function.
In that case, the library won't be able to guess the file type by way of its extension and will default to the `text/plain` value: make sure
to specify the correct `ContentType` through a `FileSubmissionOptions` instance, if you need to import a different type of file.

## Job lifecycle

Email verification jobs can take considerable processing time, depending on the number of email addresses they include, the required
quality level, the target mail exchangers under test and the Verifalia plan you are running on
(premium plans come with a faster processing speed).

Running an email verification job requires submitting it to the Verifalia API and eventually polling it until it completes:
while all the `Run*()` functions discussed so far hide all this complexity, it is possible to
manually handle this process in order to better integrate with an existing workflow or application.

### Submission

To manually handle the running process,
one can use one of the `Submit*()` functions which submit an email verification job to the Verifalia API: they
still return an `emailValidation.Job` instance like the `Run*()` functions do, with its results only available in the event
the job's `Status` field equals to `emailValidation.JobStatus.Completed`. Uncompleted jobs still
have their `Overview` available, along with its `Id` field which can be stored and later used to refer to the job
while using the Verifalia API.

```go
package main

import (
	"fmt"
	"github.com/verifalia/verifalia-go-sdk/verifalia"
    "github.com/verifalia/verifalia-go-sdk/verifalia/emailValidation"
)

func main() {
    client := verifalia.NewClient("<USERNAME>", "<PASSWORD>") // See above

    //  Submit an email address for verification

    validation, err := client.EmailValidation.Submit("batman@gmail.com")

    if err != nil {
        panic(err)
    }
	
    // Print the job Id

    fmt.Println(validation.Overview.Id)

    // Output:
    // 9ece66cf-916c-4313-9c40-b8a73f0ef872
	
    if validation.Overview.Status == emailValidation.Status.Completed {
        // validation.Entries will have the validation results!
    } else {
        // What about having a coffee?
    }
}
```

It is also possible to call the `Submit*WithOptions()` functions in order to specify 
additional processing options, like the results quality vs. processing time trade-off level Verifalia must consider while working on
your data, or the data retention policy Verifalia must obey for the verification job, a `context.Context` which can
limit the waiting time and several other processing details.

#### Completion callbacks

Along with each email validation job, it is possible to specify an URL which
Verifalia will invoke (POST) once the job completes: this URL must use the HTTPS or HTTP
scheme and be publicly accessible over the Internet.
To learn more about completion callbacks, please see https://verifalia.com/developers#email-validations-completion-callback

To specify a completion callback URL, use one of the `Submit*WithOptions()` or `Run*WithOptions`
functions and set the `CompletionCallback` of the specified `SubmissionOptions` instance.

Note that completion callbacks are invoked asynchronously and it could take up to
several seconds for your callback URL to get invoked.

### Retrieving a job

Once you have an email validation job `Id`, which is always returned by any of the `Submit*()` functions as part of the validation's `Overview` property, you can retrieve an updated snapshot of the job by way of the `Get()` function.

In the following example, we are requesting the current snapshot of a given email validation job back from Verifalia:

```go
package main

import (
    "fmt"
    "github.com/verifalia/verifalia-go-sdk/verifalia"
    "github.com/verifalia/verifalia-go-sdk/verifalia/emailValidation"
)

func main() {
    client := verifalia.NewClient("<USERNAME>", "<PASSWORD>") // See above

    // Retrieve an email validation job, given its Id

    validation, err := client.EmailValidation.Get("9ece66cf-916c-4313-9c40-b8a73f0ef872")

    if err != nil {
        panic(err)
    }
	
    if validation.Overview.Status == emailValidation.Status.Completed {
        // validation.Entries will have the validation results!
    } else {
        // What about having a coffee?
    }
}
```

### Waiting for completion

While the `Run*()` functions automatically wait of the completion of their email verification jobs,
developers who want to manually handle the running process and thus call the `Submit*()` functions
would need to keep calling `Get()` (or `GetWithOptions()`) until the job completes.

To ease this task,
this library includes the function `WaitForCompletion()`, which pauses the execution and automatically
waits until the specified job completes.

Here is an example showing how to do that:

```go
package main

import (
    "fmt"
    "github.com/verifalia/verifalia-go-sdk/verifalia"
)

func main() {
    client := verifalia.NewClient("<USERNAME>", "<PASSWORD>") // See above

    //  Submit an email address for verification

    validation, err := client.EmailValidation.Submit("batman@gmail.com")

    if err != nil {
        panic(err)
    }

    // Wait until the job completes

    validation, err = client.EmailValidation.WaitForCompletion(validation)

    if err != nil {
        panic(err)
    }

    // Print some results

    entry := validation.Entries[0]
    fmt.Printf("%v => %v", entry.EmailAddress, entry.Classification)

    // Output:
    // batman@gmail.com => Deliverable
}
```

### Don't forget to clean up, when you are done ###

Verifalia automatically deletes completed email verification jobs after a configurable
data-retention period (minimum 5 minutes, maximum 30 days) but it is strongly advisable that
you delete your completed jobs as soon as possible, for privacy and security reasons.

To do that, you can call the `Delete()` function, passing the job Id you wish to get rid of:

```go
package main

import (
    "fmt"
    "github.com/verifalia/verifalia-go-sdk/verifalia"
)

func main() {
    client := verifalia.NewClient("<USERNAME>", "<PASSWORD>") // See above

    //  Delete an email validation job, given its Id

    err := client.EmailValidation.Delete("9ece66cf-916c-4313-9c40-b8a73f0ef872")

    if err != nil {
        panic(err)
    }
}
```

Once deleted, a job is gone and there is no way to retrieve its email validation results.

## Iterating over your email validation jobs

For management and reporting purposes, you may want to obtain a detailed list of
your past email validation jobs. This SDK library allows to do that through
the `List()` function, which returns a channel that allows to iterate asynchronously
over a collection of `emailValidation.Overview` instances (the same type of the `Overview` property of the results returned by `Submit*()`, `Run*()`, `Get*()` functions and alike).

The `ListWithOptions()` function allows to specify the listing options, including the
sorting direction of the desired results.

Here is how to iterate over your jobs, from the most recent to the oldest one:

```go
package main

import (
    "fmt"
    "github.com/verifalia/verifalia-go-sdk/verifalia"
    "github.com/verifalia/verifalia-go-sdk/verifalia/common"
    "github.com/verifalia/verifalia-go-sdk/verifalia/emailValidation"
)

func main() {
    client := verifalia.NewClient("<USERNAME>", "<PASSWORD>") // See above
	
    // Configure the options to have the most recent job first

    listingOptions := emailValidation.ListingOptions{
        Direction: common.Backward,
    }

    // Proceed with the asynchronous listing

    results := client.EmailValidation.ListWithOptions(listingOptions)
    count := 0

    for result := range results {

        if result.Error != nil {
            panic(result.Error)
        }

        fmt.Printf("Id: %v, submitted: %v, status: %v, entries: %v\n",
            result.JobOverview.Id,
            result.JobOverview.SubmittedOn,
            result.JobOverview.Status,
            result.JobOverview.NoOfEntries)

        // Limit the iteration to the first 20 items

        count++

        if count > 20 {
            break
        }
    }

    // Output:
    // Id: 7a7987a3-cc86-4ae8-b3d7-ff0088620503, submitted: 2022-06-18 12:56:43.908432 +0000 UTC, status: InProgress, entries: 23
    // Id: 2c4b1d73-a7b3-40e3-a1b8-748dc499d9f7, submitted: 2022-06-18 12:56:15.698191 +0000 UTC, status: Completed, entries: 12
    // Id: b918a5cb-a853-4cb0-a591-7c8ca21978db, submitted: 2022-06-18 12:56:12.981241 +0000 UTC, status: Completed, entries: 126
    // Id: e3d769b7-e033-422a-b1d8-0a088f566f8d, submitted: 2022-06-18 12:56:02.613184 +0000 UTC, status: Completed, entries: 1
    // Id: 0001aef1-e94f-40c4-b45e-9999f7e37de4, submitted: 2022-06-18 12:56:01.602428 +0000 UTC, status: Completed, entries: 65
    // Id: 0e6c0b38-f3ce-4847-bee8-95947f772242, submitted: 2022-06-18 12:56:01.019199 +0000 UTC, status: Completed, entries: 18
    // Id: 7fedfcb8-4be8-449f-99f4-7ae09e5e8cb5, submitted: 2022-06-18 12:55:54.18652 +0000 UTC, status: Completed, entries: 1
    // ...
}
```

## Managing credits ##

To manage the Verifalia credits for your account you can use the `client.Credit` field.

### Getting the credits balance ###

One of the most common tasks you may need to perform on your account is retrieving the available number of free daily credits and credit packs.
To do that, call the `GetBalance()` function, which returns a `credit.Balance` object, as shown in the next example:

```go
package main

import (
    "fmt"
    "github.com/verifalia/verifalia-go-sdk/verifalia"
)

func main() {
    client := verifalia.NewClient("<USERNAME>", "<PASSWORD>") // See above

    //  Get the available credits balance

    balance, err := client.Credit.GetBalance()

    if err != nil {
        panic(err)
    }

    fmt.Printf("Credit packs: %v, free daily credits: %v (will reset in %v)\n",
        balance.CreditPacks,
        balance.FreeCredits,
        balance.FreeCreditsResetIn)

    // Output:
    // Credit packs: 956.332, free daily credits: 128.66 (will reset in 9h8m23s)
}
```

To add credit packs to your Verifalia account visit [https://verifalia.com/client-area#/credits/add][5].

## Changelog / What's new

### v1.0

Released on June 18<sup>th</sup>, 2022

- First public version of the library, with partial support for the Verifalia API v2.3

[0]: https://verifalia.com
[2]: https://verifalia.com/developers#authentication
[3]: https://verifalia.com/client-area#/users/new
[4]: https://verifalia.com/sign-up
[5]: https://verifalia.com/client-area#/credits/add