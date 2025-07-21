import { ThemeProvider, StyledEngineProvider } from "@mui/material";
import { render, RenderOptions, RenderResult } from "@testing-library/react";
import { theme } from "~/theme";
import MemoryRouterTest from "./MemoryRouterTest";
import QueryClientWrap from "~/contexts/query-client-provider";
import { ToastProvider } from "~/contexts/toast-context";

const customRender = (
  ui: React.ReactElement,
  renderOptions: Omit<RenderOptions, "queries"> = {}
): RenderResult => {
  const Wrapper: React.ComponentType = ({ children }) => (
    <StyledEngineProvider injectFirst>
      <ToastProvider>
        <QueryClientWrap>
          <ThemeProvider theme={theme}>{children}</ThemeProvider>
        </QueryClientWrap>
      </ToastProvider>
    </StyledEngineProvider>
  );
  return render(ui, { wrapper: Wrapper, ...renderOptions });
};

// re-export everything
export * from "@testing-library/react";
// override render method
export { customRender as render };

export { MemoryRouterTest as MemoryRouter };
