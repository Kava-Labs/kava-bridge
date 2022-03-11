 <!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [bridge/v1beta1/erc20.proto](#bridge/v1beta1/erc20.proto)
    - [ERC20BridgePair](#bridge.v1beta1.ERC20BridgePair)
  
- [bridge/v1beta1/genesis.proto](#bridge/v1beta1/genesis.proto)
    - [EnabledERC20Token](#bridge.v1beta1.EnabledERC20Token)
    - [GenesisState](#bridge.v1beta1.GenesisState)
    - [Params](#bridge.v1beta1.Params)
  
- [bridge/v1beta1/query.proto](#bridge/v1beta1/query.proto)
    - [QueryERC20BridgePairsRequest](#bridge.v1beta1.QueryERC20BridgePairsRequest)
    - [QueryERC20BridgePairsResponse](#bridge.v1beta1.QueryERC20BridgePairsResponse)
    - [QueryParamsRequest](#bridge.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#bridge.v1beta1.QueryParamsResponse)
  
    - [Query](#bridge.v1beta1.Query)
  
- [bridge/v1beta1/tx.proto](#bridge/v1beta1/tx.proto)
    - [MsgBridgeERC20FromEthereum](#bridge.v1beta1.MsgBridgeERC20FromEthereum)
    - [MsgBridgeERC20FromEthereumResponse](#bridge.v1beta1.MsgBridgeERC20FromEthereumResponse)
  
    - [Msg](#bridge.v1beta1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="bridge/v1beta1/erc20.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## bridge/v1beta1/erc20.proto



<a name="bridge.v1beta1.ERC20BridgePair"></a>

### ERC20BridgePair
ERC20BridgePair defines an ERC20 token bridged between external and Kava EVM


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `external_erc20_address` | [bytes](#bytes) |  | external_erc20_address represents the external EVM ERC20 address |
| `internal_erc20_address` | [bytes](#bytes) |  | internal_erc20_address represents the corresponding internal Kava EVM ERC20 address |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="bridge/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## bridge/v1beta1/genesis.proto



<a name="bridge.v1beta1.EnabledERC20Token"></a>

### EnabledERC20Token
EnabledERC20Token defines an external ERC20 that is allowed to be bridged to Kava


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [bytes](#bytes) |  | Address of the contract on Ethereum |
| `name` | [string](#string) |  | Name of the token. |
| `symbol` | [string](#string) |  | Symbol of the ERC20 token, usually a shorter version of the name. |
| `decimals` | [uint32](#uint32) |  | Number of decimals the ERC20 uses to get its user representation. The max value is an unsigned 8 bit integer, but is an uint32 as the smallest protobuf integer type. |






<a name="bridge.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the bridge module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#bridge.v1beta1.Params) |  | params defines all the parameters of the module. |
| `erc20_bridge_pairs` | [ERC20BridgePair](#bridge.v1beta1.ERC20BridgePair) | repeated | erc20_bridge_pairs defines all of the bridged erc20 tokens. |






<a name="bridge.v1beta1.Params"></a>

### Params
Params defines the bridge module params


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled_erc20_tokens` | [EnabledERC20Token](#bridge.v1beta1.EnabledERC20Token) | repeated | List of ERC20Tokens that are allowed to be bridged to Kava |
| `relayer` | [bytes](#bytes) |  | Permissioned relayer address that is allowed to submit bridge messages |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="bridge/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## bridge/v1beta1/query.proto



<a name="bridge.v1beta1.QueryERC20BridgePairsRequest"></a>

### QueryERC20BridgePairsRequest
QueryERC20BridgePairsRequest defines the request type for querying x/bridge ERC20 pairs.






<a name="bridge.v1beta1.QueryERC20BridgePairsResponse"></a>

### QueryERC20BridgePairsResponse
QueryERC20BridgePairsRequest defines the response type for querying x/bridge ERC20 pairs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `erc20_bridge_pairs` | [ERC20BridgePair](#bridge.v1beta1.ERC20BridgePair) | repeated | erc20_bridge_pairs defines all of the currently bridged erc20 tokens. |






<a name="bridge.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest defines the request type for querying x/bridge parameters.






<a name="bridge.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse defines the response type for querying x/bridge parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#bridge.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="bridge.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service for bridge module

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#bridge.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#bridge.v1beta1.QueryParamsResponse) | Params queries all parameters of the bridge module. | GET|/kava/bridge/v1beta1/params|
| `ERC20BridgePairs` | [QueryERC20BridgePairsRequest](#bridge.v1beta1.QueryERC20BridgePairsRequest) | [QueryERC20BridgePairsResponse](#bridge.v1beta1.QueryERC20BridgePairsResponse) | ERC20BridgePairs queries the bridge address pairs | GET|/kava/bridge/v1beta1/bridge-erc20-pairs|

 <!-- end services -->



<a name="bridge/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## bridge/v1beta1/tx.proto



<a name="bridge.v1beta1.MsgBridgeERC20FromEthereum"></a>

### MsgBridgeERC20FromEthereum
MsgBridgeERC20FromEthereum defines a ERC20 bridge transfer from Ethereum.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `relayer` | [string](#string) |  |  |
| `ethereum_erc20_address` | [string](#string) |  | Originating Ethereum ERC20 contract address |
| `amount` | [string](#string) |  | ERC20 token amount to transfer |
| `receiver` | [string](#string) |  | Receiver hex address on Kava |
| `sequence` | [string](#string) |  | Unique sequence per bridge event |






<a name="bridge.v1beta1.MsgBridgeERC20FromEthereumResponse"></a>

### MsgBridgeERC20FromEthereumResponse
MsgBridgeERC20FromEthereumResponse defines the response value from





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="bridge.v1beta1.Msg"></a>

### Msg
Msg defines the bridge Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `BridgeERC20FromEthereum` | [MsgBridgeERC20FromEthereum](#bridge.v1beta1.MsgBridgeERC20FromEthereum) | [MsgBridgeERC20FromEthereumResponse](#bridge.v1beta1.MsgBridgeERC20FromEthereumResponse) | BridgeERC20FromEthereum defines a method for bridging ERC20 tokens from Ethereum. | |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |
