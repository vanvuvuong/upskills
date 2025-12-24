# LỘ TRÌNH HỌC TẬP - Senior DevOps/SRE

_Bắt đầu: 2024-12-23_

---

## Giới thiệu & Đánh giá tổng quan

### Mục tiêu

- Thiết kế, xây dựng và vận hành hệ thống lớn, ổn định
- Có khả năng đào tạo và huấn luyện thế hệ kỹ sư tiếp theo

### Tài liệu chính

| Sách                                                       | Vai trò                           | Độ ưu tiên |
| ---------------------------------------------------------- | --------------------------------- | ---------- |
| **Designing Data-Intensive Applications (DDIA)**           | Nền tảng tư duy thiết kế hệ thống | BẮT BUỘC   |
| **Site Reliability Engineering (Google SRE)**              | Triết lý và thực hành vận hành    | BẮT BUỘC   |
| **Chaos Engineering: System Resiliency** (Rosenthal)       | Lý thuyết chaos engineering       | CAO        |
| **Chaos Engineering: Controlled Disruption** (Pawlikowski) | Thực hành chaos engineering       | CAO        |

### Tại sao thứ tự này?

```
DDIA Ch 8-9 → SRE Ch 3-6 → Chaos Theory → Chaos Practice
     ↓              ↓              ↓              ↓
 [TẠI SAO lỗi]  [ĐO LƯỜNG]    [CÁCH PHÁ]    [TỰ ĐỘNG]
```

- **DDIA** dạy TẠI SAO distributed systems fail → tiền đề hiểu chaos
- **SRE** dạy CÁCH ĐO độ tin cậy → cần trước khi chạy experiments
- **Rosenthal** dạy THIẾT KẾ experiment → cần trước khi hands-on
- **Pawlikowski** dạy TRIỂN KHAI → áp dụng ngay

### Kết quả mong đợi

Sau khi hoàn thành:

- Vẽ được trade-off diagram cho bất kỳ system nào
- Định nghĩa được SLO/SLI cho service đang vận hành
- Viết được postmortem template chuẩn
- Thiết kế và chạy được Chaos experiment
- Tổ chức được GameDay với team
- Có tài liệu đào tạo junior

### Thời gian ước tính

- Giai đoạn 1: 8-10 tuần (đọc + ghi chú)
- Giai đoạn 2: 4-6 tuần
- Giai đoạn 3: Song song với công việc

---

## Giai đoạn 1: Nền tảng tư duy (8-10 tuần)

---

### Tuần 1-2: DDIA Phần I (Ch 1-4)

**Thời gian:** 5-6 giờ

---

#### Chương 1: Ứng dụng đáng tin cậy, mở rộng, bảo trì

**Trang:** 3-22

##### Trọng tâm khi đọc

1. **Ba thuộc tính cốt lõi:**

   - Reliability: hệ thống hoạt động đúng ngay cả khi có lỗi
   - Scalability: xử lý được tải tăng lên
   - Maintainability: dễ vận hành, hiểu, thay đổi

2. **Các loại lỗi:**

   - Hardware faults: disk crash, RAM lỗi, mất điện
   - Software errors: bug, resource leak, cascading failures
   - Human errors: config sai, deploy lỗi → chiếm phần lớn outages

3. **Đo lường hiệu năng:**

   - Throughput: số request/giây
   - Response time: thời gian từ request đến response
   - Latency: thời gian chờ xử lý (không bao gồm network)
   - Percentiles: p50, p95, p99, p999

4. **Bài toán Twitter fanout:**
   - Push model: ghi vào timeline mỗi follower khi post
   - Pull model: query followers khi đọc timeline
   - Hybrid: push cho user thường, pull cho celebrity

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: Reliability vs Availability vs Durability
- [ ] Hiểu: Tail latency (p99, p999) và tại sao quan trọng
- [ ] Hiểu: Fan-out problem và trade-offs

##### Câu hỏi kiểm tra

- [ ] **Q1:** Tại sao p99 latency quan trọng hơn average?

  > _Gợi ý:_ Nghĩ về trải nghiệm user khi 1% requests chậm. User đó có thể là user quan trọng nhất (nhiều data, nhiều requests). Một user chậm = nhiều complaints.

- [ ] **Q2:** Twitter fanout: push vs pull trade-off là gì?

  > _Gợi ý:_ Push = write-heavy (ghi vào timeline triệu follower). Pull = read-heavy (query mỗi lần mở app). Celebrity với 30M followers → push không khả thi.

- [ ] **Q3:** Vertical vs horizontal scaling: khi nào chọn gì?
  > _Gợi ý:_ Vertical đơn giản nhưng có giới hạn (không thể mua máy vô hạn lớn). Horizontal phức tạp hơn (distributed systems) nhưng scale vô hạn.

##### Câu hỏi gợi mở

- Hệ thống mày đang vận hành có đo p99 không? Nếu không, tại sao?
- Khi nào "đủ reliable"? 99.9% hay 99.99%? Cost-benefit là gì?

##### Ghi chú của tao

---

#### Chương 2: Mô hình dữ liệu và ngôn ngữ truy vấn

**Trang:** 27-63

##### Trọng tâm khi đọc

1. **Relational Model (SQL):**

   - Data tổ chức thành tables với rows
   - Relationships qua foreign keys
   - Schema cố định, thay đổi cần migration
   - Tốt cho: data có cấu trúc rõ ràng, nhiều joins

2. **Document Model (MongoDB, CouchDB):**

   - Data lưu dạng JSON/BSON documents
   - Flexible schema, mỗi document có thể khác nhau
   - Tốt cho: data tự chứa, ít joins, schema hay thay đổi

3. **Graph Model (Neo4j):**

   - Nodes và edges với properties
   - Tốt cho: social networks, fraud detection, knowledge graphs

