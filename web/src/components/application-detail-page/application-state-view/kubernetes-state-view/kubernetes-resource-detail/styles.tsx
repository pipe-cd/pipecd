import {
  IconButton,
  Paper,
  Typography,
  TypographyProps,
  styled,
} from "@mui/material";

const DETAIL_WIDTH = 400;
export const PanelWrap = styled(Paper)(() => ({
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

export const PanelTitle = styled((props: TypographyProps) => (
  <Typography variant="h6" {...props} />
))(({ theme }) => ({
  paddingRight: theme.spacing(4),
  wordBreak: "break-all",
  paddingBottom: theme.spacing(2),
}));

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
