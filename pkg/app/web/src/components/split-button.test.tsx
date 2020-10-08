import { fireEvent } from "@testing-library/react";
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

  fireEvent.click(screen.getByRole("button", { name: "option1" }));

  expect(onClick).toHaveBeenCalledWith(0);

  act(() => {
    fireEvent.click(
      screen.getByRole("button", { name: "select merge strategy" })
    );
  });
  fireEvent.click(screen.getByRole("menuitem", { name: "option2" }));
  fireEvent.click(screen.getByRole("button", { name: "option2" }));

  expect(onClick).toHaveBeenCalledWith(1);
});
