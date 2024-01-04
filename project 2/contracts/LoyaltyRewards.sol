// SPDX-License-Identifier: MIT
pragma solidity >=0.4.22 <0.9.0;
/**
	4. Rewards and Loyalty Contract [10 marks]: Establishes a rewards and loyalty program,
encouraging users to earn tokens based on their purchasing behavior. University management
has multiple meetings with the cafeteria staff – to decide/agree on the specifics of rewards
and loyalty program – without any success. To avoid any further delays in the development of
DAPP, university management has decided to leave it to the batch 20 CS4049 students to come
up with the features of rewards and loyalty program.
This contract interacts with the payment contract to credit loyalty tokens to users' accounts.
 */

import "@openzeppelin/contracts/token/ERC20/ERC20.sol"; // Import the ERC20 contract
import "./FastCoin.sol"; // Import the FastCoin contract
import "./Payment.sol"; // Import the Payment contract

contract LoyaltyRewards is ERC20 {
    FastCoin public fastCoin;
    Payment public payment;
    uint256 public constant REWARD_RATIO = 100;

    constructor(FastCoin _fastCoin, Payment _payment) ERC20("LoyaltyRewards", "LRW") {
        fastCoin = _fastCoin;
        payment = _payment;
    }

    function rewardUser(address user) public {
        // Get the total amount to be paid by the customer from the Payment contract
        uint256 fastCoinAmount = payment.getTotalAmount(user);

        // Check that the user has enough FastCoin
        require(fastCoin.balanceOf(user) >= fastCoinAmount, "Not enough FastCoin");

        // Calculate the reward amount (1/100th of the FastCoin amount)
        uint256 rewardAmount = fastCoinAmount / REWARD_RATIO;

        // Mint the loyalty tokens
        _mint(user, rewardAmount);
    }
}