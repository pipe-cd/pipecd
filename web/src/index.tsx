import { ThemeProvider, StyledEngineProvider } from "@mui/material";
import CssBaseline from "@mui/material/CssBaseline";
import { render } from "react-dom";
import { useMemo, useState } from "react";
import { createAppTheme, ThemeContext, ThemeMode } from "./theme";
import { Routes } from "./routes";
import { BrowserRouter } from "react-router-dom";
import { setupDayjs } from "./utils/setup-dayjs";
import { CookiesProvider } from "react-cookie";
import QueryClientWrap from "./contexts/query-client-provider";
import { AuthProvider } from "./contexts/auth-context";
import { ToastProvider } from "./contexts/toast-context/toast-provider";
import { CommandProvider } from "./contexts/command-context";

// The root component — lives here so we can use React state for the theme toggle.
function App() {
  // Pick up whatever the user last chose, or fall back to their OS preference
  const savedMode = localStorage.getItem("theme_mode") as ThemeMode | null;
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const [mode, setMode] = useState<ThemeMode>(
    savedMode ?? (prefersDark ? "dark" : "light")
  );

  const toggleTheme = () => {
    setMode((current) => {
      const next = current === "light" ? "dark" : "light";
      localStorage.setItem("theme_mode", next); // remember across sessions
      return next;
    });
  };

  // Rebuild the MUI theme only when mode actually changes
  const appTheme = useMemo(() => createAppTheme(mode), [mode]);

  return (
    <ThemeContext.Provider value={{ mode, toggleTheme }}>
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={appTheme}>
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
                    <CssBaseline />
                    <Routes />
                  </CommandProvider>
                </AuthProvider>
              </QueryClientWrap>
            </ToastProvider>
          </BrowserRouter>
        </ThemeProvider>
      </StyledEngineProvider>
    </ThemeContext.Provider>
  );
}

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
      <App />
    </CookiesProvider>,
    document.getElementById("root")
  );
}

run();
