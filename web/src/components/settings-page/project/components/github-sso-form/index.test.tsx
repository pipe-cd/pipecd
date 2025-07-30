import { screen, render } from "~~/test-utils";
import { GithubSSOForm } from "./index";
import { server } from "~/mocks/server";

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("GithubSSOForm", () => {
  it("renders without crashing", () => {
    render(<GithubSSOForm />);
    // Check for a heading, label, or unique element in the form
    // Adjust the text below to match your component's actual content
    expect(
      screen.getByRole("heading", { name: /Single Sign-On/i })
    ).toBeInTheDocument();
    expect(
      screen.getByText(
        "Single sign-on (SSO) allows users to log in to PipeCD by relying on a trusted third party service. Currently, only GitHub is supported."
      )
    ).toBeInTheDocument();
  });
});
