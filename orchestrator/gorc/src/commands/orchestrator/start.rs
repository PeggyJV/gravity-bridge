// use crate::application::APP;
use abscissa_core::{Command, Options, Runnable};
// use clarity::address::Address as EthAddress;
// use clarity::PrivateKey as EthPrivateKey;
// use deep_space::Contact;
// use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
// use orchestrator::main_loop::orchestrator_main_loop;
// use tonic::transport::Channel;
// use web30::client::Web3;

#[derive(Command, Debug, Options)]
pub struct StartCommand {
    //
}

impl Runnable for StartCommand {
    fn run(&self) {
        //      let cosmos_key: deep_space::private_key::PrivateKey; // TODO init from keystore
        //      let ethereum_key: EthPrivateKey; // TODO init from keystore

        //      // let config = APP.config();
        //      let web3: Web3; // TODO init from config
        //      let contact: Contact; // TODO init from config
        //      let grpc_client: GravityQueryClient<Channel>; // TODO init from config
        //      let gravity_contract_address: EthAddress; // TODO init from config
        //      let pay_fees_in: String; // TODO init from config

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
