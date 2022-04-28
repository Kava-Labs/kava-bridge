 <!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [bridge/v1beta1/conversion_pair.proto](#bridge/v1beta1/conversion_pair.proto)
    - [ConversionPair](#bridge.v1beta1.ConversionPair)
  
- [bridge/v1beta1/erc20.proto](#bridge/v1beta1/erc20.proto)
    - [ERC20BridgePair](#bridge.v1beta1.ERC20BridgePair)
  
- [bridge/v1beta1/genesis.proto](#bridge/v1beta1/genesis.proto)
    - [EnabledERC20Token](#bridge.v1beta1.EnabledERC20Token)
    - [GenesisState](#bridge.v1beta1.GenesisState)
    - [Params](#bridge.v1beta1.Params)
  
- [bridge/v1beta1/query.proto](#bridge/v1beta1/query.proto)
    - [QueryConversionPairRequest](#bridge.v1beta1.QueryConversionPairRequest)
    - [QueryConversionPairResponse](#bridge.v1beta1.QueryConversionPairResponse)
    - [QueryConversionPairsRequest](#bridge.v1beta1.QueryConversionPairsRequest)
    - [QueryConversionPairsResponse](#bridge.v1beta1.QueryConversionPairsResponse)
    - [QueryERC20BridgePairRequest](#bridge.v1beta1.QueryERC20BridgePairRequest)
    - [QueryERC20BridgePairResponse](#bridge.v1beta1.QueryERC20BridgePairResponse)
    - [QueryERC20BridgePairsRequest](#bridge.v1beta1.QueryERC20BridgePairsRequest)
    - [QueryERC20BridgePairsResponse](#bridge.v1beta1.QueryERC20BridgePairsResponse)
    - [QueryParamsRequest](#bridge.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#bridge.v1beta1.QueryParamsResponse)
  
    - [Query](#bridge.v1beta1.Query)
  
- [bridge/v1beta1/tx.proto](#bridge/v1beta1/tx.proto)
    - [MsgBridgeEthereumToKava](#bridge.v1beta1.MsgBridgeEthereumToKava)
    - [MsgBridgeEthereumToKavaResponse](#bridge.v1beta1.MsgBridgeEthereumToKavaResponse)
    - [MsgConvertCoinToERC20](#bridge.v1beta1.MsgConvertCoinToERC20)
    - [MsgConvertCoinToERC20Response](#bridge.v1beta1.MsgConvertCoinToERC20Response)
  
    - [Msg](#bridge.v1beta1.Msg)
  
- [relayer/v1beta1/broadcast_message.proto](#relayer/v1beta1/broadcast_message.proto)
    - [BroadcastMessage](#relayer.v1beta1.BroadcastMessage)
    - [Echo](#relayer.v1beta1.Echo)
  
- [relayer/v1beta1/message.proto](#relayer/v1beta1/message.proto)
    - [HelloRequest](#relayer.v1beta1.HelloRequest)
    - [MessageData](#relayer.v1beta1.MessageData)
  
- [Scalar Value Types](#scalar-value-types)



<a name="bridge/v1beta1/conversion_pair.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## bridge/v1beta1/conversion_pair.proto



<a name="bridge.v1beta1.ConversionPair"></a>

### ConversionPair
ConversionPair defines a Kava ERC20 address and corresponding denom that is
allowed to be converted between ERC20 and sdk.Coin


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `kava_erc20_address` | [bytes](#bytes) |  | ERC20 address of the token on the Kava EVM |
| `denom` | [string](#string) |  | Denom of the corresponding sdk.Coin |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



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
| `minimum_withdraw_amount` | [string](#string) |  | Minimum amount of the token that can be bridged back to Ethereum to prevent outgoing transfers that are much smaller than Ethereum gas costs. |






<a name="bridge.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the bridge module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#bridge.v1beta1.Params) |  | params defines all the parameters of the module. |
| `erc20_bridge_pairs` | [ERC20BridgePair](#bridge.v1beta1.ERC20BridgePair) | repeated | erc20_bridge_pairs defines all of the bridged erc20 tokens. |
| `next_withdraw_sequence` | [string](#string) |  | next_withdraw_sequence defines the unique incrementing sequence per withdraw tx. |






<a name="bridge.v1beta1.Params"></a>

### Params
Params defines the bridge module params


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bridge_enabled` | [bool](#bool) |  | Flag for enabling incoming/outgoing bridge transactions AND Kava ERC20/sdk.Coin conversions. |
| `enabled_erc20_tokens` | [EnabledERC20Token](#bridge.v1beta1.EnabledERC20Token) | repeated | List of ERC20Tokens that are allowed to be bridged to Kava |
| `relayer` | [bytes](#bytes) |  | Permissioned relayer address that is allowed to submit bridge messages |
| `enabled_conversion_pairs` | [ConversionPair](#bridge.v1beta1.ConversionPair) | repeated | enabled_conversion_pairs defines the list of conversion pairs allowed to be converted between Kava ERC20 and sdk.Coin |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="bridge/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## bridge/v1beta1/query.proto



<a name="bridge.v1beta1.QueryConversionPairRequest"></a>

### QueryConversionPairRequest
QueryConversionPairRequest defines the request type for querying a x/bridge conversion pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address_or_denom` | [string](#string) |  | AddressOrDenom defines the ERC20 address or the sdk.Coin denom of the pair to search for. |






<a name="bridge.v1beta1.QueryConversionPairResponse"></a>

### QueryConversionPairResponse
QueryConversionPairsResponse defines the response type for querying a x/bridge conversion pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `conversion_pair` | [ConversionPair](#bridge.v1beta1.ConversionPair) |  | ConversionPair defines the queried conversion pairs. |






<a name="bridge.v1beta1.QueryConversionPairsRequest"></a>

### QueryConversionPairsRequest
QueryConversionPairsRequest defines the request type for querying x/bridge conversion pairs.






<a name="bridge.v1beta1.QueryConversionPairsResponse"></a>

### QueryConversionPairsResponse
QueryConversionPairsResponse defines the response type for querying x/bridge conversion pairs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `conversion_pairs` | [ConversionPair](#bridge.v1beta1.ConversionPair) | repeated | ConversionPairs defines the queried conversion pairs. |






<a name="bridge.v1beta1.QueryERC20BridgePairRequest"></a>

### QueryERC20BridgePairRequest
QueryERC20BridgePairRequest defines the request type for querying x/bridge ERC20 pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | Address defines the internal or external address to query for. This is a string and not bytes as bytes in the query must be base64 encoded which is not ideal for addresses where we prefer hex encoding. |






<a name="bridge.v1beta1.QueryERC20BridgePairResponse"></a>

### QueryERC20BridgePairResponse
QueryERC20BridgePairRequest defines the response type for querying x/bridge ERC20 pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `erc20_bridge_pair` | [ERC20BridgePair](#bridge.v1beta1.ERC20BridgePair) |  | ERC20BridgePair defines the queried bridged erc20 pair. |






<a name="bridge.v1beta1.QueryERC20BridgePairsRequest"></a>

### QueryERC20BridgePairsRequest
QueryERC20BridgePairsRequest defines the request type for querying x/bridge ERC20 pairs.






<a name="bridge.v1beta1.QueryERC20BridgePairsResponse"></a>

### QueryERC20BridgePairsResponse
QueryERC20BridgePairsRequest defines the response type for querying x/bridge ERC20 pairs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `erc20_bridge_pairs` | [ERC20BridgePair](#bridge.v1beta1.ERC20BridgePair) | repeated | ERC20BridgePairs defines all of the currently bridged erc20 tokens. |






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
| `ERC20BridgePairs` | [QueryERC20BridgePairsRequest](#bridge.v1beta1.QueryERC20BridgePairsRequest) | [QueryERC20BridgePairsResponse](#bridge.v1beta1.QueryERC20BridgePairsResponse) | ERC20BridgePairs queries the bridge address pairs. | GET|/kava/bridge/v1beta1/bridge-erc20-pairs|
| `ERC20BridgePair` | [QueryERC20BridgePairRequest](#bridge.v1beta1.QueryERC20BridgePairRequest) | [QueryERC20BridgePairResponse](#bridge.v1beta1.QueryERC20BridgePairResponse) | ERC20BridgePair queries a bridge address pair with either internal or external address. | GET|/kava/bridge/v1beta1/bridge-erc20-pairs/{address}|
| `ConversionPairs` | [QueryConversionPairsRequest](#bridge.v1beta1.QueryConversionPairsRequest) | [QueryConversionPairsResponse](#bridge.v1beta1.QueryConversionPairsResponse) | ConversionPairs queries the ERC20/sdk.Coin conversion pairs. | GET|/kava/bridge/v1beta1/conversion-pairs|
| `ConversionPair` | [QueryConversionPairRequest](#bridge.v1beta1.QueryConversionPairRequest) | [QueryConversionPairResponse](#bridge.v1beta1.QueryConversionPairResponse) | ConversionPair queries a conversion pair with either the ERC20 address or sdk.Coin denom. | GET|/kava/bridge/v1beta1/conversion-pairs/{address_or_denom}|

 <!-- end services -->



<a name="bridge/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## bridge/v1beta1/tx.proto



<a name="bridge.v1beta1.MsgBridgeEthereumToKava"></a>

### MsgBridgeEthereumToKava
MsgBridgeEthereumToKava defines a ERC20 bridge transfer from Ethereum to Kava.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `relayer` | [string](#string) |  | Address of the bridge relayer. |
| `ethereum_erc20_address` | [string](#string) |  | Originating Ethereum ERC20 contract address. |
| `amount` | [string](#string) |  | ERC20 token amount to transfer. |
| `receiver` | [string](#string) |  | Receiver hex address on Kava. |
| `sequence` | [string](#string) |  | Unique sequence per bridge event. |






<a name="bridge.v1beta1.MsgBridgeEthereumToKavaResponse"></a>

### MsgBridgeEthereumToKavaResponse
MsgBridgeEthereumToKavaResponse defines the response value from
Msg/BridgeEthereumToKava.






<a name="bridge.v1beta1.MsgConvertCoinToERC20"></a>

### MsgConvertCoinToERC20
MsgConvertCoinToERC20 defines a conversion from sdk.Coin to Kava ERC20.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `initiator` | [string](#string) |  | Kava bech32 address initiating the conversion. |
| `receiver` | [string](#string) |  | EVM hex address that will receive the converted Kava ERC20 tokens. |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Amount is the sdk.Coin amount to convert. |






<a name="bridge.v1beta1.MsgConvertCoinToERC20Response"></a>

### MsgConvertCoinToERC20Response
MsgConvertCoinToERC20Response defines the response value from
Msg/ConvertCoinToERC20.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="bridge.v1beta1.Msg"></a>

### Msg
Msg defines the bridge Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `BridgeEthereumToKava` | [MsgBridgeEthereumToKava](#bridge.v1beta1.MsgBridgeEthereumToKava) | [MsgBridgeEthereumToKavaResponse](#bridge.v1beta1.MsgBridgeEthereumToKavaResponse) | BridgeEthereumToKava defines a method for bridging ERC20 tokens from Ethereum to Kava. | |
| `ConvertCoinToERC20` | [MsgConvertCoinToERC20](#bridge.v1beta1.MsgConvertCoinToERC20) | [MsgConvertCoinToERC20Response](#bridge.v1beta1.MsgConvertCoinToERC20Response) | ConvertCoinToERC20 defines a method for converting sdk.Coin to Kava ERC20. | |

 <!-- end services -->



<a name="relayer/v1beta1/broadcast_message.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## relayer/v1beta1/broadcast_message.proto



<a name="relayer.v1beta1.BroadcastMessage"></a>

### BroadcastMessage
BroadcastMessage is used between peers to wrap messages for each protocol


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | Unique ID of this message. |
| `recipient_peer_ids` | [string](#string) | repeated | Selected recipients of the message, to partially restrict the broadcast to a subset a peers. |
| `payload` | [google.protobuf.Any](#google.protobuf.Any) |  | Customtype workaround for not having to use a separate protocgen.sh script |
| `created` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | Timestamp when the message was broadcasted. |






<a name="relayer.v1beta1.Echo"></a>

### Echo
Echo is a message used for testing


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `message` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="relayer/v1beta1/message.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## relayer/v1beta1/message.proto



<a name="relayer.v1beta1.HelloRequest"></a>

### HelloRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `message` | [string](#string) |  |  |






<a name="relayer.v1beta1.MessageData"></a>

### MessageData
MessageData is used between peers to wrap messages for each protocol


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `payload` | [google.protobuf.Any](#google.protobuf.Any) |  | Customtype workaround for not having to use a separate protocgen.sh script |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

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