4. **Schema-on-read vs Schema-on-write:**
   - Schema-on-write (SQL): validate khi ghi, lỗi ngay nếu sai schema
   - Schema-on-read (NoSQL): interpret khi đọc, flexible nhưng có thể đọc data lỗi

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: Relational vs Document vs Graph trade-offs
- [ ] Hiểu: Schema-on-read vs Schema-on-write
- [ ] Skim: Graph queries (Cypher, SPARQL) - ít dùng trong DevOps

##### Câu hỏi kiểm tra

- [ ] **Q1:** Khi nào chọn MongoDB thay vì PostgreSQL?
  > _Gợi ý:_ Document DB tốt khi: (1) data tự chứa, không cần joins, (2) schema thay đổi thường xuyên, (3) cần horizontal scale dễ. Relational tốt khi: (1) data có relationships phức tạp, (2) cần ACID transactions, (3) cần ad-hoc queries.

##### Câu hỏi gợi mở

- Hệ thống mày dùng DB gì? Trade-off khi chọn DB đó là gì?
- Nếu thiết kế lại từ đầu, mày có chọn DB khác không?

##### Ghi chú của tao

---

#### Chương 3: Lưu trữ và truy xuất

**Trang:** 69-103

**ĐÂY LÀ CHƯƠNG QUAN TRỌNG - ĐỌC KỸ**

##### Trọng tâm khi đọc

1. **Log-Structured Storage (LSM-Tree):**

   - Ghi append-only vào memory (memtable)
   - Khi đầy, flush ra disk thành SSTable (sorted)
   - Background compaction: merge SSTables
   - **Ưu điểm:** Write throughput cao (sequential writes)
   - **Nhược điểm:** Read có thể chậm (check multiple SSTables)
   - **Dùng bởi:** Cassandra, RocksDB, LevelDB, HBase

2. **B-Tree:**

   - Cây cân bằng, mỗi node là một page trên disk
   - Updates in-place (overwrite page)
   - **Ưu điểm:** Read nhanh, predictable performance
   - **Nhược điểm:** Write amplification, random I/O
   - **Dùng bởi:** PostgreSQL, MySQL, Oracle

3. **Write Amplification:**

   - Tỷ lệ data ghi thực tế trên disk / data ghi logic
   - LSM-Tree: compaction gây write amplification
   - B-Tree: mỗi update có thể rewrite cả page

4. **OLTP vs OLAP:**

   - OLTP (Online Transaction Processing): nhiều small transactions, read/write random rows
   - OLAP (Online Analytical Processing): ít queries nhưng scan millions of rows, aggregations

5. **Column-oriented Storage:**
   - Lưu theo cột thay vì hàng
   - Compression tốt hơn (data cùng cột similar)
   - Tốt cho OLAP queries (chỉ đọc columns cần)

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: LSM-Tree cấu trúc và flow
- [ ] Hiểu: B-Tree cấu trúc và flow
- [ ] Hiểu: OLTP vs OLAP workloads
- [ ] Hiểu: Write Amplification và impact
- [ ] Hiểu: Column-oriented storage

##### Câu hỏi kiểm tra

- [ ] **Q1:** Tại sao Cassandra/RocksDB dùng LSM-tree?

  > _Gợi ý:_ Cả hai optimize cho write-heavy workloads. Sequential writes nhanh hơn random writes 100-1000x trên HDD, 10-100x trên SSD. Time-series data, logs, events = write-heavy.

- [ ] **Q2:** Tại sao PostgreSQL/MySQL dùng B-tree?

  > _Gợi ý:_ OLTP workloads cần read nhanh và predictable. B-tree cho O(log n) lookups. Transactions cần update-in-place cho ACID guarantees.

- [ ] **Q3:** Write amplification là gì? Ảnh hưởng disk I/O thế nào?
  > _Gợi ý:_ Nếu write amplification = 10, ghi 1GB data thực tế ghi 10GB lên disk. SSD có giới hạn write cycles, write amplification cao = SSD chết nhanh hơn.

##### Câu hỏi gợi mở

- SyncoraDMP MES của mày thuộc OLTP hay OLAP? Tại sao?
- Nếu cần thêm analytics trên production data, mày sẽ thiết kế thế nào để không ảnh hưởng OLTP?

##### Ghi chú của tao

---

#### Chương 4: Mã hóa và tiến hóa

**Trang:** 111-139

##### Trọng tâm khi đọc

1. **Encoding formats:**

   - JSON/XML: human-readable, verbose, no schema
   - Protobuf/Thrift: binary, schema required, compact
   - Avro: binary, schema in file, good for Hadoop

2. **Schema Evolution:**

   - Thêm field mới: cần default value
   - Xóa field: chỉ xóa optional fields
   - Đổi type: phải compatible (int32 → int64 OK, int → string NO)

3. **Forward vs Backward Compatibility:**

   - **Forward:** Code cũ đọc được data từ code mới
   - **Backward:** Code mới đọc được data từ code cũ
   - Cần CẢ HAI khi rolling deployment

4. **Tại sao quan trọng cho DevOps:**
   - Rolling deployment: old và new code chạy cùng lúc
   - Database migration: old data phải đọc được bởi new code
   - Microservices: service A và B có thể ở versions khác nhau

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: Forward vs Backward compatibility
- [ ] Hiểu: JSON vs Protobuf vs Avro trade-offs
- [ ] Hiểu: Schema evolution rules

##### Câu hỏi kiểm tra

- [ ] **Q1:** Khi deploy service mới, tại sao cần backward compatible schema?
  > _Gợi ý:_ Rolling deployment = old instances vẫn chạy khi new instances starting. Old code phải đọc được data mà new code ghi (forward). New code phải đọc được data mà old code đã ghi (backward).

##### Câu hỏi gợi mở

- Hệ thống mày dùng format gì cho API? JSON? Protobuf?
- Đã bao giờ gặp bug do schema incompatibility khi deploy chưa?

##### Ghi chú của tao

---

### Tuần 3-4: DDIA Phần II (Ch 5-7)

**Thời gian:** 6-8 giờ

---

#### Chương 5: Replication

**Trang:** 151-192

##### Trọng tâm khi đọc

