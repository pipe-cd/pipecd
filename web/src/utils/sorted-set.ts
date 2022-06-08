export const sortedSet = <T>(arr: Array<T>): Array<T> =>
  [...new Set(arr)].sort();
