import { StyledEngineProvider } from "@mui/material";
import { render } from "react-dom";
import { Routes } from "./routes";
import { BrowserRouter } from "react-router-dom";
import { setupDayjs } from "./utils/setup-dayjs";
import { CookiesProvider } from "react-cookie";
import QueryClientWrap from "./contexts/query-client-provider";
import { AuthProvider } from "./contexts/auth-context";
import { ToastProvider } from "./contexts/toast-context/toast-provider";
import { CommandProvider } from "./contexts/command-context";
import { ThemeProvider } from "./contexts/ThemeContext";

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
        <ThemeProvider>
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
        </ThemeProvider>
      </StyledEngineProvider>
    </CookiesProvider>,
    document.getElementById("root")
  );
}

run();
