import { fireEvent } from "@testing-library/react";
import React from "react";
import { MemoryRouter } from "react-router-dom";
import { render, screen } from "../../test-utils";
import { Role } from "../modules/me";
import { Header } from "./header";

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

  fireEvent.click(screen.getByRole("button", { name: "User Menu" }));
  expect(screen.getByRole("menuitem", { name: "Logout" })).toBeInTheDocument();
});
