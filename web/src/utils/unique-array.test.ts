import { sortedSet } from "./unique-array";

test("uniqueArray", () => {
  expect(sortedSet([])).toEqual([]);
  expect(sortedSet([1, 2, 3])).toEqual([1, 2, 3]);
  expect(sortedSet([1, 2, 2, 3])).toEqual([1, 2, 3]);
});
