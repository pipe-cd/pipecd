import { act } from "react-dom/test-utils";
import { render, screen, waitFor } from "~~/test-utils";
import { APIKeyPage } from ".";
import { setupServer } from "msw/node";
import { apiKeyHandlers, getListAPIKeysEmpty } from "~/mocks/services/api-keys";

// Test suite with API mock
const server = setupServer();

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
  jest.resetModules();
});

afterAll(() => {
  server.close();
});
describe("APIKeyPage render correct row and header", () => {
  beforeEach(() => {
    server.use(...apiKeyHandlers);
  });
  it("Render table with 1 api keys", async () => {
    await act(async () => {
      await render(<APIKeyPage />);
    });

    await waitFor(() => {
      expect(screen.getByText("API_KEY_1")).toBeInTheDocument();
    });
    expect(
      screen.getByRole("columnheader", { name: "Name" })
    ).toBeInTheDocument();
    expect(
      screen.getByRole("columnheader", { name: "Role" })
    ).toBeInTheDocument();
    expect(
      screen.getByRole("columnheader", { name: "CreatedAt" })
    ).toBeInTheDocument();
    expect(
      screen.getByRole("columnheader", { name: "LastUsedAt" })
    ).toBeInTheDocument();
  });
});

// Test suite without API mock
describe("APIKeyPage without API mock", () => {
  beforeEach(() => {
    server.use(getListAPIKeysEmpty);
  });

  it("Render empty table when no API keys", async () => {
    await act(async () => {
      await render(<APIKeyPage />);
    });

    await waitFor(() => {
      expect(screen.getByText("No API Keys")).toBeInTheDocument();
    });
  });
});
