# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Terraform provider for [MyHeritage](https://www.myheritage.com) genealogy data, built on the
[terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework). It exposes a
single managed resource, `myheritage_profile`, backed by MyHeritage's private GraphQL API. The
provider address is `github.com/dmalch/myheritage` and its type name is `myheritage`.

## Commands

- `make build` — compiles to `bin/terraform-provider-myheritage`
- `make test` — `go test -v ./...`
- `make clean` — removes `bin/`
- Single test: `go test -v ./internal/myheritage -run TestGetIndividual`

Go version is pinned in `go.mod` (currently 1.25.0). CI (`.github/workflows/ci.yaml`) reads it via
`go-version-file: go.mod` — never hardcode a `go-version:` in the workflow, or a version bump in
`go.mod` will silently mismatch CI and break the build.

## Architecture

Two layers, cleanly separated:

- **`internal/` (provider layer)** — terraform-plugin-framework glue.
  - `internal/provider.go` — the provider itself: declares the `api_key` (required, sensitive)
    attribute, and registers `profile.NewProfileResource`. `Configure` stashes the config as
    `*config.MyHeritageProviderConfig` in `ResourceData`.
  - `internal/resource/profile/` — the `myheritage_profile` resource. `resource.go` holds the CRUD
    methods and the `ResourceModel`/`EventModel`/`NoteModel`/`MediaModel` structs plus their
    `types.ObjectType` definitions; `schema.go` holds the matching Terraform schema. **When you add
    or change a field you must edit it in three places that must stay in sync:** the `*Model` struct
    (`tfsdk` tags), the `...ObjectType()` `AttrTypes` map, and the schema in `schema.go`.
  - `internal/config/provider_config.go` — the provider config struct (`ApiKey`).

- **`internal/myheritage/` (API client layer)** — plain Go, no Terraform types. Each exported
  function (`GetProfileHeader`, `GetProfileDetails`, `GetIndividual`, `GetIndividualBiography`,
  `CreateProfile`, `UpdateProfile`, `DeleteProfile`) takes the `apiKey` as its first argument, builds
  a request, and POSTs to the GraphQL endpoint in `const.go` (`https://familygraphql.myheritage.com/`).
  All requests go through `doRequest` in `http_client.go`, which adds retry-on-429 (3 attempts, fixed
  2s delay via `retry-go`) and treats any non-200 as an error.

### Things worth knowing before you edit

- **Read is the only fully-implemented CRUD op.** `CreateProfile`, `UpdateProfile`, and
  `DeleteProfile` in `profile.go` are stubs that return empty/nil — so the resource's `Create`,
  `Update`, and `Delete` currently do nothing against the API. Assume the read path is real and the
  write path is a work-in-progress unless the stubs have been filled in.
- **GraphQL queries are hardcoded inline strings**, copied verbatim from MyHeritage's web client
  (fragments and all). Most requests use a JSON body (`GraphqlRequest`); `GetIndividual` uses
  `multipart/form-data` instead. Don't try to "clean up" these query strings — they mirror the
  server's expected shape.
- **Event type codes come from MyHeritage**, not us: `EventTypeBirth = "BITH"`, `EventTypeDeath =
  "DEAT"` (the `BITH` spelling is theirs). `getContent` deliberately returns null content for `DEAT`
  events.
- **`GetProfileDetails` filters out `is_fact_of_relative` events** in the resource `Read` loop —
  only the profile's own facts land in state.

## Testing

Tests use [gomega](https://github.com/onsi/gomega) (`. "github.com/onsi/gomega"` dot-import,
`RegisterTestingT(t)`). The API-client tests are **integration tests that hit the live MyHeritage
API** and are therefore `t.Skip()`'d by default; they need a real bearer token in `testApiKey`
(`const_test.go`, committed empty) and real profile/individual IDs. `make test` passes with them
skipped. To run one, drop the `t.Skip()` and set `testApiKey` locally — do not commit a real key.
