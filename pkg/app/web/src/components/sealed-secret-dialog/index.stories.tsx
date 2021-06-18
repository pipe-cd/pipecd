import { action } from "@storybook/addon-actions";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { SealedSecretDialog } from "./";
import { Provider } from "react-redux";
import { createStore } from "~~/test-utils";
import { Story } from "@storybook/react";

export default {
  title: "APPLICATION/SealedSecretDialog",
  component: SealedSecretDialog,
};

export const Overview: Story = () => (
  <Provider
    store={createStore({
      applications: {
        entities: { [dummyApplication.id]: dummyApplication },
        ids: [dummyApplication.id],
      },
    })}
  >
    <SealedSecretDialog
      open
      applicationId={dummyApplication.id}
      onClose={action("onClose")}
    />
  </Provider>
);

export const Encrypted: Story = () => (
  <Provider
    store={createStore({
      applications: {
        entities: { [dummyApplication.id]: dummyApplication },
        ids: [dummyApplication.id],
      },
      sealedSecret: {
        isLoading: false,
        // dummy data
        data:
          "AgEAk+XZY0+8hJSv5GG8rDTZGV56xe3xCROxGtVwNVY3kSg3MAvq1BcbhkIToT1q7JVYnE/MQ/nuks5MXPrFpu5/bHCehlzsFc8tnffQeYtVO3XuCi771zd4KNsAmCMvN0Fo57eqzuhU/uEYt+1thRJh3FA/tJilZ0j1JiSvrn01Zfb1Xw0QGMCO4/C75YRUa6g2JC75yLKZ06Yzequ1wSglLzF7rzUdl1+kxCBWJM1HnFkBpLGtCcH7xqBbkQMEz1jOjgpluPDUD7nCwMmZzY9sw83jOXfsCdNvmy/uStIsTDBZDPE1uWxN8MAWEieLoOvLLUFWGXmMDEJAa6nv1iI/jv8OLXzUEOeklc/OKZ7jLOnylmqQ4U7YCFsfHqAMzAFadpORxHCHt09zHzm9nKoyhOCFo9gDVyX6HiD9h3V/j1gBmr9lgonfiVC0HbsghOLJMjf+9njNiEqD7IOhBVv6TEwGpGI+1SYoBu6m8+Jlex2j5VcspMybx6p5aKmuwQpLqrEFMZBYJwAXMrAxkMiiE5QnB9/K7YAcYr3a33qGPu+JmQKgg9QRN0X5jvl/FNaYDoXzNp1FBuFEQUtu0hw8QKtf/DWVAUj/zgyviPQ2j8fy5Yg7qwxuL12EQMeAtC0pwyFNE7SlyvGoGoKc+yp0/EfpFzyH3n807+UD8ueN+bbvKn0gXkePO2Af1a61LU/lPOe68A==",
      },
    })}
  >
    <SealedSecretDialog
      open
      applicationId={dummyApplication.id}
      onClose={action("onClose")}
    />
  </Provider>
);
