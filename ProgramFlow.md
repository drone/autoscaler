# Program Flow and Architecture

## Overview

The Drone Autoscaler is a service that automatically scales CI/CD agent infrastructure based on build demand. It monitors pending builds and provisions or terminates server instances dynamically to optimize resource utilization and cost.

## Application Architecture

The application consists of several key components that work together:

### 1. Entry Point (`cmd/drone-autoscaler/main.go`)

The main entry point:
- Loads configuration from environment variables
- Initializes the hosting provider (Amazon, DigitalOcean, Google, etc.)
- Establishes database connection
- Sets up the server store
- Creates and starts the scaling engine
- Starts the HTTP server for metrics and control endpoints

### 2. Core Components

#### Provider Interface (`provider.go`)
The Provider interface is the abstraction layer between the autoscaler and cloud hosting services:

```go
type Provider interface {
    Create(context.Context, InstanceCreateOpts) (*Instance, error)
    Destroy(context.Context, *Instance) error
}
```

**Key Types:**
- `ProviderType`: Enum of supported providers (Amazon, Azure, DigitalOcean, Google, HetznerCloud, Linode, OpenStack, Packet, Scaleway, Vultr)
- `Instance`: Represents a provisioned server with metadata (ID, Name, Address, Region, Image, Size, etc.)
- `InstanceCreateOpts`: Configuration options for instance creation (TLS certificates, SSH keys, etc.)
- `InstanceError`: Captures creation errors with server logs

#### Engine Interface (`engine.go`)
The Engine orchestrates the autoscaling logic:

```go
type Engine interface {
    Start(context.Context)
    Pause()
    Paused() bool
    Resume()
}
```

The Engine manages multiple sub-components that run on configured intervals:

- **Planner**: Calculates current build demand and determines how many servers to create or destroy
- **Allocator**: Provisions new servers by calling the Provider's Create method
- **Installer**: Configures Docker and the Drone agent on newly provisioned servers
- **Pinger**: Keeps heartbeat with running servers
- **Collector**: Gracefully shuts down servers marked for termination
- **Reaper**: Purges old server records from the database

#### Server Store Interface (`server.go`)
Persists server metadata in a database:
- Track server creation/update/deletion
- Query servers by state (Pending, Creating, Created, Staging, Running, Shutdown, Stopping, Stopped, Error)
- Maintain server lifecycle history

#### Configuration (`config/config.go`)
Loads all settings from environment:
- Database connection details
- Provider credentials
- Drone server connection
- Scaling parameters (min/max servers, capacity per server, buffer)
- Agent configuration (Docker image, environment variables, volumes)

#### Metrics (`metrics/`)
Collects and exposes Prometheus metrics:
- Server creation/deletion counts and latencies
- Active server count
- Build queue depth
- Engine performance metrics

## Program Flow

### Startup Sequence

```
1. main() loads configuration
2. setupProvider() initializes cloud provider client
3. Database connection established
4. ServerStore created
5. Drone client authenticated
6. Engine created with all components
7. Engine.Start() begins scaling loop
8. HTTP server started for metrics/control
9. Application waits for shutdown signal
```

### Scaling Loop (runs at configured interval, typically 30 seconds)

```
1. Planner.Plan()
   ├─ Query Drone for pending/running builds
   ├─ Calculate required capacity
   ├─ Compare with current allocation
   └─ Create/mark servers for termination

2. Allocator.Allocate()
   ├─ Find servers in "Pending" state
   ├─ Call Provider.Create() for each
   └─ Update state to "Creating" → "Created"

3. Installer.Install()
   ├─ Find servers in "Created" state
   ├─ Connect via Docker
   ├─ Configure Drone agent
   └─ Update state to "Running"

4. Pinger.Ping()
   ├─ Health check running servers
   └─ Detect stale/failing instances

5. Collector.Collect()
   ├─ Find servers in "Shutdown" state
   ├─ Gracefully stop Drone agent
   ├─ Call Provider.Destroy()
   └─ Update state to "Stopped"

6. Reaper.Reaper()
   ├─ Find old "Stopped" servers (>24h)
   └─ Delete from database
```

## Drivers Package Interface

The `drivers/` package contains provider-specific implementations that satisfy the `Provider` interface. Each driver handles cloud-specific API interactions.

### Driver Structure

Each cloud provider driver follows a consistent pattern:

#### Provider Implementation
Located in `drivers/{provider}/provider.go`

- Contains a `provider` struct with cloud-specific configuration
- Implements the `Provider` interface methods: `Create()` and `Destroy()`
- Uses cloud SDK to communicate with the provider's API

#### Options Pattern (`drivers/{provider}/option.go`)

Drivers use functional options for flexible configuration:

```go
type Option func(*provider)

func WithRegion(region string) Option { ... }
func WithImage(image string) Option { ... }
func WithSize(size string) Option { ... }
// ... more options
```

