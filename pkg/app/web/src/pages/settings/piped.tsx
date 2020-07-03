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
  List,
  ListItem,
  ListItemText,
  makeStyles,
  ListItemSecondaryAction,
  IconButton,
  Menu,
  MenuItem,
} from "@material-ui/core";
import { Add as AddIcon, MoreVert as MoreVertIcon } from "@material-ui/icons";
import React, { FC, memo, useState } from "react";
import { AddPipedForm } from "../../components/add-piped-form";
import { useDispatch, useSelector } from "react-redux";
import {
  addPiped,
  RegisteredPiped,
  clearRegisteredPipedInfo,
  Piped,
  selectAll,
  fetchPipeds,
  disablePiped,
} from "../../modules/pipeds";
import { AppState } from "../../modules";
import { AppDispatch } from "../../store";

const useStyles = makeStyles((theme) => ({
  item: {
    backgroundColor: theme.palette.background.paper,
  },
}));

const ITEM_HEIGHT = 48;

export const SettingsPipedPage: FC = memo(function SettingsPipedPage() {
  const classes = useStyles();
  const [isOpenForm, setIsOpenForm] = useState(false);
  const [actionTargetId, setActionTargetId] = useState<string | null>(null);
  const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
  const isOpenMenu = Boolean(anchorEl);
  const dispatch = useDispatch<AppDispatch>();
  const pipeds = useSelector<AppState, Piped[]>((state) =>
    selectAll(state.pipeds)
  );
  const registeredPiped = useSelector<AppState, RegisteredPiped | null>(
    (state) => state.pipeds.registeredPiped
  );

  const handleMenuOpen = (
    event: React.MouseEvent<HTMLButtonElement>,
    pipedId: string
  ): void => {
    setActionTargetId(pipedId);
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = (): void => {
    setAnchorEl(null);
    setActionTargetId(null);
  };

  const handleDisableClick = (): void => {
    if (actionTargetId) {
      dispatch(disablePiped({ pipedId: actionTargetId })).then(() => {
        dispatch(fetchPipeds(true));
      });
    }
    setAnchorEl(null);
    setActionTargetId(null);
  };

  const handleSubmit = (props: { name: string; desc: string }): void => {
    dispatch(addPiped(props)).then(() => {
      setIsOpenForm(false);
    });
  };

  const handleClose = (): void => {
    setIsOpenForm(false);
  };

  const handleClosePipedInfo = (): void => {
    dispatch(clearRegisteredPipedInfo());
    dispatch(fetchPipeds(true));
  };

  return (
    <div>
      <Toolbar variant="dense">
        <Button
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => setIsOpenForm(true)}
        >
          ADD
        </Button>
      </Toolbar>
      <Divider />

      <List disablePadding>
        {pipeds.map((pipe) => (
          <ListItem
            key={`pipe-${pipe.id}`}
            divider
            dense
            className={classes.item}
          >
            <ListItemText
              primary={pipe.id}
              secondary={`${pipe.name}: ${pipe.desc}`}
            />
            <ListItemSecondaryAction>
              <IconButton
                edge="end"
                aria-label="open menu"
                onClick={(e) => handleMenuOpen(e, pipe.id)}
              >
                <MoreVertIcon />
              </IconButton>
            </ListItemSecondaryAction>
          </ListItem>
        ))}
      </List>

      <Menu
        id="long-menu"
        anchorEl={anchorEl}
        keepMounted
        open={isOpenMenu}
        onClose={handleMenuClose}
        PaperProps={{
          style: {
            maxHeight: ITEM_HEIGHT * 4.5,
            width: "20ch",
          },
        }}
      >
        <MenuItem onClick={handleDisableClick}>Disable Piped</MenuItem>
      </Menu>

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
          />
          <TextField
            label="secret key"
            variant="outlined"
            value={registeredPiped?.key}
            fullWidth
            margin="dense"
          />
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
