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
  Typography,
  Accordion,
  AccordionSummary,
  AccordionDetails,
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
  enablePiped,
} from "../../modules/pipeds";
import { AppState } from "../../modules";
import { AppDispatch } from "../../store";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import clsx from "clsx";

const useStyles = makeStyles((theme) => ({
  main: {
    height: "100%",
    overflow: "auto",
  },
  item: {
    backgroundColor: theme.palette.background.paper,
  },
  disabledPipedsAccordion: {
    padding: 0,
  },
  disabledItemsSummary: {
    borderBottom: "1px solid rgba(0, 0, 0, .125)",
  },
  pipedsList: {
    flex: 1,
  },
  disabledItemsSecondaryHeader: {
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(3),
  },
  disabledItem: {
    opacity: 0.6,
  },
}));

const ITEM_HEIGHT = 48;

const usePipeds = (): [Piped[], Piped[]] => {
  const pipeds = useSelector<AppState, Piped[]>((state) =>
    selectAll(state.pipeds)
  );

  const disabled: Piped[] = [];
  const enabled: Piped[] = [];

  pipeds.forEach((piped) => {
    if (piped.disabled) {
      disabled.push(piped);
    } else {
      enabled.push(piped);
    }
  });

  return [enabled, disabled];
};

export const SettingsPipedPage: FC = memo(function SettingsPipedPage() {
  const classes = useStyles();
  const [isOpenForm, setIsOpenForm] = useState(false);
  const [actionTarget, setActionTarget] = useState<Piped | null>(null);
  const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
  const isOpenMenu = Boolean(anchorEl);
  const dispatch = useDispatch<AppDispatch>();
  const [enabledPipeds, disabledPipeds] = usePipeds();

  const registeredPiped = useSelector<AppState, RegisteredPiped | null>(
    (state) => state.pipeds.registeredPiped
  );

  const handleMenuOpen = (
    event: React.MouseEvent<HTMLButtonElement>,
    piped: Piped
  ): void => {
    setActionTarget(piped);
    setAnchorEl(event.currentTarget);
  };

  const closeMenu = (): void => {
    setAnchorEl(null);
    setTimeout(() => {
      setActionTarget(null);
    }, 200);
  };

  const handleDisableClick = (): void => {
    closeMenu();
    if (!actionTarget) {
      return;
    }

    const act = actionTarget.disabled ? enablePiped : disablePiped;

    dispatch(act({ pipedId: actionTarget.id })).then(() => {
      dispatch(fetchPipeds(true));
    });
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
    <>
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

      <div className={classes.main}>
        <List disablePadding className={classes.pipedsList}>
          {enabledPipeds.map((piped) => (
            <ListItem
              key={`pipe-${piped.id}`}
              divider
              dense
              className={classes.item}
            >
              <ListItemText
                primary={piped.id}
                secondary={`${piped.name}: ${piped.desc}`}
              />
              <ListItemSecondaryAction>
                <IconButton
                  edge="end"
                  aria-label="open menu"
                  onClick={(e) => handleMenuOpen(e, piped)}
                >
                  <MoreVertIcon />
                </IconButton>
              </ListItemSecondaryAction>
            </ListItem>
          ))}
        </List>

        <Accordion>
          <AccordionSummary
            expandIcon={<ExpandMoreIcon />}
            className={classes.disabledItemsSummary}
          >
            <Typography>Disabled pipeds</Typography>
            <Typography
              className={classes.disabledItemsSecondaryHeader}
            >{`Items: ${disabledPipeds.length}`}</Typography>
          </AccordionSummary>
          <AccordionDetails className={classes.disabledPipedsAccordion}>
            <List disablePadding className={classes.pipedsList}>
              {disabledPipeds.map((piped) => (
                <ListItem
                  key={`pipe-${piped.id}`}
                  divider
                  dense
                  className={clsx(classes.item, classes.disabledItem)}
                >
                  <ListItemText
                    primary={piped.id}
                    secondary={`${piped.name}: ${piped.desc}`}
                  />
                  <ListItemSecondaryAction>
                    <IconButton
                      edge="end"
                      aria-label="open menu"
                      onClick={(e) => handleMenuOpen(e, piped)}
                    >
                      <MoreVertIcon />
                    </IconButton>
                  </ListItemSecondaryAction>
                </ListItem>
              ))}
            </List>
          </AccordionDetails>
        </Accordion>
      </div>

      <Menu
        id="piped-menu"
        anchorEl={anchorEl}
        keepMounted
        open={isOpenMenu}
        onClose={() => closeMenu()}
        PaperProps={{
          style: {
            maxHeight: ITEM_HEIGHT * 4.5,
            width: "20ch",
          },
        }}
      >
        {actionTarget && actionTarget.disabled ? (
          <MenuItem onClick={handleDisableClick}>Enable</MenuItem>
        ) : (
          <MenuItem onClick={handleDisableClick}>Disable</MenuItem>
        )}
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
    </>
  );
});
