# Tony Architecture: The Ring -3 Agentic Hypervisor

Tony is the definitive agentic operating system kernel and hypervisor layer. It provides the low-level infrastructure, cryptographic primitives, and secure boundaries required for sovereign agentic execution.

## Core Architectural Pillars

### 1. Low-Level IPC & Networking (The UDS Engine)
- **Zero TCP/HTTP Overhead**: Tony rejects local HTTP loops and loopback addresses for internal routing.
- **Unix Domain Sockets (UDS)**: All local communication utilizes UDS for raw Inter-Process Communication (IPC) passed directly through kernel memory buffers.
- **File Descriptor Passing (`SCM_RIGHTS`)**: Live, encrypted streams (SSH, terminal pipes) are passed between processes using `SCM_RIGHTS`, allowing the master shell (`vish`) to delegate streams to guest agents without leaking sensitive setup data.
- **POSIX Security Boundaries**: Access control is enforced via file system permissions (`chmod`/`chown`) on socket endpoints, restricting communication based on OS user IDs (`uid`).

### 2. Pure Non-NIST Cryptography ("Water Logic")
- **NIST Rejection**: Strictly forbids NIST-recommended primitives like `secp256r1`.
- **Curve25519 Standardization**:
    - **Authentication**: `Ed25519` for signatures and attestations.
    - **Encryption**: `X25519` for key exchanges and data sealing.
- **Unified Identity**: Leverages birational equivalence to map a single 32-byte private seed to both `Ed25519` and `X25519` coordinates.

### 3. Hierarchical Deterministic (HD) Identity
- **Tree Derivation**: A root master seed generates a deep tree of independent leaf keys (`m / purpose' / agent_id' / persona' / task_index`).
- **Stateless Personas**: Identity is treated as a microservice. Transient leaf keys are derived on-the-fly for specific tasks (Git signing, DB access) and destroyed immediately after use.

### 4. Anti-Scraping Volatile Memory Defenses
- **Memory Locking (`mlock`)**: Active key strings are pinned to physical RAM using `mlock` to prevent swapping to persistent disk.
- **Microsecond Scrubbing**: Memory is explicitly zeroed (using `sodium_memzero` or atomic equivalents) immediately after cryptographic operations conclude.

### 5. Deployment Topologies
- **Manned Master**: High-trust controller (within `vish`) running on local physical metal with the master seed under direct supervision.
- **Unmanned Workers**: Headless instances on cloud VPS, running in isolated Linux user namespaces and cgroups. They remain blind to the master seed, managing state via sharded MPC.

## Directory Structure

- `pkg/ipc/`: Unix Domain Socket engine and `SCM_RIGHTS` implementation.
- `pkg/crypto/`: Curve25519 unified identity and HD derivation.
- `pkg/memory/`: `mlock` wrappers and memory scrubbing utilities.
- `pkg/identity/`: Stateless persona management and derivation paths.
