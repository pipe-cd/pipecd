import loadable from "@loadable/component";
import { EntityId } from "@reduxjs/toolkit";
import { FC, useEffect, useState } from "react";
import {
  Redirect,
  Route,
  Switch,
  useLocation,
  RouteComponentProps,
} from "react-router-dom";
import { ApplicationIndexPage } from "~/components/applications-page";
import { WarningBanner } from "~/components/warning-banner";
import { DeploymentIndexPage } from "~/components/deployments-page";
import { DeploymentChainsIndexPage } from "~/components/deployment-chains-page";
import { Header } from "~/components/header";
import { LoginPage } from "~/components/login-page";
import { Toasts } from "~/components/toasts";
import {
  PAGE_PATH_APPLICATIONS,
  PAGE_PATH_DEPLOYMENTS,
  PAGE_PATH_DEPLOYMENT_CHAINS,
  PAGE_PATH_INSIGHTS,
  PAGE_PATH_EVENTS,
  PAGE_PATH_LOGIN,
  PAGE_PATH_SETTINGS,
  PAGE_PATH_TOP,
} from "~/constants/path";
import {
  REDIRECT_PATH_KEY,
  BANNER_VERSION_KEY,
  USER_PROJECTS,
} from "~/constants/localstorage";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { useInterval } from "~/hooks/use-interval";
import useQueryString from "./hooks/use-query-string";
import {
  fetchCommand,
  selectIds as selectCommandIds,
} from "~/modules/commands";
import { fetchPipeds } from "~/modules/pipeds";
import { sortedSet } from "~/utils/sorted-set";

const SettingsIndexPage = loadable(
  () => import(/* webpackChunkName: "settings" */ "~/components/settings-page"),
  {
    resolveComponent: (components) => components.SettingsIndexPage,
  }
);

const InsightIndexPage = loadable(
  () => import(/* webpackChunkName: "insight" */ "~/components/insight-page"),
  {
    resolveComponent: (components) => components.InsightIndexPage,
  }
);

const EventIndexPage = loadable(
  () => import(/* webpackChunkName: "events" */ "~/components/events-page"),
  {
    resolveComponent: (components) => components.EventIndexPage,
  }
);

const DeploymentDetailPage = loadable(
  () =>
    import(
      /* webpackChunkName: "deployments-detail" */ "~/components/deployments-detail-page"
    ),
  {
    resolveComponent: (components) => components.DeploymentDetailPage,
  }
);

const ApplicationDetailPage = loadable(
  () =>
    import(
      /* webpackChunkName: "applications-detail" */ "~/components/application-detail-page"
    ),
  {
    resolveComponent: (components) => components.ApplicationDetailPage,
  }
);

// Fetch commands detail periodically
const FETCH_COMMANDS_INTERVAL = 3000;
const useCommandsStatusChecking = (): void => {
  const dispatch = useAppDispatch();
  const commandIds = useAppSelector<EntityId[]>((state) =>
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

export const Routes: FC = () => {
  const dispatch = useAppDispatch();
  const me = useAppSelector((state) => state.me);
  useEffect(() => {
    if (me?.isLogin) {
      dispatch(fetchPipeds(true));
    }
  }, [dispatch, me]);
  useCommandsStatusChecking();

  const location = useLocation();
  const [, onLoadProject] = useQueryString("project", "");
  useEffect(() => {
    if (me?.isLogin) {
      onLoadProject(me.projectId);

      // Add logged in users project to localstorage.
      const projects = localStorage.getItem(USER_PROJECTS)?.split(",") || [];
      projects.push(me.projectId);
      localStorage.setItem(USER_PROJECTS, sortedSet(projects).join(","));
    }
  }, [location, me, onLoadProject]);

  const [showWarningBanner, setShowWarningBaner] = useState(
    localStorage.getItem(BANNER_VERSION_KEY) !== `${process.env.PIPECD_VERSION}`
  );

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
          <Route exact path={PAGE_PATH_LOGIN}>
            <LoginPage />
          </Route>
          <Route
            path={PAGE_PATH_TOP}
            component={(props: RouteComponentProps) => {
              localStorage.setItem(
                REDIRECT_PATH_KEY,
                `${props.location.pathname}${props.location.search}`
              );
              return (
                <Redirect to={`${PAGE_PATH_LOGIN}${props.location.search}`} />
              );
            }}
          />
        </Switch>
      </>
    );
  }

  const handleCloseWarningBanner = (): void => {
    localStorage.setItem(BANNER_VERSION_KEY, `${process.env.PIPECD_VERSION}`);
    setShowWarningBaner(false);
  };

  return (
    <>
      {showWarningBanner && (
        <WarningBanner onClose={handleCloseWarningBanner} />
      )}
      <Header />
      <Switch>
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
        <Route
          exact
          path={PAGE_PATH_DEPLOYMENT_CHAINS}
          component={DeploymentChainsIndexPage}
        />
        <Route path={PAGE_PATH_SETTINGS} component={SettingsIndexPage} />
        <Route path={PAGE_PATH_INSIGHTS} component={InsightIndexPage} />
        <Route path={PAGE_PATH_EVENTS} component={EventIndexPage} />
        <Route
          path={PAGE_PATH_TOP}
          component={() => {
            const path =
              localStorage.getItem(REDIRECT_PATH_KEY) || PAGE_PATH_APPLICATIONS;
            localStorage.removeItem(REDIRECT_PATH_KEY);
            return <Redirect to={path} />;
          }}
        />
      </Switch>
      <Toasts />
    </>
  );
};
