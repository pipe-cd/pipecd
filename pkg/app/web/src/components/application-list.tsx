import {
  IconButton,
  Link,
  makeStyles,
  Menu,
  MenuItem,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from "@material-ui/core";
import MenuIcon from "@material-ui/icons/MoreVert";
import { Dictionary } from "@reduxjs/toolkit";
import dayjs from "dayjs";
import React, { FC, memo, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "../constants";
import { APPLICATION_SYNC_STATUS_TEXT } from "../constants/application-sync-status-text";
import { AppState } from "../modules";
import {
  Application,
  selectAll,
  fetchApplications,
} from "../modules/applications";
import {
  Environment,
  selectEntities as selectEnvs,
} from "../modules/environments";
import { DisableApplicationDialog } from "./disable-application-dialog";
import { SyncStatusIcon } from "./sync-status-icon";

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
  },
  statusCell: {
    display: "flex",
    alignItems: "center",
  },
  statusText: {
    marginLeft: theme.spacing(1),
  },
}));

const NOT_AVAILABLE_TEXT = "N/A";

const EmptyDeploymentData: FC = () => (
  <>
    <TableCell>{NOT_AVAILABLE_TEXT}</TableCell>
    <TableCell>{NOT_AVAILABLE_TEXT}</TableCell>
    <TableCell>{NOT_AVAILABLE_TEXT}</TableCell>
    <TableCell>{NOT_AVAILABLE_TEXT}</TableCell>
  </>
);

export const ApplicationList: FC = memo(function ApplicationList() {
  const classes = useStyles();
  const dispatch = useDispatch();
  const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
  const isOpenMenu = Boolean(anchorEl);
  const [disableTarget, setDisableTarget] = useState<null | string>(null);
  const applications = useSelector<AppState, Application[]>((state) =>
    selectAll(state.applications)
  );
  const envs = useSelector<AppState, Dictionary<Environment>>((state) =>
    selectEnvs(state.environments)
  );
  const handleOnClickDisable = (): void => {
    if (anchorEl?.dataset.id) {
      setDisableTarget(anchorEl.dataset.id);
    }
    setAnchorEl(null);
  };
  const handleOnCancelDisable = (): void => {
    setDisableTarget(null);
    dispatch(fetchApplications());
  };
  const handleOnDisable = (): void => {
    setAnchorEl(null);
    setDisableTarget(null);
    dispatch(fetchApplications());
  };

  return (
    <div className={classes.root}>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Status</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>Environment</TableCell>
              <TableCell>Version</TableCell>
              <TableCell>Commit</TableCell>
              <TableCell>Trigger</TableCell>
              <TableCell>Last Deployment</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>
            {applications.map((app) => {
              const recentlyDeployment = app.mostRecentlySuccessfulDeployment;
              return (
                <TableRow key={`app-${app.id}`}>
                  <TableCell>
                    <div className={classes.statusCell}>
                      {app.syncState ? (
                        <>
                          <SyncStatusIcon status={app.syncState.status} />
                          <Typography className={classes.statusText}>
                            {APPLICATION_SYNC_STATUS_TEXT[app.syncState.status]}
                          </Typography>
                        </>
                      ) : (
                        NOT_AVAILABLE_TEXT
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <Link
                      component={RouterLink}
                      to={`${PAGE_PATH_APPLICATIONS}/${app.id}`}
                    >
                      {app.name}
                    </Link>
                  </TableCell>
                  <TableCell>{envs[app.envId]?.name}</TableCell>
                  {recentlyDeployment ? (
                    <>
                      <TableCell>{recentlyDeployment.version}</TableCell>
                      <TableCell>
                        {recentlyDeployment.trigger?.commit?.hash.slice(0, 8) ??
                          NOT_AVAILABLE_TEXT}
                      </TableCell>
                      <TableCell>
                        {recentlyDeployment.trigger?.commander}
                      </TableCell>
                      <TableCell>
                        {dayjs(recentlyDeployment.startedAt * 1000).fromNow()}
                      </TableCell>
                    </>
                  ) : (
                    <EmptyDeploymentData />
                  )}
                  <TableCell align="right">
                    <IconButton
                      data-id={app.id}
                      onClick={(e) => setAnchorEl(e.currentTarget)}
                    >
                      <MenuIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>

      <Menu
        id="application-menu"
        anchorEl={anchorEl}
        keepMounted
        open={isOpenMenu}
        onClose={() => {
          setAnchorEl(null);
        }}
        PaperProps={{
          style: {
            width: "20ch",
          },
        }}
      >
        <MenuItem onClick={handleOnClickDisable}>Disable</MenuItem>
      </Menu>

      <DisableApplicationDialog
        applicationId={disableTarget}
        onDisable={handleOnDisable}
        onCancel={handleOnCancelDisable}
      />
    </div>
  );
});
