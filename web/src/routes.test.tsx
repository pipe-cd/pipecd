import { setupServer } from "msw/node";
import { waitFor } from "@testing-library/react";
import { MemoryRouter, render, screen } from "~~/test-utils";
import { Routes } from "./routes";
import { CookiesProvider } from "react-cookie";
import { AuthProvider } from "./contexts/auth-context";
import { getMeUnauthenticatedHandler } from "./mocks/services/me";

const server = setupServer(getMeUnauthenticatedHandler);

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

describe("Pages", () => {
  test("redirect to login page", async () => {
    const consoleErrorSpy = jest
      .spyOn(console, "error")
      .mockImplementation(() => {});

    render(
      <CookiesProvider>
        <MemoryRouter initialEntries={["/"]} initialIndex={0}>
          <AuthProvider>
            <Routes />
          </AuthProvider>
        </MemoryRouter>
      </CookiesProvider>
    );
    await waitFor(() =>
      expect(
        screen.getByRole("textbox", { name: /project name/i })
      ).toBeInTheDocument()
    );
    consoleErrorSpy.mockRestore();
  });
});
