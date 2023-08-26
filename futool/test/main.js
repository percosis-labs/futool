const {FURY_ENDPOINT_FUTOOL, BINANCE_CHAIN_ENDPOINT_FUTOOL, LOADED_FURY_MNEMONIC,
    LOADED_BINANCE_CHAIN_MNEMONIC, BEP3_ASSETS } = require("./config.js");
const { setup, loadFuryDeputies } = require("./futool.js");
const { incomingSwap, outgoingSwap } = require("./swap.js");

var main = async () => {
    // Initialize clients compatible with futool
    const clients = await setup(FURY_ENDPOINT_FUTOOL, BINANCE_CHAIN_ENDPOINT_FUTOOL,
        LOADED_FURY_MNEMONIC, LOADED_BINANCE_CHAIN_MNEMONIC);

    // Load each Fury deputy hot wallet
    await loadFuryDeputies(clients.furyClient, BEP3_ASSETS, 100000);

    await incomingSwap(clients.furyClient, clients.bnbClient, BEP3_ASSETS, "busd", 10200005);
    // await outgoingSwap(clients.furyClient, clients.bnbClient, BEP3_ASSETS, "busd", 500005);
};

main();
