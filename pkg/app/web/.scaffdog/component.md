---
name: "comp"
questions:
  name: "Please enter a component name."
root: "src/components"
output: "**/*"
ignore: []
---

# `{{ inputs.name }}/index.ts`

```tsx
export { {{ inputs.name | pascal }} } from "./{{ inputs.name }}";
```

# `{{ inputs.name }}/{{ inputs.name }}.tsx`

```tsx
import React, { FC } from "react";
import { makeStyles } from "@material-ui/core";

const useStyles = makeStyles(() => ({}));

interface Props {
}

export const {{ inputs.name | pascal }}: FC<Props> = ({ }) => {
  const classes = useStyles();
  return (
    <div>
      {{ inputs.name }}
    </div>
  )
};
```

# `{{ inputs.name }}/{{ inputs.name }}.stories.tsx`

```tsx
import React from "react";
import { {{ inputs.name | pascal }} } from "./{{ inputs.name }}";

export default {
  title: "{{ inputs.name | pascal }}",
  component: {{ inputs.name | pascal }}
};

export const overview: React.FC = () => (
  <{{ inputs.name | pascal }} />
);
```
