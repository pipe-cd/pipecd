import React, { FC, memo, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Route, Switch } from "react-router-dom";
import { Header } from "../components/header";
import { Toasts } from "../components/toasts";
import {
  PAGE_PATH_APPLICATIONS,
  PAGE_PATH_DEPLOYMENTS,
  PAGE_PATH_SETTINGS,
  PAGE_PATH_TOP,
  PAGE_PATH_LOGIN,
} from "../constants";
import { AppState } from "../modules";
import {
  fetchCommand,
  selectIds as selectCommandIds,
} from "../modules/commands";
import { fetchEnvironments } from "../modules/environments";
import { fetchPipeds } from "../modules/pipeds";
import { useInterval } from "../utils/use-interval";
import { ApplicationDetailPage } from "./applications/detail";
import { ApplicationIndexPage } from "./applications/index";
import { DeploymentDetailPage } from "./deployments/detail";
import { DeploymentIndexPage } from "./deployments/index";
import { SettingsIndexPage } from "./settings";
import { EntityId } from "@reduxjs/toolkit";
import { LoginPage } from "./login";

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
  useEffect(() => {
    dispatch(fetchEnvironments());
    dispatch(fetchPipeds(false));
  }, [dispatch]);
  useCommandsStatusChecking();

  return (
    <>
      <Header />
      <Switch>
        <Route exact path={PAGE_PATH_TOP} component={ApplicationIndexPage} />
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
      </Switch>
      <Toasts />
    </>
  );
});
