import { server } from "../mocks/server";
import { getMe } from "./me";

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

test("getMe() call", async () => {
  await expect(getMe()).resolves.toMatchInlineSnapshot(`
    Object {
      "avatarUrl": "https://test.pipecd.dev/avatar.jpg",
      "projectId": "pipecd",
      "projectRole": 2,
      "subject": "hello-pipecd",
    }
  `);
});
