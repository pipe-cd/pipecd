import React, { memo, FC, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  fetchApplications,
  selectAll,
  Application,
  addApplication,
} from "../../modules/applications";
import { AppState } from "../../modules";
import { Link, Drawer, Toolbar, Button, Divider } from "@material-ui/core";
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "../../constants";
import { AddApplicationForm } from "../../components/add-application-form";
import { Add } from "@material-ui/icons";
import { AppDispatch } from "../../modules/index";

export const ApplicationIndexPage: FC = memo(function ApplicationIndexPage() {
  const dispatch = useDispatch<AppDispatch>();
  const [isOpenForm, setIsOpenForm] = useState(false);
  const [isAdding, applications] = useSelector<
    AppState,
    [boolean, Application[]]
  >((state) => [state.applications.adding, selectAll(state.applications)]);

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

      <ul>
        {applications.map((application) => (
          <li key={application.id}>
            <Link
              component={RouterLink}
              to={`${PAGE_PATH_APPLICATIONS}/${application.id}`}
            >
              {application.name}
            </Link>
          </li>
        ))}
      </ul>

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
