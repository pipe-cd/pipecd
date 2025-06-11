import { screen, render, waitFor } from "~~/test-utils";
import { RBACForm } from ".";
import { server } from "~/mocks/server";

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("RBACForm", () => {
  it("renders the title and description correctly", async () => {
    render(<RBACForm />);
    await waitFor(() => {
      expect(
        screen.getByRole("heading", { name: /Role-Based Access Control/i })
      ).toBeInTheDocument();
    });

    expect(
      screen.getByText(
        "Role-based access control (RBAC) allows restricting the access on PipeCD web based on the roles of user groups within the project. Before using this feature, the SSO must be configured."
      )
    ).toBeInTheDocument();
  });
});
