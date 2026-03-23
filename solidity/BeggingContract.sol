
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BeggingContract {
    // 合约所有者地址
    address public owner;

    // 记录每个捐赠者的捐赠总额
    mapping(address => uint256) public totalDonationsByUser;

    // 记录每次捐赠的详细信息
    mapping(address => Donation[]) public userDonations;

    // 捐赠者名单（用于排行榜）
    address[] public donorsList;

    // 捐赠时间限制
    uint256 public startTime;
    uint256 public endTime;

    // 捐赠总数统计
    uint256 public totalDonations;

    // 捐赠者人数统计
    uint256 public donorsCount;

    // 捐赠结构体
    struct Donation {
        uint256 amount;
        uint256 timestamp;
    }

    // 排行榜条目结构体
    struct LeaderboardEntry {
        address donor;
        uint256 amount;
    }

    // 捐赠事件
    event Donated(address indexed donor, uint256 amount, uint256 timestamp);

    // 提款事件
    event Withdrawn(address indexed owner, uint256 amount, uint256 timestamp);

    // 设置时间限制事件
    event TimeLimitSet(uint256 startTime, uint256 endTime);

    // 构造函数，初始化合约所有者和时间限制
    constructor(uint256 _startTime, uint256 _endTime) {
        owner = msg.sender;
        require(_endTime > _startTime, "End time must be after start time");
        startTime = _startTime;
        endTime = _endTime;
    }

    // 仅所有者修饰符
    modifier onlyOwner() {
        require(msg.sender == owner, "Only contract owner can call this function");
        _;
    }

    // 时间限制修饰符
    modifier withinTimeLimit() {
        require(block.timestamp >= startTime && block.timestamp <= endTime, "Donations are not allowed at this time");
        _;
    }

    /**
     * @dev 捐赠函数，允许用户向合约发送以太币（仅在时间限制内）
     */
    function donate() public payable withinTimeLimit {
        require(msg.value > 0, "Donation amount must be greater than 0");

        // 如果是新捐赠者，增加到捐赠者名单
        if(totalDonationsByUser[msg.sender] == 0) {
            donorsList.push(msg.sender);
            donorsCount++;
        }

        // 记录捐赠详情
        totalDonationsByUser[msg.sender] += msg.value;
        userDonations[msg.sender].push(Donation({
        amount: msg.value,
        timestamp: block.timestamp
        }));

        totalDonations += msg.value;

        // 触发捐赠事件
        emit Donated(msg.sender, msg.value, block.timestamp);
    }

    /**
     * @dev 允许合约所有者提取所有资金
     */
    function withdraw() public onlyOwner {
        require(address(this).balance > 0, "No balance to withdraw");

        uint256 amount = address(this).balance;
//        payable(owner).transfer(amount);

        (bool success, ) = payable(owner).call{value: amount}("");
        require(success, "Withdrawal failed");

        // 触发提款事件
        emit Withdrawn(owner, amount, block.timestamp);
    }

    /**
     * @dev 获取指定用户的捐赠次数
     * @param user 用户地址
     * @return 捐赠次数
     */
    function getUserDonationCount(address user) public view returns(uint256) {
        return userDonations[user].length;
    }

    /**
     * @dev 获取指定用户某次捐赠的详情
     * @param user 用户地址
     * @param index 捐赠索引
     * @return 捐赠金额和时间戳
     */
    function getUserDonationAt(address user, uint256 index) public view returns(uint256, uint256) {
        require(index < userDonations[user].length, "Index out of bounds");
        Donation storage donation = userDonations[user][index];
        return (donation.amount, donation.timestamp);
    }

    /**
     * @dev 获取合约余额
     * @return 合约当前余额
     */
    function getContractBalance() public view returns(uint256) {
        return address(this).balance;
    }

    /**
     * @dev 设置新的时间限制（仅所有者）
     * @param _startTime 新的开始时间
     * @param _endTime 新的结束时间
     */
    function setTimeLimit(uint256 _startTime, uint256 _endTime) public onlyOwner {
        require(_endTime > _startTime, "End time must be after start time");
        startTime = _startTime;
        endTime = _endTime;

        emit TimeLimitSet(startTime, endTime);
    }

    /**
     * @dev 获取当前时间是否在捐赠时间内
     * @return 是否在捐赠时间内
     */
    function isDonationAllowed() public view returns(bool) {
        return (block.timestamp >= startTime && block.timestamp <= endTime);
    }

    /**
     * @dev 获取捐赠排行榜（前3名）
     * @return 排行榜条目数组
     */
    function getLeaderboard() public view returns(LeaderboardEntry[] memory) {
        LeaderboardEntry[] memory leaderboard = new LeaderboardEntry[](donorsList.length);

        // 填充排行榜数据
        for(uint i = 0; i < donorsList.length; i++) {
            leaderboard[i] = LeaderboardEntry(donorsList[i], totalDonationsByUser[donorsList[i]]);
        }

        // 简单排序（选择排序），找出前三名
        for(uint i = 0; i < leaderboard.length && i < 3; i++) {
            uint maxIndex = i;
            for(uint j = i + 1; j < leaderboard.length; j++) {
                if(leaderboard[j].amount > leaderboard[maxIndex].amount) {
                    maxIndex = j;
                }
            }
            // 交换元素
            if(maxIndex != i) {
                LeaderboardEntry memory temp = leaderboard[i];
                leaderboard[i] = leaderboard[maxIndex];
                leaderboard[maxIndex] = temp;
            }
        }

        // 创建最终结果数组（最多3个元素）
        uint resultLength = leaderboard.length < 3 ? leaderboard.length : 3;
        LeaderboardEntry[] memory result = new LeaderboardEntry[](resultLength);
        for(uint i = 0; i < resultLength; i++) {
            result[i] = leaderboard[i];
        }

        return result;
    }

    /**
     * @dev 获取捐赠者名单长度
     * @return 捐赠者名单长度
     */
    function getDonorsListLength() public view returns(uint256) {
        return donorsList.length;
    }

    /**
     * @dev 获取捐赠者名单中的指定地址
     * @param index 索引
     * @return 捐赠者地址
     */
    function getDonorAtIndex(uint256 index) public view returns(address) {
        require(index < donorsList.length, "Index out of bounds");
        return donorsList[index];
    }
}
