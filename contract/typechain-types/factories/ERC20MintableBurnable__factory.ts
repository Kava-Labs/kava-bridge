/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import {
  Signer,
  utils,
  Contract,
  ContractFactory,
  Overrides,
  BigNumberish,
} from "ethers";
import { Provider, TransactionRequest } from "@ethersproject/providers";
import type {
  ERC20MintableBurnable,
  ERC20MintableBurnableInterface,
} from "../ERC20MintableBurnable";

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
        internalType: "uint8",
        name: "decimals_",
        type: "uint8",
      },
    ],
    stateMutability: "nonpayable",
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
        name: "previousOwner",
        type: "address",
      },
      {
        indexed: true,
        internalType: "address",
        name: "newOwner",
        type: "address",
      },
    ],
    name: "OwnershipTransferred",
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
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "sender",
        type: "address",
      },
      {
        indexed: true,
        internalType: "address",
        name: "toAddr",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "sequence",
        type: "uint256",
      },
    ],
    name: "Withdraw",
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
        name: "to",
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
    name: "owner",
    outputs: [
      {
        internalType: "address",
        name: "",
        type: "address",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "renounceOwnership",
    outputs: [],
    stateMutability: "nonpayable",
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
        name: "newOwner",
        type: "address",
      },
    ],
    name: "transferOwnership",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "toAddr",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
    ],
    name: "withdraw",
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
      {
        internalType: "address",
        name: "toAddr",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
    ],
    name: "withdrawFrom",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
];

