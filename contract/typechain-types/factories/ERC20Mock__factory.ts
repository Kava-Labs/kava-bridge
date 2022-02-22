/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import {
  Signer,
  utils,
  Contract,
  ContractFactory,
  PayableOverrides,
  BigNumberish,
} from "ethers";
import { Provider, TransactionRequest } from "@ethersproject/providers";
import type { ERC20Mock, ERC20MockInterface } from "../ERC20Mock";

const _abi = [
  {
    inputs: [
      {
        internalType: "string",
        name: "name",
        type: "string",
      },
      {
        internalType: "string",
        name: "symbol",
        type: "string",
      },
      {
        internalType: "address",
        name: "initialAccount",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "initialBalance",
        type: "uint256",
      },
    ],
    stateMutability: "payable",
    type: "constructor",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "owner",
        type: "address",
      },
      {
        indexed: true,
        internalType: "address",
        name: "spender",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "value",
        type: "uint256",
      },
    ],
    name: "Approval",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "from",
        type: "address",
      },
      {
        indexed: true,
        internalType: "address",
        name: "to",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "value",
        type: "uint256",
      },
    ],
    name: "Transfer",
    type: "event",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "owner",
        type: "address",
      },
      {
        internalType: "address",
        name: "spender",
        type: "address",
      },
    ],
    name: "allowance",
    outputs: [
      {
        internalType: "uint256",
        name: "",
        type: "uint256",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "spender",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
    ],
    name: "approve",
    outputs: [
      {
        internalType: "bool",
        name: "",
        type: "bool",
      },
    ],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "owner",
        type: "address",
      },
      {
        internalType: "address",
        name: "spender",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "value",
        type: "uint256",
      },
    ],
    name: "approveInternal",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "account",
        type: "address",
      },
    ],
    name: "balanceOf",
    outputs: [
      {
        internalType: "uint256",
        name: "",
        type: "uint256",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "account",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
    ],
    name: "burn",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "decimals",
    outputs: [
      {
        internalType: "uint8",
        name: "",
        type: "uint8",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "spender",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "subtractedValue",
        type: "uint256",
      },
    ],
    name: "decreaseAllowance",
    outputs: [
      {
        internalType: "bool",
        name: "",
        type: "bool",
      },
    ],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "spender",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "addedValue",
        type: "uint256",
      },
    ],
    name: "increaseAllowance",
    outputs: [
      {
        internalType: "bool",
        name: "",
        type: "bool",
      },
    ],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "account",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
    ],
    name: "mint",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "name",
    outputs: [
      {
        internalType: "string",
        name: "",
        type: "string",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "symbol",
    outputs: [
      {
        internalType: "string",
        name: "",
        type: "string",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "totalSupply",
    outputs: [
      {
        internalType: "uint256",
        name: "",
        type: "uint256",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "to",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
    ],
    name: "transfer",
    outputs: [
      {
        internalType: "bool",
        name: "",
        type: "bool",
      },
    ],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "from",
        type: "address",
      },
      {
        internalType: "address",
        name: "to",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
    ],
    name: "transferFrom",
    outputs: [
      {
        internalType: "bool",
        name: "",
        type: "bool",
      },
    ],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "from",
        type: "address",
      },
      {
        internalType: "address",
        name: "to",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "value",
        type: "uint256",
      },
    ],
    name: "transferInternal",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
];

const _bytecode =
  "0x6080604052604051620023403803806200234083398181016040528101906200002991906200058a565b83838160039080519060200190620000439291906200029d565b5080600490805190602001906200005c9291906200029d565b505050620000937fd19264395b48414e799b74dff30f3ea679f30ada74f4969e43fdad06504e8a7360001b6200011760201b60201c565b620000c77fe59d975e2d1bb701ebda22c96061b04dc37a0de0e68008a63cbb1e16ced6962e60001b6200011760201b60201c565b620000fb7fb71b3b4b16f0fbc9ce69ec1bfec441cc2a22f0e88c5c9fa684c002c7d0db4b2b60001b6200011760201b60201c565b6200010d82826200011a60201b60201c565b50505050620007dc565b50565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614156200018d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040162000184906200069b565b60405180910390fd5b620001a1600083836200029360201b60201c565b8060026000828254620001b59190620006ec565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546200020c9190620006ec565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516200027391906200075a565b60405180910390a36200028f600083836200029860201b60201c565b5050565b505050565b505050565b828054620002ab90620007a6565b90600052602060002090601f016020900481019282620002cf57600085556200031b565b82601f10620002ea57805160ff19168380011785556200031b565b828001600101855582156200031b579182015b828111156200031a578251825591602001919060010190620002fd565b5b5090506200032a91906200032e565b5090565b5b80821115620003495760008160009055506001016200032f565b5090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b620003b6826200036b565b810181811067ffffffffffffffff82111715620003d857620003d76200037c565b5b80604052505050565b6000620003ed6200034d565b9050620003fb8282620003ab565b919050565b600067ffffffffffffffff8211156200041e576200041d6200037c565b5b62000429826200036b565b9050602081019050919050565b60005b838110156200045657808201518184015260208101905062000439565b8381111562000466576000848401525b50505050565b6000620004836200047d8462000400565b620003e1565b905082815260208101848484011115620004a257620004a162000366565b5b620004af84828562000436565b509392505050565b600082601f830112620004cf57620004ce62000361565b5b8151620004e18482602086016200046c565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200051782620004ea565b9050919050565b62000529816200050a565b81146200053557600080fd5b50565b60008151905062000549816200051e565b92915050565b6000819050919050565b62000564816200054f565b81146200057057600080fd5b50565b600081519050620005848162000559565b92915050565b60008060008060808587031215620005a757620005a662000357565b5b600085015167ffffffffffffffff811115620005c857620005c76200035c565b5b620005d687828801620004b7565b945050602085015167ffffffffffffffff811115620005fa57620005f96200035c565b5b6200060887828801620004b7565b93505060406200061b8782880162000538565b92505060606200062e8782880162000573565b91505092959194509250565b600082825260208201905092915050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b600062000683601f836200063a565b915062000690826200064b565b602082019050919050565b60006020820190508181036000830152620006b68162000674565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000620006f9826200054f565b915062000706836200054f565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156200073e576200073d620006bd565b5b828201905092915050565b62000754816200054f565b82525050565b600060208201905062000771600083018462000749565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620007bf57607f821691505b60208210811415620007d657620007d562000777565b5b50919050565b611b5480620007ec6000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c806340c10f19116100975780639dc29fac116100665780639dc29fac14610286578063a457c2d7146102a2578063a9059cbb146102d2578063dd62ed3e14610302576100f5565b806340c10f191461020057806356189cb41461021c57806370a082311461023857806395d89b4114610268576100f5565b8063222f5be0116100d3578063222f5be01461016657806323b872dd14610182578063313ce567146101b257806339509351146101d0576100f5565b806306fdde03146100fa578063095ea7b31461011857806318160ddd14610148575b600080fd5b610102610332565b60405161010f919061124a565b60405180910390f35b610132600480360381019061012d9190611305565b6103c4565b60405161013f9190611360565b60405180910390f35b6101506103e7565b60405161015d919061138a565b60405180910390f35b610180600480360381019061017b91906113a5565b6103f1565b005b61019c600480360381019061019791906113a5565b610485565b6040516101a99190611360565b60405180910390f35b6101ba6104b4565b6040516101c79190611414565b60405180910390f35b6101ea60048036038101906101e59190611305565b6104bd565b6040516101f79190611360565b60405180910390f35b61021a60048036038101906102159190611305565b610567565b005b610236600480360381019061023191906113a5565b6105f9565b005b610252600480360381019061024d919061142f565b61068d565b60405161025f919061138a565b60405180910390f35b6102706106d5565b60405161027d919061124a565b60405180910390f35b6102a0600480360381019061029b9190611305565b610767565b005b6102bc60048036038101906102b79190611305565b6107f9565b6040516102c99190611360565b60405180910390f35b6102ec60048036038101906102e79190611305565b6108e3565b6040516102f99190611360565b60405180910390f35b61031c6004803603810190610317919061145c565b610906565b604051610329919061138a565b60405180910390f35b606060038054610341906114cb565b80601f016020809104026020016040519081016040528092919081815260200182805461036d906114cb565b80156103ba5780601f1061038f576101008083540402835291602001916103ba565b820191906000526020600020905b81548152906001019060200180831161039d57829003601f168201915b5050505050905090565b6000806103cf61098d565b90506103dc818585610995565b600191505092915050565b6000600254905090565b61041d7fcf002acf6a989f6241f6b861c50947109da2fb0414f2eff92dc603363fbf2b3b60001b610b60565b6104497fa7821e1f29db7ca3f47dfb1f459e02ec34458ffcdd96861fe789972bb6f533ac60001b610b60565b6104757f16110f38d4659b155a13edd306f0953524d2d33e0815f846a84316d736644c1860001b610b60565b610480838383610b63565b505050565b60008061049061098d565b905061049d858285610de4565b6104a8858585610b63565b60019150509392505050565b60006012905090565b6000806104c861098d565b905061055c818585600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610557919061152c565b610995565b600191505092915050565b6105937f3c614c12bd6b4dcd3c7dcf27df0e90c7fbbcdb96c1b3388c1b8f40afb1aed80260001b610b60565b6105bf7f4193708e7d2122d09a9e556518f3ec7836929c384a1bc093b69e3449a54814a560001b610b60565b6105eb7f8787ff5d5841094fb16a476f489a178914b198a873c83cf28d3d68ac005af98b60001b610b60565b6105f58282610e70565b5050565b6106257ff1deb146022f3d3d1fe1c57b0d5fbf35a31bf8486ccae9377c86ecec320f611e60001b610b60565b6106517f7f7d686a8b0bfd41a6d594192594a5522ba8aeb1e8142fda7234042422be4ff360001b610b60565b61067d7f48983a74522d68448120b611c55ecc336767742b456b0942cde0669e3518491e60001b610b60565b610688838383610995565b505050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6060600480546106e4906114cb565b80601f0160208091040260200160405190810160405280929190818152602001828054610710906114cb565b801561075d5780601f106107325761010080835404028352916020019161075d565b820191906000526020600020905b81548152906001019060200180831161074057829003601f168201915b5050505050905090565b6107937f6df74d9450eb5a9bd54e1cd3ddc6cf7d900560410283112d9085266d7ad4d77960001b610b60565b6107bf7f4cdd92dcd8e86636fc00acbeaa2b137e4b8b6ff91ec5cc4dbf27ac618a60934b60001b610b60565b6107eb7f2db3daee41a90603eff5f9b897ca37bddac063c48f35f77255d86fda24f3ce3760001b610b60565b6107f58282610fd0565b5050565b60008061080461098d565b90506000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050838110156108ca576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108c1906115f4565b60405180910390fd5b6108d78286868403610995565b60019250505092915050565b6000806108ee61098d565b90506108fb818585610b63565b600191505092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610a05576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109fc90611686565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610a75576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a6c90611718565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92583604051610b53919061138a565b60405180910390a3505050565b50565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610bd3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bca906117aa565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610c43576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c3a9061183c565b60405180910390fd5b610c4e8383836111a7565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610cd4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ccb906118ce565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610d67919061152c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610dcb919061138a565b60405180910390a3610dde8484846111ac565b50505050565b6000610df08484610906565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610e6a5781811015610e5c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e539061193a565b60405180910390fd5b610e698484848403610995565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610ee0576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ed7906119a6565b60405180910390fd5b610eec600083836111a7565b8060026000828254610efe919061152c565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610f53919061152c565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610fb8919061138a565b60405180910390a3610fcc600083836111ac565b5050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611040576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161103790611a38565b60405180910390fd5b61104c826000836111a7565b60008060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050818110156110d2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110c990611aca565b60405180910390fd5b8181036000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555081600260008282546111299190611aea565b92505081905550600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161118e919061138a565b60405180910390a36111a2836000846111ac565b505050565b505050565b505050565b600081519050919050565b600082825260208201905092915050565b60005b838110156111eb5780820151818401526020810190506111d0565b838111156111fa576000848401525b50505050565b6000601f19601f8301169050919050565b600061121c826111b1565b61122681856111bc565b93506112368185602086016111cd565b61123f81611200565b840191505092915050565b600060208201905081810360008301526112648184611211565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061129c82611271565b9050919050565b6112ac81611291565b81146112b757600080fd5b50565b6000813590506112c9816112a3565b92915050565b6000819050919050565b6112e2816112cf565b81146112ed57600080fd5b50565b6000813590506112ff816112d9565b92915050565b6000806040838503121561131c5761131b61126c565b5b600061132a858286016112ba565b925050602061133b858286016112f0565b9150509250929050565b60008115159050919050565b61135a81611345565b82525050565b60006020820190506113756000830184611351565b92915050565b611384816112cf565b82525050565b600060208201905061139f600083018461137b565b92915050565b6000806000606084860312156113be576113bd61126c565b5b60006113cc868287016112ba565b93505060206113dd868287016112ba565b92505060406113ee868287016112f0565b9150509250925092565b600060ff82169050919050565b61140e816113f8565b82525050565b60006020820190506114296000830184611405565b92915050565b6000602082840312156114455761144461126c565b5b6000611453848285016112ba565b91505092915050565b600080604083850312156114735761147261126c565b5b6000611481858286016112ba565b9250506020611492858286016112ba565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806114e357607f821691505b602082108114156114f7576114f661149c565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611537826112cf565b9150611542836112cf565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611577576115766114fd565b5b828201905092915050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b60006115de6025836111bc565b91506115e982611582565b604082019050919050565b6000602082019050818103600083015261160d816115d1565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b60006116706024836111bc565b915061167b82611614565b604082019050919050565b6000602082019050818103600083015261169f81611663565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b60006117026022836111bc565b915061170d826116a6565b604082019050919050565b60006020820190508181036000830152611731816116f5565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b60006117946025836111bc565b915061179f82611738565b604082019050919050565b600060208201905081810360008301526117c381611787565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b60006118266023836111bc565b9150611831826117ca565b604082019050919050565b6000602082019050818103600083015261185581611819565b9050919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b60006118b86026836111bc565b91506118c38261185c565b604082019050919050565b600060208201905081810360008301526118e7816118ab565b9050919050565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000600082015250565b6000611924601d836111bc565b915061192f826118ee565b602082019050919050565b6000602082019050818103600083015261195381611917565b9050919050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b6000611990601f836111bc565b915061199b8261195a565b602082019050919050565b600060208201905081810360008301526119bf81611983565b9050919050565b7f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360008201527f7300000000000000000000000000000000000000000000000000000000000000602082015250565b6000611a226021836111bc565b9150611a2d826119c6565b604082019050919050565b60006020820190508181036000830152611a5181611a15565b9050919050565b7f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60008201527f6365000000000000000000000000000000000000000000000000000000000000602082015250565b6000611ab46022836111bc565b9150611abf82611a58565b604082019050919050565b60006020820190508181036000830152611ae381611aa7565b9050919050565b6000611af5826112cf565b9150611b00836112cf565b925082821015611b1357611b126114fd565b5b82820390509291505056fea26469706673582212209ed425b37282588a43380fa908f61ba7c2dd45f1fb92f8129c82b85c666b3d0764736f6c63430008090033";

type ERC20MockConstructorParams =
  | [signer?: Signer]
  | ConstructorParameters<typeof ContractFactory>;

const isSuperArgs = (
  xs: ERC20MockConstructorParams
): xs is ConstructorParameters<typeof ContractFactory> => xs.length > 1;

export class ERC20Mock__factory extends ContractFactory {
  constructor(...args: ERC20MockConstructorParams) {
    if (isSuperArgs(args)) {
      super(...args);
    } else {
      super(_abi, _bytecode, args[0]);
    }
    this.contractName = "ERC20Mock";
  }

  deploy(
    name: string,
    symbol: string,
    initialAccount: string,
    initialBalance: BigNumberish,
    overrides?: PayableOverrides & { from?: string | Promise<string> }
  ): Promise<ERC20Mock> {
    return super.deploy(
      name,
      symbol,
      initialAccount,
      initialBalance,
      overrides || {}
    ) as Promise<ERC20Mock>;
  }
  getDeployTransaction(
    name: string,
    symbol: string,
    initialAccount: string,
    initialBalance: BigNumberish,
    overrides?: PayableOverrides & { from?: string | Promise<string> }
  ): TransactionRequest {
    return super.getDeployTransaction(
      name,
      symbol,
      initialAccount,
      initialBalance,
      overrides || {}
    );
  }
  attach(address: string): ERC20Mock {
    return super.attach(address) as ERC20Mock;
  }
  connect(signer: Signer): ERC20Mock__factory {
    return super.connect(signer) as ERC20Mock__factory;
  }
  static readonly contractName: "ERC20Mock";
  public readonly contractName: "ERC20Mock";
  static readonly bytecode = _bytecode;
  static readonly abi = _abi;
  static createInterface(): ERC20MockInterface {
    return new utils.Interface(_abi) as ERC20MockInterface;
  }
  static connect(
    address: string,
    signerOrProvider: Signer | Provider
  ): ERC20Mock {
    return new Contract(address, _abi, signerOrProvider) as ERC20Mock;
  }
}
