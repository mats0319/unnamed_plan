// @ts-ignore
import sha256 from "crypto-js/sha256";

export function calcSHA256(message: string): string {
  return sha256(message);
}
