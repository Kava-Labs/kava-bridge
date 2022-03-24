export const testKavaAddrs = [
  "kava1esagqd83rhqdtpy5sxhklaxgn58k2m3s3mnpea",
  "kava1mq9qxlhze029lm0frzw2xr6hem8c3k9ts54w0w",
  "kava16g8lzm86f5wwf3x3t67qrpd46sjdpxpfazskwg",
  "kava1wn74shl496ktcfgqsc6yf0vvenhgq0hwuw6z2a",
];

export const bech32Prefix = "kava";

export function tokens(
  value: bigint | number,
  decimals: bigint | number = 18
): bigint {
  return BigInt(value) * 10n ** BigInt(decimals);
}
