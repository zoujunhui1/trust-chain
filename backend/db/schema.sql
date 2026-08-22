-- TrustChain 索引器数据库 schema / Indexer database schema
-- 用法 / usage:  mysql -u root < backend/db/schema.sql
--
-- 设计要点 / design notes:
--   * 这些表由链上事件重建，是前端「读」的唯一来源（写走链上 MetaMask）。
--     These tables are rebuilt from on-chain events and are the sole read source
--     for the frontend (writes happen on-chain via MetaMask).
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
CREATE TABLE IF NOT EXISTS campaigns (
  id               BIGINT UNSIGNED NOT NULL COMMENT '链上 campaignId / on-chain campaign id',
  charity          CHAR(42)        NOT NULL COMMENT '发起机构地址 / charity address',
  goal             DECIMAL(65,0)   NOT NULL COMMENT '目标总额(wei)=各里程碑之和 / goal in wei',
  raised           DECIMAL(65,0)   NOT NULL DEFAULT 0 COMMENT '已筹(wei)=SUM(donations) / total raised',
  released         DECIMAL(65,0)   NOT NULL DEFAULT 0 COMMENT '已放款(wei)=已释放里程碑之和 / total released',
  milestone_count  INT UNSIGNED    NOT NULL COMMENT '里程碑数量 / number of milestones',
  metadata_hash    CHAR(66)        NULL COMMENT 'IPFS 元数据哈希(view 补齐) / ipfs metadata hash',
  completed        BOOLEAN         NOT NULL DEFAULT FALSE COMMENT '是否完成(收到末个收据) / completed flag',
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
