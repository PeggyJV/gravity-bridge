use crate::application::APP;
use abscissa_core::{clap::Parser, status_err, Application, Command, Runnable};
use clarity::Uint256;
use cosmos_gravity::send::{send_request_batch_tx, send_to_eth};
use deep_space::coin::Coin;
use ethers::types::Address as EthAddress;
use gravity_proto::gravity::DenomToErc20Request;
use gravity_utils::connection_prep::{check_for_fee_denom, create_rpc_connections};
use std::{process::exit, time::Duration};

const TIMEOUT: Duration = Duration::from_secs(60);

/// This command, send Cosmos to Ethereum
#[derive(Command, Debug, Default, Parser)]
pub struct CosmosToEthCmd {
    pub args: Vec<String>,

    #[clap(short, long)]
    pub flag_no_batch: bool,
}

pub fn one_eth() -> f64 {
    1000000000000000000f64
}

pub fn one_atom() -> f64 {
    1000000f64
}

// TODO(bolten): deep_space's Coin type relies internally on clarity's Uint256,
// and it would make this code super akward to try to get around that...replacing
// that here might be part of a broader deep_space replacement
pub fn print_atom(input: Uint256) -> String {
    let float: f64 = input.to_string().parse().unwrap();
    let res = float / one_atom();
    format!("{}", res)
}

pub fn print_eth(input: Uint256) -> String {
    let float: f64 = input.to_string().parse().unwrap();
    let res = float / one_eth();
    format!("{}", res)
}

impl Runnable for CosmosToEthCmd {
    fn run(&self) {
        let config = APP.config();
        let gravity_denom = self.args.get(0).expect("denom is required");
        let gravity_denom = gravity_denom.to_string();
        let is_cosmos_originated = !gravity_denom.starts_with("gravity");

        let amount = self.args.get(1).expect("amount is required");
        let amount: Uint256 = amount.parse().expect("cannot parse amount");

        let cosmos_key = self.args.get(2).expect("name is required");
        let cosmos_key = config.load_deep_space_key(cosmos_key.to_string());

        let cosmos_prefix = config.cosmos.prefix.trim();
        let cosmos_address = cosmos_key.to_address(cosmos_prefix).unwrap();
        let cosmos_grpc = config.cosmos.grpc.trim();
        let cosmos_granter = config.cosmos.granter.clone();
        println!("Sending from Cosmos address {}", cosmos_address);
        abscissa_tokio::run_with_actix(&APP, async {
        let connections = create_rpc_connections(
            cosmos_prefix.to_string(),
            Some(cosmos_grpc.to_string()),
            None,
            TIMEOUT,
        )
        .await;
        let contact = connections.contact.unwrap();
        let mut grpc = connections.grpc.unwrap();
        let res = grpc
            .denom_to_erc20(DenomToErc20Request {
                denom: gravity_denom.clone(),
            })
            .await;
        match res {
            Ok(val) => println!(
                "Asset {} has ERC20 representation {}",
                gravity_denom,
                val.into_inner().erc20
            ),
            Err(_e) => {
                println!(
                    "Asset {} has no ERC20 representation, you may need to deploy an ERC20 for it!",
                    gravity_denom
                );
                exit(1);
            }
        }

        let amount = Coin {
            amount: amount.clone(),
            denom: gravity_denom.clone(),
        };
        let bridge_fee = Coin {
            amount: 1u64.into(),
            denom: gravity_denom.clone(),
        };

        let eth_dest = self.args.get(3).expect("ethereum destination is required");
        let eth_dest: EthAddress = eth_dest.parse().expect("cannot parse ethereum address");
        check_for_fee_denom(&gravity_denom, cosmos_address, &contact).await;

        let balances = contact
            .get_balances(cosmos_address)
            .await
            .expect("Failed to get balances!");
        let mut found = None;
        for coin in balances.iter() {
            if coin.denom == gravity_denom {
                found = Some(coin);
            }
        }

        println!("Cosmos balances {:?}", balances);

        let times = self.args.get(4).expect("times is required");
        let times = times.parse::<usize>().expect("cannot parse times");

        match found {
            None => panic!("You don't have any {} tokens!", gravity_denom),
            Some(found) => {
                if amount.amount.clone() * times.into() >= found.amount && times == 1 {
                    if is_cosmos_originated {
                        panic!("Your transfer of {} {} tokens is greater than your balance of {} tokens. Remember you need some to pay for fees!", print_atom(amount.amount), gravity_denom, print_atom(found.amount.clone()));
                    } else {
                        panic!("Your transfer of {} {} tokens is greater than your balance of {} tokens. Remember you need some to pay for fees!", print_eth(amount.amount), gravity_denom, print_eth(found.amount.clone()));
                    }
                } else if amount.amount.clone() * times.into() >= found.amount {
                    if is_cosmos_originated {
                        panic!("Your transfer of {} * {} {} tokens is greater than your balance of {} tokens. Try to reduce the amount or the --times parameter", print_atom(amount.amount), times, gravity_denom, print_atom(found.amount.clone()));
                    } else {
                        panic!("Your transfer of {} * {} {} tokens is greater than your balance of {} tokens. Try to reduce the amount or the --times parameter", print_eth(amount.amount), times, gravity_denom, print_eth(found.amount.clone()));
                    }
                }
            }
        }

        let mut successful_sends = 0;
        for _ in 0..times {
            println!(
                "Locking {} / {} into the batch pool",
                amount.clone(),
                gravity_denom
            );
            let res = send_to_eth(
                cosmos_key,
                cosmos_granter.clone(),
                eth_dest,
                amount.clone(),
                bridge_fee.clone(),
                config.cosmos.gas_price.as_tuple(),
                &contact,
                1.0
            )
            .await;
            match res {
                Ok(tx_id) => {
                    successful_sends += 1;
                    println!("Send to Eth txid {}", tx_id.txhash);
                }
                Err(e) => println!("Failed to send tokens! {:?}", e),
            }
        }

        if successful_sends > 0 {
            if !self.flag_no_batch {
                println!("Requesting a batch to push transaction along immediately");
                send_request_batch_tx(cosmos_key, cosmos_granter.clone(), gravity_denom,config.cosmos.gas_price.as_tuple(), &contact,config.cosmos.gas_adjustment)
                    .await
                    .expect("Failed to request batch");
            } else {
                println!("--no-batch specified, your transfer will wait until someone requests a batch for this token type")
            }
        } else {
            println!("No successful sends, no batch will be sent")
        }
        })
        .unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
        });
    }
}
