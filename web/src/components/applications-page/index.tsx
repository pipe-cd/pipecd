import {
  Box,
  Button,
  CircularProgress,
  Divider,
  Drawer,
  Toolbar,
} from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { Add } from "@mui/icons-material";
import CloseIcon from "@mui/icons-material/Close";
import FilterIcon from "@mui/icons-material/FilterList";
import RefreshIcon from "@mui/icons-material/Refresh";
import { FC, useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "~/constants/path";
import {
  UI_TEXT_ADD,
  UI_TEXT_FILTER,
  UI_TEXT_HIDE_FILTER,
  UI_TEXT_REFRESH,
} from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import {
  clearAddedApplicationId,
  fetchApplications,
} from "~/modules/applications";
import {
  arrayFormat,
  stringifySearchParams,
  useSearchParams,
} from "~/utils/search-params";
import AddApplicationDrawer from "./add-application-drawer";
import EditApplicationDrawer from "./edit-application-drawer";
import { ApplicationAddedView } from "./application-added-view";
import { ApplicationFilter } from "./application-filter";
import { ApplicationList } from "./application-list";

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

export const ApplicationIndexPage: FC = () => {
  const classes = useStyles();
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const filterOptions = useSearchParams();
  const [openAddForm, setOpenAddForm] = useState(false);
  const [openFilter, setOpenFilter] = useState(true);
  const isAdding = useAppSelector<boolean>(
    (state) => state.applications.adding
  );
  const isLoading = useAppSelector<boolean>(
    (state) => state.applications.loading
  );
  const addedApplicationId = useAppSelector<string | null>(
    (state) => state.applications.addedApplicationId
  );

  const currentPage =
    typeof filterOptions.page === "string"
      ? parseInt(filterOptions.page, 10)
      : 1;

  const updateURL = useCallback(
    (options: Record<string, string | number | boolean | undefined>) => {
      navigate(
        `${PAGE_PATH_APPLICATIONS}?${stringifySearchParams(
          { ...options },
          { arrayFormat: arrayFormat }
        )}`,
        { replace: true }
      );
    },
    [navigate]
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
  }, [dispatch, filterOptions]);

  const handleCloseApplicationAddedView = (): void => {
    dispatch(clearAddedApplicationId());
  };

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
          {UI_TEXT_ADD}
        </Button>
        <div className={classes.toolbarSpacer} />
        <Button
          color="primary"
          startIcon={<RefreshIcon />}
          onClick={fetchApplicationsWithOptions}
          disabled={isLoading}
        >
          {UI_TEXT_REFRESH}
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
        onAdded={() => {
          setOpenAddForm(false);
          fetchApplicationsWithOptions();
        }}
      />
      <EditApplicationDrawer onUpdated={fetchApplicationsWithOptions} />

      <Drawer
        anchor="right"
        open={!!addedApplicationId}
        onClose={(_, reason) => {
          if (reason === "backdropClick" && isAdding) return;
          handleCloseApplicationAddedView();
        }}
      >
        <ApplicationAddedView onClose={handleCloseApplicationAddedView} />
      </Drawer>
    </>
  );
};
