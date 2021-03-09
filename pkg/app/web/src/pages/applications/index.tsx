import {
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
import React, { FC, memo, useCallback, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { ApplicationFilter } from "../../components/application-filter";
import { AddApplicationDrawer } from "../../components/add-application-drawer";
import { ApplicationList } from "../../components/application-list";
import { DeploymentConfigForm } from "../../components/deployment-config-form";
import { EditApplicationDrawer } from "../../components/edit-application-drawer";
import { AppState } from "../../modules";
import { fetchApplications } from "../../modules/applications";
import { clearTemplateTarget } from "../../modules/deployment-configs";
import { AppDispatch } from "../../store";
import { UI_TEXT_FILTER, UI_TEXT_HIDE_FILTER } from "../../constants/ui-text";
import {
  stringifySearchParams,
  useSearchParams,
} from "../../utils/search-params";
import { useHistory } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "../../constants/path";

const useStyles = makeStyles((theme) => ({
  main: {
    display: "flex",
    overflow: "hidden",
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
  const dispatch = useDispatch<AppDispatch>();
  const history = useHistory();
  const filterOptions = useSearchParams();
  const [openAddForm, setOpenAddForm] = useState(false);
  const [openFilter, setOpenFilter] = useState(true);
  const [isLoading, isAdding] = useSelector<AppState, [boolean, boolean]>(
    (state) => [state.applications.loading, state.applications.adding]
  );

  const addedApplicationId = useSelector<AppState, string | null>(
    (state) => state.deploymentConfigs.targetApplicationId
  );

  const handleFilterChange = useCallback(
    (options) => {
      history.replace(
        `${PAGE_PATH_APPLICATIONS}?${stringifySearchParams(options)}`
      );
    },
    [history]
  );
  const handleFilterClear = useCallback(() => {
    history.replace(PAGE_PATH_APPLICATIONS);
  }, [history]);

  const handleRefresh = useCallback(() => {
    dispatch(fetchApplications(filterOptions));
  }, [dispatch, filterOptions]);

  const handleCloseTemplateForm = (): void => {
    dispatch(clearTemplateTarget());
  };

  useEffect(() => {
    dispatch(fetchApplications(filterOptions));
  }, [dispatch, filterOptions]);

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
          onClick={handleRefresh}
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
        <ApplicationList />
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
      />
      <EditApplicationDrawer />

      <Drawer
        anchor="right"
        open={!!addedApplicationId}
        onClose={handleCloseTemplateForm}
        ModalProps={{ disableBackdropClick: isAdding }}
      >
        {addedApplicationId && (
          <DeploymentConfigForm
            applicationId={addedApplicationId}
            onSkip={handleCloseTemplateForm}
          />
        )}
      </Drawer>
    </>
  );
});
