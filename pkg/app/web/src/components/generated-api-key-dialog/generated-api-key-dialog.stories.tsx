import React from "react";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { GeneratedAPIKeyDialog } from "./generated-api-key-dialog";

export default {
  title: "SETTINGS/APIKey/GeneratedAPIKeyDialog",
  component: GeneratedAPIKeyDialog,
  decorators: [
    createDecoratorRedux({
      apiKeys: {
        disabling: false,
        error: null,
        generatedKey:
          "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.bspmf2xvt74area19iaxl0yh33jzwelq493vzil0orgzylrdb1",
        generating: false,
        loading: false,
        entities: {},
        ids: [],
      },
    }),
  ],
};

export const overview: React.FC = () => <GeneratedAPIKeyDialog />;
