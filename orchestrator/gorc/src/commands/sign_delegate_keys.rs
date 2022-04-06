use crate::{application::APP, prelude::*};
use abscissa_core::{clap::Parser, Application, Command, Runnable};
use gravity_proto::gravity as proto;
use std::time::Duration;

/// Sign delegate keys command
#[derive(Command, Debug, Default, Parser)]
pub struct SignDelegateKeysCmd {
    /// Ethereum keyname
    ethereum_key: String,

    /// Validator address
    val_address: String,

    /// nonce
    nonce: Option<u64>,
}

impl Runnable for SignDelegateKeysCmd {
    fn run(&self) {
        let config = APP.config();
        abscissa_tokio::run_with_actix(&APP, async {
            let key = config.load_clarity_key(self.ethereum_key.clone());
            let address = self.val_address.parse().expect("Could not parse address");

            let nonce: u64 = match self.nonce {
                Some(nonce) => nonce,
                None => {
                    let timeout = Duration::from_secs(10);
                    let contact = deep_space::Contact::new(
                        &config.cosmos.grpc,
                        timeout,
                        &config.cosmos.prefix,
                    )
                    .expect("Could not create contact");

                    let account_info = contact.get_account_info(address).await;
                    let account_info = account_info.expect("Did not receive account info");
                    account_info.sequence
                }
            };

            let msg = proto::DelegateKeysSignMsg {
                validator_address: self.val_address.clone(),
                nonce,
            };

            let size = prost::Message::encoded_len(&msg);
            let mut buf = bytes::BytesMut::with_capacity(size);
            prost::Message::encode(&msg, &mut buf).expect("Failed to encode DelegateKeysSignMsg!");

            let signature = key.sign_ethereum_msg(&buf);

            println!("{}", signature);
        })
        .unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}
