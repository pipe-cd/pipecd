---
name: "comp"
description: "Generate a component file"
message: "Please enter a component name."
root: "src/components"
output: "**/*"
ignore: []
---

# `{{ input }}.tsx`

```tsx
import React, { FC } from "react";
import { makeStyles } from "@material-ui/core";

const useStyles = makeStyles(() => ({}));

interface Props {
}

export const {{ input | pascal }}: FC<Props> = ({ }) => {
  const classes = useStyles();
  return (
    <div>
      Hello
    </div>
  )
};
```

# `{{ input }}.stories.tsx`

```tsx
import React from "react";
import { {{ input | pascal }} } from "./{{ input }}";

export default {
  title: "{{ input | pascal }}",
  component: {{ input | pascal }}
};

export const overview: React.FC = () => (
  <{{ input | pascal }} />
);
```
