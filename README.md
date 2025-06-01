# cloudtrail-cli

[![GitHub Actions](https://github.com/guessi/cloudtrail-cli/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/guessi/cloudtrail-cli/actions/workflows/go.yml)
[![GoDoc](https://godoc.org/github.com/guessi/cloudtrail-cli?status.svg)](https://godoc.org/github.com/guessi/cloudtrail-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/guessi/cloudtrail-cli)](https://goreportcard.com/report/github.com/guessi/cloudtrail-cli)
[![GitHub release](https://img.shields.io/github/release/guessi/cloudtrail-cli.svg)](https://github.com/guessi/cloudtrail-cli/releases/latest)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/guessi/cloudtrail-cli)](https://github.com/guessi/cloudtrail-cli/blob/main/go.mod)

Blazing fast single purpose cli for CloudTrail log filtering, written in golang

## 🔢 Prerequisites

* An IAM Role/User with [cloudtrail:LookupEvents](https://docs.aws.amazon.com/awscloudtrail/latest/APIReference/API_LookupEvents.html) permission.

## 🚀 Quick start

```bash
$ cloudtrail-cli --help
```

```bash
$ cloudtrail-cli --start-time 2025-05-12T00:00:00Z --end-time 2025-05-12T01:00:00Z --event-source sts.amazonaws.com --max-results 3
+--------------------------------------+-------------------+----------------------+--------------------------------+-------------------+-------------------+-------------------+----------------------+-----------+----------+
| EventId                              | EventName         | EventTime            | Username                       | EventSource       | UserAgent         | SourceIPAddress   | AccessKeyId          | ErrorCode | ReadOnly |
+--------------------------------------+-------------------+----------------------+--------------------------------+-------------------+-------------------+-------------------+----------------------+-----------+----------+
| 9a7304bb-fc9c-40ce-b148-25b875d5e534 | GetCallerIdentity | 2025-05-12T00:59:57Z | aws-go-sdk-1746934587741269082 | sts.amazonaws.com | eks.amazonaws.com | eks.amazonaws.com | ASIAEXAMPLE098765432 |           | true     |
| d0db6d59-3277-4297-8f73-72eb00c35c77 | GetCallerIdentity | 2025-05-12T00:59:52Z | aws-go-sdk-1746830061119273752 | sts.amazonaws.com | eks.amazonaws.com | eks.amazonaws.com | ASIAEXAMPLE098765432 |           | true     |
| ae8b7cb1-9b58-4897-be37-8f35ff077a99 | GetCallerIdentity | 2025-05-12T00:59:28Z | aws-go-sdk-1746830061119273752 | sts.amazonaws.com | eks.amazonaws.com | eks.amazonaws.com | ASIAEXAMPLE098765432 |           | true     |
+--------------------------------------+-------------------+----------------------+--------------------------------+-------------------+-------------------+-------------------+----------------------+-----------+----------+
```

## :accessibility: FAQ

Why it would return unexpected results when multiple flags are set?

* [cloudtrail-cli](https://github.com/guessi/cloudtrail-cli) leverage [LookupEvents](https://docs.aws.amazon.com/awscloudtrail/latest/APIReference/API_LookupEvents.html) to retrieve events. Howerver, despite there is a `s` in the end of the API name and it does accept a list of `LookupAttributes`, but it doesn't change the limitation that stated in the API document - [Currently the list can contain only one item](https://docs.aws.amazon.com/awscloudtrail/latest/APIReference/API_LookupEvents.html#awscloudtrail-LookupEvents-request-LookupAttributes). Make sure to pass exactly one filter at a time to guarantee your result is expected.


## 👷 Install

### For macOS/Linux users (Recommended)

Brand new install

```bash
brew tap guessi/tap && brew update && brew install cloudtrail-cli
```

To upgrade version

```bash
brew update && brew upgrade cloudtrail-cli
```

### Manually setup (Linux, Windows, macOS)

<details><!-- markdownlint-disable-line -->
<summary>Click to expand!</summary><!-- markdownlint-disable-line -->

#### For Linux users

```bash
$ curl -fsSL https://github.com/guessi/cloudtrail-cli/releases/latest/download/cloudtrail-cli-Linux-$(uname -m).tar.gz -o - | tar zxvf -
$ mv ./cloudtrail-cli /usr/local/bin/cloudtrail-cli
```

#### For macOS users

```bash
$ curl -fsSL https://github.com/guessi/cloudtrail-cli/releases/latest/download/cloudtrail-cli-Darwin-$(uname -m).tar.gz -o - | tar zxvf -
$ mv ./cloudtrail-cli /usr/local/bin/cloudtrail-cli
```

#### For Windows users

```powershell
PS> $SRC = 'https://github.com/guessi/cloudtrail-cli/releases/latest/download/cloudtrail-cli-Windows-x86_64.tar.gz'
PS> $DST = 'C:\Temp\cloudtrail-cli-Windows-x86_64.tar.gz'
PS> Invoke-RestMethod -Uri $SRC -OutFile $DST
```
</details>

## ⚖️ License

[Apache-2.0](LICENSE)
