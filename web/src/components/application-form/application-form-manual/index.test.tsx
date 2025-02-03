import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import ApplicationFormManual from ".";
import { render, screen } from "~~/test-utils";

const onClose = jest.fn();
const onFinished = jest.fn();
const TITLE = "title";

describe("ApplicationFormManual", () => {
  it("renders without crashing", () => {
    render(
      <ApplicationFormManual
        onClose={onClose}
        onFinished={onFinished}
        title="title"
      />,
      {}
    );
  });

  describe("Test ui", () => {
    beforeEach(() => {
      render(
        <ApplicationFormManual
          onClose={onClose}
          onFinished={onFinished}
          title={TITLE}
        />,
        {}
      );
    });
    it("should have correct title", () => {
      expect(screen.getByText(TITLE)).toBeInTheDocument();
    });

    it("calls onClose when cancel button is clicked", () => {
      const button = screen.getByRole("button", { name: UI_TEXT_CANCEL });
      button.click();
      expect(onClose).toHaveBeenCalledTimes(1);
    });

    it("button save is disabled initially", () => {
      const button = screen.getByRole("button", { name: UI_TEXT_SAVE });
      expect(button).toBeDisabled();
    });

    it('form contain input label "Name"', () => {
      const input = screen.getByLabelText(/^Name/i);
      expect(input).toBeInTheDocument();
    });

    it('form contain input label "Kind"', () => {
      const input = screen.getByLabelText(/^Kind/i);
      expect(input).toBeInTheDocument();
    });

    it('form contain input label "Piped"', () => {
      const input = screen.getByLabelText(/^Piped/i);
      expect(input).toBeInTheDocument();
    });

    it('form contain input label "Platform Provider"', () => {
      const input = screen.getByLabelText(/^Platform Provider/i);
      expect(input).toBeInTheDocument();
    });
    it('form contain input label "Repository"', () => {
      const input = screen.getByLabelText(/^Repository/i);
      expect(input).toBeInTheDocument();
    });
    it('form contain input label "Path"', () => {
      const input = screen.getByLabelText(/^Path/i);
      expect(input).toBeInTheDocument();
    });
    it('form contain input label "Config Filename"', () => {
      const input = screen.getByLabelText(/^Config Filename/i);
      expect(input).toBeInTheDocument();
    });
  });
});
