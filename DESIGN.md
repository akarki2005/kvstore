# Design Decisions

## Persistence

### Context

Any good implementation of a KV store requires a persistence strategy; that is, data must be stored permanently and be easily recoverable in the event of a crash or restart. The fundamental trade-off here is between durability (data survives crashes/restarts) and performance (write throughput). 

### Problem

Disk writes are time-consuming: nearly 1000x slower than writing to memory. The OS mitigates this by batching disk writes (periodically writing all data in RAM to disk), but for any KV store, this is problematic due to a significant amount of data loss if a server crashes or must restart.

### Potential Solutions

#### 1. "Business-As-Usual" (Allow OS to manage disk writes)

Without any additional disk writes (beyond what the OS already does), we get excellent throughput, but run the risk of losing large chunks of data that haven't yet been written to the disk when the server crashes or restarts.

Performance: Best

Durability: Worst

Complexity: Low

#### 2. Write to Disk on Every Operation (Synchronous fsync)

With this approach, we call `fsync` every time we have a write operation (we write to disk every time we write to memory). This ensures data is never lost, but significantly hampers performance due to the latency of calling `fsync` on every write.

Performance: Worst

Durability: Best

Complexity: Medium

#### 3. Group Commit (Middle Ground)

Here, we call `fsync` periodically to batch a group of writes together. This means more frequent disk writes than option 1 and fewer calls to `fsync` than option 2.

Performance: Good

Durability: Good

Complexity: High

### Decision

Ultimately, I chose to implement option 2. While performance is limited with this approach, the data integrity guarantees it offers (ACID compliance, zero data loss) make it a good approach for critical systems in domains like banking, healthcare or finance. 