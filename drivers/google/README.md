# Google Compute Engine driver

This driver provisions and destroys Compute Engine instances for the autoscaler. This document describes how it handles slow, failed and rate-limited calls, and the settings that shape it.

## Configuration

Setting `DRONE_GOOGLE_PROJECT` selects this driver.

| Variable | Default | Purpose |
|---|---|---|
| `DRONE_GOOGLE_PROJECT` | | Project to create instances in. Required. |
| `DRONE_GOOGLE_ZONE` | `us-central1-a` | Zones to create instances in, comma-separated. Each instance is created in a zone picked at random. |
| `DRONE_GOOGLE_MACHINE_TYPE` | `n1-standard-1` | Machine type. |
| `DRONE_GOOGLE_MACHINE_IMAGE` | `ubuntu-os-cloud/global/images/ubuntu-2004-focal-v20220712` | Boot image, as `project/global/images/name`. |
| `DRONE_GOOGLE_DISK_SIZE` | `50` | Boot disk size in GB. |
| `DRONE_GOOGLE_DISK_TYPE` | `pd-standard` | Boot disk type. |
| `DRONE_GOOGLE_NETWORK` | `global/networks/default` | Network to attach instances to. |
| `DRONE_GOOGLE_SUBNETWORK` | | Subnetwork to attach instances to. |
| `DRONE_GOOGLE_STACK_TYPE` | `IPV4_ONLY` | IP stack type. `IPV4_IPV6` also requests an external IPv6 address. |
| `DRONE_GOOGLE_PRIVATE_IP` | `false` | When true, instances get no external IPv4 address and the autoscaler connects over the private address. |
| `DRONE_GOOGLE_TAGS` | `allow-docker` | Network tags, comma-separated. With the default tags, the driver creates a `default-allow-docker` firewall rule (tcp 2376) on first use if one does not exist. |
| `DRONE_GOOGLE_LABELS` | | Labels to apply, as comma-separated `key:value` pairs. |
| `DRONE_GOOGLE_SCOPES` | see description | OAuth scopes for the instance's service account, comma-separated. Defaults to `devstorage.read_only`, `logging.write`, `monitoring.write` and `trace.append`, each under `https://www.googleapis.com/auth/`. |
| `DRONE_GOOGLE_SERVICE_ACCOUNT_EMAIL` | `default` | Service account to attach to instances. |
| `DRONE_GOOGLE_USERDATA` | | Cloud-init template, inline. |
| `DRONE_GOOGLE_USERDATA_FILE` | | Path to a cloud-init template. Takes precedence over `DRONE_GOOGLE_USERDATA` if both are set, and must be readable at startup. |
| `DRONE_GOOGLE_USERDATA_KEY` | `user-data` | Metadata key the template is stored under, so it can be set to a startup script instead. |
| `DRONE_GOOGLE_READ_RATELIMIT` | `25` | Requests per second for operation-status reads. |
| `DRONE_GOOGLE_SERVER_CREATE_TIMEOUT` | `3m` | Longest a `Create` call may take. See [Timeouts](#timeouts). |

The destroy timeout and retry timings below are provider defaults and are not exposed as environment variables.

## Operations

An insert or delete is only treated as finished once Compute Engine reports the operation `DONE`. A failed operation can report its error while it is still running, and a failed insert keeps its instance name reserved until the operation finishes, so acting on the error early could re-insert the same name and hit a conflict. Inserts carry a request ID so a retried insert is idempotent.

An operation that reaches `DONE` with an error is definitive: nothing was created.

Transient errors (408, 500, 502, 503, 504) are retried with a capped exponential backoff. A rate limit, either a 429 or a 403 with reason `rateLimitExceeded`, is retried after the server's `Retry-After` (10 seconds if it sends none).

## Timeouts

- **Create** (`DRONE_GOOGLE_SERVER_CREATE_TIMEOUT`): bounds the whole call, from the insert through its operation poll and the follow-up read.
- **Destroy**: bounded to 15 minutes, covering both the delete request and waiting for the delete operation.

Choose the create timeout above what a slow-but-successful creation takes. A timeout that fires after Google accepted the insert leaves an instance that may still come up, and that instance is then cleaned up (see below), so too low a value turns working creations into failures.

## Unconfirmed inserts

Once Google may have accepted an insert, the instance may exist even if `Create` cannot confirm it. This happens when:
- the insert gets no definitive answer from Google: the request times out or is cancelled, fails at the network level, or keeps returning a transient status (408, 500, 502, 503, 504) through every retry;
- the operation poll gives up, or the timeout elapses while polling;
- the follow-up read fails or times out.

A definitive answer from Google, such as a stockout, quota error, rate limit or conflict, is not in this group: nothing was created. An insert that was never sent, because the context had already expired, is also not in this group. In the cases above `Create`:

- does **not** retry the insert, since the same name would only conflict;
- does **not** delete anything itself;
- returns the instance's identity (ID and zone) together with an error.

The engine records the server as errored along with its ID, and the reaper then destroys it. If nothing was created, the delete returns not-found and is treated as already done. Without this, an errored server has no ID, so the reaper removes its database row and leaves the running VM behind.

## Destroy

`Destroy` retries transient errors and treats a not-found response as success. Any other failure, including a timeout while waiting on the operation, leaves the server in the errored state for the reaper to retry. A retry after a delete that actually succeeded returns not-found and is treated as done.

## Logging

- A slow response, where the create timeout elapsed before the insert finished, logs at Warn.
- A failed insert operation, a failed read-back, and a failed delete operation log at Error.
- A delete that was accepted but could not be confirmed logs at Warn.
