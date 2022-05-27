/// EVMEventVoteRecord is an event that is pending of confirmation by 2/3 of
/// the signer set. The event is then attested and executed in the state machine
/// once the required threshold is met.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct EvmEventVoteRecord {
    #[prost(message, optional, tag = "1")]
    pub event: ::core::option::Option<::prost_types::Any>,
    #[prost(string, repeated, tag = "2")]
    pub votes: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    #[prost(bool, tag = "3")]
    pub accepted: bool,
    #[prost(uint32, tag = "4")]
    pub chain_id: u32,
}
/// LatestEVMBlockHeight defines the latest observed EVM block height
/// and the corresponding timestamp value in nanoseconds.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LatestEvmBlockHeight {
    #[prost(uint64, tag = "1")]
    pub evm_height: u64,
    #[prost(uint64, tag = "2")]
    pub cosmos_height: u64,
    #[prost(uint32, tag = "3")]
    pub chain_id: u32,
}
/// EVMSigner represents a cosmos validator with its corresponding bridge
/// operator EVM address and its staking consensus power.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct EvmSigner {
    #[prost(uint64, tag = "1")]
    pub power: u64,
    #[prost(string, tag = "2")]
    pub evm_address: ::prost::alloc::string::String,
    #[prost(uint32, tag = "3")]
    pub chain_id: u32,
}
/// SignerSetTx is the EVM Bridge multisig set that relays
/// transactions the two chains. The staking validators keep EVM keys which
/// are used to check signatures on EVM in order to get significant gas
/// savings.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignerSetTx {
    #[prost(uint64, tag = "1")]
    pub nonce: u64,
    #[prost(uint64, tag = "2")]
    pub height: u64,
    #[prost(message, repeated, tag = "3")]
    pub signers: ::prost::alloc::vec::Vec<EvmSigner>,
    #[prost(uint32, tag = "4")]
    pub chain_id: u32,
}
/// BatchTx represents a batch of transactions going from Cosmos to EVM.
/// Batch txs are are identified by a unique hash and the token contract that is
/// shared by all the SendToEVM
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTx {
    #[prost(uint64, tag = "1")]
    pub batch_nonce: u64,
    #[prost(uint64, tag = "2")]
    pub timeout: u64,
    #[prost(message, repeated, tag = "3")]
    pub transactions: ::prost::alloc::vec::Vec<SendToEvm>,
    #[prost(string, tag = "4")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(uint64, tag = "5")]
    pub height: u64,
    #[prost(uint32, tag = "6")]
    pub chain_id: u32,
}
/// SendToEVM represents an individual SendToEVM from Cosmos to
/// EVM
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SendToEvm {
    #[prost(uint64, tag = "1")]
    pub id: u64,
    #[prost(string, tag = "2")]
    pub sender: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub evm_recipient: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "4")]
    pub erc20_token: ::core::option::Option<Erc20Token>,
    #[prost(message, optional, tag = "5")]
    pub erc20_fee: ::core::option::Option<Erc20Token>,
    #[prost(uint32, tag = "6")]
    pub chain_id: u32,
}
/// ContractCallTx represents an individual arbitrary logic call transaction
/// from Cosmos to EVM.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTx {
    #[prost(uint64, tag = "1")]
    pub invalidation_nonce: u64,
    #[prost(bytes = "vec", tag = "2")]
    pub invalidation_scope: ::prost::alloc::vec::Vec<u8>,
    #[prost(string, tag = "3")]
    pub address: ::prost::alloc::string::String,
    #[prost(bytes = "vec", tag = "4")]
    pub payload: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag = "5")]
    pub timeout: u64,
    #[prost(message, repeated, tag = "6")]
    pub tokens: ::prost::alloc::vec::Vec<Erc20Token>,
    #[prost(message, repeated, tag = "7")]
    pub fees: ::prost::alloc::vec::Vec<Erc20Token>,
    #[prost(uint64, tag = "8")]
    pub height: u64,
    #[prost(uint32, tag = "9")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Erc20Token {
    #[prost(string, tag = "1")]
    pub contract: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub amount: ::prost::alloc::string::String,
    #[prost(uint32, tag = "3")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct IdSet {
    #[prost(uint64, repeated, tag = "1")]
    pub ids: ::prost::alloc::vec::Vec<u64>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CommunityPoolEvmSpendProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub recipient: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "4")]
    pub amount: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(message, optional, tag = "5")]
    pub bridge_fee: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(uint32, tag = "6")]
    pub chain_id: u32,
}
/// This format of the community spend EVM proposal is specifically for
/// the CLI to allow simple text serialization.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CommunityPoolEvmSpendProposalForCli {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub recipient: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub amount: ::prost::alloc::string::String,
    #[prost(string, tag = "5")]
    pub bridge_fee: ::prost::alloc::string::String,
    #[prost(string, tag = "6")]
    pub deposit: ::prost::alloc::string::String,
    #[prost(uint32, tag = "7")]
    pub chain_id: u32,
}
/// MsgSendToEVM submits a SendToEVM attempt to bridge an asset over to
/// EVM. The SendToEVM will be stored and then included in a batch and
/// then submitted to EVM.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSendToEvm {
    #[prost(string, tag = "1")]
    pub sender: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub evm_recipient: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "3")]
    pub amount: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(message, optional, tag = "4")]
    pub bridge_fee: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(uint32, tag = "5")]
    pub chain_id: u32,
}
/// MsgSendToEVMResponse returns the SendToEVM transaction ID which
/// will be included in the batch tx.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSendToEvmResponse {
    #[prost(uint64, tag = "1")]
    pub id: u64,
}
/// MsgCancelSendToEVM allows the sender to cancel its own outgoing
/// SendToEVM tx and receive a refund of the tokens and bridge fees. This tx
/// will only succeed if the SendToEVM tx hasn't been batched to be
/// processed and relayed to EVM.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgCancelSendToEvm {
    #[prost(uint64, tag = "1")]
    pub id: u64,
    #[prost(string, tag = "2")]
    pub sender: ::prost::alloc::string::String,
    #[prost(uint32, tag = "3")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgCancelSendToEvmResponse {}
