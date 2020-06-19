import {
  Button,
  Divider,
  Drawer,
  Toolbar,
  Dialog,
  DialogTitle,
  DialogContent,
  TextField,
  Box,
} from "@material-ui/core";
import { Add } from "@material-ui/icons";
import React, { FC, memo, useState } from "react";
import { AddPipedForm } from "../../components/add-piped-form";
import { useDispatch, useSelector } from "react-redux";
import {
  addPiped,
  RegisteredPiped,
  clearRegisteredPipedInfo,
} from "../../modules/pipeds";
import { AppState } from "../../modules";

export const SettingsPipedPage: FC = memo(() => {
  const [isOpenForm, setIsOpenForm] = useState(false);
  const dispatch = useDispatch();
  const registeredPiped = useSelector<AppState, RegisteredPiped | null>(
    (state) => state.pipeds.registeredPiped
  );

  function handleSubmit(description: string) {
    dispatch(addPiped(description));
    setIsOpenForm(false);
  }

  function handleClose() {
    setIsOpenForm(false);
  }

  function handleClosePipedInfo() {
    dispatch(clearRegisteredPipedInfo());
  }

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

      <Drawer anchor="right" open={isOpenForm} onClose={handleClose}>
        <AddPipedForm
          projectName="pipe-cd"
          onSubmit={handleSubmit}
          onClose={handleClose}
        />
      </Drawer>

      <Dialog open={registeredPiped !== null}>
        <DialogTitle>Piped registered</DialogTitle>
        <DialogContent>
          <TextField
            label="id"
            variant="outlined"
            value={registeredPiped?.id}
            fullWidth
            margin="dense"
          ></TextField>
          <TextField
            label="secret key"
            variant="outlined"
            value={registeredPiped?.key}
            fullWidth
            margin="dense"
          ></TextField>
          <Box display="flex" justifyContent="flex-end" m={1} mt={2}>
            <Button color="primary" onClick={handleClosePipedInfo}>
              CLOSE
            </Button>
          </Box>
        </DialogContent>
      </Dialog>
    </div>
  );
});
