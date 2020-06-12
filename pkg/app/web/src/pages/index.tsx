import React, { FC, memo } from "react";
import { Route, Switch } from "react-router-dom";
import { Header } from "../components/header";
import { PAGE_PATH_DEPLOYMENTS, PAGE_PATH_APPLICATIONS } from "../constants";
import { DeploymentDetailPage } from "./deployments/detail";
import { ApplicationDetailPage } from "./applications/detail";
import { DeploymentIndexPage } from "./deployments";

export const Pages: FC = memo(() => {
  return (
    <div>
      <Header />
      <Switch>
        <Route
          exact
          path={`${PAGE_PATH_APPLICATIONS}/:applicationId`}
          component={ApplicationDetailPage}
        />
        <Route
          exact
          path={`${PAGE_PATH_DEPLOYMENTS}`}
          component={DeploymentIndexPage}
        />
        <Route
          exact
          path={`${PAGE_PATH_DEPLOYMENTS}/:deploymentId`}
          component={DeploymentDetailPage}
        />
      </Switch>
    </div>
  );
});
