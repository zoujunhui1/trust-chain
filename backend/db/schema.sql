-- TrustChain 索引器数据库 schema / Indexer database schema
-- 用法 / usage:  mysql -u root < backend/db/schema.sql
--
-- 设计要点 / design notes:
--   * 这些表大多由链上事件重建，是前端「读」的唯一来源（写走链上 MetaMask）。
--     Most of these tables are rebuilt from on-chain events and are the sole
--     read source for the frontend (writes happen on-chain via MetaMask).
--     唯一例外是 users 表，见下方说明。/ The one exception is `users`, see below.
--   * 金额用 DECIMAL(65,0) 存 wei（uint256）。MySQL DECIMAL 最大 65 位，
--     足够真实 ETH 金额（全网供应量约 27 位），不覆盖理论最大 uint256（78 位）。
--   * 地址统一小写存成 CHAR(42)，bytes32 存成 CHAR(66) 十六进制。
--   * 表与字段均带 COMMENT，可用 SHOW CREATE TABLE 或 information_schema 查看。

CREATE DATABASE IF NOT EXISTS trustchain
  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE trustchain;

-- 索引进度游标：记录每条链已处理到的区块，崩溃重启后从这里续扫。
-- Indexing cursor: last processed block per chain; enables crash-safe resume.
CREATE TABLE IF NOT EXISTS indexer_cursor (
  id          VARCHAR(64)     NOT NULL COMMENT '链标识，如 sepolia / chain id key',
  last_block  BIGINT UNSIGNED NOT NULL COMMENT '已处理到（含）的区块高度 / last processed block',
  updated_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
                                       COMMENT '最后更新时间 / last update time',
  PRIMARY KEY (id)
) ENGINE=InnoDB COMMENT='索引进度游标，崩溃续扫用 / indexing progress cursor';

-- 慈善机构准入状态：来自 CharityRegistry 的 CharityVerified / CharityRevoked。
-- Charity verification status from the registry events.
CREATE TABLE IF NOT EXISTS charities (
  address        CHAR(42)        NOT NULL COMMENT '小写 0x 地址 / lowercase 0x address',
  verified       BOOLEAN         NOT NULL DEFAULT FALSE COMMENT '是否已认证 / verified flag',
  updated_block  BIGINT UNSIGNED NOT NULL COMMENT '最后一次状态变更所在区块 / block of last change',
  PRIMARY KEY (address)
) ENGINE=InnoDB COMMENT='慈善机构准入状态 / charity verification status';

-- 募捐活动：CampaignCreated 时插入，metadata/里程碑金额通过合约 view 补齐；
-- raised/released 由 donations/milestones 聚合得到，保证可重复执行（幂等）。
-- confirmed 是个例外：这行也可能由 API 乐观插入（见 POST /api/campaigns，
-- backend/internal/api），交易一确认前端就直接写库，不用等 indexer 追上
-- CONFIRMATIONS 个块+下一次轮询。indexer 真正处理到这个事件时会把 confirmed
-- 置 true，并且用链上数据覆盖其它字段（以防乐观写入的数据有误）；已经
-- confirmed=true 的行不会再被乐观写入覆盖或降级——见 store.InsertCampaign。
-- confirmed is the one exception: this row may also be optimistically
-- inserted by the API (see POST /api/campaigns, backend/internal/api) the
-- moment the create-campaign tx confirms, instead of waiting for the indexer
-- to catch up (CONFIRMATIONS blocks + the next poll). When the indexer does
-- process the real event, it sets confirmed=true and overwrites the other
-- fields with on-chain truth (in case the optimistic write was wrong); a row
-- already confirmed=true is never overwritten or downgraded by a later
-- optimistic write — see store.InsertCampaign.
CREATE TABLE IF NOT EXISTS campaigns (
  id               BIGINT UNSIGNED NOT NULL COMMENT '链上 campaignId / on-chain campaign id',
  charity          CHAR(42)        NOT NULL COMMENT '发起机构地址 / charity address',
  goal             DECIMAL(65,0)   NOT NULL COMMENT '目标总额(wei)=各里程碑之和 / goal in wei',
  raised           DECIMAL(65,0)   NOT NULL DEFAULT 0 COMMENT '已筹(wei)=SUM(donations) / total raised',
  released         DECIMAL(65,0)   NOT NULL DEFAULT 0 COMMENT '已放款(wei)=已释放里程碑之和 / total released',
  milestone_count  INT UNSIGNED    NOT NULL COMMENT '里程碑数量 / number of milestones',
  metadata_hash    CHAR(66)        NULL COMMENT 'IPFS 元数据哈希(view 补齐) / ipfs metadata hash',
  completed        BOOLEAN         NOT NULL DEFAULT FALSE COMMENT '是否完成(收到末个收据) / completed flag',
  confirmed        BOOLEAN         NOT NULL DEFAULT TRUE COMMENT '是否已被索引器用链上事件确认过 / confirmed against a real on-chain event by the indexer',
  created_block    BIGINT UNSIGNED NOT NULL COMMENT '创建所在区块 / creation block',
  created_tx       CHAR(66)        NOT NULL COMMENT '创建交易哈希 / creation tx hash',
  created_at       TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '入库时间 / indexed at',
  PRIMARY KEY (id),
  KEY idx_charity (charity)
) ENGINE=InnoDB COMMENT='募捐活动主表 / fundraising campaigns';

