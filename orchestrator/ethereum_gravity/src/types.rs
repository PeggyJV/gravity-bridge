use ethers::prelude::*;
use std::sync::Arc;

pub type EthSignerMiddleware<S> = SignerMiddleware<Provider<Http>, S>;
pub type EthClient<S> = Arc<EthSignerMiddleware<S>>;