1. **Leader-based Replication:**

   - Một leader nhận writes
   - Followers replicate từ leader
   - Reads có thể từ leader hoặc followers

2. **Sync vs Async Replication:**

   - Sync: đợi follower confirm trước khi ack client → durable nhưng chậm
   - Async: ack client ngay → nhanh nhưng có thể mất data nếu leader crash

3. **Replication Lag Problems:**

   - **Read-your-writes:** User ghi xong đọc lại không thấy
   - **Monotonic reads:** User đọc data mới, refresh thấy data cũ
   - **Consistent prefix:** Thấy answer trước question

4. **Multi-leader Replication:**

   - Nhiều datacenters, mỗi DC có leader
   - Conflict resolution: last-write-wins, merge, custom logic
   - Dùng khi: multi-DC, offline clients (CouchDB)

5. **Leaderless Replication (Dynamo-style):**
   - Client ghi vào nhiều nodes
   - Quorum: w + r > n đảm bảo overlap
   - Dùng bởi: Cassandra, Riak, Voldemort

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: Leader-based replication flow
- [ ] Hiểu: Replication lag và các anomalies
- [ ] Hiểu: Quorum (w + r > n)
- [ ] Hiểu: Conflict resolution strategies

##### Câu hỏi kiểm tra

- [ ] **Q1:** Read-your-writes consistency giải quyết vấn đề gì?

  > _Gợi ý:_ User update profile, refresh page, thấy profile cũ → confusing. Giải pháp: đọc từ leader cho data user vừa ghi, hoặc track timestamp và đợi follower catch up.

- [ ] **Q2:** Khi nào cần multi-leader replication?

  > _Gợi ý:_ (1) Multi-datacenter: latency thấp cho writes ở mỗi DC. (2) Offline operation: mobile app cần hoạt động offline rồi sync sau. Trade-off: conflict resolution phức tạp.

- [ ] **Q3:** Quorum w + r > n nghĩa là gì?
  > _Gợi ý:_ n = tổng nodes, w = nodes phải ack write, r = nodes phải query khi read. w + r > n đảm bảo ít nhất 1 node có data mới nhất trong read set. Ví dụ: n=3, w=2, r=2 → guaranteed overlap.

##### Câu hỏi gợi mở

- Database của mày dùng replication model nào?
- Đã bao giờ gặp replication lag gây bug chưa?

##### Ghi chú của tao

---

#### Chương 6: Partitioning

**Trang:** 199-216

##### Trọng tâm khi đọc

1. **Tại sao partition:**

   - Data quá lớn cho 1 node
   - Throughput quá cao cho 1 node
   - Mỗi partition trên node khác nhau

2. **Partitioning Strategies:**

   - **Key range:** A-M trên node1, N-Z trên node2
     - Pro: range queries hiệu quả
     - Con: hot spots nếu keys không đều
   - **Hash:** hash(key) mod n
     - Pro: distribute đều
     - Con: range queries không hiệu quả

3. **Hot Spots:**

   - Một partition nhận traffic nhiều hơn
   - Ví dụ: celebrity user, trending topic
   - Giải pháp: add random suffix, split hot partition

4. **Consistent Hashing:**

   - Nodes và keys đều hash vào ring
   - Khi add/remove node, chỉ di chuyển ít data
   - Virtual nodes: mỗi physical node có nhiều positions

5. **Secondary Indexes:**
   - **Document-partitioned (local):** mỗi partition index riêng, query phải scatter-gather
   - **Term-partitioned (global):** index partitioned by term, write cần update nhiều partitions

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: Hash vs Range partitioning trade-offs
- [ ] Hiểu: Hot spots và cách xử lý
- [ ] Hiểu: Consistent hashing
- [ ] Hiểu: Rebalancing strategies

##### Câu hỏi kiểm tra

- [ ] **Q1:** Tại sao hash partitioning giúp tránh hot spots?

  > _Gợi ý:_ Hash function distribute keys đều qua partitions. Key "user_123" và "user_124" có thể ở partitions khác nhau. Nhưng không hoàn toàn tránh được: celebrity user có nhiều requests đến cùng partition.

- [ ] **Q2:** Consistent hashing giải quyết vấn đề gì khi rebalance?
  > _Gợi ý:_ Traditional hash: add node → hash(key) mod n thay đổi → phải di chuyển hầu hết data. Consistent hashing: add node chỉ di chuyển data từ neighbors → minimal disruption.

##### Câu hỏi gợi mở

- Hệ thống mày có bị hot spot không? Làm sao detect?
- Nếu cần scale database 10x, strategy là gì?

##### Ghi chú của tao

---

#### Chương 7: Transactions

**Trang:** 221-266

##### Trọng tâm khi đọc

1. **ACID thực sự nghĩa gì:**

   - **Atomicity:** Transaction hoàn thành hết hoặc rollback hết
   - **Consistency:** Database luôn ở valid state (application-level)
   - **Isolation:** Concurrent transactions như chạy sequential
   - **Durability:** Committed data không mất

2. **Isolation Levels:**

   - **Read Uncommitted:** Có thể đọc uncommitted data (dirty read)
   - **Read Committed:** Chỉ đọc committed data
   - **Repeatable Read (Snapshot):** Thấy snapshot tại thời điểm transaction start
   - **Serializable:** Full isolation, như chạy sequential

3. **Anomalies:**

   - **Dirty read:** Đọc data chưa commit
   - **Non-repeatable read:** Đọc 2 lần, giá trị khác nhau
   - **Phantom read:** Query 2 lần, số rows khác nhau
   - **Write skew:** 2 transactions đọc cùng data, ghi dựa trên đọc, cả hai succeed nhưng kết quả invalid

4. **Implementing Serializable:**
   - **Actual serial execution:** Chạy 1 transaction tại 1 thời điểm (VoltDB)
   - **Two-phase locking (2PL):** Lock tất cả rows cần, release khi commit
   - **Serializable Snapshot Isolation (SSI):** Optimistic, detect conflicts

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: ACID definitions
- [ ] Hiểu: Isolation levels và anomalies
- [ ] Hiểu: Write skew problem
- [ ] Hiểu: 2PL vs SSI trade-offs

