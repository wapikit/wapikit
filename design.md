# Repo management and binaries

## OSS Community Edition

- Open source version should include:

  - basic API server which is in open core repos
  - OSS frontend
  - campaign manager

- Open source binaries will be released publicly under github.com/wapikit/wapikit
- frontend can show a lock button and a button to link to the managed version

```sh
go build -tags=oss -o wapikit_oss ./cmd/main.go
```

## Cloud Edition

- Managed Cloud version build should include:

  - payments_controller
  - premium integrations code
  - subscription_controller
  - database schema

- closed source binaries will be built and deployed to cloud server

- there will be three plans in this, so

  - need to have a middleware which checks the plan the user is on and then allow the feature access on the same basis
  - frontend should show a lock button with a upgrade button for feature unlock

- FREE Plan: - 1k contacts (check in contact creation controller) - Only Up to 2 Campaigns a week (check in campaign creation controller) - Only up to 2 organization (check in org creation controller)

- PRO Plan
  - All Premium integration unlocked
  -

```sh
go build -tags=managed -o wapikit_cloud ./cmd/main.go
```

# Enterprise Edition

- should included all the managed files other than the payments_controller and subscription_controller
- enterprise binaries should be also be released but in private repository and all the tags that pushed to open core repo should be synced with the enterprise repo: github.com/wapikit/wapikit_enterprise

```sh
go build -tags=enterprise -o wapikit_enterprise ./cmd/main.go
```

# NEED OF THE HOUR

- Payment service [Closed Source]
- Subscription service handle [Closed Source]
- Custom website pages [Closed Source]
- Impose Limits as per Edition and Subscription Plan in Cloud Edition
