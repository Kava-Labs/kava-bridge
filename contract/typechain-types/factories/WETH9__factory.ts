/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import { Signer, utils, Contract, ContractFactory, Overrides } from "ethers";
import { Provider, TransactionRequest } from "@ethersproject/providers";
import type { WETH9, WETH9Interface } from "../WETH9";

const _abi = [
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "src",
        type: "address",
      },
      {
        indexed: true,
        internalType: "address",
        name: "guy",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "wad",
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
        name: "dst",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "wad",
        type: "uint256",
      },
    ],
    name: "Deposit",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "src",
        type: "address",
      },
      {
        indexed: true,
        internalType: "address",
        name: "dst",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "wad",
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
        name: "src",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "wad",
        type: "uint256",
      },
    ],
    name: "Withdrawal",
    type: "event",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "",
        type: "address",
      },
      {
        internalType: "address",
        name: "",
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
        name: "guy",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "wad",
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
        name: "",
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
    inputs: [],
    name: "deposit",
    outputs: [],
    stateMutability: "payable",
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
        name: "dst",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "wad",
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
        name: "src",
        type: "address",
      },
      {
        internalType: "address",
        name: "dst",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "wad",
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
        internalType: "uint256",
        name: "wad",
        type: "uint256",
      },
    ],
    name: "withdraw",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    stateMutability: "payable",
    type: "receive",
  },
];

const _bytecode =
  "0x60806040526040518060400160405280600d81526020017f57726170706564204574686572000000000000000000000000000000000000008152506000908051906020019062000051929190620000d0565b506040518060400160405280600481526020017f5745544800000000000000000000000000000000000000000000000000000000815250600190805190602001906200009f929190620000d0565b506012600260006101000a81548160ff021916908360ff160217905550348015620000c957600080fd5b50620001e5565b828054620000de90620001af565b90600052602060002090601f0160209004810192826200010257600085556200014e565b82601f106200011d57805160ff19168380011785556200014e565b828001600101855582156200014e579182015b828111156200014d57825182559160200191906001019062000130565b5b5090506200015d919062000161565b5090565b5b808211156200017c57600081600090555060010162000162565b5090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620001c857607f821691505b60208210811415620001df57620001de62000180565b5b50919050565b610eeb80620001f56000396000f3fe6080604052600436106100a05760003560e01c8063313ce56711610064578063313ce567146101ad57806370a08231146101d857806395d89b4114610215578063a9059cbb14610240578063d0e30db01461027d578063dd62ed3e14610287576100af565b806306fdde03146100b4578063095ea7b3146100df57806318160ddd1461011c57806323b872dd146101475780632e1a7d4d14610184576100af565b366100af576100ad6102c4565b005b600080fd5b3480156100c057600080fd5b506100c961036a565b6040516100d69190610b1c565b60405180910390f35b3480156100eb57600080fd5b5061010660048036038101906101019190610bd7565b6103f8565b6040516101139190610c32565b60405180910390f35b34801561012857600080fd5b506101316104ea565b60405161013e9190610c5c565b60405180910390f35b34801561015357600080fd5b5061016e60048036038101906101699190610c77565b6104f2565b60405161017b9190610c32565b60405180910390f35b34801561019057600080fd5b506101ab60048036038101906101a69190610cca565b610856565b005b3480156101b957600080fd5b506101c2610990565b6040516101cf9190610d13565b60405180910390f35b3480156101e457600080fd5b506101ff60048036038101906101fa9190610d2e565b6109a3565b60405161020c9190610c5c565b60405180910390f35b34801561022157600080fd5b5061022a6109bb565b6040516102379190610b1c565b60405180910390f35b34801561024c57600080fd5b5061026760048036038101906102629190610bd7565b610a49565b6040516102749190610c32565b60405180910390f35b6102856102c4565b005b34801561029357600080fd5b506102ae60048036038101906102a99190610d5b565b610a5e565b6040516102bb9190610c5c565b60405180910390f35b34600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546103139190610dca565b925050819055503373ffffffffffffffffffffffffffffffffffffffff167fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c346040516103609190610c5c565b60405180910390a2565b6000805461037790610e4f565b80601f01602080910402602001604051908101604052809291908181526020018280546103a390610e4f565b80156103f05780601f106103c5576101008083540402835291602001916103f0565b820191906000526020600020905b8154815290600101906020018083116103d357829003601f168201915b505050505081565b600081600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516104d89190610c5c565b60405180910390a36001905092915050565b600047905090565b600081600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561054057600080fd5b3373ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161415801561061857507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205414155b1561073a5781600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156106a657600080fd5b81600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546107329190610e81565b925050819055505b81600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546107899190610e81565b9250508190555081600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546107df9190610dca565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516108439190610c5c565b60405180910390a3600190509392505050565b80600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156108a257600080fd5b80600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546108f19190610e81565b925050819055503373ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f1935050505015801561093e573d6000803e3d6000fd5b503373ffffffffffffffffffffffffffffffffffffffff167f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65826040516109859190610c5c565b60405180910390a250565b600260009054906101000a900460ff1681565b60036020528060005260406000206000915090505481565b600180546109c890610e4f565b80601f01602080910402602001604051908101604052809291908181526020018280546109f490610e4f565b8015610a415780601f10610a1657610100808354040283529160200191610a41565b820191906000526020600020905b815481529060010190602001808311610a2457829003601f168201915b505050505081565b6000610a563384846104f2565b905092915050565b6004602052816000526040600020602052806000526040600020600091509150505481565b600081519050919050565b600082825260208201905092915050565b60005b83811015610abd578082015181840152602081019050610aa2565b83811115610acc576000848401525b50505050565b6000601f19601f8301169050919050565b6000610aee82610a83565b610af88185610a8e565b9350610b08818560208601610a9f565b610b1181610ad2565b840191505092915050565b60006020820190508181036000830152610b368184610ae3565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610b6e82610b43565b9050919050565b610b7e81610b63565b8114610b8957600080fd5b50565b600081359050610b9b81610b75565b92915050565b6000819050919050565b610bb481610ba1565b8114610bbf57600080fd5b50565b600081359050610bd181610bab565b92915050565b60008060408385031215610bee57610bed610b3e565b5b6000610bfc85828601610b8c565b9250506020610c0d85828601610bc2565b9150509250929050565b60008115159050919050565b610c2c81610c17565b82525050565b6000602082019050610c476000830184610c23565b92915050565b610c5681610ba1565b82525050565b6000602082019050610c716000830184610c4d565b92915050565b600080600060608486031215610c9057610c8f610b3e565b5b6000610c9e86828701610b8c565b9350506020610caf86828701610b8c565b9250506040610cc086828701610bc2565b9150509250925092565b600060208284031215610ce057610cdf610b3e565b5b6000610cee84828501610bc2565b91505092915050565b600060ff82169050919050565b610d0d81610cf7565b82525050565b6000602082019050610d286000830184610d04565b92915050565b600060208284031215610d4457610d43610b3e565b5b6000610d5284828501610b8c565b91505092915050565b60008060408385031215610d7257610d71610b3e565b5b6000610d8085828601610b8c565b9250506020610d9185828601610b8c565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610dd582610ba1565b9150610de083610ba1565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115610e1557610e14610d9b565b5b828201905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610e6757607f821691505b60208210811415610e7b57610e7a610e20565b5b50919050565b6000610e8c82610ba1565b9150610e9783610ba1565b925082821015610eaa57610ea9610d9b565b5b82820390509291505056fea2646970667358221220dc0057e2a8e66a9ab69f1eb678c0b6a617e39261f95c832c326d12cc0498fb1c64736f6c63430008090033";

