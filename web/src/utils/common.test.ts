import { sortFunc, sortDateFunc } from "./common";

test("sortFunc orders values in the requested direction", () => {
  expect(sortFunc("a", "b")).toBe(-1);
  expect(sortFunc("b", "a")).toBe(1);
  expect(sortFunc(1, 2)).toBe(-1);
  expect(sortFunc(2, 1)).toBe(1);

  expect(sortFunc("a", "b", "DESC")).toBe(1);
  expect(sortFunc("b", "a", "DESC")).toBe(-1);
});

test("sortFunc returns 0 for equal values", () => {
  expect(sortFunc("a", "a")).toBe(0);
  expect(sortFunc(1, 1)).toBe(0);
  expect(sortFunc("a", "a", "DESC")).toBe(0);
  expect(sortFunc(1, 1, "DESC")).toBe(0);
});

test("sortFunc keeps equal elements in their original order", () => {
  // Array.prototype.sort is only stable when the comparator returns 0 for
  // equal elements. Applications and pipeds can share a name, so this is the
  // ordering users actually see in the application forms.
  const items = [
    { name: "api", id: 1 },
    { name: "api", id: 2 },
    { name: "api", id: 3 },
    { name: "api", id: 4 },
    { name: "api", id: 5 },
    { name: "api", id: 6 },
    { name: "web", id: 7 },
    { name: "api", id: 8 },
    { name: "api", id: 9 },
    { name: "api", id: 10 },
    { name: "api", id: 11 },
    { name: "api", id: 12 },
  ];

  const sorted = [...items].sort((a, b) => sortFunc(a.name, b.name));

  expect(sorted.map((item) => item.id)).toEqual([
    1,
    2,
    3,
    4,
    5,
    6,
    8,
    9,
    10,
    11,
    12,
    7,
  ]);
});

test("sortDateFunc orders dates in the requested direction", () => {
  expect(sortDateFunc("2024-01-01", "2024-01-02")).toBe(-1);
  expect(sortDateFunc("2024-01-02", "2024-01-01")).toBe(1);
  expect(sortDateFunc("2024-01-01", "2024-01-02", "DESC")).toBe(1);
  expect(sortDateFunc("2024-01-02", "2024-01-01", "DESC")).toBe(-1);
});

test("sortDateFunc returns 0 for equal dates", () => {
  expect(sortDateFunc("2024/01/01", "2024/01/01")).toBe(0);
  expect(sortDateFunc("2024/01/01", "2024/01/01", "DESC")).toBe(0);
});
