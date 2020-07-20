import {
  Button,
  Divider,
  Drawer,
  makeStyles,
  Toolbar,
} from "@material-ui/core";
import { Add } from "@material-ui/icons";
import CloseIcon from "@material-ui/icons/Close";
import FilterIcon from "@material-ui/icons/FilterList";
import React, { FC, memo, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AddApplicationForm } from "../../components/add-application-form";
import { ApplicationFilter } from "../../components/application-filter";
import { ApplicationList } from "../../components/application-list";
import { AppState } from "../../modules";
import { addApplication, fetchApplications } from "../../modules/applications";
import { AppDispatch } from "../../store";

const useStyles = makeStyles(() => ({
  main: {
    display: "flex",
    overflow: "hidden",
    flex: 1,
  },
  toolbarSpacer: {
    flexGrow: 1,
  },
}));

export const ApplicationIndexPage: FC = memo(function ApplicationIndexPage() {
  const classes = useStyles();
  const dispatch = useDispatch<AppDispatch>();
  const [isOpenForm, setIsOpenForm] = useState(false);
  const [isOpenFilter, setIsOpenFilter] = useState(false);
  const isAdding = useSelector<AppState, boolean>(
    (state) => state.applications.adding
  );
  const projectName = useSelector<AppState, string>(
    (state) => state.project.name
  );

  const handleClose = (): void => {
    setIsOpenForm(false);
  };

  const handleChangeFilterOptions = (): void => {
    dispatch(fetchApplications());
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

      <Drawer
        anchor="right"
        open={isOpenForm}
        onClose={handleClose}
        ModalProps={{ disableBackdropClick: isAdding }}
      >
        <AddApplicationForm
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
      </Drawer>
    </>
  );
});
