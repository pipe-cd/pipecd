---
name: "module"
questions:
  name: "Please enter module name"
root: "src/modules"
output: "**/*"
ignore: []
---

# `{{ inputs.name }}.ts`

```ts
import { createSlice } from "@reduxjs/toolkit";

type {{ inputs.name | pascal }} = {};

const initialState: {{ inputs.name | pascal }} = {};

export const {{ inputs.name | camel }}Slice = createSlice({
  name: "{{ inputs.name | camel }}",
  initialState,
  reducers: {},
});
```

# `{{ inputs.name }}.test.ts`

```ts
import { {{ inputs.name | camel }}Slice } from "./{{ inputs.name }}";

describe("{{ inputs.name | camel }}Slice reducer", () => {
  it("should return the initial state", () => {
    expect(
      {{ inputs.name | camel }}Slice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot();
  });
});
```
