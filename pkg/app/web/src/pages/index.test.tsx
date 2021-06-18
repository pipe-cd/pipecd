import { MemoryRouter } from "react-router-dom";
import { render, screen } from "test-utils";
import { Pages } from "./index";
import { setupServer } from "msw/node";
import { GetMeResponse } from "pipe/pkg/app/web/api_client/service_pb";
import { createHandler } from "../mocks/create-handler";
import { waitFor } from "@testing-library/react";

const server = setupServer(
  createHandler<GetMeResponse>("/GetMe", () => {
    const response = new GetMeResponse();
    return response;
  })
);

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

describe("Pages", () => {
  test("redirect to login page", async () => {
    render(
      <MemoryRouter initialEntries={["/"]} initialIndex={0}>
        <Pages />
      </MemoryRouter>,
      { initialState: { me: { isLogin: false } } }
    );
    await waitFor(() =>
      expect(
        screen.getByRole("textbox", { name: /project name/i })
      ).toBeInTheDocument()
    );
  });
});
