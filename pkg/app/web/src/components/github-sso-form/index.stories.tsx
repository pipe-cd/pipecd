import * as React from "react";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { GithubSSOForm } from "./";

export default {
  title: "SETTINGS/GithubSSOForm",
  component: GithubSSOForm,
  decorators: [createDecoratorRedux({})],
};

export const overview: React.FC = () => <GithubSSOForm />;
