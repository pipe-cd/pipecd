import React, { FC, memo } from "react";
import { Route, Switch } from "react-router-dom";
import { Header } from "../components/header";
import { PAGE_PATH_DEPLOYMENTS, PAGE_PATH_APPLICATIONS } from "../constants";
import { DeploymentDetailPage } from "./deployments/detail";

export const Pages: FC = memo(() => {
  return (
    <div>
      <Header />
      <Switch>
        <Route
          exact
          path={PAGE_PATH_APPLICATIONS}
          component={() => <div>applications</div>}
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
