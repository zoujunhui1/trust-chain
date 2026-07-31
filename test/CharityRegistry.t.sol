// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {CharityRegistry} from "../src/CharityRegistry.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/// @title  CharityRegistry unit tests
/// @notice Covers the permission paths, state changes, events, and the
///         two-step ownership transfer inherited from Ownable2Step.
/// @notice 覆盖权限路径、状态变化、事件，以及继承自 Ownable2Step 的两步所有权转移。
contract CharityRegistryTest is Test {
    CharityRegistry internal registry;

    address internal owner = makeAddr("owner");
    address internal charity = makeAddr("charity");
    address internal attacker = makeAddr("attacker");

    // Re-declared here so we can assert on them with vm.expectEmit.
    // 在测试里重新声明，以便用 vm.expectEmit 断言。
    event CharityVerified(address indexed charity);
    event CharityRevoked(address indexed charity);

    function setUp() public {
        registry = new CharityRegistry(owner);
    }

    // --------------------------------------------------------------------
    // Deployment / 部署
    // --------------------------------------------------------------------

    function test_Constructor_SetsOwner() public view {
        assertEq(registry.owner(), owner);
    }

    function test_UnknownAddress_IsNotVerified() public view {
        assertFalse(registry.isVerified(charity));
    }

    // --------------------------------------------------------------------
    // verifyCharity / 审核
    // --------------------------------------------------------------------

    function test_VerifyCharity_ByOwner_SetsVerified() public {
        vm.prank(owner);
        registry.verifyCharity(charity);
        assertTrue(registry.isVerified(charity));
    }

    function test_VerifyCharity_EmitsEvent() public {
        vm.expectEmit(true, false, false, true);
        emit CharityVerified(charity);

        vm.prank(owner);
        registry.verifyCharity(charity);
    }

    function test_VerifyCharity_RevertWhen_NotOwner() public {
        vm.prank(attacker);
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, attacker)
        );
        registry.verifyCharity(charity);
    }

    function test_VerifyCharity_RevertWhen_ZeroAddress() public {
        vm.prank(owner);
        vm.expectRevert("CharityRegistry: zero address");
        registry.verifyCharity(address(0));
    }

    function test_VerifyCharity_RevertWhen_AlreadyVerified() public {
        vm.startPrank(owner);
        registry.verifyCharity(charity);
        vm.expectRevert("CharityRegistry: already verified");
        registry.verifyCharity(charity);
        vm.stopPrank();
    }

    // --------------------------------------------------------------------
    // revokeCharity / 撤销
    // --------------------------------------------------------------------

    function test_RevokeCharity_ByOwner_ClearsVerified() public {
        vm.startPrank(owner);
        registry.verifyCharity(charity);
        registry.revokeCharity(charity);
        vm.stopPrank();
        assertFalse(registry.isVerified(charity));
    }

    function test_RevokeCharity_EmitsEvent() public {
        vm.prank(owner);
        registry.verifyCharity(charity);

        vm.expectEmit(true, false, false, true);
        emit CharityRevoked(charity);

        vm.prank(owner);
        registry.revokeCharity(charity);
    }

    function test_RevokeCharity_RevertWhen_NotOwner() public {
        vm.prank(owner);
        registry.verifyCharity(charity);

        vm.prank(attacker);
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, attacker)
        );
        registry.revokeCharity(charity);
    }

    function test_RevokeCharity_RevertWhen_NotVerified() public {
        vm.prank(owner);
        vm.expectRevert("CharityRegistry: not verified");
        registry.revokeCharity(charity);
    }

    // --------------------------------------------------------------------
    // Two-step ownership (Ownable2Step) / 两步所有权转移
    // --------------------------------------------------------------------

    function test_TransferOwnership_IsTwoStep() public {
        address newOwner = makeAddr("newOwner");

        // Step 1: current owner nominates; ownership does NOT change yet.
        // 第一步：当前 owner 提名；所有权尚未改变。
        vm.prank(owner);
        registry.transferOwnership(newOwner);
        assertEq(registry.owner(), owner);
        assertEq(registry.pendingOwner(), newOwner);

        // Step 2: the nominee accepts; now ownership changes.
        // 第二步：被提名者确认；此时所有权才转移。
        vm.prank(newOwner);
        registry.acceptOwnership();
        assertEq(registry.owner(), newOwner);
        assertEq(registry.pendingOwner(), address(0));
    }

    function test_AcceptOwnership_RevertWhen_NotPendingOwner() public {
        address newOwner = makeAddr("newOwner");

        vm.prank(owner);
        registry.transferOwnership(newOwner);

        // A random address cannot accept the pending ownership.
        // 无关地址不能接受待定的所有权。
        vm.prank(attacker);
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, attacker)
        );
        registry.acceptOwnership();
    }

    function test_NewOwner_CanVerify_AfterTransfer() public {
        address newOwner = makeAddr("newOwner");
        vm.prank(owner);
        registry.transferOwnership(newOwner);
        vm.prank(newOwner);
        registry.acceptOwnership();

        // Old owner can no longer verify.
        // 旧 owner 不能再审核。
        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, owner)
        );
        registry.verifyCharity(charity);

        // New owner can.
        // 新 owner 可以。
        vm.prank(newOwner);
        registry.verifyCharity(charity);
        assertTrue(registry.isVerified(charity));
    }

    // --------------------------------------------------------------------
    // Fuzz / 模糊测试
    // --------------------------------------------------------------------

    function testFuzz_VerifyThenRevoke(address a) public {
        vm.assume(a != address(0));

        vm.startPrank(owner);
        registry.verifyCharity(a);
        assertTrue(registry.isVerified(a));
        registry.revokeCharity(a);
        assertFalse(registry.isVerified(a));
        vm.stopPrank();
    }
}
