import { uniqueArray } from "./unique-array";

test("uniqueArray", () => {
  expect(uniqueArray([])).toEqual([]);
  expect(uniqueArray([1, 2, 3])).toEqual([1, 2, 3]);
  expect(uniqueArray([1, 2, 2, 3])).toEqual([1, 2, 3]);
});
