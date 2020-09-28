import React from "react";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";
import { GithubSSOForm } from "./github-sso-form";

export default {
  title: "SETTINGS/GithubSSOForm",
  component: GithubSSOForm,
  decorators: [createDecoratorRedux({})],
};

export const overview: React.FC = () => <GithubSSOForm />;
