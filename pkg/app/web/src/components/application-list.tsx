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
  TablePagination,
  TableFooter,
} from "@material-ui/core";
import MenuIcon from "@material-ui/icons/MoreVert";
import { Dictionary } from "@reduxjs/toolkit";
import dayjs from "dayjs";
import React, { FC, memo, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "../constants";
import { APPLICATION_SYNC_STATUS_TEXT } from "../constants/application-sync-status-text";
import { AppState } from "../modules";
import {
  Application,
  enableApplication,
  fetchApplications,
  selectAll,
} from "../modules/applications";
import {
  Environment,
  selectEntities as selectEnvs,
} from "../modules/environments";
import { AppDispatch } from "../store";
import { DisableApplicationDialog } from "./disable-application-dialog";
import { SyncStatusIcon } from "./sync-status-icon";

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
    flex: 1,
    overflow: "auto",
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

const PAGER_ROWS_PER_PAGE = [20, 50, { label: "All", value: -1 }];

export const ApplicationList: FC = memo(function ApplicationList() {
  const classes = useStyles();
  const dispatch = useDispatch<AppDispatch>();
  const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
  const isOpenMenu = Boolean(anchorEl);
  const [actionTarget, setActionTarget] = useState<Application | null>(null);
  const [openDisableDialog, setOpenDisableDialog] = useState(false);
  const [page, setPage] = React.useState(0);
  const [rowsPerPage, setRowsPerPage] = React.useState(20);

  const applications = useSelector<AppState, Application[]>((state) =>
    selectAll(state.applications)
  );
  const envs = useSelector<AppState, Dictionary<Environment>>((state) =>
    selectEnvs(state.environments)
  );

  const closeMenu = (): void => {
    setAnchorEl(null);
    setTimeout(() => {
      setActionTarget(null);
    }, 200);
  };

  // Menu item event handler
  const handleOnClickDisable = (): void => {
    setAnchorEl(null);
    setOpenDisableDialog(true);
  };

  const handleOnClickEnable = (): void => {
    if (actionTarget) {
      dispatch(enableApplication({ applicationId: actionTarget.id })).then(
        () => {
          dispatch(fetchApplications());
        }
      );
    }
    closeMenu();
  };

  const handleCloseDialog = (): void => {
    closeMenu();
    setOpenDisableDialog(false);
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
              <TableCell>Running Version</TableCell>
              <TableCell>Running Commit</TableCell>
              <TableCell>Deployed By</TableCell>
              <TableCell>Deployed At</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>
            {(rowsPerPage > 0
              ? applications.slice(
                  page * rowsPerPage,
                  page * rowsPerPage + rowsPerPage
                )
              : applications
            ).map((app) => {
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
                        {recentlyDeployment.trigger?.commander ||
                          recentlyDeployment.trigger?.commit?.author ||
                          NOT_AVAILABLE_TEXT}
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
                      onClick={(e) => {
                        setAnchorEl(e.currentTarget);
                        setActionTarget(app);
                      }}
                    >
                      <MenuIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
          <TableFooter>
            <TableRow>
              <TablePagination
                rowsPerPageOptions={PAGER_ROWS_PER_PAGE}
                count={applications.length}
                rowsPerPage={rowsPerPage}
                page={page}
                colSpan={8}
                onChangePage={(_, newPage) => {
                  setPage(newPage);
                }}
                onChangeRowsPerPage={(e) => {
                  setRowsPerPage(parseInt(e.target.value, 10));
                  setPage(0);
                }}
              />
            </TableRow>
          </TableFooter>
        </Table>
      </TableContainer>

      <Menu
        id="application-menu"
        anchorEl={anchorEl}
        keepMounted
        open={isOpenMenu}
        onClose={closeMenu}
        PaperProps={{
          style: {
            width: "20ch",
          },
        }}
      >
        {actionTarget && actionTarget.disabled ? (
          <MenuItem onClick={handleOnClickEnable}>Enable</MenuItem>
        ) : (
          <MenuItem onClick={handleOnClickDisable}>Disable</MenuItem>
        )}
      </Menu>

      <DisableApplicationDialog
        open={openDisableDialog}
        applicationId={actionTarget && actionTarget.id}
        onDisable={handleCloseDialog}
        onCancel={handleCloseDialog}
      />
    </div>
  );
});
