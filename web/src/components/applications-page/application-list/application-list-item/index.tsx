import {
  Box,
  IconButton,
  Link,
  makeStyles,
  Menu,
  MenuItem,
  TableCell,
  TableRow,
} from "@material-ui/core";
import MenuIcon from "@material-ui/icons/MoreVert";
import OpenInNewIcon from "@material-ui/icons/OpenInNew";
import clsx from "clsx";
import dayjs from "dayjs";
import { FC, memo, useState, Fragment } from "react";
import { Link as RouterLink } from "react-router-dom";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";
import { PAGE_PATH_APPLICATIONS } from "~/constants/path";
import { UI_TEXT_NOT_AVAILABLE_TEXT } from "~/constants/ui-text";
import { useAppSelector } from "~/hooks/redux";
import { Application, selectById } from "~/modules/applications";
import { AppSyncStatus } from "~/components/app-sync-status";

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
    flex: 1,
    overflow: "auto",
  },
  disabled: {
    background: theme.palette.grey[200],
  },
  labels: {
    maxHeight: 200,
    overflowY: "scroll",
    "&::-webkit-scrollbar": {
      display: "none",
    },
  },
  version: {
    maxWidth: 300,
    maxHeight: 200,
    wordBreak: "break-word",
    overflowY: "scroll",
    "&::-webkit-scrollbar": {
      display: "none",
    },
  },
  deployedBy: {
    maxWidth: 150,
    overflow: "hidden",
    textOverflow: "ellipsis",
    whiteSpace: "nowrap",
  },
  linkIcon: {
    fontSize: 16,
    verticalAlign: "text-bottom",
    marginLeft: theme.spacing(0.5),
  },
  warning: {
    color: "red",
  },
}));

const EmptyDeploymentData: FC<{ displayAllProperties: boolean }> = ({
  displayAllProperties,
}) =>
  displayAllProperties ? (
    <>
      <TableCell>{UI_TEXT_NOT_AVAILABLE_TEXT}</TableCell>
      <TableCell>{UI_TEXT_NOT_AVAILABLE_TEXT}</TableCell>
      <TableCell>{UI_TEXT_NOT_AVAILABLE_TEXT}</TableCell>
      <TableCell>{UI_TEXT_NOT_AVAILABLE_TEXT}</TableCell>
    </>
  ) : (
    <>
      <TableCell>{UI_TEXT_NOT_AVAILABLE_TEXT}</TableCell>
      <TableCell>{UI_TEXT_NOT_AVAILABLE_TEXT}</TableCell>
    </>
  );

export interface ApplicationListItemProps {
  applicationId: string;
  displayAllProperties?: boolean;
  onEdit: (id: string) => void;
  onEnable: (id: string) => void;
  onDisable: (id: string) => void;
  onDelete: (id: string) => void;
  onEncryptSecret: (id: string) => void;
}