#### Server Creation (`drivers/{provider}/create.go`)

Handles instance provisioning:
1. Validates input parameters
2. Calls cloud provider API to launch instance
3. Polls for instance readiness
4. Gathers instance metadata (IP address, region, etc.)
5. Returns `Instance` object or `InstanceError`

#### Server Destruction (`drivers/{provider}/destroy.go`)

Handles instance termination:
1. Verifies instance exists
2. Calls cloud provider API to terminate
3. Polls for confirmation
4. Cleans up associated resources

#### Setup (`drivers/{provider}/setup.go`)

Initializes the provider factory:
- Creates cloud SDK clients
- Loads credentials from environment or config
- Validates configuration parameters
- Returns a configured `Provider` instance

#### Utilities (`drivers/{provider}/util.go`)

Helper functions specific to the provider:
- Instance state mapping
- API response parsing
- Error handling
- Configuration defaults

### Supported Drivers

1. **Amazon EC2** (`drivers/amazon/`)
   - Launches EC2 instances
   - Configures VPC/security groups
   - Manages EBS volumes

2. **Microsoft Azure** (`drivers/azure/`)
   - Creates VMs
   - Configures networking
   - Manages storage

3. **DigitalOcean** (`drivers/digitalocean/`)
   - Creates Droplets
   - Configures networking
   - Manages firewall rules

4. **Google Cloud** (`drivers/google/`)
   - Creates Compute Engine instances
   - Configures service accounts
   - Manages firewall rules

5. **Hetzner Cloud** (`drivers/hetznercloud/`)
   - Creates servers
   - Configures networking

6. **OpenStack** (`drivers/openstack/`)
   - Creates instances
   - Configures networks

7. **Packet.com** (`drivers/packet/`)
   - Provisions bare metal servers

8. **Scaleway** (`drivers/scaleway/`)
   - Creates Scaleway servers

### Driver Integration Flow

```
Application
    ↓
setupProvider() [main.go]
    ↓
Reads PROVIDER env var (e.g., "amazon")
    ↓
Calls driver factory (e.g., amazon.New(opts...))
    ↓
Provider struct initialized with:
    - Cloud credentials
    - API clients
    - Default configurations
    ↓
Returns Provider interface
    ↓
Engine uses Provider.Create() and Provider.Destroy()
    ↓
Driver calls cloud provider API
    ↓
Manages instance lifecycle
```

### Provider Initialization Example

From `cmd/drone-autoscaler/main.go`:

```go
provider, err := setupProvider(conf)
if err != nil {
    logrus.Fatalln("Invalid or missing hosting provider")
}

// Wrap with metrics
provider = metrics.ServerCreate(provider)
provider = metrics.ServerDelete(provider)
```

The `setupProvider()` function:
1. Reads PROVIDER environment variable
2. Loads provider-specific credentials
3. Instantiates the appropriate driver with options
4. Wraps provider with metrics collection
5. Returns Provider interface

## Server Lifecycle States

Servers transition through states managed by the Engine and ServerStore:

```
Pending
  ↓ (Allocator processes)
Creating
  ↓ (Provider.Create() called)
Created
  ↓ (Installer configures)
Staging
  ↓
Running ← Pinger health checks here
  ↓ (Planner marks for termination)
Shutdown
  ↓ (Collector processes)
Stopping
  ↓ (Provider.Destroy() called)
Stopped/Error
  ↓ (Reaper removes after TTL)
[Deleted from store]
```

## Configuration and Extensibility

### Adding a New Driver

To add support for a new cloud provider:

1. Create `drivers/{provider}/` directory
2. Implement the `Provider` interface:
   - `Create(ctx context.Context, opts InstanceCreateOpts) (*Instance, error)`
   - `Destroy(ctx context.Context, instance *Instance) error`
3. Implement setup and option functions
4. Add driver initialization in `setupProvider()` 
5. Add provider constant to `provider.go`

### Key Extension Points

- **`Provider` interface**: Add cloud support by implementing these two methods
- **`ServerStore` interface**: Swap database backend (SQL, NoSQL, etc.)
- **`Engine` components**: Customize scaling algorithm, install process, monitoring
- **Metrics**: Add custom metrics collection
- **Slack notifications**: Provide webhook for infrastructure events

## Error Handling

- **Provider errors**: Wrapped in `InstanceError` with cloud logs for debugging
- **Database errors**: Logged and propagated to prevent state inconsistency
- **API timeouts**: Retried with exponential backoff (provider-specific)
- **Instance already exists**: Handled gracefully with idempotent operations

## Graceful Shutdown

On SIGTERM/SIGINT:
1. Engine paused (no new scaling operations)
2. Running build jobs allowed to finish
3. Servers marked for shutdown are destroyed
4. Database connections closed
5. Process exits cleanly