/// MsgRequestBatchTx requests a batch of transactions with a given coin
/// denomination to send across the bridge to EVM.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgRequestBatchTx {
    #[prost(string, tag = "1")]
    pub denom: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub signer: ::prost::alloc::string::String,
    #[prost(uint32, tag = "3")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgRequestBatchTxResponse {}
/// MsgSubmitEVMTxConfirmation submits an EVM signature for a given
/// validator
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitEvmTxConfirmation {
    /// TODO: can we make this take an array?
    #[prost(message, optional, tag = "1")]
    pub confirmation: ::core::option::Option<::prost_types::Any>,
    #[prost(string, tag = "2")]
    pub signer: ::prost::alloc::string::String,
    #[prost(uint32, tag = "3")]
    pub chain_id: u32,
}
/// ContractCallTxConfirmation is a signature on behalf of a validator for a
/// ContractCallTx.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxConfirmation {
    #[prost(bytes = "vec", tag = "1")]
    pub invalidation_scope: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag = "2")]
    pub invalidation_nonce: u64,
    #[prost(string, tag = "3")]
    pub evm_signer: ::prost::alloc::string::String,
    #[prost(bytes = "vec", tag = "4")]
    pub signature: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint32, tag = "5")]
    pub chain_id: u32,
}
/// BatchTxConfirmation is a signature on behalf of a validator for a BatchTx.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxConfirmation {
    #[prost(string, tag = "1")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub batch_nonce: u64,
    #[prost(string, tag = "3")]
    pub evm_signer: ::prost::alloc::string::String,
    #[prost(bytes = "vec", tag = "4")]
    pub signature: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint32, tag = "5")]
    pub chain_id: u32,
}
/// SignerSetTxConfirmation is a signature on behalf of a validator for a
/// SignerSetTx
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignerSetTxConfirmation {
    #[prost(uint64, tag = "1")]
    pub signer_set_nonce: u64,
    #[prost(string, tag = "2")]
    pub evm_signer: ::prost::alloc::string::String,
    #[prost(bytes = "vec", tag = "3")]
    pub signature: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint32, tag = "4")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitEvmTxConfirmationResponse {}
