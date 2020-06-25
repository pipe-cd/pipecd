import { Button, Divider, Drawer, Toolbar } from "@material-ui/core";
import { Add } from "@material-ui/icons";
import React, { FC, memo, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AddApplicationForm } from "../../components/add-application-form";
import { ApplicationList } from "../../components/application-list";
import { AppState } from "../../modules";
import { addApplication, fetchApplications } from "../../modules/applications";
import { AppDispatch } from "../../store";

export const ApplicationIndexPage: FC = memo(function ApplicationIndexPage() {
  const dispatch = useDispatch<AppDispatch>();
  const [isOpenForm, setIsOpenForm] = useState(false);
  const isAdding = useSelector<AppState, boolean>(
    (state) => state.applications.adding
  );

  useEffect(() => {
    dispatch(fetchApplications());
  }, [dispatch]);

  const handleClose = (): void => {
    setIsOpenForm(false);
  };

  return (
    <div>
      <Toolbar variant="dense">
        <Button
          color="primary"
          startIcon={<Add />}
          onClick={() => setIsOpenForm(true)}
        >
          ADD
        </Button>
      </Toolbar>
      <Divider />

      <ApplicationList />

      <Drawer
        anchor="right"
        open={isOpenForm}
        onClose={handleClose}
        ModalProps={{ disableBackdropClick: isAdding }}
      >
        <AddApplicationForm
          projectName="pipe-cd"
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
    </div>
  );
});
