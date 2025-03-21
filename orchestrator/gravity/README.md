# `gravity` crate

Core utility and client functionality for interacting with Gravity implementations on Ethereum and Cosmos chains.

Previously, this crate was three different crates `cosmos_gravity`, `ethereum_gravity`, and `gravity_utils`, and was combined due to heavy dependency overlap. The code is not well organized because the three previous crates weren't, and no refactoring at the function level has yet occured after combining them.

## Future work

- Reorganize modules and helper function locations to be more intuitive
- Refactor large and complex functions into smaller units that can be more easily tested
- Unit test coverage
