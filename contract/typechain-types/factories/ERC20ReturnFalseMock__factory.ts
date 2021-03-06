/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import { Signer, utils, Contract, ContractFactory, Overrides } from "ethers";
import { Provider, TransactionRequest } from "@ethersproject/providers";
import type {
  ERC20ReturnFalseMock,
  ERC20ReturnFalseMockInterface,
} from "../ERC20ReturnFalseMock";

const _abi = [
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
  "0x608060405234801561001057600080fd5b5061041a806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c8063095ea7b31461005157806323b872dd14610081578063a9059cbb146100b1578063dd62ed3e146100e1575b600080fd5b61006b60048036038101906100669190610234565b610111565b604051610078919061028f565b60405180910390f35b61009b600480360381019061009691906102aa565b610124565b6040516100a8919061028f565b60405180910390f35b6100cb60048036038101906100c69190610234565b610138565b6040516100d8919061028f565b60405180910390f35b6100fb60048036038101906100f691906102fd565b61014b565b604051610108919061034c565b60405180910390f35b6000806001819055506000905092915050565b600080600181905550600090509392505050565b6000806001819055506000905092915050565b60008060015414610191576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610188906103c4565b60405180910390fd5b6000905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006101cb826101a0565b9050919050565b6101db816101c0565b81146101e657600080fd5b50565b6000813590506101f8816101d2565b92915050565b6000819050919050565b610211816101fe565b811461021c57600080fd5b50565b60008135905061022e81610208565b92915050565b6000806040838503121561024b5761024a61019b565b5b6000610259858286016101e9565b925050602061026a8582860161021f565b9150509250929050565b60008115159050919050565b61028981610274565b82525050565b60006020820190506102a46000830184610280565b92915050565b6000806000606084860312156102c3576102c261019b565b5b60006102d1868287016101e9565b93505060206102e2868287016101e9565b92505060406102f38682870161021f565b9150509250925092565b600080604083850312156103145761031361019b565b5b6000610322858286016101e9565b9250506020610333858286016101e9565b9150509250929050565b610346816101fe565b82525050565b6000602082019050610361600083018461033d565b92915050565b600082825260208201905092915050565b7f64756d6d7920696e76616c696400000000000000000000000000000000000000600082015250565b60006103ae600d83610367565b91506103b982610378565b602082019050919050565b600060208201905081810360008301526103dd816103a1565b905091905056fea2646970667358221220558c19985958801a2323f51496624f043ef6daa2906a9e7ef78ceb7b85d4ee2164736f6c63430008090033";

type ERC20ReturnFalseMockConstructorParams =
  | [signer?: Signer]
  | ConstructorParameters<typeof ContractFactory>;

const isSuperArgs = (
  xs: ERC20ReturnFalseMockConstructorParams
): xs is ConstructorParameters<typeof ContractFactory> => xs.length > 1;

export class ERC20ReturnFalseMock__factory extends ContractFactory {
  constructor(...args: ERC20ReturnFalseMockConstructorParams) {
    if (isSuperArgs(args)) {
      super(...args);
    } else {
      super(_abi, _bytecode, args[0]);
    }
    this.contractName = "ERC20ReturnFalseMock";
  }

  deploy(
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<ERC20ReturnFalseMock> {
    return super.deploy(overrides || {}) as Promise<ERC20ReturnFalseMock>;
  }
  getDeployTransaction(
    overrides?: Overrides & { from?: string | Promise<string> }
  ): TransactionRequest {
    return super.getDeployTransaction(overrides || {});
  }
  attach(address: string): ERC20ReturnFalseMock {
    return super.attach(address) as ERC20ReturnFalseMock;
  }
  connect(signer: Signer): ERC20ReturnFalseMock__factory {
    return super.connect(signer) as ERC20ReturnFalseMock__factory;
  }
  static readonly contractName: "ERC20ReturnFalseMock";
  public readonly contractName: "ERC20ReturnFalseMock";
  static readonly bytecode = _bytecode;
  static readonly abi = _abi;
  static createInterface(): ERC20ReturnFalseMockInterface {
    return new utils.Interface(_abi) as ERC20ReturnFalseMockInterface;
  }
  static connect(
    address: string,
    signerOrProvider: Signer | Provider
  ): ERC20ReturnFalseMock {
    return new Contract(
      address,
      _abi,
      signerOrProvider
    ) as ERC20ReturnFalseMock;
  }
}
