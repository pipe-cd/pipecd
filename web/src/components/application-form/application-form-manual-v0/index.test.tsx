import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import ApplicationFormManualV0 from ".";
import { createStore, render, screen } from "~~/test-utils";
import { server } from "~/mocks/server";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { AppState } from "~/store";

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

const baseState: Partial<AppState> = {
  pipeds: {
    entities: {
      [dummyPiped.id]: dummyPiped,
    },
    ids: [dummyPiped.id],
    registeredPiped: null,
    updating: false,
    releasedVersions: [],
    breakingChangesNote: "",
  },
};

describe("ApplicationFormManualV0", () => {
  it("renders without crashing", () => {
    render(
      <ApplicationFormManualV0
        onClose={onClose}
        onFinished={onFinished}
        title="title"
      />,
      {}
    );
  });

  describe("Test ui create application", () => {
    const store = createStore(baseState);
    beforeEach(() => {
      render(
        <ApplicationFormManualV0
          onClose={onClose}
          onFinished={onFinished}
          title={TITLE}
        />,
        { store }
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

    it('form contain input label "Name" and not disabled initially', () => {
      const input = screen.getByLabelText(/^Name/i);
      expect(input).toBeInTheDocument();
      expect(input).not.toBeDisabled();
    });

    it('form contain input label "Kind" and not disabled initially', () => {
      const input = screen.getByLabelText(/^Kind/i);
      expect(input).toBeInTheDocument();
      expect(input).not.toBeDisabled();
    });

    it('form contain input label "Piped" and not disabled initially', () => {
      const input = screen.getByRole("button", { name: "Piped" });
      expect(input).toBeInTheDocument();
      expect(input).not.toHaveAttribute("aria-disabled", "true");
    });

    it('form contain input label "Platform Provider" and disabled', () => {
      const input = screen.getByRole("button", { name: "Platform Provider" });
      expect(input).toBeInTheDocument();
      expect(input).toHaveAttribute("aria-disabled", "true");
    });

    it('form contain input label "Repository" and disabled', () => {
      const input = screen.getByRole("button", { name: "Repository" });
      expect(input).toBeInTheDocument();
      expect(input).toHaveAttribute("aria-disabled", "true");
    });

    it('form contain input label "Path" and disabled', () => {
      const input = screen.getByLabelText(/^Path/i);
      expect(input).toBeInTheDocument();
      expect(input).toBeDisabled();
    });

    it('form contain input label "Config Filename" disabled and init value = app.pipecd.yaml', () => {
      const input = screen.getByLabelText(/^Config Filename/i);
      expect(input).toBeInTheDocument();
      expect(input).toHaveValue("app.pipecd.yaml");
      expect(input).toBeDisabled();
    });
  });

  describe("Test ui edit application", () => {
    beforeEach(() => {
      const store = createStore(baseState);
      render(
        <ApplicationFormManualV0
          onClose={onClose}
          onFinished={onFinished}
          title={TITLE}
          detailApp={dummyApplication}
        />,
        { store }
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

    it('form contain input label "Kind" and disabled initially', () => {
      const input = screen.getByRole("button", { name: "Kind" });
      expect(input).toBeInTheDocument();
      expect(input).toHaveAttribute("aria-disabled", "true");
    });

    it('form contain input label "Piped"', () => {
      const input = screen.getByRole("button", { name: "Piped" });
      expect(input).toBeInTheDocument();
      expect(input).not.toHaveAttribute("aria-disabled", "true");
    });

    it('form contain input label "Platform Provider"', () => {
      const input = screen.getByRole("button", { name: "Platform Provider" });
      expect(input).toBeInTheDocument();
      expect(input).not.toHaveAttribute("aria-disabled", "true");
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
