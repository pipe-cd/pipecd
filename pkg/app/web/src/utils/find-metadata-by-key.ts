export const findMetadataByKey = (
  metadata: [string, string][],
  key: string
): string | undefined => {
  const find = metadata.find(([k]) => k === key);

  if (!find) {
    return undefined;
  }

  return find[1];
};
