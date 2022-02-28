module github.com/kava-labs/kava-bridge

go 1.17

require (
	github.com/cosmos/cosmos-sdk v0.45.0
	github.com/ethereum/go-ethereum v1.10.11
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.3.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.15
	github.com/tendermint/tm-db v0.6.6
	github.com/tharsis/ethermint v0.10.0-alpha1
)

require (
	github.com/cosmos/cosmos-proto v0.0.0-20211020182451-c7ca7198c2f8
	github.com/cosmos/ibc-go/v3 v3.0.0-beta1 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/gogo/protobuf v1.3.3
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	golang.org/x/crypto v0.0.0-20220213190939-1e6e3497d506 // indirect
	google.golang.org/genproto v0.0.0-20220211171837-173942840c17 // indirect
)

replace (
	github.com/99designs/keyring => github.com/cosmos/keyring v1.1.7-0.20210622111912-ef00f8ac3d76
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
