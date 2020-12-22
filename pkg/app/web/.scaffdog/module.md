---
name: "module"
description: "module of description"
message: "Please enter module name"
root: "src/modules"
output: "**/*"
ignore: []
---

# `{{ input }}.ts`

```ts
import { createSlice } from "@reduxjs/toolkit";

type {{ input | pascal }} = {};

const initialState: {{ input | pascal }} = {};

export const {{ input | camel }}Slice = createSlice({
  name: "{{ input | camel }}",
  initialState,
  reducers: {},
});
```

# `{{ input }}.test.ts`

```ts
import { {{ input | camel }}Slice } from "./{{ input }}";

describe("{{ input | camel }}Slice reducer", () => {
  it("should return the initial state", () => {
    expect(
      {{ input | camel }}Slice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot();
  });
});
```
