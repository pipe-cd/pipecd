import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import ApplicationFormV1 from ".";
import { createStore, render, screen } from "~~/test-utils";
import userEvent from "@testing-library/user-event";
import { AppState } from "~/store";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { server } from "~/mocks/server";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";
import { act } from "react-dom/test-utils";

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

const dummyUnregisteredApplication = {
  ...dummyApplication,
  configFilename: "app1.yaml",
  path: "app1",
  repoId: "repo1",
  labelsMap: [["env", "test"]] as [string, string][],
};

const baseState: Partial<AppState> = {
  unregisteredApplications: {
    ids: [dummyUnregisteredApplication.id],
    apps: [dummyUnregisteredApplication],
    entities: {
      [dummyUnregisteredApplication.id]: dummyUnregisteredApplication,
    },
  },
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

describe("ApplicationFormV1", () => {
  it("renders without crashing", () => {
    render(
      <ApplicationFormV1
        onClose={onClose}
        onFinished={onFinished}
        title="title"
      />,
      {}
    );
  });

  describe("Test ui", () => {
    const store = createStore(baseState);
    beforeEach(() => {
      render(
        <ApplicationFormV1
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

    it("button save is enabled initially when only 1 piped and 1 app in options", () => {
      const button = screen.getByRole("button", { name: UI_TEXT_SAVE });
      expect(button).toBeEnabled();
    });

    it("Form should have 3 step", () => {
      const step1 = screen.getByText("Select piped");
      const step2 = screen.getByText("Select application to add");
      const step3 = screen.getByText("Confirm information before adding");
      expect(step1).toBeInTheDocument();
      expect(step2).toBeInTheDocument();
      expect(step3).toBeInTheDocument();
    });

    it("form contain correct input by Step", async () => {
      // select piped
      userEvent.click(screen.getByRole("button", { name: /piped/i }));
      const optionName = `${dummyPiped.name} (${dummyPiped.id})`;
      userEvent.click(screen.getByRole("option", { name: optionName }));

      // select Application
      userEvent.click(screen.getByRole("textbox", { name: "Application" }));
      const optionApplicationName = `name: ${dummyUnregisteredApplication.name}, repo: ${dummyUnregisteredApplication.repoId}`;
      userEvent.click(
        screen.getByRole("option", { name: optionApplicationName })
      );

      // check
      expect(screen.getByRole("textbox", { name: "Kind" })).toHaveValue(
        APPLICATION_KIND_TEXT[dummyUnregisteredApplication.kind]
      );
      expect(screen.getByRole("textbox", { name: "Path" })).toHaveValue(
        dummyUnregisteredApplication.path
      );
      expect(
        screen.getByRole("textbox", { name: "Config Filename" })
      ).toHaveValue(dummyUnregisteredApplication.configFilename);
      expect(screen.getByRole("textbox", { name: "Label 0" })).toHaveValue(
        dummyUnregisteredApplication.labelsMap[0].join(": ")
      );

      // click save
      expect(screen.getByRole("button", { name: UI_TEXT_SAVE })).toBeEnabled();
      act(() => {
        userEvent.click(screen.getByRole("button", { name: UI_TEXT_SAVE }));
      });
      expect(screen.getByText("Add Application")).toBeInTheDocument();
    });
  });
});
