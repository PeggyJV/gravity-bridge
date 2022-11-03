use crate::{application::APP, prelude::*};
use abscissa_core::{clap::Parser, Application, Command, Runnable};
use ethers::{prelude::Signer, utils::keccak256};
use ocular::prelude::*;
use ocular_somm_gravity::SommGravity;

/// Sign delegate keys command
#[derive(Command, Debug, Default, Parser)]
pub struct SignDelegateKeysCmd {
    pub args: Vec<String>,
}

impl Runnable for SignDelegateKeysCmd {
    fn run(&self) {
        let config = APP.config();
        abscissa_tokio::run_with_actix(&APP, async {
            let name = self.args.get(0).expect("ethereum-key-name is required");
            let ethereum_wallet = config.load_ethers_wallet(name.clone());

            let validator_address = self
                .args
                .get(1)
                .expect("validator-address is required")
                .to_owned();

            let nonce: u64 = match self.args.get(2) {
                Some(nonce) => nonce.parse().expect("cannot parse nonce"),
                None => {
                    let mut cosmos_client = GrpcClient::new(&config.cosmos.grpc)
                        .await
                        .expect("failed to construct GrpcClient");
                    let account = cosmos_client.query_account(&validator_address).await;
                    let account = account.expect("Did not receive account info");
                    account.sequence
                }
            };

            let msg = SommGravity::DelegateKeysSignMsg {
                validator_address: &validator_address,
                nonce,
            }
            .into_any()
            .unwrap()
            .value;

            let data = keccak256(msg);
            let signature = ethereum_wallet
                .sign_message(data)
                .await
                .expect("Could not sign DelegateKeysSignMsg");

            println!("{}", signature);
        })
        .unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}
