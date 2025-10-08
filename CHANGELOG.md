# go-cas Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Changed
- this release now officially concurs to the Golang module versioning scheme
  - please import this library as `github.com/cloudogu/go-cas/v2`

### Fixed
- Allow custom override of log-out detection 
  - this makes sense when deciding POST requests to any URL is not enough
  - see the README.md for further information

## [v2.2.2] - 2019-05-21
### Fixed

- forward requests without basic auth. This can be useful for applications with anonymous users.
- caching of unauthenticated requests
- enhanced error logging
- unit-tests for the rest_handler

## [v2.2.1] - 2018-05-16
### Fixed

- Fix mistake regarding the cache

## [v2.2.0] - 2018-05-16
### Added
- Add REST request cache

## [v2.1.0] 
This and all previous versions origin from the forked repository.
