# kvstore

A lightweight distributed KV store built in Go, with TCP networking and concurrent client handling.

## Table of Contents

- [Quick Start](#quick-start)
- [Project Goals](#project-goals)
- [Milestones](#milestones)

## Quick Start

### Start the Server

First, `cd` into the project directory.

Then, run the server with `$ go run main.go`.

### Start a Client

Open up a new terminal, then run `$ telnet 127.0.0.1 8080`.

### Commands

```
GET [key] # Get the value corresponding to a key

SET [key] [value] # Set a key to a corresponding value

DELETE [key] # Remove a key
```

## Project Goals

- Understand persistence and durability guarantees
- Learn replication and consistency models
- Practice network programming and protocol design
- Foundational system design knowledge for interviews

## Milestones

### Phase 1: Basic KV Store (Week 1)
- [x] TCP server with concurrent client handling
- [x] Core operations: GET, SET, DELETE
- [ ] Thread-safe map with RWMutex
- [ ] Graceful shutdown with signal handling
- [ ] Enhanced protocol (support multi-word values)
- [ ] Error handling improvements

**Key Learnings:** Concurrency, network programming, basic client-server architecture

---

### Phase 2: Persistence (Week 1)
- [ ] Append-Only File (AOF) logging
- [ ] Replay log on startup
- [ ] Write-Ahead Log (WAL) implementation
- [ ] Log compaction/snapshotting
- [ ] Configurable fsync policies (every write vs batched)

**Key Learnings:** Durability vs performance trade-offs, crash recovery

---

### Phase 3: Replication (Week 2)
- [ ] Leader-follower architecture
- [ ] Replication protocol design
- [ ] Async replication with replication lag tracking
- [ ] Follower promotion on leader failure
- [ ] Read-your-writes consistency

**Key Learnings:** CAP theorem, consistency models, distributed state

---

### Phase 4: Advanced Features (Week 2)
- [ ] TTL (time-to-live) for keys
- [ ] Basic transactions (MULTI/EXEC)
- [ ] Pub/Sub support
- [ ] Client connection pooling
- [ ] Metrics and monitoring (Prometheus?)

**Key Learnings:** Complex distributed primitives, observability

---