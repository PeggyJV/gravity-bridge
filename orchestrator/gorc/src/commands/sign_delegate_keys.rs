use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use cosmos_gravity::query as QueryValSet;
use gravity_proto::gravity as proto;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use tonic::transport::Channel;

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

            let nonce = self.args.get(2).expect("nonce is required");
            let nonce = nonce.parse().expect("could not parse nonce");

            let client = &mut config.cosmos.grpc.clone();
            let client = client
                .parse::<GravityQueryClient<Channel>>()
                .expect("Could not parse derivation path");

            let valset = QueryValSet::get_valset(&mut client, nonce);

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
