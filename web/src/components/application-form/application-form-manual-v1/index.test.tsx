import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import ApplicationFormManualV1 from ".";
import { render, screen, waitFor } from "~~/test-utils";
import { server } from "~/mocks/server";
import { dummyApplicationPipedV1 } from "~/__fixtures__/dummy-application";

const onClose = jest.fn();
const onFinished = jest.fn();
const TITLE = "title";

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

describe("ApplicationFormManualV1", () => {
  it("renders without crashing", () => {
    render(
      <ApplicationFormManualV1
        onClose={onClose}
        onFinished={onFinished}
        title="title"
      />,
      {}
    );
  });

  describe("Test ui create application", () => {
    beforeEach(() => {
      render(
        <ApplicationFormManualV1
          onClose={onClose}
          onFinished={onFinished}
          title={TITLE}
        />
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
      expect(input).not.toBeDisabled();
    });

    it('form contain input label "Piped"', async () => {
      const input = screen.getByRole("combobox", { name: "Piped" });
      expect(input).toBeInTheDocument();
      await waitFor(() => {
        expect(input).not.toHaveAttribute("aria-disabled", "true");
      });
    });

    it('form contain input label "Repository"', () => {
      const input = screen.getByRole("combobox", { name: "Repository" });
      expect(input).toBeInTheDocument();
      expect(input).toHaveAttribute("aria-disabled", "true");
    });

    it('form contain input label "Path"', () => {
      const input = screen.getByLabelText(/^Path/i);
      expect(input).toBeInTheDocument();
      expect(input).toBeDisabled();
    });

    it('form contain input label "Config Filename"', () => {
      const input = screen.getByLabelText(/^Config Filename/i);
      expect(input).toBeInTheDocument();
      expect(input).toHaveValue("app.pipecd.yaml");
      expect(input).toBeDisabled();
    });
  });

  describe("Test ui edit application", () => {
    beforeEach(() => {
      render(
        <ApplicationFormManualV1
          onClose={onClose}
          onFinished={onFinished}
          title={TITLE}
          detailApp={dummyApplicationPipedV1}
        />
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

    it('form contain input label "Name" and disabled initially', () => {
      const input = screen.getByLabelText(/^Name/i);
      expect(input).toBeInTheDocument();
      expect(input).toBeDisabled();
    });

    it('form contain input label "Piped"', async () => {
      const input = screen.getByRole("combobox", { name: "Piped" });
      expect(input).toBeInTheDocument();
      await waitFor(() => {
        expect(input).not.toHaveAttribute("aria-disabled", "true");
      });
    });

    it('form contain input label "Repository" and disabled initially', () => {
      const input = screen.getByLabelText(/^Repository/i);
      expect(input).toBeInTheDocument();
      expect(input).toHaveAttribute("aria-disabled", "true");
    });

    it('form contain input label "Path" and disabled initially', () => {
      const input = screen.getByLabelText(/^Path/i);
      expect(input).toBeInTheDocument();
      expect(input).toBeDisabled();
    });

    it('form contain input label "Config Filename"', () => {
      const input = screen.getByLabelText(/^Config Filename/i);
      expect(input).toBeInTheDocument();
    });
  });
});
