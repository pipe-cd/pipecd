import { findMetadataByKey } from "./find-metadata-by-key";

test("findMetadataByKey", () => {
  expect(findMetadataByKey([], "key")).toBeUndefined();
  expect(findMetadataByKey([["key2", "value"]], "key")).toBeUndefined();
  expect(
    findMetadataByKey(
      [
        ["key", "value"],
        ["key2", "value2"],
      ],
      "key"
    )
  ).toBe("value");
});
