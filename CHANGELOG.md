# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased
### Added

## [1.2.0] - 2019-07-29
### Added

- Support for agent label assignment and matching, by [@logikone](https://github.com/logikone).
- Allow Hetzner to choose datacenter when none specified, by [@tboerger](https://github.com/tboerger).

### Fixed

- Upgraded zerolog to fix duplicate keys in json output, by [@krtx](https://github.com/krtx).

## [1.1.0] - 2019-05-29
### Added

- Create AWS instances with Name tag set to agent unique id, from [@bradrydzewski](https://github.com/bradrydzewski).
- Handle AWS instance not found errors, from [@andy-trimble](https://github.com/andy-trimble).
- Remove hard-coded DNS servers from the default Docker configuration, from [jones2026](https://github.com/jones2026).

## [1.0.0] - 2019-05-06
### Added

- Optional support for watchtower from [@bradrydzewski](https://github.com/bradrydzewski).
- Optional support for drone/gc from [@bradrydzewski](https://github.com/bradrydzewski). 
- Update the default agent image to 1.0 stable, from [@bradrydzewski](https://github.com/bradrydzewski).
- Configure agent environment variables from [@bradrydzewski](https://github.com/bradrydzewski).
- Configure agent host volume mounts from [@patrickjahns](https://github.com/patrickjahns).
- Update Digital Ocean default image from [@jlesage](https://github.com/jlesage).
- Fix problems using custom Digital Ocean image from [@jlesage](https://github.com/jlesage).
