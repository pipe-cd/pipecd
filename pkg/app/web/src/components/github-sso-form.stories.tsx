import React from "react";
import { GithubSSOForm } from "./github-sso-form";
import { action } from "@storybook/addon-actions";

export default {
  title: "GithubSSOForm",
  component: GithubSSOForm,
};

export const overview: React.FC = () => (
  <GithubSSOForm
    onSave={(params) => {
      action("onSave")(params);
      return Promise.resolve();
    }}
    isSaving={false}
  />
);
