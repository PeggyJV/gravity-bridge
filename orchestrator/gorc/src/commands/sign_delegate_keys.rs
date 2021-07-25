use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use cosmos_gravity::query;
use gravity_proto::gravity as proto;
use gravity_utils::connection_prep::create_rpc_connections;
use orchestrator::main_loop::{
    ETH_ORACLE_LOOP_SPEED, ETH_SIGNER_LOOP_SPEED,
};
use relayer::main_loop::LOOP_SPEED as RELAYER_LOOP_SPEED;
use std::cmp::min;

#[derive(Command, Debug, Default, Options)]
pub struct SignDelegateKeysCmd {
    // TODO(levi) make event-nonce optional: can be queried from a running chain:
    // #[options(free, help = "sign-delegate-key [ethereum-key-name] [validator-address] (event-nonce)")]
    #[options(
        free,
        help = "sign-delegate-key [ethereum-key-name] [validator-address] (nonce)"
    )]
    pub args: Vec<String>,
}

impl Runnable for SignDelegateKeysCmd {
    fn run(&self) {
        let config = APP.config();
        abscissa_tokio::run(&APP, async {
            let name = self.args.get(0).expect("ethereum-key-name is required");
            let key = config.load_clarity_key(name.clone());

            let val = self.args.get(1).expect("validator-address is required");
            // TODO(levi) ensure this is a valoper address for the next release

            let cosmos_prefix = config.cosmos.prefix.clone();

            let timeout = min(
                min(ETH_SIGNER_LOOP_SPEED, ETH_ORACLE_LOOP_SPEED),
                RELAYER_LOOP_SPEED,
            );

            let connections = create_rpc_connections(
                cosmos_prefix,
                Some(config.cosmos.grpc.clone()),
                Some(config.ethereum.rpc.clone()),
                timeout,
            )
            .await;

            let mut grpc = connections.grpc.clone().unwrap();


            let valset = query::get_latest_valset(&mut grpc).await;

            println!("{:#?}", valset);

            // This is were I have Problems, How do I get the nonce? valset.nonce doesn't work. Everyother thing works fine
            // Without the nonce, the match below will keep throwing a mismatch error

            let nonce = match self.args.get(2) {
                Some(nonce) => nonce.clone(),
                None => valset,
            };

            let nonce = nonce.parse().expect("could not parse nonce");


            let msg = proto::DelegateKeysSignMsg {
                validator_address: val.clone(),
                nonce,
            };

            let size = prost::Message::encoded_len(&msg);
            let mut buf = bytes::BytesMut::with_capacity(size);
            prost::Message::encode(&msg, &mut buf).expect("Failed to encode DelegateKeysSignMsg!");

            let signature = key.sign_ethereum_msg(&buf);

            println!("{}", signature);
        });
    }
}
