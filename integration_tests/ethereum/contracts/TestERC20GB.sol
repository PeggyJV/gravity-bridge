pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract TestERC20GB is ERC20 {
    constructor() public ERC20("DollaryDoos", "DDS") {
        // the four test validators
        _mint(0xd312f0f1B39D54Db2829537595fC1167B14d4b34, 10000);
        _mint(0x7bE2a04df4b9C3227928147461e19158eB2B11d1, 10000);
        _mint(0xb8c6886FDDa38adaa0F416722dd5554886C43055, 10000);
        _mint(0x14fdAC734De10065093C4Ed4a83C41638378005A, 10000);
        // this is the EtherBase address for our testnet miner in
        // tests/assets/ETHGenesis.json so it wil have both a lot
        // of ETH and a lot of erc20 tokens to test with
        _mint(0xBf660843528035a5A4921534E156a27e64B231fE, 100000000000000000000000000);
    }
}
