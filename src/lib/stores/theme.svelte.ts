type ThemePreference = "light" | "dark" | "system";
type ResolvedTheme = "light" | "dark";

const STORAGE_KEY = "blocky-visor-theme";

function getStoredPreference(): ThemePreference {
  if (typeof window === "undefined") return "system";
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored === "light" || stored === "dark" || stored === "system")
    return stored;
  return "system";
}

function resolveTheme(pref: ThemePreference): ResolvedTheme {
  if (pref === "system") {
    if (typeof window === "undefined") return "dark";
    return window.matchMedia("(prefers-color-scheme: light)").matches
      ? "light"
      : "dark";
  }
  return pref;
}

function apply(t: ResolvedTheme) {
  if (typeof document === "undefined") return;
  document.documentElement.classList.toggle("dark", t === "dark");
}

function createThemeStore() {
  const initialPref = getStoredPreference();
  let preference = $state<ThemePreference>(initialPref);
  let resolved = $state<ResolvedTheme>(resolveTheme(initialPref));

  // Apply on init
  apply(resolveTheme(initialPref));

  // Listen for system theme changes when preference is 'system'
  if (typeof window !== "undefined") {
    window
      .matchMedia("(prefers-color-scheme: dark)")
      .addEventListener("change", () => {
        if (preference === "system") {
          resolved = resolveTheme("system");
          apply(resolved);
        }
      });
  }

  return {
    get preference() {
      return preference;
    },
    get isDark() {
      return resolved === "dark";
    },
    setPreference(pref: ThemePreference) {
      preference = pref;
      resolved = resolveTheme(pref);
      localStorage.setItem(STORAGE_KEY, pref);
      apply(resolved);
    },
  };
}

export const themeStore = createThemeStore();
