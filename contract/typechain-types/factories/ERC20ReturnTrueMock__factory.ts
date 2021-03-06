/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import { Signer, utils, Contract, ContractFactory, Overrides } from "ethers";
import { Provider, TransactionRequest } from "@ethersproject/providers";
import type {
  ERC20ReturnTrueMock,
  ERC20ReturnTrueMockInterface,
} from "../ERC20ReturnTrueMock";

const _abi = [
  {
    inputs: [
      {
        internalType: "address",
        name: "owner",
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
        name: "",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "",
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
        internalType: "uint256",
        name: "allowance_",
        type: "uint256",
      },
    ],
    name: "setAllowance",
    outputs: [],
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
      {
        internalType: "uint256",
        name: "",
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
        name: "",
        type: "address",
      },
      {
        internalType: "address",
        name: "",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "",
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
];

const _bytecode =
  "0x608060405234801561001057600080fd5b50610430806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c8063095ea7b31461005c57806323b872dd1461008c5780633ba93f26146100bc578063a9059cbb146100d8578063dd62ed3e14610108575b600080fd5b6100766004803603810190610071919061029a565b610138565b60405161008391906102f5565b60405180910390f35b6100a660048036038101906100a19190610310565b61014b565b6040516100b391906102f5565b60405180910390f35b6100d660048036038101906100d19190610363565b61015f565b005b6100f260048036038101906100ed919061029a565b6101a5565b6040516100ff91906102f5565b60405180910390f35b610122600480360381019061011d9190610390565b6101b8565b60405161012f91906103df565b60405180910390f35b6000806001819055506001905092915050565b600080600181905550600190509392505050565b806000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555050565b6000806001819055506001905092915050565b60008060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061023182610206565b9050919050565b61024181610226565b811461024c57600080fd5b50565b60008135905061025e81610238565b92915050565b6000819050919050565b61027781610264565b811461028257600080fd5b50565b6000813590506102948161026e565b92915050565b600080604083850312156102b1576102b0610201565b5b60006102bf8582860161024f565b92505060206102d085828601610285565b9150509250929050565b60008115159050919050565b6102ef816102da565b82525050565b600060208201905061030a60008301846102e6565b92915050565b60008060006060848603121561032957610328610201565b5b60006103378682870161024f565b93505060206103488682870161024f565b925050604061035986828701610285565b9150509250925092565b60006020828403121561037957610378610201565b5b600061038784828501610285565b91505092915050565b600080604083850312156103a7576103a6610201565b5b60006103b58582860161024f565b92505060206103c68582860161024f565b9150509250929050565b6103d981610264565b82525050565b60006020820190506103f460008301846103d0565b9291505056fea2646970667358221220cfd360a922ee694cef47416d7befbbfda9e0f11f6a8f655266f860674fcafc3064736f6c63430008090033";

type ERC20ReturnTrueMockConstructorParams =
  | [signer?: Signer]
  | ConstructorParameters<typeof ContractFactory>;

const isSuperArgs = (
  xs: ERC20ReturnTrueMockConstructorParams
): xs is ConstructorParameters<typeof ContractFactory> => xs.length > 1;

export class ERC20ReturnTrueMock__factory extends ContractFactory {
  constructor(...args: ERC20ReturnTrueMockConstructorParams) {
    if (isSuperArgs(args)) {
      super(...args);
    } else {
      super(_abi, _bytecode, args[0]);
    }
    this.contractName = "ERC20ReturnTrueMock";
  }

  deploy(
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<ERC20ReturnTrueMock> {
    return super.deploy(overrides || {}) as Promise<ERC20ReturnTrueMock>;
  }
  getDeployTransaction(
    overrides?: Overrides & { from?: string | Promise<string> }
  ): TransactionRequest {
    return super.getDeployTransaction(overrides || {});
  }
  attach(address: string): ERC20ReturnTrueMock {
    return super.attach(address) as ERC20ReturnTrueMock;
  }
  connect(signer: Signer): ERC20ReturnTrueMock__factory {
    return super.connect(signer) as ERC20ReturnTrueMock__factory;
  }
  static readonly contractName: "ERC20ReturnTrueMock";
  public readonly contractName: "ERC20ReturnTrueMock";
  static readonly bytecode = _bytecode;
  static readonly abi = _abi;
  static createInterface(): ERC20ReturnTrueMockInterface {
    return new utils.Interface(_abi) as ERC20ReturnTrueMockInterface;
  }
  static connect(
    address: string,
    signerOrProvider: Signer | Provider
  ): ERC20ReturnTrueMock {
    return new Contract(address, _abi, signerOrProvider) as ERC20ReturnTrueMock;
  }
}
