pragma solidity ^0.5.0;

contract UserStorage {
    struct User {
        string name;
        string email;
        string role;
        uint256 walletBalance;
        uint256 completedSessions;
        uint256 upcomingSessions;
        uint256 totalHours;
        uint256 averageRating;
    }

    mapping(address => User) private users;
    mapping(string => address) private emailToAddress;

    event UserCreated(address userAddress, string email);
    event UserUpdated(address userAddress);

    function createUser(string memory name, string memory email, string memory role) public {
        require(emailToAddress[email] == address(0), "Email already registered");
        users[msg.sender] = User(name, email, role, 0, 0, 0, 0, 0);
        emailToAddress[email] = msg.sender;
        emit UserCreated(msg.sender, email);
    }

    function getUser(address userAddress) public view returns (User memory) {
    require(users[userAddress].exists, "User not found");
    return users[userAddress];
}

    function getUserByEmail(string memory email) public view returns (User memory) {
    address userAddress = emailToAddress[email];
    require(userAddress != address(0), string(abi.encodePacked("User not found for email: ", email)));
    return getUser(userAddress);
}

    function updateUser(string memory name, string memory email, string memory role) public {
        User storage user = users[msg.sender];
        user.name = name;
        user.email = email;
        user.role = role;
        emit UserUpdated(msg.sender);
    }

    function updateWalletBalance(uint256 newBalance) public {
        users[msg.sender].walletBalance = newBalance;
        emit UserUpdated(msg.sender);
    }

    function updateLearningStats(uint256 completedSessions, uint256 upcomingSessions, uint256 totalHours, uint256 averageRating) public {
        User storage user = users[msg.sender];
        user.completedSessions = completedSessions;
        user.upcomingSessions = upcomingSessions;
        user.totalHours = totalHours;
        user.averageRating = averageRating;
        emit UserUpdated(msg.sender);
    }
}