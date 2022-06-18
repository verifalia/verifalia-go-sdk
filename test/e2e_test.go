package main

import (
	"bytes"
	"github.com/verifalia/verifalia-go-sdk/verifalia"
	"github.com/verifalia/verifalia-go-sdk/verifalia/common"
	"github.com/verifalia/verifalia-go-sdk/verifalia/emailValidation"
	"github.com/verifalia/verifalia-go-sdk/verifalia/rest"
	"os"
	"testing"
	"time"
)

func buildClient() *verifalia.Client {
	return verifalia.NewClient(os.Getenv("VERIFALIA_USERNAME"), os.Getenv("VERIFALIA_PASSWORD"))
}

func TestFileSubmission(t *testing.T) {
	client := buildClient()

	fileData := []byte("test@example.com\njohn.petrucci@gmail.com\nhellospank@yahoo.com")
	reader := bytes.NewReader(fileData)

	fileValidation, err := client.EmailValidation.SubmitFileReaderWithOptions(reader, &emailValidation.FileSubmissionOptions{
		ContentType: rest.ContentType.TextPlain,
	}, nil)

	if err != nil {
		t.Error(err)
	}

	t.Logf("%v", *fileValidation)

	fileValidation, err = client.EmailValidation.WaitForCompletion(fileValidation)

	if err != nil {
		t.Error(err)
	}

	t.Logf("%v", *fileValidation)

	for _, entry := range fileValidation.Entries {
		t.Logf("%v => %v (%v)\n", entry.InputData, entry.Status, entry.Classification)
	}
}

func TestGetBalance(t *testing.T) {
	client := buildClient()
	balance, err := client.Credit.GetBalance()

	if err != nil {
		t.Error(err)
	}

	t.Logf("Credit packs: %s\n", balance.CreditPacks.String())
	t.Logf("Free credits: %s\n", balance.FreeCredits.String())
	t.Logf("Free credits reset in: %s\n", *balance.FreeCreditsResetIn)
}

func TestSingleValidation(t *testing.T) {
	client := buildClient()
	singleValidation, err := client.EmailValidation.Submit("hellokitty@gmail.com")

	if err != nil {
		t.Error(err)
	}

	singleValidation, err = client.EmailValidation.WaitForCompletion(singleValidation)

	if err != nil {
		t.Error(err)
	}

	for _, entry := range singleValidation.Entries {
		t.Logf("%v => %v (%v)\n", entry.InputData, entry.Status, entry.Classification)
	}
}

func TestSubmitMany(t *testing.T) {
	client := buildClient()

	validation, err := client.EmailValidation.SubmitManyWithOptions([]emailValidation.ValidationRequestEntry{
		{
			InputData: "hellokitty@gmail.com",
		},
		{
			InputData: "invalid@domain@com",
			Custom:    "test123",
		},
		{
			InputData: "invalid@domain@com",
		},
	}, &emailValidation.SubmissionOptions{
		Quality:       emailValidation.Quality.High,
		Retention:     30 * time.Minute,
		Deduplication: emailValidation.Deduplication.Relaxed,
	})

	if err != nil {
		t.Error(err)
	}

	validation, err = client.EmailValidation.WaitForCompletion(validation)

	if err != nil {
		t.Error(err)
	}

	for _, entry := range validation.Entries {
		t.Logf("%v => %v (%v)\n", entry.InputData, entry.Status, entry.Classification)
	}
}

func TestListing(t *testing.T) {
	client := buildClient()

	results := client.EmailValidation.ListWithOptions(emailValidation.ListingOptions{
		Direction: common.Backward,
	})

	count := 0

	for result := range results {
		t.Logf("Result...")

		if result.Error != nil {
			t.Error(result.Error)
		}

		t.Logf("%v => %v\n", result.JobOverview.Id, result.JobOverview.SubmittedOn)

		// Limit the iteration to the first 100 items

		count++

		if count > 100 {
			break
		}
	}
}
