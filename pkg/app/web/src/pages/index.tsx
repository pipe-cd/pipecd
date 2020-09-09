import { EntityId } from "@reduxjs/toolkit";
import React, { FC, memo, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Route, Switch, Redirect } from "react-router-dom";
import { Header } from "../components/header";
import { Toasts } from "../components/toasts";
import {
  PAGE_PATH_APPLICATIONS,
  PAGE_PATH_DEPLOYMENTS,
  PAGE_PATH_INSIGHTS,
  PAGE_PATH_LOGIN,
  PAGE_PATH_SETTINGS,
  PAGE_PATH_TOP,
} from "../constants";
import { AppState } from "../modules";
import {
  fetchCommand,
  selectIds as selectCommandIds,
} from "../modules/commands";
import { fetchEnvironments } from "../modules/environments";
import { useMe } from "../modules/me";
import { fetchPipeds } from "../modules/pipeds";
import { useInterval } from "../utils/use-interval";
import { ApplicationDetailPage } from "./applications/detail";
import { ApplicationIndexPage } from "./applications/index";
import { DeploymentDetailPage } from "./deployments/detail";
import { DeploymentIndexPage } from "./deployments/index";
import { InsightIndexPage } from "./insight";
import { LoginPage } from "./login";
import { SettingsIndexPage } from "./settings";

// Fetch commands detail periodically
const FETCH_COMMANDS_INTERVAL = 3000;
const useCommandsStatusChecking = (): void => {
  const dispatch = useDispatch();
  const commandIds = useSelector<AppState, EntityId[]>((state) =>
    selectCommandIds(state.commands)
  );

  const fetchCommands = (): void => {
    commandIds.map((id) => {
      dispatch(fetchCommand(`${id}`));
    });
  };

  useInterval(
    fetchCommands,
    commandIds.length > 0 ? FETCH_COMMANDS_INTERVAL : null
  );
};

export const Pages: FC = memo(function Pages() {
  const dispatch = useDispatch();
  const me = useMe();
  useEffect(() => {
    dispatch(fetchEnvironments());
    dispatch(fetchPipeds(false));
  }, [dispatch]);
  useCommandsStatusChecking();

  if (me === null) {
    return (
      <>
        <Header />
      </>
    );
  }

  if (me.isLogin === false) {
    return (
      <>
        <Header />
        <LoginPage />
      </>
    );
  }

  return (
    <>
      <Header />
      <Switch>
        <Route
          exact
          path={PAGE_PATH_TOP}
          component={() => <Redirect to={PAGE_PATH_APPLICATIONS} />}
        />
        <Route exact path={PAGE_PATH_LOGIN} component={LoginPage} />
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
        <Route path={PAGE_PATH_INSIGHTS} component={InsightIndexPage} />
      </Switch>
      <Toasts />
    </>
  );
});