-- 里程碑：CampaignCreated 时按数量建行（金额由 view 补齐），
-- 之后 MilestoneReleased / ReceiptSubmitted 更新状态与收据哈希。
CREATE TABLE IF NOT EXISTS milestones (
  campaign_id   BIGINT UNSIGNED  NOT NULL COMMENT '所属活动 id / owning campaign id',
  idx           INT UNSIGNED     NOT NULL COMMENT '里程碑序号(从0起) / milestone index',
  amount        DECIMAL(65,0)    NOT NULL COMMENT '本里程碑金额(wei) / milestone amount in wei',
  state         TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=Locked 1=Released 2=Proven',
  receipt_hash  CHAR(66)         NULL COMMENT '收据文件哈希(Proven 后有值) / receipt hash',
  PRIMARY KEY (campaign_id, idx)
) ENGINE=InnoDB COMMENT='活动里程碑及放款/收据状态 / campaign milestones';

-- 逐笔捐款流水：DonationReceived 一条一行。
-- (tx_hash, log_index) 唯一 → 幂等：重扫同一段区块不会重复插入。
CREATE TABLE IF NOT EXISTS donations (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键 / surrogate pk',
  campaign_id   BIGINT UNSIGNED NOT NULL COMMENT '所属活动 id / campaign id',
  donor         CHAR(42)        NOT NULL COMMENT '捐赠者地址 / donor address',
  amount        DECIMAL(65,0)   NOT NULL COMMENT '捐款金额(wei) / donation amount in wei',
  block_number  BIGINT UNSIGNED NOT NULL COMMENT '事件所在区块 / event block',
  tx_hash       CHAR(66)        NOT NULL COMMENT '交易哈希 / tx hash',
  log_index     INT UNSIGNED    NOT NULL COMMENT '日志在交易中的序号 / log index',
  created_at    TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '入库时间 / indexed at',
  PRIMARY KEY (id),
  UNIQUE KEY uniq_log (tx_hash, log_index),
  KEY idx_campaign (campaign_id),
  KEY idx_donor (donor)
) ENGINE=InnoDB COMMENT='逐笔捐款流水(幂等唯一键) / per-donation ledger';

-- 全局活动流水：每处理一个 campaign 相关事件都追加一行，专供 Transparency
-- Dashboard 的活动表用（campaigns/milestones 只存"当前状态"，历史会被覆盖）。
-- Global activity feed: one append-only row per campaign-scoped event, purely
-- for the Transparency Dashboard's activity table (campaigns/milestones only
-- hold current state and get overwritten, so this is the only history).
CREATE TABLE IF NOT EXISTS activity_events (
  id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键 / surrogate pk',
  campaign_id    BIGINT UNSIGNED NOT NULL COMMENT '所属活动 id / campaign id',
  event_type     VARCHAR(32)     NOT NULL COMMENT 'CampaignCreated/DonationReceived/MilestoneReleased/ReceiptSubmitted/CampaignCompleted',
  amount         DECIMAL(65,0)   NULL COMMENT '捐款/放款金额(wei)；其余事件类型为 NULL / wei, NULL for non-monetary events',
  milestone_idx  INT UNSIGNED    NULL COMMENT '里程碑序号；仅放款/举证事件有值 / set only for milestone-scoped events',
  block_number   BIGINT UNSIGNED NOT NULL COMMENT '事件所在区块 / event block',
  tx_hash        CHAR(66)        NOT NULL COMMENT '交易哈希 / tx hash',
  log_index      INT UNSIGNED    NOT NULL COMMENT '日志在交易中的序号 / log index',
  created_at     TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '入库时间 / indexed at',
  PRIMARY KEY (id),
  UNIQUE KEY uniq_log (tx_hash, log_index),
  KEY idx_campaign (campaign_id),
  KEY idx_block (block_number)
) ENGINE=InnoDB COMMENT='全局链上活动流水(幂等唯一键) / global on-chain activity feed';

