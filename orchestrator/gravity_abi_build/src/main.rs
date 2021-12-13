use ethers::contract::Abigen;
use std::process;

fn main() {
    // Gravity contract

    let abigen = match Abigen::new("Gravity", "../gravity_abi/Gravity.json") {
        Ok(abigen) => abigen,
        Err(e) => {
            println!("Could not open Gravity.json: {}", e);
            process::exit(1);
        }
    };

    let abi = match abigen
        .add_event_derive("serde::Deserialize")
        .add_event_derive("serde::Serialize")
        .generate()
    {
        Ok(abi) => abi,
        Err(e) => {
            println!("Could not generate abi from Gravity.json: {}", e);
            process::exit(1);
        }
    };

    match abi.write_to_file("../gravity_abi/src/gravity.rs") {
        Ok(_) => (),
        Err(e) => println!("Error writing gravity.rs: {}", e),
    }

    // OpenZeppelin ERC20 contract

    let abigen = match Abigen::new("ERC20", "../gravity_abi/ERC20.json") {
        Ok(abigen) => abigen,
        Err(e) => {
            println!("Could not open ERC20.json: {}", e);
            process::exit(1);
        }
    };

    let abi = match abigen
        .add_event_derive("serde::Deserialize")
        .add_event_derive("serde::Serialize")
        .generate()
    {
        Ok(abi) => abi,
        Err(e) => {
            println!("Could not generate abi from ERC20.json: {}", e);
            process::exit(1);
        }
    };

    match abi.write_to_file("../gravity_abi/src/erc20.rs") {
        Ok(_) => (),
        Err(e) => println!("Error writing erc20.rs: {}", e),
    }
}
