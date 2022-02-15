mod add;
mod delete;
mod import;
mod list;
mod recover;
mod rename;
mod show;

use abscissa_core::{clap::Parser, Command, Runnable};

#[derive(Command, Debug, Parser, Runnable)]
pub enum EthKeysCmd {
    Add(add::AddEthKeyCmd),

    Delete(delete::DeleteEthKeyCmd),

    Import(import::ImportEthKeyCmd),

    List(list::ListEthKeyCmd),

    Recover(recover::RecoverEthKeyCmd),

    Rename(rename::RenameEthKeyCmd),

    Show(show::ShowEthKeyCmd),
}
