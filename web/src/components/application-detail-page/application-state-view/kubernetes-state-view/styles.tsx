import {
  IconButton,
  Paper,
  Typography,
  TypographyProps,
  styled,
} from "@mui/material";

export const StateViewRoot = styled("div")(() => ({
  display: "flex",
  flex: 1,
  justifyContent: "center",
  overflow: "hidden",
}));

export const StateViewWrapper = styled("div")(() => ({
  flex: 1,
  display: "flex",
  justifyContent: "center",
  overflow: "hidden",
}));

export const StateView = styled("div")(() => ({
  position: "relative",
  overflow: "auto",
}));

// resource detail panel
export const InfoRowTitle = styled((props: TypographyProps) => (
  <Typography variant="subtitle1" {...props} />
))(({ theme }) => ({
  color: theme.palette.text.secondary,
  minWidth: 120,
}));

export const InfoRowValue = styled((props: TypographyProps) => (
  <Typography variant="body1" {...props} />
))(() => ({
  flex: 1,
  wordBreak: "break-all",
}));

const DETAIL_WIDTH = 400;
export const ResourceDetailPanel = styled(Paper)(() => ({
  width: DETAIL_WIDTH,
  padding: "16px 24px",
  height: "100%",
  overflow: "auto",
  position: "relative",
  zIndex: 2,
}));

export const CloseButton = styled(IconButton)(({ theme }) => ({
  position: "absolute",
  right: theme.spacing(1),
  top: theme.spacing(1),
  color: theme.palette.grey[500],
}));

export const ResourceDetailTitle = styled((props: TypographyProps) => (
  <Typography variant="h6" {...props} />
))(({ theme }) => ({
  paddingRight: theme.spacing(4),
  wordBreak: "break-all",
  paddingBottom: theme.spacing(2),
}));
