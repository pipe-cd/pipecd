export function isNullish(obj: unknown): boolean {
  if (obj === undefined || obj === null) {
    return true;
  }

  return false;
}