/// MsgSubmitEVMEvent
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitEvmEvent {
    #[prost(message, optional, tag = "1")]
    pub event: ::core::option::Option<::prost_types::Any>,
    #[prost(string, tag = "2")]
    pub signer: ::prost::alloc::string::String,
    #[prost(uint32, tag = "3")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitEvmEventResponse {}
/// MsgDelegateKeys allows validators to delegate their voting responsibilities
/// to a given orchestrator address. This key is then used as an optional
/// authentication method for attesting events from EVM.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgDelegateKeys {
    #[prost(string, tag = "1")]
    pub validator_address: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub orchestrator_address: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub evm_address: ::prost::alloc::string::String,
    #[prost(bytes = "vec", tag = "4")]
    pub eth_signature: ::prost::alloc::vec::Vec<u8>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgDelegateKeysResponse {}
/// DelegateKeysSignMsg defines the message structure an operator is expected to
/// sign when submitting a MsgDelegateKeys message. The resulting signature
/// should populate the eth_signature field.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysSignMsg {
    #[prost(string, tag = "1")]
    pub validator_address: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub nonce: u64,
}
/// Periodic update of latest observed EVM and Cosmos heights from the
/// orchestrator
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgEvmHeightVote {
    #[prost(uint64, tag = "1")]
    pub evm_height: u64,
    #[prost(string, tag = "2")]
    pub signer: ::prost::alloc::string::String,
    #[prost(uint32, tag = "3")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgEvmHeightVoteResponse {}
////////////
// Events //
////////////

/// SendToCosmosEvent is submitted when the SendToCosmosEvent is emitted by they
/// gravity contract. ERC20 representation coins are minted to the cosmosreceiver
/// address.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SendToCosmosEvent {
    #[prost(uint64, tag = "1")]
    pub event_nonce: u64,
    #[prost(string, tag = "2")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub amount: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub evm_sender: ::prost::alloc::string::String,
    #[prost(string, tag = "5")]
    pub cosmos_receiver: ::prost::alloc::string::String,
    #[prost(uint64, tag = "6")]
    pub evm_height: u64,
    #[prost(uint32, tag = "7")]
    pub chain_id: u32,
}
/// BatchExecutedEvent claims that a batch of BatchTxExecutedal operations on the
/// bridge contract was executed successfully on ETH
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchExecutedEvent {
    #[prost(string, tag = "1")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub event_nonce: u64,
    #[prost(uint64, tag = "3")]
    pub evm_height: u64,
    #[prost(uint64, tag = "4")]
    pub batch_nonce: u64,
    #[prost(uint32, tag = "5")]
    pub chain_id: u32,
}
// ContractCallExecutedEvent describes a contract call that has been
// successfully executed on EVM.

