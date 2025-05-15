import {
  Box,
  IconButton,
  Link,
  Menu,
  MenuItem,
  TableCell,
  TableRow,
} from "@mui/material";
import MenuIcon from "@mui/icons-material/MoreVert";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";
import dayjs from "dayjs";
import { FC, memo, useState, Fragment } from "react";
import { Link as RouterLink } from "react-router-dom";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";
import { PAGE_PATH_APPLICATIONS } from "~/constants/path";
import { UI_TEXT_NOT_AVAILABLE_TEXT } from "~/constants/ui-text";
import { useAppSelector } from "~/hooks/redux";
import { Application, selectById } from "~/modules/applications";
import { AppSyncStatus } from "~/components/app-sync-status";

enum PipedVersion {
  V0 = "v0",
  V1 = "v1",
}

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

    const pipedVersion =
      !app.platformProvider || app?.deployTargetsByPluginMap?.length > 0
        ? PipedVersion.V1
        : PipedVersion.V0;

    return (
      <>
        <TableRow
          sx={(theme) => ({
            backgroundColor: app.disabled ? theme.palette.grey[200] : "inherit",
          })}
        >
          <TableCell>
            <Box
              sx={{
                display: "flex",
                alignItems: "center",
              }}
            >
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
          <TableCell>
            {pipedVersion === PipedVersion.V0 &&
              APPLICATION_KIND_TEXT[app.kind]}
            {pipedVersion === PipedVersion.V1 && "APPLICATION"}
          </TableCell>
          <TableCell>
            <Box
              sx={{
                maxHeight: 200,
                overflowY: "scroll",
                "&::-webkit-scrollbar": {
                  display: "none",
                },
              }}
            >
              {app.labelsMap.length !== 0
                ? app.labelsMap.map(([key, value]) => (
                    <Fragment key={key}>
                      <span>{key + ": " + value}</span>
                      <br />
                    </Fragment>
                  ))
                : "-"}
            </Box>
          </TableCell>
          {recentlyDeployment ? (
            <>
              <TableCell>
                <Box
                  sx={{
                    maxWidth: 300,
                    maxHeight: 200,
                    wordBreak: "break-word",
                    overflowY: "scroll",
                    "&::-webkit-scrollbar": {
                      display: "none",
                    },
                  }}
                >
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
                            <OpenInNewIcon
                              sx={{
                                fontSize: 16,
                                verticalAlign: "text-bottom",
                                marginLeft: 0.5,
                              }}
                            />
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
                </Box>
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
                      <OpenInNewIcon
                        sx={{
                          fontSize: 16,
                          verticalAlign: "text-bottom",
                          marginLeft: 0.5,
                        }}
                      />
                    </Link>
                  )}
                </TableCell>
              )}
              {displayAllProperties && (
                <TableCell
                  sx={{
                    maxWidth: 150,
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                    whiteSpace: "nowrap",
                  }}
                >
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
              size="large"
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
          <MenuItem
            sx={{
              color: "red",
            }}
            onClick={handleDelete}
          >
            Delete
          </MenuItem>
        </Menu>
      </>
    );
  }
);
