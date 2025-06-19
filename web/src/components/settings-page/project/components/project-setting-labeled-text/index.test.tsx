import { render, screen } from "@testing-library/react";
import { ProjectSettingLabeledText } from "./index";

describe("ProjectSettingLabeledText", () => {
  it("renders the label and value", () => {
    render(<ProjectSettingLabeledText label="Client ID" value="123456" />);
    expect(screen.getByText("Client ID")).toBeInTheDocument();
    expect(screen.getByText("123456")).toBeInTheDocument();
  });
});
