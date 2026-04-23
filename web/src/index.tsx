import { ThemeProvider, StyledEngineProvider } from "@mui/material";
import CssBaseline from "@mui/material/CssBaseline";
import { createRoot } from "react-dom/client";
import { theme } from "./theme";
import { Routes } from "./routes";
import { BrowserRouter } from "react-router-dom";
import { setupDayjs } from "./utils/setup-dayjs";
import { CookiesProvider } from "react-cookie";
import QueryClientWrap from "./contexts/query-client-provider";
import { AuthProvider } from "./contexts/auth-context";
import { ToastProvider } from "./contexts/toast-context/toast-provider";
import { CommandProvider } from "./contexts/command-context";

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

  const container = document.getElementById("root");
  if (!container) {
    throw new Error("Root element not found");
  }
  const root = createRoot(container);
  root.render(
    <CookiesProvider>
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={theme}>
          <BrowserRouter>
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
    </CookiesProvider>
  );
}

run();