const _bytecode =
  "0x60a06040523480156200001157600080fd5b50604051620023f3380380620023f3833981810160405281019062000037919062000407565b600083838160039080519060200190620000539291906200017c565b5080600490805190602001906200006c9291906200017c565b5050506200008f62000083620000ae60201b60201c565b620000b660201b60201c565b80600681905550508060ff1660808160ff168152505050505062000506565b600033905090565b6000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b8280546200018a90620004d0565b90600052602060002090601f016020900481019282620001ae5760008555620001fa565b82601f10620001c957805160ff1916838001178555620001fa565b82800160010185558215620001fa579182015b82811115620001f9578251825591602001919060010190620001dc565b5b5090506200020991906200020d565b5090565b5b80821115620002285760008160009055506001016200020e565b5090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b62000295826200024a565b810181811067ffffffffffffffff82111715620002b757620002b66200025b565b5b80604052505050565b6000620002cc6200022c565b9050620002da82826200028a565b919050565b600067ffffffffffffffff821115620002fd57620002fc6200025b565b5b62000308826200024a565b9050602081019050919050565b60005b838110156200033557808201518184015260208101905062000318565b8381111562000345576000848401525b50505050565b6000620003626200035c84620002df565b620002c0565b90508281526020810184848401111562000381576200038062000245565b5b6200038e84828562000315565b509392505050565b600082601f830112620003ae57620003ad62000240565b5b8151620003c08482602086016200034b565b91505092915050565b600060ff82169050919050565b620003e181620003c9565b8114620003ed57600080fd5b50565b6000815190506200040181620003d6565b92915050565b60008060006060848603121562000423576200042262000236565b5b600084015167ffffffffffffffff8111156200044457620004436200023b565b5b620004528682870162000396565b935050602084015167ffffffffffffffff8111156200047657620004756200023b565b5b620004848682870162000396565b92505060406200049786828701620003f0565b9150509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620004e957607f821691505b602082108114156200050057620004ff620004a1565b5b50919050565b608051611ed16200052260003960006104620152611ed16000f3fe608060405234801561001057600080fd5b506004361061010b5760003560e01c8063715018a6116100a2578063a457c2d711610071578063a457c2d7146102a8578063a9059cbb146102d8578063dd62ed3e14610308578063f2fde38b14610338578063f3fef3a3146103545761010b565b8063715018a6146102465780638da5cb5b146102505780639555a9421461026e57806395d89b411461028a5761010b565b8063313ce567116100de578063313ce567146101ac57806339509351146101ca57806340c10f19146101fa57806370a08231146102165761010b565b806306fdde0314610110578063095ea7b31461012e57806318160ddd1461015e57806323b872dd1461017c575b600080fd5b610118610370565b6040516101259190611476565b60405180910390f35b61014860048036038101906101439190611531565b610402565b604051610155919061158c565b60405180910390f35b610166610425565b60405161017391906115b6565b60405180910390f35b610196600480360381019061019191906115d1565b61042f565b6040516101a3919061158c565b60405180910390f35b6101b461045e565b6040516101c19190611640565b60405180910390f35b6101e460048036038101906101df9190611531565b610486565b6040516101f1919061158c565b60405180910390f35b610214600480360381019061020f9190611531565b610530565b005b610230600480360381019061022b919061165b565b6105ba565b60405161023d91906115b6565b60405180910390f35b61024e610602565b005b61025861068a565b6040516102659190611697565b60405180910390f35b610288600480360381019061028391906115d1565b6106b4565b005b61029261073c565b60405161029f9190611476565b60405180910390f35b6102c260048036038101906102bd9190611531565b6107ce565b6040516102cf919061158c565b60405180910390f35b6102f260048036038101906102ed9190611531565b6108b8565b6040516102ff919061158c565b60405180910390f35b610322600480360381019061031d91906116b2565b6108db565b60405161032f91906115b6565b60405180910390f35b610352600480360381019061034d919061165b565b610962565b005b61036e60048036038101906103699190611531565b610a5a565b005b60606003805461037f90611721565b80601f01602080910402602001604051908101604052809291908181526020018280546103ab90611721565b80156103f85780601f106103cd576101008083540402835291602001916103f8565b820191906000526020600020905b8154815290600101906020018083116103db57829003601f168201915b5050505050905090565b60008061040d610ade565b905061041a818585610ae6565b600191505092915050565b6000600254905090565b60008061043a610ade565b9050610447858285610cb1565b610452858585610d3d565b60019150509392505050565b60007f0000000000000000000000000000000000000000000000000000000000000000905090565b600080610491610ade565b9050610525818585600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546105209190611782565b610ae6565b600191505092915050565b610538610ade565b73ffffffffffffffffffffffffffffffffffffffff1661055661068a565b73ffffffffffffffffffffffffffffffffffffffff16146105ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105a390611824565b60405180910390fd5b6105b68282610fbe565b5050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b61060a610ade565b73ffffffffffffffffffffffffffffffffffffffff1661062861068a565b73ffffffffffffffffffffffffffffffffffffffff161461067e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161067590611824565b60405180910390fd5b610688600061111e565b565b6000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6106bf833383610cb1565b6106c983826111e4565b8173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167ff341246adaac6f497bc2a656f546ab9e182111d630394f0c57c710a59a2cb567836107216113bb565b60405161072f929190611844565b60405180910390a3505050565b60606004805461074b90611721565b80601f016020809104026020016040519081016040528092919081815260200182805461077790611721565b80156107c45780601f10610799576101008083540402835291602001916107c4565b820191906000526020600020905b8154815290600101906020018083116107a757829003601f168201915b5050505050905090565b6000806107d9610ade565b90506000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508381101561089f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610896906118df565b60405180910390fd5b6108ac8286868403610ae6565b60019250505092915050565b6000806108c3610ade565b90506108d0818585610d3d565b600191505092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b61096a610ade565b73ffffffffffffffffffffffffffffffffffffffff1661098861068a565b73ffffffffffffffffffffffffffffffffffffffff16146109de576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109d590611824565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610a4e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a4590611971565b60405180910390fd5b610a578161111e565b50565b610a626113c5565b610a6c33826111e4565b8173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167ff341246adaac6f497bc2a656f546ab9e182111d630394f0c57c710a59a2cb56783610ac46113bb565b604051610ad2929190611844565b60405180910390a35050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610b56576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b4d90611a03565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610bc6576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bbd90611a95565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92583604051610ca491906115b6565b60405180910390a3505050565b6000610cbd84846108db565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610d375781811015610d29576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d2090611b01565b60405180910390fd5b610d368484848403610ae6565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610dad576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610da490611b93565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610e1d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e1490611c25565b60405180910390fd5b610e288383836113d3565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610eae576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ea590611cb7565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610f419190611782565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610fa591906115b6565b60405180910390a3610fb88484846113d8565b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141561102e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161102590611d23565b60405180910390fd5b61103a600083836113d3565b806002600082825461104c9190611782565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546110a19190611782565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161110691906115b6565b60405180910390a361111a600083836113d8565b5050565b6000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611254576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161124b90611db5565b60405180910390fd5b611260826000836113d3565b60008060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050818110156112e6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112dd90611e47565b60405180910390fd5b8181036000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816002600082825461133d9190611e67565b92505081905550600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516113a291906115b6565b60405180910390a36113b6836000846113d8565b505050565b6000600654905090565b600160065401600681905550565b505050565b505050565b600081519050919050565b600082825260208201905092915050565b60005b838110156114175780820151818401526020810190506113fc565b83811115611426576000848401525b50505050565b6000601f19601f8301169050919050565b6000611448826113dd565b61145281856113e8565b93506114628185602086016113f9565b61146b8161142c565b840191505092915050565b60006020820190508181036000830152611490818461143d565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006114c88261149d565b9050919050565b6114d8816114bd565b81146114e357600080fd5b50565b6000813590506114f5816114cf565b92915050565b6000819050919050565b61150e816114fb565b811461151957600080fd5b50565b60008135905061152b81611505565b92915050565b6000806040838503121561154857611547611498565b5b6000611556858286016114e6565b92505060206115678582860161151c565b9150509250929050565b60008115159050919050565b61158681611571565b82525050565b60006020820190506115a1600083018461157d565b92915050565b6115b0816114fb565b82525050565b60006020820190506115cb60008301846115a7565b92915050565b6000806000606084860312156115ea576115e9611498565b5b60006115f8868287016114e6565b9350506020611609868287016114e6565b925050604061161a8682870161151c565b9150509250925092565b600060ff82169050919050565b61163a81611624565b82525050565b60006020820190506116556000830184611631565b92915050565b60006020828403121561167157611670611498565b5b600061167f848285016114e6565b91505092915050565b611691816114bd565b82525050565b60006020820190506116ac6000830184611688565b92915050565b600080604083850312156116c9576116c8611498565b5b60006116d7858286016114e6565b92505060206116e8858286016114e6565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061173957607f821691505b6020821081141561174d5761174c6116f2565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061178d826114fb565b9150611798836114fb565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156117cd576117cc611753565b5b828201905092915050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b600061180e6020836113e8565b9150611819826117d8565b602082019050919050565b6000602082019050818103600083015261183d81611801565b9050919050565b600060408201905061185960008301856115a7565b61186660208301846115a7565b9392505050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b60006118c96025836113e8565b91506118d48261186d565b604082019050919050565b600060208201905081810360008301526118f8816118bc565b9050919050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b600061195b6026836113e8565b9150611966826118ff565b604082019050919050565b6000602082019050818103600083015261198a8161194e565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b60006119ed6024836113e8565b91506119f882611991565b604082019050919050565b60006020820190508181036000830152611a1c816119e0565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b6000611a7f6022836113e8565b9150611a8a82611a23565b604082019050919050565b60006020820190508181036000830152611aae81611a72565b9050919050565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000600082015250565b6000611aeb601d836113e8565b9150611af682611ab5565b602082019050919050565b60006020820190508181036000830152611b1a81611ade565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b6000611b7d6025836113e8565b9150611b8882611b21565b604082019050919050565b60006020820190508181036000830152611bac81611b70565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b6000611c0f6023836113e8565b9150611c1a82611bb3565b604082019050919050565b60006020820190508181036000830152611c3e81611c02565b9050919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b6000611ca16026836113e8565b9150611cac82611c45565b604082019050919050565b60006020820190508181036000830152611cd081611c94565b9050919050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b6000611d0d601f836113e8565b9150611d1882611cd7565b602082019050919050565b60006020820190508181036000830152611d3c81611d00565b9050919050565b7f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360008201527f7300000000000000000000000000000000000000000000000000000000000000602082015250565b6000611d9f6021836113e8565b9150611daa82611d43565b604082019050919050565b60006020820190508181036000830152611dce81611d92565b9050919050565b7f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60008201527f6365000000000000000000000000000000000000000000000000000000000000602082015250565b6000611e316022836113e8565b9150611e3c82611dd5565b604082019050919050565b60006020820190508181036000830152611e6081611e24565b9050919050565b6000611e72826114fb565b9150611e7d836114fb565b925082821015611e9057611e8f611753565b5b82820390509291505056fea26469706673582212205d5e8aa368667949106f8857889d21ed55477c897e78828f4b9ab44aaa86acfe64736f6c63430008090033";