-- 连接过钱包的用户：这是本 schema 里唯一不是由链上事件重建的表——由 API 在
-- 前端钱包连接成功时直接写入（POST /api/users/connect），而不是索引器。
-- role 是前端当时算出的角色快照（admin/charity/donor，见 frontend/src/lib/role.ts），
-- 不是权限来源——真正的权限判定始终是链上 owner()/isVerified()，这张表只是
-- "谁连过、上次是什么角色、什么时候" 的记录。
-- Wallet-connected users: the one table in this schema NOT rebuilt from
-- on-chain events — the API writes it directly when a wallet connects
-- (POST /api/users/connect), not the indexer. `role` is a snapshot of what
-- the frontend computed at connect time (admin/charity/donor); it is never
-- the source of truth for permissions — that's always the on-chain
-- owner()/isVerified() checks — this table just records who connected, as
-- what role, and when.
CREATE TABLE IF NOT EXISTS users (
  address              CHAR(42)    NOT NULL COMMENT '小写 0x 地址 / lowercase 0x address',
  role                 VARCHAR(16) NOT NULL COMMENT '连接时的角色快照 admin/charity/donor / role snapshot at connect time',
  first_connected_at   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '首次连接时间 / first-seen time',
  last_connected_at    TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
                                             COMMENT '最近一次连接时间 / most recent connect time',
  PRIMARY KEY (address)
) ENGINE=InnoDB COMMENT='连接过钱包的用户(API 直写，非索引器) / wallet-connected users (API-written, not indexer)';

-- 活动标题：第二个不是由链上事件重建的表。合约的 metadataHash 只是 bytes32
-- 哈希（校验用，不是可解析指针），链上事件从没带过标题/详情/主题这些字段，
-- 所以这些字段天生不在 campaigns 表里，靠这张独立的表补上。
-- 特意不建外键关联 campaigns(id)：创建活动的交易一确认，前端就会调用
-- POST /api/campaigns/{id}/metadata，这通常比 indexer 把该活动写进 campaigns
-- 表还快（indexer 要等 CONFIRMATIONS 个块 + 下一次轮询），建外键会导致这个
-- 请求在那段时间差内失败。读的时候用 LEFT JOIN，campaigns 表还没同步到位
-- 也不影响这些字段已经先存上。
-- theme 是可选的：前端建活动时可以选一个主题（对应 frontend/src/lib/theme.ts
-- 里某个主题的 key），没选就留 NULL，前端会退回按 id 确定性分配主题的老逻辑。
-- Campaign title/description/theme: the second table not rebuilt from
-- on-chain events. The contract's metadataHash is just a bytes32 hash (for
-- verification, not a resolvable pointer), so on-chain events never carry
-- any of these — this table fills that gap.
-- Deliberately no FK to campaigns(id): the frontend calls
-- POST /api/campaigns/{id}/metadata as soon as the create-campaign tx
-- confirms, which is usually faster than the indexer inserting that
-- campaign's row (it waits CONFIRMATIONS blocks + the next poll) — an FK
-- would make this write fail during that gap. Reads LEFT JOIN this table, so
-- campaigns not yet indexed don't block metadata from being stored first.
-- theme is optional: the charity can pick one at creation (matching a key in
-- frontend/src/lib/theme.ts); left NULL, the frontend falls back to its old
-- deterministic by-id assignment.
CREATE TABLE IF NOT EXISTS campaign_metadata (
  campaign_id  BIGINT UNSIGNED NOT NULL COMMENT '链上 campaignId / on-chain campaign id',
  title        VARCHAR(200)    NOT NULL COMMENT '活动标题(前端创建时填写) / campaign title, entered at creation',
  description  TEXT            NULL COMMENT '活动详情(可选) / campaign description (optional)',
  theme        VARCHAR(32)     NULL COMMENT '选的主题 key(可选，见 lib/theme.ts) / chosen theme key (optional, see lib/theme.ts)',
  created_at   TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '入库时间 / stored at',
  PRIMARY KEY (campaign_id)
) ENGINE=InnoDB COMMENT='活动标题/详情/主题(API 直写，非索引器) / campaign title, description, theme (API-written, not indexer)';