##### Câu hỏi kiểm tra

- [ ] **Q1:** Dirty read vs Non-repeatable read vs Phantom read khác nhau thế nào?

  > _Gợi ý:_
  >
  > - Dirty: đọc data CHƯA commit (nếu rollback thì sao?)
  > - Non-repeatable: đọc row 2 lần, VALUE khác (ai đó update)
  > - Phantom: query 2 lần, số ROWS khác (ai đó insert/delete)

- [ ] **Q2:** Tại sao serializable isolation đắt? Trade-off là gì?
  > _Gợi ý:_ 2PL: nhiều locks → deadlocks, low throughput. SSI: abort rate cao khi contention. Serial execution: không tận dụng multi-core. Trade-off: correctness vs performance.

##### Câu hỏi gợi mở

- Database của mày dùng isolation level gì default?
- Đã bao giờ gặp bug do race condition trong database chưa?

##### Ghi chú của tao

---

### Tuần 5-6: DDIA Ch 8-9 (QUAN TRỌNG NHẤT)

**Thời gian:** 6-8 giờ

**ĐÂY LÀ HAI CHƯƠNG QUAN TRỌNG NHẤT - NỀN TẢNG CHO CHAOS ENGINEERING**

---

#### Chương 8: Vấn đề với hệ thống phân tán

**Trang:** 273-310

##### Trọng tâm khi đọc

1. **Partial Failures:**

   - Trong distributed systems, một phần có thể fail trong khi phần khác hoạt động
   - Không thể biết chắc operation thành công hay thất bại
   - Khác với single machine: crash hoặc hoạt động, không có trạng thái giữa

2. **Unreliable Networks:**

   - Request có thể: lost, queued, delayed
   - Response có thể: lost, delayed
   - KHÔNG THỂ phân biệt: node crash vs network delay vs slow response
   - Timeout: quá ngắn → false positives, quá dài → slow detection

3. **Unreliable Clocks:**

   - **Time-of-day clocks:** Synchronized qua NTP, có thể jump backward
   - **Monotonic clocks:** Chỉ để đo duration, không sync giữa nodes
   - **Clock skew:** Các nodes có thời gian khác nhau (có thể hàng giây)
   - KHÔNG NÊN dùng timestamps để ordering events trong distributed systems

4. **Process Pauses:**

   - GC pause: có thể hàng giây
   - VM live migration
   - OS swap, disk I/O
   - Node có thể "đứng hình" và không biết mình đã dừng

5. **Byzantine Faults:**

   - Node có thể lie, send incorrect data
   - Thường assume non-Byzantine trong datacenter (trust machines)
   - Cần Byzantine fault tolerance trong: blockchain, untrusted environment

6. **Split Brain:**
   - Network partition chia cluster thành 2 groups
   - Mỗi group nghĩ group kia đã chết
   - Cả 2 elect leader riêng → data divergence, corruption

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: Partial failures và implications
- [ ] Hiểu: Network unreliability
- [ ] Hiểu: Clock skew problems
- [ ] Hiểu: Split brain scenario

##### Câu hỏi kiểm tra

- [ ] **Q1:** Tại sao không thể phân biệt network delay vs node crash?

  > _Gợi ý:_ Khi không nhận response, có thể: (1) node crash, (2) node đang xử lý chậm, (3) request lost, (4) response lost, (5) network congestion. Không có cách nào biết chắc. Timeout chỉ là heuristic.

- [ ] **Q2:** Clock skew gây ra vấn đề gì trong distributed systems?

  > _Gợi ý:_ Nếu dùng timestamp để order events: Node A ghi lúc 10:00:01, Node B ghi lúc 10:00:00 (clock chậm). Dù B ghi sau, timestamp B < A → wrong order. Last-write-wins có thể chọn sai winner.

- [ ] **Q3:** Split brain là gì và tại sao nguy hiểm?
  > _Gợi ý:_ Cluster 5 nodes, network split 2-3. Cả 2 groups elect leader (nếu không có quorum). 2 leaders accept writes → 2 versions of data → conflict, data loss khi rejoin.

##### Câu hỏi gợi mở

- Hệ thống mày xử lý network timeout thế nào? Timeout bao lâu?
- Đã bao giờ gặp split brain chưa? Hậu quả là gì?

##### Ghi chú của tao

---

#### Chương 9: Tính nhất quán và Đồng thuận

**Trang:** 321-372

##### Trọng tâm khi đọc

1. **Linearizability:**

   - Mọi operation có vẻ atomic và xảy ra tại một thời điểm
   - Như có một single copy of data
   - Cần cho: leader election, locks, unique constraints
   - Đắt: latency cao, không available khi partition

2. **CAP Theorem:**

   - **Consistency (Linearizability):** Mọi read thấy write gần nhất
   - **Availability:** Mọi request nhận response (không cần đúng)
   - **Partition tolerance:** Hoạt động khi network partition
   - Khi partition xảy ra, PHẢI chọn C hoặc A, không thể cả hai
   - Thực tế: "Consistency vs Latency" trade-off

3. **Eventual Consistency:**

   - Nếu không có writes mới, cuối cùng tất cả reads return cùng value
   - Weak guarantee: "eventually" có thể là mãi mãi
   - Dùng khi: có thể tolerate stale reads, cần availability

4. **Consensus:**

   - Nhiều nodes đồng ý về một value
   - Cần cho: leader election, atomic commit
   - **Paxos:** Khó hiểu, khó implement đúng
   - **Raft:** Dễ hiểu hơn, widely used (etcd, Consul)
   - **Zab:** Dùng bởi ZooKeeper

5. **Raft Basics:**

   - Leader election: term numbers, majority votes
   - Log replication: leader append, followers replicate
   - Safety: committed entries không bị overwrite
   - Failure: heartbeat timeout → new election