type ERC20MintableBurnableConstructorParams =
  | [signer?: Signer]
  | ConstructorParameters<typeof ContractFactory>;

const isSuperArgs = (
  xs: ERC20MintableBurnableConstructorParams
): xs is ConstructorParameters<typeof ContractFactory> => xs.length > 1;

export class ERC20MintableBurnable__factory extends ContractFactory {
  constructor(...args: ERC20MintableBurnableConstructorParams) {
    if (isSuperArgs(args)) {
      super(...args);
    } else {
      super(_abi, _bytecode, args[0]);
    }
    this.contractName = "ERC20MintableBurnable";
  }

  deploy(
    name: string,
    symbol: string,
    decimals_: BigNumberish,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<ERC20MintableBurnable> {
    return super.deploy(
      name,
      symbol,
      decimals_,
      overrides || {}
    ) as Promise<ERC20MintableBurnable>;
  }
  getDeployTransaction(
    name: string,
    symbol: string,
    decimals_: BigNumberish,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): TransactionRequest {
    return super.getDeployTransaction(name, symbol, decimals_, overrides || {});
  }
  attach(address: string): ERC20MintableBurnable {
    return super.attach(address) as ERC20MintableBurnable;
  }
  connect(signer: Signer): ERC20MintableBurnable__factory {
    return super.connect(signer) as ERC20MintableBurnable__factory;
  }
  static readonly contractName: "ERC20MintableBurnable";
  public readonly contractName: "ERC20MintableBurnable";
  static readonly bytecode = _bytecode;
  static readonly abi = _abi;
  static createInterface(): ERC20MintableBurnableInterface {
    return new utils.Interface(_abi) as ERC20MintableBurnableInterface;
  }
  static connect(
    address: string,
    signerOrProvider: Signer | Provider
  ): ERC20MintableBurnable {
    return new Contract(
      address,
      _abi,
      signerOrProvider
    ) as ERC20MintableBurnable;
  }
}
