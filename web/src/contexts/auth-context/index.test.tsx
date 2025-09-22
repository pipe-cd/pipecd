import { setupServer } from "msw/node";
import { MemoryRouter, render, screen, waitFor } from "~~/test-utils";
import useAuth from "./use-auth";
import { FC } from "react";
import { getMeHandler } from "~/mocks/services/me";
import useProjectName from "./use-project-name";
import { AuthProvider } from "./auth-provider";
import { Cookies, CookiesProvider } from "react-cookie";

const server = setupServer(getMeHandler);

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

const TestComponent: FC = () => {
  const { me } = useAuth();
  const projectName = useProjectName();

  return (
    <div>
      <p data-testid="avatarUrl">{me ? me.avatarUrl : "No Avatar"}</p>
      <p data-testid="projectId">{me ? me.projectId : "No ProjectId"}</p>
      <p data-testid="projectName">{projectName}</p>
      <p data-testid="subject">{me ? me.subject : "No Subject"}</p>
      <p data-testid="isLogin">{me?.isLogin ? "Logged In" : "Not Logged In"}</p>
    </div>
  );
};

describe("auth context", () => {
  beforeEach(() => {
    const cookies = new Cookies();
    cookies.set("token", "my-test-token");
    render(
      <MemoryRouter>
        <CookiesProvider cookies={cookies}>
          <AuthProvider>
            <TestComponent />
          </AuthProvider>
        </CookiesProvider>
      </MemoryRouter>
    );
  });

  it("should fetch user data on mount", async () => {
    await waitFor(() => {
      expect(screen.getByTestId("avatarUrl")).toHaveTextContent("avatar-url");
    });
    await waitFor(() => {
      expect(screen.getByTestId("subject")).toHaveTextContent("userName");
    });
    await waitFor(() =>
      expect(screen.getByTestId("isLogin")).toHaveTextContent("Logged In")
    );
    await waitFor(() =>
      expect(screen.getByTestId("projectId")).toHaveTextContent("pipecd")
    );
    await waitFor(() =>
      expect(screen.getByTestId("projectName")).toHaveTextContent("pipecd")
    );
  });
});