6. **Two-Phase Commit (2PC):**
   - Coordinator: prepare → wait all votes → commit/abort
   - Blocking: nếu coordinator crash, participants stuck
   - Không partition tolerant

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: Linearizability definition và use cases
- [ ] Hiểu: CAP theorem (đúng interpretation)
- [ ] Hiểu: Raft consensus basics
- [ ] Hiểu: 2PC và limitations

##### Câu hỏi kiểm tra

- [ ] **Q1:** Linearizability khác với Serializability thế nào?

  > _Gợi ý:_
  >
  > - Linearizability: về single object, real-time ordering
  > - Serializability: về transactions, có thể reorder miễn là kết quả như serial
  > - Có thể có serializable mà không linearizable và ngược lại

- [ ] **Q2:** CAP: tại sao phải chọn C hoặc A khi có P?

  > _Gợi ý:_ Khi partition: Node A và B không thể communicate. Client gửi write đến A. Để consistent, A phải wait B confirm → không available. Để available, A accept write không cần B → không consistent (B có stale data).

- [ ] **Q3:** Raft đạt consensus như thế nào? Leader election hoạt động ra sao?
  > _Gợi ý:_
  >
  > 1. Term numbers tăng mỗi election
  > 2. Candidate request votes, cần majority
  > 3. Voters chỉ vote 1 lần per term
  > 4. Leader với highest term wins
  > 5. Split vote → timeout → new election với higher term

##### Câu hỏi gợi mở

- Hệ thống mày dùng consensus ở đâu? etcd? ZooKeeper? Consul?
- Khi network partition, hệ thống mày chọn C hay A?

##### Ghi chú của tao

---

### Tuần 7-8: SRE Principles (Ch 1-6)

**Thời gian:** 5-6 giờ

---

#### Ch 1-2: Giới thiệu SRE và Môi trường Production

**Trang:** 3-22

##### Trọng tâm khi đọc

1. **SRE vs Traditional Ops:**

   - Sysadmin: manual, reactive, team size scales with system size
   - SRE: software engineering approach, automate, team size scales sublinearly

2. **SRE Tenets:**

   - Durable focus on engineering (≤50% ops work)
   - Pursuit of maximum change velocity without violating SLO
   - Monitoring: not just up/down, but meaningful signals
   - Emergency response: on-call, blameless postmortems
   - Change management: progressive rollouts
   - Demand forecasting and capacity planning
   - Provisioning: fast, automated
   - Efficiency and performance: resource optimization

3. **Error Budgets:**
   - 100% reliability không phải goal (quá đắt)
   - Error budget = 1 - SLO (nếu SLO 99.9%, budget = 0.1%)
   - Dùng budget cho: risky releases, maintenance windows
   - Khi budget cạn → focus stability, no risky changes

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: SRE philosophy và tenets
- [ ] Hiểu: Error budget concept

##### Câu hỏi kiểm tra

- [ ] **Q1:** SRE khác sysadmin truyền thống ở điểm nào?
  > _Gợi ý:_ (1) Code để giải quyết ops problems, (2) Share ownership với dev, (3) Error budget thay vì "zero defects", (4) Blameless culture, (5) Data-driven decisions (SLO, SLI).

##### Ghi chú của tao

---

#### Ch 3-4: Quản lý rủi ro và SLO

**Trang:** 25-47

##### Trọng tâm khi đọc

1. **Managing Risk:**

   - Reliability có cost: development, opportunity, operational
   - Mỗi nine (99% → 99.9%) đắt hơn gấp ~10x
   - Chọn target dựa trên: user expectations, revenue impact, dependencies

2. **SLI (Service Level Indicator):**

   - Metric đo lường service behavior
   - Ví dụ: request latency, error rate, throughput, availability
   - Chọn SLIs mà users care about

3. **SLO (Service Level Objective):**

   - Target value cho SLI
   - Ví dụ: p99 latency < 200ms, availability ≥ 99.9%
   - Internal commitment, có thể thay đổi

4. **SLA (Service Level Agreement):**

   - Contract với customers
   - Có consequences nếu miss (refund, penalty)
   - Thường conservative hơn SLO

5. **Error Budget:**
   - Error budget = 1 - SLO
   - Ví dụ: SLO 99.9% → budget 0.1% → 43 mins/month downtime allowed
   - Dùng cho risky changes, experiments
   - Khi cạn: freeze changes, focus stability

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: SLI vs SLO vs SLA
- [ ] Hiểu: Error budget calculation
- [ ] Hiểu: Cost của reliability

##### Câu hỏi kiểm tra

- [ ] **Q1:** SLI vs SLO vs SLA khác nhau thế nào?

  > _Gợi ý:_
  >
  > - SLI: metric (request latency = 150ms)
  > - SLO: target cho metric (p99 latency < 200ms)
  > - SLA: contract với penalty (nếu SLA miss, refund 10%)

- [ ] **Q2:** Error budget được tính như thế nào?

  > _Gợi ý:_ Error budget = 1 - SLO. Nếu SLO = 99.9% (monthly):
  >
  > - Budget = 0.1%
  > - 30 days × 24 hours × 60 mins = 43,200 mins
  > - 43,200 × 0.1% = 43.2 mins downtime allowed

- [ ] **Q3:** Khi error budget cạn, team nên làm gì?
  > _Gợi ý:_ (1) Freeze risky releases, (2) Focus on reliability work, (3) Fix bugs causing errors, (4) Improve monitoring/alerting, (5) Postmortem recent incidents. Không deploy features mới.

##### Câu hỏi gợi mở

- Hệ thống mày có SLO không? Bao nhiêu?
- Có track error budget không? Ai responsible?

##### Ghi chú của tao

---

#### Ch 5-6: Loại bỏ Toil và Giám sát

**Trang:** 49-66

##### Trọng tâm khi đọc

1. **Toil Definition:**

   - Manual: người phải làm, không automated
   - Repetitive: làm đi làm lại
   - Automatable: có thể viết script/tool
   - Tactical: reactive, interrupt-driven
   - No enduring value: không improve service long-term
   - Scales with service: nhiều traffic → nhiều toil

