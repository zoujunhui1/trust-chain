// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// ReentrancyGuard protects fund-moving functions from re-entrancy attacks.
// ReentrancyGuard 保护转账函数免受重入攻击。
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
// The registry tells us whether a caller is a verified charity.
// registry 用来判断调用者是否为已审核的慈善机构。
import {CharityRegistry} from "./CharityRegistry.sol";

/// @title  TrustChainEscrow
/// @notice Holds all donations and pays each campaign out milestone by
///         milestone. A milestone only unlocks after the previous one's
///         spending receipt has been submitted on-chain.
/// @notice 托管所有捐款，按里程碑逐笔放款；只有上一里程碑的支出凭证上链后，下一里程碑才解锁。
contract TrustChainEscrow is ReentrancyGuard {
    /// @notice Milestone lifecycle: Locked -> Released -> Proven.
    /// @notice 里程碑生命周期：Locked（锁定）-> Released（已放款）-> Proven（已举证）。
    enum MilestoneState {
        Locked,   // funds not yet released / 资金尚未放出
        Released, // funds paid to charity, receipt not yet submitted / 已放款，凭证未交
        Proven    // receipt hash submitted on-chain / 凭证哈希已上链
    }

    /// @notice A single milestone within a campaign.
    /// @notice 活动内的单个里程碑。
    struct Milestone {
        uint256 amount;         // amount this milestone pays out / 本里程碑放款金额
        MilestoneState state;   // current state / 当前状态
        bytes32 receiptHash;    // SHA-256 of the receipt file, set on submit / 收据文件的 SHA-256，提交时写入
    }

    /// @notice A fundraising campaign owned by a verified charity.
    /// @notice 由已审核慈善机构拥有的一个募捐活动。
    struct Campaign {
        address charity;        // creator and payout recipient / 创建者与收款方
        uint256 goal;           // total target = sum of milestone amounts / 目标 = 里程碑金额之和
        uint256 raised;         // total donated so far / 已筹总额
        uint256 released;       // total paid out so far / 已放出总额
        bytes32 metadataHash;   // points to off-chain details (title/desc/images) / 指向链下详情
        bool exists;            // guard for valid campaign id / 校验 campaignId 是否有效
        Milestone[] milestones; // ordered milestones / 有序的里程碑列表
    }

    /// @notice The charity registry used to check verification status.
    /// @notice 用于查询审核状态的慈善机构注册表。
    CharityRegistry public immutable registry;

    /// @notice campaignId => campaign data.
    /// @notice 活动 id => 活动数据。
    mapping(uint256 => Campaign) private _campaigns;

    /// @notice campaignId => donor => total amount donated (in wei).
    /// @notice 活动 id => 捐款人 => 累计捐款额（单位 wei）。
    mapping(uint256 => mapping(address => uint256)) public donations;

    /// @notice Number of campaigns created so far; also the next campaign id.
    /// @notice 已创建的活动数量；同时也是下一个活动的 id。
    uint256 public campaignCount;

    /// @notice Emitted when a verified charity creates a campaign.
    /// @notice 已审核机构创建活动时触发。
    event CampaignCreated(
        uint256 indexed campaignId,
        address indexed charity,
        uint256 goal,
        uint256 milestoneCount
    );

    /// @notice Emitted on every donation.
    /// @notice 每次捐款时触发。
    event DonationReceived(
        uint256 indexed campaignId,
        address indexed donor,
        uint256 amount
    );

    /// @notice Emitted when a milestone is released and paid out.
    /// @notice 里程碑放款时触发。
    event MilestoneReleased(
        uint256 indexed campaignId,
        uint256 indexed milestoneIndex,
        uint256 amount
    );

    /// @notice Emitted when a receipt hash is submitted for a milestone.
    /// @notice 里程碑的收据哈希提交上链时触发。
    event ReceiptSubmitted(
        uint256 indexed campaignId,
        uint256 indexed milestoneIndex,
        bytes32 receiptHash
    );

    /// @notice Emitted when the last milestone is proven and the campaign ends.
    /// @notice 最后一个里程碑举证完成、活动结束时触发。
    event CampaignCompleted(uint256 indexed campaignId);

    /// @notice Store the registry address at deployment (cannot change later).
    /// @notice 部署时记录 registry 地址（之后不可更改）。
    /// @param  registry_ The deployed CharityRegistry contract.
    /// @param  registry_ 已部署的 CharityRegistry 合约。
    constructor(CharityRegistry registry_) {
        registry = registry_;
    }

    /// @notice Create a campaign. Caller must be a verified charity.
    ///         The goal is the sum of all milestone amounts.
    /// @notice 创建活动。调用者必须是已审核机构；目标额为所有里程碑金额之和。
    /// @param  milestoneAmounts Ordered payout amounts, one per milestone (wei).
    /// @param  milestoneAmounts 每个里程碑的放款金额（wei），按顺序排列。
    /// @param  metadataHash Hash pointing to off-chain campaign details.
    /// @param  metadataHash 指向链下活动详情的哈希。
    /// @return campaignId The id of the newly created campaign / 新建活动的 id。
    function createCampaign(uint256[] calldata milestoneAmounts, bytes32 metadataHash)
        external
        returns (uint256 campaignId)
    {
        require(registry.isVerified(msg.sender), "TrustChainEscrow: not a verified charity");
        require(milestoneAmounts.length > 0, "TrustChainEscrow: no milestones");

        campaignId = campaignCount;
        campaignCount++;

        Campaign storage c = _campaigns[campaignId];
        c.charity = msg.sender;
        c.metadataHash = metadataHash;
        c.exists = true;

        uint256 goal;
        for (uint256 i = 0; i < milestoneAmounts.length; i++) {
            require(milestoneAmounts[i] > 0, "TrustChainEscrow: zero milestone amount");
            goal += milestoneAmounts[i];
            c.milestones.push(
                Milestone({amount: milestoneAmounts[i], state: MilestoneState.Locked, receiptHash: bytes32(0)})
            );
        }
        c.goal = goal;

        emit CampaignCreated(campaignId, msg.sender, goal, milestoneAmounts.length);
    }

    /// @notice Donate ETH to a campaign. The contract holds the funds.
    /// @notice 向活动捐款（ETH）；资金由本合约托管。
    /// @param  campaignId The campaign to donate to.
    /// @param  campaignId 要捐款的活动 id。
    function donate(uint256 campaignId) external payable {
        Campaign storage c = _campaigns[campaignId];
        require(c.exists, "TrustChainEscrow: campaign does not exist");
        require(msg.value > 0, "TrustChainEscrow: zero donation");

        c.raised += msg.value;
        donations[campaignId][msg.sender] += msg.value;

        emit DonationReceived(campaignId, msg.sender, msg.value);
    }

    /// @notice Release a milestone's funds to the charity. Charity only.
    ///         Requires the goal to be reached, the milestone to be Locked,
    ///         and the previous milestone (if any) to be Proven.
    /// @notice 将某里程碑的资金放给机构。仅限该活动机构调用。
    ///         要求目标已达成、该里程碑为 Locked，且上一里程碑（若有）已 Proven。
    /// @param  campaignId The campaign id.
    /// @param  campaignId 活动 id。
    /// @param  milestoneIndex The milestone to release.
    /// @param  milestoneIndex 要放款的里程碑序号。
    function releaseMilestone(uint256 campaignId, uint256 milestoneIndex) external nonReentrant {
        Campaign storage c = _campaigns[campaignId];
        require(c.exists, "TrustChainEscrow: campaign does not exist");
        require(msg.sender == c.charity, "TrustChainEscrow: not the campaign charity");
        require(c.raised >= c.goal, "TrustChainEscrow: goal not reached");
        require(milestoneIndex < c.milestones.length, "TrustChainEscrow: invalid milestone");

        Milestone storage m = c.milestones[milestoneIndex];
        require(m.state == MilestoneState.Locked, "TrustChainEscrow: milestone not locked");

        // Evidence chain: milestone i stays locked until milestone i-1 is Proven.
        // 证据链：里程碑 i 在 i-1 变为 Proven 之前保持锁定。
        if (milestoneIndex > 0) {
            require(
                c.milestones[milestoneIndex - 1].state == MilestoneState.Proven,
                "TrustChainEscrow: previous milestone not proven"
            );
        }

        // Effects: update state BEFORE sending funds (checks-effects-interactions).
        // 生效：先改状态，再转账（检查-生效-交互）。
        m.state = MilestoneState.Released;
        c.released += m.amount;

        // Interaction: send the funds last.
        // 交互：最后才转账。
        (bool ok, ) = c.charity.call{value: m.amount}("");
        require(ok, "TrustChainEscrow: transfer failed");

        emit MilestoneReleased(campaignId, milestoneIndex, m.amount);
    }

    /// @notice Submit the spending receipt hash for a released milestone,
    ///         moving it to Proven and unlocking the next milestone. Charity only.
    /// @notice 为已放款的里程碑提交支出收据哈希，使其变为 Proven 并解锁下一里程碑。仅限该活动机构。
    /// @param  campaignId The campaign id.
    /// @param  campaignId 活动 id。
    /// @param  milestoneIndex The milestone to prove.
    /// @param  milestoneIndex 要举证的里程碑序号。
    /// @param  receiptHash SHA-256 of the receipt file, computed off-chain.
    /// @param  receiptHash 收据文件的 SHA-256（链下计算）。
    function submitReceipt(uint256 campaignId, uint256 milestoneIndex, bytes32 receiptHash) external {
        require(receiptHash != bytes32(0), "TrustChainEscrow: empty receipt hash");

        Campaign storage c = _campaigns[campaignId];
        require(c.exists, "TrustChainEscrow: campaign does not exist");
        require(msg.sender == c.charity, "TrustChainEscrow: not the campaign charity");
        require(milestoneIndex < c.milestones.length, "TrustChainEscrow: invalid milestone");

        Milestone storage m = c.milestones[milestoneIndex];
        require(m.state == MilestoneState.Released, "TrustChainEscrow: milestone not released");

        m.receiptHash = receiptHash;
        m.state = MilestoneState.Proven;

        emit ReceiptSubmitted(campaignId, milestoneIndex, receiptHash);

        // Milestones are proven in order, so proving the last one completes the campaign.
        // 里程碑按顺序举证，故最后一个被举证即代表活动完成。
        if (milestoneIndex == c.milestones.length - 1) {
            emit CampaignCompleted(campaignId);
        }
    }

    /// @notice Read a campaign's summary. Milestones are fetched separately.
    /// @notice 读取活动摘要；里程碑单独获取。
    /// @param  campaignId The campaign id.
    /// @param  campaignId 活动 id。
    function getCampaign(uint256 campaignId)
        external
        view
        returns (
            address charity,
            uint256 goal,
            uint256 raised,
            uint256 released,
            bytes32 metadataHash,
            bool exists,
            uint256 milestoneCount
        )
    {
        Campaign storage c = _campaigns[campaignId];
        return (c.charity, c.goal, c.raised, c.released, c.metadataHash, c.exists, c.milestones.length);
    }

    /// @notice Read all milestones of a campaign.
    /// @notice 读取活动的全部里程碑。
    /// @param  campaignId The campaign id.
    /// @param  campaignId 活动 id。
    function getMilestones(uint256 campaignId) external view returns (Milestone[] memory) {
        return _campaigns[campaignId].milestones;
    }
}
