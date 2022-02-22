/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import { Signer, utils, Contract, ContractFactory, Overrides } from "ethers";
import { Provider, TransactionRequest } from "@ethersproject/providers";
import type { Bridge, BridgeInterface } from "../Bridge";

const _abi = [
  {
    inputs: [
      {
        internalType: "address",
        name: "relayer_",
        type: "address",
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
        name: "token",
        type: "address",
      },
      {
        indexed: true,
        internalType: "address",
        name: "sender",
        type: "address",
      },
      {
        indexed: true,
        internalType: "bytes32",
        name: "toAddr",
        type: "bytes32",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
    ],
    name: "Lock",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "token",
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
    ],
    name: "Unlock",
    type: "event",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "token",
        type: "address",
      },
      {
        internalType: "bytes32",
        name: "toAddr",
        type: "bytes32",
      },
      {
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
    ],
    name: "lock",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "relayer",
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
    inputs: [
      {
        internalType: "address",
        name: "token",
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
    name: "unlock",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
];

const _bytecode =
  "0x60806040523480156200001157600080fd5b506040516200109c3803806200109c833981810160405281019062000037919062000187565b6200006b7fb38e40201fae00c85198b71edf8c12c3d1a8e4bd78b227ab49bad5bb81e9f32460001b6200011a60201b60201c565b6200009f7f7cbd9c1b6458f83a0c536305e67ea2203707f2b40d130fcf1f4a470c6f09dcd660001b6200011a60201b60201c565b620000d37fe90801a62791f2d72b99cefcab2d5c7981e3c96b36d9fef840bef209db95788460001b6200011a60201b60201c565b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050620001b9565b50565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200014f8262000122565b9050919050565b620001618162000142565b81146200016d57600080fd5b50565b600081519050620001818162000156565b92915050565b600060208284031215620001a0576200019f6200011d565b5b6000620001b08482850162000170565b91505092915050565b610ed380620001c96000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806359508f8f146100465780638406c07914610062578063a80de0e814610080575b600080fd5b610060600480360381019061005b9190610993565b61009c565b005b61006a61034b565b60405161007791906109f5565b60405180910390f35b61009a60048036038101906100959190610a46565b6103f7565b005b6100c87fcea84153626e4d3f04c44617ea8a2310fb7752d9695e2bf2bcaa9792ac62e2f160001b61056b565b6100f47f880d0dab9e23def9c01249a2c0f88b3fb387808ae6c353dafcccb77c0198baee60001b61056b565b6101207f5a90b144086b42e5291359b634b6158cddf62d33c222a42fec001f4eb9a6662760001b61056b565b61014c7ff32d8b6689206f3631e1b9a367a7ed2ebf381020354b67024551dbd6e3cede2d60001b61056b565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146101da576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101d190610af6565b60405180910390fd5b6102067f3115eb0771d6a94c953edafeca99ba1b8ecb4c7b2e3cb32ba960d95979f57fed60001b61056b565b6102327ff8264bad39654750ed649de9c699e75dfc033653eb38284b852b0d09ced02dbf60001b61056b565b61025e7f874f03ef2b864ea8c904afb964fadec174dd9dddc53bf752aa9cb2faa07dc17b60001b61056b565b61028982828573ffffffffffffffffffffffffffffffffffffffff1661056e9092919063ffffffff16565b6102b57f1948b110ca6498b74330cd269a3adf70a80b8d339996caacb2ce5415ea7bec2360001b61056b565b6102e17f37d23268d7ff756ef7b5f28fcfacb847e1016d388b36858e9872768c0145e9b260001b61056b565b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fc1c90b8e0705b212262c0dbd7580efe1862c2f185bf96899226f7596beb2db098360405161033e9190610b25565b60405180910390a3505050565b60006103797ff5b6a5942b4571fe42232ac1597e811c8d313a8ad3e6ad60c443287cdabdba0960001b61056b565b6103a57ff3f99ff75d22b2346d4388d886793791a5cccd470d4ab06fdb53b7b3a4afad7860001b61056b565b6103d17f8c8fb6d670cad30151947944414076468a74c0be281987a78897ae4fd58a8b9f60001b61056b565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6104237f0b7aaf62f695dac31c32d7b979072e21b6d3a6a47d9f2a83b09f176f2b38239060001b61056b565b61044f7f1626c908e61d76824dc2d21c22d7c2eeaf1bf68924975ca99fedb31741c8b40560001b61056b565b61047b7f46a649f7a16f50fbfb07c0fb7c793cd025299ae1e4bebfe040946f98db92076160001b61056b565b6104a83330838673ffffffffffffffffffffffffffffffffffffffff166105f4909392919063ffffffff16565b6104d47f6d590ced0ea3b5d2bf077776d7782e7ca9e3ed58416e47a4b83d4bc1f0674ae060001b61056b565b6105007f8738f2538db400214bc999ba8fd5e9dbc2a9d07f5820d8dc2407658ab417995d60001b61056b565b813373ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f825e58acf499949f9c703c8d538dff2a036b7b73d23228a8bdca9c5a63d397c88460405161055e9190610b25565b60405180910390a4505050565b50565b6105ef8363a9059cbb60e01b848460405160240161058d929190610b40565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505061067d565b505050565b610677846323b872dd60e01b85858560405160240161061593929190610b69565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505061067d565b50505050565b60006106df826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166107449092919063ffffffff16565b905060008151111561073f57808060200190518101906106ff9190610bd8565b61073e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161073590610c77565b60405180910390fd5b5b505050565b6060610753848460008561075c565b90509392505050565b6060824710156107a1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161079890610d09565b60405180910390fd5b6107aa85610870565b6107e9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107e090610d75565b60405180910390fd5b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516108129190610e0f565b60006040518083038185875af1925050503d806000811461084f576040519150601f19603f3d011682016040523d82523d6000602084013e610854565b606091505b5091509150610864828286610893565b92505050949350505050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b606083156108a3578290506108f3565b6000835111156108b65782518084602001fd5b816040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108ea9190610e7b565b60405180910390fd5b9392505050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061092a826108ff565b9050919050565b61093a8161091f565b811461094557600080fd5b50565b60008135905061095781610931565b92915050565b6000819050919050565b6109708161095d565b811461097b57600080fd5b50565b60008135905061098d81610967565b92915050565b6000806000606084860312156109ac576109ab6108fa565b5b60006109ba86828701610948565b93505060206109cb86828701610948565b92505060406109dc8682870161097e565b9150509250925092565b6109ef8161091f565b82525050565b6000602082019050610a0a60008301846109e6565b92915050565b6000819050919050565b610a2381610a10565b8114610a2e57600080fd5b50565b600081359050610a4081610a1a565b92915050565b600080600060608486031215610a5f57610a5e6108fa565b5b6000610a6d86828701610948565b9350506020610a7e86828701610a31565b9250506040610a8f8682870161097e565b9150509250925092565b600082825260208201905092915050565b7f4272696467653a20756e74727573746564206164647265737300000000000000600082015250565b6000610ae0601983610a99565b9150610aeb82610aaa565b602082019050919050565b60006020820190508181036000830152610b0f81610ad3565b9050919050565b610b1f8161095d565b82525050565b6000602082019050610b3a6000830184610b16565b92915050565b6000604082019050610b5560008301856109e6565b610b626020830184610b16565b9392505050565b6000606082019050610b7e60008301866109e6565b610b8b60208301856109e6565b610b986040830184610b16565b949350505050565b60008115159050919050565b610bb581610ba0565b8114610bc057600080fd5b50565b600081519050610bd281610bac565b92915050565b600060208284031215610bee57610bed6108fa565b5b6000610bfc84828501610bc3565b91505092915050565b7f5361666545524332303a204552433230206f7065726174696f6e20646964206e60008201527f6f74207375636365656400000000000000000000000000000000000000000000602082015250565b6000610c61602a83610a99565b9150610c6c82610c05565b604082019050919050565b60006020820190508181036000830152610c9081610c54565b9050919050565b7f416464726573733a20696e73756666696369656e742062616c616e636520666f60008201527f722063616c6c0000000000000000000000000000000000000000000000000000602082015250565b6000610cf3602683610a99565b9150610cfe82610c97565b604082019050919050565b60006020820190508181036000830152610d2281610ce6565b9050919050565b7f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000600082015250565b6000610d5f601d83610a99565b9150610d6a82610d29565b602082019050919050565b60006020820190508181036000830152610d8e81610d52565b9050919050565b600081519050919050565b600081905092915050565b60005b83811015610dc9578082015181840152602081019050610dae565b83811115610dd8576000848401525b50505050565b6000610de982610d95565b610df38185610da0565b9350610e03818560208601610dab565b80840191505092915050565b6000610e1b8284610dde565b915081905092915050565b600081519050919050565b6000601f19601f8301169050919050565b6000610e4d82610e26565b610e578185610a99565b9350610e67818560208601610dab565b610e7081610e31565b840191505092915050565b60006020820190508181036000830152610e958184610e42565b90509291505056fea264697066735822122055a49a67c87ffe4c74647f2ac2fe0695ad7627a8a1dd8a00f500ab91a0216a9264736f6c63430008090033";

type BridgeConstructorParams =
  | [signer?: Signer]
  | ConstructorParameters<typeof ContractFactory>;

const isSuperArgs = (
  xs: BridgeConstructorParams
): xs is ConstructorParameters<typeof ContractFactory> => xs.length > 1;

export class Bridge__factory extends ContractFactory {
  constructor(...args: BridgeConstructorParams) {
    if (isSuperArgs(args)) {
      super(...args);
    } else {
      super(_abi, _bytecode, args[0]);
    }
    this.contractName = "Bridge";
  }

  deploy(
    relayer_: string,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<Bridge> {
    return super.deploy(relayer_, overrides || {}) as Promise<Bridge>;
  }
  getDeployTransaction(
    relayer_: string,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): TransactionRequest {
    return super.getDeployTransaction(relayer_, overrides || {});
  }
  attach(address: string): Bridge {
    return super.attach(address) as Bridge;
  }
  connect(signer: Signer): Bridge__factory {
    return super.connect(signer) as Bridge__factory;
  }
  static readonly contractName: "Bridge";
  public readonly contractName: "Bridge";
  static readonly bytecode = _bytecode;
  static readonly abi = _abi;
  static createInterface(): BridgeInterface {
    return new utils.Interface(_abi) as BridgeInterface;
  }
  static connect(address: string, signerOrProvider: Signer | Provider): Bridge {
    return new Contract(address, _abi, signerOrProvider) as Bridge;
  }
}
