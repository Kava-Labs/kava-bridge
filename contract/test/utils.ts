import { bech32 } from "bech32";
import { ethers } from "ethers";

export const testKavaAddrs = [
  "kava1esagqd83rhqdtpy5sxhklaxgn58k2m3s3mnpea",
  "kava1mq9qxlhze029lm0frzw2xr6hem8c3k9ts54w0w",
  "kava16g8lzm86f5wwf3x3t67qrpd46sjdpxpfazskwg",
  "kava1wn74shl496ktcfgqsc6yf0vvenhgq0hwuw6z2a",
];

export const baseDec = 1000000000000000000n; // 18 decimal places
export const bech32Prefix = "kava";

export function kavaAddrToBytes32(kavaAddr: string): Uint8Array {
  const decoded = bech32.decode(kavaAddr);

  if (decoded.prefix !== bech32Prefix) {
    throw new Error(`invalid prefix: ${decoded.prefix}`);
  }

  const bytes = bech32.fromWords(decoded.words);
  return ethers.utils.zeroPad(bytes, 32);
}

export function bytes32ToKavaAddr(addrBytes: Uint8Array): string {
  const encoded = bech32.toWords(addrBytes.slice(addrBytes.length - 20));
  return bech32.encode(bech32Prefix, encoded);
}

export function tokens(
  value: bigint | number,
  decimals: bigint | number = 18
): bigint {
  return BigInt(value) * 10n ** BigInt(decimals);
}
