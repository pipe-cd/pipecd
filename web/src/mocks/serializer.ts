export function serialize(message: Uint8Array): Uint8Array {
  let len = message.length;
  const bytesArray = [0, 0, 0, 0];
  const payload = new Uint8Array(len + 5);
  // Write message length as 32bit BE
  for (let i = 3; i >= 0; i--) {
    bytesArray[i] = len % 256;
    len = len >>> 8;
  }
  payload.set(new Uint8Array(bytesArray), 1);
  payload.set(message, 5);
  return payload;
}
