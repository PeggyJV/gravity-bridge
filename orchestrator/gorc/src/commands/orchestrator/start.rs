use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use clarity::address::Address as EthAddress;
// use clarity::PrivateKey as EthPrivateKey;
// use deep_space::Contact;
// use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
// use orchestrator::main_loop::orchestrator_main_loop;
// use tonic::transport::Channel;
// use web30::client::Web3;

#[derive(Command, Debug, Options)]
pub struct StartCommand {
    #[options(help = "cosmos key name")]
    cosmos_key: String,

    #[options(help = "ethereum key name")]
    ethereum_key: String,
}

impl Runnable for StartCommand {
    fn run(&self) {
        let config = APP.config();

        let cosmos_key = config.load_deep_space_key(self.cosmos_key.clone());

        let ethereum_key = config.load_clarity_key(self.ethereum_key.clone());

        let contract_address: EthAddress = config
            .gravity
            .contract
            .parse()
            .expect("Could not parse gravity contract address");

        let pay_fees_in = config.gravity.fees_denom.clone();

        let _ = cosmos_key; // XXX deleteme
        let _ = ethereum_key; // XXX deleteme
        let _ = contract_address; // XXX deleteme
        let _ = pay_fees_in; // XXX deleteme

        //      let web3: Web3; // TODO init from config
        //      let contact: Contact; // TODO init from config
        //      let grpc_client: GravityQueryClient<Channel>; // TODO init from config

        //      abscissa_tokio::run(&APP, async {
        //          orchestrator_main_loop(
        //              cosmos_key,
        //              ethereum_key,
        //              web3,
        //              contact,
        //              grpc_client,
        //              gravity_contract_address,
        //              pay_fees_in,
        //          )
        //          .await;
        //      });
    }
}
