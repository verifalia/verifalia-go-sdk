package emailValidation

import (
	"github.com/ericlagergren/decimal"
	"github.com/verifalia/verifalia-go-sdk/verifalia/rest"
	"net"
	"time"
)

type Client struct {
	RestClient rest.Client
}

// JobStatus provides enumerated-like values for the supported statuses of an email validation job.
var JobStatus = struct {
	// Unknown status, due to a value reported by the API which is missing in this SDK.
	Unknown string

	// The email validation job is being processed by Verifalia.
	InProgress string

	// The email validation job has been completed and its results are available.
	Completed string

	// The email validation job has either been deleted.
	Deleted string

	// The email validation job is expired.
	Expired string
}{
	Unknown:    "Unknown",
	InProgress: "InProgress",
	Completed:  "Completed",
	Deleted:    "Deleted",
	Expired:    "Expired",
}

// Status provides enumerated-like values for the supported validation statuses of an email validation entry.
var Status = struct {
	// Unknown validation status, due to a value reported by the API which is missing in this SDK.
	Unknown string

	// The e-mail address has been successfully validated.
	Success string

	// A quoted pair within a quoted word is not closed properly.
	UnmatchedQuotedPair string

	// An unexpected quoted pair sequence has been found within a quoted word.
	UnexpectedQuotedPairSequence string

	// A new word boundary start has been detected at an invalid position.
	InvalidWordBoundaryStart string

	// An invalid character has been detected in the provided sequence.
	InvalidCharacterInSequence string

	// The number of parenthesis used to open comments is not equal to the one used to close them.
	UnbalancedCommentParenthesis string

	// An invalid sequence of two adjacent dots has been found.
	DoubleDotSequence string

	// The local part of the e-mail address has an invalid length.
	InvalidLocalPartLength string

	// An invalid folding white space (FWS) sequence has been found.
	InvalidFoldingWhiteSpaceSequence string

	// The at sign symbol (@), used to separate the local part from the domain part of the e-mail address, has not been found.
	AtSignNotFound string

	// An invalid quoted word with no content has been found.
	InvalidEmptyQuotedWord string

	// The e-mail address has an invalid total length.
	InvalidAddressLength string

	// The domain part of the e-mail address is not compliant with the IETF standards.
	DomainPartCompliancyFailure string

	// The e-mail address is not compliant with the additional syntax rules of the e-mail service provider
	// which should eventually manage it.
	IspSpecificSyntaxFailure string

	// The local part of the e-mail address is a well-known role account.
	LocalPartIsWellKnownRoleAccount string

	// A timeout has occurred while querying the DNS server(s) for records about the e-mail address domain.
	DnsQueryTimeout string

	// Verification failed because of a socket connection error occurred while querying the DNS server.
	DnsConnectionFailure string

	// The domain of the e-mail address does not exist.
	DomainDoesNotExist string

	// The domain of the e-mail address does not have any valid DNS record and couldn't accept messages from another
	// host on the Internet.
	DomainIsMisconfigured string

	// The domain has a NULL MX (RFC 7505) resource record and can't thus accept e-mail messages.
	DomainHasNullMx string

	// The e-mail address is provided by a well-known disposable e-mail address provider (DEA).
	DomainIsWellKnownDea string

	// The mail exchanger being tested is a well-known disposable e-mail address provider (DEA).
	MailExchangerIsWellKnownDea string

	// While both the domain and the mail exchanger for the e-mail address being tested are not from a well-known
	// disposable e-mail address provider (DEA), the mailbox is actually disposable.
	MailboxIsDea string

	// A timeout has occurred while connecting to the mail exchanger which serves the e-mail address domain.
	SmtpConnectionTimeout string

	// A socket connection error occurred while connecting to the mail exchanger which serves the e-mail address domain.
	SmtpConnectionFailure string

	// The mailbox for the e-mail address does not exist.
	MailboxDoesNotExist string

	// A connection error occurred while validating the mailbox for the e-mail address.
	MailboxConnectionFailure string

	// The external mail exchanger rejected the validation request.
	LocalSenderAddressRejected string

	// A timeout occurred while verifying the existence of the mailbox.
	MailboxValidationTimeout string

	// The requested mailbox is temporarily unavailable; it could be experiencing technical issues or some other transient problem.
	MailboxTemporarilyUnavailable string

	// The external mail exchanger does not support international mailbox names. To support this feature, mail exchangers must comply with
	// RFC 5336 and support and announce both the 8BITMIME and the UTF8SMTP protocol extensions.
	ServerDoesNotSupportInternationalMailboxes string

	// The requested mailbox is currently over quota.
	MailboxHasInsufficientStorage string

	// A timeout occurred while verifying fake e-mail address rejection for the mail server.
	CatchAllValidationTimeout string

	// The external mail exchanger accepts fake, non-existent, e-mail addresses; therefore the provided emailValidation address MAY be nonexistent too.
	ServerIsCatchAll string

	// A connection error occurred while verifying the external mail exchanger rejects nonexistent e-mail addresses.
	CatchAllConnectionFailure string

	// The mail exchanger responsible for the e-mail address under test is temporarily unavailable.
	ServerTemporaryUnavailable string

	// The mail exchanger responsible for the e-mail address under test replied one or more non-standard SMTP replies which
	// caused the SMTP session to be aborted.
	SmtpDialogError string

	// The external mail exchanger responsible for the e-mail address under test rejected the local endpoint, probably because
	// of its own policy rules.
	LocalEndPointRejected string

	// One or more unhandled exceptions have been thrown during the verification process and something went wrong
	// on the Verifalia side.
	UnhandledException string

	// The mail exchanger responsible for the e-mail address under test hides a honeypot / spam trap.
	MailExchangerIsHoneypot string

	// The domain literal of the e-mail address couldn't accept messages from the Internet.
	UnacceptableDomainLiteral string

	// The item is a duplicate of another e-mail address in the list.
	// To find out the entry this item is a duplicate of, check the DuplicateOf property for the Entry
	// instance which exposes this status code.
	Duplicate string
}{
	Unknown:                                    "Unknown",
	Success:                                    "Success",
	UnmatchedQuotedPair:                        "UnmatchedQuotedPair",
	UnexpectedQuotedPairSequence:               "UnexpectedQuotedPairSequence",
	InvalidWordBoundaryStart:                   "InvalidWordBoundaryStart",
	InvalidCharacterInSequence:                 "InvalidCharacterInSequence",
	UnbalancedCommentParenthesis:               "UnbalancedCommentParenthesis",
	DoubleDotSequence:                          "DoubleDotSequence",
	InvalidLocalPartLength:                     "InvalidLocalPartLength",
	InvalidFoldingWhiteSpaceSequence:           "InvalidFoldingWhiteSpaceSequence",
	AtSignNotFound:                             "AtSignNotFound",
	InvalidEmptyQuotedWord:                     "InvalidEmptyQuotedWord",
	InvalidAddressLength:                       "InvalidAddressLength",
	DomainPartCompliancyFailure:                "DomainPartCompliancyFailure",
	IspSpecificSyntaxFailure:                   "IspSpecificSyntaxFailure",
	LocalPartIsWellKnownRoleAccount:            "LocalPartIsWellKnownRoleAccount",
	DnsQueryTimeout:                            "DnsQueryTimeout",
	DnsConnectionFailure:                       "DnsConnectionFailure",
	DomainDoesNotExist:                         "DomainDoesNotExist",
	DomainIsMisconfigured:                      "DomainIsMisconfigured",
	DomainHasNullMx:                            "DomainHasNullMx",
	DomainIsWellKnownDea:                       "DomainIsWellKnownDea",
	MailExchangerIsWellKnownDea:                "MailExchangerIsWellKnownDea",
	MailboxIsDea:                               "MailboxIsDea",
	SmtpConnectionTimeout:                      "SmtpConnectionTimeout",
	SmtpConnectionFailure:                      "SmtpConnectionFailure",
	MailboxDoesNotExist:                        "MailboxDoesNotExist",
	MailboxConnectionFailure:                   "MailboxConnectionFailure",
	LocalSenderAddressRejected:                 "LocalSenderAddressRejected",
	MailboxValidationTimeout:                   "MailboxValidationTimeout",
	MailboxTemporarilyUnavailable:              "MailboxTemporarilyUnavailable",
	ServerDoesNotSupportInternationalMailboxes: "ServerDoesNotSupportInternationalMailboxes",
	MailboxHasInsufficientStorage:              "MailboxHasInsufficientStorage",
	CatchAllValidationTimeout:                  "CatchAllValidationTimeout",
	ServerIsCatchAll:                           "ServerIsCatchAll",
	CatchAllConnectionFailure:                  "CatchAllConnectionFailure",
	ServerTemporaryUnavailable:                 "ServerTemporaryUnavailable",
	SmtpDialogError:                            "SmtpDialogError",
	LocalEndPointRejected:                      "LocalEndPointRejected",
	UnhandledException:                         "UnhandledException",
	MailExchangerIsHoneypot:                    "MailExchangerIsHoneypot",
	UnacceptableDomainLiteral:                  "UnacceptableDomainLiteral",
	Duplicate:                                  "Duplicate",
}

