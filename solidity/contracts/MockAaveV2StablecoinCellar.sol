// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import "./interfaces.sol";

contract MockAaveV2StablecoinCellar is Ownable {
    /**
     * @notice Whether or not the contract is paused in case of an emergency.
     */
    bool public isPaused;

    constructor() {

    }

    event mockAccrueFees();

    /**
     * @notice Take platform fees and performance fees off of cellar's active assets.
     */
    function accrueFees() external onlyOwner {
        emit mockAccrueFees();
    }

    event mockTransferFees();

    /**
     * @notice Transfer accrued fees to Cosmos to distribute.
     */
    function transferFees() external onlyOwner {
        emit mockTransferFees();
    }

    event mockEnterStrategy();

    /**
     * @notice Enters into the current Aave stablecoin strategy.
     */
    function enterStrategy() external onlyOwner  {
        emit mockEnterStrategy();
    }

    event mockRebalance(address[] path, uint256 amountOutMinimum);

    /**
     * @notice Rebalances current assets into a new asset strategy.
     * @param path path to swap from the current asset to new asset using Uniswap
     * @param amountOutMinimum minimum amount of assets returned after swap
     */
    function rebalance(address[] memory path, uint256 amountOutMinimum) external onlyOwner  {
        emit mockRebalance(path, amountOutMinimum);
    }

    event mockReinvest(address[] path, uint256 minAssetsOut);

    /**
     * @notice Reinvest rewards back into cellar's current strategy.
     * @dev Must be called within 2 day unstake period 10 days after `claimAndUnstake` was run.
     * @param path path to swap from AAVE to the current asset on Sushiswap
     * @param minAssetsOut minimum amount of assets cellar should receive after swap
     */
    function reinvest(address[] memory path, uint256 minAssetsOut) public  {
        emit mockReinvest(path, minAssetsOut);
    }

    event mockClaimAndUnstake();

    /**
     * @notice Claim rewards from Aave and begin cooldown period to unstake them.
     * @return claimed amount of rewards claimed from Aave
     */
    function claimAndUnstake() public  returns (uint256 claimed) {
        emit mockClaimAndUnstake();
        return 100;
    }

    event mockSweep(address token);

    /**
     * @notice Sweep tokens sent here that are not managed by the cellar.
     * @dev This may be used in case the wrong tokens are accidentally sent to this contract.
     * @param token address of token to transfer out of this cellar
     */
    function sweep(address token) external onlyOwner  {
        emit mockSweep(token);
    }

    event mockRemoveLiquidityRestriction();

    /**
     * @notice Removes initial liquidity restriction.
     */
    function removeLiquidityRestriction() external onlyOwner  {
        emit mockRemoveLiquidityRestriction();
    }

    event mockPause(bool isPaused);

    /**
     * @notice Pause the contract to prevent deposits.
     * @param _isPaused whether the contract should be paused or unpaused
     */
    function setPause(bool _isPaused) external onlyOwner  {
        isPaused = _isPaused;
        emit mockPause(isPaused);
    }

    event mockShutdown();

    function shutdown() external onlyOwner  {
        emit mockShutdown();
    }
}
