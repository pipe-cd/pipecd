import {
  Button,
  Divider,
  Drawer,
  makeStyles,
  Toolbar,
  CircularProgress,
} from "@material-ui/core";
import { Add } from "@material-ui/icons";
import CloseIcon from "@material-ui/icons/Close";
import FilterIcon from "@material-ui/icons/FilterList";
import RefreshIcon from "@material-ui/icons/Refresh";
import React, { FC, memo, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AddApplicationDrawer } from "../../components/add-application-drawer";
import { ApplicationFilter } from "../../components/application-filter";
import { ApplicationList } from "../../components/application-list";
import { AppState } from "../../modules";
import { addApplication, fetchApplications } from "../../modules/applications";
import { AppDispatch } from "../../store";
import { selectProjectName } from "../../modules/me";
import { DeploymentConfigForm } from "../../components/deployment-config-form";
import { clearTemplateTarget } from "../../modules/deployment-configs";

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
  const [isOpenForm, setIsOpenForm] = useState(false);
  const [isOpenFilter, setIsOpenFilter] = useState(false);
  const [isLoading, isAdding] = useSelector<AppState, [boolean, boolean]>(
    (state) => [state.applications.loading, state.applications.adding]
  );
  const projectName = useSelector<AppState, string>((state) =>
    selectProjectName(state.me)
  );

  const addedApplicationId = useSelector<AppState, string | null>(
    (state) => state.deploymentConfigs.targetApplicationId
  );

  const handleClose = (): void => {
    setIsOpenForm(false);
  };

  const handleChangeFilterOptions = (): void => {
    dispatch(fetchApplications());
  };

  const handleRefresh = (): void => {
    dispatch(fetchApplications());
  };

  const handleCloseTemplateForm = (): void => {
    dispatch(clearTemplateTarget());
  };

  useEffect(() => {
    dispatch(fetchApplications());
  }, [dispatch]);

  return (
    <>
      <Toolbar variant="dense">
        <Button
          color="primary"
          startIcon={<Add />}
          onClick={() => setIsOpenForm(true)}
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
          startIcon={isOpenFilter ? <CloseIcon /> : <FilterIcon />}
          onClick={() => setIsOpenFilter(!isOpenFilter)}
        >
          {isOpenFilter ? "HIDE FILTER" : "FILTER"}
        </Button>
      </Toolbar>

      <Divider />

      <div className={classes.main}>
        <ApplicationList />
        <ApplicationFilter
          open={isOpenFilter}
          onChange={handleChangeFilterOptions}
        />
      </div>

      <AddApplicationDrawer
        open={isOpenForm}
        projectName={projectName}
        onSubmit={(state) => {
          dispatch(addApplication(state)).then(() => {
            setIsOpenForm(false);
            dispatch(fetchApplications());
          });
        }}
        onClose={handleClose}
        isAdding={isAdding}
      />

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