/// NOTE: bytes.HexBytes is supposed to "help" with json encoding/decoding
/// investigate?
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallExecutedEvent {
    #[prost(uint64, tag = "1")]
    pub event_nonce: u64,
    #[prost(bytes = "vec", tag = "2")]
    pub invalidation_scope: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag = "3")]
    pub invalidation_nonce: u64,
    #[prost(uint64, tag = "4")]
    pub evm_height: u64,
    #[prost(uint32, tag = "5")]
    pub chain_id: u32,
}
/// ERC20DeployedEvent is submitted when an ERC20 contract
/// for a Cosmos SDK coin has been deployed on EVM.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Erc20DeployedEvent {
    #[prost(uint64, tag = "1")]
    pub event_nonce: u64,
    #[prost(string, tag = "2")]
    pub cosmos_denom: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub erc20_name: ::prost::alloc::string::String,
    #[prost(string, tag = "5")]
    pub erc20_symbol: ::prost::alloc::string::String,
    #[prost(uint64, tag = "6")]
    pub erc20_decimals: u64,
    #[prost(uint64, tag = "7")]
    pub evm_height: u64,
    #[prost(uint32, tag = "8")]
    pub chain_id: u32,
}
/// This informs the Cosmos module that a validator
/// set has been updated.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignerSetTxExecutedEvent {
    #[prost(uint64, tag = "1")]
    pub event_nonce: u64,
    #[prost(uint64, tag = "2")]
    pub signer_set_tx_nonce: u64,
    #[prost(uint64, tag = "3")]
    pub evm_height: u64,
    #[prost(message, repeated, tag = "4")]
    pub members: ::prost::alloc::vec::Vec<EvmSigner>,
    #[prost(uint32, tag = "5")]
    pub chain_id: u32,
}
#[doc = r" Generated client implementations."]
pub mod msg_client {
    #![allow(unused_variables, dead_code, missing_docs)]
    use tonic::codegen::*;
    #[doc = " Msg defines the state transitions possible within gravity"]
    pub struct MsgClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl MsgClient<tonic::transport::Channel> {
        #[doc = r" Attempt to create a new client by connecting to a given endpoint."]
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: std::convert::TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> MsgClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::ResponseBody: Body + HttpBody + Send + 'static,
        T::Error: Into<StdError>,
        <T::ResponseBody as HttpBody>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_interceptor(inner: T, interceptor: impl Into<tonic::Interceptor>) -> Self {
            let inner = tonic::client::Grpc::with_interceptor(inner, interceptor);
            Self { inner }
        }
        pub async fn set_delegate_keys(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgDelegateKeys>,
        ) -> Result<tonic::Response<super::MsgDelegateKeysResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Msg/SetDelegateKeys");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn send_to_evm(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgSendToEvm>,
        ) -> Result<tonic::Response<super::MsgSendToEvmResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Msg/SendToEVM");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn cancel_send_to_evm(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgCancelSendToEvm>,
        ) -> Result<tonic::Response<super::MsgCancelSendToEvmResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Msg/CancelSendToEVM");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn request_batch_tx(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgRequestBatchTx>,
        ) -> Result<tonic::Response<super::MsgRequestBatchTxResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Msg/RequestBatchTx");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn submit_evm_tx_confirmation(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgSubmitEvmTxConfirmation>,
        ) -> Result<tonic::Response<super::MsgSubmitEvmTxConfirmationResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/gravity.v2.Msg/SubmitEVMTxConfirmation");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn submit_evm_event(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgSubmitEvmEvent>,
        ) -> Result<tonic::Response<super::MsgSubmitEvmEventResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Msg/SubmitEVMEvent");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn submit_evm_height_vote(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgEvmHeightVote>,
        ) -> Result<tonic::Response<super::MsgEvmHeightVoteResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Msg/SubmitEVMHeightVote");
            self.inner.unary(request.into_request(), path, codec).await
        }
    }
    impl<T: Clone> Clone for MsgClient<T> {
        fn clone(&self) -> Self {
            Self {
                inner: self.inner.clone(),
            }
        }
    }
    impl<T> std::fmt::Debug for MsgClient<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "MsgClient {{ ... }}")
        }
    }
}
/// Params represent the Gravity genesis and store parameters
/// gravity_id:
/// a random 32 byte value to prevent signature reuse, for example if the
/// cosmos validators decided to use the same EVM keys for another chain
/// also running Gravity we would not want it to be possible to play a deposit
/// from chain A back on chain B's Gravity. This value IS USED ON EVM so
/// it must be set in your genesis.json before launch and not changed after
/// deploying Gravity
///
/// contract_hash:
/// the code hash of a known good version of the Gravity contract
/// solidity code. This can be used to verify the correct version
/// of the contract has been deployed. This is a reference value for
/// goernance action only it is never read by any Gravity code
///
/// bridge_EVM_address:
/// is address of the bridge contract on the EVM side, this is a
/// reference value for governance only and is not actually used by any
/// Gravity code
///
/// bridge_chain_id:
/// the unique identifier of the EVM chain, this is a reference value
/// only and is not actually used by any Gravity code
///
/// These reference values may be used by future Gravity client implemetnations
/// to allow for saftey features or convenience features like the Gravity address
/// in your relayer. A relayer would require a configured Gravity address if
/// governance had not set the address on the chain it was relaying for.
///
/// signed_signer_set_txs_window
/// signed_batches_window
/// signed_EVM_signatures_window
///
/// These values represent the time in blocks that a validator has to submit
/// a signature for a batch or valset, or to submit a EVM_signature for a
/// particular attestation nonce. In the case of attestations this clock starts
/// when the attestation is created, but only allows for slashing once the event
/// has passed
///
/// target_evm_tx_timeout:
///
/// This is the 'target' value for when EVM transactions time out, this is a
/// target because EVM is a probabilistic chain and you can't say for sure
/// what the block frequency is ahead of time.
///
/// average_block_time
/// average_EVM_block_time
///
/// These values are the average Cosmos block time and EVM block time
/// respectively and they are used to compute what the target batch timeout is.
/// It is important that governance updates these in case of any major, prolonged
/// change in the time it takes to produce a block
///
/// slash_fraction_signer_set_tx
/// slash_fraction_batch
/// slash_fraction_EVM_signature
/// slash_fraction_conflicting_EVM_signature
///
/// The slashing fractions for the various gravity related slashing conditions.
/// The first three refer to not submitting a particular message, the third for
/// submitting a different EVM_signature for the same EVM event
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Params {
    #[prost(map = "uint32, message", tag = "1")]
    pub chain_params: ::std::collections::HashMap<u32, ChainParams>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ChainParams {
    #[prost(string, tag = "1")]
    pub gravity_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub contract_source_hash: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub signed_signer_set_txs_window: u64,
    #[prost(uint64, tag = "4")]
    pub signed_batches_window: u64,
    #[prost(uint64, tag = "5")]
    pub evm_signatures_window: u64,
    #[prost(uint64, tag = "6")]
    pub target_evm_tx_timeout: u64,
    #[prost(uint64, tag = "7")]
    pub average_block_time: u64,
    #[prost(uint64, tag = "8")]
    pub average_evm_block_time: u64,
    /// TODO: slash fraction for contract call txs too
    #[prost(bytes = "vec", tag = "9")]
    pub slash_fraction_signer_set_tx: ::prost::alloc::vec::Vec<u8>,
    #[prost(bytes = "vec", tag = "10")]
    pub slash_fraction_batch: ::prost::alloc::vec::Vec<u8>,
    #[prost(bytes = "vec", tag = "11")]
    pub slash_fraction_evm_signature: ::prost::alloc::vec::Vec<u8>,
    #[prost(bytes = "vec", tag = "12")]
    pub slash_fraction_conflicting_evm_signature: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag = "13")]
    pub unbond_slashing_signer_set_txs_window: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisStateMultiChain {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
    #[prost(message, repeated, tag = "2")]
    pub delegate_keys: ::prost::alloc::vec::Vec<MsgDelegateKeys>,
    #[prost(message, repeated, tag = "3")]
    pub chain_genesis_states: ::prost::alloc::vec::Vec<ChainGenesisState>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ChainGenesisState {
    #[prost(uint32, tag = "1")]
    pub chain_id: u32,
    #[prost(uint64, tag = "2")]
    pub last_observed_event_nonce: u64,
    #[prost(message, repeated, tag = "3")]
    pub outgoing_txs: ::prost::alloc::vec::Vec<::prost_types::Any>,
    #[prost(message, repeated, tag = "4")]
    pub confirmations: ::prost::alloc::vec::Vec<::prost_types::Any>,
    #[prost(message, repeated, tag = "5")]
    pub evm_event_vote_records: ::prost::alloc::vec::Vec<EvmEventVoteRecord>,
    #[prost(message, repeated, tag = "6")]
    pub erc20_to_denoms: ::prost::alloc::vec::Vec<Erc20ToDenom>,
    #[prost(message, repeated, tag = "7")]
    pub unbatched_send_to_evm_txs: ::prost::alloc::vec::Vec<SendToEvm>,
}
/// This records the relationship between an ERC20 token and the denom
/// of the corresponding Cosmos originated asset
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Erc20ToDenom {
    #[prost(string, tag = "1")]
    pub erc20: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub denom: ::prost::alloc::string::String,
}
///  rpc Params
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ParamsRequest {}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ParamsResponse {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
}
///  rpc SignerSetTx
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignerSetTxRequest {
    #[prost(uint64, tag = "1")]
    pub signer_set_nonce: u64,
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LatestSignerSetTxRequest {
    #[prost(uint32, tag = "1")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignerSetTxResponse {
    #[prost(message, optional, tag = "1")]
    pub signer_set: ::core::option::Option<SignerSetTx>,
}
///  rpc BatchTx
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxRequest {
    #[prost(string, tag = "1")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub batch_nonce: u64,
    #[prost(uint32, tag = "3")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxResponse {
    #[prost(message, optional, tag = "1")]
    pub batch: ::core::option::Option<BatchTx>,
}
///  rpc ContractCallTx
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxRequest {
    #[prost(bytes = "vec", tag = "1")]
    pub invalidation_scope: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag = "2")]
    pub invalidation_nonce: u64,
    #[prost(uint32, tag = "3")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxResponse {
    #[prost(message, optional, tag = "1")]
    pub logic_call: ::core::option::Option<ContractCallTx>,
}
/// rpc SignerSetTxConfirmations
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignerSetTxConfirmationsRequest {
    #[prost(uint64, tag = "1")]
    pub signer_set_nonce: u64,
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignerSetTxConfirmationsResponse {
    #[prost(message, repeated, tag = "1")]
    pub signatures: ::prost::alloc::vec::Vec<SignerSetTxConfirmation>,
}
///  rpc SignerSetTxs
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignerSetTxsRequest {
    #[prost(message, optional, tag = "1")]
    pub pagination: ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageRequest>,
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SignerSetTxsResponse {
    #[prost(message, repeated, tag = "1")]
    pub signer_sets: ::prost::alloc::vec::Vec<SignerSetTx>,
    #[prost(message, optional, tag = "2")]
    pub pagination:
        ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageResponse>,
}
///  rpc BatchTxs
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxsRequest {
    #[prost(message, optional, tag = "1")]
    pub pagination: ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageRequest>,
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxsResponse {
    #[prost(message, repeated, tag = "1")]
    pub batches: ::prost::alloc::vec::Vec<BatchTx>,
    #[prost(message, optional, tag = "2")]
    pub pagination:
        ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageResponse>,
}
///  rpc ContractCallTxs
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxsRequest {
    #[prost(message, optional, tag = "1")]
    pub pagination: ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageRequest>,
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxsResponse {
    #[prost(message, repeated, tag = "1")]
    pub calls: ::prost::alloc::vec::Vec<ContractCallTx>,
    #[prost(message, optional, tag = "2")]
    pub pagination:
        ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageResponse>,
}
// NOTE(levi) pending queries: this is my address; what do I need to sign??
// why orchestrator key? hot, signing thing all the time so validator key can be
// safer

/// rpc UnsignedSignerSetTxs
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnsignedSignerSetTxsRequest {
    /// NOTE: this is an sdk.AccAddress and can represent either the
    /// orchestrator address or the corresponding validator address
    #[prost(string, tag = "1")]
    pub address: ::prost::alloc::string::String,
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnsignedSignerSetTxsResponse {
    #[prost(message, repeated, tag = "1")]
    pub signer_sets: ::prost::alloc::vec::Vec<SignerSetTx>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnsignedBatchTxsRequest {
    /// NOTE: this is an sdk.AccAddress and can represent either the
    /// orchestrator address or the corresponding validator address
    #[prost(string, tag = "1")]
    pub address: ::prost::alloc::string::String,
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnsignedBatchTxsResponse {
    /// Note these are returned with the signature empty
    #[prost(message, repeated, tag = "1")]
    pub batches: ::prost::alloc::vec::Vec<BatchTx>,
}
///  rpc UnsignedContractCallTxs
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnsignedContractCallTxsRequest {
    #[prost(string, tag = "1")]
    pub address: ::prost::alloc::string::String,
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnsignedContractCallTxsResponse {
    #[prost(message, repeated, tag = "1")]
    pub calls: ::prost::alloc::vec::Vec<ContractCallTx>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxFeesRequest {
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxFeesResponse {
    #[prost(message, repeated, tag = "1")]
    pub fees: ::prost::alloc::vec::Vec<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxConfirmationsRequest {
    #[prost(bytes = "vec", tag = "1")]
    pub invalidation_scope: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag = "2")]
    pub invalidation_nonce: u64,
    #[prost(uint32, tag = "3")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxConfirmationsResponse {
    #[prost(message, repeated, tag = "1")]
    pub signatures: ::prost::alloc::vec::Vec<ContractCallTxConfirmation>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxConfirmationsRequest {
    #[prost(uint64, tag = "1")]
    pub batch_nonce: u64,
    #[prost(string, tag = "2")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(uint32, tag = "3")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxConfirmationsResponse {
    #[prost(message, repeated, tag = "1")]
    pub signatures: ::prost::alloc::vec::Vec<BatchTxConfirmation>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LastSubmittedEvmEventRequest {
    #[prost(string, tag = "1")]
    pub address: ::prost::alloc::string::String,
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LastSubmittedEvmEventResponse {
    #[prost(uint64, tag = "1")]
    pub event_nonce: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Erc20ToDenomRequest {
    #[prost(string, tag = "1")]
    pub erc20: ::prost::alloc::string::String,
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Erc20ToDenomResponse {
    #[prost(string, tag = "1")]
    pub denom: ::prost::alloc::string::String,
    #[prost(bool, tag = "2")]
    pub cosmos_originated: bool,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DenomToErc20ParamsRequest {
    #[prost(string, tag = "1")]
    pub denom: ::prost::alloc::string::String,
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DenomToErc20ParamsResponse {
    #[prost(string, tag = "1")]
    pub base_denom: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub erc20_name: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub erc20_symbol: ::prost::alloc::string::String,
    #[prost(uint64, tag = "4")]
    pub erc20_decimals: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DenomToErc20Request {
    #[prost(string, tag = "1")]
    pub denom: ::prost::alloc::string::String,
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DenomToErc20Response {
    #[prost(string, tag = "1")]
    pub erc20: ::prost::alloc::string::String,
    #[prost(bool, tag = "2")]
    pub cosmos_originated: bool,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysByValidatorRequest {
    #[prost(string, tag = "1")]
    pub validator_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysByValidatorResponse {
    #[prost(string, tag = "1")]
    pub eth_address: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub orchestrator_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysByEvmSignerRequest {
    #[prost(string, tag = "1")]
    pub evm_signer: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysByEvmSignerResponse {
    #[prost(string, tag = "1")]
    pub validator_address: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub orchestrator_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysByOrchestratorRequest {
    #[prost(string, tag = "1")]
    pub orchestrator_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysByOrchestratorResponse {
    #[prost(string, tag = "1")]
    pub validator_address: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub evm_signer: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysRequest {}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysResponse {
    #[prost(message, repeated, tag = "1")]
    pub delegate_keys: ::prost::alloc::vec::Vec<MsgDelegateKeys>,
}
/// NOTE: if there is no sender address, return all
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchedSendToEvMsRequest {
    #[prost(string, tag = "1")]
    pub sender_address: ::prost::alloc::string::String,
    /// todo: figure out how to paginate given n Batches with m Send To EVMs
    ///  cosmos.base.query.v1beta1.PageRequest pagination = 2;
    #[prost(uint32, tag = "2")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchedSendToEvMsResponse {
    ///  cosmos.base.query.v1beta1.PageResponse pagination = 2;
    #[prost(message, repeated, tag = "1")]
    pub send_to_evms: ::prost::alloc::vec::Vec<SendToEvm>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnbatchedSendToEvMsRequest {
    #[prost(string, tag = "1")]
    pub sender_address: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "2")]
    pub pagination: ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageRequest>,
    #[prost(uint32, tag = "3")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UnbatchedSendToEvMsResponse {
    #[prost(message, repeated, tag = "1")]
    pub send_to_evms: ::prost::alloc::vec::Vec<SendToEvm>,
    #[prost(message, optional, tag = "2")]
    pub pagination:
        ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageResponse>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LastObservedEvmHeightRequest {
    #[prost(uint32, tag = "1")]
    pub chain_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LastObservedEvmHeightResponse {
    #[prost(message, optional, tag = "1")]
    pub last_observed_evm_height: ::core::option::Option<LatestEvmBlockHeight>,
}
#[doc = r" Generated client implementations."]
pub mod query_client {
    #![allow(unused_variables, dead_code, missing_docs)]
    use tonic::codegen::*;
    #[doc = " Query defines the gRPC querier service"]
    pub struct QueryClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl QueryClient<tonic::transport::Channel> {
        #[doc = r" Attempt to create a new client by connecting to a given endpoint."]
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: std::convert::TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> QueryClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::ResponseBody: Body + HttpBody + Send + 'static,
        T::Error: Into<StdError>,
        <T::ResponseBody as HttpBody>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_interceptor(inner: T, interceptor: impl Into<tonic::Interceptor>) -> Self {
            let inner = tonic::client::Grpc::with_interceptor(inner, interceptor);
            Self { inner }
        }
        #[doc = " Module parameters query"]
        pub async fn params(
            &mut self,
            request: impl tonic::IntoRequest<super::ParamsRequest>,
        ) -> Result<tonic::Response<super::ParamsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/Params");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " get info on individual outgoing data"]
        pub async fn signer_set_tx(
            &mut self,
            request: impl tonic::IntoRequest<super::SignerSetTxRequest>,
        ) -> Result<tonic::Response<super::SignerSetTxResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/SignerSetTx");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn latest_signer_set_tx(
            &mut self,
            request: impl tonic::IntoRequest<super::LatestSignerSetTxRequest>,
        ) -> Result<tonic::Response<super::SignerSetTxResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/LatestSignerSetTx");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn batch_tx(
            &mut self,
            request: impl tonic::IntoRequest<super::BatchTxRequest>,
        ) -> Result<tonic::Response<super::BatchTxResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/BatchTx");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn contract_call_tx(
            &mut self,
            request: impl tonic::IntoRequest<super::ContractCallTxRequest>,
        ) -> Result<tonic::Response<super::ContractCallTxResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/ContractCallTx");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " get collections of outgoing traffic from the bridge"]
        pub async fn signer_set_txs(
            &mut self,
            request: impl tonic::IntoRequest<super::SignerSetTxsRequest>,
        ) -> Result<tonic::Response<super::SignerSetTxsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/SignerSetTxs");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn batch_txs(
            &mut self,
            request: impl tonic::IntoRequest<super::BatchTxsRequest>,
        ) -> Result<tonic::Response<super::BatchTxsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/BatchTxs");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn contract_call_txs(
            &mut self,
            request: impl tonic::IntoRequest<super::ContractCallTxsRequest>,
        ) -> Result<tonic::Response<super::ContractCallTxsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/ContractCallTxs");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " TODO: can/should we group these into one endpoint?"]
        pub async fn signer_set_tx_confirmations(
            &mut self,
            request: impl tonic::IntoRequest<super::SignerSetTxConfirmationsRequest>,
        ) -> Result<tonic::Response<super::SignerSetTxConfirmationsResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/gravity.v2.Query/SignerSetTxConfirmations");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn batch_tx_confirmations(
            &mut self,
            request: impl tonic::IntoRequest<super::BatchTxConfirmationsRequest>,
        ) -> Result<tonic::Response<super::BatchTxConfirmationsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/gravity.v2.Query/BatchTxConfirmations");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn contract_call_tx_confirmations(
            &mut self,
            request: impl tonic::IntoRequest<super::ContractCallTxConfirmationsRequest>,
        ) -> Result<tonic::Response<super::ContractCallTxConfirmationsResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/gravity.v2.Query/ContractCallTxConfirmations",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " pending EVM signature queries for orchestrators to figure out which"]
        #[doc = " signatures they are missing"]
        #[doc = " TODO: can/should we group this into one endpoint?"]
        pub async fn unsigned_signer_set_txs(
            &mut self,
            request: impl tonic::IntoRequest<super::UnsignedSignerSetTxsRequest>,
        ) -> Result<tonic::Response<super::UnsignedSignerSetTxsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/gravity.v2.Query/UnsignedSignerSetTxs");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn unsigned_batch_txs(
            &mut self,
            request: impl tonic::IntoRequest<super::UnsignedBatchTxsRequest>,
        ) -> Result<tonic::Response<super::UnsignedBatchTxsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/UnsignedBatchTxs");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn unsigned_contract_call_txs(
            &mut self,
            request: impl tonic::IntoRequest<super::UnsignedContractCallTxsRequest>,
        ) -> Result<tonic::Response<super::UnsignedContractCallTxsResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/gravity.v2.Query/UnsignedContractCallTxs");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn last_submitted_evm_event(
            &mut self,
            request: impl tonic::IntoRequest<super::LastSubmittedEvmEventRequest>,
        ) -> Result<tonic::Response<super::LastSubmittedEvmEventResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/gravity.v2.Query/LastSubmittedEVMEvent");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " Queries the fees for all pending batches, results are returned in sdk.Coin"]
        #[doc = " (fee_amount_int)(contract_address) style"]
        pub async fn batch_tx_fees(
            &mut self,
            request: impl tonic::IntoRequest<super::BatchTxFeesRequest>,
        ) -> Result<tonic::Response<super::BatchTxFeesResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/BatchTxFees");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " Query for info about denoms tracked by gravity"]
        pub async fn erc20_to_denom(
            &mut self,
            request: impl tonic::IntoRequest<super::Erc20ToDenomRequest>,
        ) -> Result<tonic::Response<super::Erc20ToDenomResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/ERC20ToDenom");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " DenomToERC20Params implements a query that allows ERC-20 parameter"]
        #[doc = " information to be retrieved by a Cosmos base denomination."]
        pub async fn denom_to_erc20_params(
            &mut self,
            request: impl tonic::IntoRequest<super::DenomToErc20ParamsRequest>,
        ) -> Result<tonic::Response<super::DenomToErc20ParamsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/DenomToERC20Params");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " Query for info about denoms tracked by gravity"]
        pub async fn denom_to_erc20(
            &mut self,
            request: impl tonic::IntoRequest<super::DenomToErc20Request>,
        ) -> Result<tonic::Response<super::DenomToErc20Response>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/DenomToERC20");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " Query for batch send to EVMs"]
        pub async fn batched_send_to_ev_ms(
            &mut self,
            request: impl tonic::IntoRequest<super::BatchedSendToEvMsRequest>,
        ) -> Result<tonic::Response<super::BatchedSendToEvMsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/BatchedSendToEVMs");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " Query for unbatched send to EVMs"]
        pub async fn unbatched_send_to_ev_ms(
            &mut self,
            request: impl tonic::IntoRequest<super::UnbatchedSendToEvMsRequest>,
        ) -> Result<tonic::Response<super::UnbatchedSendToEvMsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/gravity.v2.Query/UnbatchedSendToEVMs");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " delegate keys"]
        pub async fn delegate_keys_by_validator(
            &mut self,
            request: impl tonic::IntoRequest<super::DelegateKeysByValidatorRequest>,
        ) -> Result<tonic::Response<super::DelegateKeysByValidatorResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/gravity.v2.Query/DelegateKeysByValidator");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn delegate_keys_by_evm_signer(
            &mut self,
            request: impl tonic::IntoRequest<super::DelegateKeysByEvmSignerRequest>,
        ) -> Result<tonic::Response<super::DelegateKeysByEvmSignerResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/gravity.v2.Query/DelegateKeysByEVMSigner");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn delegate_keys_by_orchestrator(
            &mut self,
            request: impl tonic::IntoRequest<super::DelegateKeysByOrchestratorRequest>,
        ) -> Result<tonic::Response<super::DelegateKeysByOrchestratorResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/gravity.v2.Query/DelegateKeysByOrchestrator",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn delegate_keys(
            &mut self,
            request: impl tonic::IntoRequest<super::DelegateKeysRequest>,
        ) -> Result<tonic::Response<super::DelegateKeysResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/gravity.v2.Query/DelegateKeys");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn last_observed_evm_height(
            &mut self,
            request: impl tonic::IntoRequest<super::LastObservedEvmHeightRequest>,
        ) -> Result<tonic::Response<super::LastObservedEvmHeightResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/gravity.v2.Query/LastObservedEVMHeight");
            self.inner.unary(request.into_request(), path, codec).await
        }
    }
    impl<T: Clone> Clone for QueryClient<T> {
        fn clone(&self) -> Self {
            Self {
                inner: self.inner.clone(),
            }
        }
    }
    impl<T> std::fmt::Debug for QueryClient<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "QueryClient {{ ... }}")
        }
    }
}
