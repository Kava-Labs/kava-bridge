/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import { Signer, utils, Contract, ContractFactory, Overrides } from "ethers";
import { Provider, TransactionRequest } from "@ethersproject/providers";
import type { Multicall, MulticallInterface } from "../Multicall";

const _abi = [
  {
    inputs: [
      {
        components: [
          {
            internalType: "address",
            name: "target",
            type: "address",
          },
          {
            internalType: "bytes",
            name: "callData",
            type: "bytes",
          },
        ],
        internalType: "struct Multicall.Call[]",
        name: "calls",
        type: "tuple[]",
      },
    ],
    name: "aggregate",
    outputs: [
      {
        internalType: "uint256",
        name: "blockNumber",
        type: "uint256",
      },
      {
        internalType: "bytes[]",
        name: "returnData",
        type: "bytes[]",
      },
    ],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "uint256",
        name: "blockNumber",
        type: "uint256",
      },
    ],
    name: "getBlockHash",
    outputs: [
      {
        internalType: "bytes32",
        name: "blockHash",
        type: "bytes32",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "getCurrentBlockCoinbase",
    outputs: [
      {
        internalType: "address",
        name: "coinbase",
        type: "address",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "getCurrentBlockDifficulty",
    outputs: [
      {
        internalType: "uint256",
        name: "difficulty",
        type: "uint256",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "getCurrentBlockGasLimit",
    outputs: [
      {
        internalType: "uint256",
        name: "gaslimit",
        type: "uint256",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "getCurrentBlockTimestamp",
    outputs: [
      {
        internalType: "uint256",
        name: "timestamp",
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
        name: "addr",
        type: "address",
      },
    ],
    name: "getEthBalance",
    outputs: [
      {
        internalType: "uint256",
        name: "balance",
        type: "uint256",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "getLastBlockHash",
    outputs: [
      {
        internalType: "bytes32",
        name: "blockHash",
        type: "bytes32",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
];

const _bytecode =
  "0x608060405234801561001057600080fd5b50610abb806100206000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c806372425d9d1161005b57806372425d9d1461012a57806386d516e814610148578063a8b0574e14610166578063ee82ac5e1461018457610088565b80630f28c97d1461008d578063252dba42146100ab57806327e86d6e146100dc5780634d2301cc146100fa575b600080fd5b6100956101b4565b6040516100a29190610381565b60405180910390f35b6100c560048036038101906100c091906106b0565b6101bc565b6040516100d3929190610843565b60405180910390f35b6100e461030f565b6040516100f1919061088c565b60405180910390f35b610114600480360381019061010f91906108a7565b610324565b6040516101219190610381565b60405180910390f35b610132610345565b60405161013f9190610381565b60405180910390f35b61015061034d565b60405161015d9190610381565b60405180910390f35b61016e610355565b60405161017b91906108e3565b60405180910390f35b61019e6004803603810190610199919061092a565b61035d565b6040516101ab919061088c565b60405180910390f35b600042905090565b60006060439150825167ffffffffffffffff8111156101de576101dd6103c6565b5b60405190808252806020026020018201604052801561021157816020015b60608152602001906001900390816101fc5790505b50905060005b83518110156103095760008085838151811061023657610235610957565b5b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff1686848151811061026b5761026a610957565b5b60200260200101516020015160405161028491906109c2565b6000604051808303816000865af19150503d80600081146102c1576040519150601f19603f3d011682016040523d82523d6000602084013e6102c6565b606091505b5091509150816102d557600080fd5b808484815181106102e9576102e8610957565b5b60200260200101819052505050808061030190610a08565b915050610217565b50915091565b600060014361031e9190610a51565b40905090565b60008173ffffffffffffffffffffffffffffffffffffffff16319050919050565b600044905090565b600045905090565b600041905090565b600081409050919050565b6000819050919050565b61037b81610368565b82525050565b60006020820190506103966000830184610372565b92915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6103fe826103b5565b810181811067ffffffffffffffff8211171561041d5761041c6103c6565b5b80604052505050565b600061043061039c565b905061043c82826103f5565b919050565b600067ffffffffffffffff82111561045c5761045b6103c6565b5b602082029050602081019050919050565b600080fd5b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006104a78261047c565b9050919050565b6104b78161049c565b81146104c257600080fd5b50565b6000813590506104d4816104ae565b92915050565b600080fd5b600067ffffffffffffffff8211156104fa576104f96103c6565b5b610503826103b5565b9050602081019050919050565b82818337600083830152505050565b600061053261052d846104df565b610426565b90508281526020810184848401111561054e5761054d6104da565b5b610559848285610510565b509392505050565b600082601f830112610576576105756103b0565b5b813561058684826020860161051f565b91505092915050565b6000604082840312156105a5576105a4610472565b5b6105af6040610426565b905060006105bf848285016104c5565b600083015250602082013567ffffffffffffffff8111156105e3576105e2610477565b5b6105ef84828501610561565b60208301525092915050565b600061060e61060984610441565b610426565b905080838252602082019050602084028301858111156106315761063061046d565b5b835b8181101561067857803567ffffffffffffffff811115610656576106556103b0565b5b808601610663898261058f565b85526020850194505050602081019050610633565b5050509392505050565b600082601f830112610697576106966103b0565b5b81356106a78482602086016105fb565b91505092915050565b6000602082840312156106c6576106c56103a6565b5b600082013567ffffffffffffffff8111156106e4576106e36103ab565b5b6106f084828501610682565b91505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b8381101561075f578082015181840152602081019050610744565b8381111561076e576000848401525b50505050565b600061077f82610725565b6107898185610730565b9350610799818560208601610741565b6107a2816103b5565b840191505092915050565b60006107b98383610774565b905092915050565b6000602082019050919050565b60006107d9826106f9565b6107e38185610704565b9350836020820285016107f585610715565b8060005b85811015610831578484038952815161081285826107ad565b945061081d836107c1565b925060208a019950506001810190506107f9565b50829750879550505050505092915050565b60006040820190506108586000830185610372565b818103602083015261086a81846107ce565b90509392505050565b6000819050919050565b61088681610873565b82525050565b60006020820190506108a1600083018461087d565b92915050565b6000602082840312156108bd576108bc6103a6565b5b60006108cb848285016104c5565b91505092915050565b6108dd8161049c565b82525050565b60006020820190506108f860008301846108d4565b92915050565b61090781610368565b811461091257600080fd5b50565b600081359050610924816108fe565b92915050565b6000602082840312156109405761093f6103a6565b5b600061094e84828501610915565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600081905092915050565b600061099c82610725565b6109a68185610986565b93506109b6818560208601610741565b80840191505092915050565b60006109ce8284610991565b915081905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610a1382610368565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610a4657610a456109d9565b5b600182019050919050565b6000610a5c82610368565b9150610a6783610368565b925082821015610a7a57610a796109d9565b5b82820390509291505056fea26469706673582212201fb995541984da8875a099a91c845b98e0f705584158106c92fbd78021d4e93864736f6c63430008090033";

type MulticallConstructorParams =
  | [signer?: Signer]
  | ConstructorParameters<typeof ContractFactory>;

const isSuperArgs = (
  xs: MulticallConstructorParams
): xs is ConstructorParameters<typeof ContractFactory> => xs.length > 1;

export class Multicall__factory extends ContractFactory {
  constructor(...args: MulticallConstructorParams) {
    if (isSuperArgs(args)) {
      super(...args);
    } else {
      super(_abi, _bytecode, args[0]);
    }
    this.contractName = "Multicall";
  }

  deploy(
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<Multicall> {
    return super.deploy(overrides || {}) as Promise<Multicall>;
  }
  getDeployTransaction(
    overrides?: Overrides & { from?: string | Promise<string> }
  ): TransactionRequest {
    return super.getDeployTransaction(overrides || {});
  }
  attach(address: string): Multicall {
    return super.attach(address) as Multicall;
  }
  connect(signer: Signer): Multicall__factory {
    return super.connect(signer) as Multicall__factory;
  }
  static readonly contractName: "Multicall";
  public readonly contractName: "Multicall";
  static readonly bytecode = _bytecode;
  static readonly abi = _abi;
  static createInterface(): MulticallInterface {
    return new utils.Interface(_abi) as MulticallInterface;
  }
  static connect(
    address: string,
    signerOrProvider: Signer | Provider
  ): Multicall {
    return new Contract(address, _abi, signerOrProvider) as Multicall;
  }
}
