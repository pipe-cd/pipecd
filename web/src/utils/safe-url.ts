/**
 * Returns the given path if it is a same-origin absolute path, otherwise the fallback.
 * Rejects protocol-relative paths ("//host") and backslash variants ("/\host")
 * that browsers and routers may resolve to another origin.
 */
export const getSafeRedirectPath = (
  path: string | null | undefined,
  fallback: string
): string => {
  if (!path || !path.startsWith("/") || path.includes("\\")) return fallback;
  if (path.startsWith("//")) return fallback;
  return path;
};

const SAFE_EXTERNAL_PROTOCOLS = ["http:", "https:"];

/**
 * Returns the given URL if it is an absolute http(s) URL, otherwise undefined.
 * Prevents rendering links with schemes such as "javascript:" or "data:".
 */
export const getSafeExternalUrl = (
  url: string | null | undefined
): string | undefined => {
  if (!url) return undefined;
  try {
    const parsed = new URL(url);
    return SAFE_EXTERNAL_PROTOCOLS.includes(parsed.protocol) ? url : undefined;
  } catch {
    return undefined;
  }
};
