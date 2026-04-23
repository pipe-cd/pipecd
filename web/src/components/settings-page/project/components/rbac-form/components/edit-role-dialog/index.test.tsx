import { server } from "~/mocks/server";
import { screen, render, waitFor } from "~~/test-utils";
import { EditRoleDialog } from ".";
import { dummyRole } from "~/__fixtures__/dummy-project";
import userEvent from "@testing-library/user-event";

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("EditRoleDialog", () => {
  it("should render without crashing", () => {
    render(
      <EditRoleDialog
        role={dummyRole}
        onClose={() => {}}
        onUpdate={(values) => console.log(values)}
      />
    );

    expect(screen.getByText("Edit Role")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Cancel" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Edit" })).toBeInTheDocument();
  });

  it("should trigger onUpdate with correct values on form submission", async () => {
    const mockOnUpdate = jest.fn();
    render(
      <EditRoleDialog
        role={dummyRole}
        onClose={() => {}}
        onUpdate={mockOnUpdate}
      />
    );

    const policiesInput = screen.getByRole("textbox", { name: "Policies" });

    await userEvent.clear(policiesInput);
    await userEvent.type(policiesInput, "resources=application,deployment;actions=*");

    await waitFor(() =>
      expect(policiesInput).toHaveValue(
        "resources=application,deployment;actions=*"
      )
    );

    await userEvent.click(screen.getByRole("button", { name: /edit/i }));

    await waitFor(() => {
      expect(mockOnUpdate).toHaveBeenCalledWith({
        name: dummyRole.name,
        policies: "resources=application,deployment;actions=*",
      });
    });
  });
});
