<!-- markdownlint-disable -->

# v1.0.17 / 2025-02-28

* Update staticcheck@2025.1
* Dependencies update
* Remove `--no-read-only` to make it easier to use
  * no `--read-only` flag set, it should not pass `ReadOnly` to the `LookupAttributes`
  * with `--read-only` flag set (or `--read-only=true`), it would pass `ReadOnly=true` to the `LookupAttributes`
  * with `--read-only=false` flag set, it would pass `ReadOnly=false` to the `LookupAttributes`

# v1.0.16 / 2025-01-04

* Biuld with golang 1.23
* Dependencies update
* Support for filtering with ResourceName
* Support for filtering with ResourceType

# v1.0.15 / 2024-12-16

* Dependencies update

# v1.0.14 / 2024-09-09

* Dependencies update

# v1.0.13 / 2024-08-02

* Dependencies update

# v1.0.12 / 2024-06-24

* Dependencies update

# v1.0.11 / 2024-03-24

* Biuld with golang 1.22
* Dependencies update

# v1.0.10 / 2024-02-16

* Dependencies update
* Implement sub-command "version"
* Bump ghr@v0.16.2
* Bump actions/cache@v4

# v1.0.9 / 2024-01-07

* Dependencies update

# v1.0.8 / 2023-11-20

* Dependencies update

# v1.0.7 / 2023-10-12

* Identical with v1.0.6

# v1.0.6 / 2023-10-12

* Biuld with golang 1.21
* Dependencies update

# v1.0.5 / 2023-08-08

* Identical with v1.0.4

# v1.0.4 / 2023-08-08

* Update staticcheck@2023.1.3
* Dependencies update

# v1.0.3 / 2023-05-12

* Bump actions/setup-go@v4
* Dependencies update

# v1.0.2 / 2023-03-14

* Add support for filtering with AccessKeyId
* Introduce `--truncate-user-name` and `--truncate-user-agent` flags

# v1.0.1 / 2023-03-12

* Introduce Homebrew install
* Dependencies update

# v1.0.0 / 2023-02-23

* Initial Release