2. **Why Reduce Toil:**

   - SRE spend ≤50% time on toil (policy)
   - Toil > 50% → burnout, no time for engineering
   - Engineering work: automation, tooling, improving systems

3. **4 Golden Signals:**

   - **Latency:** Thời gian serve request (phân biệt success vs error latency)
   - **Traffic:** Demand trên system (requests/sec, concurrent users)
   - **Errors:** Rate of failed requests
   - **Saturation:** Mức độ "full" của resources (CPU, memory, disk, network)

4. **Alerting Philosophy:**
   - Alert on symptoms (user impact), not causes
   - Every alert phải actionable
   - Ưu tiên precision (ít false positives) hơn recall
   - Pages chỉ cho urgent issues cần human intervention

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: Toil definition và examples
- [ ] Hiểu: 4 Golden Signals
- [ ] Hiểu: Alerting best practices

##### Câu hỏi kiểm tra

- [ ] **Q1:** Toil là gì? Tại sao cần giảm toil?

  > _Gợi ý:_ Toil = manual, repetitive work that doesn't add lasting value. Nếu không giảm: (1) Team burnout, (2) Không có time cho engineering, (3) Service không improve, (4) Toil scales với service size → unsustainable.

- [ ] **Q2:** 4 Golden Signals là gì? Tại sao "đủ"?
  > _Gợi ý:_
  >
  > - Latency: is it slow?
  > - Traffic: is it busy?
  > - Errors: is it broken?
  > - Saturation: is it full?
  >   Đủ vì cover: user experience (latency, errors), capacity (traffic, saturation). Mọi problem cuối cùng manifest qua 1 trong 4 signals.

##### Câu hỏi gợi mở

- Công việc hàng ngày của mày có bao nhiêu % là toil?
- Hệ thống monitoring đang track 4 Golden Signals không?

##### Ghi chú của tao

---

### Tuần 9-10: SRE Practices (Ch 11-15)

**Thời gian:** 5-6 giờ

---

#### Ch 11: On-Call

**Trang:** 125-132

##### Trọng tâm khi đọc

1. **On-Call Balance:**

   - Max 25% time on-call (1 week per month)
   - Max 2 incidents per shift (otherwise: understaffed or broken system)
   - Primary + secondary rotation

2. **Incident Response:**

   - Acknowledge alert quickly
   - Assess severity
   - Mitigate first, investigate later
   - Communicate status
   - Escalate if needed

3. **Avoiding Burnout:**
   - Balanced rotations
   - Actionable alerts only (no alert fatigue)
   - Blameless culture
   - Post-incident handoff

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: On-call rotation design
- [ ] Hiểu: Incident response basics

##### Câu hỏi kiểm tra

- [ ] **Q1:** On-call rotation nên thiết kế thế nào để tránh burnout?
  > _Gợi ý:_ (1) Max 25% time on-call, (2) Primary + secondary, (3) Compensatory time off, (4) Actionable alerts only (no noise), (5) Clear escalation path, (6) Blameless culture, (7) Handoff notes giữa shifts.

##### Ghi chú của tao

---

#### Ch 12-14: Xử lý sự cố và Quản lý Incident

**Trang:** 133-166

##### Trọng tâm khi đọc

1. **Effective Troubleshooting:**

   - Problem report → Triage → Examine → Diagnose → Test/Treat → Cure
   - Đừng assume cause trước khi có data
   - Binary search: loại trừ 50% possibilities mỗi bước

2. **Incident Command System:**

   - **Incident Commander (IC):** Điều phối tổng thể, không debug
   - **Operations Lead:** Làm việc trực tiếp với system
   - **Communications Lead:** Update stakeholders
   - Clear roles = less chaos

3. **When to Declare Incident:**

   - SLO at risk
   - Multiple teams involved
   - Customer impact
   - Long duration expected

4. **Communication:**
   - Regular updates (even if "no update")
   - Single source of truth (incident doc/channel)
   - Stakeholder management

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: Systematic troubleshooting
- [ ] Hiểu: Incident Commander role
- [ ] Hiểu: Communication trong incident

##### Câu hỏi kiểm tra

- [ ] **Q1:** Incident Commander làm gì? Tại sao cần role này?

  > _Gợi ý:_ IC điều phối, không debug. Cần vì: (1) Ai đó phải nhìn big picture, (2) Quyết định escalation/communication, (3) Tránh mọi người làm duplicate work, (4) Ensure handoff và documentation.

- [ ] **Q2:** Trong incident, communication quan trọng thế nào?
  > _Gợi ý:_ (1) Stakeholders (management, support, affected teams) cần biết status, (2) Tránh duplicate questions, (3) Build trust (transparent về progress), (4) Historical record cho postmortem.

##### Ghi chú của tao

---

#### Ch 15: Văn hóa Postmortem

**Trang:** 169-175

##### Trọng tâm khi đọc

1. **Blameless Postmortems:**

   - Focus on systems, not people
   - Người làm action X trong context Y, với information Z
   - Mục tiêu: learn và improve, không phải punish

2. **Postmortem Content:**

   - Summary: what happened
   - Impact: users affected, duration, revenue
   - Timeline: chronological events
   - Root cause: không phải "human error"
   - What went well
   - What went wrong
   - Action items: with owners and deadlines

3. **Culture:**
   - Leadership phải support blameless
   - Share postmortems widely (learn across teams)
   - Follow up on action items
   - Celebrate sharing (not hiding) incidents

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: Blameless postmortem philosophy
- [ ] Có thể viết postmortem template

##### Câu hỏi kiểm tra

- [ ] **Q1:** Tại sao postmortem phải blameless?

  > _Gợi ý:_ (1) Nếu blame → người hide mistakes → không learn, (2) "Human error" không phải root cause (tại sao system cho phép human error?), (3) Fear → không report → problems accumulate, (4) Goal là improve system, không phải punish.

