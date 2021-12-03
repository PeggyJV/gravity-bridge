use ethers::prelude::*;
use std::sync::Arc;

pub type EthSignerMiddleware = SignerMiddleware<Provider<Http>, LocalWallet>;
pub type EthClient = Arc<EthSignerMiddleware>;
