import { StyledEngineProvider, ThemeProvider as MuiThemeProvider } from "@mui/material";
import CssBaseline from "@mui/material/CssBaseline";
import { render } from "react-dom";
import { Routes } from "./routes";
import { BrowserRouter } from "react-router-dom";
import { setupDayjs } from "./utils/setup-dayjs";
import { CookiesProvider } from "react-cookie";
import QueryClientWrap from "./contexts/query-client-provider";
import { AuthProvider } from "./contexts/auth-context";
import { ToastProvider } from "./contexts/toast-context/toast-provider";
import { CommandProvider } from "./contexts/command-context";
import { lightTheme, darkTheme } from "./theme";
import { FC, useState, useEffect, ReactNode, createContext, useContext } from "react";

type ThemeMode = "light" | "dark";
const ThemeContext = createContext<{ mode: ThemeMode; toggleTheme: () => void } | null>(null);

export const useTheme = () => {
  const ctx = useContext(ThemeContext);
  if (!ctx) throw new Error("useTheme must be used within ThemeProvider");
  return ctx;
};

const AppWrapper: FC<{ children: ReactNode }> = ({ children }) => {
  const [mode, setMode] = useState<ThemeMode>(() => {
    if (typeof window !== "undefined") {
      const stored = localStorage.getItem("pipecd-theme") as ThemeMode | null;
      if (stored) return stored;
      return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
    }
    return "light";
  });

  useEffect(() => {
    localStorage.setItem("pipecd-theme", mode);
    if (typeof window !== "undefined") {
      document.documentElement.setAttribute("data-theme", mode);
    }
  }, [mode]);

  const theme = mode === "dark" ? darkTheme : lightTheme;

  return (
    <ThemeContext.Provider value={{ mode, toggleTheme: () => setMode(m => m === "light" ? "dark" : "light") }}>
      <MuiThemeProvider theme={theme}>
        <CssBaseline />
        {children}
      </MuiThemeProvider>
    </ThemeContext.Provider>
  );
};

async function run(): Promise<void> {
  if (process.env.ENABLE_MOCK === "true") {
    // NOTE: Ignore check exists this module, because this module exclude from production build.
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    const { worker } = await import("./mocks/browser");
    worker.start();
  }

  // Message to visitors.
  console.log(`
---------------------------------------------------
Hi there, fellow developer! Thanks for visiting.                ╓▄▓▓▓▓▓▓▓▄╖
As any other OSS projects out there, we highly                 ╙▀▓▓▓▓▓▓▓▓▓▓▓▄
appreciate your support. We seek for any kind               ▄▓▓▓▄ ▀▓▓▀▓▓▓▓▓▓▓▓
of contributions and feedback. If you feel                 ▐▓▓▓▓▓▓       ╟▓▓▓▓▓▌
interested, feel free to open up issues or PRs.            ▓▓▓▓▓▓         ▓▓▓▓▓▌
                                                           ▓▓▓▓▓▓▄        ╔φ, └└
The PipeCD official links:                                  ▓▓▓▓▓▓▀     ╠▒▒▒▒▒╠
  Documentation: https://pipecd.dev                           ▀▓▌ ╓▓▓▓▌ ╚▒╠^
  Github: https://github.com/pipe-cd/pipecd                     ▓▓▓▓▓▓▓▓,
  Twitter: https://twitter.com/pipecd_dev                        ╙▓▓▓▓▓▓▓▓µ
                                                                  └▓▓▓▓▓▓▓▌
Love to contribute to PipeCD? We're HIRING, so                      └▀▓▓▓╙
don't hesitate to ping us on GitHub or Twitter.

Happy PipeCD-ing 🙌
---------------------------------------------------
`);

  setupDayjs();

  // store.dispatch(fetchMe());

  render(
    <CookiesProvider>
      <StyledEngineProvider injectFirst>
        <AppWrapper>
          <BrowserRouter
            future={{
              v7_startTransition: false,
              v7_relativeSplatPath: false,
            }}
          >
            <ToastProvider>
              <QueryClientWrap>
                <AuthProvider>
                  <CommandProvider>
                    <Routes />
                  </CommandProvider>
                </AuthProvider>
              </QueryClientWrap>
            </ToastProvider>
          </BrowserRouter>
        </AppWrapper>
      </StyledEngineProvider>
    </CookiesProvider>,
    document.getElementById("root")
  );
}

run();