type WETH9ConstructorParams =
  | [signer?: Signer]
  | ConstructorParameters<typeof ContractFactory>;

const isSuperArgs = (
  xs: WETH9ConstructorParams
): xs is ConstructorParameters<typeof ContractFactory> => xs.length > 1;

export class WETH9__factory extends ContractFactory {
  constructor(...args: WETH9ConstructorParams) {
    if (isSuperArgs(args)) {
      super(...args);
    } else {
      super(_abi, _bytecode, args[0]);
    }
    this.contractName = "WETH9";
  }

  deploy(
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<WETH9> {
    return super.deploy(overrides || {}) as Promise<WETH9>;
  }
  getDeployTransaction(
    overrides?: Overrides & { from?: string | Promise<string> }
  ): TransactionRequest {
    return super.getDeployTransaction(overrides || {});
  }
  attach(address: string): WETH9 {
    return super.attach(address) as WETH9;
  }
  connect(signer: Signer): WETH9__factory {
    return super.connect(signer) as WETH9__factory;
  }
  static readonly contractName: "WETH9";
  public readonly contractName: "WETH9";
  static readonly bytecode = _bytecode;
  static readonly abi = _abi;
  static createInterface(): WETH9Interface {
    return new utils.Interface(_abi) as WETH9Interface;
  }
  static connect(address: string, signerOrProvider: Signer | Provider): WETH9 {
    return new Contract(address, _abi, signerOrProvider) as WETH9;
  }
}
