mod add;
mod delete;
mod list;
mod recover;
mod rename;
mod show;

use abscissa_core::{clap::Parser, Command, Runnable};

#[derive(Command, Debug, Parser, Runnable)]
pub enum CosmosKeysCmd {
    Add(add::AddCosmosKeyCmd),

    Delete(delete::DeleteCosmosKeyCmd),

    Recover(recover::RecoverCosmosKeyCmd),

    Rename(rename::RenameCosmosKeyCmd),

    List(list::ListCosmosKeyCmd),

    Show(show::ShowCosmosKeyCmd),
}
