import { sortedSet } from "./sorted-set";

test("uniqueArray", () => {
  expect(sortedSet([])).toEqual([]);
  expect(sortedSet([1, 2, 3])).toEqual([1, 2, 3]);
  expect(sortedSet([1, 2, 2, 3])).toEqual([1, 2, 3]);
});