- [ ] **Q2:** Postmortem template cần có những gì?
  > _Gợi ý:_
  >
  > 1. Summary (1-2 sentences)
  > 2. Impact (users, duration, severity)
  > 3. Timeline (detailed chronology)
  > 4. Root cause (systemic, not "human error")
  > 5. Contributing factors
  > 6. What went well
  > 7. What went poorly
  > 8. Action items (owner + deadline)
  > 9. Lessons learned

##### Câu hỏi gợi mở

- Team mày có viết postmortem không? Blameless?
- Postmortems có được share và follow up không?

##### Ghi chú của tao

---

## Giai đoạn 2: Chaos Engineering (4-6 tuần)

---

### Tuần 1-2: Rosenthal Ch 1-5 (Theory)

**Sách:** Chaos Engineering: System Resiliency in Practice

---

#### Chương 1-2: Chaos Engineering là gì

##### Trọng tâm khi đọc

1. **Definition:**

   - "The discipline of experimenting on a system to build confidence in the system's capability to withstand turbulent conditions in production"
   - Không phải "breaking things randomly"
   - Controlled experiments với hypothesis

2. **Why Chaos Engineering:**

   - Distributed systems có emergent behavior
   - Testing trong dev/staging không đủ
   - Find weaknesses before they find you
   - Build confidence, not fear

3. **History:**

   - Netflix migration to AWS (2008-2016)
   - Chaos Monkey (2010): randomly kill instances
   - Simian Army: Latency Monkey, Conformity Monkey, etc.
   - Chaos Kong: simulate region failure

4. **Chaos vs Testing:**
   - Testing: verify known properties
   - Chaos: discover unknown properties
   - Testing: pass/fail
   - Chaos: learn about system behavior

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: Chaos Engineering definition
- [ ] Hiểu: History (Netflix origin)
- [ ] Hiểu: Chaos vs traditional testing

##### Ghi chú của tao

---

#### Chương 3-4: Nguyên tắc Chaos Engineering

##### Trọng tâm khi đọc

1. **Principles:**

   - Build hypothesis around steady state behavior
   - Vary real-world events
   - Run experiments in production
   - Automate experiments to run continuously
   - Minimize blast radius

2. **Steady State:**

   - Define normal behavior (metrics, SLIs)
   - Hypothesis: steady state duy trì dù có turbulence
   - Focus on OUTPUT (user-facing metrics), not internal state

3. **Real-world Events:**

   - Server crash
   - Network latency/partition
   - Disk full
   - Clock skew
   - Traffic spike
   - Dependency failure

4. **Blast Radius:**
   - Start small (1 instance, 1% traffic)
   - Có kill switch
   - Monitor closely
   - Expand gradually khi confident

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: 5 principles của Chaos Engineering
- [ ] Hiểu: Steady state concept
- [ ] Hiểu: Blast radius control

##### Câu hỏi kiểm tra

- [ ] **Q1:** Tại sao cần run experiments in production?

  > _Gợi ý:_ (1) Staging ≠ Production (traffic patterns, data, scale), (2) Emergent behaviors only appear in prod, (3) Dependencies behave differently, (4) "Tested in staging" doesn't mean it works in prod.

- [ ] **Q2:** Làm sao minimize blast radius?
  > _Gợi ý:_ (1) Start with 1 instance/1% traffic, (2) Kill switch để stop ngay, (3) Monitor SLIs closely, (4) Run during low-traffic periods first, (5) Have rollback plan, (6) Gradually increase scope.

##### Ghi chú của tao

---

#### Chương 5: Thiết kế Experiment

##### Trọng tâm khi đọc

1. **Experiment Design:**

   - Hypothesis: "If X happens, system still maintains steady state"
   - Independent variable: what you change (kill instance)
   - Dependent variable: what you measure (latency, errors)
   - Control group: baseline behavior

2. **Choosing What to Test:**

   - Critical paths (user-facing)
   - Known weaknesses
   - Recent changes
   - Dependencies
   - High-risk components

3. **Running the Experiment:**

   - Document hypothesis
   - Define steady state metrics
   - Define experiment (what, when, how long)
   - Run in production
   - Observe and record
   - Analyze results
   - Fix or note findings

4. **Common Experiments:**
   - Instance failure
   - Network latency injection
   - Dependency unavailability
   - Resource exhaustion (CPU, memory, disk)
   - Traffic spike

##### Checklist

- [ ] Đọc xong chương
- [ ] Hiểu: Experiment design process
- [ ] Hiểu: Choosing what to test
- [ ] Có thể design experiment cho hệ thống mình

##### Câu hỏi kiểm tra

- [ ] **Q1:** Hypothesis cho chaos experiment nên viết như thế nào?

  > _Gợi ý:_ "Given [steady state], when [event], the system will [expected behavior]". Ví dụ: "Given p99 latency < 200ms, when 1 database replica fails, p99 latency remains < 500ms and no errors visible to users."

- [ ] **Q2:** Làm sao chọn experiment đầu tiên?
  > _Gợi ý:_ (1) Critical user path, (2) Known weakness team đã worry, (3) Recent incident (prevent recurrence), (4) High-risk dependency, (5) Start với đơn giản nhất (kill 1 instance).

##### Ghi chú của tao

---

### Tuần 3: Rosenthal Ch 6-10 (Case Studies)

##### Trọng tâm khi đọc

1. **Netflix:**

   - Chaos Monkey: random instance termination
   - Chaos Kong: region evacuation
   - ChAP: Chaos Automation Platform

2. **Slack:**

   - Disasterpiece Theater: planned chaos events
   - Focus on database failures
   - Team buy-in critical

3. **Capital One:**

   - Financial industry: cautious approach
   - Started with read-only experiments
   - Compliance requirements

4. **Lessons Learned:**
   - Culture > Tools
   - Start small, build trust
   - Executive sponsorship helps
   - Make it easy to run experiments
   - Share findings widely

##### Checklist

- [ ] Đọc xong case studies
- [ ] Note patterns across companies
- [ ] Identify applicable lessons cho context mình

##### Ghi chú của tao

---

