import userEvent from "@testing-library/user-event";
import React from "react";
import { render, screen, act } from "../../test-utils";
import { SplitButton } from "./split-button";

it("calls onClick handler with option's index if clicked", () => {
  const onClick = jest.fn();
  render(
    <SplitButton
      loading={false}
      onClick={onClick}
      options={["option1", "option2"]}
    />,
    {}
  );

  userEvent.click(screen.getByRole("button", { name: "option1" }));

  expect(onClick).toHaveBeenCalledWith(0);

  act(() => {
    userEvent.click(
      screen.getByRole("button", { name: "select merge strategy" })
    );
  });
  userEvent.click(screen.getByRole("menuitem", { name: "option2" }));
  userEvent.click(screen.getByRole("button", { name: "option2" }));

  expect(onClick).toHaveBeenCalledWith(1);
});
