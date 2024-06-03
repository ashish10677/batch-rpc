const { ethers } = require("hardhat");

async function main() {
    const RequestManager = await ethers.getContractFactory("RequestManager");
    const requestManager = await RequestManager.deploy();
    await requestManager.deployed();
    console.log("RequestManager deployed to:", requestManager.address);

    // Create a few requests
    for (let i = 0; i < 3; i++) {
        const requestId = await requestManager.generateRequestId(i);
        const targetAddress = "0x0000000000000000000000000000000000000000"; // Dummy address
        const targetChainDomain = i;
        const message = ethers.utils.formatBytes32String(`Request message ${i}`);
        const tx = await requestManager.createRequest(requestId, targetAddress, targetChainDomain, message);
        await tx.wait();
        console.log(`Request ${i} created with ID:`, requestId);
    }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
