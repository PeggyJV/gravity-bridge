pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Capped.sol";

contract CosmosERC20 is ERC20Burnable, ERC20Capped {
	uint256 MAX_UINT = 2**256 - 1;

	address gravity;

	constructor(
		address _gravityAddress,
		string memory _name,
		string memory _symbol,
		uint8 _decimals,
		uint256 _maxSupply
	) public ERC20(_name, _symbol) ERC20Capped(_maxSupply) {
		_setupDecimals(_decimals);

		gravity = _gravityAddress;
	}

	/**
     * @dev See {IERC20-transfer}.
 	 * @notice Gravity-specific: transfers to the Gravity contract result in burns.
     *
     * Requirements:
     *
     * - `recipient` cannot be the zero address.
     * - the caller must have a balance of at least `amount`.
     */
    function transfer(address recipient, uint256 amount) public virtual override returns (bool) {
		if (recipient == gravity) {
			_burn(_msgSender(), amount);
		} else {
			_transfer(_msgSender(), recipient, amount);
		}

		return true;
    }


    /**
     * @dev See {IERC20-transferFrom}.
     * @notice Gravity-specific: transfers from the Gravity contract result in mints.
	 * Gravity still needs to provide an allowance to the contract.
     *
     * Emits an {Approval} event indicating the updated allowance. This is not
     * required by the EIP. See the note at the beginning of {ERC20}.
     *
     * Requirements:
     *
     * - `sender` and `recipient` cannot be the zero address.
     * - `sender` must have a balance of at least `amount`.
     * - the caller must have allowance for ``sender``'s tokens of at least
     * `amount`.
     */
    function transferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) public virtual override returns (bool) {
		if (sender == gravity) {
			_mint(gravity, amount);
		}

		_transfer(sender, recipient, amount);

		uint256 currentAllowance = _allowances[sender][_msgSender()];
		require(currentAllowance >= amount, "ERC20: transfer amount exceeds allowance");
		unchecked {
			_approve(sender, _msgSender(), currentAllowance - amount);
		}

		return true;
    }


}
