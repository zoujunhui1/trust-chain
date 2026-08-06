// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script} from "forge-std/Script.sol";
import {console2} from "forge-std/console2.sol";
import {CharityRegistry} from "../src/CharityRegistry.sol";
import {TrustChainEscrow} from "../src/TrustChainEscrow.sol";

/// @title Deploy — 部署脚本
/// @notice Deploys CharityRegistry then TrustChainEscrow (wired to the registry).
///         先部署 CharityRegistry，再部署 TrustChainEscrow（并注入 registry 地址）。
/// @dev The broadcasting account (passed via CLI --account / --private-key)
///      becomes the registry owner.
///      广播账户（通过命令行 --account / --private-key 提供）会成为 registry 的 owner。
contract Deploy is Script {
    function run() external returns (CharityRegistry registry, TrustChainEscrow escrow) {
        // Everything between start/stop is broadcast as real on-chain transactions.
        // start 与 stop 之间的内容会作为真实链上交易广播出去。
        vm.startBroadcast();

        // 1) Access-control contract; deployer (msg.sender) is the initial owner.
        //    准入控制合约；部署者（msg.sender）为初始 owner。
        registry = new CharityRegistry(msg.sender);

        // 2) Core escrow; immutably bound to the registry above.
        //    核心托管合约；不可变地绑定上面的 registry。
        escrow = new TrustChainEscrow(registry);

        vm.stopBroadcast();

        // Print the deployed addresses so they can be recorded / verified.
        // 打印部署地址，便于记录与验证。
        console2.log("CharityRegistry deployed at:", address(registry));
        console2.log("TrustChainEscrow deployed at:", address(escrow));
        console2.log("Owner (deployer):", msg.sender);
    }
}
