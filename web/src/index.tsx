import { ThemeProvider } from "@material-ui/core";
import CssBaseline from "@material-ui/core/CssBaseline";
import { render } from "react-dom";
import { theme } from "./theme";
import { Provider } from "react-redux";
import { store } from "./store";
import { Routes } from "./routes";
import { Router } from "react-router-dom";
import { history } from "./history";
import { setupDayjs } from "./utils/setup-dayjs";
import { fetchMe } from "./modules/me";

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
of contributions and feedbacks. If you feel                ▐▓▓▓▓▓▓       ╟▓▓▓▓▓▌
interested, feel free to open up issues or PRs.            ▓▓▓▓▓▓         ▓▓▓▓▓▌
                                                           ▓▓▓▓▓▓▄        ╔φ, └└
The PipeCD official site is located at                      ▓▓▓▓▓▓▀     ╠▒▒▒▒▒╠
          https://pipecd.dev                                  ▀▓▌ ╓▓▓▓▌ ╚▒╠^
                                                                ▓▓▓▓▓▓▓▓,
Love to contribute to PipeCD? we're HIRING, so                   ╙▓▓▓▓▓▓▓▓µ
don't hesitate to ping us on GitHub.                              └▓▓▓▓▓▓▓▌
                                                                    └▀▓▓▓╙
Happy PipeCD-ing 🙌
---------------------------------------------------
`);

  setupDayjs();

  store.dispatch(fetchMe());

  render(
    <Provider store={store}>
      <ThemeProvider theme={theme}>
        <Router history={history}>
          <CssBaseline />
          <Routes />
        </Router>
      </ThemeProvider>
    </Provider>,
    document.getElementById("root")
  );
}

run();
