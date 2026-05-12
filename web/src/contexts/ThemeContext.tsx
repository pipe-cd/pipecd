import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { ThemeProvider as MuiThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import { getTheme } from "../theme";

type ThemeMode = "light" | "dark";

interface ThemeContextType {
  mode: ThemeMode;
  toggleTheme: () => void;
  setTheme: (mode: ThemeMode) => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

const THEME_STORAGE_KEY = "pipecd-theme-mode";

/**
 * Detects the system's preferred theme mode
 */
const getSystemThemeMode = (): ThemeMode => {
  if (typeof window === "undefined") return "light";
  return window.matchMedia("(prefers-color-scheme: dark)").matches
    ? "dark"
    : "light";
};

interface ThemeProviderProps {
  children: ReactNode;
  defaultMode?: ThemeMode;
}

export const ThemeProvider: React.FC<ThemeProviderProps> = ({
  children,
  defaultMode,
}) => {
  const [mode, setMode] = useState<ThemeMode>(() => {
    // Try to load from localStorage
    if (typeof window !== "undefined") {
      const stored = localStorage.getItem(THEME_STORAGE_KEY) as ThemeMode | null;
      if (stored && (stored === "light" || stored === "dark")) {
        return stored;
      }
    }
    // Fall back to default or system preference
    return defaultMode || getSystemThemeMode();
  });

  // Persist theme mode to localStorage
  useEffect(() => {
    localStorage.setItem(THEME_STORAGE_KEY, mode);
    // Update HTML attribute for CSS-based styling if needed
    if (typeof window !== "undefined") {
      document.documentElement.setAttribute("data-theme", mode);
    }
  }, [mode]);

  // Listen to system theme changes
  useEffect(() => {
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    const handleChange = (e: MediaQueryListEvent) => {
      const stored = localStorage.getItem(THEME_STORAGE_KEY);
      // Only update if user hasn't explicitly set a preference
      if (!stored) {
        setMode(e.matches ? "dark" : "light");
      }
    };

    mediaQuery.addEventListener("change", handleChange);
    return () => mediaQuery.removeEventListener("change", handleChange);
  }, []);

  const toggleTheme = () => {
    setMode((prevMode) => (prevMode === "light" ? "dark" : "light"));
  };

  const currentTheme = getTheme(mode);

  const value: ThemeContextType = {
    mode,
    toggleTheme,
    setTheme: setMode,
  };

  return (
    <ThemeContext.Provider value={value}>
      <MuiThemeProvider theme={currentTheme}>
        <CssBaseline />
        {children}
      </MuiThemeProvider>
    </ThemeContext.Provider>
  );
};

/**
 * Hook to access theme context
 * @throws Error if used outside ThemeProvider
 */
export const useTheme = (): ThemeContextType => {
  const context = useContext(ThemeContext);
  if (context === undefined) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }
  return context;
};

export default ThemeProvider;