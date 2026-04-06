# Email Verifier - Go Library

[![Go Reference](https://pkg.go.dev/badge/github.com/EnrowAPI/email-verifier-go.svg)](https://pkg.go.dev/github.com/EnrowAPI/email-verifier-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![GitHub stars](https://img.shields.io/github/stars/EnrowAPI/email-verifier-go)](https://github.com/EnrowAPI/email-verifier-go)
[![Last commit](https://img.shields.io/github/last-commit/EnrowAPI/email-verifier-go)](https://github.com/EnrowAPI/email-verifier-go/commits)

Verify email addresses in real time. Check deliverability, detect disposable and catch-all mailboxes, and clean your lists before sending. Integrate email verification into your sales pipeline, CRM sync, or outreach workflow.

Powered by [Enrow](https://enrow.io) -- each verification costs just 0.25 credits.

## Installation

```bash
go get github.com/EnrowAPI/email-verifier-go
```

Requires Go 1.21+. Zero dependencies.

## Simple Usage

```go
package main

import (
	"fmt"
	"log"

	emailverifier "github.com/EnrowAPI/email-verifier-go"
)

func main() {
	verification, err := emailverifier.VerifyEmail("your_api_key", "tcook@apple.com", "")
	if err != nil {
		log.Fatal(err)
	}

	result, err := emailverifier.GetVerificationResult("your_api_key", verification.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.Email)         // tcook@apple.com
	fmt.Println(result.Qualification) // valid
}
```

`VerifyEmail` returns a verification ID. The check runs asynchronously -- call `GetVerificationResult` to retrieve the result once it is ready. You can also pass a webhook URL to get notified automatically.

## Webhook notifications

Pass a webhook URL as the third argument to receive results automatically when verification completes.

```go
verification, err := emailverifier.VerifyEmail("your_api_key", "tcook@apple.com", "https://example.com/webhook")
```

## Bulk verification

```go
batch, err := emailverifier.VerifyEmails("your_api_key", []string{
	"tcook@apple.com",
	"satya@microsoft.com",
	"jensen@nvidia.com",
}, "")
if err != nil {
	log.Fatal(err)
}

// batch.BatchID, batch.Total, batch.Status

results, err := emailverifier.GetVerificationResults("your_api_key", batch.BatchID)
if err != nil {
	log.Fatal(err)
}
// results.Results -- slice of VerificationResult
```

Up to 5,000 verifications per batch. Pass a webhook URL to get notified when the batch completes.

## Error handling

```go
_, err := emailverifier.VerifyEmail("bad_key", "test@test.com", "")
if err != nil {
	fmt.Println(err)
	// Common errors:
	// - "Invalid or missing API key" (401)
	// - "Your credit balance is insufficient." (402)
	// - "Rate limit exceeded" (429)
}
```

## Getting an API key

Register at [app.enrow.io](https://app.enrow.io) to get your API key. You get **50 free credits** (= 200 verifications) with no credit card required.

Paid plans start at **$17/mo** for 1,000 emails up to **$497/mo** for 100,000 emails. See [pricing](https://enrow.io/pricing).

## Documentation

- [Enrow API documentation](https://docs.enrow.io)
- [Full Enrow SDK](https://github.com/EnrowAPI/enrow-go) -- includes email finder, phone finderand more

## License

MIT -- see [LICENSE](LICENSE) for details.
