import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { render, screen } from "~~/test-utils";
import { Role } from "~/modules/me";
import { Header } from "./";

it("shows login link if user state is not exists", () => {
  render(
    <MemoryRouter>
      <Header />
    </MemoryRouter>,
    {}
  );

  expect(screen.getByRole("link", { name: "Login" })).toBeInTheDocument();
});

it("shows logout link if opened user menu", () => {
  render(
    <MemoryRouter>
      <Header />
    </MemoryRouter>,
    {
      initialState: {
        me: {
          avatarUrl: "",
          subject: "user",
          isLogin: true,
          projectId: "pipecd",
          projectRole: Role.ProjectRole.ADMIN,
        },
      },
    }
  );

  userEvent.click(screen.getByRole("button", { name: "User Menu" }));
  expect(screen.getByRole("menuitem", { name: "Logout" })).toBeInTheDocument();
});
