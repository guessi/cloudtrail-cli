# cloudtrail-cli

[![GitHub Actions](https://github.com/guessi/cloudtrail-cli/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/guessi/cloudtrail-cli/actions/workflows/go.yml)
[![GoDoc](https://godoc.org/github.com/guessi/cloudtrail-cli?status.svg)](https://godoc.org/github.com/guessi/cloudtrail-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/guessi/cloudtrail-cli)](https://goreportcard.com/report/github.com/guessi/cloudtrail-cli)
[![GitHub release](https://img.shields.io/github/release/guessi/cloudtrail-cli.svg)](https://github.com/guessi/cloudtrail-cli/releases/latest)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/guessi/cloudtrail-cli)](https://github.com/guessi/cloudtrail-cli/blob/master/go.mod)

Blazing fast single purpose cli for CloudTrail log filtering, written in golang

## 🚀 Quick start

```bash
$ cloudtrail-cli --help
```

```bash
$ cloudtrail-cli --start-time 2023-02-01T00:00:00 --end-time 2023-02-01T01:00:00 --event-name AssumeRole --max-results 5
+--------------------------------------+------------+----------------------+----------+-------------------+-------------------------------+-------------------------------+-------------+-----------+----------+
| EventId                              | EventName  | EventTime            | Username | EventSource       | UserAgent                     | SourceIPAddress               | AccessKeyId | ErrorCode | ReadOnly |
+--------------------------------------+------------+----------------------+----------+-------------------+-------------------------------+-------------------------------+-------------+-----------+----------+
| 998a47f3-fb53-48e0-83f1-111111111111 | AssumeRole | 2023-02-01T00:58:28Z | -        | sts.amazonaws.com | eks.amazonaws.com             | eks.amazonaws.com             |             |           | true     |
| 56018bd8-d0f4-41d3-a718-111111111111 | AssumeRole | 2023-02-01T00:57:51Z | -        | sts.amazonaws.com | internetmonitor.amazonaws.com | internetmonitor.amazonaws.com |             |           | true     |
| d5f7ff3f-af90-4f05-9050-111111111111 | AssumeRole | 2023-02-01T00:55:22Z | -        | sts.amazonaws.com | ssm.amazonaws.com             | ssm.amazonaws.com             |             |           | true     |
| 139dd66c-d192-47fc-9158-111111111111 | AssumeRole | 2023-02-01T00:40:38Z | -        | sts.amazonaws.com | lambda.amazonaws.com          | lambda.amazonaws.com          |             |           | true     |
| 8af6dc45-fd58-4ad5-9e95-111111111111 | AssumeRole | 2023-02-01T00:35:06Z | -        | sts.amazonaws.com | lambda.amazonaws.com          | lambda.amazonaws.com          |             |           | true     |
+--------------------------------------+------------+----------------------+----------+-------------------+-------------------------------+-------------------------------+-------------+-----------+----------+
```

## 👷 Install

### For macOS users (Recommended)

```bash
$ brew tap guessi/tap && brew update && brew install cloudtrail-cli
```

### Manually setup (Linux, Windows, macOS)

<details><!-- markdownlint-disable-line -->
<summary>Click to expand!</summary><!-- markdownlint-disable-line -->

### For Linux users

```bash
$ curl -fsSL https://github.com/guessi/cloudtrail-cli/releases/latest/download/cloudtrail-cli-Linux-$(uname -m).tar.gz -o - | tar zxvf -
$ mv ./cloudtrail-cli /usr/local/bin/cloudtrail-cli
```

### For macOS users

```bash
$ curl -fsSL https://github.com/guessi/cloudtrail-cli/releases/latest/download/cloudtrail-cli-Darwin-$(uname -m).tar.gz -o - | tar zxvf -
$ mv ./cloudtrail-cli /usr/local/bin/cloudtrail-cli
```

### For Windows users

```powershell
PS> $SRC = 'https://github.com/guessi/cloudtrail-cli/releases/latest/download/cloudtrail-cli-Windows-x86_64.tar.gz'
PS> $DST = 'C:\Temp\cloudtrail-cli-Windows-x86_64.tar.gz'
PS> Invoke-RestMethod -Uri $SRC -OutFile $DST
```
</details>

## ⚖️ License

[Apache-2.0](LICENSE)