// Classification provides enumerated-like values for the classifications of the supported validation statuses of an
// email validation entry.
var Classification = struct {
	// Refers to an e-mail address which is deliverable.
	Deliverable string

	// Refers to an e-mail address which could be no longer valid.
	Risky string

	// Refers to an e-mail address which is either invalid or no longer deliverable.
	Undeliverable string

	// Contains an e-mail address whose deliverability is unknown.
	Unknown string
}{
	Deliverable:   "Deliverable",
	Risky:         "Risky",
	Undeliverable: "Undeliverable",
	Unknown:       "Unknown",
}

// Overview contains information about an e-mail validation job.
type Overview struct {
	// The unique identifier for the validation job.
	Id string

	// The date and time this validation job has been created in Verifalia.
	CreatedOn time.Time

	// The date and time this validation job has been submitted to Verifalia.
	SubmittedOn time.Time

	// The date and time this validation job has been eventually completed.
	CompletedOn *time.Time

	// The eventual priority (speed) of the validation job, relative to the parent Verifalia account. In the event of an account
	// with many concurrent validation jobs, this value allows to increase the processing speed of a job with respect to the others.
	Priority *uint8

	// An optional user-defined name for the validation job, for your own reference.
	Name string

	// The unique ID of the Verifalia user who submitted the validation job.
	Owner string

	// The IP address of the client which submitted the validation job.
	ClientIP net.IP

	// A reference to the quality level this job was validated against. The Quality enum-like object contains
	// the supported values, for example: Quality.High
	Quality string

	// The maximum data retention period Verifalia observes for this verification job, after which the job will be
	// automatically deleted.
	// A verification job can be deleted anytime prior to its retention period through the Delete() function.
	Retention time.Duration

	// The deduplication mode for the validation job. The Deduplication enum-like object contains
	// the supported values, for example: Deduplication.Relaxed
	Deduplication string

	// The processing status for the validation job. The JobStatus enum-like object contains
	// the supported values, for example: JobStatus.InProgress
	Status string

	// The number of entries the validation job contains.
	NoOfEntries uint

	// The eventual completion progress for the validation job.
	Progress *Progress
}

