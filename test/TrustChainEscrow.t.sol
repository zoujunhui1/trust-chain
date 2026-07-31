// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {TrustChainEscrow} from "../src/TrustChainEscrow.sol";
import {CharityRegistry} from "../src/CharityRegistry.sol";

/// @title  TrustChainEscrow unit tests
/// @notice Covers campaign creation, donations, milestone release, receipt
///         submission, the evidence-chain ordering rule, and reentrancy safety.
/// @notice 覆盖建活动、捐款、里程碑放款、提交凭证、证据链顺序规则、以及重入安全。
contract TrustChainEscrowTest is Test {
    CharityRegistry internal registry;
    TrustChainEscrow internal escrow;

    address internal owner = makeAddr("owner");
    address internal charity = makeAddr("charity");
    address internal donor = makeAddr("donor");
    address internal donor2 = makeAddr("donor2");
    address internal attacker = makeAddr("attacker");

    bytes32 internal constant META = keccak256("ipfs://campaign-details");
    bytes32 internal constant RECEIPT = keccak256("receipt-file-contents");

    // Re-declared for vm.expectEmit. / 为 vm.expectEmit 重新声明。
    event CampaignCreated(uint256 indexed campaignId, address indexed charity, uint256 goal, uint256 milestoneCount);
    event DonationReceived(uint256 indexed campaignId, address indexed donor, uint256 amount);
    event MilestoneReleased(uint256 indexed campaignId, uint256 indexed milestoneIndex, uint256 amount);
    event ReceiptSubmitted(uint256 indexed campaignId, uint256 indexed milestoneIndex, bytes32 receiptHash);
    event CampaignCompleted(uint256 indexed campaignId);

    function setUp() public {
        registry = new CharityRegistry(owner);
        escrow = new TrustChainEscrow(registry);

        // Verify the charity so it can create campaigns.
        // 审核该机构，使其可创建活动。
        vm.prank(owner);
        registry.verifyCharity(charity);
    }

    // --------------------------------------------------------------------
    // Helpers / 辅助函数
    // --------------------------------------------------------------------

    /// Two milestones: 1 ETH then 2 ETH, goal = 3 ETH.
    function _amounts() internal pure returns (uint256[] memory a) {
        a = new uint256[](2);
        a[0] = 1 ether;
        a[1] = 2 ether;
    }

    function _createCampaign() internal returns (uint256 id) {
        vm.prank(charity);
        id = escrow.createCampaign(_amounts(), META);
    }

    function _donate(uint256 id, address from, uint256 amount) internal {
        vm.deal(from, amount);
        vm.prank(from);
        escrow.donate{value: amount}(id);
    }

    // --------------------------------------------------------------------
    // createCampaign / 建活动
    // --------------------------------------------------------------------

    function test_CreateCampaign_StoresData() public {
        uint256 id = _createCampaign();

        (
            address c,
            uint256 goal,
            uint256 raised,
            uint256 released,
            bytes32 meta,
            bool exists,
            uint256 mc
        ) = escrow.getCampaign(id);

        assertEq(id, 0);
        assertEq(c, charity);
        assertEq(goal, 3 ether);
        assertEq(raised, 0);
        assertEq(released, 0);
        assertEq(meta, META);
        assertTrue(exists);
        assertEq(mc, 2);
        assertEq(escrow.campaignCount(), 1);

        TrustChainEscrow.Milestone[] memory ms = escrow.getMilestones(id);
        assertEq(ms[0].amount, 1 ether);
        assertEq(ms[1].amount, 2 ether);
        assertEq(uint256(ms[0].state), uint256(TrustChainEscrow.MilestoneState.Locked));
        assertEq(uint256(ms[1].state), uint256(TrustChainEscrow.MilestoneState.Locked));
        assertEq(ms[0].receiptHash, bytes32(0));
    }

    function test_CreateCampaign_EmitsEvent() public {
        vm.expectEmit(true, true, false, true);
        emit CampaignCreated(0, charity, 3 ether, 2);
        vm.prank(charity);
        escrow.createCampaign(_amounts(), META);
    }

    function test_CreateCampaign_IncrementsId() public {
        uint256 first = _createCampaign();
        uint256 second = _createCampaign();
        assertEq(first, 0);
        assertEq(second, 1);
        assertEq(escrow.campaignCount(), 2);
    }

    function test_CreateCampaign_RevertWhen_NotVerified() public {
        vm.prank(attacker);
        vm.expectRevert("TrustChainEscrow: not a verified charity");
        escrow.createCampaign(_amounts(), META);
    }

    function test_CreateCampaign_RevertWhen_NoMilestones() public {
        uint256[] memory empty = new uint256[](0);
        vm.prank(charity);
        vm.expectRevert("TrustChainEscrow: no milestones");
        escrow.createCampaign(empty, META);
    }

    function test_CreateCampaign_RevertWhen_ZeroMilestoneAmount() public {
        uint256[] memory a = new uint256[](2);
        a[0] = 1 ether;
        a[1] = 0;
        vm.prank(charity);
        vm.expectRevert("TrustChainEscrow: zero milestone amount");
        escrow.createCampaign(a, META);
    }

    // --------------------------------------------------------------------
    // donate / 捐款
    // --------------------------------------------------------------------

    function test_Donate_RecordsFunds() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 1 ether);

        (, , uint256 raised, , , , ) = escrow.getCampaign(id);
        assertEq(raised, 1 ether);
        assertEq(escrow.donations(id, donor), 1 ether);
        assertEq(address(escrow).balance, 1 ether);
    }

    function test_Donate_MultipleDonors_Accumulate() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 1 ether);
        _donate(id, donor2, 2 ether);
        _donate(id, donor, 0.5 ether);

        (, , uint256 raised, , , , ) = escrow.getCampaign(id);
        assertEq(raised, 3.5 ether);
        assertEq(escrow.donations(id, donor), 1.5 ether);
        assertEq(escrow.donations(id, donor2), 2 ether);
    }

    function test_Donate_EmitsEvent() public {
        uint256 id = _createCampaign();
        vm.deal(donor, 1 ether);
        vm.expectEmit(true, true, false, true);
        emit DonationReceived(id, donor, 1 ether);
        vm.prank(donor);
        escrow.donate{value: 1 ether}(id);
    }

    function test_Donate_RevertWhen_CampaignDoesNotExist() public {
        vm.deal(donor, 1 ether);
        vm.prank(donor);
        vm.expectRevert("TrustChainEscrow: campaign does not exist");
        escrow.donate{value: 1 ether}(999);
    }

    function test_Donate_RevertWhen_ZeroValue() public {
        uint256 id = _createCampaign();
        vm.prank(donor);
        vm.expectRevert("TrustChainEscrow: zero donation");
        escrow.donate{value: 0}(id);
    }

    // --------------------------------------------------------------------
    // releaseMilestone / 放款
    // --------------------------------------------------------------------

    function test_ReleaseMilestone_PaysCharity() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 3 ether);

        uint256 before = charity.balance;
        vm.prank(charity);
        escrow.releaseMilestone(id, 0);

        assertEq(charity.balance, before + 1 ether);
        (, , , uint256 released, , , ) = escrow.getCampaign(id);
        assertEq(released, 1 ether);

        TrustChainEscrow.Milestone[] memory ms = escrow.getMilestones(id);
        assertEq(uint256(ms[0].state), uint256(TrustChainEscrow.MilestoneState.Released));
    }

    function test_ReleaseMilestone_EmitsEvent() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 3 ether);
        vm.expectEmit(true, true, false, true);
        emit MilestoneReleased(id, 0, 1 ether);
        vm.prank(charity);
        escrow.releaseMilestone(id, 0);
    }

    function test_ReleaseMilestone_RevertWhen_NotCharity() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 3 ether);
        vm.prank(attacker);
        vm.expectRevert("TrustChainEscrow: not the campaign charity");
        escrow.releaseMilestone(id, 0);
    }

    function test_ReleaseMilestone_RevertWhen_GoalNotReached() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 1 ether); // below 3 ETH goal
        vm.prank(charity);
        vm.expectRevert("TrustChainEscrow: goal not reached");
        escrow.releaseMilestone(id, 0);
    }

    function test_ReleaseMilestone_RevertWhen_InvalidIndex() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 3 ether);
        vm.prank(charity);
        vm.expectRevert("TrustChainEscrow: invalid milestone");
        escrow.releaseMilestone(id, 5);
    }

    function test_ReleaseMilestone_RevertWhen_AlreadyReleased() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 3 ether);
        vm.startPrank(charity);
        escrow.releaseMilestone(id, 0);
        vm.expectRevert("TrustChainEscrow: milestone not locked");
        escrow.releaseMilestone(id, 0);
        vm.stopPrank();
    }

    function test_ReleaseMilestone_RevertWhen_PreviousNotProven() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 3 ether);
        vm.startPrank(charity);
        escrow.releaseMilestone(id, 0); // released but not proven
        vm.expectRevert("TrustChainEscrow: previous milestone not proven");
        escrow.releaseMilestone(id, 1);
        vm.stopPrank();
    }

    // --------------------------------------------------------------------
    // submitReceipt / 提交凭证
    // --------------------------------------------------------------------

    function test_SubmitReceipt_ProvesMilestone() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 3 ether);
        vm.startPrank(charity);
        escrow.releaseMilestone(id, 0);

        vm.expectEmit(true, true, false, true);
        emit ReceiptSubmitted(id, 0, RECEIPT);
        escrow.submitReceipt(id, 0, RECEIPT);
        vm.stopPrank();

        TrustChainEscrow.Milestone[] memory ms = escrow.getMilestones(id);
        assertEq(uint256(ms[0].state), uint256(TrustChainEscrow.MilestoneState.Proven));
        assertEq(ms[0].receiptHash, RECEIPT);
    }

    function test_SubmitReceipt_UnlocksNextMilestone() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 3 ether);
        uint256 before = charity.balance;
        vm.startPrank(charity);
        escrow.releaseMilestone(id, 0);
        escrow.submitReceipt(id, 0, RECEIPT);
        escrow.releaseMilestone(id, 1); // now allowed
        vm.stopPrank();
        assertEq(charity.balance, before + 3 ether);
    }

    function test_SubmitReceipt_LastMilestone_CompletesCampaign() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 3 ether);
        vm.startPrank(charity);
        escrow.releaseMilestone(id, 0);
        escrow.submitReceipt(id, 0, RECEIPT);
        escrow.releaseMilestone(id, 1);

        vm.expectEmit(true, false, false, false);
        emit CampaignCompleted(id);
        escrow.submitReceipt(id, 1, RECEIPT);
        vm.stopPrank();
    }

    function test_SubmitReceipt_RevertWhen_EmptyHash() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 3 ether);
        vm.startPrank(charity);
        escrow.releaseMilestone(id, 0);
        vm.expectRevert("TrustChainEscrow: empty receipt hash");
        escrow.submitReceipt(id, 0, bytes32(0));
        vm.stopPrank();
    }

    function test_SubmitReceipt_RevertWhen_NotCharity() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 3 ether);
        vm.prank(charity);
        escrow.releaseMilestone(id, 0);
        vm.prank(attacker);
        vm.expectRevert("TrustChainEscrow: not the campaign charity");
        escrow.submitReceipt(id, 0, RECEIPT);
    }

    function test_SubmitReceipt_RevertWhen_NotReleased() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 3 ether);
        vm.prank(charity);
        vm.expectRevert("TrustChainEscrow: milestone not released");
        escrow.submitReceipt(id, 0, RECEIPT); // milestone 0 still Locked
    }

    // --------------------------------------------------------------------
    // Full lifecycle / 完整生命周期
    // --------------------------------------------------------------------

    function test_FullLifecycle() public {
        uint256 id = _createCampaign();
        _donate(id, donor, 3 ether);

        vm.startPrank(charity);
        escrow.releaseMilestone(id, 0);
        escrow.submitReceipt(id, 0, RECEIPT);
        escrow.releaseMilestone(id, 1);
        escrow.submitReceipt(id, 1, RECEIPT);
        vm.stopPrank();

        (, uint256 goal, uint256 raised, uint256 released, , , ) = escrow.getCampaign(id);
        assertEq(goal, 3 ether);
        assertEq(raised, 3 ether);
        assertEq(released, 3 ether);
        assertEq(address(escrow).balance, 0); // all funds paid out
        assertEq(charity.balance, 3 ether);

        TrustChainEscrow.Milestone[] memory ms = escrow.getMilestones(id);
        assertEq(uint256(ms[0].state), uint256(TrustChainEscrow.MilestoneState.Proven));
        assertEq(uint256(ms[1].state), uint256(TrustChainEscrow.MilestoneState.Proven));
    }

    // --------------------------------------------------------------------
    // Reentrancy / 重入攻击
    // --------------------------------------------------------------------

    function test_ReleaseMilestone_ReentrancyIsBlocked() public {
        // Deploy a malicious charity that re-enters on receiving funds.
        // 部署一个恶意机构，在收款时尝试重入。
        ReentrantCharity malicious = new ReentrantCharity(escrow);
        vm.prank(owner);
        registry.verifyCharity(address(malicious));

        uint256 id = malicious.createCampaign(_amounts(), META);
        _donate(id, donor, 3 ether);

        malicious.setAttack(true);
        // The reentrant call is blocked; the outer transfer fails and reverts.
        // 重入被拦截，外层转账失败并回滚。
        vm.expectRevert("TrustChainEscrow: transfer failed");
        malicious.release(id, 0);

        // No funds moved: escrow still holds everything.
        // 资金未被转移：escrow 仍持有全部资金。
        assertEq(address(escrow).balance, 3 ether);
    }
}

/// @dev Malicious charity contract used to test the reentrancy guard.
/// @dev 用于测试重入保护的恶意机构合约。
contract ReentrantCharity {
    TrustChainEscrow internal immutable escrow;
    bool internal attack;

    constructor(TrustChainEscrow escrow_) {
        escrow = escrow_;
    }

    function createCampaign(uint256[] calldata amounts, bytes32 meta) external returns (uint256) {
        return escrow.createCampaign(amounts, meta);
    }

    function setAttack(bool a) external {
        attack = a;
    }

    function release(uint256 id, uint256 index) external {
        escrow.releaseMilestone(id, index);
    }

    receive() external payable {
        if (attack) {
            // Try to re-enter and drain a second time.
            // 尝试重入以再取一次。
            escrow.releaseMilestone(0, 0);
        }
    }
}
