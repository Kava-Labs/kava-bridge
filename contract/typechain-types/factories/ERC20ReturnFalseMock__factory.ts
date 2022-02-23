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
  "0x608060405234801561001057600080fd5b506107e9806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c8063095ea7b31461005157806323b872dd14610081578063a9059cbb146100b1578063dd62ed3e146100e1575b600080fd5b61006b60048036038101906100669190610603565b610111565b604051610078919061065e565b60405180910390f35b61009b60048036038101906100969190610679565b610201565b6040516100a8919061065e565b60405180910390f35b6100cb60048036038101906100c69190610603565b6102f2565b6040516100d8919061065e565b60405180910390f35b6100fb60048036038101906100f691906106cc565b6103e2565b604051610108919061071b565b60405180910390f35b600061013f7f7471d140fc6d460795e86a7cf6ed2ecf418e45115b38a3306891c84032dc5bba60001b610567565b61016b7f13f257225ced87a5f788e549072b2e285a5dc560059759e4ce8d06b0d186a7eb60001b610567565b6101977f6d412df5eb496419ad76eec1b902c5b55b546a323d18502fe45e4255ea17747c60001b610567565b60006001819055506101cb7fc747b7bc0a5c0eab01a7688886e513bd602e6487ce1b7788f5a50d550f0332e260001b610567565b6101f77fbb9f5ce92b2e3931af61a004b24a8f777fdc83d311e9caea7adf9193ea35495060001b610567565b6000905092915050565b600061022f7f5c7df9630ec1a1c4b8d770af7658d75e835b05eff21ba188e2d995bc79f8020860001b610567565b61025b7fd345f22397fde5473fb6d3e7077bbdacc08979dfcbaf672436beb7e1ff159e4c60001b610567565b6102877fbd69c39c86da6f36a9548faf51e30c64dd962bc92093079da1be322bcd892f7060001b610567565b60006001819055506102bb7f0211f716a7576e993edf788c84abf40033f335c086645960b9403d7b9046d5ac60001b610567565b6102e77fb1379376875d8a7f7c789d2e2f7931e3179b8755e6b0b508ab11b0ec05c8b9b260001b610567565b600090509392505050565b60006103207f5dc1e21a6c35eefffd0831996440cda0d8d30dc07b47fcf7400ae15dda4e92af60001b610567565b61034c7f48db685d0b0f5f359c779e872d7c5a3df30b3220723a92682b1c2ddb449f7df660001b610567565b6103787fd43633b123dde5e7cf35f50d093ff146450157c48fccb250e6d16e6a7f8e6bdf60001b610567565b60006001819055506103ac7fd90a5a2caa3c40431259e31c39e10df0af628fbc6e7980365dfbe38167738a1860001b610567565b6103d87fd5e24586a2b137b6001c651a22d41563a0805733a9f8a430514024354a9b896760001b610567565b6000905092915050565b60006104107f2d9d5c465dcc5201870b4b79798377c85139cf9b48a2d934b0d7fd47a2eea81260001b610567565b61043c7f5200c85e1f0cc60293b78857358d6a369a316f08b10d8285b3c7c13541f6131860001b610567565b6104687fa52709296c4c6a1be876ae827d556579a92917fe32b1d9b3293c78f742b2579d60001b610567565b6104947f9ed736629898598e0277c5b8e2f78e9fe67257367581801c59ad1f5703e21dce60001b610567565b6000600154146104d9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104d090610793565b60405180910390fd5b6105057f72fd1cfcd339fcf22a6cf82dcad233a991723b4b4f2e50568a685e15898c32a360001b610567565b6105317fa978d8fc272b84a23c0b3bcd40bd51f599b2a27a396b74f77ec22c7336133fa360001b610567565b61055d7fc20cfe5cc727b73bb362c298b344c356f5758d45b8c10dcf1daa178e60a23c7c60001b610567565b6000905092915050565b50565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061059a8261056f565b9050919050565b6105aa8161058f565b81146105b557600080fd5b50565b6000813590506105c7816105a1565b92915050565b6000819050919050565b6105e0816105cd565b81146105eb57600080fd5b50565b6000813590506105fd816105d7565b92915050565b6000806040838503121561061a5761061961056a565b5b6000610628858286016105b8565b9250506020610639858286016105ee565b9150509250929050565b60008115159050919050565b61065881610643565b82525050565b6000602082019050610673600083018461064f565b92915050565b6000806000606084860312156106925761069161056a565b5b60006106a0868287016105b8565b93505060206106b1868287016105b8565b92505060406106c2868287016105ee565b9150509250925092565b600080604083850312156106e3576106e261056a565b5b60006106f1858286016105b8565b9250506020610702858286016105b8565b9150509250929050565b610715816105cd565b82525050565b6000602082019050610730600083018461070c565b92915050565b600082825260208201905092915050565b7f64756d6d7920696e76616c696400000000000000000000000000000000000000600082015250565b600061077d600d83610736565b915061078882610747565b602082019050919050565b600060208201905081810360008301526107ac81610770565b905091905056fea2646970667358221220855f7216f17240d7237bba5d873bd17bc7b20f29e66645a52ecd4aabf94d3d3764736f6c63430008090033";

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