-- 里程碑计划说明：跟 campaign_metadata 同一类——合约只存金额，链上事件从没
-- 带过"这笔钱打算做什么"。建活动时可选填，跟 milestone_receipts 的 note
-- 字段是一对"计划 vs 实际"：这张表是创建活动时写的"打算做什么"，
-- milestone_receipts.note 是举证时写的"实际花在哪了"，故意分成两张表，
-- 别混在一起。同样不建外键，读的时候 LEFT JOIN。
-- Per-milestone plan description: the same kind of gap as campaign_metadata
-- — the contract only stores an amount, on-chain events never carried "what
-- this money is for". Optional, filled in at campaign creation. Pairs with
-- milestone_receipts.note as "planned vs actual": this table is the plan
-- written at creation time, milestone_receipts.note is what was actually
-- spent, written at receipt time — kept as two tables on purpose, not merged.
-- Same no-FK, LEFT JOIN-on-read pattern as everything else here.
CREATE TABLE IF NOT EXISTS milestone_metadata (
  campaign_id    BIGINT UNSIGNED NOT NULL COMMENT '所属活动 id / owning campaign id',
  milestone_idx  INT UNSIGNED    NOT NULL COMMENT '里程碑序号 / milestone index',
  description    TEXT            NOT NULL COMMENT '这个里程碑打算做什么(可选，前端建活动时填) / what this milestone will accomplish (optional, entered at creation)',
  created_at     TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '入库时间 / stored at',
  PRIMARY KEY (campaign_id, milestone_idx)
) ENGINE=InnoDB COMMENT='里程碑计划说明(API 直写，非索引器) / milestone plan descriptions (API-written, not indexer)';

-- 里程碑支出凭证文件：milestones.receipt_hash 只是链上那个 bytes32 哈希，
-- 从来不是真正的文件——文件本体从没地方存。这张表补上"真正能打开看的凭证"，
-- 存的是本地磁盘上的文件路径，不是文件内容本身（见 config.UploadsDir /
-- api/receipt_upload.go）。
-- sha256 是后端收到文件后自己算的（不信任前端声称的哈希），前端展示时会
-- 拿它和 milestones.receipt_hash（链上值，经 indexer 同步）比对："两个哈希
-- 一致" 才是这份文件真的对应链上那次举证的证明，不是靠后端说了算。
-- 同样不建外键：里程碑行在 CampaignCreated 时就建好了，理论上不会有
-- campaign_metadata 那种"活动还没同步"的时间差，但保持和其它链下表一致的
-- 写法（LEFT JOIN 读取），少一个特例要记。
-- Milestone spending-receipt files: milestones.receipt_hash is only the
-- on-chain bytes32 hash — never an actual file, because there was never
-- anywhere to put one. This table adds the real, openable document, storing
-- a path on local disk rather than the file's bytes (see config.UploadsDir /
-- api/receipt_upload.go).
-- sha256 is computed by the backend itself from the received bytes (never
-- trusting a client-declared hash) — the frontend compares it against
-- milestones.receipt_hash (the on-chain value, synced by the indexer) so
-- "the two hashes match" is a fact anyone can check, not something the
-- backend merely asserts.
-- No FK either: milestone rows exist from CampaignCreated onward, so there's
-- no campaign_metadata-style race in practice, but this keeps the same
-- LEFT JOIN read pattern as every other off-chain table — one fewer special
-- case to remember.
CREATE TABLE IF NOT EXISTS milestone_receipts (
  campaign_id    BIGINT UNSIGNED NOT NULL COMMENT '所属活动 id / owning campaign id',
  milestone_idx  INT UNSIGNED    NOT NULL COMMENT '里程碑序号 / milestone index',
  file_path      VARCHAR(255)    NOT NULL COMMENT '磁盘相对路径(uploads 目录下) / path under the uploads dir',
  file_name      VARCHAR(255)    NOT NULL COMMENT '原始文件名(展示用) / original filename, for display',
  content_type   VARCHAR(100)    NOT NULL COMMENT 'MIME 类型 / MIME type',
  file_size      BIGINT UNSIGNED NOT NULL COMMENT '文件大小(字节) / file size in bytes',
  sha256         CHAR(64)        NOT NULL COMMENT '后端计算的文件哈希(十六进制，无0x) / backend-computed file hash, hex, no 0x prefix',
  note           TEXT            NULL COMMENT '机构填写的支出说明(可选) / charity-written expense note (optional)',
  uploaded_at    TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间 / uploaded at',
  PRIMARY KEY (campaign_id, milestone_idx)
) ENGINE=InnoDB COMMENT='里程碑支出凭证文件(API 直写，非索引器) / milestone spending-receipt files (API-written, not indexer)';
