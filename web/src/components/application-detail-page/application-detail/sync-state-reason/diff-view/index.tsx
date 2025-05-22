import { FC, memo } from "react";
import { Box } from "@mui/material";
import { AddedLine, ChangedLine, DeletedLine, LineWrap } from "./styles";

export interface DiffViewProps {
  content: string;
}

export const DiffView: FC<DiffViewProps> = memo(function DiffView({ content }) {
  return (
    <Box>
      {content.split("\n").map((line, i) => {
        switch (line[0]) {
          case "+":
            return (
              <LineWrap key={i} data-testid="added-line">
                <AddedLine>{line}</AddedLine>
              </LineWrap>
            );
          case "-":
            return (
              <LineWrap key={i} data-testid="deleted-line">
                <DeletedLine>{line}</DeletedLine>
              </LineWrap>
            );
          case "~":
            return (
              <LineWrap key={i} data-testid="changed-line">
                <ChangedLine>{line}</ChangedLine>
              </LineWrap>
            );
          default:
            return <LineWrap key={i}>{line}</LineWrap>;
        }
      })}
    </Box>
  );
});
