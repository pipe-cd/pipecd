export const uniqueArray = <T>(arr: Array<T>): Array<T> =>
  [...new Set(arr)].sort();