// Entry represents a single validated entry within an email verification job.
type Entry struct {
	// The index of this entry within its containing job. This property is mostly useful in the event
	// the API returns a filtered view of the items.
	Index int `json:"index"`

	// The input string being validated.
	InputData string `json:"inputData"`

	// A custom, optional string which is passed back upon completing the validation. To pass back and forth
	// a custom value, use the Custom field of the ValidationRequestEntry struct
	Custom string `json:"custom"`

	// The date this entry has been completed, if available.
	CompletedOn *time.Time `json:"completedOn"`

	// Gets the normalized email address, without any eventual comment or folding white space.
	EmailAddress string `json:"emailAddress"`

	// Gets the domain part of the email address, converted to ASCII if needed and with comments and folding
	// white spaces stripped off.
	AsciiEmailAddressDomainPart string `json:"asciiEmailAddressDomainPart"`

	// Gets the local part of the email address, without comments and folding white spaces.
	EmailAddressLocalPart string `json:"emailAddressLocalPart"`

	// Gets the domain part of the email address, without comments and folding white spaces.
	EmailAddressDomainPart string `json:"emailAddressDomainPart"`

	// If true, the email address has an international domain name.
	HasInternationalDomainName *bool `json:"hasInternationalDomainName"`

	// If true, the email address has an international mailbox name.
	HasInternationalMailboxName *bool `json:"hasInternationalMailboxName"`

	// If true, the email address comes from a disposable email address (DEA) provider.
	// See https://verifalia.com/help/email-validations/what-is-a-disposable-email-address-dea for additional information
	// about disposable email addresses.
	IsDisposableEmailAddress *bool `json:"isDisposableEmailAddress"`

	// If true, the email address comes from a free email address provider (e.g. gmail, yahoo, outlook / hotmail, ...).
	IsFreeEmailAddress *bool `json:"isFreeEmailAddress"`

	// If true, the local part of the email address is a well-known role account.
	IsRoleAccount *bool `json:"isRoleAccount"`

	// The validation status for this entry. The Status enum-like object contains
	// the supported values, for example: Status.MailboxHasInsufficientStorage
	Status string `json:"status"`

	// The classification for the status of this email address. The Classification enum-like object contains
	//	// the supported values, for example: Classification.Deliverable
	Classification string `json:"classification"`

	// The zero-based position of the character in the input data string that eventually caused the syntax validation to fail.
	SyntaxFailureIndex *int `json:"syntaxFailureIndex"`

	// The zero-based index of the first occurrence of this email address in the parent validation job, in the event the
	// Status for this entry is equal to Status.Duplicate; duplicated items do not expose any result detail apart from this and the
	// eventual Custom values.
	DuplicateOf *int `json:"duplicateOf"`
}

