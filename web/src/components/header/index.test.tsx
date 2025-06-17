import userEvent from "@testing-library/user-event";
import { MemoryRouter, render, screen, waitFor } from "~~/test-utils";
import { Header } from "./";
import { Cookies, CookiesProvider } from "react-cookie";
import { server } from "~/mocks/server";
import { AuthProvider } from "~/contexts/auth-context";

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

it("shows login link if user state is not exists", () => {
  render(
    <MemoryRouter>
      <Header />
    </MemoryRouter>,
    {}
  );

  expect(screen.getByRole("link", { name: "Login" })).toBeInTheDocument();
});

it("shows logout link if opened user menu", async () => {
  const cookies = new Cookies();
  cookies.set("token", "my-test-token");

  render(
    <MemoryRouter>
      <CookiesProvider cookies={cookies}>
        <AuthProvider>
          <Header />
        </AuthProvider>
      </CookiesProvider>
    </MemoryRouter>
  );

  await waitFor(() => {
    userEvent.click(screen.getByRole("button", { name: "User Menu" }));
  });
  expect(screen.getByRole("menuitem", { name: "Logout" })).toBeInTheDocument();
});
