# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased
### Fixed
- pending instance count is excluded from determining available capacity when reducing pool size, by [bradrydzewski](https://github.com/bradrydzewski).

## [1.7.2]
### Fixed
- captuare instance private IP when google compute private IP is enabled, by [ademariag](https://github.com/ademariag). [#68](https://github.com/drone/autoscaler/pull/68).

## [1.7.2]
### Added
- support for aws fallback instance types, by [bradrydzewski](https://github.com/bradrydzewski). [d524689b].(https://github.com/drone/autoscaler/commit/d524689bbd1ed73ef8ee77cb3e0c5e6e6f786158).

## [1.7.1]
### Added
- support for google compute private ip, by [swjclarke](https://github.com/swjclarke).
- support for google compute service accounts, by [ademariag](https://github.com/ademariag).

### Fixed
- google compute instance scopes being ignored, by [ademariag](https://github.com/ademariag).

## [1.7.0]
### Added
- parameter to configure docker stop timeout duration.
- parameter to configure aws volume iops, by [ttousai](https://github.com/ttousai).
- parameter to configure gcp scopes, by [imranismail](https://github.com/imranismail).
- metrics to track server boot errors
- metrics to track server boot time
- metrics to track server installation errors
- metrics to track server installation time
- metrics to track server creation errors
- metrics to track server creation time

### Fixed
- do not run docker stop if the instance was not created.
- do not run docker stop if the instance was not assigned an IP.

## [1.6.1]
### Added
- support for instance not found errors in gcp, by [frebib](https://github.com/frebib).

### Fixed
- resume instance removal when autoscaler unexpectedly restarted, by [@bradrydzewski](https://github.com/bradrydzewski).

## [1.6.0]
### Changed
- Use logrus for logging instead of zerolog, by [@bradrydzewski](https://github.com/bradrydzewski).

### Added
- Read only user interface to visualize servers and logs, by [@bradrydzewski](https://github.com/bradrydzewski). 
- Support for configuring subnetworks with GCP, by [@nsigarora](https://github.com/nsigarora).
- Support for handling  ErrInstanceNotFound with Hetzner, by [@tboerger](https://github.com/tboerger).

## [1.5.0]
### Changed
- Use the new Docker runner image and deprecate the agent, by [@bradrydzewski](https://github.com/bradrydzewski).
- Enable Digital Ocean private IP addresses, by [@barrypeng6](https://github.com/barrypeng6).

## [1.4.3]
### Fixed
- Expired context preventing database updates, by [@bradrydzewski](https://github.com/bradrydzewski).

## [1.4.2]
### Added
- Log errors updating the instance state, by [@bradrydzewski](https://github.com/bradrydzewski).
- Add mutex to database operations for sqlite, by [@bradrydzewski](https://github.com/bradrydzewski).

## [1.4.1] - 2019-10-10
### Fixed
- Support for arm machines on Scaleway, by [@tboerger](https://github.com/tboerger).

## [1.4.0] - 2019-09-23
### Added
- Ability to configure the reaper internal, by [@msaizar](https://github.com/msaizar).
- Ability to configure the install check deadline, by [@bradrydzewski](https://github.com/bradrydzewski).
- Ability to configure the install check interval, by [@bradrydzewski](https://github.com/bradrydzewski).

## [1.3.0] - 2019-09-11
### Added

- Added support for Scaleway, by [@frebib](https://github.com/frebib). [#45](https://github.com/drone/autoscaler/pull/45).

### Fixed

- Fixed issue where non-existing instance could not be destroyed, by [@jlesage](https://github.com/jlesage). [#50](https://github.com/drone/autoscaler/pull/50).
- Added timeout when attempting to ping the instance, by [@bradrydzewski](https://github.com/bradrydzewski).

## [1.2.2] - 2019-08-29
### Added

- Support for loading runner environment variables from file, by [@bradrydzewski](https://github.com/bradrydzewski).
- Basic support for configuring windows agents, by [@bradrydzewski](https://github.com/bradrydzewski).

### Fixed

- Pull garbage collector image before creating the container, by [@msaizar](https://github.com/msaizar).
- Handle nil pointer caused by empty or missing interface in AWS driver, by [@bradrydzewski](https://github.com/bradrydzewski).

## [1.2.1] - 2019-08-14
### Added

- Added postgres driver, by [@mmuehlberger](https://github.com/mmuehlberger).
- Support for capacity buffer, by [@jones2026](https://github.com/jones2026). [#39](https://github.com/drone/autoscaler/pull/39).

### Fixed

- Close docker client after server ping, by [@msaizar](https://github.com/msaizar), [#42](https://github.com/drone/autoscaler/pull/42).

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
