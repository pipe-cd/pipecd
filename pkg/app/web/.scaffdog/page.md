---
name: "page"
description: "Generate a page component"
message: "Please enter page name"
root: "src/pages"
output: "**/*"
ignore: []
---

# `{{ input }}.tsx`

```tsx
import React, { memo, FC, useEffect } from "react";
import { useParams } from "react-router-dom";

export const {{ input | pascal }}Page: FC = memo(() => {
  return <div>hello</div>;
});
```
