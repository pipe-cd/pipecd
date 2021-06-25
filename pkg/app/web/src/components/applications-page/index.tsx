import {
  Box,
  Button,
  CircularProgress,
  Divider,
  Drawer,
  makeStyles,
  Toolbar,
} from "@material-ui/core";
import { Add } from "@material-ui/icons";
import CloseIcon from "@material-ui/icons/Close";
import FilterIcon from "@material-ui/icons/FilterList";
import RefreshIcon from "@material-ui/icons/Refresh";
import { FC, memo, useCallback, useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "~/constants/path";
import { UI_TEXT_FILTER, UI_TEXT_HIDE_FILTER } from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { fetchApplicationCount } from "~/modules/application-counts";
import { ApplicationKind, fetchApplications } from "~/modules/applications";
import { clearTemplateTarget } from "~/modules/deployment-configs";
import { stringifySearchParams, useSearchParams } from "~/utils/search-params";
import { AddApplicationDrawer } from "./add-application-drawer";
import { ApplicationCounts } from "./application-counts";
import { ApplicationFilter } from "./application-filter";
import { ApplicationList } from "./application-list";
import { DeploymentConfigForm } from "./deployment-config-form";
import { EditApplicationDrawer } from "./edit-application-drawer";

const useStyles = makeStyles((theme) => ({
  main: {
    display: "flex",
    overflowY: "hidden",
    overflowX: "auto",
    flex: 1,
  },
  toolbarSpacer: {
    flexGrow: 1,
  },
  buttonProgress: {
    color: theme.palette.primary.main,
    position: "absolute",
    top: "50%",
    left: "50%",
    marginTop: -12,
    marginLeft: -12,
  },
}));

export const ApplicationIndexPage: FC = memo(function ApplicationIndexPage() {
  const classes = useStyles();
  const dispatch = useAppDispatch();
  const history = useHistory();
  const filterOptions = useSearchParams();
  const [openAddForm, setOpenAddForm] = useState(false);
  const [openFilter, setOpenFilter] = useState(true);
  const [isLoading, isAdding] = useAppSelector<[boolean, boolean]>((state) => [
    state.applications.loading,
    state.applications.adding,
  ]);
  const addedApplicationId = useAppSelector<string | null>(
    (state) => state.deploymentConfigs.targetApplicationId
  );
  const currentPage =
    typeof filterOptions.page === "string"
      ? parseInt(filterOptions.page, 10)
      : 1;

  const updateURL = useCallback(
    (options: Record<string, string | number | boolean | undefined>) => {
      history.replace(
        `${PAGE_PATH_APPLICATIONS}?${stringifySearchParams({
          ...options,
        })}`
      );
    },
    [history]
  );

  const handleFilterChange = useCallback(
    (options) => {
      updateURL({ ...options, page: 1 });
    },
    [updateURL]
  );
  const handleFilterClear = useCallback(() => {
    updateURL({ page: currentPage });
  }, [updateURL, currentPage]);

  const fetchApplicationsWithOptions = useCallback(() => {
    dispatch(fetchApplications(filterOptions));
    dispatch(fetchApplicationCount());
  }, [dispatch, filterOptions]);

  const handleCloseTemplateForm = (): void => {
    dispatch(clearTemplateTarget());
  };

  const handleApplicationCountClick = useCallback(
    (kind: ApplicationKind) => {
      updateURL({ ...filterOptions, kind });
    },
    [updateURL, filterOptions]
  );

  const handlePageChange = useCallback(
    (page: number) => {
      updateURL({ ...filterOptions, page });
    },
    [updateURL, filterOptions]
  );

  useEffect(() => {
    fetchApplicationsWithOptions();
  }, [fetchApplicationsWithOptions]);

  return (
    <>
      <Toolbar variant="dense">
        <Button
          color="primary"
          startIcon={<Add />}
          onClick={() => setOpenAddForm(true)}
        >
          ADD
        </Button>
        <div className={classes.toolbarSpacer} />
        <Button
          color="primary"
          startIcon={<RefreshIcon />}
          onClick={fetchApplicationsWithOptions}
          disabled={isLoading}
        >
          {"REFRESH"}
          {isLoading && (
            <CircularProgress size={24} className={classes.buttonProgress} />
          )}
        </Button>
        <Button
          color="primary"
          startIcon={openFilter ? <CloseIcon /> : <FilterIcon />}
          onClick={() => setOpenFilter(!openFilter)}
        >
          {openFilter ? UI_TEXT_HIDE_FILTER : UI_TEXT_FILTER}
        </Button>
      </Toolbar>

      <Divider />

      <div className={classes.main}>
        <Box display="flex" flexDirection="column" flex={1} p={2}>
          <ApplicationCounts onClick={handleApplicationCountClick} />
          <ApplicationList
            currentPage={currentPage}
            onPageChange={handlePageChange}
            onRefresh={fetchApplicationsWithOptions}
          />
        </Box>
        {openFilter && (
          <ApplicationFilter
            options={filterOptions}
            onChange={handleFilterChange}
            onClear={handleFilterClear}
          />
        )}
      </div>

      <AddApplicationDrawer
        open={openAddForm}
        onClose={() => setOpenAddForm(false)}
        onAdded={fetchApplicationsWithOptions}
      />
      <EditApplicationDrawer onUpdated={fetchApplicationsWithOptions} />

      <Drawer
        anchor="right"
        open={!!addedApplicationId}
        onClose={handleCloseTemplateForm}
        ModalProps={{ disableBackdropClick: isAdding }}
      >
        {addedApplicationId && (
          <DeploymentConfigForm onSkip={handleCloseTemplateForm} />
        )}
      </Drawer>
    </>
  );
});
