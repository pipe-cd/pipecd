import React from "react";
import { StaticAdminForm } from "./static-admin-form";
import { action } from "@storybook/addon-actions";

export default {
  title: "StaticAdminForm",
  component: StaticAdminForm,
};

export const overview: React.FC = () => (
  <StaticAdminForm
    staticAdminDisabled={false}
    username="User"
    onUpdatePassword={(password: string) => {
      action("onUpdatePassword")(password);
      return Promise.resolve();
    }}
    onUpdateUsername={action("onUpdateUsername")}
    onToggleAvailability={action("onToggleAvailability")}
    isUpdatingPassword={false}
    isUpdatingUsername={false}
  />
);
