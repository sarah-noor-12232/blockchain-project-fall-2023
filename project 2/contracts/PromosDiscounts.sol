// SPDX-License-Identifier: MIT
pragma solidity >=0.4.22 <0.9.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./OrderProcessing.sol";

contract PromosDiscounts is ERC20 {
    OrderProcessing public orderProcessing;

    // Mapping from item ID to discount percentage
    mapping(uint256 => uint256) private discounts;

    constructor(OrderProcessing _orderProcessing) ERC20("PromosDiscounts", "PRD") {
        orderProcessing = _orderProcessing;
    }

    // Function to add a discount
    function addDiscount(uint256 itemId, uint256 discount) public {
        require(discount <= 100, "Discount cannot be more than 100%");
        discounts[itemId] = discount;
    }

    // Function to remove a discount
    function removeDiscount(uint256 itemId) public {
        delete discounts[itemId];
    }

    // Function to calculate the adjusted order amount
    function calculateAdjustedAmount(uint256 orderId) public view returns (uint256) {
        uint256 orderAmount = orderProcessing.getOrderAmount(orderId);
        uint256 itemId = orderProcessing.getItemId(orderId);
        uint256 discount = discounts[itemId];
        uint256 discountAmount = orderAmount * discount / 100;
        return orderAmount - discountAmount;
    }
}