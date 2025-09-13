import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import ApplicationFormV1 from ".";
import { render, screen, waitFor } from "~~/test-utils";
import userEvent from "@testing-library/user-event";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { server } from "~/mocks/server";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { act } from "react-dom/test-utils";
import { ApplicationInfo } from "~/types/applications";

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
  configFilename: dummyApplication.gitPath?.configFilename || "",
  path: dummyApplication.gitPath?.path,
  repoId: dummyApplication.gitPath?.repo?.id,
  labelsMap: [] as [string, string][],
} as ApplicationInfo.AsObject;

describe("ApplicationFormV1", () => {
  it("renders without crashing", () => {
    render(
      <ApplicationFormV1
        onClose={onClose}
        onFinished={onFinished}
        title="title"
      />
    );
    expect(screen.getByText("title")).toBeInTheDocument();
  });

  describe("Test ui", () => {
    beforeEach(() => {
      render(
        <ApplicationFormV1
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

    it("Form should have 3 step", () => {
      const step1 = screen.getByText("Select piped and deploy targets");
      const step2 = screen.getByText("Select application to add");
      const step3 = screen.getByText("Confirm information before adding");
      expect(step1).toBeInTheDocument();
      expect(step2).toBeInTheDocument();
      expect(step3).toBeInTheDocument();
    });

    it("form contain correct input by Step", async () => {
      // select piped
      userEvent.click(screen.getByRole("combobox", { name: /piped/i }));
      const optionName = `${dummyPiped.name} (${dummyPiped.id})`;
      await waitFor(() => {
        userEvent.click(screen.getByRole("option", { name: optionName }));
      });

      // select deploy target
      userEvent.click(screen.getByRole("combobox", { name: /deploy target/i }));
      const optionDeployTarget = `${dummyPiped.pluginsList[0].deployTargetsList[0]} - ${dummyPiped.pluginsList[0].name}`;
      userEvent.click(screen.getByRole("option", { name: optionDeployTarget }));

      // select Application
      userEvent.click(screen.getByRole("combobox", { name: "Application" }));
      const optionApplicationName = `name: ${dummyUnregisteredApplication.name}, repo: ${dummyUnregisteredApplication.repoId}`;
      userEvent.click(
        screen.getByRole("option", { name: optionApplicationName })
      );

      // check
      expect(screen.getByRole("textbox", { name: "Path" })).toHaveValue(
        dummyUnregisteredApplication.path
      );
      expect(
        screen.getByRole("textbox", { name: "Config Filename" })
      ).toHaveValue(dummyUnregisteredApplication.configFilename);

      // click save
      expect(screen.getByRole("button", { name: UI_TEXT_SAVE })).toBeEnabled();
      act(() => {
        userEvent.click(screen.getByRole("button", { name: UI_TEXT_SAVE }));
      });
      expect(screen.getByText("Add Application")).toBeInTheDocument();
    });
  });
});
