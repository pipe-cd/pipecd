import { EntityId } from "@reduxjs/toolkit";
import React, { FC, memo, useEffect } from "react";
import loadable from "@loadable/component";
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
} from "../constants/path";
import { AppState } from "../modules";
import {
  fetchCommand,
  selectIds as selectCommandIds,
} from "../modules/commands";
import { fetchEnvironments } from "../modules/environments";
import { useMe } from "../modules/me";
import { fetchPipeds } from "../modules/pipeds";
import { useInterval } from "../hooks/use-interval";
import { ApplicationIndexPage } from "./applications/index";
import { DeploymentIndexPage } from "./deployments/index";
import { LoginPage } from "./login";

const SettingsIndexPage = loadable(
  () => import(/* webpackChunkName: "settings" */ "./settings"),
  {
    resolveComponent: (components) => components.SettingsIndexPage,
  }
);

const InsightIndexPage = loadable(
  () => import(/* webpackChunkName: "insight" */ "./insight"),
  {
    resolveComponent: (components) => components.InsightIndexPage,
  }
);

const DeploymentDetailPage = loadable(
  () =>
    import(/* webpackChunkName: "deployments-detail" */ "./deployments/detail"),
  {
    resolveComponent: (components) => components.DeploymentDetailPage,
  }
);

const ApplicationDetailPage = loadable(
  () =>
    import(
      /* webpackChunkName: "applications-detail" */ "./applications/detail"
    ),
  {
    resolveComponent: (components) => components.ApplicationDetailPage,
  }
);

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
    if (me?.isLogin) {
      dispatch(fetchEnvironments());
      dispatch(fetchPipeds(false));
    }
  }, [dispatch, me]);
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
        <Switch>
          <Route
            exact
            path={PAGE_PATH_TOP}
            component={() => <Redirect to={PAGE_PATH_LOGIN} />}
          />
          <Route path={`${PAGE_PATH_LOGIN}/:projectName?`}>
            <LoginPage />
          </Route>
        </Switch>
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
