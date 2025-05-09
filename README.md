[![CircleCI](https://circleci.com/gh/alexfalkowski/sashactl.svg?style=svg)](https://circleci.com/gh/alexfalkowski/sashactl)
[![codecov](https://codecov.io/gh/alexfalkowski/sashactl/graph/badge.svg?token=QSRFU8VNST)](https://codecov.io/gh/alexfalkowski/sashactl)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/sashactl)](https://goreportcard.com/report/github.com/alexfalkowski/sashactl)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/sashactl.svg)](https://pkg.go.dev/github.com/alexfalkowski/sashactl)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# Sasha's Adventures Client

This client publishes articles that are visible at https://sasha.sasha-adventures.com/ using https://github.com/alexfalkowski/sasha.

## Background

We wanted a way to record all the awesome adventures we were going to. Then we had a baby so less adventures!

### Why a client?

So my wife learns to use the command line.

## Client

The client is broken down into multiple commands.

## Configuration

This client uses S3 (R2) so we need to configure it:

```yaml
aws:
  accessKeyID: secrets/access_key_id
  accessKeySecret: secrets/access_key_secret
  address: R2 url
  region: auto
  retry:
    attempts: 3
    backoff: 100ms
    timeout: 10s
  timeout: 5s
```

We also need to know where we will store the articles locally:

```yaml
articles:
  path: path to articles
```

### Format

The article format is a [YAML](https://yaml.org/) file, with the following format:

```yaml
name: This is a great article
body: Add your story here
slug: this-is-a-great-article
images:
  - name: filename.jpeg
    description: Description of the image.
```

> [!CAUTION]
>  The name/slug can't be changed as it is treaded as an ID.

### Images

The system only supports [JPEG](https://en.wikipedia.org/wiki/JPEG) files.

If you are using photos from iPhone, this will usually use [HEIC](https://en.wikipedia.org/wiki/High_Efficiency_Image_File_Format) files.

So to covert them with [ImageMagick](https://github.com/ImageMagick/ImageMagick), use the following:

```bash
magick mogrify -format jpeg *.heic
```

### New

To create a new article, proceed with the following:

```bash
./sashactl new -n "This is a great article"
```

Then go edit the article (body, images, etc.)

### Publish

To publish a new article, proceed with the following:

```bash
./sashactl publish -s this-is-a-great-article
```

> [!CAUTION]
>  This command takes a slug not a name!

### Unpublish

To unpublish an article, proceed with the following:

```bash
./sashactl unpublish -s this-is-a-great-article
```

> [!CAUTION]
>  This command takes a slug not a name!

## Development

If you would like to contribute, here is how you can get started.

### Structure

The project follows the structure in [golang-standards/project-layout](https://github.com/golang-standards/project-layout).

### Dependencies

Please make sure that you have the following installed:
- [Ruby](https://www.ruby-lang.org/en/)
- [Golang](https://go.dev/)

### Style

This project favours the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

### Setup

Check out [CI](.circleci/config.yml).

### Changes

To see what has changed, please have a look at `CHANGELOG.md`