### Tuần 4-5: Pawlikowski Part 1-2 (Hands-on)

**Sách:** Chaos Engineering: Site Reliability through Controlled Disruption

---

#### Part 1: Fundamentals (Ch 1-4)

##### Trọng tâm khi đọc

1. **Observability First:**

   - Không thể chaos nếu không observe
   - Metrics, logs, traces
   - Know your steady state

2. **USE Method:**

   - Utilization: % time resource busy
   - Saturation: queue length, wait time
   - Errors: error counts

3. **Linux Forensics:**
   - Process management
   - Resource monitoring
   - Network analysis

##### Checklist

- [ ] Đọc xong Part 1
- [ ] Hiểu: USE method
- [ ] Có thể check system health với Linux tools

##### Ghi chú của tao

---

#### Part 2: Chaos in Action (Ch 5-9)

##### Trọng tâm khi đọc

1. **Docker Chaos (Ch 5):**

   - Kill containers
   - Network manipulation
   - Resource constraints

2. **Syscall Chaos (Ch 6):**

   - Inject failures at syscall level
   - Tools: strace, perf
   - Low-level but powerful

3. **JVM Chaos (Ch 7):**

   - Byteman: inject failures vào Java
   - Memory pressure
   - Exception injection

4. **Application-level (Ch 8):**

   - Fault injection trong code
   - Feature flags for chaos
   - SDK-based injection

5. **Browser Chaos (Ch 9):**
   - Frontend resilience
   - API failure simulation
   - Offline behavior

##### Checklist

- [ ] Đọc xong Part 2
- [ ] Hands-on: chạy thử ít nhất 1 experiment
- [ ] Note tools applicable cho tech stack mình

##### Câu hỏi kiểm tra

- [ ] **Q1:** Khi nào dùng syscall-level vs application-level chaos?
  > _Gợi ý:_ Syscall: low-level, affect mọi thứ, powerful but blunt. Application: precise, targeted, need code changes. Bắt đầu với application-level (dễ control), syscall cho edge cases.

##### Ghi chú của tao

---

### Tuần 6: Pawlikowski Part 3 (Kubernetes Chaos)

---

#### Ch 10-12: Kubernetes Chaos

##### Trọng tâm khi đọc

1. **Pod Chaos (Ch 10):**

   - Pod deletion
   - Container kill
   - Resource limits

2. **Automation (Ch 11):**

   - LitmusChaos
   - Chaos Mesh
   - Gremlin
   - Scheduled experiments

3. **Under the Hood (Ch 12):**
   - Kubelet behavior
   - Control plane failures
   - etcd chaos

##### Checklist

- [ ] Đọc xong Part 3
- [ ] Hiểu: LitmusChaos hoặc Chaos Mesh basics
- [ ] Có thể set up automated chaos experiment

##### Câu hỏi kiểm tra

- [ ] **Q1:** Làm sao automate chaos experiments trong K8s?
  > _Gợi ý:_ (1) CRDs define experiment, (2) Chaos operator execute, (3) Cron schedule, (4) CI/CD integration (run after deploy), (5) Prometheus alerts as abort condition.

##### Ghi chú của tao

---

#### Ch 13: Chaos Engineering for People

**CHƯƠNG QUAN TRỌNG**

##### Trọng tâm khi đọc

1. **GameDay:**

   - Planned chaos event với team
   - Practice incident response
   - Find process gaps

2. **Designing GameDay:**

   - Clear objectives
   - Defined scope
   - Safety mechanisms
   - Observers/facilitators
   - Debrief after

3. **Benefits:**
   - Team readiness
   - Process improvement
   - Build confidence
   - Knowledge sharing

##### Checklist

- [ ] Đọc xong Chapter 13
- [ ] Có thể design GameDay cho team
- [ ] Plan first GameDay

##### Câu hỏi kiểm tra

- [ ] **Q1:** GameDay khác với automated chaos experiments thế nào?

  > _Gợi ý:_ Automated: test system resilience. GameDay: test PEOPLE and PROCESS resilience. Cả hai cần. Automated tìm technical issues, GameDay tìm communication/coordination issues.

- [ ] **Q2:** Làm sao convince team/management chạy GameDay?
  > _Gợi ý:_ (1) Start small (1 service, non-critical), (2) Show other company successes, (3) Frame as learning not testing, (4) Clear safety boundaries, (5) Demonstrate value with first GameDay results.

##### Ghi chú của tao

---

## Giai đoạn 3: Reference (Song song với công việc)

Đọc khi cần:

| Chủ đề                              | Nguồn                  | Khi nào              |
| ----------------------------------- | ---------------------- | -------------------- |
| Load balancing / Cascading failures | SRE Ch 19-22           | Gặp vấn đề tải       |
| Distributed consensus deep dive     | SRE Ch 23              | Làm việc với etcd/ZK |
| Multi-cluster K8s                   | Production K8s Ch 6-10 | Scale K8s            |
| Data pipelines                      | DDIA Ch 10-12          | Thiết kế ETL         |

---

## Ghi chú & Thảo luận

_(Thêm notes và discussions với mentor/Claude ở đây)_

---

## Checklist tổng hợp

### Giai đoạn 1

- [ ] DDIA Ch 1-4 hoàn thành
- [ ] DDIA Ch 5-7 hoàn thành
- [ ] DDIA Ch 8-9 hoàn thành
- [ ] SRE Ch 1-6 hoàn thành
- [ ] SRE Ch 11-15 hoàn thành
- [ ] Viết được SLO cho hệ thống mình
- [ ] Viết được postmortem template

### Giai đoạn 2

- [ ] Rosenthal Ch 1-5 hoàn thành
- [ ] Rosenthal Ch 6-10 hoàn thành
- [ ] Pawlikowski Part 1-2 hoàn thành
- [ ] Pawlikowski Part 3 hoàn thành
- [ ] Design được chaos experiment
- [ ] Chạy được 1 experiment
- [ ] Tổ chức được 1 GameDay

---

_Cập nhật lần cuối: \_\_\_\__
