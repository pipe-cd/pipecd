import { FC, memo, useEffect, useRef } from "react";
import { CircularProgress, Box } from "@mui/material";
import { LogLine } from "../log-line";
import { DEFAULT_BACKGROUND_COLOR } from "~/constants/term-colors";
import { LogBlock } from "~~/model/logblock_pb";
import useAuth from "~/contexts/auth-context/use-auth";

export interface LogProps {
  logs: LogBlock.AsObject[];
  loading: boolean;
}

export const Log: FC<LogProps> = memo(function Log({ logs, loading }) {
  const bottomRef = useRef<null | HTMLDivElement>(null);
  const { me } = useAuth();
  const username = me?.subject || "user";

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [logs]);

  return (
    <Box
      sx={(theme) => ({
        fontFamily: theme.typography.fontFamilyMono,
        backgroundColor: DEFAULT_BACKGROUND_COLOR,
        paddingTop: theme.spacing(2),
        paddingBottom: theme.spacing(2),
        height: "100%",
      })}
    >
      {/* Powerline style prompt */}
      <Box sx={{ ml: "5rem", mb: 2, userSelect: "none", display: "flex" }}>
        <Box
          sx={{
            display: "inline-flex",
            alignItems: "center",
            height: "28px",
            fontSize: "0.8rem",
            lineHeight: 1,
            overflow: "hidden",
            borderRadius: "14px",
          }}
        >
          {/* Segment 1: Username */}
          <Box
            sx={{
              backgroundColor: "#df5b18",
              color: "#fff",
              pl: "16px",
              pr: "26px",
              height: "100%",
              display: "flex",
              alignItems: "center",
              fontWeight: "bold",
              position: "relative",
              zIndex: 6,
              borderTopLeftRadius: "14px",
              borderBottomLeftRadius: "14px",
              clipPath: "polygon(0% 0%, calc(100% - 10px) 0%, 100% 50%, calc(100% - 10px) 100%, 0% 100%)",
            }}
          >
            {username}
          </Box>

          {/* Segment 2: Directory (~) */}
          <Box
            sx={{
              backgroundColor: "#dca018",
              color: "#fff",
              pl: "24px",
              pr: "24px",
              height: "100%",
              display: "flex",
              alignItems: "center",
              fontWeight: "bold",
              fontSize: "1.1rem",
              position: "relative",
              zIndex: 5,
              marginLeft: "-10px",
              clipPath: "polygon(0% 0%, 10px 50%, 0% 100%, calc(100% - 10px) 100%, 100% 50%, calc(100% - 10px) 0%)",
            }}
          >
            ~
          </Box>

          {/* Segment 3: Accent Arrow 1 (Green) */}
          <Box
            sx={{
              backgroundColor: "#50b577",
              width: "22px",
              height: "100%",
              position: "relative",
              zIndex: 4,
              marginLeft: "-10px",
              clipPath: "polygon(0% 0%, 10px 50%, 0% 100%, calc(100% - 10px) 100%, 100% 50%, calc(100% - 10px) 0%)",
            }}
          />

          {/* Segment 4: Accent Arrow 2 (Teal) */}
          <Box
            sx={{
              backgroundColor: "#3ea5a2",
              width: "22px",
              height: "100%",
              position: "relative",
              zIndex: 3,
              marginLeft: "-10px",
              clipPath: "polygon(0% 0%, 10px 50%, 0% 100%, calc(100% - 10px) 100%, 100% 50%, calc(100% - 10px) 0%)",
            }}
          />

          {/* Segment 5: Accent Arrow 3 (Dark Forest Green) */}
          <Box
            sx={{
              backgroundColor: "#1b594b",
              width: "22px",
              height: "100%",
              position: "relative",
              zIndex: 2,
              marginLeft: "-10px",
              clipPath: "polygon(0% 0%, 10px 50%, 0% 100%, calc(100% - 10px) 100%, 100% 50%, calc(100% - 10px) 0%)",
            }}
          />

          {/* Segment 6: stage.log */}
          <Box
            sx={{
              backgroundColor: "#141819",
              color: "#ecf0f1",
              pl: "24px",
              pr: "16px",
              height: "100%",
              display: "flex",
              alignItems: "center",
              fontWeight: "bold",
              position: "relative",
              zIndex: 1,
              marginLeft: "-10px",
              borderTopRightRadius: "14px",
              borderBottomRightRadius: "14px",
              clipPath: "polygon(0% 0%, 10px 50%, 0% 100%, 100% 100%, 100% 0%)",
            }}
          >
            stage.log
          </Box>
        </Box>
      </Box>

      {logs.map((log, i) => (
        <LogLine
          key={`log-${log.index}`}
          severity={log.severity}
          body={log.log}
          lineNumber={i + 1}
          createdAt={log.createdAt}
        />
      ))}
      {loading && (
        <Box
          sx={{
            display: "flex",
            justifyContent: "center",
            p: 1,
          }}
        >
          <CircularProgress color="secondary" />
        </Box>
      )}
      <Box
        sx={(theme) => ({
          height: theme.spacing(1),
          backgroundColor: DEFAULT_BACKGROUND_COLOR,
        })}
        ref={bottomRef}
      />
    </Box>
  );
});




