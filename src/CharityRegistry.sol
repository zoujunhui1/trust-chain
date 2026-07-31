// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// Import OpenZeppelin's two-step ownership base contract.
// 导入 OpenZeppelin 的"两步转移所有权"基础合约。
import {Ownable2Step, Ownable} from "@openzeppelin/contracts/access/Ownable2Step.sol";

/// @title  CharityRegistry
/// @notice The platform admin (owner) verifies or revokes charity addresses.
///         Only verified addresses may create campaigns in the escrow contract.
/// @notice 平台管理员(owner)审核或撤销慈善机构地址；只有已审核地址才能在托管合约中创建活动。
contract CharityRegistry is Ownable2Step {
    /// @notice Whether an address is a verified charity. Defaults to false.
    /// @notice 某地址是否为已审核的慈善机构；默认 false（未审核）。
    mapping(address => bool) private _verified;

    /// @notice Emitted when the admin verifies a charity.
    /// @notice 管理员审核通过一个慈善机构时触发。
    /// @param  charity The address that was verified.
    /// @param  charity 被审核通过的地址。
    event CharityVerified(address indexed charity);

    /// @notice Emitted when the admin revokes a charity's verification.
    /// @notice 管理员撤销某慈善机构的审核资格时触发。
    /// @param  charity The address whose verification was revoked.
    /// @param  charity 被撤销资格的地址。
    event CharityRevoked(address indexed charity);

    /// @notice Set the initial owner (the platform admin) at deployment.
    /// @notice 部署时设置初始 owner（平台管理员）。
    /// @param  initialOwner The address that will be the platform admin.
    /// @param  initialOwner 将成为平台管理员的地址。
    constructor(address initialOwner) Ownable(initialOwner) {}

    /// @notice Verify a charity so it can create campaigns. Owner only.
    /// @notice 审核通过一个慈善机构，使其可创建活动。仅限 owner 调用。
    /// @param  charity The charity address to verify.
    /// @param  charity 要审核的慈善机构地址。
    function verifyCharity(address charity) external onlyOwner {
        require(charity != address(0), "CharityRegistry: zero address");
        require(!_verified[charity], "CharityRegistry: already verified");
        _verified[charity] = true;
        emit CharityVerified(charity);
    }

    /// @notice Revoke a charity's verification. Owner only.
    ///         Existing campaigns are unaffected; the charity simply cannot
    ///         create new ones after this.
    /// @notice 撤销某慈善机构的审核资格。仅限 owner 调用。
    ///         已存在的活动不受影响，机构此后只是无法再创建新活动。
    /// @param  charity The charity address to revoke.
    /// @param  charity 要撤销资格的慈善机构地址。
    function revokeCharity(address charity) external onlyOwner {
        require(_verified[charity], "CharityRegistry: not verified");
        _verified[charity] = false;
        emit CharityRevoked(charity);
    }

    /// @notice Returns true if `charity` is currently a verified charity.
    /// @notice 若 `charity` 当前为已审核的慈善机构则返回 true。
    function isVerified(address charity) external view returns (bool) {
        return _verified[charity];
    }
}