// Job represents a snapshot of an e-mail validation job, along with its overview and eventual validated entries.
type Job struct {
	// Overview information for this e-mail validation job.
	Overview Overview

	// The eventual validated items for this e-mail validation job.
	Entries []Entry
}

// Progress contain the progress details for an e-mail validation job.
type Progress struct {
	// The percentage of completed entries, ranging from 0 to 1.
	Percentage decimal.Big

	// An eventual estimated required time span needed to complete the whole job.
	EstimatedTimeRemaining *time.Duration
}

// Priority includes useful values for the speed of a validation job, relative to the parent Verifalia account. In the event of an account
// with many concurrent validation jobs, this value allows to increase the processing speed of a job with respect to the others.
var Priority = struct {
	// The lowest possible processing priority (speed) for a validation job.
	Lowest uint8

	// Normal processing priority (speed) for a validation job.
	Normal uint8

	// The highest possible processing priority (speed) for a validation job.
	Highest uint8
}{
	Lowest:  0,
	Normal:  127,
	Highest: 255,
}

// Quality is a reference to a Verifalia quality level. Quality levels determine how Verifalia validates email addresses, including whether
// and how the automatic reprocessing logic occurs (for transient statuses) and the verification timeouts settings.
var Quality = struct {
	// The Standard quality level. Suitable for most businesses, provides good results for the vast majority of email addresses;
	// features a single validation pass and 5 second anti-tarpit time; less suitable for validating email addresses with temporary
	// issues (mailbox over quota, greylisting, etc.) and slower mail exchangers.
	Standard string

	// The High quality level. Much higher quality, featuring 3 validation passes and 50 seconds of anti-tarpit time, so you can
	// even validate most addresses with temporary issues, or slower mail exchangers.
	High string

	// The Extreme quality level. Unbeatable, top-notch quality for professionals who need the best results the industry can offer:
	// performs email validations at the highest level, with 9 validation passes and 2 minutes of anti-tarpit time.
	Extreme string
}{
	Standard: "Standard",
	High:     "High",
	Extreme:  "Extreme",
}

// Deduplication contains the supported strategies Verifalia follows while determining which email addresses are duplicates,
// within an e-mail verification job with multiple items.
// Duplicated items (after the first occurrence) will have the Status.Duplicate status.
var Deduplication = struct {
	// Duplicates detection is turned off.
	Off string

	// Identifies duplicates using an algorithm with safe rules-only, which guarantee no false duplicates.
	Safe string

	// Identifies duplicates using a set of relaxed rules which assume the target email service providers
	// are configured with modern settings only (instead of the broader options the RFCs from the '80s allow).
	Relaxed string
}{
	Off:     "Off",
	Safe:    "Safe",
	Relaxed: "Relaxed",
}
