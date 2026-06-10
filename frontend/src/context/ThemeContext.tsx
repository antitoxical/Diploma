import React, { createContext, useContext, useState, useEffect, useCallback, ReactNode } from "react";

type ThemeMode = "light" | "dark" | "system";

interface ThemeContextType {
  mode: ThemeMode;
  isDark: boolean;
  setMode: (m: ThemeMode) => void;
}

const ThemeContext = createContext<ThemeContextType>({
  mode: "system",
  isDark: false,
  setMode: () => {},
});

const LIGHT_VARS: Record<string, string> = {
  "--tg-theme-bg-color": "#f6f6f6",
  "--tg-theme-text-color": "#222222",
  "--tg-theme-hint-color": "#999999",
  "--tg-theme-link-color": "#a349a4",
  "--tg-theme-button-color": "#a349a4",
  "--tg-theme-button-text-color": "#ffffff",
  "--tg-theme-secondary-bg-color": "#ffffff",
  "--tg-theme-secondary-bg-text-color": "#222222",
};

const DARK_VARS: Record<string, string> = {
  "--tg-theme-bg-color": "#1c1c1e",
  "--tg-theme-text-color": "#f5f5f5",
  "--tg-theme-hint-color": "#aeaeb2",
  "--tg-theme-link-color": "#bf7ff0",
  "--tg-theme-button-color": "#8e44ad",
  "--tg-theme-button-text-color": "#ffffff",
  "--tg-theme-secondary-bg-color": "#2c2c2e",
  "--tg-theme-secondary-bg-text-color": "#f5f5f5",
};

/** Telegram keys we may overlay on top of the base palette in «Авто» mode */
const TG_PARAM_MAP: Record<string, string> = {
  bg_color: "--tg-theme-bg-color",
  text_color: "--tg-theme-text-color",
  hint_color: "--tg-theme-hint-color",
  link_color: "--tg-theme-link-color",
  button_color: "--tg-theme-button-color",
  button_text_color: "--tg-theme-button-text-color",
  secondary_bg_color: "--tg-theme-secondary-bg-color",
};

function getTgWebApp(): any | null {
  try {
    return (window as any).Telegram?.WebApp ?? null;
  } catch {
    return null;
  }
}

function getTgColorScheme(): "dark" | "light" | null {
  const scheme = getTgWebApp()?.colorScheme;
  return scheme === "dark" || scheme === "light" ? scheme : null;
}

function applyVars(vars: Record<string, string>) {
  const root = document.documentElement;
  Object.entries(vars).forEach(([k, v]) => root.style.setProperty(k, v));
}

function setDocumentTheme(dark: boolean) {
  const root = document.documentElement;
  root.setAttribute("data-theme", dark ? "dark" : "light");
  root.style.colorScheme = dark ? "dark" : "light";
}

/** Base palette first, then optional Telegram accents — never partial theme */
function applyThemeVars(dark: boolean, tg: any | null) {
  const base = dark ? DARK_VARS : LIGHT_VARS;
  applyVars({ ...base });

  const params = tg?.themeParams;
  if (!params) return;

  const root = document.documentElement;
  Object.entries(TG_PARAM_MAP).forEach(([tgKey, cssVar]) => {
    const value = params[tgKey];
    if (value) root.style.setProperty(cssVar, value);
  });

  // Text on secondary surfaces must follow main text, not a mismatched TG fallback
  const textOnSecondary =
    params.text_color ?? base["--tg-theme-text-color"];
  root.style.setProperty("--tg-theme-secondary-bg-text-color", textOnSecondary);
}

export const ThemeProvider = ({ children }: { children: ReactNode }) => {
  const [mode, setModeState] = useState<ThemeMode>(() => {
    const stored = localStorage.getItem("theme-mode") as ThemeMode;
    return stored === "light" || stored === "dark" || stored === "system" ? stored : "system";
  });

  const resolveIsDark = useCallback((m: ThemeMode): boolean => {
    if (m === "dark") return true;
    if (m === "light") return false;
    const tgScheme = getTgColorScheme();
    if (tgScheme) return tgScheme === "dark";
    return window.matchMedia?.("(prefers-color-scheme: dark)").matches || false;
  }, []);

  const [isDark, setIsDark] = useState(() => resolveIsDark(mode));

  const setMode = useCallback((m: ThemeMode) => {
    localStorage.setItem("theme-mode", m);
    setModeState(m);
  }, []);

  const applyTheme = useCallback(() => {
    const dark = resolveIsDark(mode);
    setIsDark(dark);
    setDocumentTheme(dark);

    const tg = getTgWebApp();

    if (mode === "system" && tg?.themeParams && Object.keys(tg.themeParams).length > 0) {
      applyThemeVars(dark, tg);
    } else {
      applyVars(dark ? DARK_VARS : LIGHT_VARS);
    }
  }, [mode, resolveIsDark]);

  useEffect(() => {
    applyTheme();

    const tg = getTgWebApp();
    const onTgThemeChanged = () => {
      if (mode === "system") applyTheme();
    };
    tg?.onEvent?.("themeChanged", onTgThemeChanged);

    const mq = window.matchMedia?.("(prefers-color-scheme: dark)");
    const onSystemSchemeChange = () => {
      if (mode === "system" && !getTgColorScheme()) applyTheme();
    };
    mq?.addEventListener?.("change", onSystemSchemeChange);

    return () => {
      tg?.offEvent?.("themeChanged", onTgThemeChanged);
      mq?.removeEventListener?.("change", onSystemSchemeChange);
    };
  }, [mode, applyTheme]);

  return (
    <ThemeContext.Provider value={{ mode, isDark, setMode }}>
      {children}
    </ThemeContext.Provider>
  );
};

export const useTheme = () => useContext(ThemeContext);
