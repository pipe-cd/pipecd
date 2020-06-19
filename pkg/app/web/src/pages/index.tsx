import React, { FC, memo } from "react";
import { Route, Switch } from "react-router-dom";
import { Header } from "../components/header";
import {
  PAGE_PATH_DEPLOYMENTS,
  PAGE_PATH_APPLICATIONS,
  PAGE_PATH_TOP,
  PAGE_PATH_SETTINGS,
} from "../constants";
import { ApplicationIndexPage } from "./applications/index";
import { ApplicationDetailPage } from "./applications/detail";
import { DeploymentIndexPage } from "./deployments/index";
import { DeploymentDetailPage } from "./deployments/detail";
import { SettingsIndexPage } from "./settings";

export const Pages: FC = memo(() => {
  return (
    <>
      <Header />
      <Switch>
        <Route exact path={PAGE_PATH_TOP} component={ApplicationIndexPage} />
        <Route
          exact
          path={PAGE_PATH_APPLICATIONS}
          component={ApplicationIndexPage}
        />
        <Route
          exact
          path={`${PAGE_PATH_APPLICATIONS}/:applicationId`}
          component={ApplicationDetailPage}
        />
        <Route
          exact
          path={PAGE_PATH_DEPLOYMENTS}
          component={DeploymentIndexPage}
        />
        <Route
          exact
          path={`${PAGE_PATH_DEPLOYMENTS}/:deploymentId`}
          component={DeploymentDetailPage}
        />
        <Route path={PAGE_PATH_SETTINGS} component={SettingsIndexPage} />
      </Switch>
    </>
  );
});
