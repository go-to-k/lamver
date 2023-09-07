# lamver

[![Go Report Card](https://goreportcard.com/badge/github.com/go-to-k/lamver)](https://goreportcard.com/report/github.com/go-to-k/lamver) ![GitHub](https://img.shields.io/github/license/go-to-k/lamver) ![GitHub](https://img.shields.io/github/v/release/go-to-k/lamver) [![ci](https://github.com/go-to-k/lamver/actions/workflows/ci.yml/badge.svg)](https://github.com/go-to-k/lamver/actions/workflows/ci.yml)

The description in **Japanese** is available on the following blog page. -> [Blog](https://go-to-k.hatenablog.com/entry/lamver)

The description in **English** is available on the following blog page. -> [Blog](https://dev.to/aws-builders/lambda-runtimeversion-search-tool-across-regions-41l0)

## What is

CLI tool to search AWS Lambda runtime values and versions.

By filtering by **regions, runtime and versions, or part of the function name**, you can find a list of functions **across regions**.

This will allow you can see the following.

- What Lambda functions exist in which regions
- Whether there are any functions that **have reached EOL**
- Whether there is a function in **an unexpected region**
- Whether a function exists based on **a specific naming rule**

Also this tool can support output results **as a CSV file.**

## Install

- Homebrew
  ```sh
  brew install go-to-k/tap/lamver
  ```
- Binary
  - [Releases](https://github.com/go-to-k/lamver/releases)
- Git Clone and install(for developers)
  ```sh
  git clone https://github.com/go-to-k/lamver.git
  cd lamver
  make install
  ```

## How to use
  ```sh
  lamver [-p <profile>] [-r <default region>] [-o <output file path>] [-k <keyword for function name>]
  ```

### options

- -p, --profile: optional
  - AWS profile name
- -r, --region: optional
  - Default AWS region
    - The region to output is selected interactively and does not need to be specified.
- -o, --output: optional
  - Output file path for CSV format
- -k, --keyword: optional
  - Keyword for function name filtering (case-insensitive)

## Input flow

### Enter `lamver`

```sh
❯ lamver
```

You can specify `-k, --keyword` option. This is a keyword for **function name filtering (case-insensitive)**.

```sh
❯ lamver -k goto
```

### Choose regions

```sh
? Select regions you want to search.
  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]
  [x]  ap-northeast-1
  [ ]  ap-northeast-2
  [ ]  ap-northeast-3
  [ ]  ap-south-1
  [ ]  ap-southeast-1
  [ ]  ap-southeast-2
  [ ]  ca-central-1
  [ ]  eu-central-1
  [ ]  eu-north-1
  [ ]  eu-west-1
  [ ]  eu-west-2
  [ ]  eu-west-3
  [ ]  sa-east-1
  [x]  us-east-1
> [x]  us-east-2
  [ ]  us-west-1
  [ ]  us-west-2
```

### Choose runtime values

```sh
? Select runtime values you want to search.
  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]
> [ ]  dotnet6
  [ ]  dotnetcore1.0
  [ ]  dotnetcore2.0
  [ ]  dotnetcore2.1
  [ ]  dotnetcore3.1
  [x]  go1.x
  [ ]  java8
  [ ]  java8.al2
  [ ]  java11
  [ ]  java17
  [ ]  nodejs
  [ ]  nodejs4.3
  [ ]  nodejs4.3-edge
  [ ]  nodejs6.10
  [ ]  nodejs8.10
  [ ]  nodejs10.x
  [x]  nodejs12.x
  [ ]  nodejs14.x
  [ ]  nodejs16.x
  [ ]  nodejs18.x
  [ ]  provided
  [x]  provided.al2
  [ ]  python2.7
  [ ]  python3.6
  [ ]  python3.7
  [ ]  python3.8
  [ ]  python3.9
  [ ]  python3.10
  [ ]  python3.11
  [ ]  ruby2.5
  [ ]  ruby2.7
  [ ]  ruby3.2
```

### Enter part of the function name

You can search function names in a **case-insensitive**.

**Empty** input will output **all functions**.

This phase is skipped if you specify `-k` option.

```sh
Filter a keyword of function names(case-insensitive): test-goto
```

### The result will be output

```sh
+--------------+----------------+----------------------+------------------------------+
|   RUNTIME    |     REGION     |     FUNCTIONNAME     |         LASTMODIFIED         |
+--------------+----------------+----------------------+------------------------------+
| go1.x        | ap-northeast-1 | Test-goto-function2  | 2023-01-07T14:54:23.406+0000 |
+--------------+----------------+----------------------+------------------------------+
| go1.x        | ap-northeast-1 | test-Goto-function10 | 2023-01-07T15:29:11.658+0000 |
+--------------+----------------+----------------------+------------------------------+
| go1.x        | us-east-2      | TEST-goto-function6  | 2023-01-07T15:28:08.507+0000 |
+--------------+----------------+----------------------+------------------------------+
| nodejs12.x   | ap-northeast-1 | test-GOTO-function1  | 2023-01-07T14:53:49.141+0000 |
+--------------+----------------+----------------------+------------------------------+
| nodejs12.x   | us-east-1      | TEST-GOTO-function4  | 2023-01-07T15:18:14.191+0000 |
+--------------+----------------+----------------------+------------------------------+
| nodejs12.x   | us-east-1      | test-goto-function7  | 2023-01-07T15:28:20.921+0000 |
+--------------+----------------+----------------------+------------------------------+
| nodejs12.x   | us-east-2      | test-goto-function5  | 2023-01-07T15:18:34.408+0000 |
+--------------+----------------+----------------------+------------------------------+
| provided.al2 | ap-northeast-1 | test-goto-function8  | 2023-01-07T15:28:34.968+0000 |
+--------------+----------------+----------------------+------------------------------+
| provided.al2 | us-east-1      | test-goto-function3  | 2023-01-07T15:17:35.965+0000 |
+--------------+----------------+----------------------+------------------------------+
| provided.al2 | us-east-2      | test-goto-function9  | 2023-01-07T15:29:16.107+0000 |
+--------------+----------------+----------------------+------------------------------+
INF 10 counts hit!
```

## CSV output mode

By default, results are output as table format on the screen.

If you add `-o` option, then results can be output **as a CSV file**.

```sh
lamver -o ./result.csv
```
