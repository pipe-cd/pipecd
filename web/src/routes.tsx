import loadable from "@loadable/component";
import { EntityId } from "@reduxjs/toolkit";
import { FC, useEffect, useState } from "react";
import {
  Route,
  useLocation,
  Routes as ReactRoutes,
  Navigate,
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
  PAGE_PATH_DEPLOYMENT_TRACE,
  PAGE_PATH_INSIGHTS,
  PAGE_PATH_EVENTS,
  PAGE_PATH_LOGIN,
  PAGE_PATH_SETTINGS,
  PAGE_PATH_TOP,
  PAGE_PATH_SETTINGS_PIPED,
  PAGE_PATH_SETTINGS_PROJECT,
  PAGE_PATH_SETTINGS_API_KEY,
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
import DeploymentTracePage from "./components/deployment-trace-page";

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

const SettingsPipedPage = loadable(
  () =>
    import(
      /* webpackChunkName: "settings-piped" */ "~/components/settings-page/piped"
    ),
  {
    resolveComponent: (components) => components.SettingsPipedPage,
  }
);

const SettingsProjectPage = loadable(
  () =>
    import(
      /* webpackChunkName: "settings-project" */ "~/components/settings-page/project"
    ),
  {
    resolveComponent: (components) => components.SettingsProjectPage,
  }
);

const APIKeyPage = loadable(
  () =>
    import(
      /* webpackChunkName: "settings-api-key" */ "~/components/settings-page/api-key"
    ),
  {
    resolveComponent: (components) => components.APIKeyPage,
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
        <ReactRoutes>
          <Route path={PAGE_PATH_LOGIN} element={<LoginPage />} />
          <Route
            path={PAGE_PATH_TOP}
            Component={() => {
              localStorage.setItem(
                REDIRECT_PATH_KEY,
                `${location.pathname}${location.search}`
              );
              return (
                <Navigate to={`${PAGE_PATH_LOGIN}${location.search}`} replace />
              );
            }}
          />
        </ReactRoutes>
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
      <ReactRoutes>
        <Route
          path={PAGE_PATH_APPLICATIONS}
          element={<ApplicationIndexPage />}
        />
        <Route
          path={`${PAGE_PATH_APPLICATIONS}/:applicationId`}
          element={<ApplicationDetailPage />}
        />
        <Route path={PAGE_PATH_DEPLOYMENTS} element={<DeploymentIndexPage />} />
        <Route
          path={PAGE_PATH_DEPLOYMENT_TRACE}
          element={<DeploymentTracePage />}
        />
        <Route
          path={`${PAGE_PATH_DEPLOYMENTS}/:deploymentId`}
          element={<DeploymentDetailPage />}
        />
        <Route
          path={PAGE_PATH_DEPLOYMENT_CHAINS}
          element={<DeploymentChainsIndexPage />}
        />
        <Route path={PAGE_PATH_SETTINGS} element={<SettingsIndexPage />}>
          <Route
            path={PAGE_PATH_SETTINGS}
            element={<Navigate to={PAGE_PATH_SETTINGS_PIPED} replace />}
          />
          <Route
            path={PAGE_PATH_SETTINGS_PIPED}
            element={<SettingsPipedPage />}
          />
          <Route
            path={PAGE_PATH_SETTINGS_PROJECT}
            element={<SettingsProjectPage />}
          />
          <Route path={PAGE_PATH_SETTINGS_API_KEY} element={<APIKeyPage />} />
        </Route>
        <Route path={PAGE_PATH_INSIGHTS} element={<InsightIndexPage />} />
        <Route path={PAGE_PATH_EVENTS} element={<EventIndexPage />} />
        <Route
          path={PAGE_PATH_TOP}
          Component={() => {
            const path =
              localStorage.getItem(REDIRECT_PATH_KEY) || PAGE_PATH_APPLICATIONS;
            localStorage.removeItem(REDIRECT_PATH_KEY);
            return <Navigate to={path} replace />;
          }}
        />
      </ReactRoutes>
      <Toasts />
    </>
  );
};