export const ApplicationListItem: FC<ApplicationListItemProps> = memo(
  function ApplicationListItem({
    applicationId,
    displayAllProperties = true,
    onDisable,
    onEdit,
    onEnable,
    onDelete,
    onEncryptSecret,
  }) {
    const classes = useStyles();
    const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
    const app = useAppSelector<Application.AsObject | undefined>((state) =>
      selectById(state.applications, applicationId)
    );

    const handleEdit = (): void => {
      setAnchorEl(null);
      onEdit(applicationId);
    };

    const handleDisable = (): void => {
      setAnchorEl(null);
      onDisable(applicationId);
    };

    const handleEnable = (): void => {
      setAnchorEl(null);
      onEnable(applicationId);
    };

    const handleDelete = (): void => {
      setAnchorEl(null);
      onDelete(applicationId);
    };

    const handleGenerateSecret = (): void => {
      setAnchorEl(null);
      onEncryptSecret(applicationId);
    };

    if (!app) {
      return null;
    }

    const recentlyDeployment = app.mostRecentlySuccessfulDeployment;

    return (
      <>
        <TableRow className={clsx({ [classes.disabled]: app.disabled })}>
          <TableCell>
            <Box display="flex" alignItems="center">
              <AppSyncStatus
                syncState={app.syncState}
                deploying={app.deploying}
              />
            </Box>
          </TableCell>
          <TableCell>
            <Link
              component={RouterLink}
              to={`${PAGE_PATH_APPLICATIONS}/${app.id}`}
            >
              {app.name}
            </Link>
          </TableCell>
          <TableCell>{APPLICATION_KIND_TEXT[app.kind]}</TableCell>
          <TableCell>
            <div className={classes.labels}>
              {app.labelsMap.length !== 0
                ? app.labelsMap.map(([key, value]) => (
                    <Fragment key={key}>
                      <span>{key + ": " + value}</span>
                      <br />
                    </Fragment>
                  ))
                : "-"}
            </div>
          </TableCell>
          {recentlyDeployment ? (
            <>
              <TableCell>
                <div className={classes.version}>
                  {recentlyDeployment.versionsList.length !== 0 ? (
                    recentlyDeployment.versionsList.map((v) =>
                      v.name === "" ? (
                        <Fragment key={v.version}>
                          <span>{v.version}</span>
                          <br />
                        </Fragment>
                      ) : (
                        <Fragment key={`${v.name}:${v.version}`}>
                          <Link
                            href={v.url.includes("://") ? v.url : `//${v.url}`}
                            target="_blank"
                            rel="noreferrer"
                          >
                            {/* Trim first 7 characters like commit hash in case it is too long */}
                            {v.name}:
                            {v.version.length > 7
                              ? `${v.version.slice(0, 7)}...`
                              : v.version}
                            <OpenInNewIcon className={classes.linkIcon} />
                          </Link>
                          <br />
                        </Fragment>
                      )
                    )
                  ) : recentlyDeployment.version.includes(",") ? (
                    recentlyDeployment.version
                      .split(",")
                      .filter((item, index, arr) => arr.indexOf(item) === index)
                      .map((v) => (
                        <>
                          <span>{v}</span>
                          <br />
                        </>
                      ))
                  ) : (
                    <span>{recentlyDeployment.version}</span>
                  )}
                </div>
              </TableCell>
              {displayAllProperties && (
                <TableCell>
                  {recentlyDeployment.trigger?.commit && (
                    <Link
                      href={recentlyDeployment.trigger.commit.url}
                      target="_blank"
                      rel="noreferrer"
                    >
                      {recentlyDeployment.trigger.commit.hash.slice(0, 8)}
                      <OpenInNewIcon className={classes.linkIcon} />
                    </Link>
                  )}
                </TableCell>
              )}
              {displayAllProperties && (
                <TableCell className={classes.deployedBy}>
                  {recentlyDeployment.trigger?.commander ||
                    recentlyDeployment.trigger?.commit?.author ||
                    UI_TEXT_NOT_AVAILABLE_TEXT}
                </TableCell>
              )}
              <TableCell>{dayjs(app.updatedAt * 1000).fromNow()}</TableCell>
            </>
          ) : (
            <EmptyDeploymentData displayAllProperties={displayAllProperties} />
          )}
          <TableCell align="right">
            <IconButton
              aria-label="Open menu"
              onClick={(e) => {
                setAnchorEl(e.currentTarget);
              }}
            >
              <MenuIcon />
            </IconButton>
          </TableCell>
        </TableRow>

        <Menu
          id="application-menu"
          anchorEl={anchorEl}
          open={Boolean(anchorEl)}
          onClose={() => setAnchorEl(null)}
          PaperProps={{
            style: {
              width: "20ch",
            },
          }}
        >
          {app && app.disabled ? (
            <MenuItem onClick={handleEnable}>Enable</MenuItem>
          ) : (
            <div>
              <MenuItem onClick={handleEdit}>Edit</MenuItem>
              <MenuItem onClick={handleGenerateSecret}>Encrypt Secret</MenuItem>
              <MenuItem onClick={handleDisable}>Disable</MenuItem>
            </div>
          )}
          <MenuItem className={classes.warning} onClick={handleDelete}>
            Delete
          </MenuItem>
        </Menu>
      </>
    );
  }
);
