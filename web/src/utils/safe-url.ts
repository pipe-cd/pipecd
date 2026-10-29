// ASCII control characters (tab, newline, ...) are stripped by URL parsing,
// so a path such as "/\t/host" would resolve to "//host".
const hasControlCharacter = (value: string): boolean =>
  Array.from(value).some((char) => {
    const code = char.charCodeAt(0);
    return code <= 0x1f || code === 0x7f;
  });

/**
 * Returns the given path if it is a same-origin absolute path, otherwise the fallback.
 * Rejects protocol-relative paths ("//host"), backslash variants ("/\host") and
 * paths containing control characters, which browsers and routers may resolve
 * to another origin.
 */
export const getSafeRedirectPath = (
  path: string | null | undefined,
  fallback: string
): string => {
  if (!path || !path.startsWith("/") || path.includes("\\")) return fallback;
  if (path.startsWith("//")) return fallback;
  if (hasControlCharacter(path)) return fallback;
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
