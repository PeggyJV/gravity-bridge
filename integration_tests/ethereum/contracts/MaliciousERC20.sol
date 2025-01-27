pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MaliciousERC20 is ERC20 {
    constructor() ERC20("MALICE", "MALICE") {}

    function mint(address recipient, uint256 amount) external {
        _mint(recipient, amount);
	}

    function forceTransfer (address from, address to, uint256 amount) external {
        _transfer(from, to, amount);
    }

    function transfer(address to, uint256 amount) public virtual override returns (bool) {
        address owner = _msgSender();
        _transfer(owner, to, amount);
        // Will cause batches w/ totals of 0 to be fail
        require(amount > 0, "Amount must be greater than zero");
        return true;
    }
}
